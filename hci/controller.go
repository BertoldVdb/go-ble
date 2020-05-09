package hci

import (
	"encoding/hex"
	"log"
	"sync"

	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/BertoldVdb/go-misc/closeflag"

	"github.com/sirupsen/logrus"
)

type ReadyCallback func() error
type CloseCallback func()

type Controller struct {
	logger *logrus.Entry
	dev    hciinterface.HCIInterface
	close  closeflag.CloseFlag

	Hcicmdmgr *hcicmdmgr.CommandManager
	Cmds      *hcicommands.Commands
	Events    *hcievents.EventHandler

	Info ControllerInfo

	cbReady ReadyCallback
	cbClose CloseCallback
}

type ControllerInfo struct {
	SupportedCommands   *hcicommands.InformationalReadLocalSupportedCommandsOutput
	SupportedFeatures   *hcicommands.InformationalReadLocalSupportedFeaturesOutput
	LESupportedFeatures *hcicommands.LEReadLocalSupportedFeaturesOutput
	BdAddr              *hcicommands.InformationalReadBDADDROutput
}

func New(logger *logrus.Entry, dev hciinterface.HCIInterface, cbReady ReadyCallback, cbClose CloseCallback) *Controller {
	c := &Controller{
		logger:  logger,
		dev:     dev,
		cbReady: cbReady,
		cbClose: cbClose,
	}

	c.Hcicmdmgr = hcicmdmgr.New(bleutil.LogWithPrefix(logger, "cmdmgr"), []int{10}, true, func(data []byte) error {
		if logger != nil && logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
			logger.WithFields(logrus.Fields{
				"0data": hex.EncodeToString(data),
			}).Trace("Writing HCI TX Packet to hardware")
		}

		return c.dev.SendPacket(hciinterface.HCITxPacket{Data: data})
	})
	c.Cmds = hcicommands.New(bleutil.LogWithPrefix(logger, "cmds"), c.Hcicmdmgr)
	c.Events = hcievents.New(bleutil.LogWithPrefix(logger, "events"), c.Hcicmdmgr, c.Cmds)

	c.dev.SetRecvHandler(func(rxPkt hciinterface.HCIRxPacket) error {
		if rxPkt.Received == false {
			return nil
		}

		if logger != nil && logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
			logger.WithFields(logrus.Fields{
				"0data": hex.EncodeToString(rxPkt.Data),
			}).Trace("Received HCI RX packet from hardware")
		}

		if c.Events.HandleEvent(rxPkt) {
			return nil
		}

		log.Printf("Extra rx packet %+v\n", rxPkt)

		return nil
	})

	return c
}

func (c *Controller) Run() error {
	var wg sync.WaitGroup
	var errDevice, errCmdMgr error

	wg.Add(2)
	go func() {
		defer c.Close()
		errDevice = c.dev.Run()
		wg.Done()
	}()
	go func() {
		defer c.Close()
		errCmdMgr = c.Hcicmdmgr.Run()
		wg.Done()
	}()

	errCmd := c.Cmds.BasebandResetSync()
	if errCmd != nil {
		c.Close()
	}

	c.Info.BdAddr, errCmd = c.Cmds.InformationalReadBDADDRSync(nil)
	if errCmd != nil {
		c.Close()
	}

	c.Info.SupportedCommands, errCmd = c.Cmds.InformationalReadLocalSupportedCommandsSync(nil)
	if errCmd != nil {
		c.Close()
	}

	c.Info.SupportedFeatures, errCmd = c.Cmds.InformationalReadLocalSupportedFeaturesSync(nil)
	if errCmd != nil {
		c.Close()
	}

	c.Info.LESupportedFeatures, errCmd = c.Cmds.LEReadLocalSupportedFeaturesSync(nil)
	if errCmd != nil {
		c.Close()
	}

	if errCmd == nil && c.cbReady != nil {
		wg.Add(1)
		go func() {
			defer c.Close()
			errCmd = c.cbReady()
			wg.Done()
		}()
	}

	wg.Wait()

	if errCmd != nil {
		return errCmd
	} else if errDevice != nil {
		return errDevice
	}
	return errCmdMgr
}

func (c *Controller) Close() error {
	if c.close.Close() == nil {
		c.dev.Close()
		c.Hcicmdmgr.Close()

		if c.cbClose != nil {
			c.cbClose()
		}
	}

	return nil
}
