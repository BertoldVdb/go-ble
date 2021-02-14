package bleconnecter

import (
	"context"
	"errors"
	"net"
	"sync"

	"github.com/BertoldVdb/go-ble/bleadvertiser"
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

type BLEConnectionRole struct {
	singleChan   chan (struct{})
	peerMutex    sync.Mutex
	peerValid    bool
	peer         *BLEConnection
	peerResponse chan (struct{})
}

type BLEConnecter struct {
	logger     *logrus.Entry
	config     *BLEConnecterConfig
	ctrl       *hci.Controller
	advertiser *bleadvertiser.BLEAdvertiser

	closeflag closeflag.CloseFlag

	roles [2]BLEConnectionRole
}

type BLEConnection struct {
	*hciconnmgr.Connection

	connecter *BLEConnecter
	event     *hcievents.LEConnectionCompleteEvent

	peerAddrCandidates []bleutil.BLEAddr
	peerAddr           bleutil.BLEAddr

	isCentral   bool
	ownAddrType bleutil.MacAddrType

	parametersMutex  sync.RWMutex
	parametersActual BLEConnectionParametersActual
}

func (c *BLEConnection) LocalAddr() net.Addr {
	return c.connecter.ctrl.GetOwnAddress(c.ownAddrType)
}

func (c *BLEConnection) RemoteAddr() net.Addr {
	return c.peerAddr
}

func (c *BLEConnection) Encrypt(ediv uint16, rand uint64, ltk [16]byte) error {
	return c.connecter.ctrl.Cmds.LEEnableEncryptionSync(hcicommands.LEEnableEncryptionInput{
		ConnectionHandle:     c.event.ConnectionHandle,
		EncryptedDiversifier: ediv,
		RandomNumber:         rand,
		LongTermKey:          ltk,
	})
}

func New(logger *logrus.Entry, ctrl *hci.Controller, advertiser *bleadvertiser.BLEAdvertiser, config *BLEConnecterConfig) *BLEConnecter {
	e := &BLEConnecter{
		logger:     logger,
		config:     config,
		ctrl:       ctrl,
		advertiser: advertiser,
	}

	for i := range e.roles {
		e.roles[i] = BLEConnectionRole{
			singleChan:   make(chan (struct{}), 1),
			peerResponse: make(chan (struct{}), 1),
		}
	}

	return e
}

func (c *BLEConnecter) leConnectionCompleteHandler(event *hcievents.LEConnectionCompleteEvent) *hcievents.LEConnectionCompleteEvent {
	if event.Role > 1 {
		return event
	}

	var hwConn *hciconnmgr.Connection
	if event.Status == 0 {
		hwConn = c.ctrl.ConnMgr.ConnectionNew(event.ConnectionHandle, func() error {
			if c.advertiser == nil {
				return nil
			}
			return c.advertiser.StateChanged()
		})
	}

	remoteAddr := bleutil.BLEAddr{
		MacAddr:     event.PeerAddress,
		MacAddrType: event.PeerAddressType,
	}

	role := &c.roles[event.Role]

	role.peerMutex.Lock()
	rightPeer := false

	if role.peerValid {
		if role.peer.peerAddrCandidates == nil {
			rightPeer = true
		} else {
			for _, m := range role.peer.peerAddrCandidates {
				if m == remoteAddr {
					rightPeer = true
					break
				}
			}
		}
	}

	if rightPeer {
		role.peerValid = false
		role.peer.peerAddr = remoteAddr
		role.peer.event = event
		role.peer.Connection = hwConn
		if hwConn != nil {
			hwConn.AppConn = role.peer
		}
		role.peerResponse <- struct{}{}
	}

	role.peerMutex.Unlock()

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
	for i := range c.roles {
		c.roles[i].singleChan <- struct{}{}
	}

	<-c.closeflag.Chan()

	return nil
}

func (c *BLEConnecter) Close() error {
	return c.closeflag.Close()
}

func (c *BLEConnecter) Connect(ctx context.Context, isCentral bool, peerAddrs []bleutil.BLEAddr, request BLEConnectionParametersRequested) (*BLEConnection, []bleutil.BLEAddr, error) {
	var err error

	roleID := 0
	if !isCentral {
		roleID = 1
	}
	role := &c.roles[roleID]

	/* Ensure only one connect call can be outstanding */
	select {
	case <-c.closeflag.Chan():
		return nil, peerAddrs, ErrorClosed
	case <-ctx.Done():
		return nil, peerAddrs, ctx.Err()
	case <-role.singleChan:
	}

	defer func() { role.singleChan <- struct{}{} }()

	c.logger.WithField("0addr", peerAddrs).Debug("Starting connection")

	var ownAddrType bleutil.MacAddrType
	if isCentral {
		ownAddrType = c.ctrl.GetLERecommenedOwnAddrType(hci.LEAddrUsageConnect)
	} else {
		ownAddrType = c.ctrl.GetLERecommenedOwnAddrType(hci.LEAddrUsageAdvertise)
	}

	conn := &BLEConnection{
		connecter:          c,
		isCentral:          isCentral,
		peerAddrCandidates: peerAddrs,
		ownAddrType:        ownAddrType,
	}

	request.makeValid()

	role.peerMutex.Lock()
	role.peerValid = true
	role.peer = conn
	role.peerMutex.Unlock()

	stopPeer := func() {
		role.peerMutex.Lock()
		role.peerValid = false
		role.peerMutex.Unlock()

		/* Clear the channel in case we raced */
		select {
		case <-role.peerResponse:
		default:
		}
	}

	var advCancelFunc func() error

	if conn.isCentral {
		err = c.ctrl.Cmds.LEClearWhiteListSync()
		if err != nil {
			return nil, peerAddrs, err
		}

		for _, m := range peerAddrs {
			err = c.ctrl.Cmds.LEAddDeviceToWhiteListSync(hcicommands.LEAddDeviceToWhiteListInput{
				AddressType: m.MacAddrType,
				Address:     m.MacAddr,
			})

			if err != nil {
				return nil, peerAddrs, err
			}
		}

		param := hcicommands.LECreateConnectionInput{
			LEScanInterval:        0x10, /* Scan all the time */
			LEScanWindow:          0x10,
			InitiatorFilterPolicy: 1, /* Use allowlist */
			OwnAddressType:        conn.ownAddrType,

			ConnectionIntervalMin: request.ConnectionIntervalMin,
			ConnectionIntervalMax: request.ConnectionIntervalMax,
			ConnectionLatency:     request.ConnectionLatency,
			SupervisionTimeout:    request.SupervisionTimeout,
			MinCELength:           request.MinCELength,
			MaxCELength:           request.MaxCELength,
		}

		err = c.ctrl.Cmds.LECreateConnectionSync(param)
		if err != nil {
			stopPeer()
			return nil, peerAddrs, err
		}
	} else {
		/* Enable the advertiser in the right mode.
		   If there is only one peerAddr use directed adv.
		   Do not use allowlist in peripheral mode
		*/
		if c.advertiser != nil {
			var target *bleutil.BLEAddr
			if len(peerAddrs) == 1 {
				target = &peerAddrs[0]
			}
			advCancelFunc, err = c.advertiser.LegacyAdvertisingSetConnection(false, target)
			if err != nil {
				stopPeer()
				return nil, peerAddrs, err
			}
		}
	}

	/* Wait for the connection complete event or timeout */
	select {
	case <-role.peerResponse:
		err = nil
	case <-ctx.Done():
		err = ctx.Err()
	case <-c.closeflag.Chan():
		err = ErrorClosed
	}

	stopPeer()

	if !conn.isCentral && advCancelFunc != nil {
		advCancelFunc()
	}

	if conn.event == nil {
		c.logger.WithError(err).WithField("0addr", conn.peerAddr).Debug("No event received, cancelling")

		/* This can fail, log error but don't act on it */
		if peerAddrs != nil && conn.isCentral {
			c.ctrl.Cmds.LECreateConnectionCancelSync()
		}

		return nil, peerAddrs, err
	} else if err != nil {
		c.logger.WithError(err).WithField("0addr", conn.peerAddr).Debug("Event received, still cancelling")

		/* We got the event, but the system is closing down or it came too late */
		if conn.event.Status == 0 {
			conn.Close()
		}
		return nil, peerAddrs, err
	}

	/* Was it succesful? */
	if conn.Connection == nil {
		c.logger.WithFields(logrus.Fields{
			"0addr":   conn.peerAddr,
			"1status": conn.event.Status}).Debug("Connection failed")
		return nil, peerAddrs, hcicommands.HciErrorToGo([]byte{conn.event.Status}, nil)
	}

	/* conn contains a valid hardware conection, that didn't timeout and we are ready to use it */
	c.logger.WithFields(logrus.Fields{
		"0addr":   conn.peerAddr,
		"1handle": conn.event.ConnectionHandle}).Info("Connection established")

	conn.parametersMutex.Lock()
	conn.parametersActual = BLEConnectionParametersActual{
		Interval: conn.event.ConnectionInterval,
		Latency:  conn.event.ConnectionLatency,
		Timeout:  conn.event.SupervisionTimeout,
	}
	conn.parametersMutex.Unlock()

	if !conn.isCentral {
		conn.UpdateParams(request)
	} else {
		role.peer.peerAddrCandidates = nil
		newPeers := []bleutil.BLEAddr{}

		for _, m := range peerAddrs {
			if m != role.peer.peerAddr {
				newPeers = append(newPeers, m)
			}
		}
		peerAddrs = newPeers
	}

	return conn, peerAddrs, nil
}

func (c *BLEConnection) IsCentral() bool {
	return c.isCentral
}
