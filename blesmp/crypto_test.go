package blesmp

import (
	"io"
	"testing"

	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

func newTestSMPLogger() *logrus.Entry {
	l := logrus.New()
	l.Out = io.Discard
	l.Level = logrus.PanicLevel
	return logrus.NewEntry(l)
}

func TestCryptoShortenKey(t *testing.T) {
	in := [16]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}
	got := CryptoShortenKey(in, 7)
	for i := 0; i < 7; i++ {
		if got[i] != in[i] {
			t.Errorf("byte %d: got %#x want %#x", i, got[i], in[i])
		}
	}
	for i := 7; i < 16; i++ {
		if got[i] != 0 {
			t.Errorf("byte %d should be zero, got %#x", i, got[i])
		}
	}

	got = CryptoShortenKey(in, 16)
	if got != in {
		t.Errorf("16-byte key should be unchanged, got %x", got)
	}

	got = CryptoShortenKey(in, 0)
	for i := 0; i < 16; i++ {
		if got[i] != 0 {
			t.Errorf("0-byte key: byte %d not zero: %#x", i, got[i])
		}
	}
}

// CryptoFuncC1/CryptoFuncS1 spec test vectors live in Vol 3 Part H, App. D.
// Tested via local round-trip: c1(TK, rand, pres, preq, ia, ra) on both
// sides with consistent role flags must agree.
func TestCryptoFuncC1RoundTrip(t *testing.T) {
	var tk [16]byte
	rand := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	pres := [7]byte{0x05, 0x00, 0x00, 0x09, 0x10, 0x00, 0x01}
	preq := [7]byte{0x07, 0x07, 0x10, 0x00, 0x00, 0x01, 0x01}

	// Same address pair, same role: c1 must be deterministic.
	c1 := CryptoFuncC1(true, tk, rand, pres, preq, exampleAddr(0xA0), exampleAddr(0xB0))
	c2 := CryptoFuncC1(true, tk, rand, pres, preq, exampleAddr(0xA0), exampleAddr(0xB0))
	if c1 != c2 {
		t.Fatalf("c1 not deterministic: %x vs %x", c1, c2)
	}

	// Central and peripheral with swapped IA/RA must agree (the role
	// flag in CryptoFuncC1 swaps IA/RA internally for non-centrals).
	cs := CryptoFuncC1(true, tk, rand, pres, preq, exampleAddr(0xA0), exampleAddr(0xB0))
	cp := CryptoFuncC1(false, tk, rand, pres, preq, exampleAddr(0xB0), exampleAddr(0xA0))
	if cs != cp {
		t.Fatalf("central/peripheral c1 disagree: %x vs %x", cs, cp)
	}
}

func exampleAddr(seed byte) bleutil.BLEAddr {
	return bleutil.BLEAddr{
		MacAddr:     bleutil.MacAddr(seed) | bleutil.MacAddr(seed)<<8 | bleutil.MacAddr(seed)<<16,
		MacAddrType: 0,
	}
}

// ReversedAESCipher round-trip: encrypting then decrypting the same
// buffer must return the original bytes.
func TestReversedAESCipherRoundTrip(t *testing.T) {
	key := make([]byte, 16)
	for i := range key {
		key[i] = byte(i)
	}
	block := NewReversedAESCipher(key)

	plain := make([]byte, 16)
	for i := range plain {
		plain[i] = byte(0xFF - i)
	}
	cipher := make([]byte, 16)
	block.Encrypt(cipher, plain)

	out := make([]byte, 16)
	block.Decrypt(out, cipher)
	for i := range out {
		if out[i] != plain[i] {
			t.Fatalf("byte %d: got %#x want %#x", i, out[i], plain[i])
		}
	}

	// In-place encrypt: should not crash and must round-trip.
	buf := make([]byte, 16)
	copy(buf, plain)
	block.Encrypt(buf, buf)
	block.Decrypt(buf, buf)
	for i := range buf {
		if buf[i] != plain[i] {
			t.Fatalf("in-place byte %d: got %#x want %#x", i, buf[i], plain[i])
		}
	}

	if block.BlockSize() != 16 {
		t.Errorf("block size: got %d want 16", block.BlockSize())
	}
}

// CryptoGeneratePassKey returns a uniformly random 6-digit value
// (range [0, 999999]). Calling it many times must always produce a
// value in that range — never out of bounds.
func TestCryptoGeneratePassKeyRange(t *testing.T) {
	cfg := &SMPConnConfig{StaticPasscode: -1}
	logger := newTestSMPLogger()
	conn := &SMPConn{config: cfg, logger: logger}
	for i := 0; i < 1000; i++ {
		k, err := conn.CryptoGeneratePassKey()
		if err != nil {
			t.Fatal(err)
		}
		if k > 999999 {
			t.Fatalf("out of range: %d", k)
		}
	}
}

// StaticPasscode short-circuits the random generator. The returned
// value must be exactly the configured value.
func TestCryptoGeneratePassKeyStatic(t *testing.T) {
	cfg := &SMPConnConfig{StaticPasscode: 123456}
	conn := &SMPConn{config: cfg, logger: newTestSMPLogger()}
	for i := 0; i < 5; i++ {
		k, err := conn.CryptoGeneratePassKey()
		if err != nil {
			t.Fatal(err)
		}
		if k != 123456 {
			t.Fatalf("static passcode: got %d", k)
		}
	}
}

// CryptoFuncS1: STK = AES_TK(rand1[0..7] || rand2[0..7]). Verify that
// it produces the same output regardless of which side computes it
// (modulo the role-dependent input ordering).
func TestCryptoFuncS1Symmetric(t *testing.T) {
	var tk [16]byte
	r1 := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	r2 := [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

	a := CryptoFuncS1(true, tk, r1, r2)
	b := CryptoFuncS1(false, tk, r2, r1)
	if a != b {
		t.Fatalf("S1 disagreement central/peripheral: %x vs %x", a, b)
	}
}
