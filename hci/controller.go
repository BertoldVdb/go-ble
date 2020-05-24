package hci

import (
	"encoding/hex"
	"log"

	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/BertoldVdb/go-misc/closeflag"
	"github.com/BertoldVdb/go-misc/multirun"

	"github.com/sirupsen/logrus"
)

type ReadyCallback func() error
type CloseCallback func()

type Controller struct {
	logger *logrus.Entry
	dev    hciinterface.HCIInterface
	close  closeflag.CloseFlag
	config *ControllerConfig

	Hcicmdmgr *hcicmdmgr.CommandManager
	Cmds      *hcicommands.Commands
	Events    *hcievents.EventHandler
	ConnMgr   *hciconnmgr.ConnectionManager

	Info ControllerInfo

	multirun multirun.MultiRun
}

type ControllerConfig struct {
	AwaitStartup     bool
	LERandomAddrBits int
}

func DefaultConfig() *ControllerConfig {
	return &ControllerConfig{
		AwaitStartup:     false,
		LERandomAddrBits: 24,
	}
}

type ControllerInfo struct {
	SupportedCommands   *hcicommands.InformationalReadLocalSupportedCommandsOutput
	SupportedFeatures   *hcicommands.InformationalReadLocalSupportedFeaturesOutput
	LESupportedFeatures *hcicommands.LEReadLocalSupportedFeaturesOutput
	BdAddr              *hcicommands.InformationalReadBDADDROutput
	RandomAddr          bleutil.MacAddr
}

func New(logger *logrus.Entry, dev hciinterface.HCIInterface, config *ControllerConfig) *Controller {
	c := &Controller{
		logger: logger,
		dev:    dev,
		config: config,
	}

	sendFunc := func(data []byte) error {
		if logger != nil && logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
			logger.WithFields(logrus.Fields{
				"0data": hex.EncodeToString(data),
			}).Trace("Writing HCI TX Packet to hardware")
		}

		return c.dev.SendPacket(hciinterface.HCITxPacket{Data: data})
	}

	c.Hcicmdmgr = hcicmdmgr.New(bleutil.LogWithPrefix(logger, "cmdmgr"), []int{10}, config.AwaitStartup, sendFunc)
	c.Cmds = hcicommands.New(bleutil.LogWithPrefix(logger, "cmds"), c.Hcicmdmgr)
	c.Events = hcievents.New(bleutil.LogWithPrefix(logger, "events"), c.Hcicmdmgr, c.Cmds)
	c.ConnMgr = hciconnmgr.New(bleutil.LogWithPrefix(logger, "connmgr"), c.Cmds, c.Events, sendFunc)

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

		if c.ConnMgr.HandleData(rxPkt) {
			return nil
		}

		log.Printf("Extra rx packet %+v\n", rxPkt)
		c.dev.SendPacket(hciinterface.HCITxPacket{Data: rxPkt.Data})

		return nil
	})

	c.multirun.RegisterRunnable(c.dev)
	c.multirun.RegisterRunnable(c.Hcicmdmgr)
	c.multirun.RegisterFunc(func() error {
		return c.configureDevice()
	}, func() error {
		return c.Cmds.BasebandResetSync()
	})
	c.multirun.RegisterRunnable(c.ConnMgr)

	return c
}

func (c *Controller) configureDevice() error {
	err := c.Cmds.BasebandResetSync()
	if err != nil {
		return err
	}

	/* Quirk: Should not be needed unless we support EDR as well, but some controllers
	   require it to work with extended LE commands at all */
	c.Cmds.BasebandWriteLEHostSupportSync(hcicommands.BasebandWriteLEHostSupportInput{
		LESupportedHost: 1,
	})

	/* Packet based flow control */
	c.Cmds.BasebandWriteFlowControlModeSync(hcicommands.BasebandWriteFlowControlModeInput{})

	c.Info.BdAddr, err = c.Cmds.InformationalReadBDADDRSync(nil)
	if err != nil {
		return err
	}

	/* Setup the privacy address */
	err = c.setLERandomAddress()
	if err != nil {
		return err
	}

	/*
		//TOOD: Just safekeeping the minimum advertising code to put in the right place later
		c.Cmds.LESetAdvertisingDataSync(hcicommands.LESetAdvertisingDataInput{
			AdvertisingDataLength: 3,
			AdvertisingData:       [31]byte{2, 1, 6},
		})

		c.Cmds.LESetScanResponseDataSync(hcicommands.LESetScanResponseDataInput{
			ScanResponseDataLength: 4,
			ScanResponseData:       [31]byte{3, 9, 'B', 'e'},
		})

		c.Cmds.LESetAdvertisingParametersSync(hcicommands.LESetAdvertisingParametersInput{
			AdvertisingIntervalMin:  0x20,
			AdvertisingIntervalMax:  0x40,
			AdvertisingType:         2,
			OwnAddressType:          1,
			AdvertisingChannelMap:   7,
			AdvertisingFilterPolicy: 0,
		})

		c.Cmds.LESetAdvertisingEnableSync(hcicommands.LESetAdvertisingEnableInput{
			AdvertisingEnable: 1,
		})*/

	c.Info.SupportedCommands, err = c.Cmds.InformationalReadLocalSupportedCommandsSync(nil)
	if err != nil {
		return err
	}

	c.Info.SupportedFeatures, err = c.Cmds.InformationalReadLocalSupportedFeaturesSync(nil)
	if err != nil {
		return err
	}

	c.Info.LESupportedFeatures, err = c.Cmds.LEReadLocalSupportedFeaturesSync(nil)
	return err
}

func (c *Controller) Run(ready func()) error {
	return c.multirun.Run(ready)
}

func (c *Controller) Close() error {
	return c.multirun.Close()
}
