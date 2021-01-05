package blesmp

import (
	"encoding/hex"
	"log"
	"time"

	"github.com/BertoldVdb/go-ble/bleconnecter"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	blel2cap "github.com/BertoldVdb/go-ble/l2cap"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/BertoldVdb/go-misc/gobpersist"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"
)

type smpStoredLTKMapKey [6 + 6 + 2]byte

type smpStoredLTK struct {
	EDIV uint16
	Rand uint64
	LTK  [16]byte
}

func makeSMPStoredLTKMapKey(ia bleutil.BLEAddr, ra bleutil.BLEAddr) smpStoredLTKMapKey {
	var result smpStoredLTKMapKey

	ra.MacAddr.Encode(result[:])
	ia.MacAddr.Encode(result[6:])
	result[12] = byte(ia.MacAddrType & 1)
	result[13] = byte(ra.MacAddrType & 1)

	return result
}

type SMP struct {
	ioCapability byte

	storedKeys        map[smpStoredLTKMapKey]smpStoredLTK
	storedKeysPersist *gobpersist.GobPersist
}

func New() *SMP {
	s := &SMP{
		storedKeys: make(map[smpStoredLTKMapKey]smpStoredLTK),
	}

	s.storedKeysPersist = &gobpersist.GobPersist{
		Target:   &s.storedKeys,
		Filename: "/tmp/smp.gob",
	}

	log.Println(s.storedKeysPersist.Load())
	log.Println(s.storedKeysPersist.Save())

	return s
}

type SMPConn struct {
	parent *SMP

	conn   hciconnmgr.BufferConn
	logger *logrus.Entry

	pduRx chan (*pdu.PDU)

	protocol smpProtocol

	timeout *time.Timer

	addrLERemote bleutil.BLEAddr
	addrLELocal  bleutil.BLEAddr
}

func (c *SMPConn) updateTimeout(run bool) {
	if !run {
		c.timeout.Stop()
	} else {
		c.timeout.Reset(30 * time.Second)
	}
}

func (c *SMPConn) smpHandler() {
	defer func() {
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

	timedOut := false

	for {
		select {
		case <-c.timeout.C:
			c.logger.Warn("Security manager timeout")
			timedOut = true

		case pdu, ok := <-c.pduRx:
			if !ok {
				return
			}

			if c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				c.logger.WithField("0pdu", hex.EncodeToString(pdu.Buf())).Debug("PSM Received")
			}

			keepBuffer := false
			if !timedOut {
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

func (c *SMPConn) encryptWithLTK() (bool, error) {
	c.parent.storedKeysPersist.Lock()
	ltk, ok := c.parent.storedKeys[makeSMPStoredLTKMapKey(c.addrLELocal, c.addrLERemote)]
	c.parent.storedKeysPersist.Unlock()
	if !ok {
		return false, nil
	}

	raw := c.rawConnLE()
	return true, raw.Encrypt(ltk.EDIV, ltk.Rand, ltk.LTK)
}

func (p *SMP) AddConn(conn hciconnmgr.BufferConn) *SMPConn {
	c := &SMPConn{
		parent:  p,
		conn:    conn,
		logger:  bleutil.LogWithPrefix(conn.GetLogger(), "psm"),
		pduRx:   make(chan (*pdu.PDU), 1),
		timeout: time.NewTimer(0),
	}

	<-c.timeout.C

	raw := c.rawConnLE()
	c.addrLERemote = raw.RemoteAddr().(bleutil.BLEAddr)
	c.addrLELocal = raw.LocalAddr().(bleutil.BLEAddr)

	go c.reader()
	go c.smpHandler()

	c.parent.ioCapability = 0x2
	ok, err := c.encryptWithLTK()
	log.Println("Using ltk", ok, err)
	//c.sendPairingRequest(0x5) //MTIM, bonding

	return c
}
