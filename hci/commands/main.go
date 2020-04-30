package hcicommands

import hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"

type Commands struct {
	hcicmdmgr *hcicmdmgr.CommandManager
}

func New(hcicmdmgr *hcicmdmgr.CommandManager) *Commands {
	return &Commands{
		hcicmdmgr: hcicmdmgr,
	}
}
