package hcidriverserial

import (
	"bytes"
	"errors"
	"io"
	"sync"
	"testing"
	"time"

	hciconst "github.com/BertoldVdb/go-ble/hci/const"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
)

// Historical bug: `5 + lo | (hi << 8)` parsed as `(5+lo) | (hi<<8)` due to
// Go operator precedence — wrong whenever `lo` and `hi` shared bits. The
// current code parenthesizes the OR expression.
func TestGetHCIPacketLengthACL(t *testing.T) {
	cases := []struct {
		name string
		pkt  []byte
		want int
	}{
		{"len=0", []byte{hciconst.MsgTypeACL, 0x40, 0x00, 0x00, 0x00}, 5},
		{"len=1", []byte{hciconst.MsgTypeACL, 0x40, 0x00, 0x01, 0x00}, 6},
		{"len=255", []byte{hciconst.MsgTypeACL, 0x40, 0x00, 0xFF, 0x00}, 260},
		// 0xFF + (0x01<<8) = 511. Buggy expression yields 261; correct yields 516.
		{"len=511 (precedence trap)", []byte{hciconst.MsgTypeACL, 0x40, 0x00, 0xFF, 0x01}, 516},
		{"len=0x1234", []byte{hciconst.MsgTypeACL, 0, 0, 0x34, 0x12}, 5 + 0x1234},
		{"len=0xFFFF", []byte{hciconst.MsgTypeACL, 0, 0, 0xFF, 0xFF}, 5 + 0xFFFF},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := getHCIPacketLength(c.pkt); got != c.want {
				t.Fatalf("got %d, want %d", got, c.want)
			}
		})
	}
}

func TestGetHCIPacketLengthEvent(t *testing.T) {
	if got := getHCIPacketLength([]byte{hciconst.MsgTypeEvent, 0x0E, 0x04}); got != 7 {
		t.Fatalf("event len: got %d, want 7", got)
	}
}

func TestGetHCIPacketLengthCommand(t *testing.T) {
	if got := getHCIPacketLength([]byte{hciconst.MsgTypeCommand, 0x03, 0x0C, 0x00}); got != 4 {
		t.Fatalf("cmd len: got %d, want 4", got)
	}
}

func TestGetHCIPacketLengthSCO(t *testing.T) {
	if got := getHCIPacketLength([]byte{hciconst.MsgTypeSCO, 0x40, 0x00, 0x10}); got != 20 {
		t.Fatalf("sco len: got %d, want 20", got)
	}
}

// ISO uses 14-bit length: bits 0..5 of byte[4] are the high 6 bits;
// bits 6..7 are reserved and must not be counted as length.
func TestGetHCIPacketLengthISO(t *testing.T) {
	cases := []struct {
		name string
		pkt  []byte
		want int
	}{
		{"low byte only", []byte{hciconst.MsgTypeISO, 0x40, 0x00, 0x40, 0x00}, 5 + 0x40},
		{"14-bit max", []byte{hciconst.MsgTypeISO, 0x40, 0x00, 0xFF, 0x3F}, 5 + 0x3FFF},
		{"reserved bits ignored", []byte{hciconst.MsgTypeISO, 0x40, 0x00, 0x00, 0xC0}, 5},
		{"wide len with status bits set", []byte{hciconst.MsgTypeISO, 0, 0, 0x34, 0xD2}, 5 + 0x1234},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := getHCIPacketLength(c.pkt); got != c.want {
				t.Fatalf("got %d, want %d", got, c.want)
			}
		})
	}
}

func TestGetHCIPacketLengthInvalid(t *testing.T) {
	if got := getHCIPacketLength([]byte{0xAA, 0, 0, 0, 0}); got != -1 {
		t.Fatalf("unknown type: got %d, want -1", got)
	}
	if got := getHCIPacketLength([]byte{hciconst.MsgTypeACL, 0x40}); got != 0 {
		t.Fatalf("short buf: got %d, want 0", got)
	}
	if got := getHCIPacketLength(nil); got != 0 {
		t.Fatalf("empty buf: got %d, want 0", got)
	}
}

