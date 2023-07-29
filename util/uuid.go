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

func UUIDFromBytesValid(value []byte) (UUID, bool) {
	var result UUID

	l := len(value)
	if !(l == 2 || l == 4 || l == 16) {
		return result, false
	}

	if l == 16 {
		copy(result[:], value)
		return result, true
	}

	copy(result[:12], UUIDBase[:12])
	copy(result[12:], value)
	return result, true
}

func UUIDFromBytes(value []byte) UUID {
	uuid, valid := UUIDFromBytesValid(value)

	Assert(valid, "UUID needs to be 2, 4 or 16 bytes long")

	return uuid
}

func (u UUID) IsZero() bool {
	for _, v := range u {
		if v != 0 {
			return false
		}
	}
	return true
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

func (u *UUID) UUIDToBytes() []byte {
	length := u.GetLength()

	if length != 2 {
		return u[:]
	}

	return u[12:14]
}

func (u UUID) String() string {
	var sizes = []int{4, 2, 2, 2, 6}
	var blocks [5]string
	var index int
	skip := 0
	length := u.GetLength()

	if length <= 4 {
		sizes = sizes[:1]
	}
	if length <= 2 {
		skip = 4
	}

	ReverseSlice(u[:])

	for i, m := range sizes {
		blocks[i] = hex.EncodeToString(u[index : index+m])
		index += m
	}

	return strings.Join(blocks[:len(sizes)], "-")[skip:]
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

func (base UUID) CreateVariant(key uint8) UUID {
	if base.GetLength() != 16 {
		panic("Variant creation is only possible for random UUID")
	}

	base[0] += key

	return base
}

func (base UUID) CreateVariantAlt(key uint8) UUID {
	if base.GetLength() != 16 {
		panic("Variant creation is only possible for random UUID")
	}

	base[12] += key

	return base
}
