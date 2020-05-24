package bleutil

import (
	"bytes"
	"encoding/hex"
	"strings"
)

type UUID [16]byte

var (
	UUIDBase UUID = [16]byte{0xFB, 0x34, 0x9B, 0x5F, 0x80, 0x00, 0x00, 0x80, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)

func UUIDFromBytes(value []byte) UUID {
	l := len(value)
	Assert(l == 2 || l == 4 || l == 16, "Invalid length, needs to be 16, 32 or 128 bits")

	var result UUID

	if l == 16 {
		copy(result[:], value)
		return result
	}

	copy(result[:12], UUIDBase[:12])
	copy(result[12:], value)
	return result
}

func (u *UUID) GetLength() int {
	if !bytes.Equal(u[:12], UUIDBase[:12]) {
		return 16
	}

	if u[14] == 0 && u[15] == 0 {
		return 2
	}

	return 4
}

func (u UUID) String() string {
	var sizes = []int{4, 2, 2, 2, 6}
	var blocks [5]string
	var index int

	ReverseSlice(u[:])

	for i, m := range sizes {
		blocks[i] = hex.EncodeToString(u[index : index+m])
		index += m
	}

	return strings.Join(blocks[:], "-")
}

func UUIDFromString(value string) (UUID, error) {
	value = strings.ReplaceAll(value, "-", "")
	value = strings.ReplaceAll(value, "_", "")

	bytes, err := hex.DecodeString(value)
	if err != nil {
		return UUIDBase, err
	}

	ReverseSlice(bytes)

	return UUIDFromBytes(bytes), nil
}

func UUIDFromStringPanic(value string) UUID {
	result, err := UUIDFromString(value)
	if err != nil {
		panic(err)
	}
	return result
}
