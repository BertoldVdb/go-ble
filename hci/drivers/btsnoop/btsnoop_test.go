package btsnoop

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	hciconst "github.com/BertoldVdb/go-ble/hci/const"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
)

// fakeIface is a minimal HCIInterface used so we can wrap with btsnoop
// without spinning up a real driver.
type fakeIface struct {
	cb hciinterface.HCIRxHandler
}

func (f *fakeIface) Run() error                                    { return nil }
func (f *fakeIface) Close() error                                  { return nil }
func (f *fakeIface) SetRecvHandler(cb hciinterface.HCIRxHandler) error {
	f.cb = cb
	return nil
}
func (f *fakeIface) SendPacket(_ hciinterface.HCITxPacket) error { return nil }

type closableBuf struct {
	*bytes.Buffer
}

func (c *closableBuf) Close() error { return nil }

func TestWrapWritesHeader(t *testing.T) {
	buf := &closableBuf{Buffer: &bytes.Buffer{}}
	if _, err := Wrap(&fakeIface{}, buf); err != nil {
		t.Fatal(err)
	}
	got := buf.Bytes()
	if len(got) != 16 {
		t.Fatalf("header length: got %d want 16", len(got))
	}
	if !bytes.Equal(got[:8], []byte("btsnoop\x00")) {
		t.Errorf("magic: got %q want %q", got[:8], "btsnoop\x00")
	}
	if v := binary.BigEndian.Uint32(got[8:]); v != 1 {
		t.Errorf("version: got %d want 1", v)
	}
	if v := binary.BigEndian.Uint32(got[12:]); v != 1002 {
		t.Errorf("datalink: got %d want 1002", v)
	}
}

// shortWriter returns less than the full input (and no error) on the first
// call, then writes the rest. btsnoop must finish the record despite the
// short write — otherwise records become unparseable.
type shortWriter struct {
	buf  bytes.Buffer
	once sync.Once
}

func (s *shortWriter) Write(p []byte) (int, error) {
	short := false
	s.once.Do(func() { short = true })
	if short && len(p) > 1 {
		// accept only the first byte, force the caller to retry
		return s.buf.Write(p[:1])
	}
	return s.buf.Write(p)
}

func (s *shortWriter) Close() error { return nil }

func TestWrapRetriesShortWrite(t *testing.T) {
	w := &shortWriter{}
	if _, err := Wrap(&fakeIface{}, w); err != nil {
		t.Fatal(err)
	}
	got := w.buf.Bytes()
	if len(got) != 16 {
		t.Fatalf("header truncated: got %d bytes, want 16: %x", len(got), got)
	}
}

// erroringWriter fails after writing n bytes. Wrap must propagate the error
// rather than silently producing a half-written file.
type erroringWriter struct {
	buf      bytes.Buffer
	failAt   int
	written  int
	closed   bool
}

func (e *erroringWriter) Write(p []byte) (int, error) {
	if e.written >= e.failAt {
		return 0, errors.New("disk full")
	}
	remaining := e.failAt - e.written
	if remaining > len(p) {
		remaining = len(p)
	}
	n, _ := e.buf.Write(p[:remaining])
	e.written += n
	if remaining < len(p) {
		return n, errors.New("disk full")
	}
	return n, nil
}

func (e *erroringWriter) Close() error { e.closed = true; return nil }

func TestWrapPropagatesWriteError(t *testing.T) {
	w := &erroringWriter{failAt: 4}
	if _, err := Wrap(&fakeIface{}, w); err == nil {
		t.Fatal("expected write error")
	}
}

func TestWrapFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file mode is unix-specific")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "trace.btsnoop")

	_, err := WrapFile(&fakeIface{}, path)
	if err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	mode := st.Mode().Perm()
	// Must be 0600: snoop logs may contain LTKs, IRKs, OOB pairing data.
	if mode != 0o600 {
		t.Fatalf("file mode: got %#o want 0600", mode)
	}
}

// TestLogPacketRecordFormat checks that each packet produces a complete
// 24-byte record header followed by the data, and that flags encode
// CMD/EVT (data flag) and direction.
func TestLogPacketRecordFormat(t *testing.T) {
	buf := &closableBuf{Buffer: &bytes.Buffer{}}
	dev, err := Wrap(&fakeIface{}, buf)
	if err != nil {
		t.Fatal(err)
	}
	l := dev.(*logger)

	// HCI Command (data flag should be set, direction = TX => bit 0 clear)
	cmdData := []byte{hciconst.MsgTypeCommand, 0x03, 0x0C, 0x00}
	l.logPacket(true, time.Unix(0, 0), cmdData)

	// HCI ACL (data flag should NOT be set, direction = RX => bit 0 set)
	aclData := []byte{hciconst.MsgTypeACL, 0x40, 0x00, 0x00, 0x00}
	l.logPacket(false, time.Unix(0, 0), aclData)

	got := buf.Bytes()
	// header (16) + record1 (24+4) + record2 (24+5) = 73
	if got, want := len(got), 16+24+len(cmdData)+24+len(aclData); got != want {
		t.Fatalf("total bytes: got %d want %d", got, want)
	}

	// Inspect record 1
	rec1 := got[16:]
	origLen := binary.BigEndian.Uint32(rec1[:4])
	flags := binary.BigEndian.Uint32(rec1[8:12])
	if int(origLen) != len(cmdData) {
		t.Errorf("cmd record length: got %d want %d", origLen, len(cmdData))
	}
	// flags bit 0 = direction (1 = host→controller? no — code: !isTransmit ⇒ |=1)
	// We sent isTransmit=true, so bit 0 must be 0.
	if flags&1 != 0 {
		t.Errorf("cmd flags: direction bit set on TX")
	}
	// flags bit 1 = command/event
	if flags&2 == 0 {
		t.Errorf("cmd flags: data flag (bit 1) not set on Command")
	}

	// Inspect record 2
	rec2 := rec1[24+len(cmdData):]
	flags2 := binary.BigEndian.Uint32(rec2[8:12])
	if flags2&1 == 0 {
		t.Errorf("acl flags: direction bit clear on RX")
	}
	if flags2&2 != 0 {
		t.Errorf("acl flags: data flag set on ACL")
	}
}

func TestLogPacketStopsAfterError(t *testing.T) {
	w := &erroringWriter{failAt: 16} // header succeeds, first record fails
	dev, err := Wrap(&fakeIface{}, w)
	if err != nil {
		t.Fatal(err)
	}
	l := dev.(*logger)

	// First record triggers the failure path; logger marks itself failed.
	l.logPacket(true, time.Unix(0, 0), []byte{hciconst.MsgTypeCommand, 0, 0, 0})
	written1 := w.buf.Len()
	// Subsequent record must not write anything (failed flag set).
	l.logPacket(true, time.Unix(0, 0), []byte{hciconst.MsgTypeCommand, 0, 0, 0})
	if w.buf.Len() != written1 {
		t.Fatalf("logger continued writing after failure: %d -> %d", written1, w.buf.Len())
	}
}

// Compile-time interface check.
var _ io.WriteCloser = (*closableBuf)(nil)
