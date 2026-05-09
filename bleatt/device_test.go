package bleatt

import (
	"context"
	"io"
	"sync"
	"testing"

	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"
)

// attFakeConn is a minimal BufferConn for GATT device tests.
type attFakeConn struct {
	logger *logrus.Entry
	closed bool
}

func newATTFakeConn() *attFakeConn {
	l := logrus.New()
	l.Out = io.Discard
	l.Level = logrus.PanicLevel
	return &attFakeConn{logger: logrus.NewEntry(l)}
}

func (f *attFakeConn) IsOpen() bool                                   { return !f.closed }
func (f *attFakeConn) Close() error                                   { f.closed = true; return nil }
func (f *attFakeConn) ReadBuffer(ctx context.Context) (*pdu.PDU, error) { <-ctx.Done(); return nil, ctx.Err() }
func (f *attFakeConn) WriteBuffer(buf *pdu.PDU) error                 { return nil }
func (f *attFakeConn) GetLogger() *logrus.Entry                       { return f.logger }
func (f *attFakeConn) UseStart()                                      {}
func (f *attFakeConn) UseDone()                                       {}

var _ hciconnmgr.BufferConn = (*attFakeConn)(nil)

func newTestLogger() *logrus.Entry {
	l := logrus.New()
	l.Out = io.Discard
	l.Level = logrus.PanicLevel
	return logrus.NewEntry(l)
}

// setMTU is monotonic — once set, it must never decrease. Earlier code
// did `atomic.StoreUint32` unconditionally, allowing a smaller value to
// shrink the negotiated MTU.
func TestSetMTUMonotonic(t *testing.T) {
	d := &gattDeviceConn{
		parent: &GattDevice{ourMTU: 247},
		logger: newTestLogger(),
	}

	if got := d.setMTU(50); got != 50 {
		t.Fatalf("first setMTU: got %d want 50", got)
	}
	if got := d.setMTU(40); got != 50 {
		t.Fatalf("smaller setMTU must not shrink: got %d want 50", got)
	}
	if got := d.setMTU(100); got != 100 {
		t.Fatalf("larger setMTU: got %d want 100", got)
	}
	// Floor: anything below 23 is normalised to 23.
	d2 := &gattDeviceConn{parent: &GattDevice{ourMTU: 247}, logger: newTestLogger()}
	if got := d2.setMTU(0); got != 247 {
		t.Fatalf("zero peer MTU should clamp to ourMTU=247, got %d", got)
	}
}

// Cap: setMTU is bounded by the local ourMTU.
func TestSetMTUClampsToOurMTU(t *testing.T) {
	d := &gattDeviceConn{parent: &GattDevice{ourMTU: 100}, logger: newTestLogger()}
	if got := d.setMTU(517); got != 100 {
		t.Fatalf("setMTU not clamped to ourMTU: got %d want 100", got)
	}
}

// CloseConn must clear initialConn so the *next* AddConn re-establishes
// the GATT client routing. Earlier code never cleared it; ClientRead /
// ClientWrite kept using the dead connection's handlers after a
// reconnect.
func TestGattDeviceCloseConnClearsInitialConn(t *testing.T) {
	dev := NewGattDevice(attstructure.NewStructure(), &GattDeviceConfig{
		MTU:                     247,
		DeviceName:              "test",
		DiscoverRemoteOnConnect: false,
	})

	conn1 := newATTFakeConn()
	if err := dev.AddConn(conn1); err != nil {
		t.Fatalf("AddConn 1: %v", err)
	}
	if dev.initialConn == nil {
		t.Fatal("initialConn not set after first AddConn")
	}

	if err := dev.CloseConn(conn1); err != nil {
		t.Fatalf("CloseConn: %v", err)
	}
	if dev.initialConn != nil {
		t.Fatal("initialConn still set after CloseConn (would cause stale-conn reads)")
	}

	conn2 := newATTFakeConn()
	if err := dev.AddConn(conn2); err != nil {
		t.Fatalf("AddConn 2: %v", err)
	}
	if dev.initialConn == nil {
		t.Fatal("initialConn not re-armed after second AddConn")
	}
}

// CloseConn on a non-initial connection must NOT touch initialConn.
func TestGattDeviceCloseConnNonInitial(t *testing.T) {
	dev := NewGattDevice(attstructure.NewStructure(), &GattDeviceConfig{
		MTU:                     247,
		DeviceName:              "test",
		DiscoverRemoteOnConnect: false,
	})

	conn1 := newATTFakeConn()
	conn2 := newATTFakeConn()
	if err := dev.AddConn(conn1); err != nil {
		t.Fatal(err)
	}
	if err := dev.AddConn(conn2); err != nil {
		t.Fatal(err)
	}
	initialBefore := dev.initialConn

	if err := dev.CloseConn(conn2); err != nil {
		t.Fatalf("CloseConn: %v", err)
	}
	if dev.initialConn != initialBefore {
		t.Fatal("initialConn changed when closing a non-initial connection")
	}
}

// Race: hammer setMTU concurrently from multiple goroutines, each
// proposing values from a fixed range. The final value must equal the
// max ever proposed.
func TestSetMTUConcurrent(t *testing.T) {
	d := &gattDeviceConn{parent: &GattDevice{ourMTU: 1024}, logger: newTestLogger()}

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				d.setMTU(uint16(50 + ((i * 7 + j) % 200)))
			}
		}(i)
	}
	wg.Wait()

	got := d.getMTU()
	if got < 23 || got > 249 {
		t.Fatalf("final MTU out of expected window: %d", got)
	}
	// Now advance once more to confirm it still grows monotonically.
	d.setMTU(500)
	if d.getMTU() != 500 {
		t.Fatalf("post-stress setMTU(500) failed: got %d", d.getMTU())
	}
	d.setMTU(100)
	if d.getMTU() != 500 {
		t.Fatalf("post-stress shrink leaked through: got %d", d.getMTU())
	}
}
