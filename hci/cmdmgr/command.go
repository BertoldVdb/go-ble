package hcicmdmgr

import (
	"sync"

	"github.com/BertoldVdb/go-misc/closeflag"
	"github.com/sirupsen/logrus"
)

type CommandCompleteCallback func(err error, params []byte) error
type TransmitCallback func(pkt []byte) error

type CommandManager struct {
	logger *logrus.Entry
	sync.Mutex

	closeFlag closeflag.CloseFlag

	commandMaxIssue        int
	commandMaxIssueChanged chan (struct{})

	queues []commandQueue

	transmitFunc TransmitCallback
}

func New(logger *logrus.Entry, maxSlots []int, awaitStartup bool, tx TransmitCallback) *CommandManager {
	s := &CommandManager{
		logger:       logger,
		transmitFunc: tx,
	}

	s.commandMaxIssueChanged = make(chan (struct{}), 1)

	/* Normally the device sends an event when it is done booting, but depending on system configuration
	   this may never be received. If awaitStartup is false we assume the device is ready at start.
	   Some devices also don't send it at all, instead requiring a fixed delay after power up. */
	if !awaitStartup {
		s.commandMaxIssue = 1
	}

	s.queues = make([]commandQueue, len(maxSlots))
	for i := range s.queues {
		s.queues[i].Init(s, maxSlots[i], i)
	}

	return s
}

func (s *CommandManager) Run() error {
	var wg sync.WaitGroup
	var err error

	wg.Add(len(s.queues))

	s.logger.Debug("Starting command queue workers")

	for i := range s.queues {
		go func(worker int) {
			err2 := s.queues[worker].Worker()
			if err2 != nil {
				err = err2
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	return err
}

func (s *CommandManager) Close() error {
	if s.closeFlag.Close() == nil {
		for i := range s.queues {
			s.queues[i].closeQueue()
		}
	}
	return nil
}
