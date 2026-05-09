package bleatt

import (
	"context"
	"encoding/binary"
	"io"
	"sync"
	"testing"

	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"
)

// fakeConn is a minimal hciconnmgr.BufferConn for ATT-server tests.
type fakeConn struct {
	mu     sync.Mutex
	tx     []*pdu.PDU
	logger *logrus.Entry
	closed bool
}

func newFakeConn() *fakeConn {
	l := logrus.New()
	l.Out = io.Discard
	l.Level = logrus.PanicLevel
	return &fakeConn{logger: logrus.NewEntry(l)}
}

func (f *fakeConn) IsOpen() bool { return !f.closed }
func (f *fakeConn) Close() error { f.closed = true; return nil }
func (f *fakeConn) ReadBuffer(ctx context.Context) (*pdu.PDU, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}
func (f *fakeConn) WriteBuffer(buf *pdu.PDU) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.tx = append(f.tx, buf)
	return nil
}
func (f *fakeConn) GetLogger() *logrus.Entry { return f.logger }
func (f *fakeConn) UseStart()                {}
func (f *fakeConn) UseDone()                 {}

func (f *fakeConn) takeTx() []*pdu.PDU {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := f.tx
	f.tx = nil
	return out
}

// Compile-time check that fakeConn satisfies the BufferConn interface
// expected by bleatt.
var _ hciconnmgr.BufferConn = (*fakeConn)(nil)

// buildTestServer wires up an attServer with a small structure: one
// service (180a, Device Information) with one read-only characteristic
// (2a29, "manufacturer") and one writable characteristic (2a30).
func buildTestServer(t *testing.T) (*attServer, *gattDeviceConn, *fakeConn) {
	t.Helper()

	external := attstructure.NewStructure()
	svc := external.AddPrimaryService(bleutil.UUIDFromStringPanic("180a"))
	svc.AddCharacteristicReadOnly(bleutil.UUIDFromStringPanic("2a29"), []byte("manufacturer"))
	svc.AddCharacteristic(
		bleutil.UUIDFromStringPanic("2a30"),
		attstructure.CharacteristicRead|attstructure.CharacteristicWriteAck,
		attstructure.ValueConfig{LengthMax: 32},
	)

	dev := NewGattDevice(external, &GattDeviceConfig{
		MTU:                     247,
		DeviceName:              "test",
		DiscoverRemoteOnConnect: false,
	})

	conn := newFakeConn()
	gconn := &gattDeviceConn{
		parent: dev,
		conn:   conn,
		logger: bleutil.LogWithPrefix(conn.GetLogger(), "att-test"),
		mtu:    247,
	}
	gconn.client.init(gconn)
	return &dev.server, gconn, conn
}

func makePDU(b ...byte) *pdu.PDU {
	p := &pdu.PDU{}
	p.Append(b...)
	return p
}

// Read the last response code byte from the server's outgoing PDU.
func lastTxOpcode(t *testing.T, conn *fakeConn) byte {
	t.Helper()
	all := conn.takeTx()
	if len(all) == 0 {
		t.Fatal("no PDU sent")
	}
	last := all[len(all)-1]
	if last.Len() < 1 {
		t.Fatal("empty PDU sent")
	}
	return last.Buf()[0]
}

// ExchangeMTU: peer sends a request, server must reply with ExchangeMTURsp.
func TestServerExchangeMTUReq(t *testing.T) {
	srv, conn, fc := buildTestServer(t)

	// Body: 2-byte requested MTU
	body := &pdu.PDU{}
	body.Append(0xF7, 0x00) // requested 247
	keep, err := srv.handleMTUReq(conn, body)
	if err != nil {
		t.Fatal(err)
	}
	if !keep {
		// the implementation may keep the buffer; either is acceptable
	}
	if op := lastTxOpcode(t, fc); op != byte(ATTExchangeMTURsp) {
		t.Errorf("response opcode: got %#x want %#x", op, ATTExchangeMTURsp)
	}
}

// ExchangeMTU with wrong-length body must return ErrorProtocolViolation.
func TestServerExchangeMTUReqWrongLength(t *testing.T) {
	srv, conn, _ := buildTestServer(t)
	if _, err := srv.handleMTUReq(conn, makePDU(0xF7)); err == nil {
		t.Error("expected protocol violation on 1-byte MTU body")
	}
}

