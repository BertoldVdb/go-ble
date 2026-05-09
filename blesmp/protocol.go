package blesmp

import (
	"encoding/binary"

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

	// LE Secure Connections sub-states
	smpSCWaitPublicKey      smpFsmState = 10
	smpSCWaitConfirm        smpFsmState = 11
	smpSCWaitRandom         smpFsmState = 12
	smpSCWaitNumericConfirm smpFsmState = 13
	smpSCWaitDHKeyCheck     smpFsmState = 14
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

	/* LE Secure Connections — populated only when both sides agree on
	   the SC bit during the Pairing Request/Response exchange. All
	   multi-byte values are stored in big-endian (the spec's natural
	   byte order); reversal happens at the wire boundary. */
	scActive       bool
	scKeyPair      *scKeyPair
	scPeerPubX     [32]byte
	scPeerPubY     [32]byte
	scDHKey        [32]byte
	scLocalNonce   [16]byte
	scRemoteNonce  [16]byte
	scLocalConfirm [16]byte
	scPeerConfirm  [16]byte
	scAlgorithm    int    // 0=JustWorks, 1=Passkey, 2=NumericComparison, 3=OOB
	scPasskey      uint32 // user-entered or randomly displayed
	scPasskeyBit   int
	scMacKey       [16]byte
	scLTK          [16]byte
	scNumericValue uint32 // result of g2() for the user prompt
	scPeerSentRand bool   // peripheral has received initiator's random
}

func (c *SMPConn) sendBuf(pdu *pdu.PDU) error {
	if c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		opcode := byte(0)
		if pdu.Len() > 0 {
			opcode = pdu.Buf()[0]
		}
		c.logger.WithFields(logrus.Fields{
			"0opcode": opcode,
			"1len":    pdu.Len(),
		}).Debug("SMP send")
	}

	return c.conn.WriteBuffer(pdu)
}

func (c *SMPConn) getIOCapability() smpIOCapability {
	if c.config.InputNumeric != nil {
		if c.config.DisplayNumeric != nil {
			return cIOKeyboardDisplay
		}

		return cIOKeyboardOnly
	}

	if c.config.DisplayNumeric != nil {
		if c.config.InputYesNo != nil {
			return cIODisplayYesNo
		}

		return cIODisplayOnly
	}

	return cIONoInputNoOutput
}

func getLegacyAlgorithmType(initiatorIO smpIOCapability, responderIO smpIOCapability, initiatorMITM bool, responderMITM bool) (int, int) {
	if !initiatorMITM && !responderMITM {
		return 0, 0
	}

	hasKeyboard := func(cap smpIOCapability) bool {
		return cap == cIOKeyboardDisplay || cap == cIOKeyboardOnly
	}
	hasDisplay := func(cap smpIOCapability) bool {
		return cap == cIODisplayOnly || cap == cIODisplayYesNo
	}

	if hasKeyboard(initiatorIO) {
		/* Initiator can enter a code, passkey can be used unless the responder has nothing at all */
		if hasDisplay(responderIO) {
			return 1, 2
		}
		if hasKeyboard(responderIO) {
			return 1, 1
		}
	} else if hasDisplay(initiatorIO) {
		/* Initiator can display, but this only works if the responder has a keyboard */
		if hasKeyboard(responderIO) {
			return 2, 1
		}
	}

	return 0, 0
}

func (c *SMPConn) sendPairingRequestResponse(initiator bool) {
	c.updateTimeout(true)

	c.protocol.state = smpWaitPairingResponse
	c.protocol.pairingLTKValid = false
	c.protocol.pairingEDIVValid = false
	c.protocol.pairingLTKComplete = false
	c.protocol.pairingKeySize = 16
	c.protocol.pairingLTK.Authenticated = false

	authReq := c.secureAuthReq
	// Always advertise SC capability when the implementation can do it
	// (the LE SC code path is implemented; AuthReq bit 3 = SC).
	authReq |= 0x08

	req := [7]byte{
		byte(opcodePairingRequest),
		uint8(c.getIOCapability()),
		0,
		authReq,
		byte(c.protocol.pairingKeySize),
	}

	/* Ask for EncKey if we will bond */
	if c.secureAuthReq&0x1 > 0 {
		if initiator {
			req[6] = 1
		} else {
			/* Only if other side asked for it */
			req[6] = c.protocol.pairingRequest[6] & 1
		}
	}

	if initiator {
		c.protocol.pairingRequest = req
	} else {
		req[0] = byte(opcodePairingResponse)
		c.protocol.pairingResponse = req
	}

	buf := bleutil.GetBuffer(len(req))
	copy(buf.Buf(), req[:])
	c.sendBuf(buf)
}

