package bleconnecter

import (
	"context"
	"errors"
	"sync"

	"github.com/BertoldVdb/go-ble/hci"
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/BertoldVdb/go-misc/closeflag"
	"github.com/sirupsen/logrus"
)

var (
	ErrorClosed = errors.New("Connecter is closed")
)

type BLEConnecterConfig struct {
}

type BLEConnecter struct {
	logger *logrus.Entry
	config *BLEConnecterConfig
	ctrl   *hci.Controller

	closeflag closeflag.CloseFlag

	connectionSingleChan   chan (struct{})
	connectionPeerMutex    sync.Mutex
	connectionPeerValid    bool
	connectionPeer         *BLEConnection
	connectionPeerResponse chan (struct{})
}

type BLEConnection struct {
	connecter *BLEConnecter
	event     *hcievents.LEConnectionCompleteEvent

	HwConn   *hciconnmgr.Connection
	PeerAddr bleutil.BLEAddr
	IsMaster bool

	parametersMutex  sync.RWMutex
	parametersActual BLEConnectionParametersActual
}

func New(logger *logrus.Entry, ctrl *hci.Controller, config *BLEConnecterConfig) *BLEConnecter {
	e := &BLEConnecter{
		logger:                 logger,
		config:                 config,
		ctrl:                   ctrl,
		connectionSingleChan:   make(chan (struct{}), 1),
		connectionPeerResponse: make(chan (struct{}), 1),
	}

	return e
}

func (c *BLEConnecter) leConnectionCompleteHandler(event *hcievents.LEConnectionCompleteEvent) *hcievents.LEConnectionCompleteEvent {
	var hwConn *hciconnmgr.Connection
	if event.Status == 0 {
		hwConn = c.ctrl.ConnMgr.ConnectionNew(event.ConnectionHandle)
	}

	remoteAddr := bleutil.BLEAddr{
		MacAddr:     event.PeerAddress,
		MacAddrType: event.PeerAddressType,
	}

	c.connectionPeerMutex.Lock()
	rightPeer := c.connectionPeerValid &&
		((event.Role == 0 && c.connectionPeer.IsMaster && c.connectionPeer.PeerAddr == remoteAddr) ||
			(event.Role == 1 && !c.connectionPeer.IsMaster))

	if rightPeer {
		c.connectionPeerValid = false
		c.connectionPeer.PeerAddr = remoteAddr
		c.connectionPeer.event = event
		c.connectionPeer.HwConn = hwConn
		hwConn.AppConn = c.connectionPeer
		c.connectionPeerResponse <- struct{}{}
	}
	c.connectionPeerMutex.Unlock()

	if !rightPeer {
		c.logger.WithField("0event", event).Debug("Received event for wrong target address")

		/* If we tried to cancel this connection but raced against the complete event, close it.
		   Can't run blocking things here so we just make a goroutine as this should be a very rare event */
		if hwConn != nil {
			go hwConn.Close()
		}
		return event

	} else {
		c.logger.WithField("0event", event).Debug("Received event for desired target address")
	}

	return nil
}

func (c *BLEConnecter) Run() error {
	defer c.Close()

	err := c.ctrl.Events.SetLEConnectionCompleteEventCallback(c.leConnectionCompleteHandler)
	if err != nil {
		return err
	}

	/* Failure of these two is harmless */
	c.ctrl.Events.SetLERemoteConnectionParameterRequestEventCallback(c.leConnectionParameterRequestHandler)
	c.ctrl.Events.SetLEConnectionUpdateCompleteEventCallback(c.leConnectionUpdateCompleteHandler)

	/* Activate the connect calls */
	c.connectionSingleChan <- struct{}{}

	<-c.closeflag.Chan()

	return nil
}

func (c *BLEConnecter) Close() error {
	return c.closeflag.Close()
}

