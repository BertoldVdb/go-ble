package blesmp

import (
	"encoding/binary"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

// LE Secure Connections crypto toolbox (Bluetooth Core Spec v5.x,
// Vol 3, Part H, §2.2.5–2.2.9). All byte-strings here are in the
// big-endian / MSB-first order specified by the spec; callers must
// reverse incoming little-endian PDU fields before calling these
// functions and reverse outgoing results before transmitting.

// smpFuncF4 — Confirm value generation.
//
//	f4(U, V, X, Z) = AES-CMAC_X (U || V || Z)
//
// U, V are 32-byte X-coordinates of public keys; X is the 16-byte
// nonce/key; Z is a single byte ("0" for JustWorks/NumericComparison,
// 0x80 | passkey-bit for Passkey Entry).
func smpFuncF4(U, V [32]byte, X [16]byte, Z byte) [16]byte {
	var msg [32 + 32 + 1]byte
	copy(msg[:32], U[:])
	copy(msg[32:64], V[:])
	msg[64] = Z
	return aesCMAC(X[:], msg[:])
}

// smpFuncF5 — MacKey + LTK derivation.
//
//	T = AES-CMAC_SALT (W)
//	MacKey = AES-CMAC_T (Counter=0 || keyID || N1 || N2 || A1 || A2 || Length)
//	LTK    = AES-CMAC_T (Counter=1 || keyID || N1 || N2 || A1 || A2 || Length)
//
// SALT is the constant 0x6C888391AAF5A53860370BDB5A6083BE.
// keyID is the ASCII "btle" (0x62 0x74 0x6c 0x65).
// Length is 0x0100 (256 bits requested).
//
// W is the 32-byte ECDH shared secret (X coordinate).
func smpFuncF5(W [32]byte, N1, N2 [16]byte, A1, A2 [7]byte) (mac [16]byte, ltk [16]byte) {
	salt := [16]byte{0x6C, 0x88, 0x83, 0x91, 0xAA, 0xF5, 0xA5, 0x38, 0x60, 0x37, 0x0B, 0xDB, 0x5A, 0x60, 0x83, 0xBE}
	T := aesCMAC(salt[:], W[:])

	var msg [1 + 4 + 16 + 16 + 7 + 7 + 2]byte
	msg[1] = 'b'
	msg[2] = 't'
	msg[3] = 'l'
	msg[4] = 'e'
	copy(msg[5:21], N1[:])
	copy(msg[21:37], N2[:])
	copy(msg[37:44], A1[:])
	copy(msg[44:51], A2[:])
	binary.BigEndian.PutUint16(msg[51:53], 256)

	msg[0] = 0
	mac = aesCMAC(T[:], msg[:])

	msg[0] = 1
	ltk = aesCMAC(T[:], msg[:])
	return mac, ltk
}

// smpFuncF6 — DHKey check value.
//
//	f6(W, N1, N2, R, IOcap, A1, A2) = AES-CMAC_W (N1 || N2 || R || IOcap || A1 || A2)
//
// W is the 16-byte MacKey produced by f5. R is 16 bytes of association
// data (zero for JustWorks/NC, the passkey r-value for Passkey Entry,
// the OOB value for OOB).
func smpFuncF6(W [16]byte, N1, N2, R [16]byte, IOcap [3]byte, A1, A2 [7]byte) [16]byte {
	var msg [16 + 16 + 16 + 3 + 7 + 7]byte
	copy(msg[:16], N1[:])
	copy(msg[16:32], N2[:])
	copy(msg[32:48], R[:])
	copy(msg[48:51], IOcap[:])
	copy(msg[51:58], A1[:])
	copy(msg[58:65], A2[:])
	return aesCMAC(W[:], msg[:])
}

// smpFuncG2 — 6-digit numeric comparison value.
//
//	g2(U, V, X, Y) = AES-CMAC_X (U || V || Y) mod 2^32
//
// The result is taken mod 1,000,000 to obtain the 6-digit value the
// user is asked to compare across both devices.
func smpFuncG2(U, V [32]byte, X, Y [16]byte) uint32 {
	var msg [32 + 32 + 16]byte
	copy(msg[:32], U[:])
	copy(msg[32:64], V[:])
	copy(msg[64:80], Y[:])
	out := aesCMAC(X[:], msg[:])
	// mod 2^32 → low-order 32 bits of the 128-bit big-endian output (right-most 4 bytes)
	val := binary.BigEndian.Uint32(out[12:16])
	return val % 1000000
}

// smpAddrToA encodes a BLE address as the 7-byte A1/A2 form used by
// f5 / f6:  byte 0 = address-type bit (0 = public, 1 = random),
// bytes 1..6 = BD_ADDR with byte 1 = the most-significant address byte.
func smpAddrToA(addr bleutil.BLEAddr) [7]byte {
	var out [7]byte
	if addr.MacAddrType&1 == 1 {
		out[0] = 1
	}
	// Encode produces little-endian; reverse to get the spec's big-endian.
	var le [6]byte
	addr.MacAddr.Encode(le[:])
	for i := 0; i < 6; i++ {
		out[1+i] = le[5-i]
	}
	return out
}

// smpIOCapBytes packs the AuthReq, OOB and IOCapability fields into the
// 3-byte IOcap parameter accepted by f6 (Core Spec Vol 3 Part H §2.2.8).
//
//	IOcap[0] = AuthReq
//	IOcap[1] = OOB
//	IOcap[2] = IO capability
func smpIOCapBytes(authReq byte, oob byte, ioCap byte) [3]byte {
	return [3]byte{authReq, oob, ioCap}
}
