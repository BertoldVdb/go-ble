package hciconnmgr

import (
	"encoding/hex"
	"math"
	"sync"
	"sync/atomic"

	hcievents "github.com/BertoldVdb/go-ble/hci/events"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

type txSlotManager struct {
	sync.Mutex
	connmgr *ConnectionManager

	channel string

	slotBufferLength int

	newSlotsChan   chan (struct{})
	maxSlots       int
	availableSlots int

	newFragmentsChan chan (struct{})
}

func (s *txSlotManager) GetBufferLength() int {
	return s.slotBufferLength
}

func (s *txSlotManager) ReleaseSlots(numSlots int) {
	s.Lock()
	s.availableSlots += numSlots
	bleutil.Assert(s.availableSlots <= s.maxSlots, "Released so many slots there are now more available than the maximum")

	if s.connmgr.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		s.connmgr.logger.WithFields(logrus.Fields{
			"0available": s.availableSlots,
			"1max":       s.maxSlots,
			"2released":  numSlots,
		}).Trace("Slot manager updated (put)")
	}

	select {
	case s.newSlotsChan <- struct{}{}:
	default:
	}
	s.Unlock()
}

func (s *txSlotManager) WaitSlot() bool {
	s.Lock()
	for s.availableSlots == 0 {
		s.Unlock()
		select {
		case <-s.connmgr.closeflag.Chan():
			return false
		case <-s.newSlotsChan:
		}
		s.Lock()
	}

	s.availableSlots--

	if s.connmgr.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		s.connmgr.logger.WithFields(logrus.Fields{
			"0available": s.availableSlots,
			"1max":       s.maxSlots,
			"2taken":     1,
		}).Trace("Slot manager updated (take)")
	}
	s.Unlock()

	return true
}

func createSlotManager(c *ConnectionManager, channel string, slotBufferLength int, maxSlots int) *txSlotManager {
	if maxSlots == 0 {
		return nil
	}

	return &txSlotManager{
		connmgr:          c,
		newSlotsChan:     make(chan (struct{}), 1),
		newFragmentsChan: make(chan (struct{}), 1),

		channel: channel,

		maxSlots:         maxSlots,
		availableSlots:   maxSlots,
		slotBufferLength: slotBufferLength,
	}
}

func (s *txSlotManager) txWorker() error {
	s.connmgr.logger.WithField("0channel", s.channel).Debug("TX Datapump worker started")
	defer func() {
		s.connmgr.logger.WithField("0channel", s.channel).Debug("TX Datapump worker stopped")
	}()

	for {
		select {
		case <-s.connmgr.closeflag.Chan():
			return ErrorClosed
		case <-s.newFragmentsChan:
		}

	sendingLoop:
		for {
			var conn *Connection

			s.connmgr.RLock()
			minOutstanding := int32(math.MaxInt32)
			for _, m := range s.connmgr.connections {
				if m.txSlotManager != s {
					continue
				}

				if m.txFIFO.Len() > 0 {
					outstanding := atomic.LoadInt32(&m.txOutstanding)
					if outstanding < minOutstanding {
						conn = m
						minOutstanding = outstanding
					}
				}
			}
			s.connmgr.RUnlock()

			if conn == nil {
				break sendingLoop
			}

			/* Wait for a slot */
			if !s.WaitSlot() {
				return ErrorClosed
			}

			buf := conn.txFIFO.Pop()
			if buf == nil {
				/* Give back the slot since we cant't use it */
				s.ReleaseSlots(1)

				s.connmgr.logger.WithFields(logrus.Fields{
					"0channel": s.channel,
					"1handle":  conn.handle,
				}).Debug("Requested a slot but was not able to use it")
			} else {
				newOutstanding := atomic.AddInt32(&conn.txOutstanding, 1)

				if s.connmgr.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
					s.connmgr.logger.WithFields(logrus.Fields{
						"0channel":     s.channel,
						"1handle":      conn.handle,
						"2outstanding": newOutstanding,
						"3data":        hex.EncodeToString(buf),
					}).Trace("Sending fragment")
				}

				/* Send the packet */
				err := s.connmgr.sendFunc(buf)
				s.connmgr.rxtxFreeBuffers.Push(buf)

				if err != nil {
					return err
				}
			}
		}
	}
}

