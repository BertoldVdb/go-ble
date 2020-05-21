package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BertoldVdb/go-ble"
	blescannerjson "github.com/BertoldVdb/go-ble/blescanner/json"
	hcidrivers "github.com/BertoldVdb/go-ble/hci/drivers"
	bleutilparam "github.com/BertoldVdb/go-ble/util/param"
	"github.com/BertoldVdb/go-misc/httplog"
	"github.com/BertoldVdb/go-misc/logrusconfig"
)

func main() {
	listenPort := flag.String("listen", "8080", "The port to listen on")
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	stack := ble.New(logger, dev)
	if stack == nil {
		logger.Fatalln("Could not make stack")
	}

	go func() {
		<-c
		stack.Close()
	}()

	loggerHTTP := logger.WithField("prefix", "http")

	go func() {
		for i := 1; i <= 10; i++ {
			time.Sleep(1 * time.Second)
			loggerHTTP.Warnf("Visit http://127.0.0.1:%s/ to see the output of this program (%d/10)", *listenPort, i)
		}
	}()

	go func() {
		traceHTTP := httplog.HTTPLog{
			LogOut:            loggerHTTP.Debugf,
			CorrelationHeader: "X-Request-ID",
			SkipInfo:          true,
		}

		http.HandleFunc("/ble/scan", blescannerjson.New(stack.BLEScanner).HTTPHandler)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write(MustAsset("html/index.html"))
		})

		loggerHTTP.Errorln(http.ListenAndServe(":"+*listenPort, traceHTTP.GetHandler(http.DefaultServeMux)))
		stack.Close()
	}()

	logger.Fatalln(stack.Run())
}
