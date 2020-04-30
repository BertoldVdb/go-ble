package hcidrivers

import (
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	hcidriveros "github.com/BertoldVdb/go-ble/hci/drivers/os"
	hcidriverserial "github.com/BertoldVdb/go-ble/hci/drivers/serial"
)

var listFunctions = [](func() ([]string, error)){
	hcidriveros.ListDevices, hcidriverserial.ListDevices,
}
var openFunctions = [](func(string) (hciinterface.HCIInterface, error)){
	hcidriveros.Open, hcidriverserial.Open,
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
