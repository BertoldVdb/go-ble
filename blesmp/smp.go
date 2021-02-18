package blesmp

import (
	"context"
	"encoding/hex"
	"os"
	"path"
	"sync"
	"time"

	"github.com/BertoldVdb/go-ble/bleconnecter"
	"github.com/BertoldVdb/go-ble/hci"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	blel2cap "github.com/BertoldVdb/go-ble/l2cap"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/BertoldVdb/go-misc/gobpersist"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/BertoldVdb/go-misc/waitstate"
	"github.com/sirupsen/logrus"
)

type smpStoredLTKMapKey [6 + 6 + 2 + 1]byte

type smpStoredLTK struct {
	EDIV          uint16
	Rand          uint64
	LTK           [16]byte
	Authenticated bool
	Bonded        bool
}

func makeSMPStoredLTKMapKey(isCentral bool, ia bleutil.BLEAddr, ra bleutil.BLEAddr) smpStoredLTKMapKey {
	var result smpStoredLTKMapKey

	ra.MacAddr.Encode(result[:])
	ia.MacAddr.Encode(result[6:])
	result[12] = byte(ia.MacAddrType & 1)
	result[13] = byte(ra.MacAddrType & 1)
	if isCentral {
		result[14] = 1
	}

	return result
}

type SMPConnConfig struct {
	DisplayNumeric func(conn *SMPConn, number uint32) error
	InputYesNo     func(conn *SMPConn) (bool, error)
	InputNumeric   func(conn *SMPConn) (uint32, error)

	SecureOnConnect bool
	AuthReq         uint8
	StaticPasscode  int32
}

type SMPConfig struct {
	StoredKeysPath    string
	DefaultConnConfig *SMPConnConfig
}

func DefaultConfig() *SMPConfig {
	return &SMPConfig{
		StoredKeysPath: path.Join(os.TempDir(), "/smp.gob"),
		DefaultConnConfig: &SMPConnConfig{
			AuthReq:        4,
			StaticPasscode: -1,
		},
	}
}

type SMP struct {
	config *SMPConfig

	controller *hci.Controller

	storedKeys        map[smpStoredLTKMapKey]smpStoredLTK
	storedKeysPersist *gobpersist.GobPersist
}

func (s *SMP) connmgrEncryptionChanged(conn *hciconnmgr.Connection, event *hcievents.EncryptionChangeEvent) *hcievents.EncryptionChangeEvent {
	smpConn := conn.SMPConn.(*SMPConn)
	if event.Status != 0 {
		event.EncryptionEnabled = 0
	}

	select {
	case smpConn.encUpdateChan <- event.EncryptionEnabled > 0:
	default:
	}

	return event
}

func (s *SMP) connmgrEncryptionRefresh(conn *hciconnmgr.Connection, event *hcievents.EncryptionKeyRefreshCompleteEvent) *hcievents.EncryptionKeyRefreshCompleteEvent {
	smpConn := conn.SMPConn.(*SMPConn)

	select {
	case smpConn.encUpdateChan <- event.Status == 0:
	default:
	}

	return event
}

func New(logger *logrus.Entry, controller *hci.Controller, config *SMPConfig) *SMP {
	s := &SMP{
		controller: controller,
		config:     config,
		storedKeys: make(map[smpStoredLTKMapKey]smpStoredLTK),
	}

	controller.ConnMgr.SetEventsSMP(hciconnmgr.ConnectionMangerEventsSMP{
		LEEncryptionGetKey: s.connmgrLEGetKey,
		EncryptionChanged:  s.connmgrEncryptionChanged,
		EncryptionRefresh:  s.connmgrEncryptionRefresh,
	})

	s.storedKeysPersist = &gobpersist.GobPersist{
		Target:   &s.storedKeys,
		Filename: config.StoredKeysPath,
	}

	logger.WithError(s.storedKeysPersist.Load()).Info("Loading LTK database")
	logger.WithError(s.storedKeysPersist.Save()).Debug("Saving LTK database")

	return s
}

type SMPConn struct {
	parent *SMP

	config *SMPConnConfig
	conn   hciconnmgr.BufferConn
	logger *logrus.Entry

	pduRx chan (*pdu.PDU)

	protocol smpProtocol

	timeout *time.Timer

	addrLERemote bleutil.BLEAddr
	addrLELocal  bleutil.BLEAddr

	isCentral bool

	secureState     SMPState
	secureAuthReq   uint8
	secureStateWait waitstate.WaitState
	secureChan      chan (struct{})

	encUpdateChan chan (bool)

	keyMutex           sync.Mutex
	keyIsAuthenticated bool
	keyIsBonded        bool
}

