package ble

import (
	"github.com/BertoldVdb/go-ble/bleadvertiser"
	"github.com/BertoldVdb/go-ble/bleconnecter"
	"github.com/BertoldVdb/go-ble/blescanner"
	"github.com/BertoldVdb/go-ble/blesmp"
	"github.com/BertoldVdb/go-misc/multirun"

	hci "github.com/BertoldVdb/go-ble/hci"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

type BluetoothStackConfig struct {
	BLEScannerUse       bool
	BLEScannerConfig    *blescanner.BLEScannerConfig
	BLEConnecterUse     bool
	BLEConnecterConfig  *bleconnecter.BLEConnecterConfig
	BLEAdvertiserUse    bool
	BLEAdvertiserConfig *bleadvertiser.BLEAdvertiserConfig
	BLEWatchdogUse      bool

	SMPConfig           *blesmp.SMPConfig
	HCIControllerConfig *hci.ControllerConfig
}

func DefaultConfig() *BluetoothStackConfig {
	return &BluetoothStackConfig{
		BLEScannerUse:       true,
		BLEScannerConfig:    &blescanner.BLEScannerConfig{StoreGAPMap: true},
		BLEConnecterUse:     true,
		BLEConnecterConfig:  &bleconnecter.BLEConnecterConfig{},
		BLEAdvertiserUse:    true,
		BLEAdvertiserConfig: bleadvertiser.DefaultConfig(),
		SMPConfig:           blesmp.DefaultConfig(),
		HCIControllerConfig: hci.DefaultConfig(),
	}
}

type BluetoothStack struct {
	logger   *logrus.Entry
	config   *BluetoothStackConfig
	multirun multirun.MultiRun

	Controller    *hci.Controller
	BLEScanner    *blescanner.BLEScanner
	BLEAdvertiser *bleadvertiser.BLEAdvertiser
	BLEConnecter  *bleconnecter.BLEConnecter
	SMP           *blesmp.SMP
}

func New(logger *logrus.Entry, config *BluetoothStackConfig, dev hciinterface.HCIInterface) *BluetoothStack {
	if config == nil {
		config = DefaultConfig()
	}

	/* Fill in any sub-configs the caller left as nil. Each component's
	   constructor dereferences its config (HCI controller reads
	   AwaitStartup at controller.go:90, scanner uses
	   ScanCycleDurationMs at scanner.go:117, advertiser reads
	   DeviceName at legacy.go:262, connecter reads
	   BLEUpdateParametersVerify at params.go:31, SMP reads
	   StoredKeysPath at smp.go:147). Without this defaulting, a
	   caller passing a partially-populated BluetoothStackConfig
	   nil-panicked deep inside a sub-component instead of getting a
	   working stack. */
	if config.HCIControllerConfig == nil {
		config.HCIControllerConfig = hci.DefaultConfig()
	}
	if config.BLEScannerConfig == nil {
		config.BLEScannerConfig = &blescanner.BLEScannerConfig{}
	}
	if config.BLEAdvertiserConfig == nil {
		config.BLEAdvertiserConfig = bleadvertiser.DefaultConfig()
	}
	if config.BLEConnecterConfig == nil {
		config.BLEConnecterConfig = &bleconnecter.BLEConnecterConfig{}
	}
	if config.SMPConfig == nil {
		config.SMPConfig = blesmp.DefaultConfig()
	}

	s := &BluetoothStack{
		config: config,
		logger: logger,
	}

	s.Controller = hci.New(bleutil.LogWithPrefix(logger, "hci"), dev, s.config.HCIControllerConfig)
	s.BLEScanner = blescanner.New(bleutil.LogWithPrefix(logger, "scanner"), s.Controller, config.BLEScannerConfig)
	s.BLEAdvertiser = bleadvertiser.New(bleutil.LogWithPrefix(logger, "advertiser"), s.Controller, config.BLEAdvertiserConfig)
	s.BLEConnecter = bleconnecter.New(bleutil.LogWithPrefix(logger, "connecter"), s.Controller, s.BLEAdvertiser, config.BLEConnecterConfig)

	if s.Controller.ConnMgr != nil {
		s.SMP = blesmp.New(bleutil.LogWithPrefix(logger, "smp"), s.Controller, config.SMPConfig)
	}

	s.multirun.RegisterRunnableReady(s.Controller)
	if config.BLEScannerUse {
		s.multirun.RegisterRunnable(s.BLEScanner)
	}
	if config.BLEAdvertiserUse {
		s.multirun.RegisterRunnable(s.BLEAdvertiser)
	}
	if config.BLEConnecterUse {
		s.multirun.RegisterRunnable(s.BLEConnecter)
	}

	return s
}

func (s *BluetoothStack) Run(ready func()) error {
	return s.multirun.Run(ready)
}

func (s *BluetoothStack) Close() error {
	return s.multirun.Close()
}
