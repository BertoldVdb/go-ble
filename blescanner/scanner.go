package blescanner

import (
	"sync"
	"time"

	"github.com/BertoldVdb/go-ble/hci"
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/BertoldVdb/go-misc/closeflag"
	"github.com/sirupsen/logrus"
)

type GAPCallback func(*BLEDevice, *GAPRecord)
type DeviceUpdatedCallback func(*BLEDevice)

type BLEScannerConfig struct {
	StoreGAPMap         bool
	ScanCycleDurationMs int
	ScanCycleActiveDuty float32

	LEScanInterval uint16
	LEScanWindow   uint16
}

type BLEScanner struct {
	sync.RWMutex
	logger *logrus.Entry
	config *BLEScannerConfig
	ctrl   *hci.Controller
	close  closeflag.CloseFlag

	devices                      map[uint64]*BLEDevice
	manufacturerSpecificCallback map[uint16]GAPCallback
	deviceUpdatedCallbacks       []DeviceUpdatedCallback
	scanType                     int

	nextCleanup time.Time
}

func New(logger *logrus.Entry, ctrl *hci.Controller, config *BLEScannerConfig) *BLEScanner {
	e := &BLEScanner{
		logger:                       logger,
		config:                       config,
		ctrl:                         ctrl,
		devices:                      make(map[uint64]*BLEDevice),
		manufacturerSpecificCallback: make(map[uint16]GAPCallback),
	}

	return e
}

func (s *BLEScanner) configureScan(scanType int, durationMs int) error {
	s.Lock()
	s.scanType = scanType
	s.Unlock()

	if s.logger != nil {
		str := ""
		switch scanType {
		case -1:
			str = "Stopping scan"
		case 1:
			str = "Starting active scan"
		case 0:
			str = "Starting passive scan"
		}
		s.logger.WithFields(logrus.Fields{
			"0scanType":   scanType,
			"1durationMs": durationMs,
		}).Info(str)
	}

	s.ctrl.Cmds.LESetScanEnableSync(hcicommands.LESetScanEnableInput{
		LEScanEnable:     0,
		FilterDuplicates: 0,
	})

	if scanType < 0 {
		return nil
	}

	params := hcicommands.LESetScanParametersInput{
		LEScanInterval:       s.config.LEScanInterval,
		LEScanWindow:         s.config.LEScanWindow,
		OwnAddressType:       s.ctrl.GetLERecommenedOwnAddrType(hci.LEAddrUsageScan),
		ScanningFilterPolicy: 0,
	}
	if scanType >= 1 {
		params.LEScanType = 1
	}
	err := s.ctrl.Cmds.LESetScanParametersSync(params)
	if err != nil {
		return err
	}

	err = s.ctrl.Cmds.LESetScanEnableSync(hcicommands.LESetScanEnableInput{
		LEScanEnable:     1,
		FilterDuplicates: 0,
	})
	return err
}

func (s *BLEScanner) Run() error {
	defer s.Close()
	defer func() {
		s.configureScan(-1, -1)
	}()

	err := s.ctrl.Events.SetLEAdvertisingReportEventCallback(s.handleScanResult)
	if err != nil {
		return err
	}

	if s.config.ScanCycleDurationMs == 0 {
		s.config.ScanCycleDurationMs = 10000
		s.config.ScanCycleActiveDuty = 0.25
	}

	if s.config.LEScanInterval == 0 {
		s.config.LEScanInterval = 64
		s.config.LEScanWindow = 12
	}

	if s.config.ScanCycleActiveDuty <= 0 {
		err = s.configureScan(0, -1)
	} else if s.config.ScanCycleActiveDuty >= 1 {
		err = s.configureScan(1, -1)
	} else {
		timer := time.NewTimer(0)
		defer timer.Stop()

		active := 1
		for {
			select {
			case <-timer.C:
			case <-s.close.Chan():
				return nil
			}

			dutycycle := s.config.ScanCycleActiveDuty
			if active == 0 {
				dutycycle = 1 - dutycycle
			}
			duration := dutycycle * float32(bleutil.RandomRange(2*s.config.ScanCycleDurationMs/3, 4*s.config.ScanCycleDurationMs/3))
			timer.Reset(time.Duration(duration) * time.Millisecond)

			err := s.configureScan(active, int(duration))
			if err != nil {
				return err
			}

			active = 1 - active
			s.handleTimeout()
		}
	}

	if err == nil {
		<-s.close.Chan()
	}

	return err
}

func (s *BLEScanner) Close() error {
	return s.close.Close()
}

func (s *BLEScanner) RegisterDeviceUpdateCallback(cb DeviceUpdatedCallback) {
	s.Lock()
	defer s.Unlock()

	s.deviceUpdatedCallbacks = append(s.deviceUpdatedCallbacks, cb)
}

func (s *BLEScanner) SetManufacturerSpecificCallback(id uint16, cb GAPCallback) {
	s.Lock()
	defer s.Unlock()

	if cb == nil {
		delete(s.manufacturerSpecificCallback, id)
	} else {
		s.manufacturerSpecificCallback[id] = cb
	}
}

func (s *BLEScanner) GetScanType() int {
	s.RLock()
	defer s.RUnlock()

	return s.scanType
}
