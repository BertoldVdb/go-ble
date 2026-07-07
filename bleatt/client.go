package bleatt

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/BertoldVdb/go-misc/slotset"
	"github.com/sirupsen/logrus"
)

type attClientCmdData struct {
	method   ATTCommand
	buf      *pdu.PDU
	expected ATTCommand // expected response opcode (request opcode + 1)
}

type attClient struct {
	ctxExpired context.Context

	parent *gattDeviceConn
	cmdmgr *slotset.SlotSet

	// 	dbHashRequest sync.Once
	// 	dbHash        []byte

	timeoutTimerMutex sync.Mutex
	timeoutTimer      *time.Timer

	// preDiscoveryIND captures the handles of any indications received
	// before discoverRemoteDeviceStructure has settled. The cache
	// validation path consults this set so that a Service Changed
	// indication that fires immediately on encrypted reconnect — before
	// our Read by Type validation completes — still triggers eviction.
	// preDiscoveryINDOpen flips to false once discovery completes; from
	// that point inbound indications dispatch through the structure
	// normally and the map is no longer touched. Bounded at a small
	// size to keep a chatty peer from growing it without bound.
	preDiscoveryINDMutex sync.Mutex
	preDiscoveryIND      map[uint16]struct{}
	preDiscoveryINDOpen  bool

	// notifyDebugSeq is a per-connection monotonic counter for the
	// "Got notification" debug log. Incremented only when Debug is
	// enabled, so non-debug builds pay nothing. Lets the operator
	// distinguish "peer spammed N distinct PDUs in a tight loop" from
	// "library re-delivered the same buffer" when paired with the
	// logged PDU pointer.
	notifyDebugSeq uint64
}

const preDiscoveryINDMaxEntries = 32

