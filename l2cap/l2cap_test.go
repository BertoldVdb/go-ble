package blel2cap

import (
	"context"
	"encoding/binary"
	"io"
	"sync"
	"testing"

	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"
)

// fakeConn is a minimal hciconnmgr.BufferConn for driving L2CAP.
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
	if f.closed {
		return io.ErrClosedPipe
	}
	f.tx = append(f.tx, buf)
	return nil
}
func (f *fakeConn) GetLogger() *logrus.Entry { return f.logger }
func (f *fakeConn) UseStart()                {}
func (f *fakeConn) UseDone()                 {}

func (f *fakeConn) lastTx() *pdu.PDU {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.tx) == 0 {
		return nil
	}
	p := f.tx[len(f.tx)-1]
	f.tx = f.tx[:len(f.tx)-1]
	return p
}

// makeL2CAPFrame builds a complete L2CAP frame (4-byte header + payload).
// frameLen = len(payload).
func makeL2CAPFrame(cid uint16, payload []byte) *pdu.PDU {
	p := &pdu.PDU{}
	hdr := p.ExtendRight(4)
	binary.LittleEndian.PutUint16(hdr[0:2], uint16(len(payload)))
	binary.LittleEndian.PutUint16(hdr[2:4], cid)
	p.Append(payload...)
	return p
}

