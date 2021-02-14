package bleatt

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/BertoldVdb/go-misc/slotset"
	"github.com/sirupsen/logrus"
)

type attClientCmdData struct {
	method ATTCommand
	buf    *pdu.PDU
}

type attClient struct {
	ctxExpired context.Context

	parent *gattDeviceConn
	cmdmgr *slotset.SlotSet

	// 	dbHashRequest sync.Once
	// 	dbHash        []byte

	timeoutTimerMutex sync.Mutex
	timeoutTimer      *time.Timer
}

func (a *attClient) init(parent *gattDeviceConn) error {
	*a = attClient{
		parent: parent,

		cmdmgr: slotset.New(1, func(slot *slotset.Slot) {
			slot.Data = &attClientCmdData{}
		}),

		timeoutTimer: time.NewTimer(0),
	}

	ctxExpired, cancel := context.WithCancel(context.Background())
	cancel()

	a.ctxExpired = ctxExpired

	<-a.timeoutTimer.C
	go a.handleTimeout()

	return nil
}

func (a *attClient) sendCommand(ctx context.Context, cmd *pdu.PDU, withReply bool) (ATTCommand, *pdu.PDU, error) {
	slot, err := a.cmdmgr.Get(ctx)
	if err != nil {
		return 0, nil, err
	}
	defer a.cmdmgr.Put(slot)
	if !withReply {
		return 0, nil, a.write(cmd, false)
	}

	slot.Activate()

	err = a.write(cmd, true)
	if err != nil {
		slot.Deactivate()
		return 0, nil, err
	}

	_, err = slot.WaitCtx(ctx)
	slot.Deactivate()
	if err != nil {
		return 0, nil, err
	}

	data := slot.Data.(*attClientCmdData)
	return data.method, data.buf, nil
}

func (a *attClient) sendCommandErrRsp(ctx context.Context, req *pdu.PDU) (ATTCommand, *pdu.PDU, ATTError, error) {
	method := req.Buf()[0]
	cmd, response, err := a.sendCommand(ctx, req, true)
	if err == nil && cmd == ATTErrorRsp {
		data := response.DropLeft(4)
		bleutil.ReleaseBuffer(response)

		if data == nil || data[0] != method {
			return cmd, nil, 0, ErrorProtocolViolation
		}

		return cmd, nil, ATTError(data[3]), nil
	}

	return cmd, response, 0, err
}

func (a *attClient) write(buf *pdu.PDU, expectReply bool) error {
	if a.parent.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		a.parent.logger.WithFields(logrus.Fields{
			"0buf":         buf,
			"1expectReply": expectReply,
		}).Trace("ATT Client Write")
	}

	if expectReply {
		a.timeoutTimerMutex.Lock()
		a.timeoutTimer.Reset(30 * time.Second)
		a.timeoutTimerMutex.Unlock()
	}

	return a.parent.conn.WriteBuffer(buf)
}

func (a *attClient) handleNotify(handle uint16, data []byte) {
	structure := a.parent.parent.ClientGetStructure(a.ctxExpired)

	if structure != nil {
		structure.InjectNotify(handle, data)
	}
}

func (a *attClient) handleNTFIND(method ATTCommand, buf *pdu.PDU) (bool, error) {
	if method == ATTMultipleHandleValueNTF {
		for {
			hdr := buf.DropLeft(4)
			if hdr == nil {
				break
			}

			handle := binary.LittleEndian.Uint16(hdr)
			dlen := binary.LittleEndian.Uint16(hdr[2:])

			data := buf.DropLeft(int(dlen))
			if hdr == nil {
				break
			}

			if a.parent.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				a.parent.logger.WithFields(logrus.Fields{
					"0handle": handle,
					"1data":   hex.EncodeToString(data),
				}).Debug("Got mutli notification")
			}

			a.handleNotify(handle, data)
		}

		return false, nil
	}

	handleBuf := buf.DropLeft(2)
	if handleBuf != nil {
		handle := binary.LittleEndian.Uint16(handleBuf)
		isIndication := method == ATTHandleValueIND

		if a.parent.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
			a.parent.logger.WithFields(logrus.Fields{
				"0handle":       handle,
				"1isIndication": isIndication,
				"2data":         buf,
			}).Debug("Got notification")
		}

		a.handleNotify(handle, buf.Buf())

		if isIndication {
			/* Confirm notification */
			buf.Reset()
			buf.Append(byte(ATTHandleValueCNF))
			a.write(buf, false)
			return true, nil
		}
	}

	return false, nil
}

func (a *attClient) handleTimeout() {
	defer a.parent.parent.CloseConn(a.parent.conn)

	a.timeoutTimerMutex.Lock()
	t := a.timeoutTimer
	a.timeoutTimerMutex.Unlock()

	<-t.C
}