// Read of an unknown handle must produce ATTErrorRsp with InvalidHandle.
func TestServerReadReqUnknownHandle(t *testing.T) {
	srv, conn, fc := buildTestServer(t)

	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0xFFFE)
	if _, err := srv.handleReadReq(conn, ATTReadReq, body); err != nil {
		t.Fatal(err)
	}
	tx := fc.takeTx()
	if len(tx) == 0 {
		t.Fatal("no response")
	}
	if got := tx[0].Buf()[0]; got != byte(ATTErrorRsp) {
		t.Errorf("expected ATTErrorRsp, got %#x", got)
	}
}

// Read of a known handle returns its value.
func TestServerReadReqValid(t *testing.T) {
	srv, conn, fc := buildTestServer(t)

	// Find the handle of the 2a29 characteristic *value* (manufacturer).
	var valueHandle uint16
	for _, h := range srv.localStructure.Handles {
		if h.Info.UUID == bleutil.UUIDFromStringPanic("2a29") && string(h.Value) == "manufacturer" {
			valueHandle = h.Info.Handle
			break
		}
	}
	if valueHandle == 0 {
		t.Fatal("could not locate test characteristic")
	}

	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), valueHandle)
	if _, err := srv.handleReadReq(conn, ATTReadReq, body); err != nil {
		t.Fatal(err)
	}
	tx := fc.takeTx()
	if len(tx) == 0 {
		t.Fatal("no response")
	}
	resp := tx[0].Buf()
	if resp[0] != byte(ATTReadRsp) {
		t.Fatalf("expected ATTReadRsp, got %#x", resp[0])
	}
	if string(resp[1:]) != "manufacturer" {
		t.Errorf("payload: got %q want %q", resp[1:], "manufacturer")
	}
}

// Read with body length wrong → ErrorProtocolViolation.
func TestServerReadReqWrongLength(t *testing.T) {
	srv, conn, _ := buildTestServer(t)
	if _, err := srv.handleReadReq(conn, ATTReadReq, makePDU(0xAA)); err == nil {
		t.Error("expected error on 1-byte read body")
	}
}

// Write of an unknown handle: ATTErrorRsp.
func TestServerWriteReqUnknownHandle(t *testing.T) {
	srv, conn, fc := buildTestServer(t)
	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0xFFFE)
	body.Append(0x42)
	srv.handleWriteReq(conn, ATTWriteReq, body)
	if got := lastTxOpcode(t, fc); got != byte(ATTErrorRsp) {
		t.Errorf("expected ATTErrorRsp, got %#x", got)
	}
}

// Write to a writable handle and then read it back.
func TestServerWriteReqRoundTrip(t *testing.T) {
	srv, conn, fc := buildTestServer(t)

	var writeHandle uint16
	for _, h := range srv.localStructure.Handles {
		if h.Info.UUID == bleutil.UUIDFromStringPanic("2a30") && h.Info.Flags&attstructure.CharacteristicWriteAck != 0 {
			writeHandle = h.Info.Handle
			break
		}
	}
	if writeHandle == 0 {
		t.Fatal("could not locate writable characteristic")
	}

	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), writeHandle)
	body.Append('H', 'i')
	if _, err := srv.handleWriteReq(conn, ATTWriteReq, body); err != nil {
		t.Fatal(err)
	}
	if got := lastTxOpcode(t, fc); got != byte(ATTWriteRsp) {
		t.Fatalf("expected ATTWriteRsp, got %#x", got)
	}

	// Read it back.
	body = &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), writeHandle)
	if _, err := srv.handleReadReq(conn, ATTReadReq, body); err != nil {
		t.Fatal(err)
	}
	tx := fc.takeTx()
	if len(tx) == 0 {
		t.Fatal("no response")
	}
	if string(tx[0].Buf()[1:]) != "Hi" {
		t.Errorf("read-back payload: got %q want %q", tx[0].Buf()[1:], "Hi")
	}
}

// ExecuteWrite with an invalid flag byte must produce ATTErrorRsp.
func TestServerExecuteWriteInvalidFlag(t *testing.T) {
	srv, conn, fc := buildTestServer(t)
	if _, err := srv.handleExecuteWriteReq(conn, makePDU(0x42)); err != nil {
		t.Fatal(err)
	}
	if got := lastTxOpcode(t, fc); got != byte(ATTErrorRsp) {
		t.Errorf("expected ATTErrorRsp on invalid flag, got %#x", got)
	}
}

// ReadByGroupTypeReq with start > end must return InvalidHandle error.
func TestServerDiscoveryStartGTEnd(t *testing.T) {
	srv, conn, fc := buildTestServer(t)
	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 100)
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 50)
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0x2800) // PrimaryService UUID
	srv.handleDiscovery(conn, ATTReadByGroupTypeReq, body)
	if got := lastTxOpcode(t, fc); got != byte(ATTErrorRsp) {
		t.Errorf("expected ATTErrorRsp on start>end, got %#x", got)
	}
}