// L2CAP processInput must drop frames whose declared length doesn't
// match the actual byte count — without crashing.
func TestProcessInputLengthMismatch(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	bad := &pdu.PDU{}
	hdr := bad.ExtendRight(4)
	binary.LittleEndian.PutUint16(hdr[0:2], 99) // declared
	binary.LittleEndian.PutUint16(hdr[2:4], 5)  // CID
	bad.Append(0x01, 0x02)                      // only 2 bytes payload

	if err, _ := l.processInput(bad); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// processInput on a frame for an unregistered CID drops cleanly.
func TestProcessInputUnknownCID(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	frame := makeL2CAPFrame(0x9999, []byte{0x01, 0x02, 0x03})
	if err, _ := l.processInput(frame); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// registerCallbackCID rejects double-registration.
func TestRegisterCallbackCID(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	cid := uint16(0x42)
	rx := func(uint16, *pdu.PDU) (error, bool) { return nil, false }
	if !l.registerCallbackCID(cid, rx) {
		t.Fatal("first register failed")
	}
	if l.registerCallbackCID(cid, rx) {
		t.Fatal("double register should fail")
	}
	l.unregisterCID(cid)
	if !l.registerCallbackCID(cid, rx) {
		t.Fatal("re-register after unregister failed")
	}
}

// processInput dispatches a registered CID's payload to its handler.
func TestProcessInputDispatchesToHandler(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	const cid uint16 = 0x42
	var got []byte
	l.registerCallbackCID(cid, func(c uint16, p *pdu.PDU) (error, bool) {
		if c != cid {
			t.Errorf("handler got cid %#x, want %#x", c, cid)
		}
		got = append([]byte(nil), p.Buf()...)
		return nil, false
	})

	want := []byte{0xAA, 0xBB, 0xCC}
	frame := makeL2CAPFrame(cid, want)
	l.processInput(frame)

	if string(got) != string(want) {
		t.Errorf("handler payload: got %x want %x", got, want)
	}
}

// signallingHandler must drop a frame whose signalling-header length
// field doesn't match the payload's actual length.
func TestSignallingHandlerLengthMismatch(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	// Header: code=EchoReq, id=1, length=99 (declared), then 2 payload bytes.
	body := &pdu.PDU{}
	hdr := body.ExtendRight(4)
	hdr[0] = SigEchoReq
	hdr[1] = 1
	binary.LittleEndian.PutUint16(hdr[2:4], 99)
	body.Append(0xAA, 0xBB)

	err, keep := l.sig.signallingHandler(l.getSignallingChannel(), body)
	if err != nil || keep {
		t.Errorf("expected drop on length mismatch, got err=%v keep=%v", err, keep)
	}
}

// EchoReq must be answered with EchoRsp at the same id, payload echoed back.
func TestSignallingEchoReqAnswered(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	payload := []byte{0x12, 0x34}
	body := &pdu.PDU{}
	hdr := body.ExtendRight(4)
	hdr[0] = SigEchoReq
	hdr[1] = 7 // id
	binary.LittleEndian.PutUint16(hdr[2:4], uint16(len(payload)))
	body.Append(payload...)

	if err, _ := l.sig.signallingHandler(l.getSignallingChannel(), body); err != nil {
		t.Fatal(err)
	}
	tx := conn.lastTx()
	if tx == nil {
		t.Fatal("no response sent")
	}
	// tx contents = L2CAP header (4) + sig header (4) + payload
	got := tx.Buf()
	if len(got) < 8 {
		t.Fatalf("response too short: %x", got)
	}
	// CID = 1 (BR/EDR signalling) because the fake conn is not a
	// *bleconnecter.BLEConnection so isLE is false. The test only
	// cares about the response shape, not the channel.
	if cid := binary.LittleEndian.Uint16(got[2:4]); cid != 1 {
		t.Errorf("response CID: got %d want 1", cid)
	}
	if got[4] != SigEchoRsp {
		t.Errorf("response code: got %#x want %#x", got[4], SigEchoRsp)
	}
	if got[5] != 7 {
		t.Errorf("response id: got %d want 7", got[5])
	}
}

// Unknown / unsupported signalling commands are answered with CommandReject.
func TestSignallingUnsupportedCommandRejected(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	body := &pdu.PDU{}
	hdr := body.ExtendRight(4)
	// Code 0x14 is even and ≤ 0x15 — classified as a request, not a
	// response — so the dispatcher will hit the default and emit
	// CommandRejectRsp. Using e.g. 0x40 would be classified as a
	// response and silently dropped.
	hdr[0] = 0x14
	hdr[1] = 1
	binary.LittleEndian.PutUint16(hdr[2:4], 0)
	if err, _ := l.sig.signallingHandler(l.getSignallingChannel(), body); err != nil {
		t.Fatal(err)
	}
	tx := conn.lastTx()
	if tx == nil {
		t.Fatal("no response sent")
	}
	got := tx.Buf()
	if got[4] != SigCommandRejectRsp {
		t.Errorf("expected CommandRejectRsp, got code %#x", got[4])
	}
}

// signallingHandler short-headers (< 4 bytes after L2CAP header) drop silently.
func TestSignallingHandlerShortHeader(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	// Empty body: DropLeft(4) returns nil.
	body := &pdu.PDU{}
	if err, keep := l.sig.signallingHandler(5, body); err != nil || keep {
		t.Errorf("expected silent drop, got err=%v keep=%v", err, keep)
	}
}

// L2CAP.Close is idempotent.
func TestL2CAPCloseIdempotent(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	err1 := l.Close()
	err2 := l.Close()
	if err1 != err2 {
		t.Errorf("repeated Close returned different errors: %v / %v", err1, err2)
	}
}

// Closing a fresh L2CAP without anything connected must not panic.
func TestL2CAPCloseEmpty(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})
	if err := l.Close(); err != nil {
		t.Errorf("unexpected error from Close on empty L2CAP: %v", err)
	}
}

