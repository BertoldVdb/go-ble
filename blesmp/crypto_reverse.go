package blesmp

import (
	"crypto/aes"
	"crypto/cipher"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

type ReversedCipher struct {
	block cipher.Block
}

func NewReversedAESCipher(key []byte) cipher.Block {
	key = bleutil.CopySlice(key)
	bleutil.ReverseSlice(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("cryptoFuncE: Failed to create AES cipher: " + err.Error())
	}

	return &ReversedCipher{block: block}
}

func (r *ReversedCipher) BlockSize() int {
	return r.block.BlockSize()
}

func (r *ReversedCipher) Decrypt(dst []byte, src []byte) {
	if !bleutil.SameSlice(dst, src) {
		src = bleutil.CopySlice(src)
	}
	bleutil.ReverseSlice(src)
	r.block.Decrypt(dst, src)
	bleutil.ReverseSlice(dst)
}

func (r *ReversedCipher) Encrypt(dst []byte, src []byte) {
	if !bleutil.SameSlice(dst, src) {
		src = bleutil.CopySlice(src)
	}
	bleutil.ReverseSlice(src)
	r.block.Encrypt(dst, src)
	bleutil.ReverseSlice(dst)
}
