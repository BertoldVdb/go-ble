package hcicmdmgr

import (
	"encoding/binary"
)

func (s *CommandManager) updateIssueList(numCommands int) {
	/* According to the standard this can go down at arbitrary moments.
	   I don't really know how to handle that, potentially a command
	   is already going out while receiving the lower maximum...
	   I have not seen this behavior in practice. */

	s.Lock()
	s.commandMaxIssue = numCommands
	select {
	case s.commandMaxIssueChanged <- struct{}{}:
	default:
	}
	s.Unlock()
}

func (s *CommandManager) HandleEventCommandComplete(params []byte) error {
	var err error
	if len(params) >= 3 {
		numCommands := int(params[0])
		opcode := binary.LittleEndian.Uint16(params[1:3])

		if opcode > 0 {
			for i := range s.queues {
				err2 := s.queues[i].commandComplete(false, opcode, params[3:])
				if err2 != nil {
					err = err2
				}
			}
		}

		s.updateIssueList(numCommands)
	}
	return err
}

func (s *CommandManager) HandleEventCommandStatus(params []byte) error {
	var err error

	if len(params) == 4 {
		status := params[0:1]
		numCommands := int(params[1])
		opcode := binary.LittleEndian.Uint16(params[2:4])

		if opcode > 0 {
			for i := range s.queues {
				err2 := s.queues[i].commandComplete(true, opcode, status)
				if err2 != nil {
					err = err2
				}
			}
		}

		s.updateIssueList(numCommands)
	}
	return err
}
