package hciconnmgr

import (
	"encoding/binary"
	"encoding/hex"

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
		conn.rxBuffer = conn.rxBuffer[:0]
	}

	/* Append fragment */
	conn.rxBuffer = append(conn.rxBuffer, payload...)

	/* Complete L2CAP packet? */
	conn.handleACLData()

	return true
}

func (c *Connection) handleACLData() error {
	workBuf := c.rxBuffer

	for {
		bufLen := len(workBuf)
		if bufLen < 4 {
			break
		}

		pktLen := int(binary.LittleEndian.Uint16(workBuf[0:2])) + 4
		if pktLen < bufLen {
			break
		}

		pktBuf := c.connmgr.rxtxFreeBuffers.PopOrCreate(pktLen)
		copy(pktBuf, workBuf[:pktLen])

		fifoLen := c.rxFIFO.Push(pktBuf)
		select {
		case c.rxNewDataChan <- struct{}{}:
		default:
		}

		if c.connmgr.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
			c.connmgr.logger.WithFields(logrus.Fields{
				"0handle":  c.handle,
				"1fifolen": fifoLen,
				"2l2cap":   hex.EncodeToString(pktBuf),
			}).Trace("Received valid L2CAP/ACL packet")
		}

		workBuf = workBuf[pktLen:]
	}

	copy(c.rxBuffer, workBuf)
	c.rxBuffer = c.rxBuffer[:len(workBuf)]
	return nil
}

func (c *Connection) queueACLPacket(flagPB int, flagBC int, payload []byte) {
	handle := c.handle
	handle |= uint16(flagPB&0x3) << 12
	handle |= uint16(flagBC&0x3) << 14

	/* Add HCI header */
	buf := c.connmgr.rxtxFreeBuffers.PopOrCreate(5 + len(payload))
	buf[0] = 2
	binary.LittleEndian.PutUint16(buf[1:3], handle)
	binary.LittleEndian.PutUint16(buf[3:5], uint16(len(payload)))

	/* Add payload */
	copy(buf[5:], payload)

	fifoLen := c.txFIFO.Push(buf)

	if c.connmgr.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.connmgr.logger.WithFields(logrus.Fields{
			"0handle":  c.handle,
			"1fifolen": fifoLen,
			"2l2cap":   hex.EncodeToString(buf),
		}).Trace("Queued ACL fragment for TX")
	}
}

func (c *Connection) encodeACL(l2cap []byte) {
	if c.connmgr.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.connmgr.logger.WithFields(logrus.Fields{
			"0handle": c.handle,
			"1l2cap":  hex.EncodeToString(l2cap),
		}).Trace("Preparing L2CAP/ACL fragments for TX")
	}

	mtu := c.txSlotManager.GetBufferLength()

	flagPB := 0
	for len(l2cap) > 0 {
		fragmentLen := len(l2cap)
		if fragmentLen > mtu {
			fragmentLen = mtu
		}
		c.queueACLPacket(flagPB, 0, l2cap[:fragmentLen])
		l2cap = l2cap[fragmentLen:]

		/* The next fragments are continuation */
		flagPB = 1
	}

	select {
	case c.txSlotManager.newFragmentsChan <- struct{}{}:
	default:
	}
}
