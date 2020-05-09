package hcievents

import (
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	hciconst "github.com/BertoldVdb/go-ble/hci/const"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	"github.com/sirupsen/logrus"
)

func (e *EventHandler) HandleEvent(pkt hciinterface.HCIRxPacket) bool {
	if len(pkt.Data) < 3 || pkt.Data[0] != hciconst.MsgTypeEvent {
		return false
	}

	eventCode := int(pkt.Data[1])
	subEvent := 0
	paramLen := int(pkt.Data[2])
	params := pkt.Data[3:]

	if paramLen != len(params) {
		return false

	}

	if eventCode == 0x3E {
		if len(params) == 0 {
			return false
		}

		subEvent = int(params[0])
	}

	e.handleEventInternal(uint16(eventCode<<8|subEvent), params)
	return true
}

func New(logger *logrus.Entry, hcicmdmgr *hcicmdmgr.CommandManager, cmds *hcicommands.Commands) *EventHandler {
	e := &EventHandler{
		logger:    logger,
		hcicmdmgr: hcicmdmgr,
		cmds:      cmds,
	}

	/* Enable events for the HCI command manager if we use it */
	if hcicmdmgr != nil {
		e.SetCommandCompleteEventCallback(func(result *CommandCompleteEvent) *CommandCompleteEvent {
			hcicmdmgr.HandleEventCommandComplete(result.CommandOpcode, result.NumHCICommandPackets, result.ReturnParameters)
			return result
		})
		e.SetCommandStatusEventCallback(func(result *CommandStatusEvent) *CommandStatusEvent {
			hcicmdmgr.HandleEventCommandStatus(result.CommandOpcode, result.NumHCICommandPackets, result.Status)
			return result
		})
	}

	return e
}
