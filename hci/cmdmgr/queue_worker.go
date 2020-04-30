package hcicmdmgr

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/BertoldVdb/go-misc/tokenqueue"
)

type commandQueue struct {
	parent *CommandManager

	commandQueue  *tokenqueue.Queue
	commandActive []*commandToken

	tokenRemoved chan (struct{})

	closed bool
}

type commandToken struct {
	sync.Mutex

	opcode uint16
	data   []byte

	sync      bool
	cb        CommandCompleteCallback
	completed chan (error)

	timeoutTime time.Time
}

var (
	ErrorWorkerClosed = errors.New("The worker was stopped")
	ErrorTimeout      = errors.New("A command timed out.")
)

func (c *commandToken) Cleanup() {
	close(c.completed)
}

func (s *commandQueue) Init(parent *CommandManager, numSlots int) {
	s.parent = parent
	s.commandQueue = tokenqueue.NewQueue(numSlots, numSlots, func() tokenqueue.Token {
		return &commandToken{
			completed: make(chan (error), 1),
		}
	})
	s.commandActive = make([]*commandToken, numSlots)
	s.tokenRemoved = make(chan (struct{}), 1)
}

func (s *commandQueue) findToken(opcode uint16, clear bool) *commandToken {
	if s.parent.closed {
		return nil
	}

	for i, m := range s.commandActive {
		if m != nil && m.opcode == opcode {
			token := s.commandActive[i]
			if clear {
				s.commandActive[i] = nil
				select {
				case s.tokenRemoved <- struct{}{}:
				default:
				}
			}
			return token
		}
	}

	return nil
}

func (s *commandQueue) closeQueue() {
	/* Return active tokens */
	s.parent.Lock()
	if s.closed {
		s.parent.Unlock()
		return
	}
	s.closed = true

	var toComplete []*commandToken

	for i, m := range s.commandActive {
		if m != nil {
			toComplete = append(toComplete, m)
			s.commandActive[i] = nil
		}
	}

	close(s.tokenRemoved)
	s.parent.Unlock()

	for _, m := range toComplete {
		s.tokenComplete(m, ErrorWorkerClosed, nil)
	}

	s.commandQueue.Close()
}

func (s *commandQueue) waitCondition(condition chan (struct{}), timeout <-chan (time.Time)) error {
	s.parent.Unlock()

	stuckCnt := 0
	for {
		select {
		case _, ok := <-condition:
			if !ok {
				return ErrorWorkerClosed
			}

			s.parent.Lock()
			return nil
		case <-timeout:
			stuckCnt++
			if stuckCnt >= 2 {
				return ErrorTimeout
			}
		}
	}
}

func (s *commandQueue) Worker() error {
	defer s.parent.Close()

	stuckTimer := time.NewTicker(time.Second / 2)
	defer stuckTimer.Stop()

	var workingToken *commandToken
	defer func() {
		if workingToken != nil {
			/* This token was not yet activated, so it can be signalled complete without risk */
			s.tokenComplete(workingToken, ErrorWorkerClosed, nil)
		}
	}()

	commitChan := s.commandQueue.GetCommittedTokenChan(context.Background())

	for {
		select {
		case now := <-stuckTimer.C:
			s.parent.Lock()
			for _, m := range s.commandActive {
				if m != nil {
					if now.After(m.timeoutTime) {
						s.parent.Unlock()

						return ErrorTimeout
					}
				}
			}
			s.parent.Unlock()
		case tokenRaw, ok := <-commitChan:
			if !ok {
				return ErrorWorkerClosed
			}
			workingToken = tokenRaw.(*commandToken)

			/* Can we send more commands? */
			s.parent.Lock()
			for s.parent.commandMaxIssue == 0 {
				err := s.waitCondition(s.parent.commandMaxIssueChanged, stuckTimer.C)
				if err != nil {
					return err
				}
			}

			/* The standard does not specify if commands with identical opcodes will get
			   sequential responses. As such we can't pair them and we cannot have multiple
			   issued commands with the same opcode. Here we wait to ensure that there are no
			   tokens with the same opcode */
			for s.findToken(workingToken.opcode, false) != nil {
				err := s.waitCondition(s.tokenRemoved, stuckTimer.C)
				if err != nil {
					return err
				}
			}

			txData := workingToken.data

			/* It is guaranteed there will be at least one nil element */
			for i := range s.commandActive {
				if s.commandActive[i] == nil {
					s.commandActive[i] = workingToken
					workingToken = nil
					break
				}
			}

			s.parent.commandMaxIssue--

			/* Check if we can unlock another queue */
			if !s.parent.closed && s.parent.commandMaxIssue > 0 {
				select {
				case s.parent.commandMaxIssueChanged <- struct{}{}:
				default:
				}
			}
			s.parent.Unlock()

			err := s.parent.transmitFunc(txData)
			if err != nil {
				return err
			}
		}
	}
}
