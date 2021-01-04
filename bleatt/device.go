package bleatt

import (
	"context"
	"encoding/binary"
	"errors"
	"strings"
	"sync"
	"sync/atomic"

	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"
)

type GattDevice struct {
	config *GattDeviceConfig

	connsMutex  sync.Mutex
	conns       map[hciconnmgr.BufferConn](*gattDeviceConn)
	initialConn *gattDeviceConn

	ourMTU uint16

	server attServer

	//You only need to use either the done channel or the mutex
	clientStructureMutex        sync.Mutex
	clientStructureDiscoverOnce sync.Once
	clientDiscoveryDone         chan (struct{})
	clientStructure             *attstructure.Structure
}

type gattDeviceConn struct {
	parent *GattDevice

	conn   hciconnmgr.BufferConn
	logger *logrus.Entry

	mtuRequest sync.Once
	mtu        uint32

	client attClient
}

type GattDeviceConfig struct {
	ConnCb                  func(numConnections int)
	DeviceName              string
	Appearance              uint16
	DiscoverRemoteOnConnect bool
	MTU                     uint16
}

func DefaultConfig() *GattDeviceConfig {
	return &GattDeviceConfig{
		MTU:                     0xFFFF,
		DiscoverRemoteOnConnect: true,
		DeviceName:              "go-ble device",
	}
}

func NewGattDevice(externalStructure *attstructure.Structure, config *GattDeviceConfig) *GattDevice {
	if config == nil {
		config = DefaultConfig()
	}

	dev := &GattDevice{
		conns:               make(map[hciconnmgr.BufferConn](*gattDeviceConn)),
		ourMTU:              0xFFFF,
		config:              config,
		clientDiscoveryDone: make(chan (struct{})),
	}

	gattStructure := attstructure.NewStructure()
	pble := gattStructure.AddPrimaryService(bleutil.UUIDFromStringPanic("1800"))
	pble.AddCharacteristicReadOnly(bleutil.UUIDFromStringPanic("2a00"), []byte(config.DeviceName)) /* Device name */
	var apBuf [2]byte
	binary.LittleEndian.PutUint16(apBuf[:], config.Appearance)
	pble.AddCharacteristicReadOnly(bleutil.UUIDFromStringPanic("2a01"), apBuf[:]) /* Appearance: Generic network device */
	gattStructure.AddPrimaryService(bleutil.UUIDFromStringPanic("1801"))          /* Empty */

	exportedStructure := &attstructure.ExportedStructure{}
	exportedStructure.Append(gattStructure)
	exportedStructure.Append(externalStructure)

	/* Init GATT server */
	dev.server.init(dev, exportedStructure)

	return dev
}

func (d *gattDeviceConn) handlePDU(buf *pdu.PDU) (bool, error) {
	if buf == nil {
		d.client.close()
		return false, nil
	}

	valid, method, isAuthenticated, isForServer := getOpcode(buf)
	if !valid {
		d.logger.WithFields(logrus.Fields{
			"0method":      method,
			"1isAuth":      isAuthenticated,
			"2isForServer": isForServer,
			"3buf":         buf,
		}).Info("ATT Invalid received")
		return false, nil
	}

	if d.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		d.logger.WithFields(logrus.Fields{
			"0method":      method,
			"1isAuth":      isAuthenticated,
			"2isForServer": isForServer,
			"3buf":         buf,
		}).Trace("ATT Received")
	}

	buf.DropLeft(1)

	if isForServer {
		return d.parent.server.handlePDU(d, method, isAuthenticated, buf)
	}

	return d.client.handlePDU(method, isAuthenticated, buf)
}

func (d *gattDeviceConn) handleConn() error {
	defer d.parent.CloseConn(d.conn)

	ctx := context.Background()

	for {
		buf, err := d.conn.ReadBuffer(ctx)
		if err != nil {
			d.handlePDU(nil)
			return err
		}

		keepBuffer, err := d.handlePDU(buf)
		if err != nil {
			return err
		}

		if !keepBuffer {
			bleutil.ReleaseBuffer(buf)
		}
	}
}

func (g *GattDevice) CloseConn(conn hciconnmgr.BufferConn) error {
	g.connsMutex.Lock()
	defer g.connsMutex.Unlock()

	_, ok := g.conns[conn]
	if !ok {
		return errors.New("Connection not found")
	}

	delete(g.conns, conn)
	if g.config.ConnCb != nil {
		g.config.ConnCb(len(g.conns))
	}

	conn.UseDone()
	return conn.Close()
}

