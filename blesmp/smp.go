package blesmp

import (
	"encoding/hex"

	"github.com/BertoldVdb/go-ble/bleconnecter"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	blel2cap "github.com/BertoldVdb/go-ble/l2cap"
	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"
)

type SMP struct {
	ioCapability byte
}

func New() *SMP {
	return &SMP{}
}

type SMPConn struct {
	parent *SMP

	conn   hciconnmgr.BufferConn
	logger *logrus.Entry

	pduRx chan (*pdu.PDU)

	protocol smpProtocol
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

	for {
		select {
		case pdu, ok := <-c.pduRx:
			if !ok {
				return
			}

			if c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				c.logger.WithField("0pdu", hex.EncodeToString(pdu.Buf())).Debug("PSM Received")
			}

			keepBuffer := c.handleMessage(pdu)
			if !keepBuffer {
				bleutil.ReleaseBuffer(pdu)
			}
		}
	}
}

func (c *SMPConn) rawConnLE() *bleconnecter.BLEConnection {
	return c.conn.(*blel2cap.L2Connection).ParentConn().(*bleconnecter.BLEConnection)
}

func (p *SMP) AddConn(conn hciconnmgr.BufferConn) *SMPConn {
	c := &SMPConn{
		parent: p,
		conn:   conn,
		logger: bleutil.LogWithPrefix(conn.GetLogger(), "psm"),
		pduRx:  make(chan (*pdu.PDU), 1),
	}

	go c.reader()
	go c.smpHandler()

	c.parent.ioCapability = 0x2
	c.sendPairingRequest(0x5) //MTIM, bonding

	return c
}
