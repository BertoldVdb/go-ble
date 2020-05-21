package main

import (
	"flag"
	"log"
	"net"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	hcidrivers "github.com/BertoldVdb/go-ble/hci/drivers"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	hcidriverserial "github.com/BertoldVdb/go-ble/hci/drivers/serial"
	bleutilparam "github.com/BertoldVdb/go-ble/util/param"
)

type config struct {
	listenAddr *string
	deviceName string

	connectCommand    *string
	disconnectCommand *string

	resetEvent *bool
}

func runShell(cmd string) error {
	if cmd == "" {
		return nil
	}
	return exec.Command("/bin/sh", "-c", cmd).Run()
}

func clientHandler(parentWg *sync.WaitGroup, currentConn net.Conn, config *config) {
	parentWg.Add(1)
	defer parentWg.Done()

	log.Println("Accepted connection from", currentConn.RemoteAddr())
	defer func() {
		log.Println("Closed connection from", currentConn.RemoteAddr())
		currentConn.Close()
	}()

	runShell(*config.connectCommand)
	defer runShell(*config.disconnectCommand)

	clientDev, err := hcidriverserial.OpenPort(currentConn)
	if err != nil {
		log.Println("  Failed to make H4 protocol interface:", err)
		return
	}

	if *config.resetEvent {
		resetPkt := []byte{0x04, 0x0E, 0x03, 0x01, 0x00, 0x00}
		if clientDev.SendPacket(hciinterface.HCITxPacket{Data: resetPkt}) != nil {
			return
		}
	}

	hwDev, err := hcidrivers.Open(config.deviceName)
	if err != nil {
		log.Println("  Failed to open hardware interface:", err)
		return
	}

	var clientRx, hwRx uint64

	clientDev.SetRecvHandler(func(pkt hciinterface.HCIRxPacket) error {
		atomic.AddUint64(&clientRx, 1)
		return hwDev.SendPacket(hciinterface.HCITxPacket{Data: pkt.Data})
	})

	hwDev.SetRecvHandler(func(pkt hciinterface.HCIRxPacket) error {
		atomic.AddUint64(&hwRx, 1)
		return clientDev.SendPacket(hciinterface.HCITxPacket{Data: pkt.Data})
	})

	var doneFlag uint64
	go func() {
		for {
			time.Sleep(time.Second)
			if atomic.LoadUint64(&doneFlag) > 0 {
				return
			}
			log.Printf("  Forwarded %d messages from client and %d messages from hardware.", atomic.LoadUint64(&clientRx), atomic.LoadUint64(&hwRx))
		}
	}()
	defer atomic.AddUint64(&doneFlag, 1)

	var handlerWg sync.WaitGroup
	handlerWg.Add(2)

	go func() {
		defer handlerWg.Done()
		defer hwDev.Close()

		err := clientDev.Run()
		if err != nil {
			log.Println("  Client interface stopped:", err)
		}
	}()
	go func() {
		defer handlerWg.Done()
		defer clientDev.Close()

		err := hwDev.Run()
		if err != nil {
			log.Println("  Hardware interface stopped:", err)
		}
	}()

	handlerWg.Wait()
}

func main() {
	config := config{}
	config.listenAddr = flag.String("listen", ":3000", "The address to listen on")
	config.connectCommand = flag.String("connect", "", "Command to execute upon connection")
	config.disconnectCommand = flag.String("disconnect", "", "Command to execute upon disconnection")
	config.resetEvent = flag.Bool("reset", false, "Send reset complete event on connect")

	bleutilparam.Init()

	flag.Parse()

	var err error
	config.deviceName, err = bleutilparam.GetDeviceName()
	if err != nil {
		return
	}

	l, err := net.Listen("tcp", *config.listenAddr)
	if err != nil {
		log.Fatalln("Listen failed", err)
		return
	}
	defer l.Close()

	log.Println("Listening on", l.Addr())

	var conn net.Conn
	var wg sync.WaitGroup
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatalln("Accept failed", err)
		}

		if conn != nil {
			conn.Close()
			wg.Wait()
		}

		conn = c
		go clientHandler(&wg, conn, &config)
	}
}
