package bleutil

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrorMacAddrMalformed = errors.New("Mac address is malformed")
)

type MacAddr uint64

func (m MacAddr) Encode(data []byte) {
	data[0] = byte(m >> 0)
	data[1] = byte(m >> 8)
	data[2] = byte(m >> 16)
	data[3] = byte(m >> 24)
	data[4] = byte(m >> 32)
	data[5] = byte(m >> 40)
}

func (m *MacAddr) Decode(data []byte) {
	result := uint64(0)
	result |= uint64(data[0]) << 0
	result |= uint64(data[1]) << 8
	result |= uint64(data[2]) << 16
	result |= uint64(data[3]) << 24
	result |= uint64(data[4]) << 32
	result |= uint64(data[5]) << 40
	*m = MacAddr(result)
}

func (m MacAddr) String() string {
	var bytes [6]byte
	m.Encode(bytes[:])

	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", bytes[5], bytes[4], bytes[3], bytes[2], bytes[1], bytes[0])
}

func (m MacAddr) Network() string {
	return "Bluetooth"
}

func MacAddrFromString(mac string) (MacAddr, error) {
	mac = strings.ReplaceAll(mac, ":", "")
	mac = strings.ReplaceAll(mac, "/", "")
	mac = strings.ReplaceAll(mac, "-", "")

	var result MacAddr

	bytes, err := hex.DecodeString(mac)
	if err != nil {
		return result, err
	}
	if len(bytes) != 6 {
		return result, ErrorMacAddrMalformed
	}

	for i := 0; i < 3; i++ {
		bytes[i], bytes[5-i] = bytes[5-i], bytes[i]
	}

	result.Decode(bytes)

	return result, nil
}

func MacAddrFromStringPanic(mac string) MacAddr {
	addr, err := MacAddrFromString(mac)
	if err != nil {
		panic(err)
	}

	return addr
}

type MacAddrType int

const (
	MacAddrPublic         MacAddrType = 0
	MacAddrRandom         MacAddrType = 1
	MacAddrPublicIdentity MacAddrType = 2
	MacAddrStaticIdentity MacAddrType = 3
)

func (m MacAddrType) String() string {
	switch m {
	case MacAddrPublic:
		return "Public"
	case MacAddrRandom:
		return "Random"
	case MacAddrPublicIdentity:
		return "Public Identity"
	case MacAddrStaticIdentity:
		return "Static Identity"
	}

	return "Unknown"
}

type BLEAddr struct {
	MacAddr     MacAddr
	MacAddrType MacAddrType
}

func (m BLEAddr) GetUint64() uint64 {
	return uint64(m.MacAddr) | (uint64(m.MacAddrType) << 62)
}

func (m BLEAddr) String() string {
	return m.MacAddr.String() + " (" + m.MacAddrType.String() + ")"
}

func (a BLEAddr) IsLess(b BLEAddr) bool {
	if a.MacAddr < b.MacAddr {
		return true
	}
	if a.MacAddr == b.MacAddr {
		return a.MacAddrType < b.MacAddrType
	}
	return false
}

func (a BLEAddr) Network() string {
	return "BLE"
}