// ReadByGroupTypeReq for primary services succeeds.
func TestServerDiscoveryPrimaryServices(t *testing.T) {
	srv, conn, fc := buildTestServer(t)
	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 1)
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0xFFFF)
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0x2800)
	if _, err := srv.handleDiscovery(conn, ATTReadByGroupTypeReq, body); err != nil {
		t.Fatal(err)
	}
	if got := lastTxOpcode(t, fc); got != byte(ATTReadByGroupTypeRsp) {
		t.Errorf("expected ATTReadByGroupTypeRsp, got %#x", got)
	}
}

// ReadMultiple with odd-length body must produce InvalidPDU error.
func TestServerReadMultipleOddLength(t *testing.T) {
	srv, conn, fc := buildTestServer(t)
	srv.handleReadReqMultiple(conn, ATTReadMultipleReq, makePDU(0x01, 0x00, 0x02))
	if got := lastTxOpcode(t, fc); got != byte(ATTErrorRsp) {
		t.Errorf("expected ATTErrorRsp on odd-length body, got %#x", got)
	}
}

// ReadMultiple with empty body → InvalidPDU.
func TestServerReadMultipleEmpty(t *testing.T) {
	srv, conn, fc := buildTestServer(t)
	srv.handleReadReqMultiple(conn, ATTReadMultipleReq, makePDU())
	if got := lastTxOpcode(t, fc); got != byte(ATTErrorRsp) {
		t.Errorf("expected ATTErrorRsp on empty body, got %#x", got)
	}
}

// PrepareWrite to an unknown handle.
func TestServerPrepareWriteUnknown(t *testing.T) {
	srv, conn, fc := buildTestServer(t)
	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0xFFFE) // handle
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0)     // offset
	body.Append('x', 'y')
	srv.handlePrepateWriteReq(conn, body)
	if got := lastTxOpcode(t, fc); got != byte(ATTErrorRsp) {
		t.Errorf("expected ATTErrorRsp, got %#x", got)
	}
}

// PrepareWrite + ExecuteWrite (commit) round-trip.
func TestServerPrepareExecuteCommit(t *testing.T) {
	srv, conn, fc := buildTestServer(t)

	var writeHandle uint16
	for _, h := range srv.localStructure.Handles {
		if h.Info.UUID == bleutil.UUIDFromStringPanic("2a30") && h.Info.Flags&attstructure.CharacteristicWriteAck != 0 {
			writeHandle = h.Info.Handle
			break
		}
	}
	if writeHandle == 0 {
		t.Fatal("no writable characteristic")
	}

	// Two prepare-writes covering different offsets.
	for i, payload := range [][]byte{[]byte("AB"), []byte("CD")} {
		body := &pdu.PDU{}
		binary.LittleEndian.PutUint16(body.ExtendRight(2), writeHandle)
		binary.LittleEndian.PutUint16(body.ExtendRight(2), uint16(i*2)) // offset 0, then 2
		body.Append(payload...)
		if _, err := srv.handlePrepateWriteReq(conn, body); err != nil {
			t.Fatal(err)
		}
		if got := lastTxOpcode(t, fc); got != byte(ATTPrepareWriteRsp) {
			t.Fatalf("prepare %d: got opcode %#x", i, got)
		}
	}

	// Execute (commit).
	if _, err := srv.handleExecuteWriteReq(conn, makePDU(0x01)); err != nil {
		t.Fatal(err)
	}
	if got := lastTxOpcode(t, fc); got != byte(ATTExecuteWriteRsp) {
		t.Fatalf("execute commit: got opcode %#x", got)
	}

	// Read back: the value should now be "ABCD".
	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), writeHandle)
	if _, err := srv.handleReadReq(conn, ATTReadReq, body); err != nil {
		t.Fatal(err)
	}
	tx := fc.takeTx()
	if len(tx) == 0 {
		t.Fatal("no read response")
	}
	if string(tx[0].Buf()[1:]) != "ABCD" {
		t.Errorf("committed value: got %q want %q", tx[0].Buf()[1:], "ABCD")
	}
}