// fakePort feeds a static byte stream to Run() and captures Writes.
type fakePort struct {
	mu       sync.Mutex
	rxBuf    *bytes.Buffer
	txBuf    *bytes.Buffer
	closed   bool
	rxNotify chan struct{}
}

func newFakePort(input []byte) *fakePort {
	return &fakePort{
		rxBuf:    bytes.NewBuffer(input),
		txBuf:    &bytes.Buffer{},
		rxNotify: make(chan struct{}, 1),
	}
}

func (f *fakePort) Read(p []byte) (int, error) {
	for {
		f.mu.Lock()
		if f.closed {
			f.mu.Unlock()
			return 0, io.EOF
		}
		if f.rxBuf.Len() > 0 {
			n, err := f.rxBuf.Read(p)
			f.mu.Unlock()
			return n, err
		}
		f.mu.Unlock()
		select {
		case <-f.rxNotify:
		case <-time.After(2 * time.Second):
			return 0, io.EOF
		}
	}
}

func (f *fakePort) Write(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return 0, errors.New("closed")
	}
	return f.txBuf.Write(p)
}

func (f *fakePort) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.closed = true
	select {
	case f.rxNotify <- struct{}{}:
	default:
	}
	return nil
}

// Close is idempotent — multiple Close() calls (including the deferred
// one inside Run) must invoke the underlying port.Close exactly once.
func TestSerialCloseIdempotent(t *testing.T) {
	port := newFakePort(nil)
	dev, _ := OpenPort(port)

	var calls int
	type wrapped struct{ *fakePort }
	// Wrap the port so we can count Close calls.
	hci := dev.(*HCISerial)
	hci.port = &countingPort{fakePort: port, calls: &calls}

	if err := hci.Close(); err != nil {
		t.Fatal(err)
	}
	if err := hci.Close(); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("port.Close called %d times, want 1", calls)
	}
}

func TestSerialSendPacketAfterClose(t *testing.T) {
	port := newFakePort(nil)
	dev, _ := OpenPort(port)

	dev.Close()
	if err := dev.SendPacket(hciinterface.HCITxPacket{Data: []byte{1, 2, 3}}); err == nil {
		t.Fatal("expected error from SendPacket on closed driver")
	}
}

type countingPort struct {
	*fakePort
	calls *int
}

func (c *countingPort) Close() error {
	*c.calls++
	return c.fakePort.Close()
}

// Integration: feed an event and an ACL packet through Run() and verify
// the H4 framer reassembles each correctly.
func TestSerialFraming(t *testing.T) {
	evt := []byte{hciconst.MsgTypeEvent, 0x0E, 0x04, 0x01, 0x02, 0x03, 0x04}
	acl := []byte{hciconst.MsgTypeACL, 0x40, 0x00, 0x03, 0x00, 0xAA, 0xBB, 0xCC}

	in := append(append([]byte{}, evt...), acl...)
	port := newFakePort(in)

	dev, err := OpenPort(port)
	if err != nil {
		t.Fatal(err)
	}

	var got [][]byte
	var gotMu sync.Mutex
	done := make(chan struct{})
	dev.SetRecvHandler(func(pkt hciinterface.HCIRxPacket) error {
		gotMu.Lock()
		got = append(got, append([]byte(nil), pkt.Data...))
		if len(got) == 2 {
			close(done)
		}
		gotMu.Unlock()
		return nil
	})

	runDone := make(chan error, 1)
	go func() { runDone <- dev.Run() }()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for packets")
	}

	dev.Close()
	<-runDone

	gotMu.Lock()
	defer gotMu.Unlock()
	if len(got) != 2 {
		t.Fatalf("got %d packets, want 2", len(got))
	}
	if !bytes.Equal(got[0], evt) {
		t.Errorf("event mismatch: got %x want %x", got[0], evt)
	}
	if !bytes.Equal(got[1], acl) {
		t.Errorf("acl mismatch: got %x want %x", got[1], acl)
	}
}
