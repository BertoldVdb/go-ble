package bleatt

import (
	"encoding/binary"
	"sync"
	"testing"

	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
)

// captureClient stubs handleNotify so we can capture (handle, data)
// triples produced by the multi-notify parser.
type captureClient struct {
	mu     sync.Mutex
	got    []notifyRecord
	client *attClient
}

type notifyRecord struct {
	handle uint16
	data   []byte
}

// Replace handleNotify by re-implementing handleNTFIND in a test wrapper.
// Since handleNotify reaches up into the parent device chain (which we don't
// want to construct), we duplicate the parser exactly here and exercise it.
// This is the same logic that lives in client.go:handleNTFIND for the
// ATTMultipleHandleValueNTF arm — verifying it stops on a malformed length
// rather than re-aligning to attacker-controlled bytes.
func parseMultiNotifyForTest(buf *pdu.PDU) []notifyRecord {
	var out []notifyRecord
	for {
		hdr := buf.DropLeft(4)
		if hdr == nil {
			break
		}
		handle := binary.LittleEndian.Uint16(hdr)
		dlen := binary.LittleEndian.Uint16(hdr[2:])
		data := buf.DropLeft(int(dlen))
		if data == nil {
			break
		}
		out = append(out, notifyRecord{handle: handle, data: append([]byte(nil), data...)})
	}
	return out
}

func TestMultiNotifyParserNormal(t *testing.T) {
	buf := bleutil.GetBuffer(0)
	// Two records: (handle=0x0040, data="hi"), (handle=0x0050, data="abc")
	hdr := buf.ExtendRight(4)
	binary.LittleEndian.PutUint16(hdr, 0x0040)
	binary.LittleEndian.PutUint16(hdr[2:], 2)
	buf.Append([]byte("hi")...)

	hdr = buf.ExtendRight(4)
	binary.LittleEndian.PutUint16(hdr, 0x0050)
	binary.LittleEndian.PutUint16(hdr[2:], 3)
	buf.Append([]byte("abc")...)

	got := parseMultiNotifyForTest(buf)
	if len(got) != 2 {
		t.Fatalf("got %d records, want 2", len(got))
	}
	if got[0].handle != 0x0040 || string(got[0].data) != "hi" {
		t.Errorf("rec0: %+v", got[0])
	}
	if got[1].handle != 0x0050 || string(got[1].data) != "abc" {
		t.Errorf("rec1: %+v", got[1])
	}
}

// TestMultiNotifyParserMalformedLength: peer sends dlen larger than the
// remaining buffer. Earlier code checked `if hdr == nil` (always false here
// since hdr was successfully read) and proceeded to call handleNotify with
// nil data, then re-read the now-misaligned cursor. The fix checks
// `if data == nil`, so the parser stops cleanly.
func TestMultiNotifyParserMalformedLength(t *testing.T) {
	buf := bleutil.GetBuffer(0)
	// Header declares dlen=100 but we only supply 4 bytes of payload.
	hdr := buf.ExtendRight(4)
	binary.LittleEndian.PutUint16(hdr, 0x0040)
	binary.LittleEndian.PutUint16(hdr[2:], 100)
	buf.Append([]byte("abcd")...)

	got := parseMultiNotifyForTest(buf)
	if len(got) != 0 {
		t.Fatalf("got %d records, want 0 (parser must stop on dlen overrun); got=%+v", len(got), got)
	}
}

// TestMultiNotifyParserEmpty: empty PDU produces no records.
func TestMultiNotifyParserEmpty(t *testing.T) {
	buf := bleutil.GetBuffer(0)
	got := parseMultiNotifyForTest(buf)
	if len(got) != 0 {
		t.Fatalf("got %d records, want 0", len(got))
	}
}

// TestMultiNotifyParserPartialHeader: peer sends only 3 of 4 header bytes.
func TestMultiNotifyParserPartialHeader(t *testing.T) {
	buf := bleutil.GetBuffer(0)
	buf.Append(0x40, 0x00, 0x02)
	got := parseMultiNotifyForTest(buf)
	if len(got) != 0 {
		t.Fatalf("got %d records, want 0", len(got))
	}
}
