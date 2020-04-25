package hcidrivers

import (
	hciinterface "github.com/BertoldVdb/go-ble/drivers/interface"
	oshci "github.com/BertoldVdb/go-ble/drivers/os"
	serialhci "github.com/BertoldVdb/go-ble/drivers/serial"
)

var listFunctions = [](func() ([]string, error)){
	oshci.ListDevices, serialhci.ListDevices,
}
var openFunctions = [](func(string) (hciinterface.HCIInterface, error)){
	oshci.Open, serialhci.Open,
}

// ListDevices returns all found HCI devices
func ListDevices() ([]string, error) {
	var result []string

	for _, f := range listFunctions {
		tmp, _ := f()
		result = append(result, tmp...)
	}

	return result, nil
}

// Open does nothing on unsupported OSes
func Open(deviceName string) (hciinterface.HCIInterface, error) {
	for _, f := range openFunctions {
		tmp, err := f(deviceName)
		if err == nil {
			return tmp, err
		}
		if err != hciinterface.ErrorDeviceNotFound {
			return nil, err
		}
	}

	return nil, hciinterface.ErrorDeviceNotFound
}
