package hciconnmgr

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/BertoldVdb/go-misc/bufferfifo"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"
)

type AppConn interface {
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

type BufferConn interface {
	IsOpen() bool
	Close() error

	ReadBuffer(ctx context.Context) (*pdu.PDU, error)
	WriteBuffer(buf *pdu.PDU) error

	GetLogger() *logrus.Entry
	UseStart()
	UseDone()
}

type Connection struct {
	net.PacketConn

	connmgr *ConnectionManager
	handle  uint16

	txFIFO        *bufferfifo.FIFO
	txSlotManager *txSlotManager
	txOutstanding int32

	rxPDU          *pdu.PDU
	rxFIFO         *bufferfifo.FIFO
	rxNewDataChan  chan (struct{})
	rxContextMutex sync.Mutex
	rxContext      context.Context
	rxCancelFunc   context.CancelFunc

	closeChan chan (struct{})
	closeOnce sync.Once
	closeFunc func() error

	disconnectedOnce sync.Once

	AppConn AppConn
}

func (c *Connection) UseStart() {
	c.connmgr.useWg.Add(1)
}

func (c *Connection) UseDone() {
	c.connmgr.useWg.Done()
}

func (c *Connection) GetLogger() *logrus.Entry {
	return c.connmgr.logger.WithField("zhandle", c.handle)
}

func (c *Connection) GetHandle() uint16 {
	return c.handle
}

func (c *Connection) GetConnectionManager() *ConnectionManager {
	return c.connmgr
}

func (c *Connection) IsOpen() bool {
	select {
	case <-c.closeChan:
		return false
	default:
	}

	return true
}

// Close closes the connection.
// Any blocked ReadFrom or WriteTo operations will be unblocked and return errors.
func (c *Connection) Close() error {
	var err error

	c.closeOnce.Do(func() {
		select {
		case <-c.closeChan:
			return
		default:
		}

		err = c.connmgr.Cmds.LinkControlDisconnectSync(hcicommands.LinkControlDisconnectInput{
			ConnectionHandle: c.handle,
			Reason:           0x13,
		})

		c.connmgr.logger.WithError(err).WithField("0handle", c.handle).Trace("Requested disconnect")

		if err == nil {
			<-c.closeChan
		}

		c.connmgr.logger.WithError(err).WithField("0handle", c.handle).Debug("Completed disconnect")
		if c.connmgr.config.HookConnectionStateChange != nil {
			c.connmgr.config.HookConnectionStateChange(c, false)
		}
	})

	return err
}

func (c *Connection) disconnected() {
	c.disconnectedOnce.Do(func() {
		/* Throw away all buffers and return the slots */
		for {
			buf := c.txFIFO.Pop()
			if buf == nil {
				break
			}
			bleutil.ReleaseBuffer(buf)
		}
		c.txSlotManager.ReleaseSlots(int(atomic.SwapInt32(&c.txOutstanding, 0)))

		c.connmgr.logger.WithField("0handle", c.handle).Info("Connection lost")

		bleutil.ReleaseBuffer(c.rxPDU)
		c.rxPDU = nil

		close(c.closeChan)

		go func() {
			if c.closeFunc != nil {
				c.connmgr.logger.WithError(c.closeFunc()).WithField("0handle", c.handle).Debug("Close function completed")
			}
		}()
	})
}

func (c *ConnectionManager) ConnectionNew(handle uint16, closeFunc func() error) *Connection {
	c.Lock()
	_, ok := c.connections[handle]
	bleutil.Assert(!ok, "Connection handle already exists")

	conn := &Connection{
		connmgr:       c,
		handle:        handle,
		txFIFO:        bufferfifo.New(16),
		txSlotManager: c.txSlotManagerLEACL, //TODO: Make this dynamic based on the connection type...
		rxFIFO:        bufferfifo.New(16),
		rxNewDataChan: make(chan (struct{}), 1),
		closeChan:     make(chan (struct{})),
		rxContext:     context.Background(),
		rxPDU:         bleutil.GetBuffer(64),
		closeFunc:     closeFunc,
	}
	c.connections[handle] = conn
	c.Unlock()

	c.logger.WithField("0handle", handle).Debug("Created new connection")

	if c.config.HookConnectionStateChange != nil {
		c.config.HookConnectionStateChange(conn, true)
	}

	return conn
}

func (c *Connection) ReadBuffer(ctx context.Context) (*pdu.PDU, error) {
	if !c.IsOpen() {
		return nil, ErrorConnectionClosed
	}

	for {
		buf := c.rxFIFO.Pop()
		if buf != nil {
			select {
			case c.rxNewDataChan <- struct{}{}:
			default:
			}

			return buf, nil
		}

		select {
		case <-c.rxNewDataChan:
		case <-c.connmgr.closeflag.Chan():
			return nil, ErrorClosed
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-c.closeChan:
			return nil, ErrorConnectionClosed
		}
	}
}

func (c *Connection) WriteBuffer(buf *pdu.PDU) error {
	if !c.IsOpen() {
		bleutil.ReleaseBuffer(buf)
		return ErrorConnectionClosed
	}

	//TODO: Check what type of encoder to use if we support more than ACL
	c.encodeACL(buf)
	return nil
}

// WriteTo writes a packet with payload p to addr.
// WriteTo can be made to time out and return
// an Error with Timeout() == true after a fixed time limit;
// see SetDeadline and SetWriteDeadline.
// On packet-oriented connections, write timeouts are rare.
func (c *Connection) WriteTo(l2cap []byte, addr net.Addr) (int, error) {
	err := c.WriteBuffer(bleutil.CopyBufferFromSlice(l2cap))
	if err != nil {
		return 0, err
	}
	return len(l2cap), err
}

// ReadFrom reads a packet from the connection,
// copying the payload into p. It returns the number of
// bytes copied into p and the return address that
// was on the packet.
// It returns the number of bytes read (0 <= n <= len(p))
// and any error encountered. Callers should always process
// the n > 0 bytes returned before considering the error err.
// ReadFrom can be made to time out and return
// an Error with Timeout() == true after a fixed time limit;
// see SetDeadline and SetReadDeadline.
func (c *Connection) ReadFrom(buf []byte) (int, net.Addr, error) {
	if len(buf) == 0 {
		return 0, nil, nil
	}

	c.rxContextMutex.Lock()
	rx, err := c.ReadBuffer(c.rxContext)
	c.rxContextMutex.Unlock()

	if err != nil {
		return 0, nil, err
	}

	copyLen := rx.Len()
	if copyLen > len(buf) {
		copyLen = len(buf)
	}

	copy(buf[:copyLen], rx.Buf()[:copyLen])

	bleutil.ReleaseBuffer(rx)

	addr := net.Addr(nil)
	if c.AppConn != nil {
		addr = c.AppConn.RemoteAddr()
	}

	return copyLen, addr, nil
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to ReadFrom or
// WriteTo. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful ReadFrom or WriteTo calls.
//
// A zero value for t means I/O operations will not time out.
func (c *Connection) SetDeadline(t time.Time) error {
	c.SetReadDeadline(t)
	c.SetWriteDeadline(t)
	return nil
}

// SetReadDeadline sets the deadline for future ReadFrom calls
// and any currently-blocked ReadFrom call.
// A zero value for t means ReadFrom will not time out.
func (c *Connection) SetReadDeadline(t time.Time) error {
	c.rxContextMutex.Lock()
	defer c.rxContextMutex.Unlock()

	if c.rxCancelFunc != nil {
		c.rxCancelFunc()
		c.rxCancelFunc = nil
	}

	if t.IsZero() {
		c.rxContext = context.Background()
		return nil
	}

	c.rxContext, c.rxCancelFunc = context.WithDeadline(context.Background(), t)
	return nil
}

// SetWriteDeadline sets the deadline for future WriteTo calls
// and any currently-blocked WriteTo call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means WriteTo will not time out.
func (c *Connection) SetWriteDeadline(t time.Time) error {
	return nil
}

// LocalAddr returns the local network address.
func (c *Connection) LocalAddr() net.Addr {
	if c.AppConn == nil {
		return nil
	}
	return c.AppConn.LocalAddr()
}
