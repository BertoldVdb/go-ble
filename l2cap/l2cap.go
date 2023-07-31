package blel2cap

import (
	"context"
	"encoding/binary"
	"errors"
	"sync"
	"time"

	"github.com/BertoldVdb/go-ble/bleconnecter"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"
)

type RxHandler func(cid uint16, buf *pdu.PDU) (error, bool)

var (
	ErrorRxFailed = errors.New("Failed to read data from connection")
)

type l2cid struct {
	rxHandler RxHandler
}

type L2CAP struct {
	sync.RWMutex

	conn hciconnmgr.BufferConn

	cidMap   map[uint16]*l2cid
	paramBuf []uint16
	logger   *logrus.Entry

	isLE      bool
	newConnCb L2CAPConnCB

	sig    *signalling
	config *L2CAPConfig
}

type L2CAPConfig struct {
	BLEUpdateParametersVerify func(c *bleconnecter.BLEConnection, intervalMin uint16, intervalMax uint16, latency uint16, timeout uint16) bool
}

type L2CAPConnAccepter func() hciconnmgr.BufferConn
type L2CAPConnCB func(psm PSMType, accepter L2CAPConnAccepter)

func New(conn hciconnmgr.BufferConn, config *L2CAPConfig, newConnCb L2CAPConnCB) *L2CAP {
	l := &L2CAP{
		conn: conn,

		cidMap: make(map[uint16]*l2cid),
		logger: bleutil.LogWithPrefix(conn.GetLogger(), "l2"),
		config: config,

		newConnCb: newConnCb,
	}

	_, l.isLE = l.conn.(*bleconnecter.BLEConnection)
	l.sig = l.signallingInit()

	return l
}

func (l *L2CAP) processInput(buf *pdu.PDU) (error, bool) {
	header := buf.DropLeft(4)
	if header == nil {
		return nil, false
	}

	frameLen := int(binary.LittleEndian.Uint16(header[0:2]))
	cid := binary.LittleEndian.Uint16(header[2:4])

	if frameLen != buf.Len() {
		return nil, false
	}

	l.RLock()
	c, ok := l.cidMap[cid]
	l.RUnlock()
	if !ok {
		return nil, false
	}

	if l.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		l.logger.WithFields(logrus.Fields{
			"0cid":     cid,
			"1payload": buf,
		}).Trace("Read L2CAP data")
	}

	if c.rxHandler != nil {
		return c.rxHandler(cid, buf)
	}

	return nil, false
}

func (l *L2CAP) unregisterCID(cid uint16) {
	l.Lock()
	defer l.Unlock()

	delete(l.cidMap, cid)
}

func (l *L2CAP) registerCallbackCID(cid uint16, rxCb RxHandler) bool {
	l.Lock()
	defer l.Unlock()

	_, ok := l.cidMap[cid]
	if ok {
		return false
	}

	l.cidMap[cid] = &l2cid{
		rxHandler: rxCb,
	}

	return true
}

func (l *L2CAP) WriteBufferCID(cid uint16, buf *pdu.PDU) error {
	header := buf.ExtendLeft(4)
	binary.LittleEndian.PutUint16(header[0:2], uint16(buf.Len()-4))
	binary.LittleEndian.PutUint16(header[2:4], cid)

	if l.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		l.logger.WithFields(logrus.Fields{
			"0cid":     cid,
			"1payload": buf,
		}).Trace("Writing L2CAP data")
	}

	return l.conn.WriteBuffer(buf)
}

func (l *L2CAP) readTask(c chan (*pdu.PDU), ctx context.Context) {
	defer close(c)

	for {
		buf, err := l.conn.ReadBuffer(ctx)
		if err != nil {
			return
		}

		c <- buf
	}
}

func (l *L2CAP) Run() error {
	defer l.conn.UseDone()
	l.conn.UseStart()
	defer l.Close()
	defer l.sig.failActiveToken()

	l.connectionDefaultInit()

	ctx := context.Background()

	rxChan := make(chan (*pdu.PDU))
	go l.readTask(rxChan, ctx)

	tokenChan := l.sig.cmdQueue.GetCommittedTokenChan(ctx)

	for {
		select {
		case buf, ok := <-rxChan:
			if !ok {
				return ErrorRxFailed
			}

			err, keepBuffer := l.processInput(buf)
			if err != nil {
				return err
			}

			if !keepBuffer {
				bleutil.ReleaseBuffer(buf)
			}

		case t, ok := <-tokenChan:
			if !ok {
				return ErrorRxFailed
			}

			l.sig.handleTokenTx(t)

		case <-l.sig.tokenTimeoutChan():
			l.sig.failActiveToken()
		}
	}
}

func (l *L2CAP) Close() error {
	var conn []*l2cid
	l.Lock()
	for _, m := range l.cidMap {
		conn = append(conn, m)
	}
	l.Unlock()

	for _, m := range conn {
		m.rxHandler(0, nil)
	}

	l.sig.Close()
	return l.conn.Close()
}

func (l *L2CAP) SendCommandUint16(ctx context.Context, code uint8, result []uint16, params ...uint16) (uint8, []uint16, error) {
	return l.sig.SendCommandUint16(ctx, code, result, params...)
}

func (l *L2CAP) Ping() (time.Duration, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	/* Few LE devices support the echo command, but they will then reply with the error command, which is also fine.
	   EDR devices mostly support the echo command. */
	before := time.Now()
	_, _, err := l.SendCommandUint16(ctx, SigEchoReq, nil)
	duration := time.Now().Sub(before)
	cancel()

	return duration, err
}
