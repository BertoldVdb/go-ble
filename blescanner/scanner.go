package blescanner

import (
	"sort"
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
}

type BLEScanner struct {
	sync.RWMutex
	logger *logrus.Entry
	config BLEScannerConfig
	ctrl   *hci.Controller
	close  closeflag.CloseFlag

	devices                      map[uint64]*BLEDevice
	manufacturerSpecificCallback map[uint16]GAPCallback
	deviceUpdatedCallbacks       []DeviceUpdatedCallback
}

func New(logger *logrus.Entry, ctrl *hci.Controller, config BLEScannerConfig) *BLEScanner {
	e := &BLEScanner{
		logger:                       logger,
		config:                       config,
		ctrl:                         ctrl,
		devices:                      make(map[uint64]*BLEDevice),
		manufacturerSpecificCallback: make(map[uint16]GAPCallback),
	}

	return e
}

func (s *BLEScanner) configureScan(active bool) error {
	if s.logger != nil {
		if active {
			s.logger.Info("Starting active scan")
		} else {
			s.logger.Info("Starting passive scan")
		}
	}

	s.ctrl.Cmds.LESetScanEnableSync(hcicommands.LESetScanEnableInput{
		LEScanEnable:     0,
		FilterDuplicates: 0,
	})

	params := hcicommands.LESetScanParametersInput{
		LEScanInterval:       16,
		LEScanWindow:         16,
		OwnAddressType:       0,
		ScanningFilterPolicy: 0,
	}
	if active {
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

	err := s.ctrl.Events.SetLEAdvertisingReportEventCallback(s.handleScanResult)
	if err != nil {
		return err
	}

	if s.config.ScanCycleDurationMs == 0 {
		s.config.ScanCycleDurationMs = 10000
		s.config.ScanCycleActiveDuty = 0.25
	}

	if s.config.ScanCycleActiveDuty <= 0 {
		err = s.configureScan(false)
	} else if s.config.ScanCycleActiveDuty >= 1 {
		err = s.configureScan(true)
	} else {
		timer := time.NewTimer(0)
		defer timer.Stop()
		active := true
		for {
			select {
			case <-timer.C:
			case <-s.close.Chan():
				return nil
			}

			dutycycle := s.config.ScanCycleActiveDuty
			if !active {
				dutycycle = 1 - dutycycle
			}
			durationCycle := dutycycle * float32(bleutil.RandomRange(2*s.config.ScanCycleDurationMs/3, 4*s.config.ScanCycleDurationMs/3))
			timer.Reset(time.Duration(durationCycle) * time.Millisecond)

			err := s.configureScan(active)
			if err != nil {
				return err
			}

			active = !active
		}
	}

	return err
}

func (s *BLEScanner) Close() {
	s.close.Close()
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

func (s *BLEScanner) StringSummary() string {
	devices := s.KnownDevicesAddresses()
	sort.SliceStable(devices, func(i, j int) bool {
		return devices[i].IsLess(devices[j])
	})

	result := ""
	for _, m := range devices {
		dev := s.GetDevice(m)
		if dev != nil {
			result += dev.StringSummary() + "\n\n"
			dev.Release()
		}
	}

	return result
}
