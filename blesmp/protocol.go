package blesmp

import (
	"encoding/binary"
	"encoding/hex"

	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
	"github.com/sirupsen/logrus"

	crand "crypto/rand"
	"crypto/subtle"
)

type smpFsmState int

const (
	smpWaitStart           smpFsmState = 0
	smpWaitPairingResponse smpFsmState = 1
	smpWaitPairingConfirm  smpFsmState = 2
	smpWaitPairingRandom   smpFsmState = 3
	smpKeyDistribution     smpFsmState = 4
)

type smpProtocol struct {
	state smpFsmState

	pairingRequest  [7]byte
	pairingResponse [7]byte

	/* Legacy pairing */
	pairingTK       [16]byte
	pairingIRand    [16]byte
	pairingIConfirm [16]byte
	pairingRConfirm [16]byte
	pairingKeySize  int
	pairingSTK      [16]byte

	/* Key distribution */
	pairingLTKComplete bool
	pairingLTKValid    bool
	pairingEDIVValid   bool
	pairingLTK         smpStoredLTK
}

func (c *SMPConn) sendPairingRequest(authReq byte) {
	c.updateTimeout(true)

	c.protocol.state = smpWaitPairingResponse
	c.protocol.pairingLTKValid = false
	c.protocol.pairingEDIVValid = false

	c.protocol.pairingKeySize = 16

	c.protocol.pairingRequest = [7]byte{
		byte(opcodePairingRequest),
		c.parent.ioCapability,
		0,
		authReq,
		byte(c.protocol.pairingKeySize),
		0,
		1, /*EncKey */
	}

	buf := bleutil.GetBuffer(len(c.protocol.pairingRequest))
	copy(buf.Buf(), c.protocol.pairingRequest[:])
	c.conn.WriteBuffer(buf)
}

func (c *SMPConn) signalPairingFailed(reason psmFailedReason) {
	c.logger.WithFields(logrus.Fields{"0reason": reason, "1state": c.protocol.state}).Warn("Pairing failed")
	c.protocol.state = smpWaitStart
	c.updateTimeout(false)

	//TODO:signal higher layers
}

func (c *SMPConn) sendPairingFailed(reason psmFailedReason) {
	buf := bleutil.GetBuffer(2)
	buf.Buf()[0] = byte(opcodePairingFailed)
	buf.Buf()[1] = byte(reason)
	c.conn.WriteBuffer(buf)

	c.signalPairingFailed(reason)
}

func (c *SMPConn) handleSecurityRequest(authReq byte) {

}

func (c *SMPConn) handlePairingResponse(resp []byte) {
	copy(c.protocol.pairingResponse[:], resp)

	//TODO: check IO capabilities

	ks := int(resp[4])
	if ks < c.protocol.pairingKeySize {
		c.protocol.pairingKeySize = ks
	}
	if c.protocol.pairingKeySize < 7 {
		c.sendPairingFailed(failedEncryptionKeySize)
		return
	}

	_, err := crand.Read(c.protocol.pairingIRand[:])
	if err != nil {
		c.sendPairingFailed(failedUnspecifiedReason)
		return
	}

	c.protocol.pairingIConfirm = CryptoFuncC1(c.protocol.pairingTK, c.protocol.pairingIRand, c.protocol.pairingRequest, c.protocol.pairingResponse, c.addrLELocal, c.addrLERemote)

	buf := bleutil.GetBuffer(1)
	buf.Buf()[0] = byte(opcodePairingConfirm)
	buf.Append(c.protocol.pairingIConfirm[:]...)
	c.conn.WriteBuffer(buf)

	c.protocol.state = smpWaitPairingConfirm
}

func (c *SMPConn) handlePairingConfirm(resp []byte) {
	copy(c.protocol.pairingRConfirm[:], resp)

	buf := bleutil.GetBuffer(1)
	buf.Buf()[0] = byte(opcodePairingRandom)
	buf.Append(c.protocol.pairingIRand[:]...)
	c.conn.WriteBuffer(buf)

	c.protocol.state = smpWaitPairingRandom
}