func (a *attClient) close() {
	a.cmdmgr.Close()
}

func (a *attClient) handlePDU(method ATTCommand, isAuthenticated bool, buf *pdu.PDU) (bool, error) {
	switch method {
	case ATTMultipleHandleValueNTF:
		fallthrough
	case ATTHandleValueNTF:
		fallthrough
	case ATTHandleValueIND:
		return a.handleNTFIND(method, buf)

	default:
		a.timeoutTimerMutex.Lock()
		a.timeoutTimer.Stop()
		a.timeoutTimerMutex.Unlock()

		keepBuffer := false
		err := a.cmdmgr.IterateActive(func(slot *slotset.Slot) (bool, error) {
			data := slot.Data.(*attClientCmdData)
			data.method = method
			data.buf = buf

			slot.PostWithoutLock(nil)
			keepBuffer = true
			return false, nil
		})
		return keepBuffer, err
	}
}

func (a *attClient) findInformation(ctx context.Context, startingHandle uint16, endingHandle uint16) ([]attstructure.HandleInfo, error) {
	buf := bleutil.GetBuffer(5)
	buf.Buf()[0] = byte(ATTFindInformationReq)
	binary.LittleEndian.PutUint16(buf.Buf()[1:], startingHandle)
	binary.LittleEndian.PutUint16(buf.Buf()[3:], endingHandle)

	cmd, response, err := a.sendCommand(ctx, buf, true)
	defer bleutil.ReleaseBuffer(response)

	if err != nil || cmd != ATTFindInformationRsp {
		return nil, err
	}

	header := response.DropLeft(1)
	if header == nil {
		return nil, ErrorProtocolViolation
	}

	width := 2
	if header[0] == 2 {
		width = 16
	}

	var result []attstructure.HandleInfo

	for {
		record := response.DropLeft(2 + width)
		if record == nil {
			break
		}

		result = append(result, attstructure.HandleInfo{
			UUIDWidth: width,
			Handle:    binary.LittleEndian.Uint16(record),
			UUID:      bleutil.UUIDFromBytes(record[2:]),
		})
	}

	return result, nil
}

func (a *attClient) findInformationAll(ctx context.Context, startingHandle uint16, endingHandle uint16) ([]attstructure.HandleInfo, error) {
	currentHandle := uint16(1)

	var result []attstructure.HandleInfo

main:
	for {
		handles, err := a.findInformation(ctx, currentHandle, 0xFFFF)
		if err != nil || handles == nil {
			return result, err
		}

		for _, m := range handles {
			result = append(result, m)

			/* Make sure the handles keep increasing */
			if m.Handle == 0xFFFF {
				break main
			}
			if m.Handle < currentHandle {
				break main
			}
			currentHandle = m.Handle + 1
		}
	}

	return result, nil
}

func (a *attClient) writeHandle(ctx context.Context, handle uint16, value []byte, withRsp bool) (int, ATTError, error) {
	first := true
retry:
	mtu := a.parent.getMTU()

	if len(value) > mtu-3 {
		value = value[:mtu-3]
	}

	buf := bleutil.GetBuffer(3 + len(value))
	binary.LittleEndian.PutUint16(buf.Buf()[1:], handle)
	copy(buf.Buf()[3:], value)

	if withRsp {
		buf.Buf()[0] = byte(ATTWriteReq)

		_, buf, atterr, err := a.sendCommandErrRsp(ctx, buf)
		bleutil.ReleaseBuffer(buf)

		if first && err == nil && (atterr == ATTErrorInsufficientEncryption || atterr == ATTErrorInsufficientAuthentication) {
			_, err := a.parent.parent.smpConn.GoSecure(ctx, true)
			if err != nil {
				return len(value), atterr, err
			}
			first = false
			goto retry
		}

		return len(value), atterr, err
	}

	buf.Buf()[0] = byte(ATTWriteCMD)
	_, _, err := a.sendCommand(ctx, buf, false)
	return len(value), 0, err
}

func (a *attClient) readHandle(ctx context.Context, handle uint16, result []byte) ([]byte, ATTError, error) {
	first := true

retry:
	buf := bleutil.GetBuffer(3)
	buf.Buf()[0] = byte(ATTReadReq)
	binary.LittleEndian.PutUint16(buf.Buf()[1:], handle)

	cmd, response, atterr, err := a.sendCommandErrRsp(ctx, buf)
	if response != nil {
		result = append(result[:0], response.Buf()...)
		bleutil.ReleaseBuffer(response)
	}

	if first && err == nil && (atterr == ATTErrorInsufficientEncryption || atterr == ATTErrorInsufficientAuthentication) {
		_, err := a.parent.parent.smpConn.GoSecure(ctx, true)
		if err != nil {
			return result, atterr, err
		}
		first = false
		goto retry
	}

	if err != nil || atterr != 0 {
		return nil, atterr, err
	}

	if cmd != ATTReadRsp {
		return nil, 0, ErrorProtocolViolation
	}

	return result, 0, err
}

