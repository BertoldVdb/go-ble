package attperipheral

import (
	"context"
	"net"

	"github.com/BertoldVdb/go-ble"
	"github.com/BertoldVdb/go-ble/bleatt"
	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	"github.com/BertoldVdb/go-ble/bleconnecter"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	blel2cap "github.com/BertoldVdb/go-ble/l2cap"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/BertoldVdb/go-misc/closeflag"
)

type PeripheralImplementation interface {
	CreateStructure(*attstructure.Structure) error
	Connected(conn hciconnmgr.BufferConn) error
	Disconnected()
}

type PeripheralImplementationFactory func() PeripheralImplementation

type PeripheralHelper struct {
	stack *ble.BluetoothStack
	impl  []PeripheralImplementationFactory

	ctx    context.Context
	cancel context.CancelFunc

	config PeripheralHelperConfig
}

func (p *PeripheralHelper) handleConn(conn hciconnmgr.BufferConn, remoteAddr net.Addr) error {
	defer conn.Close()

	impl := make([]PeripheralImplementation, len(p.impl))
	for i, m := range p.impl {
		impl[i] = m()
	}

	structure := attstructure.NewStructure()

	for _, m := range impl {
		err := m.CreateStructure(structure)
		if err != nil {
			return err
		}
	}

	var cf closeflag.CloseFlag

	gattConfig := bleatt.DefaultConfig()

	gattConfig.ConnCb = func(num int) {
		if num == 0 {
			cf.Close()
		}
	}

	gattConfig.DeviceName = p.config.DeviceName
	gattConfig.Appearance = p.config.Appearance
	gattConfig.DiscoverRemoteOnConnect = false

	dev := bleatt.NewGattDevice(structure, gattConfig)

	var err error

	for _, m := range impl {
		err = m.Connected(conn)
		if err != nil {
			cf.Close()
		}
	}

	if err == nil {
		l2 := blel2cap.New(conn, nil, func(psm blel2cap.PSMType, accept blel2cap.L2CAPConnAccepter) {
			switch psm {
			case blel2cap.PSMTypeATT:
				dev.AddConn(accept())
			}
		})
		go l2.Run()
	}

	<-cf.Chan()

	for _, m := range impl {
		m.Disconnected()
	}

	return err
}

type PeripheralHelperConfig struct {
	MACFilter        []bleutil.BLEAddr
	ConnectionParams bleconnecter.BLEConnectionParametersRequested

	DeviceName string
	Appearance uint16
}

func DefaultConfig() PeripheralHelperConfig {
	return PeripheralHelperConfig{
		DeviceName: "Unset",
	}
}

func New(stack *ble.BluetoothStack, config PeripheralHelperConfig) *PeripheralHelper {
	p := &PeripheralHelper{
		stack:  stack,
		config: config,
	}

	p.ctx, p.cancel = context.WithCancel(context.Background())

	return p
}

func (p *PeripheralHelper) RegisterImplementation(impl PeripheralImplementationFactory) {
	p.impl = append(p.impl, impl)
}

func (p *PeripheralHelper) Run() error {
	for {
		conn, _, err := p.stack.BLEConnecter.Connect(p.ctx, false, p.config.MACFilter, p.config.ConnectionParams)
		if err != nil {
			return err
		}

		go p.handleConn(conn, conn.RemoteAddr())
	}
}

func (p *PeripheralHelper) Close() error {
	p.cancel()
	return nil
}
