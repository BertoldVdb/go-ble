package main

import (
	"context"
	"encoding/hex"
	"flag"
	"log"

	"github.com/BertoldVdb/go-ble"
	"github.com/BertoldVdb/go-ble/bleconnecter"
	hcidrivers "github.com/BertoldVdb/go-ble/hci/drivers"
	bleutil "github.com/BertoldVdb/go-ble/util"
	bleutilparam "github.com/BertoldVdb/go-ble/util/param"
	"github.com/BertoldVdb/go-misc/logrusconfig"
	"github.com/BertoldVdb/go-misc/multirun"
)

func main() {
	logrusconfig.InitParam()
	bleutilparam.Init()

	flag.Parse()

	deviceName, err := bleutilparam.GetDeviceName()
	if err != nil {
		return
	}

	logger := logrusconfig.GetLogger(0)

	dev, err := hcidrivers.Open(deviceName)
	if err != nil {
		logger.Fatalln(err)
	}

	m := multirun.MultiRun{}
	m.HandleSIGTERM()

	stack := ble.New(logger, nil, dev)
	if stack == nil {
		logger.Fatalln("Could not make stack")
	}

	m.RegisterRunnableReady(stack)

	/*E3:16:D4:3E:6C:AD SmartSolar
	  EA:EF:FA:BD:F2:C2 SmartBatterySense HQ18279HYBG
	  EA:33:20:2B:5F:27 VE.Direct Smart
	*/

	logger.Fatalln(m.Run(func() {
		peerAddr := bleutil.BLEAddr{
			MacAddr:     bleutil.MacAddrFromStringPanic("E3:16:D4:3E:6C:AD"),
			MacAddrType: 1,
		}

		conn, err := stack.BLEConnecter.Connect(context.Background(), peerAddr, bleconnecter.BLEConnectionParametersRequested{})
		if err != nil {
			log.Fatalln(err)
		}

		var rxBuf [1024]byte
		for {
			n, err := conn.HwConn.Read(rxBuf[:])
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(hex.EncodeToString(rxBuf[:n]))
		}
		//conn.Close()
	}))
}
