package bleatt

import (
	"bytes"
	"context"
	"encoding/binary"

	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"
)

type attServerWriteQueueEntry struct {
	idx     uint16
	handle  *attstructure.GATTHandle
	offset  uint16
	payload []byte
}

const (
	// maxFindByTypeValueLength caps the length of the value-to-match in
	// ATTFindByTypeValueReq. The spec allows up to (MTU-7), but allowing
	// the full MTU forces every probed handle through a bytes.Equal of
	// peer-controlled bytes — large enough to amplify CPU on the server.
	maxFindByTypeValueLength = 64

	// maxPrepareWriteQueueBytes bounds the total bytes a single peer can
	// have queued across all PrepareWrite fragments before commit.
	maxPrepareWriteQueueBytes = 4096
)

type attServer struct {
	parent *GattDevice

	localStructure *attstructure.ExportedStructure
}

func normalizeATTMTU(mtu uint16) uint16 {
	if mtu == 0 {
		return 0xFFFF
	}
	if mtu < 23 {
		return 23
	}
	return mtu
}

func (a *attServer) init(parent *GattDevice, localStructure *attstructure.ExportedStructure) error {
	a.parent = parent
	a.localStructure = localStructure

	localStructure.HandleSet = func(c *attstructure.Characteristic, value []byte) (int, error) {
		return a.characteristicNotify(context.Background(), c, value)

	}

	return nil
}

func sendError(conn *gattDeviceConn, method ATTCommand, handle uint16, errorCode ATTError) error {
	buf := bleutil.GetBuffer(5)
	buf.Buf()[0] = byte(ATTErrorRsp)
	buf.Buf()[1] = byte(method)
	binary.LittleEndian.PutUint16(buf.Buf()[2:4], handle)
	buf.Buf()[4] = byte(errorCode)

	if conn.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		conn.logger.WithFields(logrus.Fields{
			"2buf":    buf,
			"0handle": handle,
			"1code":   errorCode,
		}).Debug("ATT Server Error")
	}

	return conn.conn.WriteBuffer(buf)
}

func (a *attServer) handleMTUReq(conn *gattDeviceConn, buf *pdu.PDU) (bool, error) {
	if buf.Len() != 2 {
		return false, ErrorProtocolViolation
	}

	requestedMTU := binary.LittleEndian.Uint16(buf.Buf())

	buf.Reset()
	rsp := buf.ExtendRight(3)
	rsp[0] = byte(ATTExchangeMTURsp)
	binary.LittleEndian.PutUint16(rsp[1:], a.parent.ourMTU)

	conn.setMTU(requestedMTU)

	return true, a.write(conn, buf)
}

//Try to combine ATTFindInformationReq ATTReadByGroupTypeReq ATTReadByTypeReq ATTFindByValueTypeReq
//ATTFindInformationReq: Starting, Ending               -> Format, [Handle, UUID]
//ATTReadByGroupTypeReq: Starting, Ending, UUID         -> Length, [Handle, GroupHandle, Value] (check if thing is grouping element)
//ATTReadByTypeReq:      Starting, Ending, UUID         -> Length, [Handle, Value]
//ATTFindByValueTypeReq: Starting, Ending, UUID2, Value -> [Handle, GroupHandle]

func (a *attServer) addPayload(conn *gattDeviceConn, buf *pdu.PDU, payload []byte) (bool, int) {
	mtu := conn.getMTU()
	maxLen := mtu - buf.Len()
	if maxLen <= 0 {
		return true, 0
	}
	if len(payload) < maxLen {
		maxLen = len(payload)
	}
	value := buf.ExtendRight(maxLen)
	bytes := copy(value, payload)

	return buf.Len() == int(mtu), bytes
}

