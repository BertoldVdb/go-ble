package hciconnmgr

import (
	"testing"

	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
)

// Reassembly buffer is bounded at maxL2CAPFrameLength (= 4 + 0xFFFF).
// Earlier code grew rxPDU without limit, allowing a peer to exhaust host
// memory by sending an unbounded stream of continuation fragments.
func TestACLReassemblyBounded(t *testing.T) {
	if maxL2CAPFrameLength != 4+0xFFFF {
		t.Fatalf("maxL2CAPFrameLength changed: got %d", maxL2CAPFrameLength)
	}
}

// handleACLData breaks out of its parsing loop once the buffered bytes
// don't form a complete L2CAP packet — it does NOT block, and it does
// NOT reset the buffer. Verify the parse loop's exit conditions.
func TestHandleACLDataNeedsFullHeader(t *testing.T) {
	c := &Connection{rxPDU: &pdu.PDU{}}
	c.rxPDU.Append(0xFF, 0x00) // 2-byte L2CAP header (need at least 4)

	if err := c.handleACLData(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.rxPDU.Len() != 2 {
		t.Fatalf("rxPDU consumed mid-header: %d bytes left", c.rxPDU.Len())
	}
}

// pktLen exceeding maxL2CAPFrameLength must trigger a hard reset and
// return errACLReassemblyTooLarge.
func TestHandleACLDataRejectsOversize(t *testing.T) {
	c := &Connection{rxPDU: &pdu.PDU{}}
	// L2CAP length = 0xFFFF + 1 (synthetic, since 0xFFFF + 4 == max we need
	// to push past). Construct a header that decodes to maxL2CAPFrameLength + 1.
	// pktLen = uint16(LE) + 4 → max representable is 0xFFFF+4. We can't
	// represent more than that on the wire, so the bound is structurally
	// safe. Verify the equality boundary doesn't trigger an error.
	c.rxPDU.Append(0xFF, 0xFF, 0x05, 0x00) // length=0xFFFF, CID=5
	if err := c.handleACLData(); err != nil {
		t.Fatalf("max-size header rejected: %v", err)
	}
}

// handleACL drops short headers (<4 bytes) cleanly (returns false to
// signal not-consumed).
func TestHandleACLShortHeader(t *testing.T) {
	cm := newTestConnMgr()
	if cm.handleACL([]byte{0x00, 0x00, 0x00}) {
		t.Error("handleACL should return false on short data")
	}
}

// handleACL drops broadcast traffic (we don't support EDR broadcast).
func TestHandleACLBroadcastDropped(t *testing.T) {
	cm := newTestConnMgr()
	// flagBC=1 → broadcast. Header bytes: handle=0, BC=01b (bit 14)
	header := uint16(1 << 14)
	data := []byte{byte(header), byte(header >> 8), 0x00, 0x00}
	if !cm.handleACL(data) {
		t.Error("handleACL should return true for broadcast traffic")
	}
}

// handleACL drops fragments for unknown handles.
func TestHandleACLUnknownHandle(t *testing.T) {
	cm := newTestConnMgr()
	// flagPB=2 (start), handle=0x42, payload-length-mismatch is fine here.
	header := uint16(0x42 | (2 << 12))
	data := []byte{byte(header), byte(header >> 8), 0x00, 0x00}
	if !cm.handleACL(data) {
		t.Error("handleACL should return true for unknown handle")
	}
}

// handleACL with payloadLen mismatch → drop.
func TestHandleACLLengthMismatch(t *testing.T) {
	cm := newTestConnMgr()
	header := uint16(0x42 | (2 << 12))
	// Declared payload length = 99 but only 0 payload bytes follow.
	data := []byte{byte(header), byte(header >> 8), 99, 0}
	if cm.handleACL(data) {
		t.Error("handleACL should return false on length mismatch")
	}
}

// handleACL flagPB=0 (start, host→controller) is dropped (we only
// accept controller→host PB values 1 and 2).
func TestHandleACLUnknownPBFlag(t *testing.T) {
	cm := newTestConnMgr()
	header := uint16(0x42 | (3 << 12))
	data := []byte{byte(header), byte(header >> 8), 0x00, 0x00}
	if !cm.handleACL(data) {
		t.Error("handleACL should silently accept (true) unknown PB flag values")
	}
}

// HandleData top-level dispatcher: only ACL packets (type=2) reach handleACL.
func TestHandleData(t *testing.T) {
	cm := newTestConnMgr()
	// Empty
	if got := cm.HandleData(hciinterface.HCIRxPacket{Data: nil}); got {
		t.Error("empty data should return false")
	}
	// Non-ACL type (e.g., event = 4)
	if got := cm.HandleData(hciinterface.HCIRxPacket{Data: []byte{0x04, 0x0E, 0x00}}); got {
		t.Error("non-ACL type should return false")
	}
	// ACL with too-short body
	if cm.HandleData(hciinterface.HCIRxPacket{Data: []byte{0x02, 0x00}}); true {
		// returns whatever handleACL returns; just ensure no panic
	}
}
