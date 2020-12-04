package blel2cap

import (
	"context"
	"encoding/binary"

	"github.com/BertoldVdb/go-ble/bleconnecter"
	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/BertoldVdb/go-misc/tokenqueue"
	"github.com/sirupsen/logrus"
)

type signalling struct {
	l *L2CAP

	paramBuf []uint16

	cmdQueue    *tokenqueue.Queue
	activeToken *signallingCommandToken

	currentOutgoing uint8
}

type signallingCommandToken struct {
	ctx          context.Context
	cid          uint16
	code         uint8
	payload      *pdu.PDU
	completeChan chan (struct{})
}

func (sc *signallingCommandToken) Cleanup() {
	close(sc.completeChan)
}

func (s *signalling) signallingHandler(cid uint16, buf *pdu.PDU) (error, bool) {
	if buf == nil {
		s.cmdQueue.Close()
		return nil, false
	}

	header := buf.DropLeft(4)
	if header == nil {
		return nil, false
	}

	code := header[0]
	id := header[1]
	length := int(binary.LittleEndian.Uint16(header[2:4]))

	if length != buf.Len() {
		return nil, false
	}

	return s.signallingProcess(cid, code, id, buf)
}

func (s *signalling) signallingCommandUint16(payload *pdu.PDU, cid uint16, id uint8, responseCode uint8, params ...uint16) (error, bool) {
	payload.Reset()
	data := payload.ExtendRight(2 * len(params))
	for i, m := range params {
		binary.LittleEndian.PutUint16(data[2*i:], m)
	}

	return s.signallingWriteResponse(cid, responseCode, id, payload), true
}

func signallingErrorToUint16(err error, success uint16, failed uint16) uint16 {
	if err == nil {
		return success
	}
	return failed
}

const (
	SigCommandRejectRsp              = 0x1
	SigConnectionReq                 = 0x2
	SigConnectionRsp                 = 0x3
	SigConfigurationReq              = 0x4
	SigConfigurationRsp              = 0x5
	SigDisconnectionReq              = 0x6
	SigDisconnectionRsp              = 0x7
	SigEchoReq                       = 0x8
	SigEchoRsp                       = 0x9
	SigInformationReq                = 0xa
	SigInformationRsp                = 0xb
	SigConnectionParametereUpdateReq = 0x12
	SigConnectionParametereUpdateRsp = 0x13
)

func (s *signalling) signallingProcess(cid uint16, code uint8, id uint8, payload *pdu.PDU) (error, bool) {
	isResponse := (code <= 0x15 && code&1 == 1) || (code > 0x16 && code&1 == 0)
	isUint16Cmd := code != SigEchoReq

	if isResponse {
		if id == s.currentOutgoing && s.activeToken != nil {
			s.currentOutgoing++
			if s.currentOutgoing == 0 {
				s.currentOutgoing = 1
			}

			s.activeToken.code = code
			s.activeToken.payload = payload
			s.activeToken.completeChan <- struct{}{}
			s.activeToken = nil
			return nil, true
		}
		return nil, false
	} else {
		if s.l.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
			s.l.logger.WithFields(logrus.Fields{
				"0cid":     cid,
				"1code":    code,
				"2id":      id,
				"3payload": payload,
			}).Debug("Received L2CAP command")
		}

		if isUint16Cmd {
			s.paramBuf = s.paramBuf[:0]
			for payload.Len() >= 2 {
				s.paramBuf = append(s.paramBuf, binary.LittleEndian.Uint16(payload.DropLeft(2)))
			}
			if payload.Len() != 0 {
				return nil, false
			}
		}

		switch code {
		case SigConnectionReq:
			if len(s.paramBuf) == 2 {
				/* This is only for EDR */
				return s.signallingCommandUint16(payload, cid, id, SigConnectionRsp, 0, s.paramBuf[1], 4, 0)
			}

		case SigConfigurationReq:
			if len(s.paramBuf) >= 2 {
				return s.signallingCommandUint16(payload, cid, id, SigConfigurationRsp, 0, 0, 2)
			}

		case SigEchoReq:
			return s.signallingWriteResponse(cid, SigEchoRsp, id, payload), true

		case SigInformationReq:
			//Not used in LE

		case SigConnectionParametereUpdateReq:
			conn, ok := s.l.conn.(*bleconnecter.BLEConnection)
			if ok && len(s.paramBuf) == 4 {
				result := conn.UpdateParams(bleconnecter.BLEConnectionParametersRequested{
					ConnectionIntervalMin: s.paramBuf[0],
					ConnectionIntervalMax: s.paramBuf[1],
					ConnectionLatency:     s.paramBuf[2],
					SupervisionTimeout:    s.paramBuf[3],
				})

				return s.signallingCommandUint16(payload, cid, id, SigConnectionParametereUpdateRsp, signallingErrorToUint16(result, 0, 1))
			}
		}

		return s.signallingCommandUint16(payload, cid, id, SigCommandRejectRsp, 0)
	}

	return nil, false
}

