package main

import (
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
)

type EventPacket struct {
	rxPkt     hciinterface.HCIRxPacket
	eventCode int
	params    []byte
}

const (
	eventCodeCommandComplete = 0x0E
	eventCodeCommandStatus   = 0x0F
	eventCodeLEMeta          = 0x3E
)

func (s *BluetoothStack) handleEvent(pkt EventPacket) error {
	switch pkt.eventCode {
	case eventCodeCommandComplete:
		s.hcicmdmgr.HandleEventCommandComplete(pkt.params)
	case eventCodeCommandStatus:
		s.hcicmdmgr.HandleEventCommandStatus(pkt.params)
	}
	return nil
}
