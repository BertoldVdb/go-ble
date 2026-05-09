package blescanner

import (
	"testing"

	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	bleutil "github.com/BertoldVdb/go-ble/util"
)

// helper: drive one advertisement through handleScanResult and return
// the resulting BLEDevice.
func feedAdv(t *testing.T, s *BLEScanner, addr bleutil.MacAddr, data []byte) *BLEDevice {
	t.Helper()
	ev := &hcievents.LEAdvertisingReportEvent{
		NumReports:  1,
		EventType:   []uint8{0x00},
		AddressType: []bleutil.MacAddrType{0},
		Address:     []bleutil.MacAddr{addr},
		Data:        [][]uint8{data},
		RSSI:        []uint8{0xCE},
	}
	s.handleScanResult(ev)
	dev := s.GetDevice(bleutil.BLEAddr{MacAddr: addr})
	if dev == nil {
		t.Fatal("device not tracked")
	}
	return dev
}

func TestGAPParseFlagsAndName(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{StoreGAPMap: true})
	// Flags (0x01) = 0x06; Complete Local Name (0x09) = "Hi"
	data := []byte{
		0x02, 0x01, 0x06,
		0x03, 0x09, 'H', 'i',
	}
	dev := feedAdv(t, s, 0x010203040506, data)

	if got := dev.GetFlags(); got != 0x06 {
		t.Errorf("flags: got %#x want 0x06", got)
	}
	if got := dev.GetName(); got != "Hi" {
		t.Errorf("name: got %q want %q", got, "Hi")
	}
	name, state := dev.GetNameType()
	if state != 2 || name != "Hi" {
		t.Errorf("name+state: %q/%d", name, state)
	}
}

func TestGAPShortNameDoesNotOverwriteFull(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{StoreGAPMap: true})
	full := []byte{0x06, 0x09, 'A', 'B', 'C', 'D', 'E'} // Complete name
	short := []byte{0x03, 0x08, 'X', 'Y'}               // Shortened name
	dev := feedAdv(t, s, 0x010203040507, full)
	if dev.GetName() != "ABCDE" {
		t.Fatalf("expected 'ABCDE', got %q", dev.GetName())
	}
	feedAdv(t, s, 0x010203040507, short)
	if dev.GetName() != "ABCDE" {
		t.Fatalf("short name overwrote full: %q", dev.GetName())
	}
}

func TestGAPParseTXPower(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{StoreGAPMap: true})
	data := []byte{0x02, 0x0a, 0xfb} // TXPower = -5
	dev := feedAdv(t, s, 0x010203040508, data)
	if dev.GetTXPower() != -5 {
		t.Errorf("TXPower: got %d want -5", dev.GetTXPower())
	}
}

func TestGAPParseUUIDList(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{StoreGAPMap: true})
	// Complete list of 16-bit Service UUIDs (type 0x03): 0x180a, 0x180f
	data := []byte{0x05, 0x03, 0x0a, 0x18, 0x0f, 0x18}
	dev := feedAdv(t, s, 0x010203040509, data)

	uuids := dev.GetServices(0, nil)
	if len(uuids) != 2 {
		t.Fatalf("got %d services, want 2", len(uuids))
	}
}

func TestGAPManufacturerSpecific(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{StoreGAPMap: true})

	got := make(chan *GAPRecord, 1)
	s.SetManufacturerSpecificCallback(0x004C, func(d *BLEDevice, r *GAPRecord) {
		select {
		case got <- r:
		default:
		}
	})
	// Manufacturer-specific (0xff): manufacturer ID 0x004C (Apple), payload "Hi"
	data := []byte{0x05, 0xff, 0x4c, 0x00, 'H', 'i'}
	feedAdv(t, s, 0x01020304050a, data)

	select {
	case r := <-got:
		if r.Type != 0xff {
			t.Errorf("manufacturer-specific type: got %#x", r.Type)
		}
	default:
		t.Fatal("manufacturer-specific callback did not fire")
	}

	// Unregister
	s.SetManufacturerSpecificCallback(0x004C, nil)
}

func TestGAPRecordCopyTo(t *testing.T) {
	src := &GAPRecord{Type: 0x09, Data: []byte("hello")}
	var buf GAPRecord
	out := src.copyTo(&buf)
	if out.Type != 0x09 || string(out.Data) != "hello" {
		t.Errorf("copyTo: %+v", *out)
	}
	// Append to source must not affect the copy.
	src.Data = append(src.Data, '!')
	if string(out.Data) != "hello" {
		t.Errorf("copyTo did not isolate: %q", out.Data)
	}
}

func TestGAPRecordCopyToNilDest(t *testing.T) {
	src := &GAPRecord{Type: 0x09, Data: []byte("x")}
	out := src.copyTo(nil)
	if out == nil {
		t.Fatal("copyTo(nil) returned nil")
	}
	if string(out.Data) != "x" {
		t.Errorf("copyTo(nil) data: got %q", out.Data)
	}
}

func TestGAPRecordString(t *testing.T) {
	r := GAPRecord{Type: 0x09, EventType: EventTypeInd, Data: []byte{0xab, 0xcd}}
	got := r.String()
	if len(got) == 0 {
		t.Error("empty String()")
	}
}

func TestKnownDevicesAddresses(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{StoreGAPMap: true})
	feedAdv(t, s, 0x111111111111, []byte{0x02, 0x01, 0x06})
	feedAdv(t, s, 0x222222222222, []byte{0x02, 0x01, 0x06})
	addrs := s.KnownDevicesAddresses(nil)
	if len(addrs) != 2 {
		t.Errorf("got %d addrs, want 2", len(addrs))
	}
}

// Empty AD records / zero-length records must not infinite-loop.
func TestGAPMalformed(t *testing.T) {
	s := New(nil, nil, &BLEScannerConfig{StoreGAPMap: true})
	// length=0 → handlePDU returns. length > remaining → also returns.
	feedAdv(t, s, 0x333333333333, []byte{0x00})
	feedAdv(t, s, 0x444444444444, []byte{0x10, 0x09, 0x41}) // length=16 but only 1 byte of data
	// Length = 1 means just the type byte, no data.
	feedAdv(t, s, 0x555555555555, []byte{0x01, 0x09})
}
