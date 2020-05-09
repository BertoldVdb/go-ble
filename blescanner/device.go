package blescanner

import (
	"fmt"
	"sort"
	"sync"
	"time"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

type BLEDevice struct {
	scanner  *BLEScanner
	lastSeen time.Time

	cbMutex sync.Mutex
	cb      DeviceUpdatedCallback

	/* All next fields are safe to access with the mutex */
	sync.RWMutex
	addr            bleutil.BLEAddr
	rssi            int8
	name            string
	nameState       int
	flags           uint8
	txPower         int8
	lastConnectable time.Time
	lastSeenDev     time.Time

	gapFields map[uint8]*GAPRecord
	gapSingle *GAPRecord

	services [3][]bleutil.UUID
}

func (d *BLEDevice) StringSummary() string {
	txPower := ""
	if d.txPower > -128 {
		txPower = fmt.Sprintf("txPower=%ddBm ", d.txPower)
	}

	result := fmt.Sprintf("Addr='%s' Name='%s' Flags=%02x Connectable=%v RSSI=%ddBm %sLastSeen=%dms",
		d.addr, d.name, d.flags, d.IsConnectable(), d.rssi, txPower, time.Now().Sub(d.lastSeenDev).Milliseconds())
	if len(d.services[0]) > 0 {
		result += fmt.Sprintf("\nUUID-16 services:   %v", d.services[0])
	}
	if len(d.services[1]) > 0 {
		result += fmt.Sprintf("\nUUID-32 services:   %v", d.services[1])
	}
	if len(d.services[2]) > 0 {
		result += fmt.Sprintf("\nUUID-128 services:  %v", d.services[2])
	}

	if d.gapFields != nil {
		var types sort.IntSlice
		for i := range d.gapFields {
			types = append(types, int(i))
		}
		types.Sort()
		for _, i := range types {
			m := d.gapFields[uint8(i)]
			result += fmt.Sprintf("\n%02X: %s", i, m)
		}
	} else {
		result += fmt.Sprintf("\n%02X: %s", d.gapSingle.Type, d.gapSingle)
	}

	return result
}

func (s *BLEScanner) getDevice(addr bleutil.BLEAddr, create bool) (*BLEDevice, bool) {
	if create {
		s.Lock()
		defer s.Unlock()
	} else {
		s.RLock()
		defer s.RUnlock()
	}

	now := time.Now()
	key := addr.GetUint64()

	device, ok := s.devices[key]
	if ok {
		if create {
			device.lastSeen = now
		}

		return device, false
	}

	/* Avoid DoS by someone spamming random devices */
	if !create || len(s.devices) >= 256 {
		return nil, false
	}

	device = &BLEDevice{
		addr:        addr,
		scanner:     s,
		name:        "Unknown",
		txPower:     -128,
		lastSeen:    now,
		lastSeenDev: now,
	}

	if device.scanner.config.StoreGAPMap {
		device.gapFields = make(map[uint8]*GAPRecord)
	} else {
		device.gapSingle = &GAPRecord{}
	}

	s.devices[key] = device

	return device, true
}

func (s *BLEScanner) handleTimeout() {
	s.Lock()
	defer s.Unlock()

	now := time.Now()

	for i, m := range s.devices {
		if now.After(m.lastSeen.Add(30 * time.Second)) {
			delete(s.devices, i)

			if s.logger != nil {
				s.logger.WithField("0addr", m.addr).Info("Device expired")
			}
		}
	}
}

func (s *BLEScanner) GetDevice(addr bleutil.BLEAddr) *BLEDevice {
	dev, _ := s.getDevice(addr, false)
	if dev != nil {
		dev.RLock()
	}
	return dev
}

func (dev *BLEDevice) Release() {
	dev.RUnlock()
}

func (s *BLEScanner) KnownDevicesAddresses() []bleutil.BLEAddr {
	s.RLock()
	defer s.RUnlock()

	var result []bleutil.BLEAddr

	for _, m := range s.devices {
		result = append(result, m.addr)
	}

	return result
}

func (dev *BLEDevice) IsConnectable() bool {
	return !time.Now().After(dev.lastConnectable.Add(2 * time.Second))
}

func (dev *BLEDevice) SetUpdateCallback(cb DeviceUpdatedCallback) {
	dev.cbMutex.Lock()
	dev.cb = cb
	dev.cbMutex.Unlock()
}
