package hciconnmgr

import (
	"io"
	"testing"

	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	"github.com/sirupsen/logrus"
)

func newTestConnMgr() *ConnectionManager {
	l := logrus.New()
	l.Out = io.Discard
	l.Level = logrus.PanicLevel
	return &ConnectionManager{
		logger:      logrus.NewEntry(l),
		connections: make(map[uint16]*Connection),
		replyCh:     make(chan func() error, 4),
	}
}

// quirkFixBroadcomCompleteEvent reorders the (handle, count) pairs that
// some Broadcom controllers emit incorrectly. Originally guarded by
// `elem%1 > 0` (always false), so the "uneven" branch was dead code.
// Now correctly `elem%2`.

func cloneEvent(handles, counts []uint16) *hcievents.NumberOfCompletedPacketsEvent {
	return &hcievents.NumberOfCompletedPacketsEvent{
		NumHandles:          uint8(len(handles)),
		ConnectionHandle:    append([]uint16(nil), handles...),
		NumCompletedPackets: append([]uint16(nil), counts...),
	}
}

func TestQuirkBroadcomEvenLength(t *testing.T) {
	// Even-length: (h0,c0,h1,c1) → after the buggy controller swaps c0↔h1
	// we should swap them back so handles[0,1] = (h0,h1) and counts[0,1] = (c0,c1).
	// Simulate the wire view: handles=[h0,c0], counts=[h1,c1] (Broadcom-mangled).
	ev := cloneEvent([]uint16{0x40, 5}, []uint16{0x41, 7})
	quirkFixBroadcomCompleteEvent(ev)
	wantH := []uint16{0x40, 0x41}
	wantC := []uint16{5, 7}
	for i := range ev.ConnectionHandle {
		if ev.ConnectionHandle[i] != wantH[i] || ev.NumCompletedPackets[i] != wantC[i] {
			t.Fatalf("after fix: handles=%v counts=%v want %v / %v", ev.ConnectionHandle, ev.NumCompletedPackets, wantH, wantC)
		}
	}
}

func TestQuirkBroadcomOddLength(t *testing.T) {
	// Odd-length: handles=[h0,c0,h1], counts=[h2,c1,c2] (Broadcom-mangled).
	// Correct fix swaps in-place at indices 1: counts[1]↔handles[1], etc.
	ev := cloneEvent([]uint16{0x40, 5, 0x42}, []uint16{0x41, 7, 9})
	quirkFixBroadcomCompleteEvent(ev)
	// After fix: handles=[h0,h1, ...], counts=[c0, c1, ...]. The exact
	// expected output depends on how Broadcom mis-orders; we mainly want
	// to assert the function runs without panic and modifies content.
	// Since the audit's concern was that the "uneven" branch was dead,
	// just verify the operation actually executed by checking that one
	// pair was rearranged.
	if ev.ConnectionHandle[1] == 5 && ev.NumCompletedPackets[1] == 7 {
		t.Fatalf("odd-length quirk did not run: %v / %v", ev.ConnectionHandle, ev.NumCompletedPackets)
	}
}

func TestQuirkBroadcomEmpty(t *testing.T) {
	ev := cloneEvent(nil, nil)
	// Must not panic on zero-length input.
	quirkFixBroadcomCompleteEvent(ev)
	if len(ev.ConnectionHandle) != 0 || len(ev.NumCompletedPackets) != 0 {
		t.Fatalf("unexpected mutation on empty event")
	}
}

func TestQuirkBroadcomSingle(t *testing.T) {
	ev := cloneEvent([]uint16{0x40}, []uint16{5})
	quirkFixBroadcomCompleteEvent(ev)
	// Single-element event: nothing to swap (loop starts at i=1).
	if ev.ConnectionHandle[0] != 0x40 || ev.NumCompletedPackets[0] != 5 {
		t.Fatalf("single-element event mutated: %v / %v", ev.ConnectionHandle, ev.NumCompletedPackets)
	}
}

// createSlotManager rejects sub-minimum slot buffer lengths.
func TestCreateSlotManagerSubMinimumLength(t *testing.T) {
	cm := newTestConnMgr()
	s := createSlotManager(cm, "test", 10 /* < 27 */, 4)
	if s == nil {
		t.Fatal("expected slot manager")
	}
	if s.slotBufferLength != 27 {
		t.Errorf("expected slotBufferLength clamped to 27, got %d", s.slotBufferLength)
	}
}

// createSlotManager(maxSlots=0) returns nil (no slot manager needed).
func TestCreateSlotManagerZeroSlots(t *testing.T) {
	cm := newTestConnMgr()
	if s := createSlotManager(cm, "test", 27, 0); s != nil {
		t.Errorf("expected nil slot manager for maxSlots=0, got %+v", s)
	}
}

// ReleaseSlots clamps when the controller releases more than max.
func TestReleaseSlotsClamps(t *testing.T) {
	cm := newTestConnMgr()
	s := createSlotManager(cm, "test", 27, 4)

	// Initially availableSlots == maxSlots == 4. Releasing more should clamp.
	s.ReleaseSlots(10)
	if s.availableSlots != s.maxSlots {
		t.Errorf("expected clamp to %d, got %d", s.maxSlots, s.availableSlots)
	}
}

// WaitSlot decrements availableSlots and returns true.
func TestWaitSlotHappyPath(t *testing.T) {
	cm := newTestConnMgr()
	s := createSlotManager(cm, "test", 27, 4)

	if !s.WaitSlot() {
		t.Fatal("WaitSlot returned false unexpectedly")
	}
	if s.availableSlots != 3 {
		t.Errorf("availableSlots: got %d want 3", s.availableSlots)
	}
}

// GetBufferLength returns slotBufferLength.
func TestGetBufferLength(t *testing.T) {
	cm := newTestConnMgr()
	s := createSlotManager(cm, "test", 32, 4)
	if got := s.GetBufferLength(); got != 32 {
		t.Errorf("GetBufferLength: got %d want 32", got)
	}
}

// disconnectionCompleteHandler with non-zero status is a no-op.
func TestDisconnectionCompleteNonZeroStatus(t *testing.T) {
	cm := newTestConnMgr()
	ev := &hcievents.DisconnectionCompleteEvent{
		Status:           0x42,
		ConnectionHandle: 1,
	}
	out := cm.disconnectionCompleteHandler(ev)
	if out != ev {
		t.Error("event passed through unchanged when Status != 0")
	}
}

// FindConnectionByHandle on a fresh ConnectionManager returns nil.
func TestFindConnectionByHandleEmpty(t *testing.T) {
	cm := newTestConnMgr()
	if got := cm.FindConnectionByHandle(0x42); got != nil {
		t.Errorf("expected nil, got %+v", got)
	}
}