func (c *SMPConn) sendPairingRequest() {
	c.sendPairingRequestResponse(true)
}

func (c *SMPConn) signalPairingFailed(reason smpFailedReason) {
	c.logger.WithFields(logrus.Fields{"0reason": reason, "1state": c.protocol.state}).Warn("Pairing failed")
	c.protocol.state = smpWaitStart
	c.updateTimeout(false)

	/* Wipe in-progress LTK material so a stale partial LTK from this
	   pairing attempt cannot leak into a subsequent attempt or be
	   mistakenly accepted as bonded. */
	c.protocol.pairingLTK = smpStoredLTK{}
	c.protocol.pairingLTKValid = false
	c.protocol.pairingEDIVValid = false
	c.protocol.pairingLTKComplete = false
	c.protocol.pairingTK = [16]byte{}
	c.protocol.pairingSTK = [16]byte{}

	/* Wipe LE Secure Connections state too. */
	c.protocol.scActive = false
	c.protocol.scKeyPair = nil
	c.protocol.scDHKey = [32]byte{}
	c.protocol.scLocalNonce = [16]byte{}
	c.protocol.scRemoteNonce = [16]byte{}
	c.protocol.scLocalConfirm = [16]byte{}
	c.protocol.scPeerConfirm = [16]byte{}
	c.protocol.scMacKey = [16]byte{}
	c.protocol.scLTK = [16]byte{}
	c.protocol.scPasskey = 0
	c.protocol.scPasskeyBit = 0

	c.setState(StateFailed)
}

func (c *SMPConn) sendPairingFailed(reason smpFailedReason) {
	buf := bleutil.GetBuffer(2)
	buf.Buf()[0] = byte(opcodePairingFailed)
	buf.Buf()[1] = byte(reason)
	c.sendBuf(buf)

	c.signalPairingFailed(reason)
}

func (c *SMPConn) handleSecurityRequest(authReq byte) {
	select {
	case c.secureChan <- struct{}{}:
	default:
	}
}

func (c *SMPConn) sendSecurityRequest() {
	c.updateTimeout(true)

	buf := bleutil.GetBuffer(2)
	buf.Buf()[0] = byte(opcodeKDSecurityRequest)
	buf.Buf()[1] = byte(c.secureAuthReq)
	c.sendBuf(buf)
}

func (c *SMPConn) handleStageConfirm(initiator bool) bool {
	algorithm, algorithmResponder := getLegacyAlgorithmType(
		smpIOCapability(c.protocol.pairingRequest[1]),
		smpIOCapability(c.protocol.pairingResponse[1]),
		c.protocol.pairingRequest[3]&0x4 > 0,
		c.protocol.pairingResponse[3]&0x4 > 0)

	var key uint32
	var err error
	if !initiator {
		algorithm = algorithmResponder
	}

	c.protocol.pairingLTK.Authenticated = algorithm > 0

	if algorithm == 1 {
		/* Ask input */
		key, err = c.config.InputNumeric(c)
	} else if algorithm == 2 {
		/* Display */

		key, err = c.CryptoGeneratePassKey()
		if err != nil {
			c.sendPairingFailed(failedUnspecifiedReason)
			return false
		}

		err = c.config.DisplayNumeric(c, key)
	}

	if err != nil {
		c.sendPairingFailed(failedPasskeyEntryFailed)
		return false
	}

	c.protocol.pairingTK = [16]byte{}
	binary.LittleEndian.PutUint32(c.protocol.pairingTK[:], key)

	kl1 := int(c.protocol.pairingRequest[4])
	kl2 := int(c.protocol.pairingResponse[4])

	if kl1 > kl2 {
		kl1 = kl2
	}

	c.protocol.pairingKeySize = kl1

	/* The spec floor is 7 bytes (56 bits), but accepting it silently is the
	   KNOB attack surface. Apply our local MinKeySize policy on top. */
	minKey := c.config.MinKeySize
	if minKey <= 0 {
		minKey = 16
	}
	if minKey < 7 {
		minKey = 7
	}
	if c.protocol.pairingKeySize < minKey || c.protocol.pairingKeySize > 16 {
		c.logger.WithFields(logrus.Fields{
			"0negotiated": c.protocol.pairingKeySize,
			"1minimum":    minKey,
		}).Warn("Refusing to accept downgraded LTK key size")
		c.sendPairingFailed(failedEncryptionKeySize)
		return false
	}

	_, err = crand.Read(c.protocol.pairingIRand[:])
	if err != nil {
		c.sendPairingFailed(failedUnspecifiedReason)
		return false
	}

	c.protocol.pairingIConfirm = CryptoFuncC1(c.isCentral, c.protocol.pairingTK, c.protocol.pairingIRand, c.protocol.pairingRequest, c.protocol.pairingResponse, c.addrLELocal, c.addrLERemote)

	return true
}

