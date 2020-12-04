package blel2cap

import (
	"context"
	"errors"

	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/BertoldVdb/go-misc/bufferfifo"
	"github.com/BertoldVdb/go-misc/closeflag"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"
)

var (
	ErrorL2ChannelClosed = errors.New("L2CAP channel is closed")
)

type L2Connection struct {
	l         *L2CAP
	canClose  bool
	localCID  uint16
	remoteCID uint16
	rxBuffer  *bufferfifo.FIFO
	rxChan    chan (struct{})

	closed closeflag.CloseFlag
}

/*type BufferConn interface {
	IsOpen() bool
	Close() error

	ReadBuffer(ctx context.Context) (*pdu.PDU, error)
	WriteBuffer(buf *pdu.PDU) error

	GetLogger() *logrus.Entry
	UseStart()
	UseDone()
}
*/

func (c *L2Connection) IsOpen() bool {
	return !c.closed.IsClosed()
}

func (c *L2Connection) Close() error {
	err := c.closed.Close()
	if err == nil {
		c.l.unregisterCID(c.localCID)

		if c.canClose {
			var result [2]uint16
			_, _, err = c.l.SendCommandUint16(context.Background(), SigDisconnectionReq, result[:], c.remoteCID, c.localCID)
		} else {
			/* If one the internal channels needs to be closed, we need to close the while l2 */
			return c.l.Close()
		}
	}

	return err
}

func (c *L2Connection) ReadBuffer(ctx context.Context) (*pdu.PDU, error) {
	if !c.IsOpen() || !c.l.conn.IsOpen() {
		return nil, ErrorL2ChannelClosed
	}

	for {
		buf := c.rxBuffer.Pop()
		if buf != nil {
			select {
			case c.rxChan <- struct{}{}:
			default:
			}
			return buf, nil
		}

		select {
		case <-c.rxChan:
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-c.closed.Chan():
			return nil, ErrorL2ChannelClosed
		}
	}
}

func (c *L2Connection) WriteBuffer(buf *pdu.PDU) error {
	if !c.IsOpen() {
		bleutil.ReleaseBuffer(buf)
		return ErrorL2ChannelClosed
	}

	return c.l.WriteBufferCID(c.remoteCID, buf)
}

func (c *L2Connection) GetLogger() *logrus.Entry {
	return c.l.logger
}

func (c *L2Connection) UseStart() {
	c.l.conn.UseStart()
}

func (c *L2Connection) UseDone() {
	c.l.conn.UseDone()
}

func (c *L2Connection) connectionRxHandler(cid uint16, buf *pdu.PDU) (error, bool) {
	if buf == nil {
		return c.Close(), false
	}
	c.rxBuffer.Push(buf)

	select {
	case c.rxChan <- struct{}{}:
	default:
	}

	return nil, true
}

func (l *L2CAP) connectionCreateInternal(canClose bool, localCID uint16, remoteCID uint16) *L2Connection {
	c := &L2Connection{
		l:         l,
		canClose:  canClose,
		localCID:  localCID,
		remoteCID: remoteCID,
		rxBuffer:  bufferfifo.New(16),
		rxChan:    make(chan (struct{}), 1),
	}

	if !l.registerCallbackCID(localCID, c.connectionRxHandler) {
		return nil
	}

	return c
}

//TODO: Implement when we need it
func (l *L2CAP) Dial() *L2Connection {
	return nil
}

func (l *L2CAP) connectionDefaultInit() {
	if l.isLE {
		l.newConnCb(PSMTypeATT, func() hciconnmgr.BufferConn {
			return l.connectionCreateInternal(false, 4, 4) /* ATT */
		})
		l.newConnCb(PSMTypeSecurityManager, func() hciconnmgr.BufferConn {
			return l.connectionCreateInternal(false, 6, 6) /* SMP */
		})
	}
}