func (c *SMPConn) updateTimeout(run bool) {
	if !run {
		c.timeout.Stop()
	} else {
		c.timeout.Reset(30 * time.Second)
	}
}

func (c *SMPConn) setState(state SMPState) {
	c.secureState = state
	c.secureStateWait.Set(state)
}

func (c *SMPConn) smpHandler() {
	defer func() {
		c.setState(StatePermanentlyFailed)

		c.conn.Close()
		/* Drain channel */
		for {
			pdu, ok := <-c.pduRx
			if !ok {
				break
			}
			bleutil.ReleaseBuffer(pdu)
		}
	}()

	for {
		select {
		case <-c.timeout.C:
			c.logger.Warn("Security manager timeout")
			c.setState(StatePermanentlyFailed)

		case value := <-c.encUpdateChan:
			if value {
				c.setState(StateSecure)
			}

		case <-c.secureChan:
			if c.secureState == StateInsecure || c.secureState == StateFailed {
				c.setState(StateBusy)

				if c.isCentral {
					c.sendPairingRequest()
				} else {
					c.sendSecurityRequest()
				}
			}

		case pdu, ok := <-c.pduRx:
			if !ok {
				return
			}

			if c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				c.logger.WithField("0pdu", hex.EncodeToString(pdu.Buf())).Debug("SMP Received")
			}

			keepBuffer := false
			if c.secureState != StatePermanentlyFailed {
				keepBuffer = c.handleMessage(pdu)
			}
			if !keepBuffer {
				bleutil.ReleaseBuffer(pdu)
			}
		}
	}
}

func (c *SMPConn) rawConnLE() *bleconnecter.BLEConnection {
	return c.conn.(*blel2cap.L2Connection).ParentConn().(*bleconnecter.BLEConnection)
}

func (p *SMP) AddConn(conn hciconnmgr.BufferConn, config *SMPConnConfig) *SMPConn {
	if config == nil {
		config = p.config.DefaultConnConfig
	}

	c := &SMPConn{
		config:        config,
		parent:        p,
		conn:          conn,
		logger:        bleutil.LogWithPrefix(conn.GetLogger(), "smp"),
		pduRx:         make(chan (*pdu.PDU), 1),
		timeout:       time.NewTimer(0),
		secureState:   StateInsecure,
		secureChan:    make(chan (struct{}), 1),
		encUpdateChan: make(chan (bool), 1),
	}

	c.secureAuthReq = config.AuthReq

	<-c.timeout.C

	raw := c.rawConnLE()
	c.addrLERemote = raw.RemoteAddr().(bleutil.BLEAddr)
	c.addrLELocal = raw.LocalAddr().(bleutil.BLEAddr)
	c.isCentral = raw.IsCentral()
	raw.Connection.SMPConn = c

	c.setState(StateInsecure)

	/* Try to use the key we already have */
	if c.isCentral {
		err := c.leTryEncryptLTK()
		if err != nil {
			conn.Close()
			return nil
		}
	}

	go c.reader()
	go c.smpHandler()

	if c.isCentral && config.SecureOnConnect {
		/* While this code works for non-central as well, it would rebond on each connection...
		   Therefore, this is disabled. */
		go c.GoSecure(context.Background(), true)
	}

	return c
}

func (c *SMPConn) GetSecurity() (bool, bool, bool) {
	_, v, _ := c.secureStateWait.Get(nil, nil)
	state := v.(SMPState)

	if state != StateSecure {
		return false, false, false
	}

	c.keyMutex.Lock()
	defer c.keyMutex.Unlock()

	return true, c.keyIsAuthenticated, c.keyIsBonded
}

func (c *SMPConn) GoSecure(ctx context.Context, allowStart bool) (SMPState, error) {
	first := true

	_, state, err := c.secureStateWait.Get(ctx, func(cnt uint64, value interface{}) bool {
		wasFirst := first
		first = false

		state := value.(SMPState)
		/* If we already are secure, this is enough */
		if state == StateSecure {
			return true
		}

		/* If we are already busy, just wait */
		if state == StateBusy {
			return false
		}

		/* If we are in StateFailed or StateInsecure return or retry */
		if !wasFirst || !allowStart {
			return true
		}

		select {
		case c.secureChan <- struct{}{}:
		default:
		}
		return false
	})

	return state.(SMPState), err
}
