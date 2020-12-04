package hciconnmgr

import (
	"encoding/binary"
	"encoding/hex"

	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"
)

func (c *ConnectionManager) handleACL(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	payloadLen := int(binary.LittleEndian.Uint16(data[2:4]))
	if payloadLen+4 != len(data) {
		return false
	}

	handle := binary.LittleEndian.Uint16(data[0:2])
	flagPB := (handle >> 12) & 0x3
	flagBC := (handle >> 14) & 0x3
	handle &= 0xFFF
	payload := data[4:]

	/* Ignore EDR broadcast traffic since we don't support it */
	if flagBC != 0 {
		return true
	}

	/* flagPB=2 -> Start of burst from controller to host, 1 -> new fragment */
	if flagPB != 1 && flagPB != 2 {
		return true
	}

	/* Look up the connection */
	c.RLock()
	conn, ok := c.connections[handle]
	c.RUnlock()

	if c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			"0handle":  handle,
			"1known":   ok,
			"2flagPB":  flagPB,
			"3flagBC":  flagBC,
			"4payload": hex.EncodeToString(payload),
		}).Trace("Received ACL fragment")
	}

	if !ok {
		return true
	}

	if flagPB == 2 {
		/* New fragment, clear reassembly buffer */
		conn.rxPDU.Reset()
	}

	/* Append fragment */
	conn.rxPDU.Append(payload...)

	/* Complete L2CAP packet? */
	conn.handleACLData()

	return true
}

func (c *Connection) handleACLData() error {
	for {
		if c.rxPDU.Len() < 4 {
			break
		}

		pktLen := int(binary.LittleEndian.Uint16(c.rxPDU.Buf()[0:2])) + 4
		payload := c.rxPDU.DropLeft(pktLen)
		if payload == nil {
			break
		}

		pktBuf := bleutil.CopyBufferFromSlice(payload)

		pktPrint := ""
		if c.connmgr.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
			pktPrint = pktBuf.String()
		}

		fifoLen := c.rxFIFO.Push(pktBuf)
		select {
		case c.rxNewDataChan <- struct{}{}:
		default:
		}

		if c.connmgr.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
			c.connmgr.logger.WithFields(logrus.Fields{
				"0handle":  c.handle,
				"1fifolen": fifoLen,
				"2l2cap":   pktPrint,
			}).Trace("Received valid L2CAP/ACL packet")
		}
	}

	/* Check that the leftcap is not growing too much to avoid memory DoS */
	if c.rxPDU.LeftCap() > 8192 {
		c.rxPDU.NormalizeLeft(0)
	}

	return nil
}

func (c *Connection) queueACLPacket(flagPB int, flagBC int, payload *pdu.PDU) {
	handle := c.handle
	handle |= uint16(flagPB&0x3) << 12
	handle |= uint16(flagBC&0x3) << 14

	/* Add HCI header */
	pl := payload.Len()
	header := payload.ExtendLeft(5)
	header[0] = 2
	binary.LittleEndian.PutUint16(header[1:3], handle)
	binary.LittleEndian.PutUint16(header[3:5], uint16(pl))

	pktPrint := ""
	if c.connmgr.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		pktPrint = payload.String()
	}

	fifoLen := c.txFIFO.Push(payload)

	if c.connmgr.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.connmgr.logger.WithFields(logrus.Fields{
			"0handle":  c.handle,
			"1fifolen": fifoLen,
			"2l2cap":   pktPrint,
		}).Trace("Queued ACL fragment for TX")
	}
}

func (c *Connection) encodeACL(l2capP *pdu.PDU) {
	if c.connmgr.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.connmgr.logger.WithFields(logrus.Fields{
			"0handle": c.handle,
			"1l2cap":  l2capP,
		}).Trace("Preparing L2CAP/ACL fragments for TX")
	}

	mtu := c.txSlotManager.GetBufferLength()
	flagPB := 0
	for {
		if l2capP.Len() > mtu {
			/* If the buffer is too long, create a new one and take it from the l2capP */
			frag := bleutil.CopyBufferFromSlice(l2capP.DropLeft(mtu))
			c.queueACLPacket(flagPB, 0, frag)
		} else {
			/* Finally, just pass on l2capP itself */
			c.queueACLPacket(flagPB, 0, l2capP)
			break
		}

		/* The next fragments are continuation */
		flagPB = 1
	}

	select {
	case c.txSlotManager.newFragmentsChan <- struct{}{}:
	default:
	}
}
