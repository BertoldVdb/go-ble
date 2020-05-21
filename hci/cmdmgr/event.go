package hcicmdmgr

import (
	"encoding/hex"

	"github.com/sirupsen/logrus"
)

func (s *CommandManager) updateIssueList(numCommands int) {
	/* According to the standard this can go down at arbitrary moments.
	   I don't really know how to handle that, potentially a command
	   is already going out while receiving the lower maximum...
	   I have not seen this behavior in practice. */

	s.Lock()
	debugOldMaxIssue := s.commandMaxIssue
	s.commandMaxIssue = numCommands
	select {
	case s.commandMaxIssueChanged <- struct{}{}:
	default:
	}
	s.Unlock()

	if s.logger != nil && s.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		s.logger.WithFields(logrus.Fields{
			"0slots":    numCommands,
			"1oldslots": debugOldMaxIssue,
		}).Trace("Number of slots updated")
	}
}

func (s *CommandManager) commandComplete(CommandOpcode uint16, NumHCICommandPackets uint8, ReturnParameters []byte) error {
	var err error

	if s.logger != nil && s.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		s.logger.WithFields(logrus.Fields{
			"0opcode": CommandOpcode,
			"1slots":  NumHCICommandPackets,
			"2return": hex.EncodeToString(ReturnParameters),
		}).Trace("Command complete/status event received")
	}

	if CommandOpcode > 0 {
		for i := range s.queues {
			err2 := s.queues[i].commandComplete(false, CommandOpcode, ReturnParameters)
			if err2 != nil {
				err = err2
			}
		}
	}

	s.updateIssueList(int(NumHCICommandPackets))

	return err
}

func (s *CommandManager) HandleEventCommandComplete(CommandOpcode uint16, NumHCICommandPackets uint8, ReturnParameters []byte) error {
	return s.commandComplete(CommandOpcode, NumHCICommandPackets, ReturnParameters)
}

func (s *CommandManager) HandleEventCommandStatus(CommandOpcode uint16, NumHCICommandPackets uint8, Status uint8) error {
	return s.commandComplete(CommandOpcode, NumHCICommandPackets, []byte{Status})

}
