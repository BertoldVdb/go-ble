package blesmp

import (
	"bytes"
	"crypto/ecdh"
	crand "crypto/rand"
	"errors"
)

var (
	errPeerKeyEqualsOurs = errors.New("peer LE SC public key equals our own (debug-key swap rejected)")
	errPeerKeyInvalid    = errors.New("peer LE SC public key failed validation")
	errDHKeyZero         = errors.New("ECDH shared secret is zero")
)

// scKeyPair holds an ephemeral P-256 keypair for one LE Secure
// Connections pairing attempt. The private scalar is stored in the
// stdlib ecdh.PrivateKey; the public point is split into 32-byte
// big-endian X and Y coordinates so the wire encoding (which is
// little-endian) is a simple reversal.
type scKeyPair struct {
	priv *ecdh.PrivateKey
	pubX [32]byte
	pubY [32]byte
}

// scGenerateKeyPair returns a fresh P-256 ephemeral keypair sourced
// from crypto/rand. crypto/ecdh validates that the generated public
// key is on the curve.
func scGenerateKeyPair() (*scKeyPair, error) {
	priv, err := ecdh.P256().GenerateKey(crand.Reader)
	if err != nil {
		return nil, err
	}
	pub := priv.PublicKey().Bytes()
	// Uncompressed SEC1 point: 0x04 || X (32) || Y (32) = 65 bytes
	if len(pub) != 65 || pub[0] != 0x04 {
		return nil, errors.New("ecdh: unexpected public key encoding")
	}
	kp := &scKeyPair{priv: priv}
	copy(kp.pubX[:], pub[1:33])
	copy(kp.pubY[:], pub[33:65])
	return kp, nil
}

// scComputeDHKey derives the ECDH shared secret for the given peer
// public key. Returns the 32-byte X coordinate of the shared point in
// big-endian. Performs the LE SC validation requirements:
//   - peer point must be on the P-256 curve (delegated to crypto/ecdh
//     NewPublicKey, which rejects identity and off-curve points;
//     covers CVE-2018-5383 "Invalid Curve Attack");
//   - peer point must NOT equal our own public key (the "debug key swap"
//     scenario where an attacker echoes our public key);
//   - DHKey must not be zero (defense in depth — already implied by
//     the on-curve check, since the only way to get DHKey=0 is for the
//     peer point to be the identity).
func scComputeDHKey(kp *scKeyPair, peerX, peerY [32]byte) ([32]byte, error) {
	var zero [32]byte

	if bytes.Equal(peerX[:], kp.pubX[:]) && bytes.Equal(peerY[:], kp.pubY[:]) {
		return zero, errPeerKeyEqualsOurs
	}

	pubBytes := make([]byte, 65)
	pubBytes[0] = 0x04
	copy(pubBytes[1:33], peerX[:])
	copy(pubBytes[33:65], peerY[:])

	pub, err := ecdh.P256().NewPublicKey(pubBytes)
	if err != nil {
		return zero, errPeerKeyInvalid
	}

	secret, err := kp.priv.ECDH(pub)
	if err != nil {
		return zero, errPeerKeyInvalid
	}
	if len(secret) != 32 {
		return zero, errPeerKeyInvalid
	}

	var dhkey [32]byte
	copy(dhkey[:], secret)
	if dhkey == zero {
		return zero, errDHKeyZero
	}
	return dhkey, nil
}
