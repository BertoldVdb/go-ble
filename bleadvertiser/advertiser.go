package bleadvertiser

import (
	"sync"

	bleutil "github.com/BertoldVdb/go-ble/util"

	"github.com/BertoldVdb/go-ble/hci"
	"github.com/BertoldVdb/go-misc/closeflag"
	"github.com/sirupsen/logrus"
)

type BLEAdvertiser struct {
	logger *logrus.Entry
	ctrl   *hci.Controller
	config *BLEAdvertiserConfig

	closeflag closeflag.CloseFlag

	legacyAdvertisingMutex      sync.Mutex
	legacyAdvertisingInitOnce   sync.Once
	legacyAdvertisingAlwaysOn   bool
	legacyAdvertisingData       []byte
	legacyAdvertisingScanData   []byte
	legacyAdvertisingUpdateChan chan (int)

	legacyAdvertisingSlots    []LegacyAdvertisingSlot
	legacyAdvertisingBaseSlot *LegacyAdvertisingSlot
}

type BLEAdvertiserConfig struct {
	AlwaysAdvertising bool
	DeviceName        string
	DeviceService     bleutil.UUID
	DeviceFlags       uint8

	LegacyBaseIntervalMin uint16
	LegacyBaseIntervalMax uint16
}

func DefaultConfig() *BLEAdvertiserConfig {
	return &BLEAdvertiserConfig{
		AlwaysAdvertising: true,
		DeviceName:        "go-ble device",
		DeviceService:     bleutil.UUIDBase,
		DeviceFlags:       6,
	}
}

func New(logger *logrus.Entry, ctrl *hci.Controller, config *BLEAdvertiserConfig) *BLEAdvertiser {
	a := &BLEAdvertiser{
		logger: logger,
		ctrl:   ctrl,
		config: config,

		legacyAdvertisingUpdateChan: make(chan (int), 1),
	}

	return a
}

func (a *BLEAdvertiser) Run() error {
	defer a.Close()

	return a.legacyAdvertisingManager()
}

func (a *BLEAdvertiser) Close() error {
	return a.closeflag.Close()
}

func (a *BLEAdvertiser) StateChanged() error {
	select {
	case a.legacyAdvertisingUpdateChan <- -1:
	default:
	}
	return nil
}
