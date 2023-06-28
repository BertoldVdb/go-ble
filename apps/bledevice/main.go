package main

import (
	"flag"
	"io"
	"net"
	"runtime"

	"github.com/BertoldVdb/go-ble"
	attperipheral "github.com/BertoldVdb/go-ble/bleatt/helpers/peripheral"
	hcidrivers "github.com/BertoldVdb/go-ble/hci/drivers"
	"github.com/BertoldVdb/go-ble/hci/drivers/btsnoop"
	bleutil "github.com/BertoldVdb/go-ble/util"
	bleutilparam "github.com/BertoldVdb/go-ble/util/param"
	"github.com/BertoldVdb/go-misc/logrusconfig"
	"github.com/BertoldVdb/go-misc/multirun"

	servicedeviceinformation "github.com/BertoldVdb/go-ble/bleatt/service/deviceinformation"
	serviceserial "github.com/BertoldVdb/go-ble/bleatt/service/serial"
)

/* Main function */
func main() {
	secure := flag.Bool("secure", false, "Use encryption")
	deviceName := flag.String("name", "SerialExample", "The beacon name")
	serviceUUIDString := flag.String("uuid", "6e400001-b5a3-f393-e0a9-e50e24dcca9e", "Service UUID to use")
	destination := flag.String("destination", "", "Address to connect to. When unset, return received data")
	multiple := flag.Bool("multi", false, "Try to accept multiple connections")
	logfile := flag.String("btsnoop", "", "Write btsnoop file to path")

	logrusconfig.InitParam()
	bleutilparam.Init()

	flag.Parse()

	logger := logrusconfig.GetLogger(0)

	serviceUUID, err := bleutil.UUIDFromString(*serviceUUIDString)
	if err != nil {
		logger.Fatalln(err)
	}

	devName, err := bleutilparam.GetDeviceName()
	if err != nil {
		logger.Fatalln(err)
	}

	dev, err := hcidrivers.Open(devName)
	if err != nil {
		logger.Fatalln(err)
	}

	dev, err = btsnoop.WrapFile(dev, *logfile)
	if err != nil {
		logger.Fatalln(err)
	}

	m := multirun.MultiRun{}
	m.HandleSIGTERM()

	config := ble.DefaultConfig()
	config.BLEScannerUse = false
	config.BLEAdvertiserConfig.DeviceName = *deviceName
	config.BLEAdvertiserConfig.DeviceService = serviceUUID
	config.BLEAdvertiserConfig.LegacyBaseIntervalMin = 40
	config.BLEAdvertiserConfig.LegacyBaseIntervalMax = 200
	config.HCIControllerConfig.PrivacyAdvertise = false
	config.SMPConfig.DefaultConnConfig.AuthReq = 0

	stack := ble.New(logger, config, dev)
	if stack == nil {
		logger.Fatalln("Could not make stack")
	}

	m.RegisterRunnableReady(stack)

	peripheralConfig := attperipheral.DefaultConfig()
	peripheralConfig.Appearance = 0
	peripheralConfig.DeviceName = *deviceName
	peripheralConfig.ConnectionParams.ConnectionIntervalMin = 6
	peripheralConfig.ConnectionParams.ConnectionIntervalMax = 12
	peripheralConfig.ConnectionParams.ConnectionLatency = 1
	peripheralConfig.AcceptMultipleConnections = *multiple
	peripheralHelper := attperipheral.New(stack, peripheralConfig)

	deviceInfoConfig := servicedeviceinformation.DefaultConfig()
	deviceInfoConfig.ManufacturerName = "go-ble"
	deviceInfoConfig.ModelNumber = "serial example"
	deviceInfoConfig.FirmwareRevision = runtime.Version()

	peripheralHelper.RegisterImplementation(servicedeviceinformation.CreateService(deviceInfoConfig))

	serialConfig := serviceserial.DefaultConfig()
	serialConfig.ServiceUUID = serviceUUID
	serialConfig.Secure = *secure
	serialConfig.Connect = func() (io.ReadWriteCloser, error) {
		if *destination != "" {
			return net.Dial("tcp", *destination)
		}

		c1, c2 := net.Pipe()
		go io.Copy(c2, c2)

		return c1, nil
	}

	peripheralHelper.RegisterImplementation(serviceserial.CreateService(serialConfig))

	m.RegisterRunnable(peripheralHelper)

	logger.Fatalln(m.Run(nil))
}
