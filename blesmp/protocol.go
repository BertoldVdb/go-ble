package blesmp

import (
	"encoding/binary"
	"encoding/hex"
	"log"

	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"

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
	pairingLTK         [16]byte
	pairingEDIVValid   bool
	pairingEDIV        uint16
	pairingRand        uint64
}

func (c *SMPConn) sendPairingRequest(authReq byte) {
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

func (c *SMPConn) sendPairingFailed(reason psmFailedReason) {
	log.Println("Pairing failed", reason, c.protocol.state)
	//TODO: signal higher layers
	c.protocol.state = smpWaitStart

	buf := bleutil.GetBuffer(2)
	buf.Buf()[0] = byte(opcodePairingFailed)
	buf.Buf()[1] = byte(reason)
	c.conn.WriteBuffer(buf)

}

func (c *SMPConn) handleSecurityRequest(authReq byte) {

}

func (c *SMPConn) handlePairingResponse(resp []byte) {
	c.protocol.state = smpWaitPairingConfirm

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

	raw := c.rawConnLE()
	remote := raw.RemoteAddr().(bleutil.BLEAddr)
	local := raw.LocalAddr().(bleutil.BLEAddr)

	_, err := crand.Read(c.protocol.pairingIRand[:])
	if err != nil {
		c.sendPairingFailed(failedUnspecifiedReason)
		return
	}

	c.protocol.pairingIConfirm = CryptoFuncC1(c.protocol.pairingTK, c.protocol.pairingIRand, c.protocol.pairingRequest, c.protocol.pairingResponse, local, remote)

	buf := bleutil.GetBuffer(1)
	buf.Buf()[0] = byte(opcodePairingConfirm)
	buf.Append(c.protocol.pairingIConfirm[:]...)
	c.conn.WriteBuffer(buf)
}

func (c *SMPConn) handlePairingConfirm(resp []byte) {
	c.protocol.state = smpWaitPairingRandom

	copy(c.protocol.pairingRConfirm[:], resp)

	buf := bleutil.GetBuffer(1)
	buf.Buf()[0] = byte(opcodePairingRandom)
	buf.Append(c.protocol.pairingIRand[:]...)
	c.conn.WriteBuffer(buf)
}

func (c *SMPConn) handlePairingRandom(rand []byte) {
	c.protocol.state = smpKeyDistribution

	var pairingRRand [16]byte
	copy(pairingRRand[:], rand)

	raw := c.rawConnLE()
	remote := raw.RemoteAddr().(bleutil.BLEAddr)
	local := raw.LocalAddr().(bleutil.BLEAddr)

	tmp := CryptoFuncC1(c.protocol.pairingTK, pairingRRand, c.protocol.pairingRequest, c.protocol.pairingResponse, local, remote)

	if subtle.ConstantTimeCompare(tmp[:], c.protocol.pairingRConfirm[:]) == 0 {
		c.sendPairingFailed(failedConfirmValueFailed)
		return
	}

	/* The legacy pairing is completed, calculate STK and encrypt link to continue */
	c.protocol.pairingSTK = CryptoFuncS1(c.protocol.pairingTK, c.protocol.pairingIRand, pairingRRand)
	c.protocol.pairingSTK = CryptoShortenKey(c.protocol.pairingSTK, c.protocol.pairingKeySize)

	raw.Encrypt(0, 0, c.protocol.pairingSTK)
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

		log.Println("Received LTK", c.protocol.pairingEDIV, c.protocol.pairingRand, hex.EncodeToString(c.protocol.pairingLTK[:]))

		//TODO: Using the LTK on this connection makes no sense (it is already encrypted!), it is just to test that it is correct!
		raw := c.rawConnLE()
		raw.Encrypt(c.protocol.pairingEDIV, c.protocol.pairingRand, c.protocol.pairingLTK)
	}

	if opcode == opcodeKDInitiatorIdentification && len(data) == 10 {
		c.protocol.pairingEDIVValid = true
		c.protocol.pairingEDIV = binary.LittleEndian.Uint16(data)
		c.protocol.pairingRand = binary.LittleEndian.Uint64(data[2:])

		handleLTK()
		return true
	}
	if opcode == opcodeKDEncryptionInformation && len(data) == 16 {
		c.protocol.pairingLTKValid = true
		copy(c.protocol.pairingLTK[:], data)

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

		c.protocol.state = smpWaitStart

		//TODO: notify if needed
		log.Println("Pairing failed", reason)
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