func (a *attClient) readHandleBlob(ctx context.Context, handle uint16, result []byte) ([]byte, ATTError, error) {
	first := true

retry:
	buf := bleutil.GetBuffer(5)
	buf.Buf()[0] = byte(ATTReadBlobReq)
	binary.LittleEndian.PutUint16(buf.Buf()[1:], handle)
	binary.LittleEndian.PutUint16(buf.Buf()[3:], uint16(len(result)))

	cmd, response, atterr, err := a.sendCommandErrRsp(ctx, buf)
	if response != nil {
		result = append(result, response.Buf()...)
		bleutil.ReleaseBuffer(response)
	}

	if first && err == nil && (atterr == ATTErrorInsufficientEncryption || atterr == ATTErrorInsufficientAuthentication) {
		_, err := a.parent.parent.smpConn.GoSecure(ctx, true)
		if err != nil {
			return result, atterr, err
		}
		first = false
		goto retry
	}

	if err != nil || atterr != 0 {
		return nil, atterr, err
	}

	if cmd != ATTReadBlobRsp {
		return result, 0, ErrorProtocolViolation
	}
	return result, 0, err
}

func (a *attClient) readHandleAll(ctx context.Context, handle uint16, result []byte) ([]byte, ATTError, error) {
	mtu := a.parent.getMTUBlocking()
	result, at, err := a.readHandle(ctx, handle, result)
	if err != nil || at > 0 || (len(result) <= mtu-1) {
		return result, at, err
	}

	for {
		l1 := len(result)
		result, at, err := a.readHandleBlob(ctx, handle, result)
		if err != nil || at > 0 || len(result) == l1 {
			/* If this is not a long element, ignore this error as the read was succesful */
			if at == ATTErrorAttributeNotLong {
				at = 0
			}
			return result, at, err
		}
	}
}

func (a *attClient) readByUUID(ctx context.Context, uuid bleutil.UUID, result []byte) ([]byte, ATTError, error) {
	ub := uuid.UUIDToBytes()

	buf := bleutil.GetBuffer(5 + len(ub))
	buf.Buf()[0] = byte(ATTReadByTypeReq)
	binary.LittleEndian.PutUint16(buf.Buf()[1:], 0x1)
	binary.LittleEndian.PutUint16(buf.Buf()[3:], 0xFF)
	copy(buf.Buf()[4:], ub)

	cmd, response, aterr, err := a.sendCommandErrRsp(ctx, buf)
	defer bleutil.ReleaseBuffer(response)
	if err != nil || aterr != 0 {
		return nil, aterr, err
	}

	header := response.DropLeft(3)

	if cmd != ATTReadByTypeRsp || header == nil {
		return nil, 0, ErrorProtocolViolation
	}

	handle := binary.LittleEndian.Uint16(header[1:])

	return a.readHandleAll(ctx, handle, result)
}

func attErrorToError(atterr ATTError) error {
	if atterr == 0 {
		return nil
	}

	return fmt.Errorf("ATT Error: %d", atterr)
}

func (a *attClient) discoverRemoteDeviceStructure() (*attstructure.Structure, error) {
	/* With a high MTU this goes so much faster */
	a.parent.getMTUBlocking()

	ctx := context.Background()

	handles, err := a.findInformationAll(ctx, 1, 0xFFFF)
	if err != nil {
		return nil, err
	}

	var gattHandles []*attstructure.GATTHandle
	for _, m := range handles {
		handle := &attstructure.GATTHandle{
			Info: m,
		}

		/* If it is descriptive, try to read it */
		if isPartOfGATTDatabase(m.UUID) > 0 {
			value, attErr, err := a.readHandleAll(ctx, m.Handle, nil)
			if attErr > 0 || err != nil {
				return nil, err
			}
			handle.Value = value
		}

		gattHandles = append(gattHandles, handle)
	}

	for _, m := range gattHandles {
		a.parent.logger.WithFields(logrus.Fields{
			"1uuid":   m.Info.UUID,
			"0handle": m.Info.Handle,
			"2data":   hex.EncodeToString(m.Value),
		}).Debug("Discovered characteristic")
	}

	return attstructure.ImportStructure(gattHandles, a.parent.parent.ClientRead, a.parent.parent.ClientWrite)
}
