package hciconnmgr

import (
	"errors"
	"sync"

	"github.com/BertoldVdb/go-misc/closeflag"
	"github.com/BertoldVdb/go-misc/multirun"
	"github.com/sirupsen/logrus"

	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	hcievents "github.com/BertoldVdb/go-ble/hci/events"
)

type ConnectionManager struct {
	sync.RWMutex
	multirun.Runnable

	closeflag closeflag.CloseFlag

	logger   *logrus.Entry
	Cmds     *hcicommands.Commands
	Events   *hcievents.EventHandler
	sendFunc func(data []byte) error

	connections map[uint16]*Connection

	txSlotManagerEDRACL *txSlotManager
	txSlotManagerEDRSDO *txSlotManager //Not used for now
	txSlotManagerLEACL  *txSlotManager

	useWg sync.WaitGroup

	useBroadcomQuirk bool
}

var (
	ErrorClosed           = errors.New("Connection manager is closed")
	ErrorConnectionClosed = errors.New("The connection is not open")
)

func New(logger *logrus.Entry, cmds *hcicommands.Commands, events *hcievents.EventHandler, sendFunc func(data []byte) error) *ConnectionManager {
	return &ConnectionManager{
		logger:   logger,
		Cmds:     cmds,
		Events:   events,
		sendFunc: sendFunc,

		connections: make(map[uint16]*Connection),
	}
}

func (c *ConnectionManager) disconnectionCompleteHandler(event *hcievents.DisconnectionCompleteEvent) *hcievents.DisconnectionCompleteEvent {
	if event.Status != 0 {
		return event
	}

	c.Lock()
	conn, ok := c.connections[event.ConnectionHandle]
	if ok {
		delete(c.connections, event.ConnectionHandle)
		conn.disconnected()
	}
	c.Unlock()

	return event
}

func (c *ConnectionManager) FindConnectionByHandle(handle uint16) *Connection {
	c.RLock()
	conn := c.connections[handle]
	c.RUnlock()

	return conn
}

func (c *ConnectionManager) closeAll() {
	var conns []*Connection
	c.RLock()
	for _, m := range c.connections {
		conns = append(conns, m)
	}
	c.RUnlock()

	for _, m := range conns {
		m.Close()
	}
}

func (c *ConnectionManager) Run() error {
	defer c.useWg.Wait()

	err := c.Events.SetDisconnectionCompleteEventCallback(c.disconnectionCompleteHandler)
	if err != nil {
		return err
	}

	defer c.closeAll()

	/* Broadcom chips seem to encode the ack message incorrectly, detect that here */
	version, err := c.Cmds.InformationalReadLocalVersionInformationSync(nil)
	if err == nil {
		if version.ManufacturerName == 15 {
			c.logger.Warn("Detected Broadcom chip: using TX quirk.")
			c.useBroadcomQuirk = true
		}
	}

	err = c.runSlotManagers()
	if err != nil {
		return err
	}

	return err
}

func (c *ConnectionManager) Close() error {
	return c.closeflag.Close()
}

func (c *ConnectionManager) HandleData(rxPkt hciinterface.HCIRxPacket) bool {
	if len(rxPkt.Data) < 1 {
		return false
	}

	if rxPkt.Data[0] == 2 {
		return c.handleACL(rxPkt.Data[1:])
	}

	return false
}
