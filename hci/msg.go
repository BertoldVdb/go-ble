package main

import (
	hciconst "github.com/BertoldVdb/go-ble/hci/const"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
)

func (s *BluetoothStack) handleMessage(pkt hciinterface.HCIRxPacket) error {
	if len(pkt.Data) < 1 {
		return nil
	}

	msgType := pkt.Data[0]
	data := pkt.Data[1:]
	switch msgType {
	case hciconst.MsgTypeCommand:
		//This should not be received from the controller
	case 2:
	case 3:
		//This is used only in EDR mode
	case hciconst.MsgTypeEvent:
		if len(data) >= 2 {
			eventCode := int(data[0])
			paramLen := int(data[1])
			if paramLen+2 == len(data) {
				s.handleEvent(EventPacket{
					rxPkt:     pkt,
					params:    data[2:],
					eventCode: eventCode,
				})
			}
		}
	case 5:
	}

	return nil
}
