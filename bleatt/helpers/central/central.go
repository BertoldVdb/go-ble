package attcentral

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/BertoldVdb/go-ble"
	"github.com/BertoldVdb/go-ble/bleatt"
	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	"github.com/BertoldVdb/go-ble/bleconnecter"
	"github.com/BertoldVdb/go-ble/blesmp"
	blel2cap "github.com/BertoldVdb/go-ble/l2cap"
	bleutil "github.com/BertoldVdb/go-ble/util"
)

type CentralHelperConfig struct {
	ConnectionParametersRequested bleconnecter.BLEConnectionParametersRequested
	GATTConfig                    *bleatt.GattDeviceConfig
	SMPConnConfig                 *blesmp.SMPConnConfig
	L2CAPConfig                   *blel2cap.L2CAPConfig
}

func DefaultConfig() CentralHelperConfig {
	return CentralHelperConfig{
		ConnectionParametersRequested: bleconnecter.BLEConnectionParametersRequested{
			ConnectionIntervalMin: 6,
			ConnectionIntervalMax: 15,
			ConnectionLatency:     20,
			SupervisionTimeout:    100,
		},
		GATTConfig: bleatt.DefaultConfig(),
	}
}

type Peer struct {
	addr      bleutil.BLEAddr
	reconnect bool
	deadline  time.Time
	connect   bool
	handler   DevHandler
	conn      *bleconnecter.BLEConnection
}

type CentralHelper struct {
	sync.Mutex

	stack  *ble.BluetoothStack
	config CentralHelperConfig

	ctx    context.Context
	cancel context.CancelFunc

	desiredPeersUpdated chan (struct{})
	desiredPeers        map[uint64]*Peer
	cancelConnect       context.CancelFunc
}

func New(stack *ble.BluetoothStack, config CentralHelperConfig) *CentralHelper {
	ctx, cancel := context.WithCancel(context.Background())

	return &CentralHelper{
		stack:  stack,
		config: config,

		ctx:    ctx,
		cancel: cancel,

		desiredPeersUpdated: make(chan struct{}, 1),
		desiredPeers:        make(map[uint64]*Peer),
	}
}

var ErrorPeerAlreadyExists = errors.New("peer already exists")

func (p *CentralHelper) peersUpdated() {
	select {
	case p.desiredPeersUpdated <- struct{}{}:
	default:
	}

	if p.cancelConnect != nil {
		p.cancelConnect()
	}
}

type DevHandler func(ctx context.Context, dev *bleatt.GattDevice)

func (p *CentralHelper) PeerAdd(addr bleutil.BLEAddr, reconnect bool, deadline time.Time, handler DevHandler) (*Peer, error) {
	p.Lock()
	defer p.Unlock()

	key := addr.GetUint64()
	if p.desiredPeers[key] != nil {
		return nil, ErrorPeerAlreadyExists
	}

	peer := &Peer{
		addr:      addr,
		reconnect: reconnect,
		deadline:  deadline,
		connect:   true,
		handler:   handler,
	}

	p.desiredPeers[key] = peer
	p.peersUpdated()
	return peer, nil
}

func (p *CentralHelper) peerRemove(peer *Peer) {
	if peer == nil {
		return
	} else if peer.conn != nil {
		peer.conn.Close()
		peer.conn = nil
	}
	peer.handler(context.Background(), nil)

	key := peer.addr.GetUint64()
	if p.desiredPeers[key] == peer {
		delete(p.desiredPeers, key)
	}
	p.peersUpdated()

}

func (p *CentralHelper) PeerRemoveAddr(addr bleutil.BLEAddr) {
	p.Lock()
	defer p.Unlock()

	p.peerRemove(p.desiredPeers[addr.GetUint64()])
}

func (p *CentralHelper) PeerRemove(peer *Peer) {
	p.Lock()
	defer p.Unlock()

	p.peerRemove(peer)
}

func (p *CentralHelper) Run() error {
	for {
		if err := p.ctx.Err(); err != nil {
			return err
		}

		p.Lock()
		var deadline time.Time
		var connectlist []bleutil.BLEAddr
		now := time.Now()
		for _, peer := range p.desiredPeers {
			if peer.conn != nil || !peer.connect {
				continue
			}

			if !peer.deadline.IsZero() {
				if now.After(peer.deadline) {
					p.peerRemove(peer)
					continue
				}
				if deadline.IsZero() || peer.deadline.Before(deadline) {
					deadline = peer.deadline
				}
			}
			connectlist = append(connectlist, peer.addr)
		}

		var connCtx context.Context
		if len(connectlist) > 0 {
			if !deadline.IsZero() {
				connCtx, p.cancelConnect = context.WithDeadline(p.ctx, deadline)
			} else {
				connCtx, p.cancelConnect = context.WithCancel(p.ctx)
			}
		}
		p.Unlock()

		if len(connectlist) == 0 {
			select {
			case <-p.ctx.Done():
			case <-p.desiredPeersUpdated:
			}
			continue
		}

		conn, _, err := p.stack.BLEConnecter.Connect(connCtx, true, connectlist, p.config.ConnectionParametersRequested)
		p.cancelConnect()

		if err != nil {
			continue
		}

		p.Lock()
		key := conn.RemoteAddr().(bleutil.BLEAddr).GetUint64()
		if peer := p.desiredPeers[key]; peer != nil {
			peer.conn = conn
			handler := peer.handler
			peer.connect = peer.reconnect
			go p.handleConn(conn, peer, handler)
		}
		p.Unlock()
	}
}

func (p *CentralHelper) handleConn(conn *bleconnecter.BLEConnection, peer *Peer, handler DevHandler) {
	defer func() {
		conn.Close()
		p.PeerRemove(peer)
	}()

	dev := bleatt.NewGattDeviceWithConn(conn, attstructure.NewStructure(), p.config.GATTConfig)

	l2 := blel2cap.New(conn, p.config.L2CAPConfig, func(psm blel2cap.PSMType, accept blel2cap.L2CAPConnAccepter) {
		switch psm {
		case blel2cap.PSMTypeATT:
			dev.AddConn(accept())
		case blel2cap.PSMTypeSecurityManager:
			dev.SetSMP(p.stack.SMP.AddConn(accept(), p.config.SMPConnConfig))
		}
	})

	ctx, cancel := context.WithCancel(p.ctx)
	defer cancel()

	go func() {
		handler(ctx, dev)
		conn.Close()
	}()

	l2.Run()
}

func (p *CentralHelper) Close() error {
	p.cancel()
	return nil
}
