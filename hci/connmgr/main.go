package hciconnmgr

import (
	"errors"
	"sync"
	"time"

	"github.com/BertoldVdb/go-misc/closeflag"
	"github.com/BertoldVdb/go-misc/multirun"
	"github.com/sirupsen/logrus"

	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	hciinterface "github.com/BertoldVdb/go-ble/hci/drivers/interface"
	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	deviceinfo "github.com/BertoldVdb/go-ble/hci/information"
)

type ConnectionManagerConfig struct {
	HookConnectionStateChange func(c *Connection, open bool)

	// DisconnectTimeout bounds how long Connection.Close blocks waiting
	// for the controller's Disconnection-Complete event. Zero means use
	// the default (5 seconds). After the timeout elapses the host
	// completes a synthetic local teardown so shutdown isn't held up by
	// an unresponsive controller.
	DisconnectTimeout time.Duration
}

type ConnectionMangerEventsSMP struct {
	LEEncryptionGetKey func(conn *Connection, event *hcievents.LELongTermKeyRequestEvent) ([]byte, *hcievents.LELongTermKeyRequestEvent)
	EncryptionChanged  func(conn *Connection, event *hcievents.EncryptionChangeEvent) *hcievents.EncryptionChangeEvent
	EncryptionRefresh  func(conn *Connection, event *hcievents.EncryptionKeyRefreshCompleteEvent) *hcievents.EncryptionKeyRefreshCompleteEvent
}

type ConnectionManager struct {
	sync.RWMutex
	multirun.Runnable

	config *ConnectionManagerConfig

	closeflag closeflag.CloseFlag

	logger   *logrus.Entry
	Cmds     *hcicommands.Commands
	Events   *hcievents.EventHandler
	info     *deviceinfo.ControllerInfo
	sendFunc func(data []byte) error

	connections map[uint16]*Connection

	txSlotManagerEDRACL *txSlotManager
	txSlotManagerEDRSDO *txSlotManager //Not used for now
	txSlotManagerLEACL  *txSlotManager

	txNewConnBlockTime time.Time

	useWg sync.WaitGroup

	useBroadcomQuirk bool

	cb ConnectionMangerEventsSMP

	/* Bounded queue for HCI replies issued in response to peer events
	   (LTK requests, etc). The reply worker drains this serially in
	   Run(); handlers that need to send a sync command back to the
	   controller post a closure here instead of spawning an unbounded
	   `go func()`. */
	replyCh chan func() error
}

var (
	ErrorClosed           = errors.New("Connection manager is closed")
	ErrorConnectionClosed = errors.New("The connection is not open")
	ErrorNoTXSlotManager  = errors.New("No TX slot manager is available")
)

func DefaultConfig() *ConnectionManagerConfig {
	return &ConnectionManagerConfig{
		DisconnectTimeout: 5 * time.Second,
	}
}

func New(logger *logrus.Entry, cmds *hcicommands.Commands, events *hcievents.EventHandler, config *ConnectionManagerConfig, info *deviceinfo.ControllerInfo, sendFunc func(data []byte) error) *ConnectionManager {
	return &ConnectionManager{
		config: config,

		logger:   logger,
		Cmds:     cmds,
		Events:   events,
		sendFunc: sendFunc,
		info:     info,

		connections: make(map[uint16]*Connection),
		replyCh:     make(chan func() error, 16),
	}
}

func (c *ConnectionManager) SetEventsSMP(cb ConnectionMangerEventsSMP) {
	c.Lock()
	defer c.Unlock()

	c.cb = cb
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

func (c *ConnectionManager) Run(readyCb func()) error {
	defer c.useWg.Wait()
	defer c.closeAll()

	err := c.Events.SetDisconnectionCompleteEventCallback(c.disconnectionCompleteHandler)
	if err != nil {
		return err
	}

	err = c.Events.SetEncryptionChangeEventCallback(c.encryptionChangeHandler)
	if err != nil {
		return err
	}

	err = c.Events.SetEncryptionKeyRefreshCompleteEventCallback(c.encryptionKeyRefreshHandler)
	if err != nil {
		return err
	}

	err = c.Events.SetLELongTermKeyRequestEventCallback(c.encryptionLELongTermKeyRequestHandler)
	if err != nil {
		return err
	}

	/* Broadcom chips seem to encode the ack message incorrectly, detect that here */
	if c.info.LocalVersionInformation.ManufacturerName == 15 {
		c.logger.Warn("Detected Broadcom chip: using TX quirk.")
		c.useBroadcomQuirk = true
	}

	/* Reply worker — drains the bounded queue used by event handlers
	   to issue sync HCI replies (LTK request, etc) without spawning
	   unbounded goroutines. */
	replyDone := make(chan struct{})
	go func() {
		defer close(replyDone)
		for {
			select {
			case <-c.closeflag.Chan():
				return
			case fn := <-c.replyCh:
				if err := fn(); err != nil {
					c.logger.WithError(err).Debug("ConnectionManager reply HCI command failed")
				}
			}
		}
	}()
	defer func() { <-replyDone }()

	err = c.runSlotManagers(readyCb)
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
