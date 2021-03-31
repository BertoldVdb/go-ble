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
	} else {
		s := dev.scanner
		s.Lock()
		for i := 0; i < len(s.deviceUpdatedCallbacks); i++ {
			cb = s.deviceUpdatedCallbacks[i]
			s.Unlock()
			if cb != nil {
				cb(dev)
			}
			s.Lock()
		}
		s.Unlock()
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

		s.Lock()
		for _, m := range s.advertisingReportCallbacks {
			m(&pkt)
		}
		s.Unlock()

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
		dev.signalUpdatedCallbacks()

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
	}
	return ad
}