func (s *signalling) signallingWriteResponse(cid uint16, code uint8, id uint8, payload *pdu.PDU) error {
	header := payload.ExtendLeft(4)
	header[0] = code
	header[1] = id
	binary.LittleEndian.PutUint16(header[2:4], uint16(payload.Len()-4))

	if s.l.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		s.l.logger.WithFields(logrus.Fields{
			"0cid":     cid,
			"1code":    code,
			"2id":      id,
			"3payload": payload,
		}).Debug("Sending L2CAP signalling")
	}

	return s.l.WriteBufferCID(cid, payload)
}

func (l *L2CAP) getSignallingChannel() uint16 {
	if l.isLE {
		return 5
	}

	return 1
}

func (l *L2CAP) signallingInit() *signalling {
	s := &signalling{
		l: l,
		cmdQueue: tokenqueue.NewQueue(1, 1, func() tokenqueue.Token {
			return &signallingCommandToken{
				completeChan: make(chan (struct{})),
			}
		}),
		currentOutgoing: 0xBE,
	}

	l.registerCallbackCID(l.getSignallingChannel(), s.signallingHandler)

	return s
}

func (s *signalling) Close() error {
	s.cmdQueue.Close()
	return nil
}

func (s *signalling) sendCommand(ctx context.Context, cid uint16, code uint8, payload *pdu.PDU) (uint8, *pdu.PDU, error) {
	t, err := s.cmdQueue.GetAvailableToken(ctx)
	if err != nil {
		return 0, nil, err
	}
	token := t.(*signallingCommandToken)

	token.ctx = ctx
	token.cid = cid
	token.code = code
	if payload == nil {
		payload = bleutil.GetBuffer(0)
	}
	token.payload = payload

	err = s.cmdQueue.CommitToken(t)
	if err != nil {
		return 0, nil, err
	}

	_, ok := <-token.completeChan
	if !ok {
		return 0, nil, ErrorRxFailed
	}

	if s.l.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		s.l.logger.WithError(ctx.Err()).WithFields(logrus.Fields{
			"0cid":     token.cid,
			"1code":    token.code,
			"2payload": token.payload,
		}).Debug("Completed L2CAP command")
	}

	code = token.code
	payload = token.payload
	err = s.cmdQueue.ReleaseToken(t)
	if err != nil {
		return 0, nil, err
	}

	return code, payload, ctx.Err()
}

func (s *signalling) handleTokenTx(t tokenqueue.Token) {
	token := t.(*signallingCommandToken)
	if token.payload != nil {
		s.activeToken = token
		s.signallingWriteResponse(token.cid, token.code, s.currentOutgoing, token.payload)

	} else {
		token.completeChan <- struct{}{}
	}
}

func (s *signalling) tokenTimeoutChan() <-chan (struct{}) {
	if s.activeToken != nil {
		return s.activeToken.ctx.Done()
	}
	return nil
}

func (s *signalling) failActiveToken() {
	if s.activeToken != nil {
		s.activeToken.payload = nil
		s.activeToken.completeChan <- struct{}{}
		s.activeToken = nil
	}
}

func (s *signalling) SendCommandUint16(ctx context.Context, code uint8, result []uint16, params ...uint16) (uint8, []uint16, error) {
	cid := s.l.getSignallingChannel()
	result = result[:]

	payload := bleutil.GetBuffer(2 * len(params))
	for i, m := range params {
		binary.LittleEndian.PutUint16(payload.Buf()[2*i:], m)
	}

	responseCode, responsePDU, err := s.sendCommand(ctx, cid, code, payload)

	if err == nil && responsePDU != nil {
		for responsePDU.Len() >= 2 {
			result = append(result, binary.LittleEndian.Uint16(responsePDU.DropLeft(2)))
		}
	}

	/* If we get any error here the L2CAP state is damaged and we need to close the link */
	if err != nil {
		s.l.Close()
	}

	return responseCode, result, err
}
