package main

import (
	"log"

	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
)

func (s *BluetoothStack) Init() error {
	return nil
}

func (s *BluetoothStack) Run() error {
	s.dev.SetRecvHandler(func(rxPkt hciinterface.HCIRxPacket) error {
		if !rxPkt.Received {
			return nil
		}

		//log.Println("Packet", hex.EncodeToString(rxPkt.Data))
		s.handleMessage(rxPkt)
		return nil
	})

	go func() {
		log.Println(s.hcicmdmgr.Run())
		s.dev.Close()
	}()
	return s.dev.Run()
}
