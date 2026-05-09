package bleadvertiser

import (
	"testing"
)

// TestSlotPointerStableAcrossAppend verifies that slot pointers handed out
// by LegacyAdvertisingGetSlot remain valid after subsequent allocations.
// Earlier code stored slots by value in a slice and returned `&slice[i]`;
// once `append` reallocated, every previously-returned pointer became
// stale.
func TestSlotPointerStableAcrossAppend(t *testing.T) {
	a := &BLEAdvertiser{}

	first := a.LegacyAdvertisingGetSlot()
	first.data.IntervalMin = 0x1234

	// Allocate enough additional slots to force the underlying slice to
	// grow several times.
	for i := 0; i < 32; i++ {
		_ = a.LegacyAdvertisingGetSlot()
	}

	if first.data.IntervalMin != 0x1234 {
		t.Fatalf("slot data corrupted across slice growth: got %#x want 0x1234", first.data.IntervalMin)
	}
	// Confirm the pointer still refers to the live slot (lookup by index).
	if a.legacyAdvertisingSlots[0] != first {
		t.Fatalf("slot 0 pointer mismatch: live=%p returned=%p", a.legacyAdvertisingSlots[0], first)
	}
}

// TestSlotIndexAssignment confirms that newly-allocated slots are stored at
// their declared index and that closing one frees it for reuse.
func TestSlotIndexAssignment(t *testing.T) {
	a := &BLEAdvertiser{}

	s0 := a.LegacyAdvertisingGetSlot()
	s1 := a.LegacyAdvertisingGetSlot()
	if s0.index != 0 || s1.index != 1 {
		t.Fatalf("indices: got %d, %d want 0, 1", s0.index, s1.index)
	}

	// Initialize update channel so Close() can do its non-blocking send.
	a.legacyAdvertisingUpdateChan = make(chan int, 1)
	s0.Close()
	if s0.valid {
		t.Fatal("slot remained valid after Close")
	}

	s2 := a.LegacyAdvertisingGetSlot()
	// The freed slot 0 should be reused.
	if s2.index != 0 {
		t.Fatalf("expected reused index 0, got %d", s2.index)
	}
	if !s2.valid {
		t.Fatal("reused slot not marked valid")
	}
	if a.legacyAdvertisingSlots[0] != s2 {
		t.Fatal("reused slot pointer mismatch")
	}
	_ = s1 // silence unused
}

// LegacyAdvertisingSlot.GetData returns the value stored last.
func TestSlotGetData(t *testing.T) {
	a := &BLEAdvertiser{
		legacyAdvertisingUpdateChan: make(chan int, 1),
	}
	slot := a.LegacyAdvertisingGetSlot()

	// Default value
	d, err := slot.GetData()
	if err != nil {
		t.Fatal(err)
	}
	if d.Active != false {
		t.Error("default slot should be inactive")
	}

	// After replace, GetData reflects the new state.
	_, err = slot.ReplaceData(true, LegacyAdvertisingData{
		Active:      true,
		IntervalMin: 0x100,
	})
	if err != nil {
		t.Fatal(err)
	}
	d, err = slot.GetData()
	if err != nil {
		t.Fatal(err)
	}
	if !d.Active || d.IntervalMin != 0x100 {
		t.Errorf("GetData mismatch: %+v", d)
	}
}

// SlotClose marks the slot invalid and is reused by the next GetSlot.
func TestSlotCloseAndReuse(t *testing.T) {
	a := &BLEAdvertiser{
		legacyAdvertisingUpdateChan: make(chan int, 1),
	}

	s1 := a.LegacyAdvertisingGetSlot()
	idx1 := s1.index
	s1.Close()
	if s1.valid {
		t.Error("slot still valid after Close")
	}

	s2 := a.LegacyAdvertisingGetSlot()
	if s2.index != idx1 {
		t.Errorf("expected reused index %d, got %d", idx1, s2.index)
	}
}

// TestReplaceDataVersion checks the version-monotonicity guard used by
// LegacyAdvertisingSetConnection's cancel function: replacing with a
// stale version (lower than current) without `force` must fail.
func TestReplaceDataVersion(t *testing.T) {
	a := &BLEAdvertiser{
		legacyAdvertisingUpdateChan: make(chan int, 1),
	}

	slot := a.LegacyAdvertisingGetSlot()
	d1, err := slot.ReplaceData(true, LegacyAdvertisingData{Active: true})
	if err != nil {
		t.Fatal(err)
	}
	if d1.version != 1 {
		t.Fatalf("first replace version: got %d want 1", d1.version)
	}

	// Force-replace with a stale version still works.
	d2, err := slot.ReplaceData(true, LegacyAdvertisingData{version: 0})
	if err != nil {
		t.Fatal(err)
	}
	if d2.version != 2 {
		t.Fatalf("second replace version: got %d want 2", d2.version)
	}

	// Non-force with stale version is rejected.
	_, err = slot.ReplaceData(false, LegacyAdvertisingData{version: 1})
	if err != ErrorExpired {
		t.Fatalf("stale replace: got %v want %v", err, ErrorExpired)
	}
}
