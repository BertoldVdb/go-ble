package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BertoldVdb/go-ble"
	hcidriverserial "github.com/BertoldVdb/go-ble/hci/drivers/serial"
	prefixed "github.com/BertoldVdb/logrus-prefixed-formatter"
	"github.com/sirupsen/logrus"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	logrus.ErrorKey = " error"
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	customFormatter := new(prefixed.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	customFormatter.PrefixPadding = 20
	//  customFormatter.PadLevelText = true
	customFormatter.SpacePadding = 50
	logger.SetFormatter(customFormatter)
	entry := logrus.NewEntry(logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	conn, err := net.Dial("tcp", "192.168.0.30:3001")
	//conn, err := net.Dial("tcp", "192.168.0.23:3000")
	//conn, err := net.Dial("tcp", "127.0.0.1:3001")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	dev, err := hcidriverserial.OpenPort(conn)
	if err != nil {
		log.Println("  Failed to make H4 protocol interface:", err)
		return
	}

	stack := ble.New(entry, dev)
	if stack == nil {
		log.Fatalln("Could not make stack")
	}

	go func() {
		<-c
		stack.Close()
		//os.Exit(1)
	}()

	go func() {
		log.Fatalln(stack.Run())
	}()

	for {
		time.Sleep(500 * time.Millisecond)
		fmt.Fprintf(os.Stderr, "\033[32m\033[H\033[2JTime: %s\n\n%s\n", time.Now(), stack.BLEScanner.StringSummary())
	}
}
