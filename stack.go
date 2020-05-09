package ble

import (
	"github.com/BertoldVdb/go-ble/blescanner"
	hci "github.com/BertoldVdb/go-ble/hci"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

type BluetoothStack struct {
	logger     *logrus.Entry
	Controller *hci.Controller

	BLEScanner *blescanner.BLEScanner
}

func New(logger *logrus.Entry, dev hciinterface.HCIInterface) *BluetoothStack {
	s := &BluetoothStack{
		logger: logger,
	}

	s.Controller = hci.New(bleutil.LogWithPrefix(logger, "hci"), dev,
		func() error {
			if s.BLEScanner != nil {
				return s.BLEScanner.Run()
			}

			return nil
		}, s.Close)
	s.BLEScanner = blescanner.New(bleutil.LogWithPrefix(logger, "scanner"), s.Controller, blescanner.BLEScannerConfig{StoreGAPMap: true})

	return s
}

func (s *BluetoothStack) Run() error {
	return s.Controller.Run()
}

func (s *BluetoothStack) Close() {
	s.Controller.Close()
	if s.BLEScanner != nil {
		s.BLEScanner.Close()
	}
}