func (c *SMPConn) handlePairingRandom(rand []byte) {
	var pairingRRand [16]byte
	copy(pairingRRand[:], rand)

	tmp := CryptoFuncC1(c.protocol.pairingTK, pairingRRand, c.protocol.pairingRequest, c.protocol.pairingResponse, c.addrLELocal, c.addrLERemote)

	if subtle.ConstantTimeCompare(tmp[:], c.protocol.pairingRConfirm[:]) == 0 {
		c.sendPairingFailed(failedConfirmValueFailed)
		return
	}

	/* The legacy pairing is completed, calculate STK and encrypt link to continue */
	c.protocol.pairingSTK = CryptoFuncS1(c.protocol.pairingTK, c.protocol.pairingIRand, pairingRRand)
	c.protocol.pairingSTK = CryptoShortenKey(c.protocol.pairingSTK, c.protocol.pairingKeySize)

	c.logger.WithFields(logrus.Fields{
		"0stk": hex.EncodeToString(c.protocol.pairingSTK[:]),
	}).Info("STK Calculated")

	raw := c.rawConnLE()
	raw.Encrypt(0, 0, c.protocol.pairingSTK)

	c.protocol.state = smpKeyDistribution
}

func (c *SMPConn) handleKeyDistribution(opcode psmOpcode, data []byte) bool {
	handleLTK := func() {
		if !c.protocol.pairingEDIVValid || !c.protocol.pairingLTKValid {
			return
		}

		if c.protocol.pairingLTKComplete {
			return
		}

		c.protocol.pairingLTKComplete = true

		c.parent.storedKeysPersist.Lock()
		c.parent.storedKeys[makeSMPStoredLTKMapKey(c.addrLELocal, c.addrLERemote)] = c.protocol.pairingLTK
		c.parent.storedKeysPersist.Unlock()
		err := c.parent.storedKeysPersist.Save()

		c.logger.WithError(err).WithFields(logrus.Fields{
			"0ediv": c.protocol.pairingLTK.EDIV,
			"1rand": c.protocol.pairingLTK.Rand,
			"2ltk":  hex.EncodeToString(c.protocol.pairingLTK.LTK[:]),
		}).Info("LTK Saved")

		//TODO: Using the LTK on this connection makes no sense (it is already encrypted!), it is just to test that it is correct!
		//Just store it for next time: OriginatorMac/DestinationMac -> EDIV, RAND, LTK
		//raw := c.rawConnLE()
		//raw.Encrypt(c.protocol.pairingEDIV, c.protocol.pairingRand, c.protocol.pairingLTK)
	}

	if opcode == opcodeKDInitiatorIdentification && len(data) == 10 {
		c.protocol.pairingEDIVValid = true
		c.protocol.pairingLTK.EDIV = binary.LittleEndian.Uint16(data)
		c.protocol.pairingLTK.Rand = binary.LittleEndian.Uint64(data[2:])

		handleLTK()
		return true
	}
	if opcode == opcodeKDEncryptionInformation && len(data) == 16 {
		c.protocol.pairingLTKValid = true
		copy(c.protocol.pairingLTK.LTK[:], data)

		handleLTK()
		return true
	}
	return false
}

func (c *SMPConn) handleMessage(pdu *pdu.PDU) bool {
	if pdu.Len() < 1 {
		return false
	}

	opcode := psmOpcode(pdu.Buf()[0])
	/* This opcode is supported in all states */
	if opcode == opcodePairingFailed {
		reason := psmFailedReason(0)
		if pdu.Len() >= 2 {
			reason = psmFailedReason(pdu.Buf()[1])
		}

		c.signalPairingFailed(reason)
		return false
	}

	switch c.protocol.state {
	case smpWaitStart:
		if opcode == opcodeKDSecurityRequest && pdu.Len() == 2 {
			c.handleSecurityRequest(pdu.Buf()[1])
			return false
		}
	case smpWaitPairingResponse:
		if opcode == opcodePairingResponse && pdu.Len() == 7 {
			c.handlePairingResponse(pdu.Buf())
			return false
		}
	case smpWaitPairingConfirm:
		if opcode == opcodePairingConfirm && pdu.Len() == 17 {
			c.handlePairingConfirm(pdu.Buf()[1:])
			return false
		}
	case smpWaitPairingRandom:
		if opcode == opcodePairingRandom && pdu.Len() == 17 {
			c.handlePairingRandom(pdu.Buf()[1:])
			return false
		}
	case smpKeyDistribution:
		if c.handleKeyDistribution(opcode, pdu.Buf()[1:]) {
			return false
		}
	}

	c.sendPairingFailed(failedCommandNotSupported)

	return false
}
