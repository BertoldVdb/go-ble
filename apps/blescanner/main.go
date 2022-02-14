package main

import (
	"flag"
	"net/http"

	"github.com/BertoldVdb/go-ble"
	blescannerjson "github.com/BertoldVdb/go-ble/blescanner/json"
	hcidrivers "github.com/BertoldVdb/go-ble/hci/drivers"
	bleutilparam "github.com/BertoldVdb/go-ble/util/param"
	"github.com/BertoldVdb/go-misc/logrusconfig"
	"github.com/BertoldVdb/go-misc/multirun"
	"github.com/BertoldVdb/go-misc/multirunhttp"
)

func main() {
	listenPort := flag.Int("listen", 8080, "The port to listen on")
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

	config := ble.DefaultConfig()
	config.BLEScannerUse = true
	config.BLEScannerConfig.LEScanInterval = 64
	config.BLEScannerConfig.LEScanWindow = 64

	stack := ble.New(logger, config, dev)
	if stack == nil {
		logger.Fatalln("Could not make stack")
	}

	http.HandleFunc("/ble/scan", blescannerjson.New(stack.BLEScanner).HTTPHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(MustAsset("html/index.html"))
	})

	m.RegisterRunnableReady(stack)
	m.RegisterRunnable(&multirunhttp.MultiRunHTTP{
		Server:     &http.Server{},
		LoggerHTTP: logger.WithField("prefix", "http"),
		ListenPort: *listenPort,
	})
	logger.Fatalln(m.Run(nil))
}