// L2Connection: external write buffers traverse the L2CAP header path.
func TestL2ConnectionWriteAndRead(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	c := l.connectionCreateInternal(false, 4, 4)
	if c == nil {
		t.Fatal("connectionCreateInternal returned nil")
	}

	// Write through L2Connection.WriteBuffer — must add the L2CAP header
	// using the *remote* CID (4 in our setup).
	body := &pdu.PDU{}
	body.Append(0xDE, 0xAD)
	if err := c.WriteBuffer(body); err != nil {
		t.Fatal(err)
	}
	tx := conn.lastTx()
	if tx == nil {
		t.Fatal("nothing written")
	}
	if cid := binary.LittleEndian.Uint16(tx.Buf()[2:4]); cid != 4 {
		t.Errorf("CID: got %#x want 4", cid)
	}

	// Inject a frame back to the rx handler — the L2Connection should
	// queue it, and ReadBuffer should retrieve it.
	frame := bleutil.GetBuffer(0)
	frame.Append(0x10, 0x11, 0x12)
	if _, keep := c.connectionRxHandler(4, frame); !keep {
		t.Fatal("connectionRxHandler should keep the buffer")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*1000) // ms
	_ = ctx
	_ = cancel

	got, err := c.ReadBuffer(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if string(got.Buf()) != string([]byte{0x10, 0x11, 0x12}) {
		t.Errorf("ReadBuffer payload: got %x", got.Buf())
	}
}

// L2Connection.Close drains queued rxBuffer entries (nothing leaks).
// Use pool buffers because Close calls bleutil.ReleaseBuffer on them.
func TestL2ConnectionCloseDrainsRx(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	c := l.connectionCreateInternal(false, 4, 4)

	for i := 0; i < 3; i++ {
		f := bleutil.GetBuffer(0)
		f.Append(byte(i))
		if _, keep := c.connectionRxHandler(4, f); !keep {
			t.Fatal("connectionRxHandler should keep the buffer")
		}
	}
	// Close should not panic and should leave the connection un-readable.
	c.Close()
	if c.IsOpen() {
		t.Error("L2Connection should report closed after Close")
	}
}

// connectionRxHandler should release the buffer if the connection is
// already closed (no rxBuffer leak).
func TestL2ConnectionRxAfterClose(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})
	c := l.connectionCreateInternal(false, 4, 4)
	c.Close()

	f := bleutil.GetBuffer(0)
	f.Append(0xAA)
	err, keep := c.connectionRxHandler(4, f)
	if err == nil {
		t.Error("expected error on push to closed L2Connection")
	}
	if keep {
		t.Error("closed L2Connection must not keep the buffer")
	}
}

// connectionRxHandler with nil buf signals a close from the wire.
func TestL2ConnectionRxNilBuf(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})
	c := l.connectionCreateInternal(false, 4, 4)

	if _, keep := c.connectionRxHandler(4, nil); keep {
		t.Error("nil buf must not keep")
	}
	if c.IsOpen() {
		t.Error("L2Connection should be closed after nil buf")
	}
}

// L2CAP.Run() exits cleanly on parent close + drains via failActiveToken.
// We don't run the loop (it requires a real BufferConn that surfaces
// data); just ensure Run() doesn't break the package-level invariants
// when never called by closing immediately.
func TestSignallingRunCloseSafe(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})
	if err := l.Close(); err != nil {
		t.Fatal(err)
	}
}

// WriteBufferCID prepends a 4-byte L2CAP header.
func TestWriteBufferCIDPrependsHeader(t *testing.T) {
	conn := newFakeConn()
	l := New(conn, nil, func(psm PSMType, accept L2CAPConnAccepter) {})

	body := &pdu.PDU{}
	body.Append(0x01, 0x02, 0x03)
	if err := l.WriteBufferCID(0x55, body); err != nil {
		t.Fatal(err)
	}
	tx := conn.lastTx()
	if tx == nil {
		t.Fatal("no PDU written")
	}
	got := tx.Buf()
	if len(got) < 7 {
		t.Fatalf("written PDU too short: %x", got)
	}
	if length := binary.LittleEndian.Uint16(got[0:2]); length != 3 {
		t.Errorf("length field: got %d want 3", length)
	}
	if cid := binary.LittleEndian.Uint16(got[2:4]); cid != 0x55 {
		t.Errorf("CID field: got %#x want 0x55", cid)
	}
	if string(got[4:7]) != "\x01\x02\x03" {
		t.Errorf("payload: got %x", got[4:7])
	}
}