func (a *attClient) init(parent *gattDeviceConn) error {
	*a = attClient{
		parent: parent,

		cmdmgr: slotset.New(1, func(slot *slotset.Slot) {
			slot.Data = &attClientCmdData{}
		}),

		timeoutTimer: time.AfterFunc(time.Hour, func() {
			parent.parent.CloseConn(parent.conn)
		}),

		preDiscoveryIND:     make(map[uint16]struct{}),
		preDiscoveryINDOpen: true,
	}

	a.timeoutTimer.Stop()

	ctxExpired, cancel := context.WithCancel(context.Background())
	cancel()

	a.ctxExpired = ctxExpired

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

	/* Stamp the expected response opcode (req+1) so handlePDU can
	   discard a stale response that arrives after a previous request
	   timed out — without this check, writeHandle/findInformation/...
	   silently treat any non-error response as success. */
	expected := ATTCommand(0)
	if cmd.Len() > 0 {
		expected = ATTCommand(cmd.Buf()[0]) + 1
	}
	slot.Data.(*attClientCmdData).expected = expected

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
			if data == nil {
				/* Peer declared a length larger than the buffer; the cursor
				   was not advanced, so there is no way to recover the parser.
				   Stop processing this PDU. */
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
			seq := atomic.AddUint64(&a.notifyDebugSeq, 1)
			a.parent.logger.WithFields(logrus.Fields{
				"0handle":       handle,
				"1isIndication": isIndication,
				"2data":         buf,
				"3pdu":          fmt.Sprintf("%p", buf),
				"4seq":          seq,
			}).Debug("Got notification")
		}

		if isIndication {
			a.recordPreDiscoveryIndication(handle)
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

func (a *attClient) close() {
	a.timeoutTimerMutex.Lock()
	a.timeoutTimer.Stop()
	a.timeoutTimerMutex.Unlock()

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
			/* Drop stale responses that don't correspond to the active
			   request. ATTErrorRsp is always allowed (peer can refuse
			   any request); otherwise the response opcode must be the
			   expected one. */
			if data.expected != 0 && method != ATTErrorRsp && method != data.expected {
				if a.parent.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					a.parent.logger.WithFields(logrus.Fields{
						"0got":      method,
						"1expected": data.expected,
					}).Debug("Dropping ATT response with unexpected opcode")
				}
				return false, nil
			}
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
	first := true
retry:

	buf := bleutil.GetBuffer(5)
	buf.Buf()[0] = byte(ATTFindInformationReq)
	binary.LittleEndian.PutUint16(buf.Buf()[1:], startingHandle)
	binary.LittleEndian.PutUint16(buf.Buf()[3:], endingHandle)

	_, response, atterr, err := a.sendCommandErrRsp(ctx, buf)

	if first && err == nil && (atterr == ATTErrorInsufficientEncryption || atterr == ATTErrorInsufficientAuthentication) {
		_, err := a.parent.parent.smpConn.GoSecure(ctx, true)
		if err != nil {
			return nil, err
		}
		first = false

		goto retry
	}

	if err != nil || response == nil {
		return nil, err
	}

	defer bleutil.ReleaseBuffer(response)

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

// recordPreDiscoveryIndication is called from handleNTFIND for every
// inbound indication. It captures the handle into a small set as long
// as discovery is still in progress; once discovery completes the set
// is no longer touched and stops growing. The cache validation path
// reads this set after Read by Type SC settles, so a Service Changed
// indication that fires before validation completes still triggers
// eviction.
func (a *attClient) recordPreDiscoveryIndication(handle uint16) {
	a.preDiscoveryINDMutex.Lock()
	defer a.preDiscoveryINDMutex.Unlock()
	if !a.preDiscoveryINDOpen {
		return
	}
	if len(a.preDiscoveryIND) >= preDiscoveryINDMaxEntries {
		return
	}
	a.preDiscoveryIND[handle] = struct{}{}
}

func (a *attClient) closePreDiscoveryWindow() {
	a.preDiscoveryINDMutex.Lock()
	a.preDiscoveryINDOpen = false
	a.preDiscoveryIND = nil
	a.preDiscoveryINDMutex.Unlock()
}

func (a *attClient) sawPreDiscoveryIndication(handle uint16) bool {
	if handle == 0 {
		return false
	}
	a.preDiscoveryINDMutex.Lock()
	defer a.preDiscoveryINDMutex.Unlock()
	_, ok := a.preDiscoveryIND[handle]
	return ok
}

// findGATTServiceRange walks the handle list (in handle order) and
// returns the inclusive handle range of the GATT Service (UUID 0x1801),
// or (0, 0) if the peer does not publish one. It identifies the service
// by matching UUIDPrimaryService entries and reading the service-UUID
// payload, then scanning forward to the next service declaration.
func findGATTServiceRange(handles []*attstructure.GATTHandle) (uint16, uint16) {
	const gattServiceUUID16 = 0x1801

	startIdx := -1
	for i, h := range handles {
		if h.Info.UUID != attstructure.UUIDPrimaryService {
			continue
		}
		if len(h.Value) != 2 {
			continue
		}
		uuid16 := binary.LittleEndian.Uint16(h.Value)
		if uuid16 == gattServiceUUID16 {
			startIdx = i
			break
		}
	}
	if startIdx < 0 {
		return 0, 0
	}

	start := handles[startIdx].Info.Handle
	end := uint16(0xFFFF)
	for i := startIdx + 1; i < len(handles); i++ {
		u := handles[i].Info.UUID
		if u == attstructure.UUIDPrimaryService || u == attstructure.UUIDSecondaryService {
			end = handles[i].Info.Handle - 1
			break
		}
	}
	return start, end
}

// findServiceChangedHandle scans the cached attribute list for the
// Service Changed value handle (UUID 0x2A05) within the GATT Service
// range. Returns 0 if not present.
func findServiceChangedHandle(handles []*attstructure.GATTHandle, start, end uint16) uint16 {
	scUUID := bleutil.UUIDFromStringPanic("2a05")
	if start == 0 && end == 0 {
		return 0
	}
	for _, h := range handles {
		if h.Info.Handle < start || h.Info.Handle > end {
			continue
		}
		if h.Info.UUID == scUUID {
			return h.Info.Handle
		}
	}
	return 0
}

// validateServiceChanged issues a Read by Type request for UUID 0x2A05
// over the cached GATT Service range and returns true when the response
// reports the same handle the cache recorded. A mismatch (or any error
// short of a clean ATTErrorAttributeNotFound, which is a definitive
// "the peer's GATT Service no longer hosts Service Changed") signals
// the cache must be evicted. ATTErrorAttributeNotFound is treated the
// same way: if the cache said SC was at handle X and the peer no
// longer has it, the database has changed.
func (a *attClient) validateServiceChanged(ctx context.Context, c *CachedGATT) bool {
	scUUID := bleutil.UUIDFromStringPanic("2a05")
	ub := scUUID.UUIDToBytes()

	buf := bleutil.GetBuffer(5 + len(ub))
	buf.Buf()[0] = byte(ATTReadByTypeReq)
	binary.LittleEndian.PutUint16(buf.Buf()[1:], c.GATTServiceStart)
	binary.LittleEndian.PutUint16(buf.Buf()[3:], c.GATTServiceEnd)
	copy(buf.Buf()[4:], ub)

	cmd, response, aterr, err := a.sendCommandErrRsp(ctx, buf)
	defer bleutil.ReleaseBuffer(response)
	if err != nil {
		return false
	}
	if aterr != 0 || cmd != ATTReadByTypeRsp {
		return false
	}

	header := response.DropLeft(3)
	if header == nil {
		return false
	}
	gotHandle := binary.LittleEndian.Uint16(header[1:])
	return gotHandle == c.ServiceChangedHandle
}

func (a *attClient) discoverRemoteDeviceStructure() (*attstructure.Structure, error) {
	/* With a high MTU this goes so much faster */
	a.parent.getMTUBlocking()

	defer a.closePreDiscoveryWindow()

	cfg := a.parent.parent.config

	var (
		gattHandles []*attstructure.GATTHandle
		cached      *CachedGATT
		fromCache   bool
	)

	if cfg.DiscoveryCacheGet != nil {
		cached = cfg.DiscoveryCacheGet(a.parent.parent)
	}

	if cached != nil && len(cached.Handles) > 0 {
		ctx := context.Background()

		valid := true
		if cached.ServiceChangedHandle != 0 {
			if !a.validateServiceChanged(ctx, cached) {
				a.parent.logger.Debug("GATT cache invalid: Service Changed handle moved or missing")
				valid = false
			} else if a.sawPreDiscoveryIndication(cached.ServiceChangedHandle) {
				a.parent.logger.Debug("GATT cache invalid: Service Changed indication observed during validation")
				valid = false
			}
		}

		if valid {
			gattHandles = cached.Handles
			fromCache = true
		} else {
			if cfg.DiscoveryCacheSet != nil {
				cfg.DiscoveryCacheSet(a.parent.parent, nil)
			}
			cached = nil
		}
	}

	if !fromCache {
		ctx := context.Background()

		handles, err := a.findInformationAll(ctx, 1, 0xFFFF)
		if err != nil {
			return nil, err
		}

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
	}

	cacheStr := ""
	if fromCache {
		cacheStr = " (cached)"
	}
	for _, m := range gattHandles {
		a.parent.logger.WithFields(logrus.Fields{
			"1uuid":   m.Info.UUID,
			"0handle": m.Info.Handle,
			"2data":   hex.EncodeToString(m.Value),
		}).Debug("Discovered characteristic" + cacheStr)
	}

	result, err := attstructure.ImportStructure(gattHandles, a.parent.parent.ClientRead, a.parent.parent.ClientWrite)
	if err != nil {
		return result, err
	}

	/* SkipCCCDWrite is a per-peer assertion the caller persists in the
	   cache. On a fresh discovery (no cache, or post-eviction) the flag
	   defaults to false — a freshly bonded peer needs the first CCCD
	   write to actually arrive on the wire. The caller's
	   DiscoveryCacheSet implementation decides what value to store
	   going forward; if they set BondedCCCDPersistent=true, subsequent
	   reconnects will hit the cache hot path and skip the write. */
	if fromCache && cached != nil && cached.BondedCCCDPersistent {
		result.SkipCCCDWrite = true
	}

	if cfg.DiscoveryCacheSet != nil {
		next := &CachedGATT{
			Handles: gattHandles,
		}
		next.GATTServiceStart, next.GATTServiceEnd = findGATTServiceRange(gattHandles)
		next.ServiceChangedHandle = findServiceChangedHandle(gattHandles, next.GATTServiceStart, next.GATTServiceEnd)
		if cached != nil {
			/* Carry over caller-set policy across rediscovery — losing
			   BondedCCCDPersistent because the database changed shape
			   would needlessly demote a peer the user already trusts. */
			next.BondedCCCDPersistent = cached.BondedCCCDPersistent
		}
		cfg.DiscoveryCacheSet(a.parent.parent, next)
	}

	return result, nil
}
