package blescanner

import (
	"time"

	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

type BLEAdvertisingReport struct {
	Addr    bleutil.BLEAddr
	RSSI    int8
	PktType EventType
	Data    []byte
}

func (dev *BLEDevice) handlePDU(event EventType, data []byte) {
	for {
		if len(data) < 2 {
			return
		}

		recordLen := int(data[0])
		if 1+recordLen > len(data) {
			return
		}
		record := data[1:(1 + recordLen)]
		data = data[1+recordLen:]

		dev.handleRecord(event, record)
	}
}

func (dev *BLEDevice) signalUpdatedCallbacks() {
	dev.cbMutex.Lock()
	cb := dev.cb
	dev.cbMutex.Unlock()

	if cb != nil {
		cb(dev)
		return
	}

	/* Snapshot the global callback list so user callbacks may safely
	   register/unregister callbacks (the previous index-based iteration
	   could miss or repeat entries when the slice mutated under it). */
	s := dev.scanner
	s.Lock()
	cbs := append([]registeredDeviceUpdateCB(nil), s.deviceUpdatedCallbacks...)
	s.Unlock()

	for _, e := range cbs {
		if e.cb != nil {
			e.cb(dev)
		}
	}
}

func (s *BLEScanner) handleScanResult(ad *hcievents.LEAdvertisingReportEvent) *hcievents.LEAdvertisingReportEvent {
	now := time.Now()

	for i := 0; i < int(ad.NumReports); i++ {
		bleaddr := bleutil.BLEAddr{
			MacAddr:     ad.Address[i],
			MacAddrType: ad.AddressType[i],
		}

		event := EventType(ad.EventType[i])

		pkt := BLEAdvertisingReport{
			Addr:    bleaddr,
			RSSI:    int8(ad.RSSI[i]),
			PktType: event,
			Data:    ad.Data[i],
		}

		skip := false

		/* Snapshot the callback list so callbacks can safely re-enter
		   scanner methods (the previous design held s.Lock() across the
		   callback call, which would deadlock if the user touched the
		   scanner from inside the callback). */
		s.Lock()
		cbs := append([]registeredAdvReportCB(nil), s.advertisingReportCallbacks...)
		s.Unlock()

		for _, m := range cbs {
			if skip = m.cb(&pkt); skip {
				break
			}
		}
		if skip {
			continue
		}

		dev, isNew := s.getDevice(bleaddr, true)
		if dev == nil {
			continue
		}

		dev.Lock()
		dev.lastSeenDev = now
		if event == EventTypeInd || event == EventTypeDirectInd {
			dev.lastConnectable = now
		}
		dev.rssi = int8(ad.RSSI[i])
		dev.handlePDU(event, ad.Data[i])

		if s.logger != nil {
			if isNew {
				s.logger.WithFields(logrus.Fields{
					"0addr": dev.addr,
					"1rssi": dev.rssi,
				}).Info("Found new device")

			} else if s.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
				s.logger.WithFields(logrus.Fields{
					"0addr": dev.addr,
					"1rssi": dev.rssi,
				}).Trace("Device updated")
			}
		}

		dev.Unlock()

		/* Invoke user-update callbacks after the device write lock has
		   been released; otherwise, callbacks that call self-locking
		   getters on the same goroutine would deadlock. */
		dev.signalUpdatedCallbacks()
	}
	return ad
}
