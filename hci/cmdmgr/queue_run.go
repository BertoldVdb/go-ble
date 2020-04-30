package hcicmdmgr

import (
	"context"
	"log"
	"time"

	hciconst "github.com/BertoldVdb/go-ble/hci/const"
)

func (s *commandQueue) commandRunGetToken(cmd HCICommand, sync bool, cb CommandCompleteCallback) (*commandToken, error) {
	tokenRaw, err := s.commandQueue.GetAvailableToken(context.Background())
	if err != nil {
		return nil, err
	}

	token := tokenRaw.(*commandToken)

	token.opcode = uint16(cmd.OCF | cmd.OGF<<10)

	tmp := token.data[:0]
	tmp = append(tmp, hciconst.MsgTypeCommand)
	tmp = append(tmp, byte(token.opcode))
	tmp = append(tmp, byte(token.opcode>>8))
	tmp = append(tmp, 0)
	token.data = append(tmp, cmd.Params...)
	token.sync = sync
	token.cb = cb

	return token, nil
}

func (s *commandQueue) commandRunPutToken(token *commandToken) ([]byte, error) {
	token.timeoutTime = time.Now().Add(time.Second)
	token.data[3] = byte(len(token.data) - 4)

	log.Println(token.data)
	err := s.commandQueue.CommitToken(token)
	if err != nil {
		return nil, err
	}

	/* Wait for completion */
	if token.sync {
		cberr, ok := <-token.completed
		if !ok {
			return nil, ErrorWorkerClosed
		}
		log.Println(token.data)

		if cberr != nil {
			cberr = s.commandRunReleaseToken(token)
			return nil, cberr
		}

		return token.data, cberr
	}

	return nil, nil
}

func (s *commandQueue) commandRunReleaseToken(token *commandToken) error {
	if token.sync {
		return s.commandQueue.ReleaseToken(token)
	}

	return nil
}

func (s *commandQueue) commandRun(cmd HCICommand, output []byte, sync bool, cb CommandCompleteCallback) ([]byte, error) {
	token, err := s.commandRunGetToken(cmd, sync, cb)
	if err != nil {
		return nil, err
	}
	buf, err := s.commandRunPutToken(token)
	if err != nil {
		return nil, err
	}

	output = append(output, buf...)
	err = s.commandRunReleaseToken(token)
	return output, err
}

func (s *commandQueue) tokenComplete(token *commandToken, cberr error, params []byte) error {
	var err error

	if !token.sync {
		cb := token.cb

		if cb != nil {
			err = cb(cberr, params)
		}

		err2 := s.commandQueue.ReleaseToken(token)
		if err2 != nil {
			err = err2
		}
	} else {
		token.Lock()
		token.data = append(token.data[:0], params...)
		token.completed <- cberr
		token.Unlock()
	}

	return err
}

func (s *commandQueue) commandComplete(status bool, opcode uint16, params []byte) error {
	var err error

	s.parent.Lock()
	token := s.findToken(opcode, true)
	s.parent.Unlock()

	if token == nil {
		return err
	}

	err = s.tokenComplete(token, nil, params)

	return err
}
