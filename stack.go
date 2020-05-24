package ble

import (
	"github.com/BertoldVdb/go-ble/bleconnecter"
	"github.com/BertoldVdb/go-ble/blescanner"
	"github.com/BertoldVdb/go-misc/multirun"

	hci "github.com/BertoldVdb/go-ble/hci"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

type BluetoothStackConfig struct {
	BLEScannerConfig   *blescanner.BLEScannerConfig
	BLEConnecterConfig *bleconnecter.BLEConnecterConfig

	HCIControllerConfig *hci.ControllerConfig
}

func DefaultConfig() *BluetoothStackConfig {
	return &BluetoothStackConfig{
		BLEScannerConfig:    &blescanner.BLEScannerConfig{StoreGAPMap: true},
		BLEConnecterConfig:  &bleconnecter.BLEConnecterConfig{},
		HCIControllerConfig: hci.DefaultConfig(),
	}
}

type BluetoothStack struct {
	logger   *logrus.Entry
	config   *BluetoothStackConfig
	multirun multirun.MultiRun

	Controller   *hci.Controller
	BLEScanner   *blescanner.BLEScanner
	BLEConnecter *bleconnecter.BLEConnecter
}

func New(logger *logrus.Entry, config *BluetoothStackConfig, dev hciinterface.HCIInterface) *BluetoothStack {
	if config == nil {
		config = DefaultConfig()
	}

	s := &BluetoothStack{
		config: config,
		logger: logger,
	}

	s.Controller = hci.New(bleutil.LogWithPrefix(logger, "hci"), dev, s.config.HCIControllerConfig)
	s.BLEScanner = blescanner.New(bleutil.LogWithPrefix(logger, "scanner"), s.Controller, config.BLEScannerConfig)
	s.BLEConnecter = bleconnecter.New(bleutil.LogWithPrefix(logger, "connecter"), s.Controller, config.BLEConnecterConfig)

	s.multirun.RegisterRunnableReady(s.Controller)
	s.multirun.RegisterRunnable(s.BLEScanner)
	s.multirun.RegisterRunnable(s.BLEConnecter)

	return s
}

func (s *BluetoothStack) Run(ready func()) error {
	return s.multirun.Run(ready)
}

func (s *BluetoothStack) Close() error {
	return s.multirun.Close()
}
