package bleatt

import (
	"context"
	"encoding/binary"
	"encoding/hex"
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

	<-a.timeoutTimer.C
	go a.handleTimeout()

	return nil
}

func (a *attClient) sendCommand(ctx context.Context, cmd *pdu.PDU) (ATTCommand, *pdu.PDU, error) {
	slot, err := a.cmdmgr.Get(ctx)
	if err != nil {
		return 0, nil, err
	}
	defer a.cmdmgr.Put(slot)
	slot.Activate()

	err = a.write(cmd, true)
	if err != nil {
		return 0, nil, err
	}

	_, err = slot.WaitCtx(ctx)
	if err != nil {
		return 0, nil, err
	}

	slot.Deactivate()

	data := slot.Data.(*attClientCmdData)
	return data.method, data.buf, nil
}

func (a *attClient) sendCommandErrRsp(ctx context.Context, req *pdu.PDU) (ATTCommand, *pdu.PDU, ATTError, error) {
	method := req.Buf()[0]
	cmd, response, err := a.sendCommand(ctx, req)
	if cmd == ATTErrorRsp {
		data := response.DropLeft(4)
		defer bleutil.ReleaseBuffer(response)

		if data == nil || data[0] != method {
			return cmd, nil, 0, ErrorProtocolViolation
		}

		return cmd, nil, ATTError(data[3]), nil
	}

	return cmd, response, 0, err
}

func (a *attClient) write(buf *pdu.PDU, expectReply bool) error {
	if a.parent.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		a.parent.logger.WithFields(logrus.Fields{
			"0buf":         buf,
			"1expectReply": expectReply,
		}).Debug("ATT Client Write")
	}

	if expectReply {
		a.timeoutTimerMutex.Lock()
		a.timeoutTimer.Reset(10 * time.Second)
		a.timeoutTimerMutex.Unlock()
	}

	return a.parent.conn.WriteBuffer(buf)
}

func (a *attClient) handleNTFIND(method ATTCommand, buf *pdu.PDU) (bool, error) {
	if method == ATTMultipleHandleValueNTF {
		for data := buf.DropLeft(4); data != nil; data = buf.DropLeft(4) {
			handle := binary.LittleEndian.Uint16(data)
			// The other 2 bytes contain the length of the data. It is not important for us

			if a.parent.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
				a.parent.logger.WithFields(logrus.Fields{
					"0handle": handle,
				}).Trace("Got mutli notification")
			}

		}

		return false, nil
	}

	handleBuf := buf.DropLeft(2)
	if handleBuf != nil {
		handle := binary.LittleEndian.Uint16(handleBuf)
		isIndication := method == ATTHandleValueIND

		if a.parent.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
			a.parent.logger.WithFields(logrus.Fields{
				"0handle":       handle,
				"1isIndication": isIndication,
				"2data":         buf,
			}).Trace("Got notification")
		}

		return true, nil
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

func (a *attClient) handlePDU(method ATTCommand, isAuthenticated bool, buf *pdu.PDU) (bool, error) {
	switch method {
	case ATTMultipleHandleValueNTF:
		return a.handleNTFIND(method, buf)
	case ATTHandleValueNTF:
		return a.handleNTFIND(method, buf)

	case ATTHandleValueIND:
		keepBuffer, err := a.handleNTFIND(method, buf)
		if err != nil {
			return false, err
		}

		/* Confirm notification */
		rsp := bleutil.GetBuffer(1)
		rsp.Buf()[0] = byte(ATTHandleValueCNF)
		return keepBuffer, a.write(rsp, false)
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

	cmd, response, err := a.sendCommand(ctx, buf)
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

func (a *attClient) readHandle(ctx context.Context, handle uint16, result []byte) ([]byte, ATTError, error) {
	buf := bleutil.GetBuffer(3)
	buf.Buf()[0] = byte(ATTReadReq)
	binary.LittleEndian.PutUint16(buf.Buf()[1:], handle)

	cmd, response, aterr, err := a.sendCommandErrRsp(ctx, buf)
	defer bleutil.ReleaseBuffer(response)

	if err != nil || aterr != 0 {
		return nil, aterr, err
	}

	if cmd != ATTReadRsp {
		return nil, 0, ErrorProtocolViolation
	}

	result = append(result[:0], response.Buf()...)
	return result, 0, err
}

func (a *attClient) readHandleBlob(ctx context.Context, handle uint16, result []byte) ([]byte, ATTError, error) {
	buf := bleutil.GetBuffer(5)
	buf.Buf()[0] = byte(ATTReadBlobReq)
	binary.LittleEndian.PutUint16(buf.Buf()[1:], handle)
	binary.LittleEndian.PutUint16(buf.Buf()[3:], uint16(len(result)))

	cmd, response, aterr, err := a.sendCommandErrRsp(ctx, buf)
	defer bleutil.ReleaseBuffer(response)

	if err != nil || aterr != 0 {
		return nil, aterr, err
	}

	if cmd != ATTReadBlobRsp {
		return result, 0, ErrorProtocolViolation
	}
	result = append(result, response.Buf()...)
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

func (a *attClient) discoverRemoteDeviceStructure() {
	/* This is the same for all clients on one device */
	a.parent.parent.remoteDiscoverOnce.Do(func() {
		/* With a high MTU this goes so much faster */
		a.parent.getMTUBlocking()

		ctx := context.Background()

		handles, err := a.findInformationAll(ctx, 1, 0xFFFF)
		if err != nil {
			return
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
					return
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

		//TODO: Removed for quick push (this code is not releasable yet)
		//a.parent.notifyStructure(gattHandles)
	})
}
