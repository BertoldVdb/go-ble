package blescanner

import (
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

func (s *BLEScanner) getDevice(addr bleutil.BLEAddr, create bool) (*BLEDevice, bool) {
	s.handleTimeout()

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
	now := time.Now()

	s.RLock()
	needClean := now.After(s.nextCleanup)
	s.RUnlock()

	if needClean {
		s.Lock()
		defer s.Unlock()

		s.nextCleanup = now.Add(15 * time.Second)
		if s.logger != nil {
			s.logger.WithField("0numDevices", len(s.devices)).Debug("Cleaning expired devices")
		}

		for i, m := range s.devices {
			if m.isExpired() {
				delete(s.devices, i)

				if s.logger != nil {
					s.logger.WithField("0addr", m.addr).Info("Device removed due to expiry")
				}
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

func (s *BLEScanner) KnownDevicesAddresses(result []bleutil.BLEAddr) []bleutil.BLEAddr {
	s.handleTimeout()

	s.RLock()
	defer s.RUnlock()

	result = result[:0]

	for _, m := range s.devices {
		if !m.isExpired() {
			result = append(result, m.addr)
		}
	}

	return result
}

/* Warning: this function needs to be called with the scanner lock, not the
 * device lock */
func (dev *BLEDevice) isExpired() bool {
	return time.Now().After(dev.lastSeen.Add(30 * time.Second))
}

func (dev *BLEDevice) LastSeen() time.Time {
	return dev.lastSeenDev
}

func (dev *BLEDevice) GetRSSI() int8 {
	return dev.rssi
}

func (dev *BLEDevice) IsConnectable() bool {
	return !time.Now().After(dev.lastConnectable.Add(2 * time.Second))
}

func (dev *BLEDevice) SetUpdateCallback(cb DeviceUpdatedCallback) {
	dev.cbMutex.Lock()
	dev.cb = cb
	dev.cbMutex.Unlock()
}
