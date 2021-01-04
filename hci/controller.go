package hci

import (
	"encoding/hex"
	"log"
	"sync"

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
	logger   *logrus.Entry
	devMutex sync.Mutex
	dev      hciinterface.HCIInterface
	close    closeflag.CloseFlag
	config   *ControllerConfig

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

	PrivacyConnect   bool
	PrivacyScan      bool
	PrivacyAdvertise bool
}

func DefaultConfig() *ControllerConfig {
	return &ControllerConfig{
		AwaitStartup:     false,
		LERandomAddrBits: 24,

		PrivacyConnect:   true,
		PrivacyScan:      true,
		PrivacyAdvertise: true,
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

		c.devMutex.Lock()
		defer c.devMutex.Unlock()
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
	c.multirun.RegisterRunnableReady(c.ConnMgr)

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