func (a *attServer) handleDiscovery(conn *gattDeviceConn, method ATTCommand, buf *pdu.PDU) (bool, error) {
	if buf.Len() < 4 {
		return false, ErrorProtocolViolation
	}

	startHandle := binary.LittleEndian.Uint16(buf.Buf()[0:2])
	endHandle := binary.LittleEndian.Uint16(buf.Buf()[2:4])

	if startHandle > endHandle || startHandle == 0 {
		return false, sendError(conn, method, startHandle, ATTErrorInvalidHandle)
	}

	var uuid bleutil.UUID
	checkUUID := false
	if method == ATTReadByGroupTypeReq || method == ATTReadByTypeReq {
		var valid bool
		uuid, valid = bleutil.UUIDFromBytesValid(buf.Buf()[4:])
		if !valid {
			return false, ErrorProtocolViolation
		}

		checkUUID = true
	}

	var checkValue []byte
	if method == ATTFindByTypeValueReq {
		if buf.Len() < 6 {
			return false, ErrorProtocolViolation
		}

		uuid = bleutil.UUIDFromBytes(buf.Buf()[4:6])
		checkValue = buf.Buf()[6:]
		/* Cap the value-to-match length: the spec only requires up to
		   the negotiated MTU but we don't actually serve attributes
		   bigger than this and a peer could otherwise force a
		   bytes.Equal scan over very large payloads. */
		if len(checkValue) > maxFindByTypeValueLength {
			return false, sendError(conn, method, startHandle, ATTErrorInvalidPDU)
		}
	}

	buf.Reset()

	hasResults := false
	header := byte(0)

	for _, m := range a.localStructure.Handles {
		if m.Info.Handle > endHandle {
			break
		}

		if m.Info.Handle < startHandle {
			continue
		}

		if checkUUID && m.Info.UUID != uuid {
			continue
		}

		if method == ATTReadByGroupTypeReq && m.Info.GroupEndHandle == 0 {
			continue
		}

		mHeader := byte(0)
		needHeader := 1
		extra := 0
		addValue := false
		switch method {
		case ATTFindInformationReq:
			mHeader = byte(1)
			extra = 2 + m.Info.UUID.GetLength()
			if m.Info.UUID.GetLength() > 2 {
				mHeader = byte(2)
			}
		case ATTReadByTypeReq:
			//Todo check max len
			extra = 2
			mHeader = byte(extra + len(m.Value))
			addValue = true
		case ATTReadByGroupTypeReq:
			extra = 2 + 2
			mHeader = byte(extra + len(m.Value))
			addValue = true
		case ATTFindByTypeValueReq:
			needHeader = 0
			extra = 2 + 2
		}

		/* The security check has to run on every method that touches m.Value.
		   ATTFindByTypeValueReq does a bytes.Equal against the attribute value;
		   without this check, an unencrypted peer could probe the contents of
		   read-encryption-required attributes by guess-and-check. */
		needSecCheck := addValue || (method == ATTFindByTypeValueReq)
		if needSecCheck {
			secErr := a.checkSecurity(conn, true, m.Info.Flags)
			if secErr != ATTErrorNone {
				return false, sendError(conn, method, m.Info.Handle, secErr)
			}
		}

		if checkValue != nil && !bytes.Equal(m.Value, checkValue) {
			continue
		}

		if !hasResults {
			hdr := buf.ExtendRight(1 + needHeader)
			hdr[0] = byte(method + 1)
			header = mHeader
			if needHeader > 0 {
				hdr[1] = header
			}
			hasResults = true
		}

		if needHeader > 0 && mHeader != header {
			break
		}

		if buf.Len()+extra > conn.getMTU() {
			break
		}

		data := buf.ExtendRight(extra)
		switch method {
		case ATTFindInformationReq:
			binary.LittleEndian.PutUint16(data[0:2], m.Info.Handle)
			copy(data[2:], m.Info.UUID.UUIDToBytes())
		case ATTReadByTypeReq:
			binary.LittleEndian.PutUint16(data[0:2], m.Info.Handle)
		case ATTReadByGroupTypeReq:
			binary.LittleEndian.PutUint16(data[0:2], m.Info.Handle)
			binary.LittleEndian.PutUint16(data[2:4], m.Info.GroupEndHandle)
		case ATTFindByTypeValueReq:
			binary.LittleEndian.PutUint16(data[0:2], m.Info.Handle)
			binary.LittleEndian.PutUint16(data[2:4], m.Info.GroupEndHandle)
		}

		if addValue {
			a.localStructure.Lock()
			if m.ValueConfig.ValueBeforeReadCb != nil {
				m.ValueConfig.ValueBeforeReadCb(m, 0)
			}
			_, bytes := a.addPayload(conn, buf, m.Value)
			if m.ValueConfig.ValueAfterReadCb != nil {
				m.ValueConfig.ValueAfterReadCb(m, 0, bytes)
			}
			a.localStructure.Unlock()
		}
	}

	if hasResults {
		return true, a.write(conn, buf)
	}

	return false, sendError(conn, method, startHandle, ATTErrorAttributeNotFound)
}