func (g *GattDevice) AddConn(conn hciconnmgr.BufferConn) error {
	d := &gattDeviceConn{
		parent: g,
		conn:   conn,
		logger: bleutil.LogWithPrefix(conn.GetLogger(), "att"),

		mtu: 23,
	}

	conn.UseStart()

	d.client.init(d)

	g.connsMutex.Lock()
	initial := false
	if g.initialConn == nil { //This is the connection that is always present
		g.initialConn = d
		initial = true
	}
	g.conns[conn] = d
	if g.config.ConnCb != nil {
		g.config.ConnCb(len(g.conns))
	}
	g.connsMutex.Unlock()

	go d.handleConn()
	go d.getMTUBlocking()
	if initial && d.parent.config.DiscoverRemoteOnConnect {
		go g.clientDiscover(d)
	}

	return nil
}

func (d *gattDeviceConn) getMTUBlocking() int {
	d.mtuRequest.Do(func() {
		buf := bleutil.GetBuffer(3)
		buf.Buf()[0] = byte(ATTExchangeMTUReq)
		binary.LittleEndian.PutUint16(buf.Buf()[1:], d.parent.ourMTU)
		cmd, response, err := d.client.sendCommand(context.Background(), buf, true)
		defer bleutil.ReleaseBuffer(response)

		if err == nil && cmd == ATTExchangeMTURsp && response.Len() == 2 {
			newMTU := binary.LittleEndian.Uint16(response.Buf())
			d.setMTU(newMTU)
		}

	})

	return d.getMTU()
}

func (d *gattDeviceConn) getMTU() int {
	mtu := atomic.LoadUint32(&d.mtu)

	return int(mtu)
}

func (d *gattDeviceConn) setMTU(new uint16) int {
	new32 := uint32(new)
	mtu := atomic.LoadUint32(&d.mtu)

	if new32 > mtu {
		d.logger.WithField("0old", mtu).WithField("1new", new).Debug("Update MTU")
		mtu = new32
	}

	atomic.StoreUint32(&d.mtu, new32)
	return int(new32)
}

func (d *GattDevice) getConnWithHighestMTU() *gattDeviceConn {
	d.connsMutex.Lock()
	var conn *gattDeviceConn
	var maxMTU int
	for _, m := range d.conns {
		cmtu := m.getMTU()
		if cmtu > maxMTU {
			maxMTU = cmtu
			conn = m
		}
	}
	d.connsMutex.Unlock()
	return conn
}

func (d *GattDevice) clientDiscover(conn *gattDeviceConn) {
	s, err := conn.client.discoverRemoteDeviceStructure()

	if err != nil {
		conn.logger.WithError(err).Warn("Failed to parse remote GATT definition")
	} else {
		conn.logger.Info("Discovery completed")
		for _, m := range strings.Split(s.String(), "\n") {
			conn.logger.Debugf("GATT: %s", strings.Trim(m, "\r\n"))
		}
	}

	d.clientStructureMutex.Lock()
	defer d.clientStructureMutex.Unlock()
	d.clientStructure = s
	close(d.clientDiscoveryDone)
}

func (d *GattDevice) ClientGetStructure(ctx context.Context) *attstructure.Structure {
	select {
	case <-d.clientDiscoveryDone:
		return d.clientStructure
	case <-ctx.Done():
		return nil
	}
}

func (d *GattDevice) ClientRead(ctx context.Context, handle uint16, buf []byte) ([]byte, error) {
	//TODO: Allow using another connection for client requests
	conn := d.initialConn

	result, atterr, err := conn.client.readHandleAll(ctx, handle, buf)
	if err != nil {
		return nil, err
	}
	err = attErrorToError(atterr)
	return result, err
}

func (d *GattDevice) ClientWrite(ctx context.Context, handle uint16, buf []byte, withRsp bool) (int, error) {
	conn := d.initialConn

	l, atterr, err := conn.client.writeHandle(ctx, handle, buf, withRsp)
	if err != nil {
		return l, err
	}
	err = attErrorToError(atterr)
	return l, err
}

/* Note that this is just an indication, the real MTU may be different */
func (d *GattDevice) ServerGetNotifyMTU(characteristic *attstructure.Characteristic) int {
	conn := d.getConnWithHighestMTU()

	if conn == nil {
		return 23
	}

	return conn.getMTUBlocking()
}

func (d *GattDevice) HasConnections() bool {
	d.connsMutex.Lock()
	defer d.connsMutex.Unlock()

	return len(d.conns) > 0
}
