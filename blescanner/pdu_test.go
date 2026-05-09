package blescanner

import (
	"sync"
	"testing"
	"time"

	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	bleutil "github.com/BertoldVdb/go-ble/util"
)

// Verify that handleScanResult does NOT hold s.Lock() while invoking
// advertisingReportCallbacks. The callback re-enters the scanner via
// RegisterAdvertisingReportCallback (which itself takes s.Lock()) — if
// the dispatcher held the lock during the callback, this would deadlock.
func TestAdvertisingReportCallbackReentrant(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{})

	var reentered int
	var mu sync.Mutex

	s.RegisterAdvertisingReportCallback(func(pkt *BLEAdvertisingReport) bool {
		// Re-entering the scanner from within the callback used to
		// deadlock; with the snapshot-then-call fix it must succeed.
		s.RegisterAdvertisingReportCallback(func(pkt *BLEAdvertisingReport) bool {
			return false
		})
		mu.Lock()
		reentered++
		mu.Unlock()
		return false
	})

	ev := &hcievents.LEAdvertisingReportEvent{
		NumReports:  1,
		EventType:   []uint8{0x00},
		AddressType: []bleutil.MacAddrType{0},
		Address:     []bleutil.MacAddr{0x010203040506},
		Data:        [][]uint8{{0x02, 0x01, 0x06}},
		RSSI:        []uint8{0xCE},
	}

	done := make(chan struct{})
	go func() {
		s.handleScanResult(ev)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("deadlock: callback re-entry into scanner timed out")
	}

	mu.Lock()
	if reentered != 1 {
		t.Fatalf("callback ran %d times, want 1", reentered)
	}
	mu.Unlock()
}

// Test that callbacks registered from within a DeviceUpdatedCallback do
// not corrupt iteration. The snapshot-before-call fix replaced
// index-based iteration over a slice that could grow under it.
func TestDeviceUpdatedCallbackSnapshot(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{StoreGAPMap: true})

	var calls int
	var mu sync.Mutex
	s.RegisterDeviceUpdateCallback(func(dev *BLEDevice) {
		mu.Lock()
		calls++
		mu.Unlock()
		// Register a fresh callback during dispatch — it must NOT be
		// invoked during this same dispatch (snapshot semantics).
		s.RegisterDeviceUpdateCallback(func(dev *BLEDevice) {
			mu.Lock()
			calls++
			mu.Unlock()
		})
	})

	ev := &hcievents.LEAdvertisingReportEvent{
		NumReports:  1,
		EventType:   []uint8{0x00},
		AddressType: []bleutil.MacAddrType{0},
		Address:     []bleutil.MacAddr{0x111213141516},
		Data:        [][]uint8{{0x02, 0x01, 0x06}},
		RSSI:        []uint8{0xCE},
	}
	s.handleScanResult(ev)

	mu.Lock()
	defer mu.Unlock()
	if calls != 1 {
		t.Fatalf("got %d calls, want 1 (newly-added callback should not fire on the same dispatch)", calls)
	}
}

// Multiple callbacks must all fire, in registration order.
func TestAdvertisingReportCallbackMultipleFire(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{})

	var order []int
	var mu sync.Mutex
	for i := 0; i < 3; i++ {
		i := i
		s.RegisterAdvertisingReportCallback(func(pkt *BLEAdvertisingReport) bool {
			mu.Lock()
			order = append(order, i)
			mu.Unlock()
			return false
		})
	}

	ev := &hcievents.LEAdvertisingReportEvent{
		NumReports:  1,
		EventType:   []uint8{0x00},
		AddressType: []bleutil.MacAddrType{0},
		Address:     []bleutil.MacAddr{0x212223242526},
		Data:        [][]uint8{{0x02, 0x01, 0x06}},
		RSSI:        []uint8{0xCE},
	}
	s.handleScanResult(ev)

	mu.Lock()
	defer mu.Unlock()
	if len(order) != 3 || order[0] != 0 || order[1] != 1 || order[2] != 2 {
		t.Fatalf("call order: got %v want [0 1 2]", order)
	}
}
