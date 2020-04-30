package main

import (
	"log"
	"net"
	"time"

	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	hcidriverserial "github.com/BertoldVdb/go-ble/hci/drivers/serial"
)

type BluetoothStack struct {
	dev hciinterface.HCIInterface

	hcicmdmgr *hcicmdmgr.CommandManager
	cmds      *hcicommands.Commands
}

func main() {
	conn, err := net.Dial("tcp", "192.168.0.30:3000")
	//conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	clientDev, err := hcidriverserial.OpenPort(conn)
	if err != nil {
		log.Println("  Failed to make H4 protocol interface:", err)
		return
	}

	stack := BluetoothStack{
		dev: clientDev,
		hcicmdmgr: hcicmdmgr.New([]int{10, 1}, true, func(data []byte) error {
			return clientDev.SendPacket(hciinterface.HCITxPacket{Data: data})
		}),
	}
	stack.cmds = hcicommands.New(stack.hcicmdmgr)

	stack.Init()

	go func() {
		log.Fatalln(stack.Run())
	}()

	for {
		time.Sleep(time.Second)
		log.Println("Sending")
		/*cmd := hcicmdmgr.HCICommand{
			OGF: 3,
			OCF: 3,
		}

		log.Println("run async", stack.hcicmdmgr.CommandRunAsync(0, cmd, func(err error, params []byte) error {
			log.Println("Callback", params, err)
			return err
		}))
		log.Println("Post async")
		log.Println(stack.hcicmdmgr.CommandRunSync(0, cmd, nil))*/
		log.Println(stack.cmds.BasebandResetSync())
		//log.Println(stack.cmds.InformationalReadLocalVersion())

		params := hcicommands.LinkControlPeriodicInquiryModeInput{
			MaxPeriodLength: 0x3000,
			MinPeriodLength: 0x1000,
			LAP:             0x9e8b01,
			InquiryLength:   10,
			NumResponses:    0,
		}
		log.Println(stack.cmds.LinkControlPeriodicInquiryModeSync(params), params)
		log.Println(stack.cmds.InformationalReadBDADDRSync(nil))
		result, err := stack.cmds.InformationalReadLocalSupportedCommandsSync(nil)
		log.Printf("%+v %v\n", result, err)
	}
}
