package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"

	hcidrivers "github.com/BertoldVdb/go-ble/drivers"
	hciinterface "github.com/BertoldVdb/go-ble/drivers/interface"
	serialhci "github.com/BertoldVdb/go-ble/drivers/serial"
)

type config struct {
	listenAddr *string
	deviceName *string

	connectCommand    *string
	disconnectCommand *string
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

	clientDev, err := serialhci.OpenPort(currentConn)
	if err != nil {
		log.Println("  Failed to make H4 protocol interface:", err)
		return
	}

	hwDev, err := hcidrivers.Open(*config.deviceName)
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
			log.Println("  Client interface closed:", err)
		}
	}()
	go func() {
		defer handlerWg.Done()
		defer clientDev.Close()

		err := hwDev.Run()
		if err != nil {
			log.Println("  Hardware interface closed:", err)
		}
	}()

	handlerWg.Wait()
}

func main() {
	config := config{}

	config.listenAddr = flag.String("listen", ":3000", "The address to listen on")
	config.deviceName = flag.String("device", "", "The device to use. If not specified, a list of all devices is printed")
	config.connectCommand = flag.String("connect", "", "Command to execute upon connection")
	config.disconnectCommand = flag.String("disconnect", "", "Command to execute upon disconnection")

	flag.Parse()

	if *config.deviceName == "" {
		fmt.Println("No device specified (--device).")
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
