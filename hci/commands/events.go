package hcicommands

import (
	"encoding/binary"
	//hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
)

// InquiryResultEvent represents the event specified in Section 7.7.2
type InquiryResultEvent struct {
	NumResponses           uint8
	BDADDR                 [][6]byte
	PageScanRepetitionMode []uint8
	Reserved               []uint16
	ClassOfDevice          []uint32
	ClockOffset            []uint16
}

func (o *InquiryResultEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.NumResponses = r.GetOne()
	if cap(o.BDADDR) < int(o.NumResponses) {
		o.BDADDR = make([][6]byte, 0, int(o.NumResponses))
	}
	o.BDADDR = o.BDADDR[:int(o.NumResponses)]
	for j := 0; j < int(o.NumResponses); j++ {
		copy(o.BDADDR[j][:], r.Get(6))
	}
	if cap(o.PageScanRepetitionMode) < int(o.NumResponses) {
		o.PageScanRepetitionMode = make([]uint8, 0, int(o.NumResponses))
	}
	o.PageScanRepetitionMode = o.PageScanRepetitionMode[:int(o.NumResponses)]
	for j := 0; j < int(o.NumResponses); j++ {
		o.PageScanRepetitionMode[j] = r.GetOne()
	}
	if cap(o.Reserved) < int(o.NumResponses) {
		o.Reserved = make([]uint16, 0, int(o.NumResponses))
	}
	o.Reserved = o.Reserved[:int(o.NumResponses)]
	for j := 0; j < int(o.NumResponses); j++ {
		o.Reserved[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.ClassOfDevice) < int(o.NumResponses) {
		o.ClassOfDevice = make([]uint32, 0, int(o.NumResponses))
	}
	o.ClassOfDevice = o.ClassOfDevice[:int(o.NumResponses)]
	for j := 0; j < int(o.NumResponses); j++ {
		o.ClassOfDevice[j] = decodeUint24(r.Get(3))
	}
	if cap(o.ClockOffset) < int(o.NumResponses) {
		o.ClockOffset = make([]uint16, 0, int(o.NumResponses))
	}
	o.ClockOffset = o.ClockOffset[:int(o.NumResponses)]
	for j := 0; j < int(o.NumResponses); j++ {
		o.ClockOffset[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// ConnectionCompleteEvent represents the event specified in Section 7.7.3
type ConnectionCompleteEvent struct {
	Status            uint8
	ConnectionHandle  uint16
	BDADDR            [6]byte
	LinkType          uint8
	EncryptionEnabled uint8
}

func (o *ConnectionCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	copy(o.BDADDR[:], r.Get(6))
	o.LinkType = r.GetOne()
	o.EncryptionEnabled = r.GetOne()
	return r.Valid()
}

// ConnectionRequestEvent represents the event specified in Section 7.7.4
type ConnectionRequestEvent struct {
	BDADDR        [6]byte
	ClassOfDevice uint32
	LinkType      uint8
}

func (o *ConnectionRequestEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	o.ClassOfDevice = decodeUint24(r.Get(3))
	o.LinkType = r.GetOne()
	return r.Valid()
}

// DisconnectionCompleteEvent represents the event specified in Section 7.7.5
type DisconnectionCompleteEvent struct {
	Status           uint8
	ConnectionHandle uint16
	Reason           uint8
}

func (o *DisconnectionCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Reason = r.GetOne()
	return r.Valid()
}

// AuthenticationCompleteEvent represents the event specified in Section 7.7.6
type AuthenticationCompleteEvent struct {
	Status           uint8
	ConnectionHandle uint16
}

func (o *AuthenticationCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// RemoteNameRequestCompleteEvent represents the event specified in Section 7.7.7
type RemoteNameRequestCompleteEvent struct {
	Status     uint8
	BDADDR     [6]byte
	RemoteName [248]byte
}

func (o *RemoteNameRequestCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	copy(o.RemoteName[:], r.Get(248))
	return r.Valid()
}

// EncryptionChangeEvent represents the event specified in Section 7.7.8
type EncryptionChangeEvent struct {
	Status            uint8
	ConnectionHandle  uint16
	EncryptionEnabled uint8
}

func (o *EncryptionChangeEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.EncryptionEnabled = r.GetOne()
	return r.Valid()
}

// ChangeConnectionLinkKeyCompleteEvent represents the event specified in Section 7.7.9
type ChangeConnectionLinkKeyCompleteEvent struct {
	Status           uint8
	ConnectionHandle uint16
}

func (o *ChangeConnectionLinkKeyCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// MasterLinkKeyCompleteEvent represents the event specified in Section 7.7.10
type MasterLinkKeyCompleteEvent struct {
	Status           uint8
	ConnectionHandle uint16
	KeyFlag          uint8
}

func (o *MasterLinkKeyCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.KeyFlag = r.GetOne()
	return r.Valid()
}

// ReadRemoteSupportedFeaturesCompleteEvent represents the event specified in Section 7.7.11
type ReadRemoteSupportedFeaturesCompleteEvent struct {
	Status           uint8
	ConnectionHandle uint16
	LMPFeatures      uint64
}

func (o *ReadRemoteSupportedFeaturesCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LMPFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// ReadRemoteVersionInformationCompleteEvent represents the event specified in Section 7.7.12
type ReadRemoteVersionInformationCompleteEvent struct {
	Status           uint8
	ConnectionHandle uint16
	Version          uint8
	ManufacturerName uint16
	Subversion       uint16
}

func (o *ReadRemoteVersionInformationCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Version = r.GetOne()
	o.ManufacturerName = binary.LittleEndian.Uint16(r.Get(2))
	o.Subversion = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// QoSSetupCompleteEvent represents the event specified in Section 7.7.13
type QoSSetupCompleteEvent struct {
	Status           uint8
	ConnectionHandle uint16
	Unused           uint8
	ServiceType      uint8
	TokenRate        uint32
	PeakBandwidth    uint32
	Latency          uint32
	DelayVariation   uint32
}

func (o *QoSSetupCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Unused = r.GetOne()
	o.ServiceType = r.GetOne()
	o.TokenRate = binary.LittleEndian.Uint32(r.Get(4))
	o.PeakBandwidth = binary.LittleEndian.Uint32(r.Get(4))
	o.Latency = binary.LittleEndian.Uint32(r.Get(4))
	o.DelayVariation = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// CommandCompleteEvent represents the event specified in Section 7.7.14
type CommandCompleteEvent struct {
	NumHCICommandPackets uint8
	CommandOpcode        uint16
}

func (o *CommandCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.NumHCICommandPackets = r.GetOne()
	o.CommandOpcode = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// CommandStatusEvent represents the event specified in Section 7.7.15
type CommandStatusEvent struct {
	Status               uint8
	NumHCICommandPackets uint8
	CommandOpcode        uint16
}

func (o *CommandStatusEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.NumHCICommandPackets = r.GetOne()
	o.CommandOpcode = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// HardwareErrorEvent represents the event specified in Section 7.7.16
type HardwareErrorEvent struct {
	HardwareCode uint8
}

func (o *HardwareErrorEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.HardwareCode = r.GetOne()
	return r.Valid()
}

// FlushOccurredEvent represents the event specified in Section 7.7.17
type FlushOccurredEvent struct {
	Handle uint16
}

func (o *FlushOccurredEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// RoleChangeEvent represents the event specified in Section 7.7.18
type RoleChangeEvent struct {
	Status  uint8
	BDADDR  [6]byte
	NewRole uint8
}

func (o *RoleChangeEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	o.NewRole = r.GetOne()
	return r.Valid()
}

// NumberOfCompletedPacketsEvent represents the event specified in Section 7.7.19
type NumberOfCompletedPacketsEvent struct {
	NumHandles          uint8
	ConnectionHandle    []uint16
	NumCompletedPackets []uint16
}

func (o *NumberOfCompletedPacketsEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.NumHandles = r.GetOne()
	if cap(o.ConnectionHandle) < int(o.NumHandles) {
		o.ConnectionHandle = make([]uint16, 0, int(o.NumHandles))
	}
	o.ConnectionHandle = o.ConnectionHandle[:int(o.NumHandles)]
	for j := 0; j < int(o.NumHandles); j++ {
		o.ConnectionHandle[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.NumCompletedPackets) < int(o.NumHandles) {
		o.NumCompletedPackets = make([]uint16, 0, int(o.NumHandles))
	}
	o.NumCompletedPackets = o.NumCompletedPackets[:int(o.NumHandles)]
	for j := 0; j < int(o.NumHandles); j++ {
		o.NumCompletedPackets[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// ModeChangeEvent represents the event specified in Section 7.7.20
type ModeChangeEvent struct {
	Status           uint8
	ConnectionHandle uint16
	CurrentMode      uint8
	Interval         uint16
}

func (o *ModeChangeEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.CurrentMode = r.GetOne()
	o.Interval = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// ReturnLinkKeysEvent represents the event specified in Section 7.7.21
type ReturnLinkKeysEvent struct {
	NumKeys uint8
	BDADDR  [][6]byte
	LinkKey [][16]byte
}

func (o *ReturnLinkKeysEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.NumKeys = r.GetOne()
	if cap(o.BDADDR) < int(o.NumKeys) {
		o.BDADDR = make([][6]byte, 0, int(o.NumKeys))
	}
	o.BDADDR = o.BDADDR[:int(o.NumKeys)]
	for j := 0; j < int(o.NumKeys); j++ {
		copy(o.BDADDR[j][:], r.Get(6))
	}
	if cap(o.LinkKey) < int(o.NumKeys) {
		o.LinkKey = make([][16]byte, 0, int(o.NumKeys))
	}
	o.LinkKey = o.LinkKey[:int(o.NumKeys)]
	for j := 0; j < int(o.NumKeys); j++ {
		copy(o.LinkKey[j][:], r.Get(16))
	}
	return r.Valid()
}

// PINCodeRequestEvent represents the event specified in Section 7.7.22
type PINCodeRequestEvent struct {
	BDADDR [6]byte
}

func (o *PINCodeRequestEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkKeyRequestEvent represents the event specified in Section 7.7.23
type LinkKeyRequestEvent struct {
	BDADDR [6]byte
}

func (o *LinkKeyRequestEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkKeyNotificationEvent represents the event specified in Section 7.7.24
type LinkKeyNotificationEvent struct {
	BDADDR  [6]byte
	LinkKey [16]byte
	KeyType uint8
}

func (o *LinkKeyNotificationEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	copy(o.LinkKey[:], r.Get(16))
	o.KeyType = r.GetOne()
	return r.Valid()
}

// DataBufferOverflowEvent represents the event specified in Section 7.7.26
type DataBufferOverflowEvent struct {
	LinkType uint8
}

func (o *DataBufferOverflowEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.LinkType = r.GetOne()
	return r.Valid()
}

// MaxSlotsChangeEvent represents the event specified in Section 7.7.27
type MaxSlotsChangeEvent struct {
	ConnectionHandle uint16
	LMPMaxSlots      uint8
}

func (o *MaxSlotsChangeEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LMPMaxSlots = r.GetOne()
	return r.Valid()
}

// ReadClockOffsetCompleteEvent represents the event specified in Section 7.7.28
type ReadClockOffsetCompleteEvent struct {
	Status           uint8
	ConnectionHandle uint16
	ClockOffset      uint16
}

func (o *ReadClockOffsetCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ClockOffset = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// ConnectionPacketTypeChangedEvent represents the event specified in Section 7.7.29
type ConnectionPacketTypeChangedEvent struct {
	Status           uint8
	ConnectionHandle uint16
	PacketType       uint16
}

func (o *ConnectionPacketTypeChangedEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PacketType = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// QoSViolationEvent represents the event specified in Section 7.7.30
type QoSViolationEvent struct {
	Handle uint16
}

func (o *QoSViolationEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// PageScanRepetitionModeChangeEvent represents the event specified in Section 7.7.31
type PageScanRepetitionModeChangeEvent struct {
	BDADDR                 [6]byte
	PageScanRepetitionMode uint8
}

func (o *PageScanRepetitionModeChangeEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	o.PageScanRepetitionMode = r.GetOne()
	return r.Valid()
}

// FlowSpecificationCompleteEvent represents the event specified in Section 7.7.32
type FlowSpecificationCompleteEvent struct {
	Status           uint8
	ConnectionHandle uint16
	Unused           uint8
	FlowDirection    uint8
	ServiceType      uint8
	TokenRate        uint32
	TokenBucketSize  uint32
	PeakBandwidth    uint32
	AccessLatency    uint32
}

func (o *FlowSpecificationCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Unused = r.GetOne()
	o.FlowDirection = r.GetOne()
	o.ServiceType = r.GetOne()
	o.TokenRate = binary.LittleEndian.Uint32(r.Get(4))
	o.TokenBucketSize = binary.LittleEndian.Uint32(r.Get(4))
	o.PeakBandwidth = binary.LittleEndian.Uint32(r.Get(4))
	o.AccessLatency = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// InquiryResultwithRSSIEvent represents the event specified in Section 7.7.33
type InquiryResultwithRSSIEvent struct {
	NumResponses           uint8
	BDADDR                 [][6]byte
	PageScanRepetitionMode []uint8
	Reserved               []uint8
	ClassOfDevice          []uint32
	ClockOffset            []uint16
	RSSI                   []uint8
}

func (o *InquiryResultwithRSSIEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.NumResponses = r.GetOne()
	if cap(o.BDADDR) < int(o.NumResponses) {
		o.BDADDR = make([][6]byte, 0, int(o.NumResponses))
	}
	o.BDADDR = o.BDADDR[:int(o.NumResponses)]
	for j := 0; j < int(o.NumResponses); j++ {
		copy(o.BDADDR[j][:], r.Get(6))
	}
	if cap(o.PageScanRepetitionMode) < int(o.NumResponses) {
		o.PageScanRepetitionMode = make([]uint8, 0, int(o.NumResponses))
	}
	o.PageScanRepetitionMode = o.PageScanRepetitionMode[:int(o.NumResponses)]
	for j := 0; j < int(o.NumResponses); j++ {
		o.PageScanRepetitionMode[j] = r.GetOne()
	}
	if cap(o.Reserved) < int(o.NumResponses) {
		o.Reserved = make([]uint8, 0, int(o.NumResponses))
	}
	o.Reserved = o.Reserved[:int(o.NumResponses)]
	for j := 0; j < int(o.NumResponses); j++ {
		o.Reserved[j] = r.GetOne()
	}
	if cap(o.ClassOfDevice) < int(o.NumResponses) {
		o.ClassOfDevice = make([]uint32, 0, int(o.NumResponses))
	}
	o.ClassOfDevice = o.ClassOfDevice[:int(o.NumResponses)]
	for j := 0; j < int(o.NumResponses); j++ {
		o.ClassOfDevice[j] = decodeUint24(r.Get(3))
	}
	if cap(o.ClockOffset) < int(o.NumResponses) {
		o.ClockOffset = make([]uint16, 0, int(o.NumResponses))
	}
	o.ClockOffset = o.ClockOffset[:int(o.NumResponses)]
	for j := 0; j < int(o.NumResponses); j++ {
		o.ClockOffset[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.RSSI) < int(o.NumResponses) {
		o.RSSI = make([]uint8, 0, int(o.NumResponses))
	}
	o.RSSI = o.RSSI[:int(o.NumResponses)]
	for j := 0; j < int(o.NumResponses); j++ {
		o.RSSI[j] = r.GetOne()
	}
	return r.Valid()
}

// ReadRemoteExtendedFeaturesCompleteEvent represents the event specified in Section 7.7.34
type ReadRemoteExtendedFeaturesCompleteEvent struct {
	Status              uint8
	ConnectionHandle    uint16
	PageNumber          uint8
	MaximumPageNumber   uint8
	ExtendedLMPFeatures uint64
}

func (o *ReadRemoteExtendedFeaturesCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PageNumber = r.GetOne()
	o.MaximumPageNumber = r.GetOne()
	o.ExtendedLMPFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// SynchronousConnectionCompleteEvent represents the event specified in Section 7.7.35
type SynchronousConnectionCompleteEvent struct {
	Status               uint8
	ConnectionHandle     uint16
	BDADDR               [6]byte
	LinkType             uint8
	TransmissionInterval uint8
	RetransmissionWindow uint8
	RXPacketLength       uint16
	TXPacketLength       uint16
	AirMode              uint8
}

func (o *SynchronousConnectionCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	copy(o.BDADDR[:], r.Get(6))
	o.LinkType = r.GetOne()
	o.TransmissionInterval = r.GetOne()
	o.RetransmissionWindow = r.GetOne()
	o.RXPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.TXPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.AirMode = r.GetOne()
	return r.Valid()
}

// SniffSubratingEvent represents the event specified in Section 7.7.37
type SniffSubratingEvent struct {
	Status           uint8
	ConnectionHandle uint16
	MaxTXLatency     uint16
	MaxRXLatency     uint16
	MinRemoteTimeout uint16
	MinLocalTimeout  uint16
}

func (o *SniffSubratingEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxTXLatency = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxRXLatency = binary.LittleEndian.Uint16(r.Get(2))
	o.MinRemoteTimeout = binary.LittleEndian.Uint16(r.Get(2))
	o.MinLocalTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// ExtendedInquiryResultEvent represents the event specified in Section 7.7.38
type ExtendedInquiryResultEvent struct {
	NumResponses            uint8
	BDADDR                  [6]byte
	PageScanRepetitionMode  uint8
	Reserved                uint8
	ClassOfDevice           uint32
	ClockOffset             uint16
	RSSI                    uint8
	ExtendedInquiryResponse [240]byte
}

func (o *ExtendedInquiryResultEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.NumResponses = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	o.PageScanRepetitionMode = r.GetOne()
	o.Reserved = r.GetOne()
	o.ClassOfDevice = decodeUint24(r.Get(3))
	o.ClockOffset = binary.LittleEndian.Uint16(r.Get(2))
	o.RSSI = r.GetOne()
	copy(o.ExtendedInquiryResponse[:], r.Get(240))
	return r.Valid()
}

// EncryptionKeyRefreshCompleteEvent represents the event specified in Section 7.7.39
type EncryptionKeyRefreshCompleteEvent struct {
	Status           uint8
	ConnectionHandle uint16
}

func (o *EncryptionKeyRefreshCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// IOCapabilityRequestEvent represents the event specified in Section 7.7.40
type IOCapabilityRequestEvent struct {
	BDADDR [6]byte
}

func (o *IOCapabilityRequestEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// IOCapabilityResponseEvent represents the event specified in Section 7.7.41
type IOCapabilityResponseEvent struct {
	BDADDR                     [6]byte
	IOCapability               uint8
	OOBDataPresent             uint8
	AuthenticationRequirements uint8
}

func (o *IOCapabilityResponseEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	o.IOCapability = r.GetOne()
	o.OOBDataPresent = r.GetOne()
	o.AuthenticationRequirements = r.GetOne()
	return r.Valid()
}

// UserConfirmationRequestEvent represents the event specified in Section 7.7.42
type UserConfirmationRequestEvent struct {
	BDADDR       [6]byte
	NumericValue uint32
}

func (o *UserConfirmationRequestEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	o.NumericValue = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// UserPasskeyRequestEvent represents the event specified in Section 7.7.43
type UserPasskeyRequestEvent struct {
	BDADDR [6]byte
}

func (o *UserPasskeyRequestEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// RemoteOOBDataRequestEvent represents the event specified in Section 7.7.44
type RemoteOOBDataRequestEvent struct {
	BDADDR [6]byte
}

func (o *RemoteOOBDataRequestEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// SimplePairingCompleteEvent represents the event specified in Section 7.7.45
type SimplePairingCompleteEvent struct {
	Status uint8
	BDADDR [6]byte
}

func (o *SimplePairingCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkSupervisionTimeoutChangedEvent represents the event specified in Section 7.7.46
type LinkSupervisionTimeoutChangedEvent struct {
	ConnectionHandle       uint16
	LinkSupervisionTimeout uint16
}

func (o *LinkSupervisionTimeoutChangedEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LinkSupervisionTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// EnhancedFlushCompleteEvent represents the event specified in Section 7.7.47
type EnhancedFlushCompleteEvent struct {
	Handle uint16
}

func (o *EnhancedFlushCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// UserPasskeyNotificationEvent represents the event specified in Section 7.7.48
type UserPasskeyNotificationEvent struct {
	BDADDR  [6]byte
	Passkey uint32
}

func (o *UserPasskeyNotificationEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	o.Passkey = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// KeypressNotificationEvent represents the event specified in Section 7.7.49
type KeypressNotificationEvent struct {
	BDADDR           [6]byte
	NotificationType uint8
}

func (o *KeypressNotificationEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	o.NotificationType = r.GetOne()
	return r.Valid()
}

// RemoteHostSupportedFeaturesNotificationEvent represents the event specified in Section 7.7.50
type RemoteHostSupportedFeaturesNotificationEvent struct {
	BDADDR                [6]byte
	HostSupportedFeatures uint64
}

func (o *RemoteHostSupportedFeaturesNotificationEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	o.HostSupportedFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// PhysicalLinkCompleteEvent represents the event specified in Section 7.7.51
type PhysicalLinkCompleteEvent struct {
	Status             uint8
	PhysicalLinkHandle uint8
}

func (o *PhysicalLinkCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.PhysicalLinkHandle = r.GetOne()
	return r.Valid()
}

// ChannelSelectedEvent represents the event specified in Section 7.7.52
type ChannelSelectedEvent struct {
	PhysicalLinkHandle uint8
}

func (o *ChannelSelectedEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.PhysicalLinkHandle = r.GetOne()
	return r.Valid()
}

// DisconnectionPhysicalLinkCompleteEvent represents the event specified in Section 7.7.53
type DisconnectionPhysicalLinkCompleteEvent struct {
	Status             uint8
	PhysicalLinkHandle uint8
	Reason             uint8
}

func (o *DisconnectionPhysicalLinkCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.PhysicalLinkHandle = r.GetOne()
	o.Reason = r.GetOne()
	return r.Valid()
}

// PhysicalLinkLossEarlyWarningEvent represents the event specified in Section 7.7.54
type PhysicalLinkLossEarlyWarningEvent struct {
	PhysicalLinkHandle uint8
	LinkLossReason     uint8
}

func (o *PhysicalLinkLossEarlyWarningEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.PhysicalLinkHandle = r.GetOne()
	o.LinkLossReason = r.GetOne()
	return r.Valid()
}

// PhysicalLinkRecoveryEvent represents the event specified in Section 7.7.55
type PhysicalLinkRecoveryEvent struct {
	PhysicalLinkHandle uint8
}

func (o *PhysicalLinkRecoveryEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.PhysicalLinkHandle = r.GetOne()
	return r.Valid()
}

// LogicalLinkCompleteEvent represents the event specified in Section 7.7.56
type LogicalLinkCompleteEvent struct {
	Status             uint8
	LogicalLinkHandle  uint16
	PhysicalLinkHandle uint8
	TXFlowSpecID       uint8
}

func (o *LogicalLinkCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.LogicalLinkHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PhysicalLinkHandle = r.GetOne()
	o.TXFlowSpecID = r.GetOne()
	return r.Valid()
}

// DisconnectionLogicalLinkCompleteEvent represents the event specified in Section 7.7.57
type DisconnectionLogicalLinkCompleteEvent struct {
	Status            uint8
	LogicalLinkHandle uint16
	Reason            uint8
}

func (o *DisconnectionLogicalLinkCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.LogicalLinkHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Reason = r.GetOne()
	return r.Valid()
}

// FlowSpecModifyCompleteEvent represents the event specified in Section 7.7.58
type FlowSpecModifyCompleteEvent struct {
	Status uint8
	Handle uint16
}

func (o *FlowSpecModifyCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// NumberOfCompletedDataBlocksEvent represents the event specified in Section 7.7.59
type NumberOfCompletedDataBlocksEvent struct {
	TotalNumDataBlocks  uint16
	NumHandles          uint8
	Handle              []uint16
	NumCompletedPackets []uint16
	NumCompletedBlocks  []uint16
}

func (o *NumberOfCompletedDataBlocksEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.TotalNumDataBlocks = binary.LittleEndian.Uint16(r.Get(2))
	o.NumHandles = r.GetOne()
	if cap(o.Handle) < int(o.NumHandles) {
		o.Handle = make([]uint16, 0, int(o.NumHandles))
	}
	o.Handle = o.Handle[:int(o.NumHandles)]
	for j := 0; j < int(o.NumHandles); j++ {
		o.Handle[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.NumCompletedPackets) < int(o.NumHandles) {
		o.NumCompletedPackets = make([]uint16, 0, int(o.NumHandles))
	}
	o.NumCompletedPackets = o.NumCompletedPackets[:int(o.NumHandles)]
	for j := 0; j < int(o.NumHandles); j++ {
		o.NumCompletedPackets[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.NumCompletedBlocks) < int(o.NumHandles) {
		o.NumCompletedBlocks = make([]uint16, 0, int(o.NumHandles))
	}
	o.NumCompletedBlocks = o.NumCompletedBlocks[:int(o.NumHandles)]
	for j := 0; j < int(o.NumHandles); j++ {
		o.NumCompletedBlocks[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// ShortRangeModeChangeCompleteEvent represents the event specified in Section 7.7.60
type ShortRangeModeChangeCompleteEvent struct {
	Status              uint8
	PhysicalLinkHandle  uint8
	ShortRangeModeState uint8
}

func (o *ShortRangeModeChangeCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.PhysicalLinkHandle = r.GetOne()
	o.ShortRangeModeState = r.GetOne()
	return r.Valid()
}

// AMPStatusChangeEvent represents the event specified in Section 7.7.61
type AMPStatusChangeEvent struct {
	Status    uint8
	AMPStatus uint8
}

func (o *AMPStatusChangeEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.AMPStatus = r.GetOne()
	return r.Valid()
}

// AMPStartTestEvent represents the event specified in Section 7.7.62
type AMPStartTestEvent struct {
	Status       uint8
	TestScenario uint8
}

func (o *AMPStartTestEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.TestScenario = r.GetOne()
	return r.Valid()
}

// AMPTestEndEvent represents the event specified in Section 7.7.63
type AMPTestEndEvent struct {
	Status       uint8
	TestScenario uint8
}

func (o *AMPTestEndEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	o.TestScenario = r.GetOne()
	return r.Valid()
}

// AMPReceiverReportEvent represents the event specified in Section 7.7.64
type AMPReceiverReportEvent struct {
	ControllerType uint8
	Reason         uint8
}

func (o *AMPReceiverReportEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.ControllerType = r.GetOne()
	o.Reason = r.GetOne()
	return r.Valid()
}

// LEConnectionCompleteEvent represents the event specified in Section 7.7.65.1
type LEConnectionCompleteEvent struct {
	SubeventCode        uint8
	Status              uint8
	ConnectionHandle    uint16
	Role                uint8
	PeerAddressType     uint8
	PeerAddress         [6]byte
	ConnectionInterval  uint16
	ConnectionLatency   uint16
	SupervisionTimeout  uint16
	MasterClockAccuracy uint8
}

func (o *LEConnectionCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Role = r.GetOne()
	o.PeerAddressType = r.GetOne()
	copy(o.PeerAddress[:], r.Get(6))
	o.ConnectionInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.ConnectionLatency = binary.LittleEndian.Uint16(r.Get(2))
	o.SupervisionTimeout = binary.LittleEndian.Uint16(r.Get(2))
	o.MasterClockAccuracy = r.GetOne()
	return r.Valid()
}

// LEAdvertisingReportEvent represents the event specified in Section 7.7.65.2
type LEAdvertisingReportEvent struct {
	SubeventCode uint8
	NumReports   uint8
	EventType    []uint8
	AddressType  []uint8
	Address      [][6]byte
	DataLength   []uint8
	Data         [][]byte
	RSSI         []uint8
}

func (o *LEAdvertisingReportEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.NumReports = r.GetOne()
	if cap(o.EventType) < int(o.NumReports) {
		o.EventType = make([]uint8, 0, int(o.NumReports))
	}
	o.EventType = o.EventType[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.EventType[j] = r.GetOne()
	}
	if cap(o.AddressType) < int(o.NumReports) {
		o.AddressType = make([]uint8, 0, int(o.NumReports))
	}
	o.AddressType = o.AddressType[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.AddressType[j] = r.GetOne()
	}
	if cap(o.Address) < int(o.NumReports) {
		o.Address = make([][6]byte, 0, int(o.NumReports))
	}
	o.Address = o.Address[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		copy(o.Address[j][:], r.Get(6))
	}
	if cap(o.DataLength) < int(o.NumReports) {
		o.DataLength = make([]uint8, 0, int(o.NumReports))
	}
	o.DataLength = o.DataLength[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.DataLength[j] = r.GetOne()
	}
	var0 := 0
	for _, m := range o.DataLength {
		var0 += int(m)
	}
	if cap(o.Data) < var0 {
		o.Data = make([][]byte, 0, var0)
	}
	o.Data = o.Data[:var0]
	for j := 0; j < var0; j++ {
		o.Data[j] = append(o.Data[j][:0], r.Get(int(o.DataLength[j]))...)
	}
	if cap(o.RSSI) < int(o.NumReports) {
		o.RSSI = make([]uint8, 0, int(o.NumReports))
	}
	o.RSSI = o.RSSI[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.RSSI[j] = r.GetOne()
	}
	return r.Valid()
}

// LEConnectionUpdateCompleteEvent represents the event specified in Section 7.7.65.3
type LEConnectionUpdateCompleteEvent struct {
	SubeventCode       uint8
	Status             uint8
	ConnectionHandle   uint16
	ConnectionInterval uint16
	ConnectionLatency  uint16
	SupervisionTimeout uint16
}

func (o *LEConnectionUpdateCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ConnectionInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.ConnectionLatency = binary.LittleEndian.Uint16(r.Get(2))
	o.SupervisionTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEReadRemoteFeaturesCompleteEvent represents the event specified in Section 7.7.65.4
type LEReadRemoteFeaturesCompleteEvent struct {
	SubeventCode     uint8
	Status           uint8
	ConnectionHandle uint16
	LEFeatures       uint64
}

func (o *LEReadRemoteFeaturesCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LEFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// LELongTermKeyRequestEvent represents the event specified in Section 7.7.65.5
type LELongTermKeyRequestEvent struct {
	SubeventCode         uint8
	ConnectionHandle     uint16
	RandomNumber         uint64
	EncryptedDiversifier uint16
}

func (o *LELongTermKeyRequestEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.RandomNumber = binary.LittleEndian.Uint64(r.Get(8))
	o.EncryptedDiversifier = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LERemoteConnectionParameterRequestEvent represents the event specified in Section 7.7.65.6
type LERemoteConnectionParameterRequestEvent struct {
	SubeventCode     uint8
	ConnectionHandle uint16
	IntervalMin      uint16
	IntervalMax      uint16
	Latency          uint16
	Timeout          uint16
}

func (o *LERemoteConnectionParameterRequestEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.IntervalMin = binary.LittleEndian.Uint16(r.Get(2))
	o.IntervalMax = binary.LittleEndian.Uint16(r.Get(2))
	o.Latency = binary.LittleEndian.Uint16(r.Get(2))
	o.Timeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEDataLengthChangeEvent represents the event specified in Section 7.7.65.7
type LEDataLengthChangeEvent struct {
	SubeventCode     uint8
	ConnectionHandle uint16
	MaxTXOctets      uint16
	MaxTXTime        uint16
	MaxRXOctets      uint16
	MaxRXTime        uint16
}

func (o *LEDataLengthChangeEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxTXOctets = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxTXTime = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxRXOctets = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxRXTime = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEReadLocalP256PublicKeyCompleteEvent represents the event specified in Section 7.7.65.8
type LEReadLocalP256PublicKeyCompleteEvent struct {
	SubeventCode       uint8
	Status             uint8
	LocalP256PublicKey [64]byte
}

func (o *LEReadLocalP256PublicKeyCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	copy(o.LocalP256PublicKey[:], r.Get(64))
	return r.Valid()
}

// LEGenerateDHKeyCompleteEvent represents the event specified in Section 7.7.65.9
type LEGenerateDHKeyCompleteEvent struct {
	SubeventCode uint8
	Status       uint8
	DHKey        [32]byte
}

func (o *LEGenerateDHKeyCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	copy(o.DHKey[:], r.Get(32))
	return r.Valid()
}

// LEEnhancedConnectionCompleteEvent represents the event specified in Section 7.7.65.10
type LEEnhancedConnectionCompleteEvent struct {
	SubeventCode                  uint8
	Status                        uint8
	ConnectionHandle              uint16
	Role                          uint8
	PeerAddressType               uint8
	PeerAddress                   [6]byte
	LocalResolvablePrivateAddress [6]byte
	PeerResolvablePrivateAddress  [6]byte
	ConnectionInterval            uint16
	ConnectionLatency             uint16
	SupervisionTimeout            uint16
	MasterClockAccuracy           uint8
}

func (o *LEEnhancedConnectionCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Role = r.GetOne()
	o.PeerAddressType = r.GetOne()
	copy(o.PeerAddress[:], r.Get(6))
	copy(o.LocalResolvablePrivateAddress[:], r.Get(6))
	copy(o.PeerResolvablePrivateAddress[:], r.Get(6))
	o.ConnectionInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.ConnectionLatency = binary.LittleEndian.Uint16(r.Get(2))
	o.SupervisionTimeout = binary.LittleEndian.Uint16(r.Get(2))
	o.MasterClockAccuracy = r.GetOne()
	return r.Valid()
}

// LEDirectedAdvertisingReportEvent represents the event specified in Section 7.7.65.11
type LEDirectedAdvertisingReportEvent struct {
	SubeventCode      uint8
	NumReports        uint8
	EventType         []uint8
	AddressType       []uint8
	Address           [][6]byte
	DirectAddressType []uint8
	DirectAddress     [][6]byte
	RSSI              []uint8
}

func (o *LEDirectedAdvertisingReportEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.NumReports = r.GetOne()
	if cap(o.EventType) < int(o.NumReports) {
		o.EventType = make([]uint8, 0, int(o.NumReports))
	}
	o.EventType = o.EventType[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.EventType[j] = r.GetOne()
	}
	if cap(o.AddressType) < int(o.NumReports) {
		o.AddressType = make([]uint8, 0, int(o.NumReports))
	}
	o.AddressType = o.AddressType[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.AddressType[j] = r.GetOne()
	}
	if cap(o.Address) < int(o.NumReports) {
		o.Address = make([][6]byte, 0, int(o.NumReports))
	}
	o.Address = o.Address[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		copy(o.Address[j][:], r.Get(6))
	}
	if cap(o.DirectAddressType) < int(o.NumReports) {
		o.DirectAddressType = make([]uint8, 0, int(o.NumReports))
	}
	o.DirectAddressType = o.DirectAddressType[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.DirectAddressType[j] = r.GetOne()
	}
	if cap(o.DirectAddress) < int(o.NumReports) {
		o.DirectAddress = make([][6]byte, 0, int(o.NumReports))
	}
	o.DirectAddress = o.DirectAddress[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		copy(o.DirectAddress[j][:], r.Get(6))
	}
	if cap(o.RSSI) < int(o.NumReports) {
		o.RSSI = make([]uint8, 0, int(o.NumReports))
	}
	o.RSSI = o.RSSI[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.RSSI[j] = r.GetOne()
	}
	return r.Valid()
}

// LEPHYUpdateCompleteEvent represents the event specified in Section 7.7.65.12
type LEPHYUpdateCompleteEvent struct {
	SubeventCode     uint8
	Status           uint8
	ConnectionHandle uint16
	TXPHY            uint8
	RXPHY            uint8
}

func (o *LEPHYUpdateCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.TXPHY = r.GetOne()
	o.RXPHY = r.GetOne()
	return r.Valid()
}

// LEExtendedAdvertisingReportEvent represents the event specified in Section 7.7.65.13
type LEExtendedAdvertisingReportEvent struct {
	SubeventCode                uint8
	NumReports                  uint8
	EventType                   []uint16
	AddressType                 []uint8
	Address                     [][6]byte
	PrimaryPHY                  []uint8
	SecondaryPHY                []uint8
	AdvertisingSID              []uint8
	TXPower                     []uint8
	RSSI                        []uint8
	PeriodicAdvertisingInterval []uint16
	DirectAddressType           []uint8
	DirectAddress               [][6]byte
	DataLength                  []uint8
	Data                        [][]byte
}

func (o *LEExtendedAdvertisingReportEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.NumReports = r.GetOne()
	if cap(o.EventType) < int(o.NumReports) {
		o.EventType = make([]uint16, 0, int(o.NumReports))
	}
	o.EventType = o.EventType[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.EventType[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.AddressType) < int(o.NumReports) {
		o.AddressType = make([]uint8, 0, int(o.NumReports))
	}
	o.AddressType = o.AddressType[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.AddressType[j] = r.GetOne()
	}
	if cap(o.Address) < int(o.NumReports) {
		o.Address = make([][6]byte, 0, int(o.NumReports))
	}
	o.Address = o.Address[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		copy(o.Address[j][:], r.Get(6))
	}
	if cap(o.PrimaryPHY) < int(o.NumReports) {
		o.PrimaryPHY = make([]uint8, 0, int(o.NumReports))
	}
	o.PrimaryPHY = o.PrimaryPHY[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.PrimaryPHY[j] = r.GetOne()
	}
	if cap(o.SecondaryPHY) < int(o.NumReports) {
		o.SecondaryPHY = make([]uint8, 0, int(o.NumReports))
	}
	o.SecondaryPHY = o.SecondaryPHY[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.SecondaryPHY[j] = r.GetOne()
	}
	if cap(o.AdvertisingSID) < int(o.NumReports) {
		o.AdvertisingSID = make([]uint8, 0, int(o.NumReports))
	}
	o.AdvertisingSID = o.AdvertisingSID[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.AdvertisingSID[j] = r.GetOne()
	}
	if cap(o.TXPower) < int(o.NumReports) {
		o.TXPower = make([]uint8, 0, int(o.NumReports))
	}
	o.TXPower = o.TXPower[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.TXPower[j] = r.GetOne()
	}
	if cap(o.RSSI) < int(o.NumReports) {
		o.RSSI = make([]uint8, 0, int(o.NumReports))
	}
	o.RSSI = o.RSSI[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.RSSI[j] = r.GetOne()
	}
	if cap(o.PeriodicAdvertisingInterval) < int(o.NumReports) {
		o.PeriodicAdvertisingInterval = make([]uint16, 0, int(o.NumReports))
	}
	o.PeriodicAdvertisingInterval = o.PeriodicAdvertisingInterval[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.PeriodicAdvertisingInterval[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.DirectAddressType) < int(o.NumReports) {
		o.DirectAddressType = make([]uint8, 0, int(o.NumReports))
	}
	o.DirectAddressType = o.DirectAddressType[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.DirectAddressType[j] = r.GetOne()
	}
	if cap(o.DirectAddress) < int(o.NumReports) {
		o.DirectAddress = make([][6]byte, 0, int(o.NumReports))
	}
	o.DirectAddress = o.DirectAddress[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		copy(o.DirectAddress[j][:], r.Get(6))
	}
	if cap(o.DataLength) < int(o.NumReports) {
		o.DataLength = make([]uint8, 0, int(o.NumReports))
	}
	o.DataLength = o.DataLength[:int(o.NumReports)]
	for j := 0; j < int(o.NumReports); j++ {
		o.DataLength[j] = r.GetOne()
	}
	var1 := 0
	for _, m := range o.DataLength {
		var1 += int(m)
	}
	if cap(o.Data) < var1 {
		o.Data = make([][]byte, 0, var1)
	}
	o.Data = o.Data[:var1]
	for j := 0; j < var1; j++ {
		o.Data[j] = append(o.Data[j][:0], r.Get(int(o.DataLength[j]))...)
	}
	return r.Valid()
}

// LEPeriodicAdvertisingSyncEstablishedEvent represents the event specified in Section 7.7.65.14
type LEPeriodicAdvertisingSyncEstablishedEvent struct {
	SubeventCode                uint8
	Status                      uint8
	SyncHandle                  uint16
	AdvertisingSID              uint8
	AdvertiserAddressType       uint8
	AdvertiserAddress           [6]byte
	AdvertiserPHY               uint8
	PeriodicAdvertisingInterval uint16
	AdvertiserClockAccuracy     uint8
}

func (o *LEPeriodicAdvertisingSyncEstablishedEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.AdvertisingSID = r.GetOne()
	o.AdvertiserAddressType = r.GetOne()
	copy(o.AdvertiserAddress[:], r.Get(6))
	o.AdvertiserPHY = r.GetOne()
	o.PeriodicAdvertisingInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.AdvertiserClockAccuracy = r.GetOne()
	return r.Valid()
}

// LEPeriodicAdvertisingReportEvent represents the event specified in Section 7.7.65.15
type LEPeriodicAdvertisingReportEvent struct {
	SubeventCode uint8
	SyncHandle   uint16
	TXPower      uint8
	RSSI         uint8
	CTEType      uint8
	DataStatus   uint8
	DataLength   uint8
	Data         []byte
}

func (o *LEPeriodicAdvertisingReportEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.TXPower = r.GetOne()
	o.RSSI = r.GetOne()
	o.CTEType = r.GetOne()
	o.DataStatus = r.GetOne()
	o.DataLength = r.GetOne()
	o.Data = append(o.Data[:0], r.GetRemainder()...)
	return r.Valid()
}

// LEPeriodicAdvertisingSyncLostEvent represents the event specified in Section 7.7.65.16
type LEPeriodicAdvertisingSyncLostEvent struct {
	SubeventCode uint8
	SyncHandle   uint16
}

func (o *LEPeriodicAdvertisingSyncLostEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEScanTimeoutEvent represents the event specified in Section 7.7.65.17
type LEScanTimeoutEvent struct {
	SubeventCode uint8
}

func (o *LEScanTimeoutEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	return r.Valid()
}

// LEAdvertisingSetTerminatedEvent represents the event specified in Section 7.7.65.18
type LEAdvertisingSetTerminatedEvent struct {
	SubeventCode                          uint8
	Status                                uint8
	AdvertisingHandle                     uint8
	ConnectionHandle                      uint16
	NumCompletedExtendedAdvertisingEvents uint8
}

func (o *LEAdvertisingSetTerminatedEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.AdvertisingHandle = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.NumCompletedExtendedAdvertisingEvents = r.GetOne()
	return r.Valid()
}

// LEScanRequestReceivedEvent represents the event specified in Section 7.7.65.19
type LEScanRequestReceivedEvent struct {
	SubeventCode       uint8
	AdvertisingHandle  uint8
	ScannerAddressType uint8
	ScannerAddress     [6]byte
}

func (o *LEScanRequestReceivedEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.AdvertisingHandle = r.GetOne()
	o.ScannerAddressType = r.GetOne()
	copy(o.ScannerAddress[:], r.Get(6))
	return r.Valid()
}

// LEChannelSelectionAlgorithmEvent represents the event specified in Section 7.7.65.20
type LEChannelSelectionAlgorithmEvent struct {
	SubeventCode              uint8
	ConnectionHandle          uint16
	ChannelSelectionAlgorithm uint8
}

func (o *LEChannelSelectionAlgorithmEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ChannelSelectionAlgorithm = r.GetOne()
	return r.Valid()
}

// LEConnectionlessIQReportEvent represents the event specified in Section 7.7.65.21
type LEConnectionlessIQReportEvent struct {
	SubeventCode         uint8
	SyncHandle           uint16
	ChannelIndex         uint8
	RSSI                 uint16
	RSSIAntennaID        uint8
	CTEType              uint8
	SlotDurations        uint8
	PacketStatus         uint8
	PeriodicEventCounter uint16
	SampleCount          uint8
	ISample              []uint8
	QSample              []uint8
}

func (o *LEConnectionlessIQReportEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ChannelIndex = r.GetOne()
	o.RSSI = binary.LittleEndian.Uint16(r.Get(2))
	o.RSSIAntennaID = r.GetOne()
	o.CTEType = r.GetOne()
	o.SlotDurations = r.GetOne()
	o.PacketStatus = r.GetOne()
	o.PeriodicEventCounter = binary.LittleEndian.Uint16(r.Get(2))
	o.SampleCount = r.GetOne()
	if cap(o.ISample) < int(o.SampleCount) {
		o.ISample = make([]uint8, 0, int(o.SampleCount))
	}
	o.ISample = o.ISample[:int(o.SampleCount)]
	for j := 0; j < int(o.SampleCount); j++ {
		o.ISample[j] = r.GetOne()
	}
	if cap(o.QSample) < int(o.SampleCount) {
		o.QSample = make([]uint8, 0, int(o.SampleCount))
	}
	o.QSample = o.QSample[:int(o.SampleCount)]
	for j := 0; j < int(o.SampleCount); j++ {
		o.QSample[j] = r.GetOne()
	}
	return r.Valid()
}

// LEConnectionIQReportEvent represents the event specified in Section 7.7.65.22
type LEConnectionIQReportEvent struct {
	SubeventCode           uint8
	ConnectionHandle       uint16
	RXPHY                  uint8
	DataChannelIndex       uint8
	RSSI                   uint16
	RSSIAntennaID          uint8
	CTEType                uint8
	SlotDurations          uint8
	PacketStatus           uint8
	ConnectionEventCounter uint16
	SampleCount            uint8
	ISample                []uint8
	QSample                []uint8
}

func (o *LEConnectionIQReportEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.RXPHY = r.GetOne()
	o.DataChannelIndex = r.GetOne()
	o.RSSI = binary.LittleEndian.Uint16(r.Get(2))
	o.RSSIAntennaID = r.GetOne()
	o.CTEType = r.GetOne()
	o.SlotDurations = r.GetOne()
	o.PacketStatus = r.GetOne()
	o.ConnectionEventCounter = binary.LittleEndian.Uint16(r.Get(2))
	o.SampleCount = r.GetOne()
	if cap(o.ISample) < int(o.SampleCount) {
		o.ISample = make([]uint8, 0, int(o.SampleCount))
	}
	o.ISample = o.ISample[:int(o.SampleCount)]
	for j := 0; j < int(o.SampleCount); j++ {
		o.ISample[j] = r.GetOne()
	}
	if cap(o.QSample) < int(o.SampleCount) {
		o.QSample = make([]uint8, 0, int(o.SampleCount))
	}
	o.QSample = o.QSample[:int(o.SampleCount)]
	for j := 0; j < int(o.SampleCount); j++ {
		o.QSample[j] = r.GetOne()
	}
	return r.Valid()
}

// LECTERequestFailedEvent represents the event specified in Section 7.7.65.23
type LECTERequestFailedEvent struct {
	SubeventCode     uint8
	Status           uint8
	ConnectionHandle uint16
}

func (o *LECTERequestFailedEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEPeriodicAdvertisingSyncTransferReceivedEvent represents the event specified in Section 7.7.65.24
type LEPeriodicAdvertisingSyncTransferReceivedEvent struct {
	SubeventCode                uint8
	Status                      uint8
	ConnectionHandle            uint16
	ServiceData                 uint16
	SyncHandle                  uint16
	AdvertisingSID              uint8
	AdvertiserAddressType       uint8
	AdvertiserAddress           [6]byte
	AdvertiserPHY               uint8
	PeriodicAdvertisingInterval uint16
	AdvertiserClockAccuracy     uint8
}

func (o *LEPeriodicAdvertisingSyncTransferReceivedEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ServiceData = binary.LittleEndian.Uint16(r.Get(2))
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.AdvertisingSID = r.GetOne()
	o.AdvertiserAddressType = r.GetOne()
	copy(o.AdvertiserAddress[:], r.Get(6))
	o.AdvertiserPHY = r.GetOne()
	o.PeriodicAdvertisingInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.AdvertiserClockAccuracy = r.GetOne()
	return r.Valid()
}

// LECISEstablishedEvent represents the event specified in Section 7.7.65.25
type LECISEstablishedEvent struct {
	SubeventCode     uint8
	Status           uint8
	ConnectionHandle uint16
	PHYMToS          uint8
	PHYSToM          uint8
	NSE              uint8
	BNMToS           uint8
	BNSToM           uint8
	FTMToS           uint8
	FTSToM           uint8
	MaxPDUMToS       uint16
	MaxPDUSToM       uint16
	ISOinterval      uint16
}

func (o *LECISEstablishedEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PHYMToS = r.GetOne()
	o.PHYSToM = r.GetOne()
	o.NSE = r.GetOne()
	o.BNMToS = r.GetOne()
	o.BNSToM = r.GetOne()
	o.FTMToS = r.GetOne()
	o.FTSToM = r.GetOne()
	o.MaxPDUMToS = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxPDUSToM = binary.LittleEndian.Uint16(r.Get(2))
	o.ISOinterval = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LECISRequestEvent represents the event specified in Section 7.7.65.26
type LECISRequestEvent struct {
	SubeventCode        uint8
	ACLConnectionHandle uint16
	CISConnectonHandle  uint16
	CIGID               uint8
	CISID               uint8
}

func (o *LECISRequestEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.ACLConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.CISConnectonHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.CIGID = r.GetOne()
	o.CISID = r.GetOne()
	return r.Valid()
}

// LECreateBIGCompleteEvent represents the event specified in Section 7.7.65.27
type LECreateBIGCompleteEvent struct {
	SubeventCode        uint8
	Status              uint8
	BIGHandle           uint8
	BIGSyncDelay        uint32
	TransportLatencyBIG uint32
	PHY                 uint8
	NSE                 uint8
	BN                  uint8
	PTO                 uint8
	IRC                 uint8
	MaxPDU              uint16
	ISOInterval         uint16
	NumBIS              uint8
	ConnectionHandle    []uint16
}

func (o *LECreateBIGCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.BIGHandle = r.GetOne()
	o.BIGSyncDelay = decodeUint24(r.Get(3))
	o.TransportLatencyBIG = decodeUint24(r.Get(3))
	o.PHY = r.GetOne()
	o.NSE = r.GetOne()
	o.BN = r.GetOne()
	o.PTO = r.GetOne()
	o.IRC = r.GetOne()
	o.MaxPDU = binary.LittleEndian.Uint16(r.Get(2))
	o.ISOInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.NumBIS = r.GetOne()
	if cap(o.ConnectionHandle) < int(o.NumBIS) {
		o.ConnectionHandle = make([]uint16, 0, int(o.NumBIS))
	}
	o.ConnectionHandle = o.ConnectionHandle[:int(o.NumBIS)]
	for j := 0; j < int(o.NumBIS); j++ {
		o.ConnectionHandle[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// LETerminateBIGCompleteEvent represents the event specified in Section 7.7.65.28
type LETerminateBIGCompleteEvent struct {
	SubeventCode uint8
	BIGHandle    uint8
	Reason       uint8
}

func (o *LETerminateBIGCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.BIGHandle = r.GetOne()
	o.Reason = r.GetOne()
	return r.Valid()
}

// LEBIGSyncEstablishedEvent represents the event specified in Section 7.7.65.29
type LEBIGSyncEstablishedEvent struct {
	SubeventCode        uint8
	Status              uint8
	BIGHandle           uint8
	TransportLatencyBIG uint32
	NSE                 uint8
	BN                  uint8
	PTO                 uint8
	IRC                 uint8
	MaxPDU              uint16
	ISOInterval         uint16
	NumBIS              uint8
	ConnectionHandle    []uint16
}

func (o *LEBIGSyncEstablishedEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.BIGHandle = r.GetOne()
	o.TransportLatencyBIG = decodeUint24(r.Get(3))
	o.NSE = r.GetOne()
	o.BN = r.GetOne()
	o.PTO = r.GetOne()
	o.IRC = r.GetOne()
	o.MaxPDU = binary.LittleEndian.Uint16(r.Get(2))
	o.ISOInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.NumBIS = r.GetOne()
	if cap(o.ConnectionHandle) < int(o.NumBIS) {
		o.ConnectionHandle = make([]uint16, 0, int(o.NumBIS))
	}
	o.ConnectionHandle = o.ConnectionHandle[:int(o.NumBIS)]
	for j := 0; j < int(o.NumBIS); j++ {
		o.ConnectionHandle[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// LEBIGSyncLostEvent represents the event specified in Section 7.7.65.30
type LEBIGSyncLostEvent struct {
	SubeventCode uint8
	BIGHandle    uint8
	Reason       uint8
}

func (o *LEBIGSyncLostEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.BIGHandle = r.GetOne()
	o.Reason = r.GetOne()
	return r.Valid()
}

// LERequestPeerSCACompleteEvent represents the event specified in Section 7.7.65.31
type LERequestPeerSCACompleteEvent struct {
	SubeventCode      uint8
	Status            uint8
	ConnectionHandle  uint16
	PeerClockAccuracy uint8
}

func (o *LERequestPeerSCACompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PeerClockAccuracy = r.GetOne()
	return r.Valid()
}

// LEPathLossThresholdEvent represents the event specified in Section 7.7.65.32
type LEPathLossThresholdEvent struct {
	SubeventCode     uint8
	ConnectionHandle uint16
	CurrentPathLoss  uint8
	ZoneEntered      uint8
}

func (o *LEPathLossThresholdEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.CurrentPathLoss = r.GetOne()
	o.ZoneEntered = r.GetOne()
	return r.Valid()
}

// LETransmitPowerReportingEvent represents the event specified in Section 7.7.65.33
type LETransmitPowerReportingEvent struct {
	SubeventCode           uint8
	Status                 uint8
	ConnectionHandle       uint16
	Reason                 uint8
	PHY                    uint8
	TransmitPowerLevel     uint8
	TransmitPowerLevelFlag uint8
	Delta                  uint8
}

func (o *LETransmitPowerReportingEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Reason = r.GetOne()
	o.PHY = r.GetOne()
	o.TransmitPowerLevel = r.GetOne()
	o.TransmitPowerLevelFlag = r.GetOne()
	o.Delta = r.GetOne()
	return r.Valid()
}

// LEBIGInfoAdvertisingReportEvent represents the event specified in Section 7.7.65.34
type LEBIGInfoAdvertisingReportEvent struct {
	SubeventCode uint8
	SyncHandle   uint16
	NumBIS       uint8
	NSE          uint8
	ISOInterval  uint16
	BN           uint8
	PTO          uint8
	IRC          uint8
	MaxPDU       uint16
	SDUInterval  uint32
	MaxSDU       uint16
	PHY          uint8
	Framing      uint8
	Encryption   uint8
}

func (o *LEBIGInfoAdvertisingReportEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.SubeventCode = r.GetOne()
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.NumBIS = r.GetOne()
	o.NSE = r.GetOne()
	o.ISOInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.BN = r.GetOne()
	o.PTO = r.GetOne()
	o.IRC = r.GetOne()
	o.MaxPDU = binary.LittleEndian.Uint16(r.Get(2))
	o.SDUInterval = decodeUint24(r.Get(3))
	o.MaxSDU = binary.LittleEndian.Uint16(r.Get(2))
	o.PHY = r.GetOne()
	o.Framing = r.GetOne()
	o.Encryption = r.GetOne()
	return r.Valid()
}

// TriggeredClockCaptureEvent represents the event specified in Section 7.7.66
type TriggeredClockCaptureEvent struct {
	ConnectionHandle uint16
	WhichClock       uint8
	Clock            uint32
	SlotOffset       uint16
}

func (o *TriggeredClockCaptureEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.WhichClock = r.GetOne()
	o.Clock = binary.LittleEndian.Uint32(r.Get(4))
	o.SlotOffset = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// SynchronizationTrainReceivedEvent represents the event specified in Section 7.7.68
type SynchronizationTrainReceivedEvent struct {
	Status                               uint8
	BDADDR                               [6]byte
	ClockOffset                          uint32
	AFHChannelMap                        [10]byte
	LTADDR                               uint8
	NextBroadcastInstant                 uint32
	ConnectionlessSlaveBroadcastInterval uint16
	ServiceData                          uint8
}

func (o *SynchronizationTrainReceivedEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	o.ClockOffset = binary.LittleEndian.Uint32(r.Get(4))
	copy(o.AFHChannelMap[:], r.Get(10))
	o.LTADDR = r.GetOne()
	o.NextBroadcastInstant = binary.LittleEndian.Uint32(r.Get(4))
	o.ConnectionlessSlaveBroadcastInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.ServiceData = r.GetOne()
	return r.Valid()
}

// ConnectionlessSlaveBroadcastReceiveEvent represents the event specified in Section 7.7.69
type ConnectionlessSlaveBroadcastReceiveEvent struct {
	BDADDR     [6]byte
	LTADDR     uint8
	CLK        uint32
	Offset     uint32
	RXStatus   uint8
	Fragment   uint8
	DataLength uint8
	Data       []byte
}

func (o *ConnectionlessSlaveBroadcastReceiveEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	o.LTADDR = r.GetOne()
	o.CLK = binary.LittleEndian.Uint32(r.Get(4))
	o.Offset = binary.LittleEndian.Uint32(r.Get(4))
	o.RXStatus = r.GetOne()
	o.Fragment = r.GetOne()
	o.DataLength = r.GetOne()
	o.Data = append(o.Data[:0], r.GetRemainder()...)
	return r.Valid()
}

// ConnectionlessSlaveBroadcastTimeoutEvent represents the event specified in Section 7.7.70
type ConnectionlessSlaveBroadcastTimeoutEvent struct {
	BDADDR [6]byte
	LTADDR uint8
}

func (o *ConnectionlessSlaveBroadcastTimeoutEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.BDADDR[:], r.Get(6))
	o.LTADDR = r.GetOne()
	return r.Valid()
}

// TruncatedPageCompleteEvent represents the event specified in Section 7.7.71
type TruncatedPageCompleteEvent struct {
	Status uint8
	BDADDR [6]byte
}

func (o *TruncatedPageCompleteEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// SlavePageResponseTimeoutEvent represents the event specified in Section 7.7.72
type SlavePageResponseTimeoutEvent struct {
	ChannelMap [10]byte
}

func (o *SlavePageResponseTimeoutEvent) decode(data []byte) bool {
	r := reader{data: data}
	copy(o.ChannelMap[:], r.Get(10))
	return r.Valid()
}

// InquiryResponseNotificationEvent represents the event specified in Section 7.7.74
type InquiryResponseNotificationEvent struct {
	LAP  uint32
	RSSI uint8
}

func (o *InquiryResponseNotificationEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.LAP = decodeUint24(r.Get(3))
	o.RSSI = r.GetOne()
	return r.Valid()
}

// AuthenticatedPayloadTimeoutExpiredEvent represents the event specified in Section 7.7.75
type AuthenticatedPayloadTimeoutExpiredEvent struct {
	ConnectionHandle uint16
}

func (o *AuthenticatedPayloadTimeoutExpiredEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// SAMStatusChangeEvent represents the event specified in Section 7.7.76
type SAMStatusChangeEvent struct {
	ConnectionHandle        uint16
	LocalSAMIndex           uint8
	LocalSAMTXAvailability  uint8
	LocalSAMRXAvailability  uint8
	RemoteSAMIndex          uint8
	RemoteSAMTXAvailability uint8
	RemoteSAMRXAvailability uint8
	LEEventMask             uint64
	Status                  uint8
}

func (o *SAMStatusChangeEvent) decode(data []byte) bool {
	r := reader{data: data}
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LocalSAMIndex = r.GetOne()
	o.LocalSAMTXAvailability = r.GetOne()
	o.LocalSAMRXAvailability = r.GetOne()
	o.RemoteSAMIndex = r.GetOne()
	o.RemoteSAMTXAvailability = r.GetOne()
	o.RemoteSAMRXAvailability = r.GetOne()
	o.LEEventMask = binary.LittleEndian.Uint64(r.Get(8))
	o.Status = r.GetOne()
	return r.Valid()
}