// PrepareWrite + ExecuteWrite (cancel) discards the queue.
func TestServerPrepareExecuteCancel(t *testing.T) {
	srv, conn, fc := buildTestServer(t)

	var writeHandle uint16
	for _, h := range srv.localStructure.Handles {
		if h.Info.UUID == bleutil.UUIDFromStringPanic("2a30") && h.Info.Flags&attstructure.CharacteristicWriteAck != 0 {
			writeHandle = h.Info.Handle
			break
		}
	}

	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), writeHandle)
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0)
	body.Append('X', 'Y')
	srv.handlePrepateWriteReq(conn, body)
	fc.takeTx() // discard ATTPrepareWriteRsp

	// Cancel.
	if _, err := srv.handleExecuteWriteReq(conn, makePDU(0x00)); err != nil {
		t.Fatal(err)
	}
	if got := lastTxOpcode(t, fc); got != byte(ATTExecuteWriteRsp) {
		t.Errorf("execute cancel: got opcode %#x", got)
	}

	// Value must be unchanged.
	body = &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), writeHandle)
	srv.handleReadReq(conn, ATTReadReq, body)
	tx := fc.takeTx()
	if string(tx[0].Buf()[1:]) == "XY" {
		t.Error("cancel should have discarded the prepared write")
	}
}

// PrepareWrite past LengthMax must be rejected per-fragment.
func TestServerPrepareWriteOverLengthMax(t *testing.T) {
	srv, conn, fc := buildTestServer(t)

	var writeHandle uint16
	for _, h := range srv.localStructure.Handles {
		if h.Info.UUID == bleutil.UUIDFromStringPanic("2a30") && h.Info.Flags&attstructure.CharacteristicWriteAck != 0 {
			writeHandle = h.Info.Handle
			break
		}
	}

	// LengthMax for this characteristic is 32 (set in buildTestServer).
	// Try a fragment at offset 30 with 5 bytes → exceeds 32.
	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), writeHandle)
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 30)
	body.Append('A', 'B', 'C', 'D', 'E')
	srv.handlePrepateWriteReq(conn, body)
	if got := lastTxOpcode(t, fc); got != byte(ATTErrorRsp) {
		t.Errorf("over-LengthMax prepare: got opcode %#x want ATTErrorRsp", got)
	}
}

// FindInformation discovery returns handle/UUID pairs.
func TestServerFindInformation(t *testing.T) {
	srv, conn, fc := buildTestServer(t)

	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 1)
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0xFFFF)
	if _, err := srv.handleDiscovery(conn, ATTFindInformationReq, body); err != nil {
		t.Fatal(err)
	}
	if got := lastTxOpcode(t, fc); got != byte(ATTFindInformationRsp) {
		t.Errorf("FindInformation: got opcode %#x", got)
	}
}

// ReadByType for the Characteristic UUID returns the characteristic
// declarations within a service range.
func TestServerReadByType(t *testing.T) {
	srv, conn, fc := buildTestServer(t)

	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 1)
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0xFFFF)
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0x2803) // Characteristic UUID
	if _, err := srv.handleDiscovery(conn, ATTReadByTypeReq, body); err != nil {
		t.Fatal(err)
	}
	if got := lastTxOpcode(t, fc); got != byte(ATTReadByTypeRsp) {
		t.Errorf("ReadByType: got opcode %#x", got)
	}
}

// FindByTypeValueReq with a value that's longer than the maxFindByTypeValueLength
// cap must be rejected with InvalidPDU.
func TestServerFindByTypeValueOverlongValue(t *testing.T) {
	srv, conn, fc := buildTestServer(t)

	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 1)
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0xFFFF)
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0x2800)
	// Large checkValue (>maxFindByTypeValueLength)
	body.Append(make([]byte, maxFindByTypeValueLength+1)...)
	srv.handleDiscovery(conn, ATTFindByTypeValueReq, body)
	if got := lastTxOpcode(t, fc); got != byte(ATTErrorRsp) {
		t.Errorf("over-length checkValue: got opcode %#x want ATTErrorRsp", got)
	}
}

// FindByTypeValue with the spec's PrimaryService discovery: looking for
// the 180a service by UUID value.
func TestServerFindByTypeValueDiscoversService(t *testing.T) {
	srv, conn, fc := buildTestServer(t)
	body := &pdu.PDU{}
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 1)      // start
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0xFFFF) // end
	binary.LittleEndian.PutUint16(body.ExtendRight(2), 0x2800) // PrimaryService UUID
	body.Append(0x0a, 0x18)                                     // value: 180a in little-endian
	if _, err := srv.handleDiscovery(conn, ATTFindByTypeValueReq, body); err != nil {
		t.Fatal(err)
	}
	if got := lastTxOpcode(t, fc); got != byte(ATTFindByTypeValueRsp) {
		t.Errorf("expected ATTFindByTypeValueRsp, got %#x", got)
	}
}
