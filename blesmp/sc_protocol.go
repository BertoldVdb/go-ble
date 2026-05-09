package blesmp

import (
	crand "crypto/rand"
	"crypto/subtle"

	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

// LE Secure Connections protocol state machine.
//
// Flow (Vol 3, Part H §2.3.5):
//
//	Initiator                                  Responder
//	---------                                  ---------
//	Pairing Request   ----------------------->
//	                  <-----------------       Pairing Response
//	(both decide SC because the SC bit is set in both AuthReq fields)
//
//	Pairing Public Key (X || Y, LE)  -------->
//	                  <----------------------- Pairing Public Key
//	(both compute DHKey; reject if peer key is invalid or equals own)
//
//	Authentication stage 1 (varies per association method):
//	  JustWorks / NumericComparison:
//	     responder picks Nb, sends Cb = f4(PKbx, PKax, Nb, 0)
//	     initiator picks Na, sends it (no confirm from initiator)
//	     responder sends Nb
//	     initiator verifies Cb
//	     (NumericComparison only): both compute g2 and prompt user
//	  Passkey: 20 iterations of confirm-random with one bit per round.
//
//	Authentication stage 2 (DHKey check) — Vol 3 Part H §2.3.5.6.5:
//	  Both compute MacKey, LTK = f5(DHKey, Na, Nb, A, B)
//	  Initiator sends Ea = f6(MacKey, Na, Nb, rb, IOcap_a, A, B)
//	  Responder verifies Ea, sends Eb = f6(MacKey, Nb, Na, ra, IOcap_b, B, A)
//	  Initiator verifies Eb
//
//	Encryption with LTK; optional key distribution after that.

// scStartFromInitiator is called after the initiator has received a
// Pairing Response with the SC bit set. We generate our keypair and
// transmit our public key.
func (c *SMPConn) scStartFromInitiator() {
	if !c.scInitKeyPair() {
		return
	}
	c.scSendPublicKey()
	c.protocol.state = smpSCWaitPublicKey
}

// scStartFromResponder is called by the peripheral when a Pairing
// Request with the SC bit set has arrived and we've already sent our
// Pairing Response. We just wait for the initiator's public key.
func (c *SMPConn) scStartFromResponder() {
	if !c.scInitKeyPair() {
		return
	}
	c.protocol.state = smpSCWaitPublicKey
}

func (c *SMPConn) scInitKeyPair() bool {
	kp, err := scGenerateKeyPair()
	if err != nil {
		c.logger.WithError(err).Warn("LE SC keypair generation failed")
		c.sendPairingFailed(failedUnspecifiedReason)
		return false
	}
	c.protocol.scKeyPair = kp

	algorithm := scChooseAlgorithm(c)
	c.protocol.scAlgorithm = algorithm

	// In Passkey Entry, the input/display happens before the confirm
	// exchange. Other algorithms generate per-round nonces lazily.
	switch algorithm {
	case scAlgorithmPasskey:
		key, err := c.scAcquirePasskey()
		if err != nil {
			c.sendPairingFailed(failedPasskeyEntryFailed)
			return false
		}
		c.protocol.scPasskey = key
	}
	return true
}

// SC association methods (Core Spec Vol 3 Part H §2.3.5.6).
const (
	scAlgorithmJustWorks         = 0
	scAlgorithmPasskey           = 1
	scAlgorithmNumericComparison = 2
	scAlgorithmOOB               = 3 // not implemented
)

// scChooseAlgorithm picks the association method based on the IO
// capabilities (and MITM bits) advertised in the Pairing Request /
// Response. Mirrors Vol 3 Part H Table 2.8.
func scChooseAlgorithm(c *SMPConn) int {
	initIO := smpIOCapability(c.protocol.pairingRequest[1])
	respIO := smpIOCapability(c.protocol.pairingResponse[1])
	initMITM := c.protocol.pairingRequest[3]&0x4 != 0
	respMITM := c.protocol.pairingResponse[3]&0x4 != 0

	if !initMITM && !respMITM {
		return scAlgorithmJustWorks
	}

	hasKeyboard := func(cap smpIOCapability) bool {
		return cap == cIOKeyboardDisplay || cap == cIOKeyboardOnly
	}
	hasDisplay := func(cap smpIOCapability) bool {
		return cap == cIODisplayOnly || cap == cIODisplayYesNo || cap == cIOKeyboardDisplay
	}
	hasYesNo := func(cap smpIOCapability) bool {
		return cap == cIODisplayYesNo || cap == cIOKeyboardDisplay
	}

	// NumericComparison is preferred when both sides can display+yes/no
	// (Vol 3 Part H §2.3.5.6.4 row "DisplayYesNo / KeyboardDisplay /
	// DisplayYesNo / KeyboardDisplay"). Slightly broader than legacy
	// because the display-only side participates in DHKey-based
	// authentication.
	if hasDisplay(initIO) && hasYesNo(respIO) && hasYesNo(initIO) && hasDisplay(respIO) {
		return scAlgorithmNumericComparison
	}

	// Passkey: at least one side has a keyboard, and the other can
	// display or also has a keyboard.
	if hasKeyboard(initIO) && (hasDisplay(respIO) || hasKeyboard(respIO)) {
		return scAlgorithmPasskey
	}
	if hasKeyboard(respIO) && (hasDisplay(initIO) || hasKeyboard(initIO)) {
		return scAlgorithmPasskey
	}

	return scAlgorithmJustWorks
}

func (c *SMPConn) scAcquirePasskey() (uint32, error) {
	algorithm := c.protocol.scAlgorithm
	if algorithm != scAlgorithmPasskey {
		return 0, nil
	}

	myIO := c.getIOCapability()
	hasKeyboard := myIO == cIOKeyboardDisplay || myIO == cIOKeyboardOnly
	hasDisplay := myIO == cIODisplayOnly || myIO == cIODisplayYesNo || myIO == cIOKeyboardDisplay

	// Per spec: if both sides have a keyboard, the initiator displays
	// and the responder enters; if only one side has a keyboard, that
	// side enters. We approximate by: prefer to enter (input) if we have
	// a keyboard, otherwise display.
	if hasKeyboard && c.config.InputNumeric != nil {
		return c.config.InputNumeric(c)
	}
	if hasDisplay && c.config.DisplayNumeric != nil {
		key, err := c.CryptoGeneratePassKey()
		if err != nil {
			return 0, err
		}
		if err := c.config.DisplayNumeric(c, key); err != nil {
			return 0, err
		}
		return key, nil
	}
	return 0, errPeerKeyInvalid // fall back to refuse
}

// scSendPublicKey wraps our 64-byte public key in the PairingPublicKey
// PDU. The wire format is 32-byte X then 32-byte Y, both little-endian.
func (c *SMPConn) scSendPublicKey() {
	buf := bleutil.GetBuffer(1)
	buf.Buf()[0] = byte(opcodePairingPublicKey)
	body := buf.ExtendRight(64)
	copy(body[:32], c.protocol.scKeyPair.pubX[:])
	reverseInPlace(body[:32])
	copy(body[32:], c.protocol.scKeyPair.pubY[:])
	reverseInPlace(body[32:])
	c.sendBuf(buf)
}

// scHandlePublicKey processes the peer's Pairing Public Key PDU.
func (c *SMPConn) scHandlePublicKey(data []byte) {
	if len(data) != 64 {
		c.sendPairingFailed(failedInvalidParameters)
		return
	}

	// Wire bytes are little-endian; reverse to big-endian internally.
	copy(c.protocol.scPeerPubX[:], data[:32])
	reverseInPlace(c.protocol.scPeerPubX[:])
	copy(c.protocol.scPeerPubY[:], data[32:64])
	reverseInPlace(c.protocol.scPeerPubY[:])

	dh, err := scComputeDHKey(c.protocol.scKeyPair, c.protocol.scPeerPubX, c.protocol.scPeerPubY)
	if err != nil {
		c.logger.WithError(err).Warn("LE SC peer public key rejected")
		c.sendPairingFailed(failedDHKeyCheckFailed)
		return
	}
	c.protocol.scDHKey = dh

	if !c.isCentral {
		// Responder: send our own public key in reply.
		c.scSendPublicKey()
	}

	if _, err := crand.Read(c.protocol.scLocalNonce[:]); err != nil {
		c.sendPairingFailed(failedUnspecifiedReason)
		return
	}

	switch c.protocol.scAlgorithm {
	case scAlgorithmJustWorks, scAlgorithmNumericComparison:
		c.scStage1JustWorks()
	case scAlgorithmPasskey:
		c.protocol.scPasskeyBit = 0
		c.scStage1PasskeyRound()
	default:
		c.sendPairingFailed(failedAuthenticationRequirements)
	}
}

// scStage1JustWorks runs Stage 1 of the JustWorks / NumericComparison
// flow: responder commits to its nonce by sending Confirm; initiator
// just waits.
func (c *SMPConn) scStage1JustWorks() {
	if c.isCentral {
		// Initiator: wait for responder to send Confirm.
		c.protocol.state = smpSCWaitConfirm
		return
	}
	// Responder: compute Cb = f4(PKbx, PKax, Nb, 0) and send.
	confirm := smpFuncF4(c.protocol.scKeyPair.pubX, c.protocol.scPeerPubX, c.protocol.scLocalNonce, 0)
	c.protocol.scLocalConfirm = confirm
	c.scSendConfirm(confirm)
	c.protocol.state = smpSCWaitRandom
}

// scStage1PasskeyRound runs one of the 20 Passkey Entry rounds.
func (c *SMPConn) scStage1PasskeyRound() {
	bit := c.protocol.scPasskeyBit
	if bit >= 20 {
		// All rounds complete; advance to DHKey check.
		c.scStage2()
		return
	}
	if _, err := crand.Read(c.protocol.scLocalNonce[:]); err != nil {
		c.sendPairingFailed(failedUnspecifiedReason)
		return
	}

	z := byte(0x80) | byte((c.protocol.scPasskey>>uint(bit))&1)

	if c.isCentral {
		// Initiator: send Cai = f4(PKax, PKbx, Nai, Z).
		confirm := smpFuncF4(c.protocol.scKeyPair.pubX, c.protocol.scPeerPubX, c.protocol.scLocalNonce, z)
		c.protocol.scLocalConfirm = confirm
		c.scSendConfirm(confirm)
		c.protocol.state = smpSCWaitConfirm
	} else {
		// Responder: wait for initiator's Confirm before computing ours.
		c.protocol.state = smpSCWaitConfirm
	}
}

func (c *SMPConn) scSendConfirm(confirm [16]byte) {
	buf := bleutil.GetBuffer(1)
	buf.Buf()[0] = byte(opcodePairingConfirm)
	body := buf.ExtendRight(16)
	copy(body, confirm[:])
	reverseInPlace(body)
	c.sendBuf(buf)
}

func (c *SMPConn) scSendRandom() {
	buf := bleutil.GetBuffer(1)
	buf.Buf()[0] = byte(opcodePairingRandom)
	body := buf.ExtendRight(16)
	copy(body, c.protocol.scLocalNonce[:])
	reverseInPlace(body)
	c.sendBuf(buf)
}

// scHandleConfirm processes the peer's PairingConfirm PDU.
func (c *SMPConn) scHandleConfirm(data []byte) {
	if len(data) != 16 {
		c.sendPairingFailed(failedInvalidParameters)
		return
	}
	var peer [16]byte
	copy(peer[:], data)
	reverseInPlace(peer[:])
	c.protocol.scPeerConfirm = peer

	switch c.protocol.scAlgorithm {
	case scAlgorithmJustWorks, scAlgorithmNumericComparison:
		// Initiator received responder's confirm; now send our random.
		c.scSendRandom()
		c.protocol.state = smpSCWaitRandom
	case scAlgorithmPasskey:
		if c.isCentral {
			// Initiator received responder's confirm; now send our random.
			c.scSendRandom()
			c.protocol.state = smpSCWaitRandom
		} else {
			// Responder received initiator's confirm; compute and send ours.
			z := byte(0x80) | byte((c.protocol.scPasskey>>uint(c.protocol.scPasskeyBit))&1)
			confirm := smpFuncF4(c.protocol.scKeyPair.pubX, c.protocol.scPeerPubX, c.protocol.scLocalNonce, z)
			c.protocol.scLocalConfirm = confirm
			c.scSendConfirm(confirm)
			c.protocol.state = smpSCWaitRandom
		}
	default:
		c.sendPairingFailed(failedUnspecifiedReason)
	}
}

// scHandleRandom processes the peer's PairingRandom PDU and verifies
// the previously-received Confirm against it.
func (c *SMPConn) scHandleRandom(data []byte) {
	if len(data) != 16 {
		c.sendPairingFailed(failedInvalidParameters)
		return
	}
	var peerNonce [16]byte
	copy(peerNonce[:], data)
	reverseInPlace(peerNonce[:])
	c.protocol.scRemoteNonce = peerNonce

	switch c.protocol.scAlgorithm {
	case scAlgorithmJustWorks, scAlgorithmNumericComparison:
		c.scHandleRandomJW()
	case scAlgorithmPasskey:
		c.scHandleRandomPasskey()
	default:
		c.sendPairingFailed(failedUnspecifiedReason)
	}
}

func (c *SMPConn) scHandleRandomJW() {
	if c.isCentral {
		// Initiator: verify Cb from Nb.
		expected := smpFuncF4(c.protocol.scPeerPubX, c.protocol.scKeyPair.pubX, c.protocol.scRemoteNonce, 0)
		if subtle.ConstantTimeCompare(expected[:], c.protocol.scPeerConfirm[:]) == 0 {
			c.sendPairingFailed(failedConfirmValueFailed)
			return
		}
	} else {
		// Responder: send our nonce so the initiator can verify our confirm.
		c.scSendRandom()
	}

	if c.protocol.scAlgorithm == scAlgorithmNumericComparison {
		// Compute g2 and prompt user. If the user accepts, advance to
		// stage 2; otherwise fail. Lookup the IO callback in the SMP
		// config — InputYesNo represents "user can confirm match".
		var Ux, Vx [32]byte
		var Na, Nb [16]byte
		if c.isCentral {
			Ux = c.protocol.scKeyPair.pubX
			Vx = c.protocol.scPeerPubX
			Na = c.protocol.scLocalNonce
			Nb = c.protocol.scRemoteNonce
		} else {
			Ux = c.protocol.scPeerPubX
			Vx = c.protocol.scKeyPair.pubX
			Na = c.protocol.scRemoteNonce
			Nb = c.protocol.scLocalNonce
		}
		val := smpFuncG2(Ux, Vx, Na, Nb)
		c.protocol.scNumericValue = val
		c.logger.WithField("0value", val).Info("LE SC NumericComparison value")

		// Show the value and ask the user (synchronous; blocks the SMP goroutine).
		if c.config.DisplayNumeric != nil {
			if err := c.config.DisplayNumeric(c, val); err != nil {
				c.sendPairingFailed(failedNumericComparisonFailed)
				return
			}
		}
		ok := true
		if c.config.InputYesNo != nil {
			yn, err := c.config.InputYesNo(c)
			if err != nil || !yn {
				ok = false
			}
		}
		if !ok {
			c.sendPairingFailed(failedNumericComparisonFailed)
			return
		}
	}

	c.scStage2()
}

func (c *SMPConn) scHandleRandomPasskey() {
	// First, the initiator (after receiving responder's random) verifies
	// the responder's Confirm. The responder verifies the initiator's
	// Confirm right here.
	z := byte(0x80) | byte((c.protocol.scPasskey>>uint(c.protocol.scPasskeyBit))&1)
	var expected [16]byte
	if c.isCentral {
		// Verify Cbi = f4(PKbx, PKax, Nbi, Z).
		expected = smpFuncF4(c.protocol.scPeerPubX, c.protocol.scKeyPair.pubX, c.protocol.scRemoteNonce, z)
	} else {
		// Verify Cai = f4(PKax, PKbx, Nai, Z).
		expected = smpFuncF4(c.protocol.scPeerPubX, c.protocol.scKeyPair.pubX, c.protocol.scRemoteNonce, z)
	}
	if subtle.ConstantTimeCompare(expected[:], c.protocol.scPeerConfirm[:]) == 0 {
		c.sendPairingFailed(failedConfirmValueFailed)
		return
	}

	if !c.isCentral {
		// Responder: send our nonce.
		c.scSendRandom()
	}

	c.protocol.scPasskeyBit++
	c.scStage1PasskeyRound()
}

// scStage2 runs Phase 2 of LE SC: derive MacKey + LTK, send/verify
// DHKey checks, then enable encryption with the new LTK.
func (c *SMPConn) scStage2() {
	addrA := smpAddrToA(c.peerInitAddr())
	addrB := smpAddrToA(c.peerRespAddr())

	var Na, Nb [16]byte
	if c.isCentral {
		Na = c.protocol.scLocalNonce
		Nb = c.protocol.scRemoteNonce
	} else {
		Na = c.protocol.scRemoteNonce
		Nb = c.protocol.scLocalNonce
	}

	mac, ltk := smpFuncF5(c.protocol.scDHKey, Na, Nb, addrA, addrB)
	c.protocol.scMacKey = mac
	c.protocol.scLTK = ltk

	if c.isCentral {
		ioA := smpIOCapBytes(c.protocol.pairingRequest[3], c.protocol.pairingRequest[2], c.protocol.pairingRequest[1])
		var rb [16]byte
		if c.protocol.scAlgorithm == scAlgorithmPasskey {
			rb[12] = byte(c.protocol.scPasskey >> 24)
			rb[13] = byte(c.protocol.scPasskey >> 16)
			rb[14] = byte(c.protocol.scPasskey >> 8)
			rb[15] = byte(c.protocol.scPasskey)
		}
		Ea := smpFuncF6(mac, Na, Nb, rb, ioA, addrA, addrB)
		c.scSendDHKeyCheck(Ea)
	}
	// Both roles: now wait for the peer's DHKey check before completing.
	c.protocol.state = smpSCWaitDHKeyCheck
}

func (c *SMPConn) scSendDHKeyCheck(value [16]byte) {
	buf := bleutil.GetBuffer(1)
	buf.Buf()[0] = byte(opcodePairingDHKeyCheck)
	body := buf.ExtendRight(16)
	copy(body, value[:])
	reverseInPlace(body)
	c.sendBuf(buf)
}

// scHandleDHKeyCheck processes the peer's PairingDHKeyCheck PDU.
func (c *SMPConn) scHandleDHKeyCheck(data []byte) {
	if len(data) != 16 {
		c.sendPairingFailed(failedInvalidParameters)
		return
	}
	var peer [16]byte
	copy(peer[:], data)
	reverseInPlace(peer[:])

	addrA := smpAddrToA(c.peerInitAddr())
	addrB := smpAddrToA(c.peerRespAddr())

	var Na, Nb [16]byte
	if c.isCentral {
		Na = c.protocol.scLocalNonce
		Nb = c.protocol.scRemoteNonce
	} else {
		Na = c.protocol.scRemoteNonce
		Nb = c.protocol.scLocalNonce
	}

	ioA := smpIOCapBytes(c.protocol.pairingRequest[3], c.protocol.pairingRequest[2], c.protocol.pairingRequest[1])
	ioB := smpIOCapBytes(c.protocol.pairingResponse[3], c.protocol.pairingResponse[2], c.protocol.pairingResponse[1])

	var ra, rb [16]byte
	if c.protocol.scAlgorithm == scAlgorithmPasskey {
		ra[12] = byte(c.protocol.scPasskey >> 24)
		ra[13] = byte(c.protocol.scPasskey >> 16)
		ra[14] = byte(c.protocol.scPasskey >> 8)
		ra[15] = byte(c.protocol.scPasskey)
		rb = ra
	}

	if c.isCentral {
		// Verify Eb = f6(MacKey, Nb, Na, ra, IOcap_b, B, A).
		expected := smpFuncF6(c.protocol.scMacKey, Nb, Na, ra, ioB, addrB, addrA)
		if subtle.ConstantTimeCompare(expected[:], peer[:]) == 0 {
			c.sendPairingFailed(failedDHKeyCheckFailed)
			return
		}
		// Encrypt with the new LTK.
		c.protocol.pairingLTK.LTK = c.protocol.scLTK
		c.protocol.pairingLTK.Authenticated = c.protocol.scAlgorithm != scAlgorithmJustWorks
		ltk := smpStoredLTK{
			LTK:           c.protocol.scLTK,
			Authenticated: c.protocol.pairingLTK.Authenticated,
		}
		var err error
		if c.testEncryptHook != nil {
			err = c.testEncryptHook(ltk)
		} else {
			err = c.leEncrypt(ltk)
		}
		if err != nil {
			c.logger.WithError(err).Warn("LE SC encryption failed")
			c.sendPairingFailed(failedUnspecifiedReason)
			return
		}
		c.scComplete()
	} else {
		// Verify Ea = f6(MacKey, Na, Nb, rb, IOcap_a, A, B).
		expected := smpFuncF6(c.protocol.scMacKey, Na, Nb, rb, ioA, addrA, addrB)
		if subtle.ConstantTimeCompare(expected[:], peer[:]) == 0 {
			c.sendPairingFailed(failedDHKeyCheckFailed)
			return
		}
		// Send Eb.
		Eb := smpFuncF6(c.protocol.scMacKey, Nb, Na, ra, ioB, addrB, addrA)
		c.scSendDHKeyCheck(Eb)
		// Encryption is initiated by the central — we just wait for the
		// LELongTermKeyRequestEvent, which we'll answer with the LTK.
		c.protocol.pairingLTK.LTK = c.protocol.scLTK
		c.protocol.pairingLTK.Authenticated = c.protocol.scAlgorithm != scAlgorithmJustWorks
		c.scComplete()
	}
}

func (c *SMPConn) scComplete() {
	wantBond := (c.protocol.pairingRequest[3]&1 != 0) && (c.protocol.pairingResponse[3]&1 != 0)
	c.protocol.pairingLTK.Bonded = wantBond
	c.protocol.pairingLTK.EDIV = 0
	c.protocol.pairingLTK.Rand = 0
	c.protocol.pairingLTKValid = true
	c.protocol.pairingEDIVValid = true
	c.protocol.pairingLTKComplete = true

	/* Stash the LTK in the in-memory store regardless of bonding. The
	   peripheral's LELongTermKeyRequest handler only consults
	   storedKeys, so without this entry a non-bonded SC session cannot
	   complete encryption — the peripheral would NegativeReply and the
	   central would see encryption fail. The on-disk persistence is
	   gated on bonding so non-bonded keys don't leak across restarts. */
	c.updateLTK()

	c.leSetKeyFlagsFromLTK(c.protocol.pairingLTK)

	c.protocol.state = smpKeyDistribution

	/* StateSecure must NOT be set until the controller confirms
	   encryption is up. The central calls leEncrypt() before this
	   function (which already waited for EncryptionChange and signaled
	   StateSecure via smpHandler), but the responder reaches this
	   point BEFORE the central has even sent LEEnableEncryption.
	   Setting StateSecure here would tell ATT's checkSecurity() the
	   link is encrypted while plaintext traffic is still possible —
	   a real security regression.

	   On both sides the smpHandler's encUpdateChan loop transitions to
	   StateSecure when EncryptionChange arrives. For the central the
	   transition has already happened (leEncryptWait drained the
	   signal). For the responder we leave the state as smpKeyDistribution
	   until the event arrives. */
	c.updateTimeout(false)

	c.logger.WithFields(logrus.Fields{
		"0algorithm":     c.protocol.scAlgorithm,
		"1authenticated": c.protocol.pairingLTK.Authenticated,
		"2bonded":        c.protocol.pairingLTK.Bonded,
		"3isCentral":     c.isCentral,
	}).Info("LE Secure Connections pairing complete")

	/* The central, having already verified encryption via leEncryptWait,
	   safely advances to Secure here. The responder waits for its own
	   EncryptionChange event. */
	if c.isCentral {
		c.setState(StateSecure)
	}
}

// scHandleMessage routes incoming SMP PDUs to the SC handlers when SC
// is the active protocol. Returns false if the PDU is not for the SC
// state machine; the caller falls through to the legacy dispatcher.
func (c *SMPConn) scHandleMessage(opcode smpOpcode, body []byte) bool {
	switch c.protocol.state {
	case smpSCWaitPublicKey:
		if opcode == opcodePairingPublicKey {
			c.scHandlePublicKey(body)
			return true
		}
	case smpSCWaitConfirm:
		if opcode == opcodePairingConfirm {
			c.scHandleConfirm(body)
			return true
		}
	case smpSCWaitRandom:
		if opcode == opcodePairingRandom {
			c.scHandleRandom(body)
			return true
		}
	case smpSCWaitDHKeyCheck:
		if opcode == opcodePairingDHKeyCheck {
			c.scHandleDHKeyCheck(body)
			return true
		}
	case smpKeyDistribution:
		// Same KD handling as legacy.
		return c.handleKeyDistribution(opcode, body)
	}
	return false
}

// peerInitAddr / peerRespAddr return the BLE address of the initiator
// or responder, picking from local/remote according to our role.
func (c *SMPConn) peerInitAddr() bleutil.BLEAddr {
	if c.isCentral {
		return c.addrLELocal
	}
	return c.addrLERemote
}

func (c *SMPConn) peerRespAddr() bleutil.BLEAddr {
	if c.isCentral {
		return c.addrLERemote
	}
	return c.addrLELocal
}