func (a *attServer) findHandle(handle uint16) *attstructure.GATTHandle {
	for i, m := range a.localStructure.Handles {
		if m.Info.Handle == handle {
			return a.localStructure.Handles[i]
		}
	}
	return nil
}

func (a *attServer) checkSecurity(conn *gattDeviceConn, isRead bool, flags attstructure.CharacteristicFlag) ATTError {
	if isRead && flags&attstructure.CharacteristicRead == 0 {
		return ATTErrorReadNotPermitted
	}

	if !isRead && (flags&(attstructure.CharacteristicWriteAck|attstructure.CharacteristicWriteNoAck) == 0) {
		return ATTErrorWriteNotPermitted
	}

	needEncryption := false
	if isRead && flags&attstructure.CharacteristicReadNeedsEncryption > 0 {
		needEncryption = true
	}
	if !isRead && flags&attstructure.CharacteristicWriteNeedsEncryption > 0 {
		needEncryption = true
	}

	needAuthentication := false
	if isRead && flags&attstructure.CharacteristicReadNeedsAuthentication > 0 {
		needAuthentication = true
	}
	if !isRead && flags&attstructure.CharacteristicWriteNeedsAuthentication > 0 {
		needAuthentication = true
	}

	if needEncryption || needAuthentication {
		if conn.smpConn == nil {
			return ATTErrorUnlikelyError
		}
		encryption, authentication, _ := conn.smpConn.GetSecurity()
		if needEncryption && !encryption {
			return ATTErrorInsufficientEncryption
		}
		if needAuthentication && !authentication {
			return ATTErrorInsufficientAuthentication
		}
	}

	return ATTErrorNone
}

func (a *attServer) handleReadReq(conn *gattDeviceConn, method ATTCommand, buf *pdu.PDU) (bool, error) {
	if (method == ATTReadReq && buf.Len() != 2) || (method == ATTReadBlobReq && buf.Len() != 4) {
		return false, ErrorProtocolViolation
	}

	idx := binary.LittleEndian.Uint16(buf.DropLeft(2))
	handle := a.findHandle(idx)
	if handle == nil {
		return false, sendError(conn, method, idx, ATTErrorInvalidHandle)
	}

	secErr := a.checkSecurity(conn, true, handle.Info.Flags)
	if secErr != ATTErrorNone {
		return false, sendError(conn, method, idx, secErr)
	}

	offset := 0
	if method == ATTReadBlobReq {
		offset = int(binary.LittleEndian.Uint16(buf.DropLeft(2)))
	}

	buf.Reset()
	buf.Append(byte(method) + 1)

	errCode := ATTError(0)

	a.localStructure.Lock()
	if handle.ValueConfig.ValueBeforeReadCb != nil {
		handle.ValueConfig.ValueBeforeReadCb(handle, offset)
	}
	if method == ATTReadBlobReq && len(handle.Value)+1 <= conn.getMTU() {
		errCode = ATTErrorAttributeNotLong
	} else if offset > len(handle.Value) {
		errCode = ATTErrorInvalidOffset
	} else {
		_, bytes := a.addPayload(conn, buf, handle.Value[offset:])
		if handle.ValueConfig.ValueAfterReadCb != nil {
			handle.ValueConfig.ValueAfterReadCb(handle, offset, bytes)
		}
	}
	a.localStructure.Unlock()

	if errCode != 0 {
		return false, sendError(conn, method, idx, errCode)
	}

	return true, a.write(conn, buf)
}

