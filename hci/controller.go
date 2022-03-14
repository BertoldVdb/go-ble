package hci

import (
	"encoding/hex"
	"sync"
	"time"

	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	deviceinfo "github.com/BertoldVdb/go-ble/hci/information"
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

	Info deviceinfo.ControllerInfo

	multirun multirun.MultiRun
}

type ControllerConfig struct {
	AwaitStartup     bool
	LERandomAddrBits int

	PrivacyConnect   bool
	PrivacyScan      bool
	PrivacyAdvertise bool

	HookInitDevice func(ctrl *Controller) error

	ConnectionManagerUsed   bool
	ConnectionManagerConfig *hciconnmgr.ConnectionManagerConfig

	WatchdogTimeout time.Duration
}

func DefaultConfig() *ControllerConfig {
	return &ControllerConfig{
		AwaitStartup:     false,
		LERandomAddrBits: 48,

		PrivacyConnect:   true,
		PrivacyScan:      true,
		PrivacyAdvertise: true,

		ConnectionManagerUsed:   true,
		ConnectionManagerConfig: hciconnmgr.DefaultConfig(),

		WatchdogTimeout: 30 * time.Second,
	}
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

	if config.ConnectionManagerUsed {
		c.ConnMgr = hciconnmgr.New(bleutil.LogWithPrefix(logger, "connmgr"), c.Cmds, c.Events, config.ConnectionManagerConfig, &c.Info, sendFunc)
	}

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

		if c.ConnMgr != nil && c.ConnMgr.HandleData(rxPkt) {
			return nil
		}

		if logger != nil && logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
			logger.WithFields(logrus.Fields{
				"0data": hex.EncodeToString(rxPkt.Data),
			}).Debug("Unsupported HCI packet received")
		}

		return nil
	})

	c.multirun.RegisterRunnable(c.dev)
	c.multirun.RegisterRunnable(c.Hcicmdmgr)
	c.multirun.RegisterFunc(func() error {
		return c.configureDevice()
	}, func() error {
		return c.Cmds.BasebandResetSync()
	})

	if config.WatchdogTimeout > 0 {
		c.multirun.RegisterFunc(func() error {
			go func() {
				var out hcicommands.InformationalReadBDADDROutput
				for {
					time.Sleep(30 * time.Second)
					_, err := c.Cmds.InformationalReadBDADDRSync(&out)
					if err != nil {
						return
					}
				}
			}()
			return nil
		}, nil)
	}

	if c.ConnMgr != nil {
		c.multirun.RegisterRunnableReady(c.ConnMgr)
	}

	return c
}

func (c *Controller) configureDevice() error {
	err := c.Cmds.BasebandResetSync()
	if err != nil {
		return err
	}

	if c.config.HookInitDevice != nil {
		err = c.config.HookInitDevice(c)
		if err != nil {
			return err
		}
	}

	err = c.Info.Read(c.Cmds)
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

	/* Setup the privacy address */
	return c.setLERandomAddress()
}

func (c *Controller) Run(ready func()) error {
	return c.multirun.Run(ready)
}

func (c *Controller) Close() error {
	return c.multirun.Close()
}
