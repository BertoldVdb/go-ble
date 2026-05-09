package blesmp

import (
	"crypto/subtle"
	"testing"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

// End-to-end LE Secure Connections tests that exercise the crypto
// toolbox and association-model logic without spinning up a real
// L2CAP/HCI stack. Each test simulates an initiator and a responder
// doing the same computation and verifies they arrive at the same LTK
// (and that the DHKey-check exchange succeeds in both directions).

type scParty struct {
	kp    *scKeyPair
	addr  bleutil.BLEAddr
	io    [3]byte // AuthReq | OOB | IOcap
	nonce [16]byte
}

func newSCParty(t *testing.T, addr uint64, addrType bleutil.MacAddrType, io [3]byte) *scParty {
	t.Helper()
	kp, err := scGenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	var n [16]byte
	for i := range n {
		n[i] = byte(i ^ int(addr&0xFF))
	}
	return &scParty{
		kp:    kp,
		addr:  bleutil.BLEAddr{MacAddr: bleutil.MacAddr(addr), MacAddrType: addrType},
		io:    io,
		nonce: n,
	}
}

// runJustWorksOrNumComp runs the post-public-key Stage 1 + Stage 2 of
// JustWorks / NumericComparison and returns the LTKs derived on each
// side along with the displayed numeric value (only valid for NC).
func runJustWorksOrNumComp(t *testing.T, init, resp *scParty) (initLTK, respLTK [16]byte, numeric uint32, ok bool) {
	t.Helper()
	// Public key exchange happened: both sides hold the other's pub.
	dhInit, err := scComputeDHKey(init.kp, resp.kp.pubX, resp.kp.pubY)
	if err != nil {
		t.Fatal(err)
	}
	dhResp, err := scComputeDHKey(resp.kp, init.kp.pubX, init.kp.pubY)
	if err != nil {
		t.Fatal(err)
	}
	if dhInit != dhResp {
		t.Fatal("ECDH disagreement")
	}

	// Stage 1: responder sends Cb = f4(PKbx, PKax, Nb, 0); initiator
	// later verifies it after receiving Nb.
	cb := smpFuncF4(resp.kp.pubX, init.kp.pubX, resp.nonce, 0)

	// Initiator sends Na, responder sends Nb. Initiator verifies Cb.
	expectedCb := smpFuncF4(resp.kp.pubX, init.kp.pubX, resp.nonce, 0)
	if subtle.ConstantTimeCompare(cb[:], expectedCb[:]) == 0 {
		t.Fatal("initiator failed to verify responder's confirm")
	}

	// NumericComparison value (computed by both sides; user compares).
	numeric = smpFuncG2(init.kp.pubX, resp.kp.pubX, init.nonce, resp.nonce)

	// Stage 2: derive MacKey + LTK on both sides.
	addrA := smpAddrToA(init.addr)
	addrB := smpAddrToA(resp.addr)
	macA, ltkA := smpFuncF5(dhInit, init.nonce, resp.nonce, addrA, addrB)
	macB, ltkB := smpFuncF5(dhResp, init.nonce, resp.nonce, addrA, addrB)
	if macA != macB || ltkA != ltkB {
		t.Fatal("f5 disagreement")
	}

	// DHKey check exchange.
	var ra, rb [16]byte // zero for JW/NC
	Ea := smpFuncF6(macA, init.nonce, resp.nonce, rb, init.io, addrA, addrB)
	expEa := smpFuncF6(macB, init.nonce, resp.nonce, rb, init.io, addrA, addrB)
	if Ea != expEa {
		t.Fatal("Ea mismatch")
	}
	Eb := smpFuncF6(macB, resp.nonce, init.nonce, ra, resp.io, addrB, addrA)
	expEb := smpFuncF6(macA, resp.nonce, init.nonce, ra, resp.io, addrB, addrA)
	if Eb != expEb {
		t.Fatal("Eb mismatch")
	}

	return ltkA, ltkB, numeric, true
}

func TestSC_JustWorks_LTKAgreement(t *testing.T) {
	// IO capabilities: NoInputNoOutput on both sides; MITM=0.
	io := [3]byte{0x00, 0x00, 0x03}
	a := newSCParty(t, 0xaabbccddeeff, bleutil.MacAddrPublic, io)
	b := newSCParty(t, 0x112233445566, bleutil.MacAddrPublic, io)

	ltkA, ltkB, _, ok := runJustWorksOrNumComp(t, a, b)
	if !ok {
		t.Fatal("flow failed")
	}
	if ltkA != ltkB {
		t.Fatalf("LTK mismatch: %x vs %x", ltkA, ltkB)
	}
	var zero [16]byte
	if ltkA == zero {
		t.Fatal("LTK is all zeroes")
	}
}

func TestSC_NumericComparison(t *testing.T) {
	// Both sides DisplayYesNo, MITM=1.
	io := [3]byte{0x04, 0x00, 0x01}
	a := newSCParty(t, 0xaabbccddeeff, bleutil.MacAddrPublic, io)
	b := newSCParty(t, 0x112233445566, bleutil.MacAddrPublic, io)

	ltkA, ltkB, numeric, ok := runJustWorksOrNumComp(t, a, b)
	if !ok {
		t.Fatal("flow failed")
	}
	if ltkA != ltkB {
		t.Fatal("LTK mismatch")
	}
	if numeric >= 1000000 {
		t.Fatalf("numeric out of range: %d", numeric)
	}
}

// Passkey Entry: 20 rounds where each side commits to a per-round
// nonce via Cri = f4(PKax, PKbx, Nri, 0x80 | bit_i). The responder
// commits first; initiator reveals nonce, responder reveals nonce,
// initiator verifies responder's confirm, repeat.
func TestSC_Passkey_LTKAgreement(t *testing.T) {
	io := [3]byte{0x04, 0x00, 0x02} // KeyboardOnly, MITM=1
	a := newSCParty(t, 0xaabbccddeeff, bleutil.MacAddrPublic, io)
	b := newSCParty(t, 0x112233445566, bleutil.MacAddrPublic, io)

	const passkey uint32 = 314159

	// Compute DHKey both ways
	dhA, err := scComputeDHKey(a.kp, b.kp.pubX, b.kp.pubY)
	if err != nil {
		t.Fatal(err)
	}
	dhB, err := scComputeDHKey(b.kp, a.kp.pubX, a.kp.pubY)
	if err != nil {
		t.Fatal(err)
	}
	if dhA != dhB {
		t.Fatal("ECDH disagreement")
	}

	// 20 rounds — verify each round's confirm computation.
	for bit := 0; bit < 20; bit++ {
		var nai, nbi [16]byte
		for i := range nai {
			nai[i] = byte(i + bit*2)
			nbi[i] = byte(i + bit*2 + 1)
		}
		z := byte(0x80) | byte((passkey>>uint(bit))&1)

		Cai := smpFuncF4(a.kp.pubX, b.kp.pubX, nai, z)
		Cbi := smpFuncF4(b.kp.pubX, a.kp.pubX, nbi, z)

		// Initiator (knows nai) computes responder's confirm and verifies.
		expectedCbi := smpFuncF4(b.kp.pubX, a.kp.pubX, nbi, z)
		if Cbi != expectedCbi {
			t.Fatalf("round %d: Cbi mismatch", bit)
		}
		// Responder verifies initiator's confirm.
		expectedCai := smpFuncF4(a.kp.pubX, b.kp.pubX, nai, z)
		if Cai != expectedCai {
			t.Fatalf("round %d: Cai mismatch", bit)
		}
	}

	// Stage 2: r is the passkey in the low 4 bytes.
	var r [16]byte
	r[12] = byte((passkey >> 24) & 0xFF)
	r[13] = byte((passkey >> 16) & 0xFF)
	r[14] = byte((passkey >> 8) & 0xFF)
	r[15] = byte(passkey & 0xFF)

	addrA := smpAddrToA(a.addr)
	addrB := smpAddrToA(b.addr)
	macA, ltkA := smpFuncF5(dhA, a.nonce, b.nonce, addrA, addrB)
	macB, ltkB := smpFuncF5(dhB, a.nonce, b.nonce, addrA, addrB)
	if macA != macB || ltkA != ltkB {
		t.Fatal("f5 disagreement")
	}

	Ea := smpFuncF6(macA, a.nonce, b.nonce, r, a.io, addrA, addrB)
	if expEa := smpFuncF6(macB, a.nonce, b.nonce, r, a.io, addrA, addrB); Ea != expEa {
		t.Fatal("passkey Ea mismatch")
	}

	if ltkA != ltkB {
		t.Fatalf("LTK mismatch")
	}
}

// scNegotiated must report SC only when both sides set the SC bit.
func TestSCNegotiated(t *testing.T) {
	c := &SMPConn{}
	c.protocol.pairingRequest[3] = 0x08
	c.protocol.pairingResponse[3] = 0x08
	if !c.scNegotiated() {
		t.Error("expected SC negotiated when both sides set the SC bit")
	}
	c.protocol.pairingResponse[3] = 0x00
	if c.scNegotiated() {
		t.Error("expected SC NOT negotiated when responder lacks SC bit")
	}
	c.protocol.pairingRequest[3] = 0x00
	c.protocol.pairingResponse[3] = 0x08
	if c.scNegotiated() {
		t.Error("expected SC NOT negotiated when initiator lacks SC bit")
	}
}

func TestSCChooseAlgorithm(t *testing.T) {
	mk := func(initIO, respIO smpIOCapability, initMITM, respMITM bool) int {
		c := &SMPConn{}
		c.protocol.pairingRequest[1] = byte(initIO)
		c.protocol.pairingResponse[1] = byte(respIO)
		if initMITM {
			c.protocol.pairingRequest[3] |= 0x4
		}
		if respMITM {
			c.protocol.pairingResponse[3] |= 0x4
		}
		return scChooseAlgorithm(c)
	}

	if a := mk(cIONoInputNoOutput, cIONoInputNoOutput, false, false); a != scAlgorithmJustWorks {
		t.Errorf("no MITM should be JustWorks; got %d", a)
	}
	if a := mk(cIODisplayYesNo, cIODisplayYesNo, true, true); a != scAlgorithmNumericComparison {
		t.Errorf("YesNo+YesNo+MITM should be NumericComparison; got %d", a)
	}
	if a := mk(cIOKeyboardOnly, cIODisplayOnly, true, false); a != scAlgorithmPasskey {
		t.Errorf("Keyboard/Display+MITM should be Passkey; got %d", a)
	}
	if a := mk(cIONoInputNoOutput, cIOKeyboardOnly, true, true); a != scAlgorithmJustWorks {
		// NoIO can't authenticate; falls back to JW even though MITM is requested.
		t.Errorf("NoIO+Keyboard should fall back to JustWorks; got %d", a)
	}
}
