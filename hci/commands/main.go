package hcicommands

import (
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
	"github.com/sirupsen/logrus"
)

type Commands struct {
	logger    *logrus.Entry
	hcicmdmgr *hcicmdmgr.CommandManager
}

func New(logger *logrus.Entry, hcicmdmgr *hcicmdmgr.CommandManager) *Commands {
	return &Commands{
		logger:    logger,
		hcicmdmgr: hcicmdmgr,
	}
}
