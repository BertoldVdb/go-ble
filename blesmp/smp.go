package blesmp

import (
	"context"
	"encoding/binary"
	"os"
	"path/filepath"
	"strconv"
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

type smpStoredLTKMapKey [19]byte

type smpStoredLTK struct {
	EDIV          uint16
	Rand          uint64
	LTK           [16]byte
	Authenticated bool
	Bonded        bool
}

func makeSMPStoredLTKMapKey(isCentral bool, ia bleutil.BLEAddr, ra bleutil.BLEAddr, EDIV uint16, Rand uint64) smpStoredLTKMapKey {
	var result smpStoredLTKMapKey

	if isCentral {
		result[0] = 1
	}

	ia.MacAddr.Encode(result[1:])
	result[7] = byte(ia.MacAddrType & 1)

	if !isCentral {
		if EDIV == 0 && Rand == 0 {
			/* LE Secure Connections: there is no EDIV/Rand to identify
			   the LTK, so the key has to identify the remote (central)
			   side. Earlier versions stored ia (= local peripheral
			   address) twice, which made every bonded peripheral LTK
			   keyed only on the peripheral's own address — different
			   centrals paired with the same peripheral would collide
			   and overwrite each other's LTKs. */
			result[8] = 1
			ra.MacAddr.Encode(result[9:])
			result[15] = byte(ra.MacAddrType & 1)
		} else {
			binary.LittleEndian.PutUint16(result[9:], EDIV)
			binary.LittleEndian.PutUint64(result[11:], Rand)
		}
	} else {
		ra.MacAddr.Encode(result[8:])
		result[14] = byte(ra.MacAddrType & 1)
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

	// MinKeySize sets the minimum LTK byte length the implementation will
	// accept during legacy pairing key-size negotiation. The Bluetooth
	// spec floor is 7 bytes (56 effective bits), which is the surface
	// exploited by KNOB-style downgrade attacks. Default 16 (128-bit)
	// when zero. Set explicitly to 7 only when interoperating with very
	// old peers that cannot negotiate higher.
	MinKeySize int
}

type SMPConfig struct {
	StoredKeysPath    string
	DefaultConnConfig *SMPConnConfig
}

func DefaultConfig() *SMPConfig {
	return &SMPConfig{
		StoredKeysPath: defaultStoredKeysPath(),
		DefaultConnConfig: &SMPConnConfig{
			AuthReq:        4,
			StaticPasscode: -1,
			MinKeySize:     16,
		},
	}
}

func defaultStoredKeysPath() string {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "go-ble", "smp.gob")
	}

	return filepath.Join(os.TempDir(), "go-ble-smp-"+strconv.Itoa(os.Getuid())+".gob")
}

type SMP struct {
	config *SMPConfig

	controller *hci.Controller

	storedKeys        map[smpStoredLTKMapKey]smpStoredLTK
	storedKeysPersist *gobpersist.GobPersist
}

// smpConnFromConn returns the SMPConn attached to a connmgr.Connection,
// or nil if the connection has not yet been wired through SMP. The
// connection's AppConn assignment can race with HCI events that arrive
// during link setup; an unchecked type assertion would crash the
// process. Callers must handle nil.
func smpConnFromConn(conn *hciconnmgr.Connection) *SMPConn {
	if conn == nil || conn.SMPConn == nil {
		return nil
	}
	c, _ := conn.SMPConn.(*SMPConn)
	return c
}

func (s *SMP) connmgrEncryptionChanged(conn *hciconnmgr.Connection, event *hcievents.EncryptionChangeEvent) *hcievents.EncryptionChangeEvent {
	smpConn := smpConnFromConn(conn)
	if smpConn == nil {
		return event
	}
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
	smpConn := smpConnFromConn(conn)
	if smpConn == nil {
		return event
	}

	select {
	case smpConn.encUpdateChan <- event.Status == 0:
	default:
	}

	return event
}

func New(logger *logrus.Entry, controller *hci.Controller, config *SMPConfig) *SMP {
	if config == nil {
		config = DefaultConfig()
	}
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

	if config.StoredKeysPath != "" {
		logger.WithError(os.MkdirAll(filepath.Dir(config.StoredKeysPath), 0o700)).Debug("Creating LTK database directory")
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

	timeout <-chan (time.Time)

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

	// testEncryptHook, if set, replaces leEncrypt for unit tests so
	// the SC state machine can be exercised without a live HCI link.
	// It is package-private and only used by *_test.go files.
	testEncryptHook func(ltk smpStoredLTK) error
}

func (c *SMPConn) updateTimeout(run bool) {
	if !run {
		c.timeout = nil
	} else {
		c.timeout = time.After(30 * time.Second)
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
		case <-c.timeout:
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
				opcode := byte(0)
				if pdu.Len() > 0 {
					opcode = pdu.Buf()[0]
				}
				c.logger.WithFields(logrus.Fields{
					"0opcode": opcode,
					"1len":    pdu.Len(),
				}).Debug("SMP received")
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
		secureState:   StateInsecure,
		secureChan:    make(chan (struct{}), 1),
		encUpdateChan: make(chan (bool), 1),
	}

	c.secureAuthReq = config.AuthReq

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
