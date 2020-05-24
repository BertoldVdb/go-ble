package bleutilparam

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"strings"

	hcidrivers "github.com/BertoldVdb/go-ble/hci/drivers"
	hcidriverserial "github.com/BertoldVdb/go-ble/hci/drivers/serial"
)

var (
	ErrorNoDevice   = errors.New("No device specified")
	ErrorOutOfRange = errors.New("Device with specified index was not given")
)

var (
	deviceName *string
	netAddr    *string
)

func Init() {
	deviceName = flag.String("device", "", "The device to use. If not specified, a list of all devices is printed")
	netAddr = flag.String("netaddr", "", "If specified, adds a network address to the list of devices (called net)")
}

func parseNetAddr() {
	if *netAddr != "" {
		addrs := strings.Split(*netAddr, ",")
		name := "net"
		for i, addr := range addrs {
			hcidriverserial.RegisterDevice(name, func(name string) (io.ReadWriteCloser, error) {
				return net.Dial("tcp", addr)
			})
			name = fmt.Sprintf("net%d", i+1)
		}
		*netAddr = ""
	}
}

func GetDeviceNameMulti(index int) (string, error) {
	parseNetAddr()

	if *deviceName == "" {
		fmt.Println("No device specified (-device).")
		fmt.Println("Listing all known devices:")
		devices, err := hcidrivers.ListDevices()
		if err != nil || len(devices) == 0 {
			fmt.Println("  No devices found")
		} else {
			for i, m := range devices {
				fmt.Printf("  Device %d: %s", i, m)
				fmt.Println()
			}
		}
		return "", ErrorNoDevice
	}

	deviceNames := strings.Split(*deviceName, ",")
	if index >= len(deviceNames) {
		return "", ErrorOutOfRange
	}

	return deviceNames[index], nil
}

func GetDeviceName() (string, error) {
	return GetDeviceNameMulti(0)
}
