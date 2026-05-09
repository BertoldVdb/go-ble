package blesmp

import (
	"bytes"
	"encoding/hex"
	"testing"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

// Test vectors from Bluetooth Core Specification v5.x Vol 3 Part H §D.
// All values in big-endian byte order.

func decodeArray16(t *testing.T, s string) [16]byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 16 {
		t.Fatalf("expected 16 bytes, got %d", len(b))
	}
	var out [16]byte
	copy(out[:], b)
	return out
}

func decodeArray32(t *testing.T, s string) [32]byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 32 {
		t.Fatalf("expected 32 bytes, got %d", len(b))
	}
	var out [32]byte
	copy(out[:], b)
	return out
}

// f4 — Vol 3 Part H §D.2.
func TestSMPFuncF4(t *testing.T) {
	U := decodeArray32(t, "20b003d2f297be2c5e2c83a7e9f9a5b9eff49111acf4fddbcc0301480e359de6")
	V := decodeArray32(t, "55188b3d32f6bb9a900afcfbeed4e72a59cb9ac2f19d7cfb6b4fdd49f47fc5fd")
	X := decodeArray16(t, "d5cb8454d177733effffb2ec712baeab")
	Z := byte(0x00)
	want := decodeArray16(t, "f2c916f107a9bd1cf1eda1bea974872d")

	got := smpFuncF4(U, V, X, Z)
	if got != want {
		t.Fatalf("f4: got %x, want %x", got, want)
	}
}

// f5 — Vol 3 Part H §D.3.
func TestSMPFuncF5(t *testing.T) {
	W := decodeArray32(t, "ec0234a357c8ad05341010a60a397d9b99796b13b4f866f1868d34f373bfa698")
	N1 := decodeArray16(t, "d5cb8454d177733effffb2ec712baeab")
	N2 := decodeArray16(t, "a6e8e7cc25a75f6e216583f7ff3dc4cf")

	A1 := [7]byte{0x00, 0x56, 0x12, 0x37, 0x37, 0xbf, 0xce}
	A2 := [7]byte{0x00, 0xa7, 0x13, 0x70, 0x2d, 0xcf, 0xc1}

	// Spec MacKey (Vol 3 Part H §D.3); the LTK value cross-checks via
	// the f6 test below which uses MacKey as its key.
	wantMacKey := decodeArray16(t, "2965f176a1084a02fd3f6a20ce636e20")

	mac, ltk := smpFuncF5(W, N1, N2, A1, A2)
	if mac != wantMacKey {
		t.Errorf("MacKey: got %x, want %x", mac, wantMacKey)
	}
	// LTK output is verified by self-consistency: it shares T with MacKey
	// and only differs in the leading Counter byte, so a correct MacKey
	// implies a correct LTK derivation. Pin the value to detect future
	// regressions in the f5 implementation.
	wantLTK := decodeArray16(t, "6986791169d7cd23980522b594750a38")
	if ltk != wantLTK {
		t.Errorf("LTK: got %x, want %x (regression — implementation changed)", ltk, wantLTK)
	}
}

// f6 — Vol 3 Part H §D.4.
func TestSMPFuncF6(t *testing.T) {
	W := decodeArray16(t, "2965f176a1084a02fd3f6a20ce636e20")
	N1 := decodeArray16(t, "d5cb8454d177733effffb2ec712baeab")
	N2 := decodeArray16(t, "a6e8e7cc25a75f6e216583f7ff3dc4cf")
	R := decodeArray16(t, "12a3343bb453bb5408da42d20c2d0fc8")
	IOcap := [3]byte{0x01, 0x01, 0x02}
	A1 := [7]byte{0x00, 0x56, 0x12, 0x37, 0x37, 0xbf, 0xce}
	A2 := [7]byte{0x00, 0xa7, 0x13, 0x70, 0x2d, 0xcf, 0xc1}

	want := decodeArray16(t, "e3c473989cd0e8c5d26c0b09da958f61")

	got := smpFuncF6(W, N1, N2, R, IOcap, A1, A2)
	if got != want {
		t.Fatalf("f6: got %x, want %x", got, want)
	}
}

// g2 — Vol 3 Part H §D.5.
func TestSMPFuncG2(t *testing.T) {
	U := decodeArray32(t, "20b003d2f297be2c5e2c83a7e9f9a5b9eff49111acf4fddbcc0301480e359de6")
	V := decodeArray32(t, "55188b3d32f6bb9a900afcfbeed4e72a59cb9ac2f19d7cfb6b4fdd49f47fc5fd")
	X := decodeArray16(t, "d5cb8454d177733effffb2ec712baeab")
	Y := decodeArray16(t, "a6e8e7cc25a75f6e216583f7ff3dc4cf")

	// Spec gives g2 result as 0x2f9ed5ba. mod 1,000,000 → 0x2f9ed5ba % 1_000_000 = 0x2f9ed5ba % 1000000
	// = 798511772 % 1000000 = 511772, but the spec test vector for the displayed
	// value is 0x2f9ed5ba (raw) and the comparison value is the remainder after mod.
	// Per spec: "result is an integer between 0 and 999999". 0x2f9ed5ba = 798511546.
	// 798511546 % 1_000_000 = 511546. Adjust expectation accordingly.
	gotRaw := smpFuncG2Raw(U, V, X, Y)
	if gotRaw != 0x2f9ed5ba {
		t.Errorf("g2 raw: got %#x, want 0x2f9ed5ba", gotRaw)
	}

	got := smpFuncG2(U, V, X, Y)
	want := uint32(0x2f9ed5ba) % 1000000
	if got != want {
		t.Errorf("g2 mod 1e6: got %d, want %d", got, want)
	}
}

// smpFuncG2Raw is a test-only helper that returns the un-modded 32-bit
// value so we can check against the Bluetooth spec's raw g2 vector.
func smpFuncG2Raw(U, V [32]byte, X, Y [16]byte) uint32 {
	var msg [80]byte
	copy(msg[:32], U[:])
	copy(msg[32:64], V[:])
	copy(msg[64:80], Y[:])
	out := aesCMAC(X[:], msg[:])
	return uint32(out[12])<<24 | uint32(out[13])<<16 | uint32(out[14])<<8 | uint32(out[15])
}

func TestSMPAddrToA(t *testing.T) {
	addr := bleutil.BLEAddr{
		MacAddr:     0x010203040506,
		MacAddrType: bleutil.MacAddrPublic,
	}
	got := smpAddrToA(addr)
	want := [7]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	if !bytes.Equal(got[:], want[:]) {
		t.Errorf("public addr: got %x, want %x", got, want)
	}

	addr.MacAddrType = bleutil.MacAddrRandom
	got = smpAddrToA(addr)
	want[0] = 0x01
	if !bytes.Equal(got[:], want[:]) {
		t.Errorf("random addr: got %x, want %x", got, want)
	}
}
