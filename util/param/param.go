package bleutilparam

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"

	hcidrivers "github.com/BertoldVdb/go-ble/hci/drivers"
	hcidriverserial "github.com/BertoldVdb/go-ble/hci/drivers/serial"
)

var (
	ErrorNoDevice = errors.New("No device specified")
)

var (
	deviceName *string
	netAddr    *string
)

func Init() {
	deviceName = flag.String("device", "", "The device to use. If not specified, a list of all devices is printed")
	netAddr = flag.String("netaddr", "", "If specified, adds a network address to the list of devices (called net)")
}

func GetDeviceName() (string, error) {
	if *netAddr != "" {
		hcidriverserial.RegisterDevice("net", func(name string) (io.ReadWriteCloser, error) {
			return net.Dial("tcp", *netAddr)
		})
	}

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

	return *deviceName, nil
}
