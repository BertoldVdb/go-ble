package blesmp

import (
	"testing"
)

func TestECDHRoundTrip(t *testing.T) {
	a, err := scGenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	b, err := scGenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	dhA, err := scComputeDHKey(a, b.pubX, b.pubY)
	if err != nil {
		t.Fatal(err)
	}
	dhB, err := scComputeDHKey(b, a.pubX, a.pubY)
	if err != nil {
		t.Fatal(err)
	}
	if dhA != dhB {
		t.Fatalf("ECDH disagreement: %x vs %x", dhA, dhB)
	}
}

// scComputeDHKey must reject the case where the peer's public key is
// identical to ours (CVE-2018-5383 surface; debug-key swap). With this
// check in place an attacker that echoes our PK cannot derive the
// session DHKey.
func TestECDHRejectsOwnKey(t *testing.T) {
	kp, err := scGenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := scComputeDHKey(kp, kp.pubX, kp.pubY); err == nil {
		t.Fatal("expected error when peer public key equals our own")
	}
}

// scComputeDHKey must reject points that aren't on the curve. crypto/ecdh
// itself enforces this; we just confirm the error path returns.
func TestECDHRejectsBogusKey(t *testing.T) {
	kp, err := scGenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}
	var bogus [32]byte // all zeroes — not a valid P-256 X coordinate
	if _, err := scComputeDHKey(kp, bogus, bogus); err == nil {
		t.Fatal("expected error on all-zero peer point")
	}
}
