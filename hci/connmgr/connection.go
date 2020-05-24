package hciconnmgr

import (
	"context"
	"sync/atomic"

	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/BertoldVdb/go-misc/bufferfifo"
)

type Connection struct {
	connmgr *ConnectionManager
	handle  uint16

	state uint32

	txFIFO        *bufferfifo.FIFO
	txSlotManager *txSlotManager
	txOutstanding int32

	rxBuffer      []byte
	rxFIFO        *bufferfifo.FIFO
	rxNewDataChan chan (struct{})

	closeChan chan (struct{})

	AppConn interface{}
}

func (c *Connection) IsOpen() bool {
	return atomic.LoadUint32(&c.state) == 0
}

func (c *Connection) Close() error {
	old := atomic.SwapUint32(&c.state, 1)
	if old != 0 {
		<-c.closeChan
		return nil
	}

	err := c.connmgr.Cmds.LinkControlDisconnectSync(hcicommands.LinkControlDisconnectInput{
		ConnectionHandle: c.handle,
		Reason:           0x13,
	})

	c.connmgr.logger.WithError(err).WithField("0handle", c.handle).Trace("Requested disconnect")

	if err == nil {
		<-c.closeChan
	}

	c.connmgr.logger.WithError(err).WithField("0handle", c.handle).Debug("Completed disconnect")

	return err
}

func (c *Connection) disconnected() {
	old := atomic.SwapUint32(&c.state, 2)

	if old != 2 {
		/* Throw away all buffers and return the slots */
		for {
			buf := c.txFIFO.Pop()
			if buf == nil {
				break
			}
			c.connmgr.rxtxFreeBuffers.Push(buf)
		}
		c.txSlotManager.ReleaseSlots(int(atomic.SwapInt32(&c.txOutstanding, 0)))

		c.connmgr.logger.WithField("0handle", c.handle).Info("Connection lost")

		close(c.closeChan)
	}
}

func (c *ConnectionManager) ConnectionNew(handle uint16) *Connection {
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
	}
	c.connections[handle] = conn
	c.Unlock()

	c.logger.WithField("0handle", handle).Debug("Created new connection")

	return conn
}

func (c *Connection) Write(l2cap []byte) (int, error) {
	if !c.IsOpen() {
		return 0, ErrorConnectionClosed
	}

	//TODO: Check what type of encoder to use if we support more than ACL
	c.encodeACL(l2cap)
	return len(l2cap), nil
}

func (c *Connection) ReadBuffer(ctx context.Context) ([]byte, error) {
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

func (c *Connection) ReadBufferReturn(buf []byte) {
	c.connmgr.rxtxFreeBuffers.Push(buf)
}

func (c *Connection) Read(buf []byte) (int, error) {
	if len(buf) == 0 {
		return 0, nil
	}

	rx, err := c.ReadBuffer(context.Background())
	if err != nil {
		return 0, err
	}

	copyLen := len(rx)
	if copyLen > len(buf) {
		copyLen = len(buf)
	}

	copy(buf[:copyLen], rx[:copyLen])

	c.ReadBufferReturn(rx)

	return copyLen, nil
}
