package blesmp

import (
	"crypto/aes"
	"crypto/cipher"
)

// aesCMAC implements AES-CMAC (RFC 4493) for arbitrary-length messages.
// The key must be 16 bytes (AES-128). All inputs and outputs are in
// the natural byte order required by the standard — the BLE LE SC
// crypto toolbox specifies byte ordering up-front (Vol 3, Part H §2.2.5)
// so we hand the function plain MSB-first inputs and reverse only when
// we cross the wire.
func aesCMAC(key, msg []byte) [16]byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		// AES-128 key must be 16 bytes; this is a programming error.
		panic("blesmp: aesCMAC requires a 16-byte key")
	}
	return aesCMACBlock(block, msg)
}

func aesCMACBlock(block cipher.Block, msg []byte) [16]byte {
	const bs = 16
	const rb = 0x87

	var L [bs]byte
	block.Encrypt(L[:], L[:])

	K1 := cmacShiftLeft(L[:])
	if L[0]&0x80 != 0 {
		K1[bs-1] ^= rb
	}
	K2 := cmacShiftLeft(K1[:])
	if K1[0]&0x80 != 0 {
		K2[bs-1] ^= rb
	}

	n := (len(msg) + bs - 1) / bs
	complete := len(msg) > 0 && len(msg)%bs == 0
	if n == 0 {
		// Empty message: pad with 10..0 and use K2.
		var last [bs]byte
		last[0] = 0x80
		for j := 0; j < bs; j++ {
			last[j] ^= K2[j]
		}
		var out [bs]byte
		block.Encrypt(out[:], last[:])
		return out
	}

	var X, Y [bs]byte
	for i := 0; i < n-1; i++ {
		for j := 0; j < bs; j++ {
			Y[j] = X[j] ^ msg[i*bs+j]
		}
		block.Encrypt(X[:], Y[:])
	}

	var last [bs]byte
	if complete {
		copy(last[:], msg[(n-1)*bs:])
		for j := 0; j < bs; j++ {
			last[j] ^= K1[j]
		}
	} else {
		rem := len(msg) % bs
		copy(last[:], msg[(n-1)*bs:(n-1)*bs+rem])
		last[rem] = 0x80
		for j := 0; j < bs; j++ {
			last[j] ^= K2[j]
		}
	}

	for j := 0; j < bs; j++ {
		Y[j] = X[j] ^ last[j]
	}
	var out [bs]byte
	block.Encrypt(out[:], Y[:])
	return out
}

func cmacShiftLeft(in []byte) [16]byte {
	var out [16]byte
	var carry byte
	for i := len(in) - 1; i >= 0; i-- {
		out[i] = (in[i] << 1) | carry
		carry = (in[i] >> 7) & 0x1
	}
	return out
}

// reverseInPlace reverses a byte slice. BLE PDUs carry multi-byte values
// in little-endian wire order; the SC crypto toolbox is defined in
// big-endian. Helper for crossing that boundary.
func reverseInPlace(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

// reversed returns a fresh reversed copy of b.
func reversed(b []byte) []byte {
	out := make([]byte, len(b))
	for i := range b {
		out[i] = b[len(b)-1-i]
	}
	return out
}
