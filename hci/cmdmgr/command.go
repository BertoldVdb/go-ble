package hcicmdmgr

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type CommandCompleteCallback func(err error, params []byte) error
type TransmitCallback func(pkt []byte) error

type CommandManager struct {
	logger *logrus.Entry
	sync.Mutex

	closed bool

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

func (s *CommandManager) Close() {
	s.Lock()

	if s.closed {
		s.Unlock()
		return
	}
	s.closed = true
	s.Unlock()

	close(s.commandMaxIssueChanged)

	for i := range s.queues {
		s.queues[i].closeQueue()
	}
}
