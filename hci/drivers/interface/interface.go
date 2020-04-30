package hciinterface

import (
	"errors"
	"time"
)

var (
	// ErrorDeviceNotFound is returned if the specified device cannot be found.
	ErrorDeviceNotFound = errors.New("Device not found")
)

// HCIRxPacket contains the HCI RX packet data
type HCIRxPacket struct {
	// Received is true for packets received from the BLE controller. Only in sniffing mode can it be false.
	Received bool

	// TimeFromHardware indicates if hardware timestamping was used.
	TimeFromHardware bool

	// RxTime contains an RX timestamp
	RxTime time.Time

	// Data contains the packet data
	Data []byte
}

// HCITxPacket contains the HCI TX packet data
type HCITxPacket struct {
	// Data contains the packet data
	Data []byte
}

// HCIRxHandler is the type of the receive callback function. It gets a
// HCIRxPacket as sole parameter, and potentially returns an error.
type HCIRxHandler func(pkt HCIRxPacket) error

// HCIInterface represents a generic HCI interface that can send and receive HCI
// packets to the hardware.
type HCIInterface interface {
	// Run is the worker function. It needs to be running before packets can be sent or received.
	Run() error

	// Close closes the the interface. It can be called at any time and multiple times as well.
	// It will terminate Run, if it was running.
	Close() error

	// SetRecvHandler configures the receive handler callback function. It will be called
	// when a HCI packet is received.
	SetRecvHandler(HCIRxHandler) error

	// SendPacket sends a HCI packet to the device.
	SendPacket(pkt HCITxPacket) error
}
