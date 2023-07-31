package main

import (
	"flag"
	"io"
	"log"
	"net"
	"time"

	"github.com/BertoldVdb/go-ble"
	attcentral "github.com/BertoldVdb/go-ble/bleatt/helpers/central"
	serviceserial "github.com/BertoldVdb/go-ble/bleatt/service/serial"
	hcidrivers "github.com/BertoldVdb/go-ble/hci/drivers"
	"github.com/BertoldVdb/go-ble/hci/drivers/btsnoop"
	bleutil "github.com/BertoldVdb/go-ble/util"
	bleutilparam "github.com/BertoldVdb/go-ble/util/param"
	"github.com/BertoldVdb/go-misc/logrusconfig"
	"github.com/BertoldVdb/go-misc/multirun"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Entry

type ServerApp struct {
	central  *attcentral.CentralHelper
	lst      net.Listener
	peerAddr bleutil.BLEAddr
	serial   serviceserial.SerialConfig
}

func (s *ServerApp) Run(ready func()) error {
	ready()

	for {
		conn, err := s.lst.Accept()
		if err != nil {
			return err
		}

		/* Close connection if there is already one */
		s.central.PeerRemoveAddr(s.peerAddr)

		/* Make a new connection */
		s.central.PeerAdd(s.peerAddr, false, time.Now().Add(5*time.Second), s.serial.ClientFactory(conn, func(conn io.ReadWriteCloser) {
			if conn == nil {
				logger.Info("Callback: connection closed")
			} else {
				logger.Info("Callback: connection open")
			}
		}))

	}
}

func (s *ServerApp) Close() error {
	return s.lst.Close()
}

func main() {
	bleaddr := flag.String("addr", "A4:C1:37:32:98:C3", "Device address")
	lstaddr := flag.String("lst", ":8899", "Address to listen for TCP connections")
	serviceUUIDString := flag.String("uuid", "6e400001-b5a3-f393-e0a9-e50e24dcca9e", "Service UUID to use")
	rdUUIDString := flag.String("rduuid", "", "Read UUID to use")
	wrUUIDString := flag.String("wruuid", "", "Write UUID to use")
	logfile := flag.String("btsnoop", "", "Write btsnoop file to path")
	bleutilparam.Init()
	logrusconfig.InitParam()
	flag.Parse()

	logger = logrusconfig.GetLogger(0)

	lst, err := net.Listen("tcp", *lstaddr)
	if err != nil {
		logger.Fatalln(err)
	}

	devName, err := bleutilparam.GetDeviceName()
	if err != nil {
		logger.Fatalln(err)
	}

	dev, err := hcidrivers.Open(devName)
	if err != nil {
		log.Fatalln(err)
	}

	if *logfile != "" {
		dev, err = btsnoop.WrapFile(dev, *logfile)
		if err != nil {
			logger.Fatalln(err)
		}
	}

	config := ble.DefaultConfig()
	config.HCIControllerConfig.PrivacyConnect = false
	config.BLEScannerUse = false

	stack := ble.New(logger, config, dev)
	if stack == nil {
		logger.Fatalln("Could not make stack")
	}

	central := attcentral.New(stack, attcentral.DefaultConfig())

	m := multirun.MultiRun{}
	m.HandleSIGTERM()

	m.RegisterRunnableReady(stack)
	m.RegisterRunnable(central)

	app := &ServerApp{
		serial: serviceserial.SerialConfig{
			ServiceUUID: bleutil.UUIDFromStringPanic(*serviceUUIDString),
		},
		peerAddr: bleutil.BLEAddr{
			MacAddr:     bleutil.MacAddrFromStringPanic(*bleaddr),
			MacAddrType: 0,
		},
		central: central,
		lst:     lst,
	}

	if *rdUUIDString != "" {
		app.serial.ReadUUID = bleutil.UUIDFromStringPanic(*rdUUIDString)
	}
	if *wrUUIDString != "" {
		app.serial.WriteUUID = bleutil.UUIDFromStringPanic(*wrUUIDString)
	}

	m.RegisterRunnableReady(app)

	logger.Fatalln(m.Run(func() {
		logger.Info("Ready!")
	}))
}
