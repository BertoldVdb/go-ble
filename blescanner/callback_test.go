package blescanner

import (
	"sync"
	"testing"

	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	bleutil "github.com/BertoldVdb/go-ble/util"
)

// TestUnregisterDeviceUpdateCallback verifies that an unregistered
// callback no longer fires.
func TestUnregisterDeviceUpdateCallback(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{StoreGAPMap: true})

	var calls int
	var mu sync.Mutex
	h := s.RegisterDeviceUpdateCallback(func(dev *BLEDevice) {
		mu.Lock()
		calls++
		mu.Unlock()
	})
	if h == 0 {
		t.Fatal("zero handle returned")
	}

	ev := func(addr bleutil.MacAddr) *hcievents.LEAdvertisingReportEvent {
		return &hcievents.LEAdvertisingReportEvent{
			NumReports:  1,
			EventType:   []uint8{0x00},
			AddressType: []bleutil.MacAddrType{0},
			Address:     []bleutil.MacAddr{addr},
			Data:        [][]uint8{{0x02, 0x01, 0x06}},
			RSSI:        []uint8{0xCE},
		}
	}

	s.handleScanResult(ev(0x010203040506))
	s.UnregisterDeviceUpdateCallback(h)
	s.handleScanResult(ev(0x111213141516))

	mu.Lock()
	defer mu.Unlock()
	if calls != 1 {
		t.Fatalf("got %d calls, want 1 (callback should not fire after Unregister)", calls)
	}
}

func TestUnregisterAdvertisingReportCallback(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{})

	var fired int
	var mu sync.Mutex
	h := s.RegisterAdvertisingReportCallback(func(pkt *BLEAdvertisingReport) bool {
		mu.Lock()
		fired++
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
	s.handleScanResult(ev)
	s.UnregisterAdvertisingReportCallback(h)
	s.handleScanResult(ev)

	mu.Lock()
	defer mu.Unlock()
	if fired != 1 {
		t.Fatalf("fired %d times, want 1", fired)
	}
}
