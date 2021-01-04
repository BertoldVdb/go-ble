package blesmp

import (
	bleutil "github.com/BertoldVdb/go-ble/util"
)

func CryptoFuncC1(tk [16]byte, rand [16]byte, pres [7]byte, preq [7]byte, ia bleutil.BLEAddr, ra bleutil.BLEAddr) [16]byte {
	block := NewReversedAESCipher(tk[:])

	/* Make p1 */
	p1 := [16]byte{}

	p1[0] = byte(ia.MacAddrType & 1)
	p1[1] = byte(ra.MacAddrType & 1)
	copy(p1[2:], pres[:])
	copy(p1[2+7:], preq[:])

	/* Xor p1 with r */
	for i := range rand {
		p1[i] ^= rand[i]
	}

	/*Calculate e(p1 xor r) */
	block.Encrypt(p1[:], p1[:])

	/* Make p2 */
	p2 := [16]byte{}
	ra.MacAddr.Encode(p2[:])
	ia.MacAddr.Encode(p2[6:])

	/* Xor e(p1 xor r) with p2 */
	for i := range rand {
		p1[i] ^= p2[i]
	}

	/* Calculate e(e(p1 xor r) xor p2) */
	block.Encrypt(p1[:], p1[:])

	var result [16]byte
	copy(result[:], p1[:])
	return result
}

func CryptoFuncS1(tk [16]byte, rand1 [16]byte, rand2 [16]byte) [16]byte {
	block := NewReversedAESCipher(tk[:])

	var r [16]byte
	copy(r[:], rand1[:8])
	copy(r[8:], rand2[:8])

	block.Encrypt(r[:], r[:])

	return r
}

func CryptoShortenKey(in [16]byte, l int) [16]byte {
	for i := l; i < len(in); i++ {
		in[i] = 0
	}

	return in
}