func (a *attServer) handleReadReqMultiple(conn *gattDeviceConn, method ATTCommand, buf *pdu.PDU) (bool, error) {
	/* Per spec, the request body is a sequence of 16-bit handles. An odd
	   number of bytes is a protocol violation. */
	if buf.Len() < 4 || buf.Len()%2 != 0 {
		return false, sendError(conn, method, 0, ATTErrorInvalidPDU)
	}

	resp := bleutil.GetBuffer(1)
	resp.Buf()[0] = byte(method + 1)

	addLen := method == ATTReadMultipleValueReq

	for {
		if buf.Len() < 2 {
			break
		}
		idx := binary.LittleEndian.Uint16(buf.DropLeft(2))
		handle := a.findHandle(idx)
		if handle == nil {
			bleutil.ReleaseBuffer(resp)
			return false, sendError(conn, method, idx, ATTErrorInvalidHandle)
		}

		secErr := a.checkSecurity(conn, true, handle.Info.Flags)
		if secErr != ATTErrorNone {
			bleutil.ReleaseBuffer(resp)
			return false, sendError(conn, method, idx, secErr)
		}

		if addLen {
			if resp.Len()+2 > conn.getMTU() {
				break
			}

			binary.LittleEndian.PutUint16(resp.ExtendRight(2), uint16(len(handle.Value)))
		}

		full, _ := a.addPayload(conn, resp, handle.Value)
		if full {
			break
		}
	}

	return false, a.write(conn, resp)
}

func (a *attServer) handleWriteReq(conn *gattDeviceConn, method ATTCommand, buf *pdu.PDU) (bool, error) {
	if buf.Len() < 2 {
		return false, ErrorProtocolViolation
	}

	idx := binary.LittleEndian.Uint16(buf.DropLeft(2))
	handle := a.findHandle(idx)
	if handle == nil {
		return false, sendError(conn, method, idx, ATTErrorInvalidHandle)
	}

	secErr := a.checkSecurity(conn, false, handle.Info.Flags)
	if secErr != ATTErrorNone {
		return false, sendError(conn, method, idx, secErr)
	}
	if method == ATTWriteReq && handle.Info.Flags&attstructure.CharacteristicWriteAck == 0 {
		return false, sendError(conn, method, idx, ATTErrorWriteNotPermitted)
	}
	if method == ATTWriteCMD && handle.Info.Flags&attstructure.CharacteristicWriteNoAck == 0 {
		return false, nil
	}
	if handle.ValueConfig.LengthFixed && buf.Len() != len(handle.Value) {
		if method == ATTWriteReq {
			return false, sendError(conn, method, idx, ATTErrorLength)
		}
		return false, nil
	}
	if handle.ValueConfig.LengthMax > 0 && buf.Len() > int(handle.ValueConfig.LengthMax) {
		if method == ATTWriteReq {
			return false, sendError(conn, method, idx, ATTErrorLength)
		}
		return false, nil
	}

	a.localStructure.Lock()
	if handle.ValueConfig.LengthFixed {
		copy(handle.Value, buf.Buf())
	} else {
		handle.Value = append(handle.Value[:0], buf.Buf()...)
	}
	if handle.ValueConfig.ValueWriteCb != nil {
		handle.ValueConfig.ValueWriteCb(handle)
	}
	a.localStructure.Unlock()

	if method == ATTWriteReq {
		buf.Reset()
		buf.Append(byte(ATTWriteRsp))
		return true, a.write(conn, buf)
	}

	return false, nil
}