func (c *SMPConn) handlePairingRequest(req []byte) {
	copy(c.protocol.pairingRequest[:], req)
	c.sendPairingRequestResponse(false)

	if c.scNegotiated() {
		c.protocol.scActive = true
		c.scStartFromResponder()
		return
	}

	if !c.handleStageConfirm(false) {
		return
	}
}

func (c *SMPConn) handlePairingResponse(resp []byte) {
	copy(c.protocol.pairingResponse[:], resp)

	if c.scNegotiated() {
		c.protocol.scActive = true
		c.scStartFromInitiator()
		return
	}

	if !c.handleStageConfirm(true) {
		return
	}

	buf := bleutil.GetBuffer(1)
	buf.Buf()[0] = byte(opcodePairingConfirm)
	buf.Append(c.protocol.pairingIConfirm[:]...)
	c.sendBuf(buf)

	c.protocol.state = smpWaitPairingConfirm
}

// scNegotiated reports whether both sides set the SC bit in AuthReq —
// the precondition for the LE Secure Connections flow.
func (c *SMPConn) scNegotiated() bool {
	return c.protocol.pairingRequest[3]&0x08 != 0 && c.protocol.pairingResponse[3]&0x08 != 0
}

func (c *SMPConn) handlePairingConfirm(resp []byte) {
	copy(c.protocol.pairingRConfirm[:], resp)

	buf := bleutil.GetBuffer(1)
	if c.isCentral {
		buf.Buf()[0] = byte(opcodePairingRandom)
		buf.Append(c.protocol.pairingIRand[:]...)
	} else {
		buf.Buf()[0] = byte(opcodePairingConfirm)
		buf.Append(c.protocol.pairingIConfirm[:]...)
	}
	c.sendBuf(buf)

	c.protocol.state = smpWaitPairingRandom
}

func (c *SMPConn) handlePairingRandom(rand []byte) {
	var pairingRRand [16]byte
	copy(pairingRRand[:], rand)

	tmp := CryptoFuncC1(c.isCentral, c.protocol.pairingTK, pairingRRand, c.protocol.pairingRequest, c.protocol.pairingResponse, c.addrLELocal, c.addrLERemote)

	if subtle.ConstantTimeCompare(tmp[:], c.protocol.pairingRConfirm[:]) == 0 {
		c.sendPairingFailed(failedConfirmValueFailed)
		return
	}

	/* The legacy pairing is completed, calculate STK and encrypt link to continue */
	c.protocol.pairingSTK = CryptoFuncS1(c.isCentral, c.protocol.pairingTK, c.protocol.pairingIRand, pairingRRand)
	c.protocol.pairingSTK = CryptoShortenKey(c.protocol.pairingSTK, c.protocol.pairingKeySize)

	c.logger.Info("STK calculated")

	c.updateTimeout(false)

	if c.isCentral {
		err := c.leEncrypt(smpStoredLTK{
			LTK:           c.protocol.pairingSTK,
			Authenticated: c.protocol.pairingLTK.Authenticated,
		})
		if err != nil {
			c.sendPairingFailed(failedUnspecifiedReason)
			return
		}

	} else {
		var bytes [10]byte
		_, err := crand.Read(bytes[:])
		if err != nil {
			c.sendPairingFailed(failedUnspecifiedReason)
			return
		}

		c.protocol.pairingLTK.EDIV = binary.LittleEndian.Uint16(bytes[:])
		c.protocol.pairingLTK.Rand = binary.LittleEndian.Uint64(bytes[2:])
		c.protocol.pairingLTK.LTK = c.protocol.pairingSTK
		c.protocol.pairingLTK.Bonded = c.protocol.pairingResponse[6]&1 > 0
		c.updateLTK()

		err = c.leEncryptWait(func() error {
			/* We are ready, send the final message so the other side will turn on encryption */
			buf := bleutil.GetBuffer(1)
			buf.Buf()[0] = byte(opcodePairingRandom)
			buf.Append(c.protocol.pairingIRand[:]...)
			return c.sendBuf(buf)
		})
		if err != nil {
			c.sendPairingFailed(failedUnspecifiedReason)
			return
		}

		if c.protocol.pairingLTK.Bonded {
			/* Send LTK, then EDIV/RAND */
			buf := bleutil.GetBuffer(1)
			buf.Buf()[0] = byte(opcodeKDEncryptionInformation)
			buf.Append(c.protocol.pairingLTK.LTK[:]...)
			c.sendBuf(buf)

			buf = bleutil.GetBuffer(1)
			buf.Buf()[0] = byte(opcodeKDIdentification)
			binary.LittleEndian.PutUint16(buf.ExtendRight(2), c.protocol.pairingLTK.EDIV)
			binary.LittleEndian.PutUint64(buf.ExtendRight(8), c.protocol.pairingLTK.Rand)
			c.sendBuf(buf)
		}
	}

	c.setState(StateSecure)

	c.protocol.state = smpKeyDistribution
}

