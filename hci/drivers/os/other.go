// +build !linux

package hcidriveros

import "github.com/BertoldVdb/go-ble/hci/drivers/interface"

// ListDevices returns no devices on unsupported OSes
func ListDevices() ([]string, error) {
	return nil, nil
}

// Open does nothing on unsupported OSes
func Open(deviceName string) (hciinterface.HCIInterface, error) {
	return nil, hciinterface.ErrorDeviceNotFound
}
