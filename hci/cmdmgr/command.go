package hcicmdmgr

import (
	"log"
	"sync"
)

type CommandCompleteCallback func(err error, params []byte) error
type TransmitCallback func(pkt []byte) error

type CommandManager struct {
	sync.Mutex

	closed bool

	commandMaxIssue        int
	commandMaxIssueChanged chan (struct{})

	queues []commandQueue

	transmitFunc TransmitCallback
}

func New(maxSlots []int, awaitStartup bool, tx TransmitCallback) *CommandManager {
	s := &CommandManager{

		transmitFunc: tx,
	}

	s.commandMaxIssueChanged = make(chan (struct{}), 1)

	if !awaitStartup {
		s.commandMaxIssue = 1
	}

	s.queues = make([]commandQueue, len(maxSlots))
	for i := range s.queues {
		s.queues[i].Init(s, maxSlots[i])
	}

	return s
}

func (s *CommandManager) Run() error {
	var wg sync.WaitGroup
	var err error

	wg.Add(len(s.queues))

	for i := range s.queues {
		go func(worker int) {
			err2 := s.queues[worker].Worker()
			log.Println("Worker quit", err2)
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
	log.Println("Closing parent")
	s.Lock()
	log.Println("Closing parent have lock")

	if s.closed {
		s.Unlock()
		return
	}
	s.closed = true
	s.Unlock()

	log.Println("Closing parent released lock")

	close(s.commandMaxIssueChanged)
	for i := range s.queues {
		s.queues[i].closeQueue()
	}

	log.Println("Closing parent done")
}