func (a *attServer) handlePrepateWriteReq(conn *gattDeviceConn, buf *pdu.PDU) (bool, error) {
	if buf.Len() < 4 {
		return false, ErrorProtocolViolation
	}

	idx := binary.LittleEndian.Uint16(buf.Buf())
	handle := a.findHandle(idx)
	if handle == nil {
		return false, sendError(conn, ATTPrepareWriteReq, idx, ATTErrorInvalidHandle)
	}

	secErr := a.checkSecurity(conn, false, handle.Info.Flags)
	if secErr != ATTErrorNone {
		return false, sendError(conn, ATTPrepareWriteReq, idx, secErr)
	}

	offset := binary.LittleEndian.Uint16(buf.Buf()[2:])
	payloadLen := buf.Len() - 4

	/* Per-fragment bounds so a chatty peer cannot exhaust memory before
	   the commit-time check fires. */
	if handle.ValueConfig.LengthMax > 0 && int(offset)+payloadLen > int(handle.ValueConfig.LengthMax) {
		return false, sendError(conn, ATTPrepareWriteReq, idx, ATTErrorLength)
	}

	a.localStructure.Lock()
	if len(conn.writeQueue) >= 64 {
		a.localStructure.Unlock()
		return false, sendError(conn, ATTPrepareWriteReq, idx, ATTErrorPrepareQueueFull)
	}

	queuedBytes := payloadLen
	for _, q := range conn.writeQueue {
		queuedBytes += len(q.payload)
	}
	if queuedBytes > maxPrepareWriteQueueBytes {
		a.localStructure.Unlock()
		return false, sendError(conn, ATTPrepareWriteReq, idx, ATTErrorPrepareQueueFull)
	}

	payload := make([]byte, payloadLen)
	copy(payload, buf.Buf()[4:])

	conn.writeQueue = append(conn.writeQueue, attServerWriteQueueEntry{
		idx:     idx,
		handle:  handle,
		offset:  offset,
		payload: payload,
	})
	a.localStructure.Unlock()

	buf.ExtendLeft(1)[0] = byte(ATTPrepareWriteRsp)
	return true, a.write(conn, buf)
}

func (a *attServer) handleExecuteWriteReq(conn *gattDeviceConn, buf *pdu.PDU) (bool, error) {
	if buf.Len() != 1 {
		return false, ErrorProtocolViolation
	}
	flag := buf.Buf()[0]
	if flag > 1 {
		return false, sendError(conn, ATTExecuteWriteReq, 0, ATTErrorInvalidPDU)
	}

	errCode := ATTError(0)
	errIdx := uint16(0)

	type savedValue struct {
		bytes []byte
		cap   int
	}

	a.localStructure.Lock()
	if flag == 1 {
		/* Save the original slice header (length + capacity) so the
		   rollback can preserve any LengthFixed-style preallocated
		   buffer. The previous version restored a slice with cap==len
		   and silently broke fixed-capacity invariants. */
		originalMap := make(map[*attstructure.GATTHandle]savedValue)
		for _, m := range conn.writeQueue {
			if _, ok := originalMap[m.handle]; !ok {
				snapshot := make([]byte, len(m.handle.Value), cap(m.handle.Value))
				copy(snapshot, m.handle.Value)
				originalMap[m.handle] = savedValue{bytes: snapshot, cap: cap(m.handle.Value)}
			}
		}

		for _, m := range conn.writeQueue {
			if int(m.offset) > len(m.handle.Value) {
				errIdx = m.idx
				errCode = ATTErrorInvalidOffset
				goto fail
			}

			minLen := int(m.offset) + len(m.payload)
			if m.handle.ValueConfig.LengthMax > 0 && minLen > int(m.handle.ValueConfig.LengthMax) {
				errIdx = m.idx
				errCode = ATTErrorLength
				goto fail
			}
			if m.handle.ValueConfig.LengthFixed && minLen > len(m.handle.Value) {
				errIdx = m.idx
				errCode = ATTErrorLength
				goto fail
			}

			if cap(m.handle.Value) < minLen {
				n := make([]byte, minLen)
				copy(n, m.handle.Value)
				m.handle.Value = n
			}

			if !m.handle.ValueConfig.LengthFixed {
				m.handle.Value = m.handle.Value[:minLen]
			}

			copy(m.handle.Value[m.offset:], m.payload)
		}

	fail:
		if errCode != 0 {
			for h, sv := range originalMap {
				/* Restore the original slice header — including the
				   original capacity — so LengthFixed invariants survive. */
				h.Value = sv.bytes[:len(sv.bytes):sv.cap]
			}
		} else {
			for i := range originalMap {
				if i.ValueConfig.ValueWriteCb != nil {
					i.ValueConfig.ValueWriteCb(i)
				}
			}
		}

	}

	conn.writeQueue = conn.writeQueue[:0]
	a.localStructure.Unlock()

	if errCode != 0 {
		return false, sendError(conn, ATTExecuteWriteReq, errIdx, errCode)
	}

	buf.Reset()
	buf.Append(byte(ATTExecuteWriteRsp))
	return true, a.write(conn, buf)
}