func (c *BLEConnecter) Connect(ctx context.Context, peerAddr bleutil.BLEAddr, request BLEConnectionParametersRequested) (*BLEConnection, error) {
	var err error

	/* Ensure only one connect call can be outstanding */
	select {
	case <-c.closeflag.Chan():
		return nil, ErrorClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.connectionSingleChan:
	}

	defer func() { c.connectionSingleChan <- struct{}{} }()

	c.logger.WithField("0addr", peerAddr).Debug("Starting connection")

	conn := &BLEConnection{
		connecter: c,
		PeerAddr:  peerAddr,
		IsMaster:  peerAddr.MacAddr != 0,
	}

	request.makeValid()

	c.connectionPeerMutex.Lock()
	c.connectionPeerValid = true
	c.connectionPeer = conn
	c.connectionPeerMutex.Unlock()

	if conn.IsMaster {
		param := hcicommands.LECreateConnectionInput{
			LEScanInterval:        0x10, /* Scan all the time */
			LEScanWindow:          0x10,
			InitiatorFilterPolicy: 0,
			PeerAddressType:       peerAddr.MacAddrType,
			PeerAddress:           peerAddr.MacAddr,
			OwnAddressType:        c.ctrl.GetLERecommenedOwnAddrType(hci.LEAddrUsageConnect),

			ConnectionIntervalMin: request.ConnectionIntervalMin,
			ConnectionIntervalMax: request.ConnectionIntervalMax,
			ConnectionLatency:     request.ConnectionLatency,
			SupervisionTimeout:    request.SupervisionTimeout,
			MinCELength:           request.MinCELength,
			MaxCELength:           request.MaxCELength,
		}

		err = c.ctrl.Cmds.LECreateConnectionSync(param)
		if err != nil {
			return nil, err
		}
	} else {
		//TODO: Set advertising to connectable, or enable it if not on at all.
		//TODO use defer to reset advertising state
	}

	/* Wait for the connection complete event or timeout */
	select {
	case <-c.connectionPeerResponse:
		err = nil
	case <-ctx.Done():
		err = ctx.Err()
	case <-c.closeflag.Chan():
		err = ErrorClosed
	}

	c.connectionPeerMutex.Lock()
	c.connectionPeerValid = false
	c.connectionPeerMutex.Unlock()

	/* Clear the channel in case we raced */
	select {
	case <-c.connectionPeerResponse:
	default:
	}

	if conn.event == nil {
		c.logger.WithError(err).WithField("0addr", conn.PeerAddr).Debug("No event received, cancelling")

		/* This can fail, log error but don't act on it */
		if conn.IsMaster {
			c.ctrl.Cmds.LECreateConnectionCancelSync()
		} else {
			//TODO: Do something with the advertiser
		}
		return nil, err
	} else if err != nil {
		c.logger.WithError(err).WithField("0addr", conn.PeerAddr).Debug("Event received, still cancelling")

		/* We got the event, but the system is closing down or it came too late */
		if conn.event.Status == 0 {
			conn.Close()
		}
		return nil, err
	}

	/* Was it succesful? */
	if conn.HwConn == nil {
		c.logger.WithFields(logrus.Fields{
			"0addr":   conn.PeerAddr,
			"1status": conn.event.Status}).Debug("Connection failed")
		return nil, hcicommands.HciErrorToGo([]byte{conn.event.Status}, nil)
	}

	/* conn contains a valid hardware conection, that didn't timeout and we are ready to use it */
	c.logger.WithFields(logrus.Fields{
		"0addr":   conn.PeerAddr,
		"1handle": conn.event.ConnectionHandle}).Info("Connection established")

	conn.parametersMutex.Lock()
	conn.parametersActual = BLEConnectionParametersActual{
		Interval: conn.event.ConnectionInterval,
		Latency:  conn.event.ConnectionLatency,
		Timeout:  conn.event.SupervisionTimeout,
	}
	conn.parametersMutex.Unlock()

	if !conn.IsMaster {
		conn.UpdateParams(request)
	}

	return conn, nil
}

func (c *BLEConnection) Close() error {
	return c.HwConn.Close()
}
