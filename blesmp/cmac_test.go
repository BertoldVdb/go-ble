package blesmp

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// AES-CMAC test vectors from RFC 4493, Appendix B.
//
// Key K = 2b7e1516 28aed2a6 abf71588 09cf4f3c
//
// Length 0:    bb1d6929 e9593728 7fa37d12 9b756746
// Length 16:   070a16b4 6b4d4144 f79bdd9d d04a287c
// Length 40:   dfa66747 de9ae630 30ca3261 1497c827
// Length 64:   51f0bebf 7e3b9d92 fc497417 79363cfe
func TestAESCMACRFC4493Vectors(t *testing.T) {
	K := mustHex(t, "2b7e151628aed2a6abf7158809cf4f3c")
	M := mustHex(t,
		"6bc1bee22e409f96e93d7e117393172a"+
			"ae2d8a571e03ac9c9eb76fac45af8e51"+
			"30c81c46a35ce411e5fbc1191a0a52ef"+
			"f69f2445df4f9b17ad2b417be66c3710")

	cases := []struct {
		name string
		msg  []byte
		want string
	}{
		{"len=0", nil, "bb1d6929e95937287fa37d129b756746"},
		{"len=16", M[:16], "070a16b46b4d4144f79bdd9dd04a287c"},
		{"len=40", M[:40], "dfa66747de9ae63030ca32611497c827"},
		{"len=64", M[:64], "51f0bebf7e3b9d92fc49741779363cfe"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := aesCMAC(K, c.msg)
			want := mustHex(t, c.want)
			if !bytes.Equal(got[:], want) {
				t.Fatalf("got %x\nwant %x", got, want)
			}
		})
	}
}

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("decode %q: %v", s, err)
	}
	return b
}

func TestReverseHelpers(t *testing.T) {
	in := []byte{0x01, 0x02, 0x03, 0x04}
	out := reversed(in)
	if !bytes.Equal(out, []byte{0x04, 0x03, 0x02, 0x01}) {
		t.Errorf("reversed: got %x", out)
	}
	if !bytes.Equal(in, []byte{0x01, 0x02, 0x03, 0x04}) {
		t.Errorf("reversed mutated input: %x", in)
	}

	cp := append([]byte(nil), in...)
	reverseInPlace(cp)
	if !bytes.Equal(cp, []byte{0x04, 0x03, 0x02, 0x01}) {
		t.Errorf("reverseInPlace: got %x", cp)
	}

	// Odd length
	cp = []byte{1, 2, 3, 4, 5}
	reverseInPlace(cp)
	if !bytes.Equal(cp, []byte{5, 4, 3, 2, 1}) {
		t.Errorf("odd reverseInPlace: got %x", cp)
	}

	// Empty / single-byte: must not panic
	reverseInPlace(nil)
	reverseInPlace([]byte{0x42})
}