func (a *attServer) handlePDU(conn *gattDeviceConn, method ATTCommand, isAuthenticated bool, buf *pdu.PDU) (bool, error) {
	switch method {
	case ATTExchangeMTUReq:
		return a.handleMTUReq(conn, buf)

	/* Reading */
	case ATTReadReq:
		fallthrough
	case ATTReadBlobReq:
		return a.handleReadReq(conn, method, buf)
	case ATTReadMultipleReq:
		fallthrough
	case ATTReadMultipleValueReq:
		return a.handleReadReqMultiple(conn, method, buf)

	/* Writing */
	case ATTWriteReq:
		fallthrough
	case ATTWriteCMD:
		return a.handleWriteReq(conn, method, buf)
	case ATTSignedWriteCMD:
		return false, nil

	/* Transaction write */
	case ATTPrepareWriteReq:
		return a.handlePrepateWriteReq(conn, buf)
	case ATTExecuteWriteReq:
		return a.handleExecuteWriteReq(conn, buf)

	/* Ack indication */
	case ATTHandleValueCNF:

	/* Discovery and special reads */
	case ATTFindInformationReq:
		fallthrough
	case ATTFindByTypeValueReq:
		fallthrough
	case ATTReadByGroupTypeReq:
		fallthrough
	case ATTReadByTypeReq:
		return a.handleDiscovery(conn, method, buf)
	}

	return false, sendError(conn, method, 0, ATTErrorRequestNotSupported)
}

func (a *attServer) characteristicNotify(ctx context.Context, characteristic *attstructure.Characteristic, value []byte) (int, error) {
	handle := characteristic.ValueHandle

	ccc := handle.CCCHandle
	idx := handle.Info.Handle

	/* If there is no client config descriptor we can't notify */
	if ccc == nil {
		return 0, nil
	}

	a.localStructure.Lock()
	flags := ccc.Value[0]
	a.localStructure.Unlock()

	cmd := ATTCommand(0)
	if (flags&2 > 0) && (handle.Info.Flags&attstructure.CharacteristicIndicate > 0) {
		//Indication
		cmd = ATTHandleValueIND
	} else if (flags&1 > 0) && (handle.Info.Flags&attstructure.CharacteristicNotify > 0) {
		//Notification
		cmd = ATTHandleValueNTF
	} else {
		return 0, nil
	}

	conn := a.parent.getConnWithHighestMTU()
	if conn == nil {
		/* Without a connection we can't notify either */
		return -1, nil
	}

	/* Before we can do any indication MTU negotiation must be finished (otherwise the receiver doesn't know the fragment size) */
	conn.getMTUBlocking()

	buf := bleutil.GetBuffer(3)

	buf.Buf()[0] = byte(cmd)
	binary.LittleEndian.PutUint16(buf.Buf()[1:], idx)

	_, bytes := a.addPayload(conn, buf, value)

	if conn.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		conn.logger.WithFields(logrus.Fields{
			"0buf": buf,
		}).Debug("ATT Notify")
	}

	_, resp, err := conn.client.sendCommand(ctx, buf, cmd == ATTHandleValueIND)
	bleutil.ReleaseBuffer(resp)
	return bytes, err
}

func (a *attServer) write(conn *gattDeviceConn, buf *pdu.PDU) error {
	if conn.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		conn.logger.WithFields(logrus.Fields{
			"0buf": buf,
		}).Debug("ATT Server Write")
	}

	return conn.conn.WriteBuffer(buf)
}