func quirkFixBroadcomCompleteEvent(event *hcievents.NumberOfCompletedPacketsEvent) {
	elem := len(event.ConnectionHandle)
	if elem%1 > 0 {
		/* Uneven number of elements */
		for i := 1; i < elem; i += 2 {
			event.ConnectionHandle[i], event.NumCompletedPackets[i] = event.NumCompletedPackets[i], event.ConnectionHandle[i]
		}
	} else {
		/* Even number of elements */
		for i := 1; i < elem; i += 2 {
			event.ConnectionHandle[i], event.NumCompletedPackets[i-1] = event.NumCompletedPackets[i-1], event.ConnectionHandle[i]
		}
	}
}

func (c *ConnectionManager) packetCompleteHandler(event *hcievents.NumberOfCompletedPacketsEvent) *hcievents.NumberOfCompletedPacketsEvent {
	//TODO: add quirk system. This message seems to be encoded incorrectly in broadcom host controllers
	//Fortunately, we can fix it.

	if true {
		quirkFixBroadcomCompleteEvent(event)
		if c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
			c.logger.WithField("0data", event).Trace("Applied Broadcom PktComplete quirk")
		}
	}

	c.RLock()
	for i := range event.ConnectionHandle {
		conn, ok := c.connections[event.ConnectionHandle[i]]
		if ok {
			conn.txSlotManager.ReleaseSlots(int(event.NumCompletedPackets[i]))

			outstanding := atomic.AddInt32(&conn.txOutstanding, -int32(event.NumCompletedPackets[i]))
			bleutil.Assert(outstanding >= 0, "Negative number of outstanding packets")

			if c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
				c.logger.WithFields(logrus.Fields{
					"0handle":      conn.handle,
					"1completed":   event.NumCompletedPackets[i],
					"2outstanding": outstanding,
				}).Trace("Buffer status update received")
			}
		}

	}
	c.RUnlock()
	return event
}

func (c *ConnectionManager) runSlotManagers() error {
	bbuf, err := c.Cmds.InformationalReadBufferSizeSync(nil)
	/* If this errored the controller supports LE only */
	if err == nil {
		c.txSlotManagerEDRACL = createSlotManager(c, "EDR/ACL", int(bbuf.ACLDataPacketLength), int(bbuf.TotalNumACLDataPackets))
		c.txSlotManagerEDRSDO = createSlotManager(c, "EDR/SDO", int(bbuf.SynchronousDataPacketLength), int(bbuf.TotalNumSynchronousDataPackets))
	}

	/* Not all controllers support this command */
	lbuf2, err := c.Cmds.LEReadBufferSizeV2Sync(nil)
	if err == nil {
		if lbuf2.TotalNumLEACLDataPackets == 0 {
			c.txSlotManagerLEACL = c.txSlotManagerEDRACL
		} else {
			c.txSlotManagerLEACL = createSlotManager(c, "LE/ACL", int(lbuf2.LEACLDataPacketLength), int(lbuf2.TotalNumLEACLDataPackets))
		}
		//	c.txSlotManagerLEISO = createSlotManager(c, int(lbuf2.ISODataPacketLength), int(lbuf2.TotalNumISODataPackets))

	} else {
		/* Fallback to old version */
		lbuf, err := c.Cmds.LEReadBufferSizeSync(nil)
		if err == nil {
			if lbuf.TotalNumLEACLDataPackets == 0 {
				c.txSlotManagerLEACL = c.txSlotManagerEDRACL
			} else {
				c.txSlotManagerLEACL = createSlotManager(c, "LE/ACL", int(lbuf.LEACLDataPacketLength), int(lbuf.TotalNumLEACLDataPackets))
			}
		}
	}

	err = c.Events.SetNumberOfCompletedPacketsEventCallback(c.packetCompleteHandler)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var errMutex sync.Mutex

	startWorker := func(a *txSlotManager) {
		if a == nil {
			return
		}

		wg.Add(1)
		go func(a *txSlotManager) {
			defer wg.Done()
			err2 := a.txWorker()

			errMutex.Lock()
			if err == nil {
				err = err2
			}
			errMutex.Unlock()
		}(a)
	}

	startWorker(c.txSlotManagerEDRACL)
	startWorker(c.txSlotManagerEDRSDO)
	if c.txSlotManagerLEACL != c.txSlotManagerEDRACL {
		startWorker(c.txSlotManagerLEACL)
	}

	wg.Wait()

	return err
}