func (c *SMPConn) updateLTK() {
	/* In production c.parent is always set. Tests construct an SMPConn
	   in isolation; for those we just update the in-memory key flags
	   and skip the shared key store. */
	if c.parent == nil {
		c.leSetKeyFlagsFromLTK(c.protocol.pairingLTK)
		return
	}

	c.parent.storedKeysPersist.Lock()
	if !c.isCentral {
		c.parent.storedKeys[makeSMPStoredLTKMapKey(c.isCentral, c.addrLELocal, c.addrLERemote, 0, 0)] = c.protocol.pairingLTK
	}
	c.parent.storedKeys[makeSMPStoredLTKMapKey(c.isCentral, c.addrLELocal, c.addrLERemote, c.protocol.pairingLTK.EDIV, c.protocol.pairingLTK.Rand)] = c.protocol.pairingLTK
	c.parent.storedKeysPersist.Unlock()

	/* Persist to disk only for bonded keys. Non-bonded LTKs are kept
	   in-memory so the LELongTermKeyRequest handler can find them for
	   the duration of the session, but should not survive a restart. */
	var err error
	if c.protocol.pairingLTK.Bonded {
		err = c.parent.storedKeysPersist.Save()
	}

	c.leSetKeyFlagsFromLTK(c.protocol.pairingLTK)

	c.logger.WithError(err).WithFields(logrus.Fields{
		"0ediv":   c.protocol.pairingLTK.EDIV,
		"1rand":   c.protocol.pairingLTK.Rand,
		"2bonded": c.protocol.pairingLTK.Bonded,
		"3auth":   c.protocol.pairingLTK.Authenticated,
	}).Info("LTK saved")
}

func (c *SMPConn) handleKeyDistribution(opcode smpOpcode, data []byte) bool {
	handleLTK := func() {
		if !c.protocol.pairingEDIVValid || !c.protocol.pairingLTKValid {
			return
		}

		if c.protocol.pairingLTKComplete {
			return
		}

		c.protocol.pairingLTKComplete = true
		c.protocol.pairingLTK.Bonded = true

		c.updateLTK()
	}

	if opcode == opcodeKDIdentification && len(data) == 10 {
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
	/* Other valid key-distribution PDUs (IRK, IdentityAddress, CSRK) are not
	   yet processed by this implementation. Silently accept them — the spec
	   permits the peer to send any keys it advertised in the pairing
	   request/response, and failing the connection here would prevent
	   pairing with stacks that always send IRK/CSRK. */
	switch opcode {
	case opcodeKDIdentityInformation:
		if len(data) == 16 {
			return true
		}
	case opcodeKDIdentityAddressInformation:
		if len(data) == 7 {
			return true
		}
	case opcodeKDSigningInformation:
		if len(data) == 16 {
			return true
		}
	}
	return false
}

func (c *SMPConn) handleMessage(pdu *pdu.PDU) bool {
	if pdu.Len() < 1 {
		return false
	}

	opcode := smpOpcode(pdu.Buf()[0])
	/* This opcode is supported in all states */
	if opcode == opcodePairingFailed {
		reason := smpFailedReason(0)
		if pdu.Len() >= 2 {
			reason = smpFailedReason(pdu.Buf()[1])
		}

		c.signalPairingFailed(reason)
		return false
	}

	// Route SC-flow PDUs to the SC state machine when SC was negotiated.
	if c.protocol.scActive {
		if c.scHandleMessage(opcode, pdu.Buf()[1:]) {
			return false
		}
		// Unrecognised opcode for current SC state.
		c.sendPairingFailed(failedCommandNotSupported)
		return false
	}

	if c.isCentral {
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

	switch c.protocol.state {
	case smpWaitStart:
		if opcode == opcodePairingRequest && pdu.Len() == 7 {
			c.handlePairingRequest(pdu.Buf())
			return false
		}
	case smpWaitPairingResponse:
		if opcode == opcodePairingConfirm && pdu.Len() == 17 {
			c.handlePairingConfirm(pdu.Buf()[1:])
			return false
		}
	case smpWaitPairingRandom:
		if opcode == opcodePairingRandom && pdu.Len() == 17 {
			c.handlePairingRandom(pdu.Buf()[1:])
			return false
		}
	}

	c.sendPairingFailed(failedCommandNotSupported)

	return false
}
