package hci

import (
	"time"

	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	"github.com/BertoldVdb/go-misc/closeflag"
)

type hciKeepAlive struct {
	closeFlag closeflag.CloseFlag
	ctrl         *Controller
}

func (k *hciKeepAlive) Run(readyCb func()) error {
	readyCb()

	for {
		var out hcicommands.InformationalReadBDADDROutput
		for {
			select {
			case <-k.closeFlag.Chan():
				return nil

			case <-time.After(k.ctrl.config.WatchdogTimeout):
			}

			_, err := k.ctrl.Cmds.InformationalReadBDADDRSync(&out)
			if err != nil {
				return err
			}
		}
	}
}

func (k *hciKeepAlive) Close() error {
	return k.closeFlag.Close()
}
