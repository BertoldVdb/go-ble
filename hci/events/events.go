package hcievents

import (
	"sync"
	"encoding/binary"
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
	bleutil "github.com/BertoldVdb/go-ble/util"
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	"github.com/sirupsen/logrus"
)

// InquiryCompleteEvent represents the event specified in Section 7.7.1
type InquiryCompleteEvent struct {
	Status uint8
}

func (o *InquiryCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	return r.Valid()
}

// InquiryCompleteEventCallbackType is the type of the callback function for InquiryCompleteEvent.
type InquiryCompleteEventCallbackType func(*InquiryCompleteEvent) *InquiryCompleteEvent

// SetInquiryCompleteEventCallback configures the callback for InquiryCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetInquiryCompleteEventCallback(cb InquiryCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.inquiryCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(0, cb != nil);
}
// InquiryResultEvent represents the event specified in Section 7.7.2
type InquiryResultEvent struct {
	NumResponses uint8
	BDADDR []bleutil.MacAddr
	PageScanRepetitionMode []uint8
	Reserved []uint16
	ClassOfDevice []uint32
	ClockOffset []uint16
}

func (o *InquiryResultEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.NumResponses = uint8(r.GetOne())
	if cap(o.BDADDR) < int(o.NumResponses) {
		o.BDADDR = make([]bleutil.MacAddr, 0, int(o.NumResponses))
	}
	o.BDADDR = o.BDADDR[:int(o.NumResponses)]
	for j:=0; j<int(o.NumResponses); j++ {
		o.BDADDR[j].Decode(r.Get(6))
	}
	if cap(o.PageScanRepetitionMode) < int(o.NumResponses) {
		o.PageScanRepetitionMode = make([]uint8, 0, int(o.NumResponses))
	}
	o.PageScanRepetitionMode = o.PageScanRepetitionMode[:int(o.NumResponses)]
	for j:=0; j<int(o.NumResponses); j++ {
		o.PageScanRepetitionMode[j] = uint8(r.GetOne())
	}
	if cap(o.Reserved) < int(o.NumResponses) {
		o.Reserved = make([]uint16, 0, int(o.NumResponses))
	}
	o.Reserved = o.Reserved[:int(o.NumResponses)]
	for j:=0; j<int(o.NumResponses); j++ {
		o.Reserved[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.ClassOfDevice) < int(o.NumResponses) {
		o.ClassOfDevice = make([]uint32, 0, int(o.NumResponses))
	}
	o.ClassOfDevice = o.ClassOfDevice[:int(o.NumResponses)]
	for j:=0; j<int(o.NumResponses); j++ {
		o.ClassOfDevice[j] = bleutil.DecodeUint24(r.Get(3))
	}
	if cap(o.ClockOffset) < int(o.NumResponses) {
		o.ClockOffset = make([]uint16, 0, int(o.NumResponses))
	}
	o.ClockOffset = o.ClockOffset[:int(o.NumResponses)]
	for j:=0; j<int(o.NumResponses); j++ {
		o.ClockOffset[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// InquiryResultEventCallbackType is the type of the callback function for InquiryResultEvent.
type InquiryResultEventCallbackType func(*InquiryResultEvent) *InquiryResultEvent

// SetInquiryResultEventCallback configures the callback for InquiryResultEvent. Passing nil will disable the callback.
func (e *EventHandler) SetInquiryResultEventCallback(cb InquiryResultEventCallbackType) error {
	e.cbMutex.Lock()
	e.inquiryResultEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(1, cb != nil);
}
// ConnectionCompleteEvent represents the event specified in Section 7.7.3
type ConnectionCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
	BDADDR bleutil.MacAddr
	LinkType uint8
	EncryptionEnabled uint8
}

func (o *ConnectionCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.BDADDR.Decode(r.Get(6))
	o.LinkType = uint8(r.GetOne())
	o.EncryptionEnabled = uint8(r.GetOne())
	return r.Valid()
}

// ConnectionCompleteEventCallbackType is the type of the callback function for ConnectionCompleteEvent.
type ConnectionCompleteEventCallbackType func(*ConnectionCompleteEvent) *ConnectionCompleteEvent

// SetConnectionCompleteEventCallback configures the callback for ConnectionCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetConnectionCompleteEventCallback(cb ConnectionCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.connectionCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(2, cb != nil);
}
// ConnectionRequestEvent represents the event specified in Section 7.7.4
type ConnectionRequestEvent struct {
	BDADDR bleutil.MacAddr
	ClassOfDevice uint32
	LinkType uint8
}

func (o *ConnectionRequestEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	o.ClassOfDevice = bleutil.DecodeUint24(r.Get(3))
	o.LinkType = uint8(r.GetOne())
	return r.Valid()
}

// ConnectionRequestEventCallbackType is the type of the callback function for ConnectionRequestEvent.
type ConnectionRequestEventCallbackType func(*ConnectionRequestEvent) *ConnectionRequestEvent

// SetConnectionRequestEventCallback configures the callback for ConnectionRequestEvent. Passing nil will disable the callback.
func (e *EventHandler) SetConnectionRequestEventCallback(cb ConnectionRequestEventCallbackType) error {
	e.cbMutex.Lock()
	e.connectionRequestEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(3, cb != nil);
}
// DisconnectionCompleteEvent represents the event specified in Section 7.7.5
type DisconnectionCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
	Reason uint8
}

func (o *DisconnectionCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Reason = uint8(r.GetOne())
	return r.Valid()
}

// DisconnectionCompleteEventCallbackType is the type of the callback function for DisconnectionCompleteEvent.
type DisconnectionCompleteEventCallbackType func(*DisconnectionCompleteEvent) *DisconnectionCompleteEvent

// SetDisconnectionCompleteEventCallback configures the callback for DisconnectionCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetDisconnectionCompleteEventCallback(cb DisconnectionCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.disconnectionCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(4, cb != nil);
}
// AuthenticationCompleteEvent represents the event specified in Section 7.7.6
type AuthenticationCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *AuthenticationCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// AuthenticationCompleteEventCallbackType is the type of the callback function for AuthenticationCompleteEvent.
type AuthenticationCompleteEventCallbackType func(*AuthenticationCompleteEvent) *AuthenticationCompleteEvent

// SetAuthenticationCompleteEventCallback configures the callback for AuthenticationCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetAuthenticationCompleteEventCallback(cb AuthenticationCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.authenticationCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(5, cb != nil);
}
// RemoteNameRequestCompleteEvent represents the event specified in Section 7.7.7
type RemoteNameRequestCompleteEvent struct {
	Status uint8
	BDADDR bleutil.MacAddr
	RemoteName [248]byte
}

func (o *RemoteNameRequestCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.BDADDR.Decode(r.Get(6))
	copy(o.RemoteName[:], r.Get(248))
	return r.Valid()
}

// RemoteNameRequestCompleteEventCallbackType is the type of the callback function for RemoteNameRequestCompleteEvent.
type RemoteNameRequestCompleteEventCallbackType func(*RemoteNameRequestCompleteEvent) *RemoteNameRequestCompleteEvent

// SetRemoteNameRequestCompleteEventCallback configures the callback for RemoteNameRequestCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetRemoteNameRequestCompleteEventCallback(cb RemoteNameRequestCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.remoteNameRequestCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(6, cb != nil);
}
// EncryptionChangeEvent represents the event specified in Section 7.7.8
type EncryptionChangeEvent struct {
	Status uint8
	ConnectionHandle uint16
	EncryptionEnabled uint8
}

func (o *EncryptionChangeEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.EncryptionEnabled = uint8(r.GetOne())
	return r.Valid()
}

// EncryptionChangeEventCallbackType is the type of the callback function for EncryptionChangeEvent.
type EncryptionChangeEventCallbackType func(*EncryptionChangeEvent) *EncryptionChangeEvent

// SetEncryptionChangeEventCallback configures the callback for EncryptionChangeEvent. Passing nil will disable the callback.
func (e *EventHandler) SetEncryptionChangeEventCallback(cb EncryptionChangeEventCallbackType) error {
	e.cbMutex.Lock()
	e.encryptionChangeEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(7, cb != nil);
}
// ChangeConnectionLinkKeyCompleteEvent represents the event specified in Section 7.7.9
type ChangeConnectionLinkKeyCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *ChangeConnectionLinkKeyCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// ChangeConnectionLinkKeyCompleteEventCallbackType is the type of the callback function for ChangeConnectionLinkKeyCompleteEvent.
type ChangeConnectionLinkKeyCompleteEventCallbackType func(*ChangeConnectionLinkKeyCompleteEvent) *ChangeConnectionLinkKeyCompleteEvent

// SetChangeConnectionLinkKeyCompleteEventCallback configures the callback for ChangeConnectionLinkKeyCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetChangeConnectionLinkKeyCompleteEventCallback(cb ChangeConnectionLinkKeyCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.changeConnectionLinkKeyCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(8, cb != nil);
}
// MasterLinkKeyCompleteEvent represents the event specified in Section 7.7.10
type MasterLinkKeyCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
	KeyFlag uint8
}

func (o *MasterLinkKeyCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.KeyFlag = uint8(r.GetOne())
	return r.Valid()
}

// MasterLinkKeyCompleteEventCallbackType is the type of the callback function for MasterLinkKeyCompleteEvent.
type MasterLinkKeyCompleteEventCallbackType func(*MasterLinkKeyCompleteEvent) *MasterLinkKeyCompleteEvent

// SetMasterLinkKeyCompleteEventCallback configures the callback for MasterLinkKeyCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetMasterLinkKeyCompleteEventCallback(cb MasterLinkKeyCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.masterLinkKeyCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(9, cb != nil);
}
// ReadRemoteSupportedFeaturesCompleteEvent represents the event specified in Section 7.7.11
type ReadRemoteSupportedFeaturesCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
	LMPFeatures uint64
}

func (o *ReadRemoteSupportedFeaturesCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LMPFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// ReadRemoteSupportedFeaturesCompleteEventCallbackType is the type of the callback function for ReadRemoteSupportedFeaturesCompleteEvent.
type ReadRemoteSupportedFeaturesCompleteEventCallbackType func(*ReadRemoteSupportedFeaturesCompleteEvent) *ReadRemoteSupportedFeaturesCompleteEvent

// SetReadRemoteSupportedFeaturesCompleteEventCallback configures the callback for ReadRemoteSupportedFeaturesCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetReadRemoteSupportedFeaturesCompleteEventCallback(cb ReadRemoteSupportedFeaturesCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.readRemoteSupportedFeaturesCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(10, cb != nil);
}
// ReadRemoteVersionInformationCompleteEvent represents the event specified in Section 7.7.12
type ReadRemoteVersionInformationCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
	Version uint8
	ManufacturerName uint16
	Subversion uint16
}

func (o *ReadRemoteVersionInformationCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Version = uint8(r.GetOne())
	o.ManufacturerName = binary.LittleEndian.Uint16(r.Get(2))
	o.Subversion = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// ReadRemoteVersionInformationCompleteEventCallbackType is the type of the callback function for ReadRemoteVersionInformationCompleteEvent.
type ReadRemoteVersionInformationCompleteEventCallbackType func(*ReadRemoteVersionInformationCompleteEvent) *ReadRemoteVersionInformationCompleteEvent

// SetReadRemoteVersionInformationCompleteEventCallback configures the callback for ReadRemoteVersionInformationCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetReadRemoteVersionInformationCompleteEventCallback(cb ReadRemoteVersionInformationCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.readRemoteVersionInformationCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(11, cb != nil);
}
// QoSSetupCompleteEvent represents the event specified in Section 7.7.13
type QoSSetupCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
	Unused uint8
	ServiceType uint8
	TokenRate uint32
	PeakBandwidth uint32
	Latency uint32
	DelayVariation uint32
}

func (o *QoSSetupCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Unused = uint8(r.GetOne())
	o.ServiceType = uint8(r.GetOne())
	o.TokenRate = binary.LittleEndian.Uint32(r.Get(4))
	o.PeakBandwidth = binary.LittleEndian.Uint32(r.Get(4))
	o.Latency = binary.LittleEndian.Uint32(r.Get(4))
	o.DelayVariation = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// QoSSetupCompleteEventCallbackType is the type of the callback function for QoSSetupCompleteEvent.
type QoSSetupCompleteEventCallbackType func(*QoSSetupCompleteEvent) *QoSSetupCompleteEvent

// SetQoSSetupCompleteEventCallback configures the callback for QoSSetupCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetQoSSetupCompleteEventCallback(cb QoSSetupCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.qoSSetupCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(12, cb != nil);
}
// CommandCompleteEvent represents the event specified in Section 7.7.14
type CommandCompleteEvent struct {
	NumHCICommandPackets uint8
	CommandOpcode uint16
	ReturnParameters []byte
}

func (o *CommandCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.NumHCICommandPackets = uint8(r.GetOne())
	o.CommandOpcode = binary.LittleEndian.Uint16(r.Get(2))
	o.ReturnParameters = append(o.ReturnParameters[:0], r.GetRemainder()...)
	return r.Valid()
}

// CommandCompleteEventCallbackType is the type of the callback function for CommandCompleteEvent.
type CommandCompleteEventCallbackType func(*CommandCompleteEvent) *CommandCompleteEvent

// SetCommandCompleteEventCallback configures the callback for CommandCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetCommandCompleteEventCallback(cb CommandCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.commandCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(13, cb != nil);
}
// CommandStatusEvent represents the event specified in Section 7.7.15
type CommandStatusEvent struct {
	Status uint8
	NumHCICommandPackets uint8
	CommandOpcode uint16
}

func (o *CommandStatusEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.NumHCICommandPackets = uint8(r.GetOne())
	o.CommandOpcode = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// CommandStatusEventCallbackType is the type of the callback function for CommandStatusEvent.
type CommandStatusEventCallbackType func(*CommandStatusEvent) *CommandStatusEvent

// SetCommandStatusEventCallback configures the callback for CommandStatusEvent. Passing nil will disable the callback.
func (e *EventHandler) SetCommandStatusEventCallback(cb CommandStatusEventCallbackType) error {
	e.cbMutex.Lock()
	e.commandStatusEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(14, cb != nil);
}
// HardwareErrorEvent represents the event specified in Section 7.7.16
type HardwareErrorEvent struct {
	HardwareCode uint8
}

func (o *HardwareErrorEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.HardwareCode = uint8(r.GetOne())
	return r.Valid()
}

// HardwareErrorEventCallbackType is the type of the callback function for HardwareErrorEvent.
type HardwareErrorEventCallbackType func(*HardwareErrorEvent) *HardwareErrorEvent

// SetHardwareErrorEventCallback configures the callback for HardwareErrorEvent. Passing nil will disable the callback.
func (e *EventHandler) SetHardwareErrorEventCallback(cb HardwareErrorEventCallbackType) error {
	e.cbMutex.Lock()
	e.hardwareErrorEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(15, cb != nil);
}
// FlushOccurredEvent represents the event specified in Section 7.7.17
type FlushOccurredEvent struct {
	Handle uint16
}

func (o *FlushOccurredEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// FlushOccurredEventCallbackType is the type of the callback function for FlushOccurredEvent.
type FlushOccurredEventCallbackType func(*FlushOccurredEvent) *FlushOccurredEvent

// SetFlushOccurredEventCallback configures the callback for FlushOccurredEvent. Passing nil will disable the callback.
func (e *EventHandler) SetFlushOccurredEventCallback(cb FlushOccurredEventCallbackType) error {
	e.cbMutex.Lock()
	e.flushOccurredEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(16, cb != nil);
}
// RoleChangeEvent represents the event specified in Section 7.7.18
type RoleChangeEvent struct {
	Status uint8
	BDADDR bleutil.MacAddr
	NewRole uint8
}

func (o *RoleChangeEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.BDADDR.Decode(r.Get(6))
	o.NewRole = uint8(r.GetOne())
	return r.Valid()
}

// RoleChangeEventCallbackType is the type of the callback function for RoleChangeEvent.
type RoleChangeEventCallbackType func(*RoleChangeEvent) *RoleChangeEvent

// SetRoleChangeEventCallback configures the callback for RoleChangeEvent. Passing nil will disable the callback.
func (e *EventHandler) SetRoleChangeEventCallback(cb RoleChangeEventCallbackType) error {
	e.cbMutex.Lock()
	e.roleChangeEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(17, cb != nil);
}
// NumberOfCompletedPacketsEvent represents the event specified in Section 7.7.19
type NumberOfCompletedPacketsEvent struct {
	NumHandles uint8
	ConnectionHandle []uint16
	NumCompletedPackets []uint16
}

func (o *NumberOfCompletedPacketsEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.NumHandles = uint8(r.GetOne())
	if cap(o.ConnectionHandle) < int(o.NumHandles) {
		o.ConnectionHandle = make([]uint16, 0, int(o.NumHandles))
	}
	o.ConnectionHandle = o.ConnectionHandle[:int(o.NumHandles)]
	for j:=0; j<int(o.NumHandles); j++ {
		o.ConnectionHandle[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.NumCompletedPackets) < int(o.NumHandles) {
		o.NumCompletedPackets = make([]uint16, 0, int(o.NumHandles))
	}
	o.NumCompletedPackets = o.NumCompletedPackets[:int(o.NumHandles)]
	for j:=0; j<int(o.NumHandles); j++ {
		o.NumCompletedPackets[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// NumberOfCompletedPacketsEventCallbackType is the type of the callback function for NumberOfCompletedPacketsEvent.
type NumberOfCompletedPacketsEventCallbackType func(*NumberOfCompletedPacketsEvent) *NumberOfCompletedPacketsEvent

// SetNumberOfCompletedPacketsEventCallback configures the callback for NumberOfCompletedPacketsEvent. Passing nil will disable the callback.
func (e *EventHandler) SetNumberOfCompletedPacketsEventCallback(cb NumberOfCompletedPacketsEventCallbackType) error {
	e.cbMutex.Lock()
	e.numberOfCompletedPacketsEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(18, cb != nil);
}
// ModeChangeEvent represents the event specified in Section 7.7.20
type ModeChangeEvent struct {
	Status uint8
	ConnectionHandle uint16
	CurrentMode uint8
	Interval uint16
}

func (o *ModeChangeEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.CurrentMode = uint8(r.GetOne())
	o.Interval = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// ModeChangeEventCallbackType is the type of the callback function for ModeChangeEvent.
type ModeChangeEventCallbackType func(*ModeChangeEvent) *ModeChangeEvent

// SetModeChangeEventCallback configures the callback for ModeChangeEvent. Passing nil will disable the callback.
func (e *EventHandler) SetModeChangeEventCallback(cb ModeChangeEventCallbackType) error {
	e.cbMutex.Lock()
	e.modeChangeEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(19, cb != nil);
}
// ReturnLinkKeysEvent represents the event specified in Section 7.7.21
type ReturnLinkKeysEvent struct {
	NumKeys uint8
	BDADDR []bleutil.MacAddr
	LinkKey [][16]byte
}

func (o *ReturnLinkKeysEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.NumKeys = uint8(r.GetOne())
	if cap(o.BDADDR) < int(o.NumKeys) {
		o.BDADDR = make([]bleutil.MacAddr, 0, int(o.NumKeys))
	}
	o.BDADDR = o.BDADDR[:int(o.NumKeys)]
	for j:=0; j<int(o.NumKeys); j++ {
		o.BDADDR[j].Decode(r.Get(6))
	}
	if cap(o.LinkKey) < int(o.NumKeys) {
		o.LinkKey = make([][16]byte, 0, int(o.NumKeys))
	}
	o.LinkKey = o.LinkKey[:int(o.NumKeys)]
	for j:=0; j<int(o.NumKeys); j++ {
		copy(o.LinkKey[j][:], r.Get(16))
	}
	return r.Valid()
}

// ReturnLinkKeysEventCallbackType is the type of the callback function for ReturnLinkKeysEvent.
type ReturnLinkKeysEventCallbackType func(*ReturnLinkKeysEvent) *ReturnLinkKeysEvent

// SetReturnLinkKeysEventCallback configures the callback for ReturnLinkKeysEvent. Passing nil will disable the callback.
func (e *EventHandler) SetReturnLinkKeysEventCallback(cb ReturnLinkKeysEventCallbackType) error {
	e.cbMutex.Lock()
	e.returnLinkKeysEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(20, cb != nil);
}
// PINCodeRequestEvent represents the event specified in Section 7.7.22
type PINCodeRequestEvent struct {
	BDADDR bleutil.MacAddr
}

func (o *PINCodeRequestEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	return r.Valid()
}

// PINCodeRequestEventCallbackType is the type of the callback function for PINCodeRequestEvent.
type PINCodeRequestEventCallbackType func(*PINCodeRequestEvent) *PINCodeRequestEvent

// SetPINCodeRequestEventCallback configures the callback for PINCodeRequestEvent. Passing nil will disable the callback.
func (e *EventHandler) SetPINCodeRequestEventCallback(cb PINCodeRequestEventCallbackType) error {
	e.cbMutex.Lock()
	e.pINCodeRequestEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(21, cb != nil);
}
// LinkKeyRequestEvent represents the event specified in Section 7.7.23
type LinkKeyRequestEvent struct {
	BDADDR bleutil.MacAddr
}

func (o *LinkKeyRequestEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	return r.Valid()
}

// LinkKeyRequestEventCallbackType is the type of the callback function for LinkKeyRequestEvent.
type LinkKeyRequestEventCallbackType func(*LinkKeyRequestEvent) *LinkKeyRequestEvent

// SetLinkKeyRequestEventCallback configures the callback for LinkKeyRequestEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLinkKeyRequestEventCallback(cb LinkKeyRequestEventCallbackType) error {
	e.cbMutex.Lock()
	e.linkKeyRequestEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(22, cb != nil);
}
// LinkKeyNotificationEvent represents the event specified in Section 7.7.24
type LinkKeyNotificationEvent struct {
	BDADDR bleutil.MacAddr
	LinkKey [16]byte
	KeyType uint8
}

func (o *LinkKeyNotificationEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	copy(o.LinkKey[:], r.Get(16))
	o.KeyType = uint8(r.GetOne())
	return r.Valid()
}

// LinkKeyNotificationEventCallbackType is the type of the callback function for LinkKeyNotificationEvent.
type LinkKeyNotificationEventCallbackType func(*LinkKeyNotificationEvent) *LinkKeyNotificationEvent

// SetLinkKeyNotificationEventCallback configures the callback for LinkKeyNotificationEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLinkKeyNotificationEventCallback(cb LinkKeyNotificationEventCallbackType) error {
	e.cbMutex.Lock()
	e.linkKeyNotificationEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(23, cb != nil);
}
// LoopbackCommandEvent represents the event specified in Section 7.7.25
type LoopbackCommandEvent struct {
	HCICommandPacket []byte
}

func (o *LoopbackCommandEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.HCICommandPacket = append(o.HCICommandPacket[:0], r.GetRemainder()...)
	return r.Valid()
}

// LoopbackCommandEventCallbackType is the type of the callback function for LoopbackCommandEvent.
type LoopbackCommandEventCallbackType func(*LoopbackCommandEvent) *LoopbackCommandEvent

// SetLoopbackCommandEventCallback configures the callback for LoopbackCommandEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLoopbackCommandEventCallback(cb LoopbackCommandEventCallbackType) error {
	e.cbMutex.Lock()
	e.loopbackCommandEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(24, cb != nil);
}
// DataBufferOverflowEvent represents the event specified in Section 7.7.26
type DataBufferOverflowEvent struct {
	LinkType uint8
}

func (o *DataBufferOverflowEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.LinkType = uint8(r.GetOne())
	return r.Valid()
}

// DataBufferOverflowEventCallbackType is the type of the callback function for DataBufferOverflowEvent.
type DataBufferOverflowEventCallbackType func(*DataBufferOverflowEvent) *DataBufferOverflowEvent

// SetDataBufferOverflowEventCallback configures the callback for DataBufferOverflowEvent. Passing nil will disable the callback.
func (e *EventHandler) SetDataBufferOverflowEventCallback(cb DataBufferOverflowEventCallbackType) error {
	e.cbMutex.Lock()
	e.dataBufferOverflowEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(25, cb != nil);
}
// MaxSlotsChangeEvent represents the event specified in Section 7.7.27
type MaxSlotsChangeEvent struct {
	ConnectionHandle uint16
	LMPMaxSlots uint8
}

func (o *MaxSlotsChangeEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LMPMaxSlots = uint8(r.GetOne())
	return r.Valid()
}

// MaxSlotsChangeEventCallbackType is the type of the callback function for MaxSlotsChangeEvent.
type MaxSlotsChangeEventCallbackType func(*MaxSlotsChangeEvent) *MaxSlotsChangeEvent

// SetMaxSlotsChangeEventCallback configures the callback for MaxSlotsChangeEvent. Passing nil will disable the callback.
func (e *EventHandler) SetMaxSlotsChangeEventCallback(cb MaxSlotsChangeEventCallbackType) error {
	e.cbMutex.Lock()
	e.maxSlotsChangeEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(26, cb != nil);
}
// ReadClockOffsetCompleteEvent represents the event specified in Section 7.7.28
type ReadClockOffsetCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
	ClockOffset uint16
}

func (o *ReadClockOffsetCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ClockOffset = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// ReadClockOffsetCompleteEventCallbackType is the type of the callback function for ReadClockOffsetCompleteEvent.
type ReadClockOffsetCompleteEventCallbackType func(*ReadClockOffsetCompleteEvent) *ReadClockOffsetCompleteEvent

// SetReadClockOffsetCompleteEventCallback configures the callback for ReadClockOffsetCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetReadClockOffsetCompleteEventCallback(cb ReadClockOffsetCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.readClockOffsetCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(27, cb != nil);
}
// ConnectionPacketTypeChangedEvent represents the event specified in Section 7.7.29
type ConnectionPacketTypeChangedEvent struct {
	Status uint8
	ConnectionHandle uint16
	PacketType uint16
}

func (o *ConnectionPacketTypeChangedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PacketType = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// ConnectionPacketTypeChangedEventCallbackType is the type of the callback function for ConnectionPacketTypeChangedEvent.
type ConnectionPacketTypeChangedEventCallbackType func(*ConnectionPacketTypeChangedEvent) *ConnectionPacketTypeChangedEvent

// SetConnectionPacketTypeChangedEventCallback configures the callback for ConnectionPacketTypeChangedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetConnectionPacketTypeChangedEventCallback(cb ConnectionPacketTypeChangedEventCallbackType) error {
	e.cbMutex.Lock()
	e.connectionPacketTypeChangedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(28, cb != nil);
}
// QoSViolationEvent represents the event specified in Section 7.7.30
type QoSViolationEvent struct {
	Handle uint16
}

func (o *QoSViolationEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// QoSViolationEventCallbackType is the type of the callback function for QoSViolationEvent.
type QoSViolationEventCallbackType func(*QoSViolationEvent) *QoSViolationEvent

// SetQoSViolationEventCallback configures the callback for QoSViolationEvent. Passing nil will disable the callback.
func (e *EventHandler) SetQoSViolationEventCallback(cb QoSViolationEventCallbackType) error {
	e.cbMutex.Lock()
	e.qoSViolationEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(29, cb != nil);
}
// PageScanRepetitionModeChangeEvent represents the event specified in Section 7.7.31
type PageScanRepetitionModeChangeEvent struct {
	BDADDR bleutil.MacAddr
	PageScanRepetitionMode uint8
}

func (o *PageScanRepetitionModeChangeEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	o.PageScanRepetitionMode = uint8(r.GetOne())
	return r.Valid()
}

// PageScanRepetitionModeChangeEventCallbackType is the type of the callback function for PageScanRepetitionModeChangeEvent.
type PageScanRepetitionModeChangeEventCallbackType func(*PageScanRepetitionModeChangeEvent) *PageScanRepetitionModeChangeEvent

// SetPageScanRepetitionModeChangeEventCallback configures the callback for PageScanRepetitionModeChangeEvent. Passing nil will disable the callback.
func (e *EventHandler) SetPageScanRepetitionModeChangeEventCallback(cb PageScanRepetitionModeChangeEventCallbackType) error {
	e.cbMutex.Lock()
	e.pageScanRepetitionModeChangeEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(31, cb != nil);
}
// FlowSpecificationCompleteEvent represents the event specified in Section 7.7.32
type FlowSpecificationCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
	Unused uint8
	FlowDirection uint8
	ServiceType uint8
	TokenRate uint32
	TokenBucketSize uint32
	PeakBandwidth uint32
	AccessLatency uint32
}

func (o *FlowSpecificationCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Unused = uint8(r.GetOne())
	o.FlowDirection = uint8(r.GetOne())
	o.ServiceType = uint8(r.GetOne())
	o.TokenRate = binary.LittleEndian.Uint32(r.Get(4))
	o.TokenBucketSize = binary.LittleEndian.Uint32(r.Get(4))
	o.PeakBandwidth = binary.LittleEndian.Uint32(r.Get(4))
	o.AccessLatency = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// FlowSpecificationCompleteEventCallbackType is the type of the callback function for FlowSpecificationCompleteEvent.
type FlowSpecificationCompleteEventCallbackType func(*FlowSpecificationCompleteEvent) *FlowSpecificationCompleteEvent

// SetFlowSpecificationCompleteEventCallback configures the callback for FlowSpecificationCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetFlowSpecificationCompleteEventCallback(cb FlowSpecificationCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.flowSpecificationCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(32, cb != nil);
}
// InquiryResultwithRSSIEvent represents the event specified in Section 7.7.33
type InquiryResultwithRSSIEvent struct {
	NumResponses uint8
	BDADDR []bleutil.MacAddr
	PageScanRepetitionMode []uint8
	Reserved []uint8
	ClassOfDevice []uint32
	ClockOffset []uint16
	RSSI []uint8
}

func (o *InquiryResultwithRSSIEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.NumResponses = uint8(r.GetOne())
	if cap(o.BDADDR) < int(o.NumResponses) {
		o.BDADDR = make([]bleutil.MacAddr, 0, int(o.NumResponses))
	}
	o.BDADDR = o.BDADDR[:int(o.NumResponses)]
	for j:=0; j<int(o.NumResponses); j++ {
		o.BDADDR[j].Decode(r.Get(6))
	}
	if cap(o.PageScanRepetitionMode) < int(o.NumResponses) {
		o.PageScanRepetitionMode = make([]uint8, 0, int(o.NumResponses))
	}
	o.PageScanRepetitionMode = o.PageScanRepetitionMode[:int(o.NumResponses)]
	for j:=0; j<int(o.NumResponses); j++ {
		o.PageScanRepetitionMode[j] = uint8(r.GetOne())
	}
	if cap(o.Reserved) < int(o.NumResponses) {
		o.Reserved = make([]uint8, 0, int(o.NumResponses))
	}
	o.Reserved = o.Reserved[:int(o.NumResponses)]
	for j:=0; j<int(o.NumResponses); j++ {
		o.Reserved[j] = uint8(r.GetOne())
	}
	if cap(o.ClassOfDevice) < int(o.NumResponses) {
		o.ClassOfDevice = make([]uint32, 0, int(o.NumResponses))
	}
	o.ClassOfDevice = o.ClassOfDevice[:int(o.NumResponses)]
	for j:=0; j<int(o.NumResponses); j++ {
		o.ClassOfDevice[j] = bleutil.DecodeUint24(r.Get(3))
	}
	if cap(o.ClockOffset) < int(o.NumResponses) {
		o.ClockOffset = make([]uint16, 0, int(o.NumResponses))
	}
	o.ClockOffset = o.ClockOffset[:int(o.NumResponses)]
	for j:=0; j<int(o.NumResponses); j++ {
		o.ClockOffset[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.RSSI) < int(o.NumResponses) {
		o.RSSI = make([]uint8, 0, int(o.NumResponses))
	}
	o.RSSI = o.RSSI[:int(o.NumResponses)]
	for j:=0; j<int(o.NumResponses); j++ {
		o.RSSI[j] = uint8(r.GetOne())
	}
	return r.Valid()
}

// InquiryResultwithRSSIEventCallbackType is the type of the callback function for InquiryResultwithRSSIEvent.
type InquiryResultwithRSSIEventCallbackType func(*InquiryResultwithRSSIEvent) *InquiryResultwithRSSIEvent

// SetInquiryResultwithRSSIEventCallback configures the callback for InquiryResultwithRSSIEvent. Passing nil will disable the callback.
func (e *EventHandler) SetInquiryResultwithRSSIEventCallback(cb InquiryResultwithRSSIEventCallbackType) error {
	e.cbMutex.Lock()
	e.inquiryResultwithRSSIEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(33, cb != nil);
}
// ReadRemoteExtendedFeaturesCompleteEvent represents the event specified in Section 7.7.34
type ReadRemoteExtendedFeaturesCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
	PageNumber uint8
	MaximumPageNumber uint8
	ExtendedLMPFeatures uint64
}

func (o *ReadRemoteExtendedFeaturesCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PageNumber = uint8(r.GetOne())
	o.MaximumPageNumber = uint8(r.GetOne())
	o.ExtendedLMPFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// ReadRemoteExtendedFeaturesCompleteEventCallbackType is the type of the callback function for ReadRemoteExtendedFeaturesCompleteEvent.
type ReadRemoteExtendedFeaturesCompleteEventCallbackType func(*ReadRemoteExtendedFeaturesCompleteEvent) *ReadRemoteExtendedFeaturesCompleteEvent

// SetReadRemoteExtendedFeaturesCompleteEventCallback configures the callback for ReadRemoteExtendedFeaturesCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetReadRemoteExtendedFeaturesCompleteEventCallback(cb ReadRemoteExtendedFeaturesCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.readRemoteExtendedFeaturesCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(34, cb != nil);
}
// SynchronousConnectionCompleteEvent represents the event specified in Section 7.7.35
type SynchronousConnectionCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
	BDADDR bleutil.MacAddr
	LinkType uint8
	TransmissionInterval uint8
	RetransmissionWindow uint8
	RXPacketLength uint16
	TXPacketLength uint16
	AirMode uint8
}

func (o *SynchronousConnectionCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.BDADDR.Decode(r.Get(6))
	o.LinkType = uint8(r.GetOne())
	o.TransmissionInterval = uint8(r.GetOne())
	o.RetransmissionWindow = uint8(r.GetOne())
	o.RXPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.TXPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.AirMode = uint8(r.GetOne())
	return r.Valid()
}

// SynchronousConnectionCompleteEventCallbackType is the type of the callback function for SynchronousConnectionCompleteEvent.
type SynchronousConnectionCompleteEventCallbackType func(*SynchronousConnectionCompleteEvent) *SynchronousConnectionCompleteEvent

// SetSynchronousConnectionCompleteEventCallback configures the callback for SynchronousConnectionCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetSynchronousConnectionCompleteEventCallback(cb SynchronousConnectionCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.synchronousConnectionCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(43, cb != nil);
}
// SynchronousConnectionChangedEvent represents the event specified in Section 7.7.36
type SynchronousConnectionChangedEvent struct {
	Status uint8
	ConnectionHandle uint16
	TransmissionInterval uint8
	RetransmissionWindow uint8
	RXPacketLength uint16
	TXPacketLength uint16
}

func (o *SynchronousConnectionChangedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.TransmissionInterval = uint8(r.GetOne())
	o.RetransmissionWindow = uint8(r.GetOne())
	o.RXPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.TXPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// SynchronousConnectionChangedEventCallbackType is the type of the callback function for SynchronousConnectionChangedEvent.
type SynchronousConnectionChangedEventCallbackType func(*SynchronousConnectionChangedEvent) *SynchronousConnectionChangedEvent

// SetSynchronousConnectionChangedEventCallback configures the callback for SynchronousConnectionChangedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetSynchronousConnectionChangedEventCallback(cb SynchronousConnectionChangedEventCallbackType) error {
	e.cbMutex.Lock()
	e.synchronousConnectionChangedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(44, cb != nil);
}
// SniffSubratingEvent represents the event specified in Section 7.7.37
type SniffSubratingEvent struct {
	Status uint8
	ConnectionHandle uint16
	MaxTXLatency uint16
	MaxRXLatency uint16
	MinRemoteTimeout uint16
	MinLocalTimeout uint16
}

func (o *SniffSubratingEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxTXLatency = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxRXLatency = binary.LittleEndian.Uint16(r.Get(2))
	o.MinRemoteTimeout = binary.LittleEndian.Uint16(r.Get(2))
	o.MinLocalTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// SniffSubratingEventCallbackType is the type of the callback function for SniffSubratingEvent.
type SniffSubratingEventCallbackType func(*SniffSubratingEvent) *SniffSubratingEvent

// SetSniffSubratingEventCallback configures the callback for SniffSubratingEvent. Passing nil will disable the callback.
func (e *EventHandler) SetSniffSubratingEventCallback(cb SniffSubratingEventCallbackType) error {
	e.cbMutex.Lock()
	e.sniffSubratingEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(45, cb != nil);
}
// ExtendedInquiryResultEvent represents the event specified in Section 7.7.38
type ExtendedInquiryResultEvent struct {
	NumResponses uint8
	BDADDR bleutil.MacAddr
	PageScanRepetitionMode uint8
	Reserved uint8
	ClassOfDevice uint32
	ClockOffset uint16
	RSSI uint8
	ExtendedInquiryResponse [240]byte
}

func (o *ExtendedInquiryResultEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.NumResponses = uint8(r.GetOne())
	o.BDADDR.Decode(r.Get(6))
	o.PageScanRepetitionMode = uint8(r.GetOne())
	o.Reserved = uint8(r.GetOne())
	o.ClassOfDevice = bleutil.DecodeUint24(r.Get(3))
	o.ClockOffset = binary.LittleEndian.Uint16(r.Get(2))
	o.RSSI = uint8(r.GetOne())
	copy(o.ExtendedInquiryResponse[:], r.Get(240))
	return r.Valid()
}

// ExtendedInquiryResultEventCallbackType is the type of the callback function for ExtendedInquiryResultEvent.
type ExtendedInquiryResultEventCallbackType func(*ExtendedInquiryResultEvent) *ExtendedInquiryResultEvent

// SetExtendedInquiryResultEventCallback configures the callback for ExtendedInquiryResultEvent. Passing nil will disable the callback.
func (e *EventHandler) SetExtendedInquiryResultEventCallback(cb ExtendedInquiryResultEventCallbackType) error {
	e.cbMutex.Lock()
	e.extendedInquiryResultEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(46, cb != nil);
}
// EncryptionKeyRefreshCompleteEvent represents the event specified in Section 7.7.39
type EncryptionKeyRefreshCompleteEvent struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *EncryptionKeyRefreshCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// EncryptionKeyRefreshCompleteEventCallbackType is the type of the callback function for EncryptionKeyRefreshCompleteEvent.
type EncryptionKeyRefreshCompleteEventCallbackType func(*EncryptionKeyRefreshCompleteEvent) *EncryptionKeyRefreshCompleteEvent

// SetEncryptionKeyRefreshCompleteEventCallback configures the callback for EncryptionKeyRefreshCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetEncryptionKeyRefreshCompleteEventCallback(cb EncryptionKeyRefreshCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.encryptionKeyRefreshCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(47, cb != nil);
}
// IOCapabilityRequestEvent represents the event specified in Section 7.7.40
type IOCapabilityRequestEvent struct {
	BDADDR bleutil.MacAddr
}

func (o *IOCapabilityRequestEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	return r.Valid()
}

// IOCapabilityRequestEventCallbackType is the type of the callback function for IOCapabilityRequestEvent.
type IOCapabilityRequestEventCallbackType func(*IOCapabilityRequestEvent) *IOCapabilityRequestEvent

// SetIOCapabilityRequestEventCallback configures the callback for IOCapabilityRequestEvent. Passing nil will disable the callback.
func (e *EventHandler) SetIOCapabilityRequestEventCallback(cb IOCapabilityRequestEventCallbackType) error {
	e.cbMutex.Lock()
	e.iOCapabilityRequestEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(48, cb != nil);
}
// IOCapabilityResponseEvent represents the event specified in Section 7.7.41
type IOCapabilityResponseEvent struct {
	BDADDR bleutil.MacAddr
	IOCapability uint8
	OOBDataPresent uint8
	AuthenticationRequirements uint8
}

func (o *IOCapabilityResponseEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	o.IOCapability = uint8(r.GetOne())
	o.OOBDataPresent = uint8(r.GetOne())
	o.AuthenticationRequirements = uint8(r.GetOne())
	return r.Valid()
}

// IOCapabilityResponseEventCallbackType is the type of the callback function for IOCapabilityResponseEvent.
type IOCapabilityResponseEventCallbackType func(*IOCapabilityResponseEvent) *IOCapabilityResponseEvent

// SetIOCapabilityResponseEventCallback configures the callback for IOCapabilityResponseEvent. Passing nil will disable the callback.
func (e *EventHandler) SetIOCapabilityResponseEventCallback(cb IOCapabilityResponseEventCallbackType) error {
	e.cbMutex.Lock()
	e.iOCapabilityResponseEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(49, cb != nil);
}
// UserConfirmationRequestEvent represents the event specified in Section 7.7.42
type UserConfirmationRequestEvent struct {
	BDADDR bleutil.MacAddr
	NumericValue uint32
}

func (o *UserConfirmationRequestEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	o.NumericValue = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// UserConfirmationRequestEventCallbackType is the type of the callback function for UserConfirmationRequestEvent.
type UserConfirmationRequestEventCallbackType func(*UserConfirmationRequestEvent) *UserConfirmationRequestEvent

// SetUserConfirmationRequestEventCallback configures the callback for UserConfirmationRequestEvent. Passing nil will disable the callback.
func (e *EventHandler) SetUserConfirmationRequestEventCallback(cb UserConfirmationRequestEventCallbackType) error {
	e.cbMutex.Lock()
	e.userConfirmationRequestEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(50, cb != nil);
}
// UserPasskeyRequestEvent represents the event specified in Section 7.7.43
type UserPasskeyRequestEvent struct {
	BDADDR bleutil.MacAddr
}

func (o *UserPasskeyRequestEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	return r.Valid()
}

// UserPasskeyRequestEventCallbackType is the type of the callback function for UserPasskeyRequestEvent.
type UserPasskeyRequestEventCallbackType func(*UserPasskeyRequestEvent) *UserPasskeyRequestEvent

// SetUserPasskeyRequestEventCallback configures the callback for UserPasskeyRequestEvent. Passing nil will disable the callback.
func (e *EventHandler) SetUserPasskeyRequestEventCallback(cb UserPasskeyRequestEventCallbackType) error {
	e.cbMutex.Lock()
	e.userPasskeyRequestEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(51, cb != nil);
}
// RemoteOOBDataRequestEvent represents the event specified in Section 7.7.44
type RemoteOOBDataRequestEvent struct {
	BDADDR bleutil.MacAddr
}

func (o *RemoteOOBDataRequestEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	return r.Valid()
}

// RemoteOOBDataRequestEventCallbackType is the type of the callback function for RemoteOOBDataRequestEvent.
type RemoteOOBDataRequestEventCallbackType func(*RemoteOOBDataRequestEvent) *RemoteOOBDataRequestEvent

// SetRemoteOOBDataRequestEventCallback configures the callback for RemoteOOBDataRequestEvent. Passing nil will disable the callback.
func (e *EventHandler) SetRemoteOOBDataRequestEventCallback(cb RemoteOOBDataRequestEventCallbackType) error {
	e.cbMutex.Lock()
	e.remoteOOBDataRequestEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(52, cb != nil);
}
// SimplePairingCompleteEvent represents the event specified in Section 7.7.45
type SimplePairingCompleteEvent struct {
	Status uint8
	BDADDR bleutil.MacAddr
}

func (o *SimplePairingCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.BDADDR.Decode(r.Get(6))
	return r.Valid()
}

// SimplePairingCompleteEventCallbackType is the type of the callback function for SimplePairingCompleteEvent.
type SimplePairingCompleteEventCallbackType func(*SimplePairingCompleteEvent) *SimplePairingCompleteEvent

// SetSimplePairingCompleteEventCallback configures the callback for SimplePairingCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetSimplePairingCompleteEventCallback(cb SimplePairingCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.simplePairingCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(53, cb != nil);
}
// LinkSupervisionTimeoutChangedEvent represents the event specified in Section 7.7.46
type LinkSupervisionTimeoutChangedEvent struct {
	ConnectionHandle uint16
	LinkSupervisionTimeout uint16
}

func (o *LinkSupervisionTimeoutChangedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LinkSupervisionTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LinkSupervisionTimeoutChangedEventCallbackType is the type of the callback function for LinkSupervisionTimeoutChangedEvent.
type LinkSupervisionTimeoutChangedEventCallbackType func(*LinkSupervisionTimeoutChangedEvent) *LinkSupervisionTimeoutChangedEvent

// SetLinkSupervisionTimeoutChangedEventCallback configures the callback for LinkSupervisionTimeoutChangedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLinkSupervisionTimeoutChangedEventCallback(cb LinkSupervisionTimeoutChangedEventCallbackType) error {
	e.cbMutex.Lock()
	e.linkSupervisionTimeoutChangedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(55, cb != nil);
}
// EnhancedFlushCompleteEvent represents the event specified in Section 7.7.47
type EnhancedFlushCompleteEvent struct {
	Handle uint16
}

func (o *EnhancedFlushCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// EnhancedFlushCompleteEventCallbackType is the type of the callback function for EnhancedFlushCompleteEvent.
type EnhancedFlushCompleteEventCallbackType func(*EnhancedFlushCompleteEvent) *EnhancedFlushCompleteEvent

// SetEnhancedFlushCompleteEventCallback configures the callback for EnhancedFlushCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetEnhancedFlushCompleteEventCallback(cb EnhancedFlushCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.enhancedFlushCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(56, cb != nil);
}
// UserPasskeyNotificationEvent represents the event specified in Section 7.7.48
type UserPasskeyNotificationEvent struct {
	BDADDR bleutil.MacAddr
	Passkey uint32
}

func (o *UserPasskeyNotificationEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	o.Passkey = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// UserPasskeyNotificationEventCallbackType is the type of the callback function for UserPasskeyNotificationEvent.
type UserPasskeyNotificationEventCallbackType func(*UserPasskeyNotificationEvent) *UserPasskeyNotificationEvent

// SetUserPasskeyNotificationEventCallback configures the callback for UserPasskeyNotificationEvent. Passing nil will disable the callback.
func (e *EventHandler) SetUserPasskeyNotificationEventCallback(cb UserPasskeyNotificationEventCallbackType) error {
	e.cbMutex.Lock()
	e.userPasskeyNotificationEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(58, cb != nil);
}
// KeypressNotificationEvent represents the event specified in Section 7.7.49
type KeypressNotificationEvent struct {
	BDADDR bleutil.MacAddr
	NotificationType uint8
}

func (o *KeypressNotificationEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	o.NotificationType = uint8(r.GetOne())
	return r.Valid()
}

// KeypressNotificationEventCallbackType is the type of the callback function for KeypressNotificationEvent.
type KeypressNotificationEventCallbackType func(*KeypressNotificationEvent) *KeypressNotificationEvent

// SetKeypressNotificationEventCallback configures the callback for KeypressNotificationEvent. Passing nil will disable the callback.
func (e *EventHandler) SetKeypressNotificationEventCallback(cb KeypressNotificationEventCallbackType) error {
	e.cbMutex.Lock()
	e.keypressNotificationEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(59, cb != nil);
}
// RemoteHostSupportedFeaturesNotificationEvent represents the event specified in Section 7.7.50
type RemoteHostSupportedFeaturesNotificationEvent struct {
	BDADDR bleutil.MacAddr
	HostSupportedFeatures uint64
}

func (o *RemoteHostSupportedFeaturesNotificationEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	o.HostSupportedFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// RemoteHostSupportedFeaturesNotificationEventCallbackType is the type of the callback function for RemoteHostSupportedFeaturesNotificationEvent.
type RemoteHostSupportedFeaturesNotificationEventCallbackType func(*RemoteHostSupportedFeaturesNotificationEvent) *RemoteHostSupportedFeaturesNotificationEvent

// SetRemoteHostSupportedFeaturesNotificationEventCallback configures the callback for RemoteHostSupportedFeaturesNotificationEvent. Passing nil will disable the callback.
func (e *EventHandler) SetRemoteHostSupportedFeaturesNotificationEventCallback(cb RemoteHostSupportedFeaturesNotificationEventCallbackType) error {
	e.cbMutex.Lock()
	e.remoteHostSupportedFeaturesNotificationEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(60, cb != nil);
}
// PhysicalLinkCompleteEvent represents the event specified in Section 7.7.51
type PhysicalLinkCompleteEvent struct {
	Status uint8
	PhysicalLinkHandle uint8
}

func (o *PhysicalLinkCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.PhysicalLinkHandle = uint8(r.GetOne())
	return r.Valid()
}

// PhysicalLinkCompleteEventCallbackType is the type of the callback function for PhysicalLinkCompleteEvent.
type PhysicalLinkCompleteEventCallbackType func(*PhysicalLinkCompleteEvent) *PhysicalLinkCompleteEvent

// SetPhysicalLinkCompleteEventCallback configures the callback for PhysicalLinkCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetPhysicalLinkCompleteEventCallback(cb PhysicalLinkCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.physicalLinkCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(100, cb != nil);
}
// ChannelSelectedEvent represents the event specified in Section 7.7.52
type ChannelSelectedEvent struct {
	PhysicalLinkHandle uint8
}

func (o *ChannelSelectedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.PhysicalLinkHandle = uint8(r.GetOne())
	return r.Valid()
}

// ChannelSelectedEventCallbackType is the type of the callback function for ChannelSelectedEvent.
type ChannelSelectedEventCallbackType func(*ChannelSelectedEvent) *ChannelSelectedEvent

// SetChannelSelectedEventCallback configures the callback for ChannelSelectedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetChannelSelectedEventCallback(cb ChannelSelectedEventCallbackType) error {
	e.cbMutex.Lock()
	e.channelSelectedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(101, cb != nil);
}
// DisconnectionPhysicalLinkCompleteEvent represents the event specified in Section 7.7.53
type DisconnectionPhysicalLinkCompleteEvent struct {
	Status uint8
	PhysicalLinkHandle uint8
	Reason uint8
}

func (o *DisconnectionPhysicalLinkCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.PhysicalLinkHandle = uint8(r.GetOne())
	o.Reason = uint8(r.GetOne())
	return r.Valid()
}

// DisconnectionPhysicalLinkCompleteEventCallbackType is the type of the callback function for DisconnectionPhysicalLinkCompleteEvent.
type DisconnectionPhysicalLinkCompleteEventCallbackType func(*DisconnectionPhysicalLinkCompleteEvent) *DisconnectionPhysicalLinkCompleteEvent

// SetDisconnectionPhysicalLinkCompleteEventCallback configures the callback for DisconnectionPhysicalLinkCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetDisconnectionPhysicalLinkCompleteEventCallback(cb DisconnectionPhysicalLinkCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.disconnectionPhysicalLinkCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(102, cb != nil);
}
// PhysicalLinkLossEarlyWarningEvent represents the event specified in Section 7.7.54
type PhysicalLinkLossEarlyWarningEvent struct {
	PhysicalLinkHandle uint8
	LinkLossReason uint8
}

func (o *PhysicalLinkLossEarlyWarningEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.PhysicalLinkHandle = uint8(r.GetOne())
	o.LinkLossReason = uint8(r.GetOne())
	return r.Valid()
}

// PhysicalLinkLossEarlyWarningEventCallbackType is the type of the callback function for PhysicalLinkLossEarlyWarningEvent.
type PhysicalLinkLossEarlyWarningEventCallbackType func(*PhysicalLinkLossEarlyWarningEvent) *PhysicalLinkLossEarlyWarningEvent

// SetPhysicalLinkLossEarlyWarningEventCallback configures the callback for PhysicalLinkLossEarlyWarningEvent. Passing nil will disable the callback.
func (e *EventHandler) SetPhysicalLinkLossEarlyWarningEventCallback(cb PhysicalLinkLossEarlyWarningEventCallbackType) error {
	e.cbMutex.Lock()
	e.physicalLinkLossEarlyWarningEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(103, cb != nil);
}
// PhysicalLinkRecoveryEvent represents the event specified in Section 7.7.55
type PhysicalLinkRecoveryEvent struct {
	PhysicalLinkHandle uint8
}

func (o *PhysicalLinkRecoveryEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.PhysicalLinkHandle = uint8(r.GetOne())
	return r.Valid()
}

// PhysicalLinkRecoveryEventCallbackType is the type of the callback function for PhysicalLinkRecoveryEvent.
type PhysicalLinkRecoveryEventCallbackType func(*PhysicalLinkRecoveryEvent) *PhysicalLinkRecoveryEvent

// SetPhysicalLinkRecoveryEventCallback configures the callback for PhysicalLinkRecoveryEvent. Passing nil will disable the callback.
func (e *EventHandler) SetPhysicalLinkRecoveryEventCallback(cb PhysicalLinkRecoveryEventCallbackType) error {
	e.cbMutex.Lock()
	e.physicalLinkRecoveryEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(104, cb != nil);
}
// LogicalLinkCompleteEvent represents the event specified in Section 7.7.56
type LogicalLinkCompleteEvent struct {
	Status uint8
	LogicalLinkHandle uint16
	PhysicalLinkHandle uint8
	TXFlowSpecID uint8
}

func (o *LogicalLinkCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.LogicalLinkHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PhysicalLinkHandle = uint8(r.GetOne())
	o.TXFlowSpecID = uint8(r.GetOne())
	return r.Valid()
}

// LogicalLinkCompleteEventCallbackType is the type of the callback function for LogicalLinkCompleteEvent.
type LogicalLinkCompleteEventCallbackType func(*LogicalLinkCompleteEvent) *LogicalLinkCompleteEvent

// SetLogicalLinkCompleteEventCallback configures the callback for LogicalLinkCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLogicalLinkCompleteEventCallback(cb LogicalLinkCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.logicalLinkCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(105, cb != nil);
}
// DisconnectionLogicalLinkCompleteEvent represents the event specified in Section 7.7.57
type DisconnectionLogicalLinkCompleteEvent struct {
	Status uint8
	LogicalLinkHandle uint16
	Reason uint8
}

func (o *DisconnectionLogicalLinkCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.LogicalLinkHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Reason = uint8(r.GetOne())
	return r.Valid()
}

// DisconnectionLogicalLinkCompleteEventCallbackType is the type of the callback function for DisconnectionLogicalLinkCompleteEvent.
type DisconnectionLogicalLinkCompleteEventCallbackType func(*DisconnectionLogicalLinkCompleteEvent) *DisconnectionLogicalLinkCompleteEvent

// SetDisconnectionLogicalLinkCompleteEventCallback configures the callback for DisconnectionLogicalLinkCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetDisconnectionLogicalLinkCompleteEventCallback(cb DisconnectionLogicalLinkCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.disconnectionLogicalLinkCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(106, cb != nil);
}
// FlowSpecModifyCompleteEvent represents the event specified in Section 7.7.58
type FlowSpecModifyCompleteEvent struct {
	Status uint8
	Handle uint16
}

func (o *FlowSpecModifyCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// FlowSpecModifyCompleteEventCallbackType is the type of the callback function for FlowSpecModifyCompleteEvent.
type FlowSpecModifyCompleteEventCallbackType func(*FlowSpecModifyCompleteEvent) *FlowSpecModifyCompleteEvent

// SetFlowSpecModifyCompleteEventCallback configures the callback for FlowSpecModifyCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetFlowSpecModifyCompleteEventCallback(cb FlowSpecModifyCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.flowSpecModifyCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(107, cb != nil);
}
// NumberOfCompletedDataBlocksEvent represents the event specified in Section 7.7.59
type NumberOfCompletedDataBlocksEvent struct {
	TotalNumDataBlocks uint16
	NumHandles uint8
	Handle []uint16
	NumCompletedPackets []uint16
	NumCompletedBlocks []uint16
}

func (o *NumberOfCompletedDataBlocksEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.TotalNumDataBlocks = binary.LittleEndian.Uint16(r.Get(2))
	o.NumHandles = uint8(r.GetOne())
	if cap(o.Handle) < int(o.NumHandles) {
		o.Handle = make([]uint16, 0, int(o.NumHandles))
	}
	o.Handle = o.Handle[:int(o.NumHandles)]
	for j:=0; j<int(o.NumHandles); j++ {
		o.Handle[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.NumCompletedPackets) < int(o.NumHandles) {
		o.NumCompletedPackets = make([]uint16, 0, int(o.NumHandles))
	}
	o.NumCompletedPackets = o.NumCompletedPackets[:int(o.NumHandles)]
	for j:=0; j<int(o.NumHandles); j++ {
		o.NumCompletedPackets[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.NumCompletedBlocks) < int(o.NumHandles) {
		o.NumCompletedBlocks = make([]uint16, 0, int(o.NumHandles))
	}
	o.NumCompletedBlocks = o.NumCompletedBlocks[:int(o.NumHandles)]
	for j:=0; j<int(o.NumHandles); j++ {
		o.NumCompletedBlocks[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// NumberOfCompletedDataBlocksEventCallbackType is the type of the callback function for NumberOfCompletedDataBlocksEvent.
type NumberOfCompletedDataBlocksEventCallbackType func(*NumberOfCompletedDataBlocksEvent) *NumberOfCompletedDataBlocksEvent

// SetNumberOfCompletedDataBlocksEventCallback configures the callback for NumberOfCompletedDataBlocksEvent. Passing nil will disable the callback.
func (e *EventHandler) SetNumberOfCompletedDataBlocksEventCallback(cb NumberOfCompletedDataBlocksEventCallbackType) error {
	e.cbMutex.Lock()
	e.numberOfCompletedDataBlocksEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(108, cb != nil);
}
// ShortRangeModeChangeCompleteEvent represents the event specified in Section 7.7.60
type ShortRangeModeChangeCompleteEvent struct {
	Status uint8
	PhysicalLinkHandle uint8
	ShortRangeModeState uint8
}

func (o *ShortRangeModeChangeCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.PhysicalLinkHandle = uint8(r.GetOne())
	o.ShortRangeModeState = uint8(r.GetOne())
	return r.Valid()
}

// ShortRangeModeChangeCompleteEventCallbackType is the type of the callback function for ShortRangeModeChangeCompleteEvent.
type ShortRangeModeChangeCompleteEventCallbackType func(*ShortRangeModeChangeCompleteEvent) *ShortRangeModeChangeCompleteEvent

// SetShortRangeModeChangeCompleteEventCallback configures the callback for ShortRangeModeChangeCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetShortRangeModeChangeCompleteEventCallback(cb ShortRangeModeChangeCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.shortRangeModeChangeCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(109, cb != nil);
}
// AMPStatusChangeEvent represents the event specified in Section 7.7.61
type AMPStatusChangeEvent struct {
	Status uint8
	AMPStatus uint8
}

func (o *AMPStatusChangeEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.AMPStatus = uint8(r.GetOne())
	return r.Valid()
}

// AMPStatusChangeEventCallbackType is the type of the callback function for AMPStatusChangeEvent.
type AMPStatusChangeEventCallbackType func(*AMPStatusChangeEvent) *AMPStatusChangeEvent

// SetAMPStatusChangeEventCallback configures the callback for AMPStatusChangeEvent. Passing nil will disable the callback.
func (e *EventHandler) SetAMPStatusChangeEventCallback(cb AMPStatusChangeEventCallbackType) error {
	e.cbMutex.Lock()
	e.aMPStatusChangeEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(110, cb != nil);
}
// AMPStartTestEvent represents the event specified in Section 7.7.62
type AMPStartTestEvent struct {
	Status uint8
	TestScenario uint8
}

func (o *AMPStartTestEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.TestScenario = uint8(r.GetOne())
	return r.Valid()
}

// AMPStartTestEventCallbackType is the type of the callback function for AMPStartTestEvent.
type AMPStartTestEventCallbackType func(*AMPStartTestEvent) *AMPStartTestEvent

// SetAMPStartTestEventCallback configures the callback for AMPStartTestEvent. Passing nil will disable the callback.
func (e *EventHandler) SetAMPStartTestEventCallback(cb AMPStartTestEventCallbackType) error {
	e.cbMutex.Lock()
	e.aMPStartTestEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(111, cb != nil);
}
// AMPTestEndEvent represents the event specified in Section 7.7.63
type AMPTestEndEvent struct {
	Status uint8
	TestScenario uint8
}

func (o *AMPTestEndEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.TestScenario = uint8(r.GetOne())
	return r.Valid()
}

// AMPTestEndEventCallbackType is the type of the callback function for AMPTestEndEvent.
type AMPTestEndEventCallbackType func(*AMPTestEndEvent) *AMPTestEndEvent

// SetAMPTestEndEventCallback configures the callback for AMPTestEndEvent. Passing nil will disable the callback.
func (e *EventHandler) SetAMPTestEndEventCallback(cb AMPTestEndEventCallbackType) error {
	e.cbMutex.Lock()
	e.aMPTestEndEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(112, cb != nil);
}
// AMPReceiverReportEvent represents the event specified in Section 7.7.64
type AMPReceiverReportEvent struct {
	ControllerType uint8
	Reason uint8
}

func (o *AMPReceiverReportEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.ControllerType = uint8(r.GetOne())
	o.Reason = uint8(r.GetOne())
	return r.Valid()
}

// AMPReceiverReportEventCallbackType is the type of the callback function for AMPReceiverReportEvent.
type AMPReceiverReportEventCallbackType func(*AMPReceiverReportEvent) *AMPReceiverReportEvent

// SetAMPReceiverReportEventCallback configures the callback for AMPReceiverReportEvent. Passing nil will disable the callback.
func (e *EventHandler) SetAMPReceiverReportEventCallback(cb AMPReceiverReportEventCallbackType) error {
	e.cbMutex.Lock()
	e.aMPReceiverReportEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(113, cb != nil);
}
// LEConnectionCompleteEvent represents the event specified in Section 7.7.65.1
type LEConnectionCompleteEvent struct {
	SubeventCode uint8
	Status uint8
	ConnectionHandle uint16
	Role uint8
	PeerAddressType bleutil.MacAddrType
	PeerAddress bleutil.MacAddr
	ConnectionInterval uint16
	ConnectionLatency uint16
	SupervisionTimeout uint16
	MasterClockAccuracy uint8
}

func (o *LEConnectionCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Role = uint8(r.GetOne())
	o.PeerAddressType = bleutil.MacAddrType(r.GetOne())
	o.PeerAddress.Decode(r.Get(6))
	o.ConnectionInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.ConnectionLatency = binary.LittleEndian.Uint16(r.Get(2))
	o.SupervisionTimeout = binary.LittleEndian.Uint16(r.Get(2))
	o.MasterClockAccuracy = uint8(r.GetOne())
	return r.Valid()
}

// LEConnectionCompleteEventCallbackType is the type of the callback function for LEConnectionCompleteEvent.
type LEConnectionCompleteEventCallbackType func(*LEConnectionCompleteEvent) *LEConnectionCompleteEvent

// SetLEConnectionCompleteEventCallback configures the callback for LEConnectionCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEConnectionCompleteEventCallback(cb LEConnectionCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEConnectionCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(200, cb != nil);
}
// LEAdvertisingReportEvent represents the event specified in Section 7.7.65.2
type LEAdvertisingReportEvent struct {
	SubeventCode uint8
	NumReports uint8
	EventType []uint8
	AddressType []bleutil.MacAddrType
	Address []bleutil.MacAddr
	DataLength []uint8
	Data [][]byte
	RSSI []uint8
}

func (o *LEAdvertisingReportEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.NumReports = uint8(r.GetOne())
	if cap(o.EventType) < int(o.NumReports) {
		o.EventType = make([]uint8, 0, int(o.NumReports))
	}
	o.EventType = o.EventType[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.EventType[j] = uint8(r.GetOne())
	}
	if cap(o.AddressType) < int(o.NumReports) {
		o.AddressType = make([]bleutil.MacAddrType, 0, int(o.NumReports))
	}
	o.AddressType = o.AddressType[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.AddressType[j] = bleutil.MacAddrType(r.GetOne())
	}
	if cap(o.Address) < int(o.NumReports) {
		o.Address = make([]bleutil.MacAddr, 0, int(o.NumReports))
	}
	o.Address = o.Address[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.Address[j].Decode(r.Get(6))
	}
	if cap(o.DataLength) < int(o.NumReports) {
		o.DataLength = make([]uint8, 0, int(o.NumReports))
	}
	o.DataLength = o.DataLength[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.DataLength[j] = uint8(r.GetOne())
	}
	if cap(o.Data) < int(o.NumReports) {
		o.Data = make([][]byte, 0, int(o.NumReports))
	}
	o.Data = o.Data[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.Data[j] = append(o.Data[j][:0], r.Get(int(o.DataLength[j]))...)
	}
	if cap(o.RSSI) < int(o.NumReports) {
		o.RSSI = make([]uint8, 0, int(o.NumReports))
	}
	o.RSSI = o.RSSI[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.RSSI[j] = uint8(r.GetOne())
	}
	return r.Valid()
}

// LEAdvertisingReportEventCallbackType is the type of the callback function for LEAdvertisingReportEvent.
type LEAdvertisingReportEventCallbackType func(*LEAdvertisingReportEvent) *LEAdvertisingReportEvent

// SetLEAdvertisingReportEventCallback configures the callback for LEAdvertisingReportEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEAdvertisingReportEventCallback(cb LEAdvertisingReportEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEAdvertisingReportEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(201, cb != nil);
}
// LEConnectionUpdateCompleteEvent represents the event specified in Section 7.7.65.3
type LEConnectionUpdateCompleteEvent struct {
	SubeventCode uint8
	Status uint8
	ConnectionHandle uint16
	ConnectionInterval uint16
	ConnectionLatency uint16
	SupervisionTimeout uint16
}

func (o *LEConnectionUpdateCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ConnectionInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.ConnectionLatency = binary.LittleEndian.Uint16(r.Get(2))
	o.SupervisionTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEConnectionUpdateCompleteEventCallbackType is the type of the callback function for LEConnectionUpdateCompleteEvent.
type LEConnectionUpdateCompleteEventCallbackType func(*LEConnectionUpdateCompleteEvent) *LEConnectionUpdateCompleteEvent

// SetLEConnectionUpdateCompleteEventCallback configures the callback for LEConnectionUpdateCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEConnectionUpdateCompleteEventCallback(cb LEConnectionUpdateCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEConnectionUpdateCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(202, cb != nil);
}
// LEReadRemoteFeaturesCompleteEvent represents the event specified in Section 7.7.65.4
type LEReadRemoteFeaturesCompleteEvent struct {
	SubeventCode uint8
	Status uint8
	ConnectionHandle uint16
	LEFeatures uint64
}

func (o *LEReadRemoteFeaturesCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LEFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// LEReadRemoteFeaturesCompleteEventCallbackType is the type of the callback function for LEReadRemoteFeaturesCompleteEvent.
type LEReadRemoteFeaturesCompleteEventCallbackType func(*LEReadRemoteFeaturesCompleteEvent) *LEReadRemoteFeaturesCompleteEvent

// SetLEReadRemoteFeaturesCompleteEventCallback configures the callback for LEReadRemoteFeaturesCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEReadRemoteFeaturesCompleteEventCallback(cb LEReadRemoteFeaturesCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEReadRemoteFeaturesCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(203, cb != nil);
}
// LELongTermKeyRequestEvent represents the event specified in Section 7.7.65.5
type LELongTermKeyRequestEvent struct {
	SubeventCode uint8
	ConnectionHandle uint16
	RandomNumber uint64
	EncryptedDiversifier uint16
}

func (o *LELongTermKeyRequestEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.RandomNumber = binary.LittleEndian.Uint64(r.Get(8))
	o.EncryptedDiversifier = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LELongTermKeyRequestEventCallbackType is the type of the callback function for LELongTermKeyRequestEvent.
type LELongTermKeyRequestEventCallbackType func(*LELongTermKeyRequestEvent) *LELongTermKeyRequestEvent

// SetLELongTermKeyRequestEventCallback configures the callback for LELongTermKeyRequestEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLELongTermKeyRequestEventCallback(cb LELongTermKeyRequestEventCallbackType) error {
	e.cbMutex.Lock()
	e.lELongTermKeyRequestEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(204, cb != nil);
}
// LERemoteConnectionParameterRequestEvent represents the event specified in Section 7.7.65.6
type LERemoteConnectionParameterRequestEvent struct {
	SubeventCode uint8
	ConnectionHandle uint16
	IntervalMin uint16
	IntervalMax uint16
	Latency uint16
	Timeout uint16
}

func (o *LERemoteConnectionParameterRequestEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.IntervalMin = binary.LittleEndian.Uint16(r.Get(2))
	o.IntervalMax = binary.LittleEndian.Uint16(r.Get(2))
	o.Latency = binary.LittleEndian.Uint16(r.Get(2))
	o.Timeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LERemoteConnectionParameterRequestEventCallbackType is the type of the callback function for LERemoteConnectionParameterRequestEvent.
type LERemoteConnectionParameterRequestEventCallbackType func(*LERemoteConnectionParameterRequestEvent) *LERemoteConnectionParameterRequestEvent

// SetLERemoteConnectionParameterRequestEventCallback configures the callback for LERemoteConnectionParameterRequestEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLERemoteConnectionParameterRequestEventCallback(cb LERemoteConnectionParameterRequestEventCallbackType) error {
	e.cbMutex.Lock()
	e.lERemoteConnectionParameterRequestEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(205, cb != nil);
}
// LEDataLengthChangeEvent represents the event specified in Section 7.7.65.7
type LEDataLengthChangeEvent struct {
	SubeventCode uint8
	ConnectionHandle uint16
	MaxTXOctets uint16
	MaxTXTime uint16
	MaxRXOctets uint16
	MaxRXTime uint16
}

func (o *LEDataLengthChangeEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxTXOctets = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxTXTime = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxRXOctets = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxRXTime = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEDataLengthChangeEventCallbackType is the type of the callback function for LEDataLengthChangeEvent.
type LEDataLengthChangeEventCallbackType func(*LEDataLengthChangeEvent) *LEDataLengthChangeEvent

// SetLEDataLengthChangeEventCallback configures the callback for LEDataLengthChangeEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEDataLengthChangeEventCallback(cb LEDataLengthChangeEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEDataLengthChangeEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(206, cb != nil);
}
// LEReadLocalP256PublicKeyCompleteEvent represents the event specified in Section 7.7.65.8
type LEReadLocalP256PublicKeyCompleteEvent struct {
	SubeventCode uint8
	Status uint8
	LocalP256PublicKey [64]byte
}

func (o *LEReadLocalP256PublicKeyCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	copy(o.LocalP256PublicKey[:], r.Get(64))
	return r.Valid()
}

// LEReadLocalP256PublicKeyCompleteEventCallbackType is the type of the callback function for LEReadLocalP256PublicKeyCompleteEvent.
type LEReadLocalP256PublicKeyCompleteEventCallbackType func(*LEReadLocalP256PublicKeyCompleteEvent) *LEReadLocalP256PublicKeyCompleteEvent

// SetLEReadLocalP256PublicKeyCompleteEventCallback configures the callback for LEReadLocalP256PublicKeyCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEReadLocalP256PublicKeyCompleteEventCallback(cb LEReadLocalP256PublicKeyCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEReadLocalP256PublicKeyCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(207, cb != nil);
}
// LEGenerateDHKeyCompleteEvent represents the event specified in Section 7.7.65.9
type LEGenerateDHKeyCompleteEvent struct {
	SubeventCode uint8
	Status uint8
	DHKey [32]byte
}

func (o *LEGenerateDHKeyCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	copy(o.DHKey[:], r.Get(32))
	return r.Valid()
}

// LEGenerateDHKeyCompleteEventCallbackType is the type of the callback function for LEGenerateDHKeyCompleteEvent.
type LEGenerateDHKeyCompleteEventCallbackType func(*LEGenerateDHKeyCompleteEvent) *LEGenerateDHKeyCompleteEvent

// SetLEGenerateDHKeyCompleteEventCallback configures the callback for LEGenerateDHKeyCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEGenerateDHKeyCompleteEventCallback(cb LEGenerateDHKeyCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEGenerateDHKeyCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(208, cb != nil);
}
// LEEnhancedConnectionCompleteEvent represents the event specified in Section 7.7.65.10
type LEEnhancedConnectionCompleteEvent struct {
	SubeventCode uint8
	Status uint8
	ConnectionHandle uint16
	Role uint8
	PeerAddressType bleutil.MacAddrType
	PeerAddress bleutil.MacAddr
	LocalResolvablePrivateAddress bleutil.MacAddr
	PeerResolvablePrivateAddress bleutil.MacAddr
	ConnectionInterval uint16
	ConnectionLatency uint16
	SupervisionTimeout uint16
	MasterClockAccuracy uint8
}

func (o *LEEnhancedConnectionCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Role = uint8(r.GetOne())
	o.PeerAddressType = bleutil.MacAddrType(r.GetOne())
	o.PeerAddress.Decode(r.Get(6))
	o.LocalResolvablePrivateAddress.Decode(r.Get(6))
	o.PeerResolvablePrivateAddress.Decode(r.Get(6))
	o.ConnectionInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.ConnectionLatency = binary.LittleEndian.Uint16(r.Get(2))
	o.SupervisionTimeout = binary.LittleEndian.Uint16(r.Get(2))
	o.MasterClockAccuracy = uint8(r.GetOne())
	return r.Valid()
}

// LEEnhancedConnectionCompleteEventCallbackType is the type of the callback function for LEEnhancedConnectionCompleteEvent.
type LEEnhancedConnectionCompleteEventCallbackType func(*LEEnhancedConnectionCompleteEvent) *LEEnhancedConnectionCompleteEvent

// SetLEEnhancedConnectionCompleteEventCallback configures the callback for LEEnhancedConnectionCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEEnhancedConnectionCompleteEventCallback(cb LEEnhancedConnectionCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEEnhancedConnectionCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(209, cb != nil);
}
// LEDirectedAdvertisingReportEvent represents the event specified in Section 7.7.65.11
type LEDirectedAdvertisingReportEvent struct {
	SubeventCode uint8
	NumReports uint8
	EventType []uint8
	AddressType []bleutil.MacAddrType
	Address []bleutil.MacAddr
	DirectAddressType []bleutil.MacAddrType
	DirectAddress []bleutil.MacAddr
	RSSI []uint8
}

func (o *LEDirectedAdvertisingReportEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.NumReports = uint8(r.GetOne())
	if cap(o.EventType) < int(o.NumReports) {
		o.EventType = make([]uint8, 0, int(o.NumReports))
	}
	o.EventType = o.EventType[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.EventType[j] = uint8(r.GetOne())
	}
	if cap(o.AddressType) < int(o.NumReports) {
		o.AddressType = make([]bleutil.MacAddrType, 0, int(o.NumReports))
	}
	o.AddressType = o.AddressType[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.AddressType[j] = bleutil.MacAddrType(r.GetOne())
	}
	if cap(o.Address) < int(o.NumReports) {
		o.Address = make([]bleutil.MacAddr, 0, int(o.NumReports))
	}
	o.Address = o.Address[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.Address[j].Decode(r.Get(6))
	}
	if cap(o.DirectAddressType) < int(o.NumReports) {
		o.DirectAddressType = make([]bleutil.MacAddrType, 0, int(o.NumReports))
	}
	o.DirectAddressType = o.DirectAddressType[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.DirectAddressType[j] = bleutil.MacAddrType(r.GetOne())
	}
	if cap(o.DirectAddress) < int(o.NumReports) {
		o.DirectAddress = make([]bleutil.MacAddr, 0, int(o.NumReports))
	}
	o.DirectAddress = o.DirectAddress[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.DirectAddress[j].Decode(r.Get(6))
	}
	if cap(o.RSSI) < int(o.NumReports) {
		o.RSSI = make([]uint8, 0, int(o.NumReports))
	}
	o.RSSI = o.RSSI[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.RSSI[j] = uint8(r.GetOne())
	}
	return r.Valid()
}

// LEDirectedAdvertisingReportEventCallbackType is the type of the callback function for LEDirectedAdvertisingReportEvent.
type LEDirectedAdvertisingReportEventCallbackType func(*LEDirectedAdvertisingReportEvent) *LEDirectedAdvertisingReportEvent

// SetLEDirectedAdvertisingReportEventCallback configures the callback for LEDirectedAdvertisingReportEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEDirectedAdvertisingReportEventCallback(cb LEDirectedAdvertisingReportEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEDirectedAdvertisingReportEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(210, cb != nil);
}
// LEPHYUpdateCompleteEvent represents the event specified in Section 7.7.65.12
type LEPHYUpdateCompleteEvent struct {
	SubeventCode uint8
	Status uint8
	ConnectionHandle uint16
	TXPHY uint8
	RXPHY uint8
}

func (o *LEPHYUpdateCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.TXPHY = uint8(r.GetOne())
	o.RXPHY = uint8(r.GetOne())
	return r.Valid()
}

// LEPHYUpdateCompleteEventCallbackType is the type of the callback function for LEPHYUpdateCompleteEvent.
type LEPHYUpdateCompleteEventCallbackType func(*LEPHYUpdateCompleteEvent) *LEPHYUpdateCompleteEvent

// SetLEPHYUpdateCompleteEventCallback configures the callback for LEPHYUpdateCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEPHYUpdateCompleteEventCallback(cb LEPHYUpdateCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEPHYUpdateCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(211, cb != nil);
}
// LEExtendedAdvertisingReportEvent represents the event specified in Section 7.7.65.13
type LEExtendedAdvertisingReportEvent struct {
	SubeventCode uint8
	NumReports uint8
	EventType []uint16
	AddressType []bleutil.MacAddrType
	Address []bleutil.MacAddr
	PrimaryPHY []uint8
	SecondaryPHY []uint8
	AdvertisingSID []uint8
	TXPower []uint8
	RSSI []uint8
	PeriodicAdvertisingInterval []uint16
	DirectAddressType []bleutil.MacAddrType
	DirectAddress []bleutil.MacAddr
	DataLength []uint8
	Data [][]byte
}

func (o *LEExtendedAdvertisingReportEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.NumReports = uint8(r.GetOne())
	if cap(o.EventType) < int(o.NumReports) {
		o.EventType = make([]uint16, 0, int(o.NumReports))
	}
	o.EventType = o.EventType[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.EventType[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.AddressType) < int(o.NumReports) {
		o.AddressType = make([]bleutil.MacAddrType, 0, int(o.NumReports))
	}
	o.AddressType = o.AddressType[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.AddressType[j] = bleutil.MacAddrType(r.GetOne())
	}
	if cap(o.Address) < int(o.NumReports) {
		o.Address = make([]bleutil.MacAddr, 0, int(o.NumReports))
	}
	o.Address = o.Address[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.Address[j].Decode(r.Get(6))
	}
	if cap(o.PrimaryPHY) < int(o.NumReports) {
		o.PrimaryPHY = make([]uint8, 0, int(o.NumReports))
	}
	o.PrimaryPHY = o.PrimaryPHY[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.PrimaryPHY[j] = uint8(r.GetOne())
	}
	if cap(o.SecondaryPHY) < int(o.NumReports) {
		o.SecondaryPHY = make([]uint8, 0, int(o.NumReports))
	}
	o.SecondaryPHY = o.SecondaryPHY[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.SecondaryPHY[j] = uint8(r.GetOne())
	}
	if cap(o.AdvertisingSID) < int(o.NumReports) {
		o.AdvertisingSID = make([]uint8, 0, int(o.NumReports))
	}
	o.AdvertisingSID = o.AdvertisingSID[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.AdvertisingSID[j] = uint8(r.GetOne())
	}
	if cap(o.TXPower) < int(o.NumReports) {
		o.TXPower = make([]uint8, 0, int(o.NumReports))
	}
	o.TXPower = o.TXPower[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.TXPower[j] = uint8(r.GetOne())
	}
	if cap(o.RSSI) < int(o.NumReports) {
		o.RSSI = make([]uint8, 0, int(o.NumReports))
	}
	o.RSSI = o.RSSI[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.RSSI[j] = uint8(r.GetOne())
	}
	if cap(o.PeriodicAdvertisingInterval) < int(o.NumReports) {
		o.PeriodicAdvertisingInterval = make([]uint16, 0, int(o.NumReports))
	}
	o.PeriodicAdvertisingInterval = o.PeriodicAdvertisingInterval[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.PeriodicAdvertisingInterval[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	if cap(o.DirectAddressType) < int(o.NumReports) {
		o.DirectAddressType = make([]bleutil.MacAddrType, 0, int(o.NumReports))
	}
	o.DirectAddressType = o.DirectAddressType[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.DirectAddressType[j] = bleutil.MacAddrType(r.GetOne())
	}
	if cap(o.DirectAddress) < int(o.NumReports) {
		o.DirectAddress = make([]bleutil.MacAddr, 0, int(o.NumReports))
	}
	o.DirectAddress = o.DirectAddress[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.DirectAddress[j].Decode(r.Get(6))
	}
	if cap(o.DataLength) < int(o.NumReports) {
		o.DataLength = make([]uint8, 0, int(o.NumReports))
	}
	o.DataLength = o.DataLength[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.DataLength[j] = uint8(r.GetOne())
	}
	if cap(o.Data) < int(o.NumReports) {
		o.Data = make([][]byte, 0, int(o.NumReports))
	}
	o.Data = o.Data[:int(o.NumReports)]
	for j:=0; j<int(o.NumReports); j++ {
		o.Data[j] = append(o.Data[j][:0], r.Get(int(o.DataLength[j]))...)
	}
	return r.Valid()
}

// LEExtendedAdvertisingReportEventCallbackType is the type of the callback function for LEExtendedAdvertisingReportEvent.
type LEExtendedAdvertisingReportEventCallbackType func(*LEExtendedAdvertisingReportEvent) *LEExtendedAdvertisingReportEvent

// SetLEExtendedAdvertisingReportEventCallback configures the callback for LEExtendedAdvertisingReportEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEExtendedAdvertisingReportEventCallback(cb LEExtendedAdvertisingReportEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEExtendedAdvertisingReportEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(212, cb != nil);
}
// LEPeriodicAdvertisingSyncEstablishedEvent represents the event specified in Section 7.7.65.14
type LEPeriodicAdvertisingSyncEstablishedEvent struct {
	SubeventCode uint8
	Status uint8
	SyncHandle uint16
	AdvertisingSID uint8
	AdvertiserAddressType bleutil.MacAddrType
	AdvertiserAddress bleutil.MacAddr
	AdvertiserPHY uint8
	PeriodicAdvertisingInterval uint16
	AdvertiserClockAccuracy uint8
}

func (o *LEPeriodicAdvertisingSyncEstablishedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.AdvertisingSID = uint8(r.GetOne())
	o.AdvertiserAddressType = bleutil.MacAddrType(r.GetOne())
	o.AdvertiserAddress.Decode(r.Get(6))
	o.AdvertiserPHY = uint8(r.GetOne())
	o.PeriodicAdvertisingInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.AdvertiserClockAccuracy = uint8(r.GetOne())
	return r.Valid()
}

// LEPeriodicAdvertisingSyncEstablishedEventCallbackType is the type of the callback function for LEPeriodicAdvertisingSyncEstablishedEvent.
type LEPeriodicAdvertisingSyncEstablishedEventCallbackType func(*LEPeriodicAdvertisingSyncEstablishedEvent) *LEPeriodicAdvertisingSyncEstablishedEvent

// SetLEPeriodicAdvertisingSyncEstablishedEventCallback configures the callback for LEPeriodicAdvertisingSyncEstablishedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEPeriodicAdvertisingSyncEstablishedEventCallback(cb LEPeriodicAdvertisingSyncEstablishedEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEPeriodicAdvertisingSyncEstablishedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(213, cb != nil);
}
// LEPeriodicAdvertisingReportEvent represents the event specified in Section 7.7.65.15
type LEPeriodicAdvertisingReportEvent struct {
	SubeventCode uint8
	SyncHandle uint16
	TXPower uint8
	RSSI uint8
	CTEType uint8
	DataStatus uint8
	DataLength uint8
	Data []byte
}

func (o *LEPeriodicAdvertisingReportEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.TXPower = uint8(r.GetOne())
	o.RSSI = uint8(r.GetOne())
	o.CTEType = uint8(r.GetOne())
	o.DataStatus = uint8(r.GetOne())
	o.DataLength = uint8(r.GetOne())
	o.Data = append(o.Data[:0], r.GetRemainder()...)
	return r.Valid()
}

// LEPeriodicAdvertisingReportEventCallbackType is the type of the callback function for LEPeriodicAdvertisingReportEvent.
type LEPeriodicAdvertisingReportEventCallbackType func(*LEPeriodicAdvertisingReportEvent) *LEPeriodicAdvertisingReportEvent

// SetLEPeriodicAdvertisingReportEventCallback configures the callback for LEPeriodicAdvertisingReportEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEPeriodicAdvertisingReportEventCallback(cb LEPeriodicAdvertisingReportEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEPeriodicAdvertisingReportEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(214, cb != nil);
}
// LEPeriodicAdvertisingSyncLostEvent represents the event specified in Section 7.7.65.16
type LEPeriodicAdvertisingSyncLostEvent struct {
	SubeventCode uint8
	SyncHandle uint16
}

func (o *LEPeriodicAdvertisingSyncLostEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEPeriodicAdvertisingSyncLostEventCallbackType is the type of the callback function for LEPeriodicAdvertisingSyncLostEvent.
type LEPeriodicAdvertisingSyncLostEventCallbackType func(*LEPeriodicAdvertisingSyncLostEvent) *LEPeriodicAdvertisingSyncLostEvent

// SetLEPeriodicAdvertisingSyncLostEventCallback configures the callback for LEPeriodicAdvertisingSyncLostEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEPeriodicAdvertisingSyncLostEventCallback(cb LEPeriodicAdvertisingSyncLostEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEPeriodicAdvertisingSyncLostEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(215, cb != nil);
}
// LEScanTimeoutEvent represents the event specified in Section 7.7.65.17
type LEScanTimeoutEvent struct {
	SubeventCode uint8
}

func (o *LEScanTimeoutEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	return r.Valid()
}

// LEScanTimeoutEventCallbackType is the type of the callback function for LEScanTimeoutEvent.
type LEScanTimeoutEventCallbackType func(*LEScanTimeoutEvent) *LEScanTimeoutEvent

// SetLEScanTimeoutEventCallback configures the callback for LEScanTimeoutEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEScanTimeoutEventCallback(cb LEScanTimeoutEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEScanTimeoutEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(216, cb != nil);
}
// LEAdvertisingSetTerminatedEvent represents the event specified in Section 7.7.65.18
type LEAdvertisingSetTerminatedEvent struct {
	SubeventCode uint8
	Status uint8
	AdvertisingHandle uint8
	ConnectionHandle uint16
	NumCompletedExtendedAdvertisingEvents uint8
}

func (o *LEAdvertisingSetTerminatedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.AdvertisingHandle = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.NumCompletedExtendedAdvertisingEvents = uint8(r.GetOne())
	return r.Valid()
}

// LEAdvertisingSetTerminatedEventCallbackType is the type of the callback function for LEAdvertisingSetTerminatedEvent.
type LEAdvertisingSetTerminatedEventCallbackType func(*LEAdvertisingSetTerminatedEvent) *LEAdvertisingSetTerminatedEvent

// SetLEAdvertisingSetTerminatedEventCallback configures the callback for LEAdvertisingSetTerminatedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEAdvertisingSetTerminatedEventCallback(cb LEAdvertisingSetTerminatedEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEAdvertisingSetTerminatedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(217, cb != nil);
}
// LEScanRequestReceivedEvent represents the event specified in Section 7.7.65.19
type LEScanRequestReceivedEvent struct {
	SubeventCode uint8
	AdvertisingHandle uint8
	ScannerAddressType bleutil.MacAddrType
	ScannerAddress bleutil.MacAddr
}

func (o *LEScanRequestReceivedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.AdvertisingHandle = uint8(r.GetOne())
	o.ScannerAddressType = bleutil.MacAddrType(r.GetOne())
	o.ScannerAddress.Decode(r.Get(6))
	return r.Valid()
}

// LEScanRequestReceivedEventCallbackType is the type of the callback function for LEScanRequestReceivedEvent.
type LEScanRequestReceivedEventCallbackType func(*LEScanRequestReceivedEvent) *LEScanRequestReceivedEvent

// SetLEScanRequestReceivedEventCallback configures the callback for LEScanRequestReceivedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEScanRequestReceivedEventCallback(cb LEScanRequestReceivedEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEScanRequestReceivedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(218, cb != nil);
}
// LEChannelSelectionAlgorithmEvent represents the event specified in Section 7.7.65.20
type LEChannelSelectionAlgorithmEvent struct {
	SubeventCode uint8
	ConnectionHandle uint16
	ChannelSelectionAlgorithm uint8
}

func (o *LEChannelSelectionAlgorithmEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ChannelSelectionAlgorithm = uint8(r.GetOne())
	return r.Valid()
}

// LEChannelSelectionAlgorithmEventCallbackType is the type of the callback function for LEChannelSelectionAlgorithmEvent.
type LEChannelSelectionAlgorithmEventCallbackType func(*LEChannelSelectionAlgorithmEvent) *LEChannelSelectionAlgorithmEvent

// SetLEChannelSelectionAlgorithmEventCallback configures the callback for LEChannelSelectionAlgorithmEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEChannelSelectionAlgorithmEventCallback(cb LEChannelSelectionAlgorithmEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEChannelSelectionAlgorithmEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(219, cb != nil);
}
// LEConnectionlessIQReportEvent represents the event specified in Section 7.7.65.21
type LEConnectionlessIQReportEvent struct {
	SubeventCode uint8
	SyncHandle uint16
	ChannelIndex uint8
	RSSI uint16
	RSSIAntennaID uint8
	CTEType uint8
	SlotDurations uint8
	PacketStatus uint8
	PeriodicEventCounter uint16
	SampleCount uint8
	ISample []uint8
	QSample []uint8
}

func (o *LEConnectionlessIQReportEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ChannelIndex = uint8(r.GetOne())
	o.RSSI = binary.LittleEndian.Uint16(r.Get(2))
	o.RSSIAntennaID = uint8(r.GetOne())
	o.CTEType = uint8(r.GetOne())
	o.SlotDurations = uint8(r.GetOne())
	o.PacketStatus = uint8(r.GetOne())
	o.PeriodicEventCounter = binary.LittleEndian.Uint16(r.Get(2))
	o.SampleCount = uint8(r.GetOne())
	if cap(o.ISample) < int(o.SampleCount) {
		o.ISample = make([]uint8, 0, int(o.SampleCount))
	}
	o.ISample = o.ISample[:int(o.SampleCount)]
	for j:=0; j<int(o.SampleCount); j++ {
		o.ISample[j] = uint8(r.GetOne())
	}
	if cap(o.QSample) < int(o.SampleCount) {
		o.QSample = make([]uint8, 0, int(o.SampleCount))
	}
	o.QSample = o.QSample[:int(o.SampleCount)]
	for j:=0; j<int(o.SampleCount); j++ {
		o.QSample[j] = uint8(r.GetOne())
	}
	return r.Valid()
}

// LEConnectionlessIQReportEventCallbackType is the type of the callback function for LEConnectionlessIQReportEvent.
type LEConnectionlessIQReportEventCallbackType func(*LEConnectionlessIQReportEvent) *LEConnectionlessIQReportEvent

// SetLEConnectionlessIQReportEventCallback configures the callback for LEConnectionlessIQReportEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEConnectionlessIQReportEventCallback(cb LEConnectionlessIQReportEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEConnectionlessIQReportEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(220, cb != nil);
}
// LEConnectionIQReportEvent represents the event specified in Section 7.7.65.22
type LEConnectionIQReportEvent struct {
	SubeventCode uint8
	ConnectionHandle uint16
	RXPHY uint8
	DataChannelIndex uint8
	RSSI uint16
	RSSIAntennaID uint8
	CTEType uint8
	SlotDurations uint8
	PacketStatus uint8
	ConnectionEventCounter uint16
	SampleCount uint8
	ISample []uint8
	QSample []uint8
}

func (o *LEConnectionIQReportEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.RXPHY = uint8(r.GetOne())
	o.DataChannelIndex = uint8(r.GetOne())
	o.RSSI = binary.LittleEndian.Uint16(r.Get(2))
	o.RSSIAntennaID = uint8(r.GetOne())
	o.CTEType = uint8(r.GetOne())
	o.SlotDurations = uint8(r.GetOne())
	o.PacketStatus = uint8(r.GetOne())
	o.ConnectionEventCounter = binary.LittleEndian.Uint16(r.Get(2))
	o.SampleCount = uint8(r.GetOne())
	if cap(o.ISample) < int(o.SampleCount) {
		o.ISample = make([]uint8, 0, int(o.SampleCount))
	}
	o.ISample = o.ISample[:int(o.SampleCount)]
	for j:=0; j<int(o.SampleCount); j++ {
		o.ISample[j] = uint8(r.GetOne())
	}
	if cap(o.QSample) < int(o.SampleCount) {
		o.QSample = make([]uint8, 0, int(o.SampleCount))
	}
	o.QSample = o.QSample[:int(o.SampleCount)]
	for j:=0; j<int(o.SampleCount); j++ {
		o.QSample[j] = uint8(r.GetOne())
	}
	return r.Valid()
}

// LEConnectionIQReportEventCallbackType is the type of the callback function for LEConnectionIQReportEvent.
type LEConnectionIQReportEventCallbackType func(*LEConnectionIQReportEvent) *LEConnectionIQReportEvent

// SetLEConnectionIQReportEventCallback configures the callback for LEConnectionIQReportEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEConnectionIQReportEventCallback(cb LEConnectionIQReportEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEConnectionIQReportEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(221, cb != nil);
}
// LECTERequestFailedEvent represents the event specified in Section 7.7.65.23
type LECTERequestFailedEvent struct {
	SubeventCode uint8
	Status uint8
	ConnectionHandle uint16
}

func (o *LECTERequestFailedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LECTERequestFailedEventCallbackType is the type of the callback function for LECTERequestFailedEvent.
type LECTERequestFailedEventCallbackType func(*LECTERequestFailedEvent) *LECTERequestFailedEvent

// SetLECTERequestFailedEventCallback configures the callback for LECTERequestFailedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLECTERequestFailedEventCallback(cb LECTERequestFailedEventCallbackType) error {
	e.cbMutex.Lock()
	e.lECTERequestFailedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(222, cb != nil);
}
// LEPeriodicAdvertisingSyncTransferReceivedEvent represents the event specified in Section 7.7.65.24
type LEPeriodicAdvertisingSyncTransferReceivedEvent struct {
	SubeventCode uint8
	Status uint8
	ConnectionHandle uint16
	ServiceData uint16
	SyncHandle uint16
	AdvertisingSID uint8
	AdvertiserAddressType bleutil.MacAddrType
	AdvertiserAddress bleutil.MacAddr
	AdvertiserPHY uint8
	PeriodicAdvertisingInterval uint16
	AdvertiserClockAccuracy uint8
}

func (o *LEPeriodicAdvertisingSyncTransferReceivedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ServiceData = binary.LittleEndian.Uint16(r.Get(2))
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.AdvertisingSID = uint8(r.GetOne())
	o.AdvertiserAddressType = bleutil.MacAddrType(r.GetOne())
	o.AdvertiserAddress.Decode(r.Get(6))
	o.AdvertiserPHY = uint8(r.GetOne())
	o.PeriodicAdvertisingInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.AdvertiserClockAccuracy = uint8(r.GetOne())
	return r.Valid()
}

// LEPeriodicAdvertisingSyncTransferReceivedEventCallbackType is the type of the callback function for LEPeriodicAdvertisingSyncTransferReceivedEvent.
type LEPeriodicAdvertisingSyncTransferReceivedEventCallbackType func(*LEPeriodicAdvertisingSyncTransferReceivedEvent) *LEPeriodicAdvertisingSyncTransferReceivedEvent

// SetLEPeriodicAdvertisingSyncTransferReceivedEventCallback configures the callback for LEPeriodicAdvertisingSyncTransferReceivedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEPeriodicAdvertisingSyncTransferReceivedEventCallback(cb LEPeriodicAdvertisingSyncTransferReceivedEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEPeriodicAdvertisingSyncTransferReceivedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(223, cb != nil);
}
// LECISEstablishedEvent represents the event specified in Section 7.7.65.25
type LECISEstablishedEvent struct {
	SubeventCode uint8
	Status uint8
	ConnectionHandle uint16
	PHYMToS uint8
	PHYSToM uint8
	NSE uint8
	BNMToS uint8
	BNSToM uint8
	FTMToS uint8
	FTSToM uint8
	MaxPDUMToS uint16
	MaxPDUSToM uint16
	ISOinterval uint16
}

func (o *LECISEstablishedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PHYMToS = uint8(r.GetOne())
	o.PHYSToM = uint8(r.GetOne())
	o.NSE = uint8(r.GetOne())
	o.BNMToS = uint8(r.GetOne())
	o.BNSToM = uint8(r.GetOne())
	o.FTMToS = uint8(r.GetOne())
	o.FTSToM = uint8(r.GetOne())
	o.MaxPDUMToS = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxPDUSToM = binary.LittleEndian.Uint16(r.Get(2))
	o.ISOinterval = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LECISEstablishedEventCallbackType is the type of the callback function for LECISEstablishedEvent.
type LECISEstablishedEventCallbackType func(*LECISEstablishedEvent) *LECISEstablishedEvent

// SetLECISEstablishedEventCallback configures the callback for LECISEstablishedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLECISEstablishedEventCallback(cb LECISEstablishedEventCallbackType) error {
	e.cbMutex.Lock()
	e.lECISEstablishedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(224, cb != nil);
}
// LECISRequestEvent represents the event specified in Section 7.7.65.26
type LECISRequestEvent struct {
	SubeventCode uint8
	ACLConnectionHandle uint16
	CISConnectonHandle uint16
	CIGID uint8
	CISID uint8
}

func (o *LECISRequestEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.ACLConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.CISConnectonHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.CIGID = uint8(r.GetOne())
	o.CISID = uint8(r.GetOne())
	return r.Valid()
}

// LECISRequestEventCallbackType is the type of the callback function for LECISRequestEvent.
type LECISRequestEventCallbackType func(*LECISRequestEvent) *LECISRequestEvent

// SetLECISRequestEventCallback configures the callback for LECISRequestEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLECISRequestEventCallback(cb LECISRequestEventCallbackType) error {
	e.cbMutex.Lock()
	e.lECISRequestEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(225, cb != nil);
}
// LECreateBIGCompleteEvent represents the event specified in Section 7.7.65.27
type LECreateBIGCompleteEvent struct {
	SubeventCode uint8
	Status uint8
	BIGHandle uint8
	BIGSyncDelay uint32
	TransportLatencyBIG uint32
	PHY uint8
	NSE uint8
	BN uint8
	PTO uint8
	IRC uint8
	MaxPDU uint16
	ISOInterval uint16
	NumBIS uint8
	ConnectionHandle []uint16
}

func (o *LECreateBIGCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.BIGHandle = uint8(r.GetOne())
	o.BIGSyncDelay = bleutil.DecodeUint24(r.Get(3))
	o.TransportLatencyBIG = bleutil.DecodeUint24(r.Get(3))
	o.PHY = uint8(r.GetOne())
	o.NSE = uint8(r.GetOne())
	o.BN = uint8(r.GetOne())
	o.PTO = uint8(r.GetOne())
	o.IRC = uint8(r.GetOne())
	o.MaxPDU = binary.LittleEndian.Uint16(r.Get(2))
	o.ISOInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.NumBIS = uint8(r.GetOne())
	if cap(o.ConnectionHandle) < int(o.NumBIS) {
		o.ConnectionHandle = make([]uint16, 0, int(o.NumBIS))
	}
	o.ConnectionHandle = o.ConnectionHandle[:int(o.NumBIS)]
	for j:=0; j<int(o.NumBIS); j++ {
		o.ConnectionHandle[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// LECreateBIGCompleteEventCallbackType is the type of the callback function for LECreateBIGCompleteEvent.
type LECreateBIGCompleteEventCallbackType func(*LECreateBIGCompleteEvent) *LECreateBIGCompleteEvent

// SetLECreateBIGCompleteEventCallback configures the callback for LECreateBIGCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLECreateBIGCompleteEventCallback(cb LECreateBIGCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.lECreateBIGCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(226, cb != nil);
}
// LETerminateBIGCompleteEvent represents the event specified in Section 7.7.65.28
type LETerminateBIGCompleteEvent struct {
	SubeventCode uint8
	BIGHandle uint8
	Reason uint8
}

func (o *LETerminateBIGCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.BIGHandle = uint8(r.GetOne())
	o.Reason = uint8(r.GetOne())
	return r.Valid()
}

// LETerminateBIGCompleteEventCallbackType is the type of the callback function for LETerminateBIGCompleteEvent.
type LETerminateBIGCompleteEventCallbackType func(*LETerminateBIGCompleteEvent) *LETerminateBIGCompleteEvent

// SetLETerminateBIGCompleteEventCallback configures the callback for LETerminateBIGCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLETerminateBIGCompleteEventCallback(cb LETerminateBIGCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.lETerminateBIGCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(227, cb != nil);
}
// LEBIGSyncEstablishedEvent represents the event specified in Section 7.7.65.29
type LEBIGSyncEstablishedEvent struct {
	SubeventCode uint8
	Status uint8
	BIGHandle uint8
	TransportLatencyBIG uint32
	NSE uint8
	BN uint8
	PTO uint8
	IRC uint8
	MaxPDU uint16
	ISOInterval uint16
	NumBIS uint8
	ConnectionHandle []uint16
}

func (o *LEBIGSyncEstablishedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.BIGHandle = uint8(r.GetOne())
	o.TransportLatencyBIG = bleutil.DecodeUint24(r.Get(3))
	o.NSE = uint8(r.GetOne())
	o.BN = uint8(r.GetOne())
	o.PTO = uint8(r.GetOne())
	o.IRC = uint8(r.GetOne())
	o.MaxPDU = binary.LittleEndian.Uint16(r.Get(2))
	o.ISOInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.NumBIS = uint8(r.GetOne())
	if cap(o.ConnectionHandle) < int(o.NumBIS) {
		o.ConnectionHandle = make([]uint16, 0, int(o.NumBIS))
	}
	o.ConnectionHandle = o.ConnectionHandle[:int(o.NumBIS)]
	for j:=0; j<int(o.NumBIS); j++ {
		o.ConnectionHandle[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// LEBIGSyncEstablishedEventCallbackType is the type of the callback function for LEBIGSyncEstablishedEvent.
type LEBIGSyncEstablishedEventCallbackType func(*LEBIGSyncEstablishedEvent) *LEBIGSyncEstablishedEvent

// SetLEBIGSyncEstablishedEventCallback configures the callback for LEBIGSyncEstablishedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEBIGSyncEstablishedEventCallback(cb LEBIGSyncEstablishedEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEBIGSyncEstablishedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(228, cb != nil);
}
// LEBIGSyncLostEvent represents the event specified in Section 7.7.65.30
type LEBIGSyncLostEvent struct {
	SubeventCode uint8
	BIGHandle uint8
	Reason uint8
}

func (o *LEBIGSyncLostEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.BIGHandle = uint8(r.GetOne())
	o.Reason = uint8(r.GetOne())
	return r.Valid()
}

// LEBIGSyncLostEventCallbackType is the type of the callback function for LEBIGSyncLostEvent.
type LEBIGSyncLostEventCallbackType func(*LEBIGSyncLostEvent) *LEBIGSyncLostEvent

// SetLEBIGSyncLostEventCallback configures the callback for LEBIGSyncLostEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEBIGSyncLostEventCallback(cb LEBIGSyncLostEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEBIGSyncLostEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(229, cb != nil);
}
// LERequestPeerSCACompleteEvent represents the event specified in Section 7.7.65.31
type LERequestPeerSCACompleteEvent struct {
	SubeventCode uint8
	Status uint8
	ConnectionHandle uint16
	PeerClockAccuracy uint8
}

func (o *LERequestPeerSCACompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PeerClockAccuracy = uint8(r.GetOne())
	return r.Valid()
}

// LERequestPeerSCACompleteEventCallbackType is the type of the callback function for LERequestPeerSCACompleteEvent.
type LERequestPeerSCACompleteEventCallbackType func(*LERequestPeerSCACompleteEvent) *LERequestPeerSCACompleteEvent

// SetLERequestPeerSCACompleteEventCallback configures the callback for LERequestPeerSCACompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLERequestPeerSCACompleteEventCallback(cb LERequestPeerSCACompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.lERequestPeerSCACompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(230, cb != nil);
}
// LEPathLossThresholdEvent represents the event specified in Section 7.7.65.32
type LEPathLossThresholdEvent struct {
	SubeventCode uint8
	ConnectionHandle uint16
	CurrentPathLoss uint8
	ZoneEntered uint8
}

func (o *LEPathLossThresholdEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.CurrentPathLoss = uint8(r.GetOne())
	o.ZoneEntered = uint8(r.GetOne())
	return r.Valid()
}

// LEPathLossThresholdEventCallbackType is the type of the callback function for LEPathLossThresholdEvent.
type LEPathLossThresholdEventCallbackType func(*LEPathLossThresholdEvent) *LEPathLossThresholdEvent

// SetLEPathLossThresholdEventCallback configures the callback for LEPathLossThresholdEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEPathLossThresholdEventCallback(cb LEPathLossThresholdEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEPathLossThresholdEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(231, cb != nil);
}
// LETransmitPowerReportingEvent represents the event specified in Section 7.7.65.33
type LETransmitPowerReportingEvent struct {
	SubeventCode uint8
	Status uint8
	ConnectionHandle uint16
	Reason uint8
	PHY uint8
	TransmitPowerLevel uint8
	TransmitPowerLevelFlag uint8
	Delta uint8
}

func (o *LETransmitPowerReportingEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Reason = uint8(r.GetOne())
	o.PHY = uint8(r.GetOne())
	o.TransmitPowerLevel = uint8(r.GetOne())
	o.TransmitPowerLevelFlag = uint8(r.GetOne())
	o.Delta = uint8(r.GetOne())
	return r.Valid()
}

// LETransmitPowerReportingEventCallbackType is the type of the callback function for LETransmitPowerReportingEvent.
type LETransmitPowerReportingEventCallbackType func(*LETransmitPowerReportingEvent) *LETransmitPowerReportingEvent

// SetLETransmitPowerReportingEventCallback configures the callback for LETransmitPowerReportingEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLETransmitPowerReportingEventCallback(cb LETransmitPowerReportingEventCallbackType) error {
	e.cbMutex.Lock()
	e.lETransmitPowerReportingEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(232, cb != nil);
}
// LEBIGInfoAdvertisingReportEvent represents the event specified in Section 7.7.65.34
type LEBIGInfoAdvertisingReportEvent struct {
	SubeventCode uint8
	SyncHandle uint16
	NumBIS uint8
	NSE uint8
	ISOInterval uint16
	BN uint8
	PTO uint8
	IRC uint8
	MaxPDU uint16
	SDUInterval uint32
	MaxSDU uint16
	PHY uint8
	Framing uint8
	Encryption uint8
}

func (o *LEBIGInfoAdvertisingReportEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.SubeventCode = uint8(r.GetOne())
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.NumBIS = uint8(r.GetOne())
	o.NSE = uint8(r.GetOne())
	o.ISOInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.BN = uint8(r.GetOne())
	o.PTO = uint8(r.GetOne())
	o.IRC = uint8(r.GetOne())
	o.MaxPDU = binary.LittleEndian.Uint16(r.Get(2))
	o.SDUInterval = bleutil.DecodeUint24(r.Get(3))
	o.MaxSDU = binary.LittleEndian.Uint16(r.Get(2))
	o.PHY = uint8(r.GetOne())
	o.Framing = uint8(r.GetOne())
	o.Encryption = uint8(r.GetOne())
	return r.Valid()
}

// LEBIGInfoAdvertisingReportEventCallbackType is the type of the callback function for LEBIGInfoAdvertisingReportEvent.
type LEBIGInfoAdvertisingReportEventCallbackType func(*LEBIGInfoAdvertisingReportEvent) *LEBIGInfoAdvertisingReportEvent

// SetLEBIGInfoAdvertisingReportEventCallback configures the callback for LEBIGInfoAdvertisingReportEvent. Passing nil will disable the callback.
func (e *EventHandler) SetLEBIGInfoAdvertisingReportEventCallback(cb LEBIGInfoAdvertisingReportEventCallbackType) error {
	e.cbMutex.Lock()
	e.lEBIGInfoAdvertisingReportEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(233, cb != nil);
}
// TriggeredClockCaptureEvent represents the event specified in Section 7.7.66
type TriggeredClockCaptureEvent struct {
	ConnectionHandle uint16
	WhichClock uint8
	Clock uint32
	SlotOffset uint16
}

func (o *TriggeredClockCaptureEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.WhichClock = uint8(r.GetOne())
	o.Clock = binary.LittleEndian.Uint32(r.Get(4))
	o.SlotOffset = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// TriggeredClockCaptureEventCallbackType is the type of the callback function for TriggeredClockCaptureEvent.
type TriggeredClockCaptureEventCallbackType func(*TriggeredClockCaptureEvent) *TriggeredClockCaptureEvent

// SetTriggeredClockCaptureEventCallback configures the callback for TriggeredClockCaptureEvent. Passing nil will disable the callback.
func (e *EventHandler) SetTriggeredClockCaptureEventCallback(cb TriggeredClockCaptureEventCallbackType) error {
	e.cbMutex.Lock()
	e.triggeredClockCaptureEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(114, cb != nil);
}
// SynchronizationTrainCompleteEvent represents the event specified in Section 7.7.67
type SynchronizationTrainCompleteEvent struct {
	Status uint8
}

func (o *SynchronizationTrainCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	return r.Valid()
}

// SynchronizationTrainCompleteEventCallbackType is the type of the callback function for SynchronizationTrainCompleteEvent.
type SynchronizationTrainCompleteEventCallbackType func(*SynchronizationTrainCompleteEvent) *SynchronizationTrainCompleteEvent

// SetSynchronizationTrainCompleteEventCallback configures the callback for SynchronizationTrainCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetSynchronizationTrainCompleteEventCallback(cb SynchronizationTrainCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.synchronizationTrainCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(115, cb != nil);
}
// SynchronizationTrainReceivedEvent represents the event specified in Section 7.7.68
type SynchronizationTrainReceivedEvent struct {
	Status uint8
	BDADDR bleutil.MacAddr
	ClockOffset uint32
	AFHChannelMap [10]byte
	LTADDR uint8
	NextBroadcastInstant uint32
	ConnectionlessSlaveBroadcastInterval uint16
	ServiceData uint8
}

func (o *SynchronizationTrainReceivedEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.BDADDR.Decode(r.Get(6))
	o.ClockOffset = binary.LittleEndian.Uint32(r.Get(4))
	copy(o.AFHChannelMap[:], r.Get(10))
	o.LTADDR = uint8(r.GetOne())
	o.NextBroadcastInstant = binary.LittleEndian.Uint32(r.Get(4))
	o.ConnectionlessSlaveBroadcastInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.ServiceData = uint8(r.GetOne())
	return r.Valid()
}

// SynchronizationTrainReceivedEventCallbackType is the type of the callback function for SynchronizationTrainReceivedEvent.
type SynchronizationTrainReceivedEventCallbackType func(*SynchronizationTrainReceivedEvent) *SynchronizationTrainReceivedEvent

// SetSynchronizationTrainReceivedEventCallback configures the callback for SynchronizationTrainReceivedEvent. Passing nil will disable the callback.
func (e *EventHandler) SetSynchronizationTrainReceivedEventCallback(cb SynchronizationTrainReceivedEventCallbackType) error {
	e.cbMutex.Lock()
	e.synchronizationTrainReceivedEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(116, cb != nil);
}
// ConnectionlessSlaveBroadcastReceiveEvent represents the event specified in Section 7.7.69
type ConnectionlessSlaveBroadcastReceiveEvent struct {
	BDADDR bleutil.MacAddr
	LTADDR uint8
	CLK uint32
	Offset uint32
	RXStatus uint8
	Fragment uint8
	DataLength uint8
	Data []byte
}

func (o *ConnectionlessSlaveBroadcastReceiveEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	o.LTADDR = uint8(r.GetOne())
	o.CLK = binary.LittleEndian.Uint32(r.Get(4))
	o.Offset = binary.LittleEndian.Uint32(r.Get(4))
	o.RXStatus = uint8(r.GetOne())
	o.Fragment = uint8(r.GetOne())
	o.DataLength = uint8(r.GetOne())
	o.Data = append(o.Data[:0], r.GetRemainder()...)
	return r.Valid()
}

// ConnectionlessSlaveBroadcastReceiveEventCallbackType is the type of the callback function for ConnectionlessSlaveBroadcastReceiveEvent.
type ConnectionlessSlaveBroadcastReceiveEventCallbackType func(*ConnectionlessSlaveBroadcastReceiveEvent) *ConnectionlessSlaveBroadcastReceiveEvent

// SetConnectionlessSlaveBroadcastReceiveEventCallback configures the callback for ConnectionlessSlaveBroadcastReceiveEvent. Passing nil will disable the callback.
func (e *EventHandler) SetConnectionlessSlaveBroadcastReceiveEventCallback(cb ConnectionlessSlaveBroadcastReceiveEventCallbackType) error {
	e.cbMutex.Lock()
	e.connectionlessSlaveBroadcastReceiveEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(117, cb != nil);
}
// ConnectionlessSlaveBroadcastTimeoutEvent represents the event specified in Section 7.7.70
type ConnectionlessSlaveBroadcastTimeoutEvent struct {
	BDADDR bleutil.MacAddr
	LTADDR uint8
}

func (o *ConnectionlessSlaveBroadcastTimeoutEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.BDADDR.Decode(r.Get(6))
	o.LTADDR = uint8(r.GetOne())
	return r.Valid()
}

// ConnectionlessSlaveBroadcastTimeoutEventCallbackType is the type of the callback function for ConnectionlessSlaveBroadcastTimeoutEvent.
type ConnectionlessSlaveBroadcastTimeoutEventCallbackType func(*ConnectionlessSlaveBroadcastTimeoutEvent) *ConnectionlessSlaveBroadcastTimeoutEvent

// SetConnectionlessSlaveBroadcastTimeoutEventCallback configures the callback for ConnectionlessSlaveBroadcastTimeoutEvent. Passing nil will disable the callback.
func (e *EventHandler) SetConnectionlessSlaveBroadcastTimeoutEventCallback(cb ConnectionlessSlaveBroadcastTimeoutEventCallbackType) error {
	e.cbMutex.Lock()
	e.connectionlessSlaveBroadcastTimeoutEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(118, cb != nil);
}
// TruncatedPageCompleteEvent represents the event specified in Section 7.7.71
type TruncatedPageCompleteEvent struct {
	Status uint8
	BDADDR bleutil.MacAddr
}

func (o *TruncatedPageCompleteEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.BDADDR.Decode(r.Get(6))
	return r.Valid()
}

// TruncatedPageCompleteEventCallbackType is the type of the callback function for TruncatedPageCompleteEvent.
type TruncatedPageCompleteEventCallbackType func(*TruncatedPageCompleteEvent) *TruncatedPageCompleteEvent

// SetTruncatedPageCompleteEventCallback configures the callback for TruncatedPageCompleteEvent. Passing nil will disable the callback.
func (e *EventHandler) SetTruncatedPageCompleteEventCallback(cb TruncatedPageCompleteEventCallbackType) error {
	e.cbMutex.Lock()
	e.truncatedPageCompleteEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(119, cb != nil);
}
// SlavePageResponseTimeoutEvent represents the event specified in Section 7.7.72
type SlavePageResponseTimeoutEvent struct {
}

func (o *SlavePageResponseTimeoutEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	return r.Valid()
}

// SlavePageResponseTimeoutEventCallbackType is the type of the callback function for SlavePageResponseTimeoutEvent.
type SlavePageResponseTimeoutEventCallbackType func(*SlavePageResponseTimeoutEvent) *SlavePageResponseTimeoutEvent

// SetSlavePageResponseTimeoutEventCallback configures the callback for SlavePageResponseTimeoutEvent. Passing nil will disable the callback.
func (e *EventHandler) SetSlavePageResponseTimeoutEventCallback(cb SlavePageResponseTimeoutEventCallbackType) error {
	e.cbMutex.Lock()
	e.slavePageResponseTimeoutEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(120, cb != nil);
}
// ConnectionlessSlaveBroadcastChannelMapChangeEvent represents the event specified in Section 7.7.73
type ConnectionlessSlaveBroadcastChannelMapChangeEvent struct {
	ChannelMap [10]byte
}

func (o *ConnectionlessSlaveBroadcastChannelMapChangeEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	copy(o.ChannelMap[:], r.Get(10))
	return r.Valid()
}

// ConnectionlessSlaveBroadcastChannelMapChangeEventCallbackType is the type of the callback function for ConnectionlessSlaveBroadcastChannelMapChangeEvent.
type ConnectionlessSlaveBroadcastChannelMapChangeEventCallbackType func(*ConnectionlessSlaveBroadcastChannelMapChangeEvent) *ConnectionlessSlaveBroadcastChannelMapChangeEvent

// SetConnectionlessSlaveBroadcastChannelMapChangeEventCallback configures the callback for ConnectionlessSlaveBroadcastChannelMapChangeEvent. Passing nil will disable the callback.
func (e *EventHandler) SetConnectionlessSlaveBroadcastChannelMapChangeEventCallback(cb ConnectionlessSlaveBroadcastChannelMapChangeEventCallbackType) error {
	e.cbMutex.Lock()
	e.connectionlessSlaveBroadcastChannelMapChangeEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(121, cb != nil);
}
// InquiryResponseNotificationEvent represents the event specified in Section 7.7.74
type InquiryResponseNotificationEvent struct {
	LAP uint32
	RSSI uint8
}

func (o *InquiryResponseNotificationEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.LAP = bleutil.DecodeUint24(r.Get(3))
	o.RSSI = uint8(r.GetOne())
	return r.Valid()
}

// InquiryResponseNotificationEventCallbackType is the type of the callback function for InquiryResponseNotificationEvent.
type InquiryResponseNotificationEventCallbackType func(*InquiryResponseNotificationEvent) *InquiryResponseNotificationEvent

// SetInquiryResponseNotificationEventCallback configures the callback for InquiryResponseNotificationEvent. Passing nil will disable the callback.
func (e *EventHandler) SetInquiryResponseNotificationEventCallback(cb InquiryResponseNotificationEventCallbackType) error {
	e.cbMutex.Lock()
	e.inquiryResponseNotificationEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(122, cb != nil);
}
// AuthenticatedPayloadTimeoutExpiredEvent represents the event specified in Section 7.7.75
type AuthenticatedPayloadTimeoutExpiredEvent struct {
	ConnectionHandle uint16
}

func (o *AuthenticatedPayloadTimeoutExpiredEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// AuthenticatedPayloadTimeoutExpiredEventCallbackType is the type of the callback function for AuthenticatedPayloadTimeoutExpiredEvent.
type AuthenticatedPayloadTimeoutExpiredEventCallbackType func(*AuthenticatedPayloadTimeoutExpiredEvent) *AuthenticatedPayloadTimeoutExpiredEvent

// SetAuthenticatedPayloadTimeoutExpiredEventCallback configures the callback for AuthenticatedPayloadTimeoutExpiredEvent. Passing nil will disable the callback.
func (e *EventHandler) SetAuthenticatedPayloadTimeoutExpiredEventCallback(cb AuthenticatedPayloadTimeoutExpiredEventCallbackType) error {
	e.cbMutex.Lock()
	e.authenticatedPayloadTimeoutExpiredEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(123, cb != nil);
}
// SAMStatusChangeEvent represents the event specified in Section 7.7.76
type SAMStatusChangeEvent struct {
	ConnectionHandle uint16
	LocalSAMIndex uint8
	LocalSAMTXAvailability uint8
	LocalSAMRXAvailability uint8
	RemoteSAMIndex uint8
	RemoteSAMTXAvailability uint8
	RemoteSAMRXAvailability uint8
	LEEventMask uint64
	Status uint8
}

func (o *SAMStatusChangeEvent) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LocalSAMIndex = uint8(r.GetOne())
	o.LocalSAMTXAvailability = uint8(r.GetOne())
	o.LocalSAMRXAvailability = uint8(r.GetOne())
	o.RemoteSAMIndex = uint8(r.GetOne())
	o.RemoteSAMTXAvailability = uint8(r.GetOne())
	o.RemoteSAMRXAvailability = uint8(r.GetOne())
	o.LEEventMask = binary.LittleEndian.Uint64(r.Get(8))
	o.Status = uint8(r.GetOne())
	return r.Valid()
}

// SAMStatusChangeEventCallbackType is the type of the callback function for SAMStatusChangeEvent.
type SAMStatusChangeEventCallbackType func(*SAMStatusChangeEvent) *SAMStatusChangeEvent

// SetSAMStatusChangeEventCallback configures the callback for SAMStatusChangeEvent. Passing nil will disable the callback.
func (e *EventHandler) SetSAMStatusChangeEventCallback(cb SAMStatusChangeEventCallbackType) error {
	e.cbMutex.Lock()
	e.sAMStatusChangeEventCallback = cb
	e.cbMutex.Unlock()

	 return e.eventChanged(124, cb != nil);
}


type EventHandler struct {
	logger *logrus.Entry
	hcicmdmgr *hcicmdmgr.CommandManager
	cmds *hcicommands.Commands
	cbMutex sync.RWMutex
	enableMutex sync.Mutex
	eventMask uint64
	eventMask2 uint64
	eventMaskLe uint64

	inquiryCompleteEventCallback InquiryCompleteEventCallbackType
	inquiryCompleteEvent *InquiryCompleteEvent

	inquiryResultEventCallback InquiryResultEventCallbackType
	inquiryResultEvent *InquiryResultEvent

	connectionCompleteEventCallback ConnectionCompleteEventCallbackType
	connectionCompleteEvent *ConnectionCompleteEvent

	connectionRequestEventCallback ConnectionRequestEventCallbackType
	connectionRequestEvent *ConnectionRequestEvent

	disconnectionCompleteEventCallback DisconnectionCompleteEventCallbackType
	disconnectionCompleteEvent *DisconnectionCompleteEvent

	authenticationCompleteEventCallback AuthenticationCompleteEventCallbackType
	authenticationCompleteEvent *AuthenticationCompleteEvent

	remoteNameRequestCompleteEventCallback RemoteNameRequestCompleteEventCallbackType
	remoteNameRequestCompleteEvent *RemoteNameRequestCompleteEvent

	encryptionChangeEventCallback EncryptionChangeEventCallbackType
	encryptionChangeEvent *EncryptionChangeEvent

	changeConnectionLinkKeyCompleteEventCallback ChangeConnectionLinkKeyCompleteEventCallbackType
	changeConnectionLinkKeyCompleteEvent *ChangeConnectionLinkKeyCompleteEvent

	masterLinkKeyCompleteEventCallback MasterLinkKeyCompleteEventCallbackType
	masterLinkKeyCompleteEvent *MasterLinkKeyCompleteEvent

	readRemoteSupportedFeaturesCompleteEventCallback ReadRemoteSupportedFeaturesCompleteEventCallbackType
	readRemoteSupportedFeaturesCompleteEvent *ReadRemoteSupportedFeaturesCompleteEvent

	readRemoteVersionInformationCompleteEventCallback ReadRemoteVersionInformationCompleteEventCallbackType
	readRemoteVersionInformationCompleteEvent *ReadRemoteVersionInformationCompleteEvent

	qoSSetupCompleteEventCallback QoSSetupCompleteEventCallbackType
	qoSSetupCompleteEvent *QoSSetupCompleteEvent

	commandCompleteEventCallback CommandCompleteEventCallbackType
	commandCompleteEvent *CommandCompleteEvent

	commandStatusEventCallback CommandStatusEventCallbackType
	commandStatusEvent *CommandStatusEvent

	hardwareErrorEventCallback HardwareErrorEventCallbackType
	hardwareErrorEvent *HardwareErrorEvent

	flushOccurredEventCallback FlushOccurredEventCallbackType
	flushOccurredEvent *FlushOccurredEvent

	roleChangeEventCallback RoleChangeEventCallbackType
	roleChangeEvent *RoleChangeEvent

	numberOfCompletedPacketsEventCallback NumberOfCompletedPacketsEventCallbackType
	numberOfCompletedPacketsEvent *NumberOfCompletedPacketsEvent

	modeChangeEventCallback ModeChangeEventCallbackType
	modeChangeEvent *ModeChangeEvent

	returnLinkKeysEventCallback ReturnLinkKeysEventCallbackType
	returnLinkKeysEvent *ReturnLinkKeysEvent

	pINCodeRequestEventCallback PINCodeRequestEventCallbackType
	pINCodeRequestEvent *PINCodeRequestEvent

	linkKeyRequestEventCallback LinkKeyRequestEventCallbackType
	linkKeyRequestEvent *LinkKeyRequestEvent

	linkKeyNotificationEventCallback LinkKeyNotificationEventCallbackType
	linkKeyNotificationEvent *LinkKeyNotificationEvent

	loopbackCommandEventCallback LoopbackCommandEventCallbackType
	loopbackCommandEvent *LoopbackCommandEvent

	dataBufferOverflowEventCallback DataBufferOverflowEventCallbackType
	dataBufferOverflowEvent *DataBufferOverflowEvent

	maxSlotsChangeEventCallback MaxSlotsChangeEventCallbackType
	maxSlotsChangeEvent *MaxSlotsChangeEvent

	readClockOffsetCompleteEventCallback ReadClockOffsetCompleteEventCallbackType
	readClockOffsetCompleteEvent *ReadClockOffsetCompleteEvent

	connectionPacketTypeChangedEventCallback ConnectionPacketTypeChangedEventCallbackType
	connectionPacketTypeChangedEvent *ConnectionPacketTypeChangedEvent

	qoSViolationEventCallback QoSViolationEventCallbackType
	qoSViolationEvent *QoSViolationEvent

	pageScanRepetitionModeChangeEventCallback PageScanRepetitionModeChangeEventCallbackType
	pageScanRepetitionModeChangeEvent *PageScanRepetitionModeChangeEvent

	flowSpecificationCompleteEventCallback FlowSpecificationCompleteEventCallbackType
	flowSpecificationCompleteEvent *FlowSpecificationCompleteEvent

	inquiryResultwithRSSIEventCallback InquiryResultwithRSSIEventCallbackType
	inquiryResultwithRSSIEvent *InquiryResultwithRSSIEvent

	readRemoteExtendedFeaturesCompleteEventCallback ReadRemoteExtendedFeaturesCompleteEventCallbackType
	readRemoteExtendedFeaturesCompleteEvent *ReadRemoteExtendedFeaturesCompleteEvent

	synchronousConnectionCompleteEventCallback SynchronousConnectionCompleteEventCallbackType
	synchronousConnectionCompleteEvent *SynchronousConnectionCompleteEvent

	synchronousConnectionChangedEventCallback SynchronousConnectionChangedEventCallbackType
	synchronousConnectionChangedEvent *SynchronousConnectionChangedEvent

	sniffSubratingEventCallback SniffSubratingEventCallbackType
	sniffSubratingEvent *SniffSubratingEvent

	extendedInquiryResultEventCallback ExtendedInquiryResultEventCallbackType
	extendedInquiryResultEvent *ExtendedInquiryResultEvent

	encryptionKeyRefreshCompleteEventCallback EncryptionKeyRefreshCompleteEventCallbackType
	encryptionKeyRefreshCompleteEvent *EncryptionKeyRefreshCompleteEvent

	iOCapabilityRequestEventCallback IOCapabilityRequestEventCallbackType
	iOCapabilityRequestEvent *IOCapabilityRequestEvent

	iOCapabilityResponseEventCallback IOCapabilityResponseEventCallbackType
	iOCapabilityResponseEvent *IOCapabilityResponseEvent

	userConfirmationRequestEventCallback UserConfirmationRequestEventCallbackType
	userConfirmationRequestEvent *UserConfirmationRequestEvent

	userPasskeyRequestEventCallback UserPasskeyRequestEventCallbackType
	userPasskeyRequestEvent *UserPasskeyRequestEvent

	remoteOOBDataRequestEventCallback RemoteOOBDataRequestEventCallbackType
	remoteOOBDataRequestEvent *RemoteOOBDataRequestEvent

	simplePairingCompleteEventCallback SimplePairingCompleteEventCallbackType
	simplePairingCompleteEvent *SimplePairingCompleteEvent

	linkSupervisionTimeoutChangedEventCallback LinkSupervisionTimeoutChangedEventCallbackType
	linkSupervisionTimeoutChangedEvent *LinkSupervisionTimeoutChangedEvent

	enhancedFlushCompleteEventCallback EnhancedFlushCompleteEventCallbackType
	enhancedFlushCompleteEvent *EnhancedFlushCompleteEvent

	userPasskeyNotificationEventCallback UserPasskeyNotificationEventCallbackType
	userPasskeyNotificationEvent *UserPasskeyNotificationEvent

	keypressNotificationEventCallback KeypressNotificationEventCallbackType
	keypressNotificationEvent *KeypressNotificationEvent

	remoteHostSupportedFeaturesNotificationEventCallback RemoteHostSupportedFeaturesNotificationEventCallbackType
	remoteHostSupportedFeaturesNotificationEvent *RemoteHostSupportedFeaturesNotificationEvent

	physicalLinkCompleteEventCallback PhysicalLinkCompleteEventCallbackType
	physicalLinkCompleteEvent *PhysicalLinkCompleteEvent

	channelSelectedEventCallback ChannelSelectedEventCallbackType
	channelSelectedEvent *ChannelSelectedEvent

	disconnectionPhysicalLinkCompleteEventCallback DisconnectionPhysicalLinkCompleteEventCallbackType
	disconnectionPhysicalLinkCompleteEvent *DisconnectionPhysicalLinkCompleteEvent

	physicalLinkLossEarlyWarningEventCallback PhysicalLinkLossEarlyWarningEventCallbackType
	physicalLinkLossEarlyWarningEvent *PhysicalLinkLossEarlyWarningEvent

	physicalLinkRecoveryEventCallback PhysicalLinkRecoveryEventCallbackType
	physicalLinkRecoveryEvent *PhysicalLinkRecoveryEvent

	logicalLinkCompleteEventCallback LogicalLinkCompleteEventCallbackType
	logicalLinkCompleteEvent *LogicalLinkCompleteEvent

	disconnectionLogicalLinkCompleteEventCallback DisconnectionLogicalLinkCompleteEventCallbackType
	disconnectionLogicalLinkCompleteEvent *DisconnectionLogicalLinkCompleteEvent

	flowSpecModifyCompleteEventCallback FlowSpecModifyCompleteEventCallbackType
	flowSpecModifyCompleteEvent *FlowSpecModifyCompleteEvent

	numberOfCompletedDataBlocksEventCallback NumberOfCompletedDataBlocksEventCallbackType
	numberOfCompletedDataBlocksEvent *NumberOfCompletedDataBlocksEvent

	shortRangeModeChangeCompleteEventCallback ShortRangeModeChangeCompleteEventCallbackType
	shortRangeModeChangeCompleteEvent *ShortRangeModeChangeCompleteEvent

	aMPStatusChangeEventCallback AMPStatusChangeEventCallbackType
	aMPStatusChangeEvent *AMPStatusChangeEvent

	aMPStartTestEventCallback AMPStartTestEventCallbackType
	aMPStartTestEvent *AMPStartTestEvent

	aMPTestEndEventCallback AMPTestEndEventCallbackType
	aMPTestEndEvent *AMPTestEndEvent

	aMPReceiverReportEventCallback AMPReceiverReportEventCallbackType
	aMPReceiverReportEvent *AMPReceiverReportEvent

	lEConnectionCompleteEventCallback LEConnectionCompleteEventCallbackType
	lEConnectionCompleteEvent *LEConnectionCompleteEvent

	lEAdvertisingReportEventCallback LEAdvertisingReportEventCallbackType
	lEAdvertisingReportEvent *LEAdvertisingReportEvent

	lEConnectionUpdateCompleteEventCallback LEConnectionUpdateCompleteEventCallbackType
	lEConnectionUpdateCompleteEvent *LEConnectionUpdateCompleteEvent

	lEReadRemoteFeaturesCompleteEventCallback LEReadRemoteFeaturesCompleteEventCallbackType
	lEReadRemoteFeaturesCompleteEvent *LEReadRemoteFeaturesCompleteEvent

	lELongTermKeyRequestEventCallback LELongTermKeyRequestEventCallbackType
	lELongTermKeyRequestEvent *LELongTermKeyRequestEvent

	lERemoteConnectionParameterRequestEventCallback LERemoteConnectionParameterRequestEventCallbackType
	lERemoteConnectionParameterRequestEvent *LERemoteConnectionParameterRequestEvent

	lEDataLengthChangeEventCallback LEDataLengthChangeEventCallbackType
	lEDataLengthChangeEvent *LEDataLengthChangeEvent

	lEReadLocalP256PublicKeyCompleteEventCallback LEReadLocalP256PublicKeyCompleteEventCallbackType
	lEReadLocalP256PublicKeyCompleteEvent *LEReadLocalP256PublicKeyCompleteEvent

	lEGenerateDHKeyCompleteEventCallback LEGenerateDHKeyCompleteEventCallbackType
	lEGenerateDHKeyCompleteEvent *LEGenerateDHKeyCompleteEvent

	lEEnhancedConnectionCompleteEventCallback LEEnhancedConnectionCompleteEventCallbackType
	lEEnhancedConnectionCompleteEvent *LEEnhancedConnectionCompleteEvent

	lEDirectedAdvertisingReportEventCallback LEDirectedAdvertisingReportEventCallbackType
	lEDirectedAdvertisingReportEvent *LEDirectedAdvertisingReportEvent

	lEPHYUpdateCompleteEventCallback LEPHYUpdateCompleteEventCallbackType
	lEPHYUpdateCompleteEvent *LEPHYUpdateCompleteEvent

	lEExtendedAdvertisingReportEventCallback LEExtendedAdvertisingReportEventCallbackType
	lEExtendedAdvertisingReportEvent *LEExtendedAdvertisingReportEvent

	lEPeriodicAdvertisingSyncEstablishedEventCallback LEPeriodicAdvertisingSyncEstablishedEventCallbackType
	lEPeriodicAdvertisingSyncEstablishedEvent *LEPeriodicAdvertisingSyncEstablishedEvent

	lEPeriodicAdvertisingReportEventCallback LEPeriodicAdvertisingReportEventCallbackType
	lEPeriodicAdvertisingReportEvent *LEPeriodicAdvertisingReportEvent

	lEPeriodicAdvertisingSyncLostEventCallback LEPeriodicAdvertisingSyncLostEventCallbackType
	lEPeriodicAdvertisingSyncLostEvent *LEPeriodicAdvertisingSyncLostEvent

	lEScanTimeoutEventCallback LEScanTimeoutEventCallbackType
	lEScanTimeoutEvent *LEScanTimeoutEvent

	lEAdvertisingSetTerminatedEventCallback LEAdvertisingSetTerminatedEventCallbackType
	lEAdvertisingSetTerminatedEvent *LEAdvertisingSetTerminatedEvent

	lEScanRequestReceivedEventCallback LEScanRequestReceivedEventCallbackType
	lEScanRequestReceivedEvent *LEScanRequestReceivedEvent

	lEChannelSelectionAlgorithmEventCallback LEChannelSelectionAlgorithmEventCallbackType
	lEChannelSelectionAlgorithmEvent *LEChannelSelectionAlgorithmEvent

	lEConnectionlessIQReportEventCallback LEConnectionlessIQReportEventCallbackType
	lEConnectionlessIQReportEvent *LEConnectionlessIQReportEvent

	lEConnectionIQReportEventCallback LEConnectionIQReportEventCallbackType
	lEConnectionIQReportEvent *LEConnectionIQReportEvent

	lECTERequestFailedEventCallback LECTERequestFailedEventCallbackType
	lECTERequestFailedEvent *LECTERequestFailedEvent

	lEPeriodicAdvertisingSyncTransferReceivedEventCallback LEPeriodicAdvertisingSyncTransferReceivedEventCallbackType
	lEPeriodicAdvertisingSyncTransferReceivedEvent *LEPeriodicAdvertisingSyncTransferReceivedEvent

	lECISEstablishedEventCallback LECISEstablishedEventCallbackType
	lECISEstablishedEvent *LECISEstablishedEvent

	lECISRequestEventCallback LECISRequestEventCallbackType
	lECISRequestEvent *LECISRequestEvent

	lECreateBIGCompleteEventCallback LECreateBIGCompleteEventCallbackType
	lECreateBIGCompleteEvent *LECreateBIGCompleteEvent

	lETerminateBIGCompleteEventCallback LETerminateBIGCompleteEventCallbackType
	lETerminateBIGCompleteEvent *LETerminateBIGCompleteEvent

	lEBIGSyncEstablishedEventCallback LEBIGSyncEstablishedEventCallbackType
	lEBIGSyncEstablishedEvent *LEBIGSyncEstablishedEvent

	lEBIGSyncLostEventCallback LEBIGSyncLostEventCallbackType
	lEBIGSyncLostEvent *LEBIGSyncLostEvent

	lERequestPeerSCACompleteEventCallback LERequestPeerSCACompleteEventCallbackType
	lERequestPeerSCACompleteEvent *LERequestPeerSCACompleteEvent

	lEPathLossThresholdEventCallback LEPathLossThresholdEventCallbackType
	lEPathLossThresholdEvent *LEPathLossThresholdEvent

	lETransmitPowerReportingEventCallback LETransmitPowerReportingEventCallbackType
	lETransmitPowerReportingEvent *LETransmitPowerReportingEvent

	lEBIGInfoAdvertisingReportEventCallback LEBIGInfoAdvertisingReportEventCallbackType
	lEBIGInfoAdvertisingReportEvent *LEBIGInfoAdvertisingReportEvent

	triggeredClockCaptureEventCallback TriggeredClockCaptureEventCallbackType
	triggeredClockCaptureEvent *TriggeredClockCaptureEvent

	synchronizationTrainCompleteEventCallback SynchronizationTrainCompleteEventCallbackType
	synchronizationTrainCompleteEvent *SynchronizationTrainCompleteEvent

	synchronizationTrainReceivedEventCallback SynchronizationTrainReceivedEventCallbackType
	synchronizationTrainReceivedEvent *SynchronizationTrainReceivedEvent

	connectionlessSlaveBroadcastReceiveEventCallback ConnectionlessSlaveBroadcastReceiveEventCallbackType
	connectionlessSlaveBroadcastReceiveEvent *ConnectionlessSlaveBroadcastReceiveEvent

	connectionlessSlaveBroadcastTimeoutEventCallback ConnectionlessSlaveBroadcastTimeoutEventCallbackType
	connectionlessSlaveBroadcastTimeoutEvent *ConnectionlessSlaveBroadcastTimeoutEvent

	truncatedPageCompleteEventCallback TruncatedPageCompleteEventCallbackType
	truncatedPageCompleteEvent *TruncatedPageCompleteEvent

	slavePageResponseTimeoutEventCallback SlavePageResponseTimeoutEventCallbackType
	slavePageResponseTimeoutEvent *SlavePageResponseTimeoutEvent

	connectionlessSlaveBroadcastChannelMapChangeEventCallback ConnectionlessSlaveBroadcastChannelMapChangeEventCallbackType
	connectionlessSlaveBroadcastChannelMapChangeEvent *ConnectionlessSlaveBroadcastChannelMapChangeEvent

	inquiryResponseNotificationEventCallback InquiryResponseNotificationEventCallbackType
	inquiryResponseNotificationEvent *InquiryResponseNotificationEvent

	authenticatedPayloadTimeoutExpiredEventCallback AuthenticatedPayloadTimeoutExpiredEventCallbackType
	authenticatedPayloadTimeoutExpiredEvent *AuthenticatedPayloadTimeoutExpiredEvent

	sAMStatusChangeEventCallback SAMStatusChangeEventCallbackType
	sAMStatusChangeEvent *SAMStatusChangeEvent

}

func (e *EventHandler) handleEventInternal(eventCode uint16, params []byte) {
	switch eventCode {
	case 0x0100:
		e.cbMutex.RLock()
		cb := e.inquiryCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.inquiryCompleteEvent == nil {
				e.inquiryCompleteEvent = &InquiryCompleteEvent{}
			}

			if e.inquiryCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.inquiryCompleteEvent).Debug("InquiryCompleteEvent decoded")
				}
				e.inquiryCompleteEvent = cb(e.inquiryCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("InquiryCompleteEvent has no callback")
			}
		}
	case 0x0200:
		e.cbMutex.RLock()
		cb := e.inquiryResultEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.inquiryResultEvent == nil {
				e.inquiryResultEvent = &InquiryResultEvent{}
			}

			if e.inquiryResultEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.inquiryResultEvent).Debug("InquiryResultEvent decoded")
				}
				e.inquiryResultEvent = cb(e.inquiryResultEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("InquiryResultEvent has no callback")
			}
		}
	case 0x0300:
		e.cbMutex.RLock()
		cb := e.connectionCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.connectionCompleteEvent == nil {
				e.connectionCompleteEvent = &ConnectionCompleteEvent{}
			}

			if e.connectionCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.connectionCompleteEvent).Debug("ConnectionCompleteEvent decoded")
				}
				e.connectionCompleteEvent = cb(e.connectionCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ConnectionCompleteEvent has no callback")
			}
		}
	case 0x0400:
		e.cbMutex.RLock()
		cb := e.connectionRequestEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.connectionRequestEvent == nil {
				e.connectionRequestEvent = &ConnectionRequestEvent{}
			}

			if e.connectionRequestEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.connectionRequestEvent).Debug("ConnectionRequestEvent decoded")
				}
				e.connectionRequestEvent = cb(e.connectionRequestEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ConnectionRequestEvent has no callback")
			}
		}
	case 0x0500:
		e.cbMutex.RLock()
		cb := e.disconnectionCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.disconnectionCompleteEvent == nil {
				e.disconnectionCompleteEvent = &DisconnectionCompleteEvent{}
			}

			if e.disconnectionCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.disconnectionCompleteEvent).Debug("DisconnectionCompleteEvent decoded")
				}
				e.disconnectionCompleteEvent = cb(e.disconnectionCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("DisconnectionCompleteEvent has no callback")
			}
		}
	case 0x0600:
		e.cbMutex.RLock()
		cb := e.authenticationCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.authenticationCompleteEvent == nil {
				e.authenticationCompleteEvent = &AuthenticationCompleteEvent{}
			}

			if e.authenticationCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.authenticationCompleteEvent).Debug("AuthenticationCompleteEvent decoded")
				}
				e.authenticationCompleteEvent = cb(e.authenticationCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("AuthenticationCompleteEvent has no callback")
			}
		}
	case 0x0700:
		e.cbMutex.RLock()
		cb := e.remoteNameRequestCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.remoteNameRequestCompleteEvent == nil {
				e.remoteNameRequestCompleteEvent = &RemoteNameRequestCompleteEvent{}
			}

			if e.remoteNameRequestCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.remoteNameRequestCompleteEvent).Debug("RemoteNameRequestCompleteEvent decoded")
				}
				e.remoteNameRequestCompleteEvent = cb(e.remoteNameRequestCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("RemoteNameRequestCompleteEvent has no callback")
			}
		}
	case 0x0800:
		e.cbMutex.RLock()
		cb := e.encryptionChangeEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.encryptionChangeEvent == nil {
				e.encryptionChangeEvent = &EncryptionChangeEvent{}
			}

			if e.encryptionChangeEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.encryptionChangeEvent).Debug("EncryptionChangeEvent decoded")
				}
				e.encryptionChangeEvent = cb(e.encryptionChangeEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("EncryptionChangeEvent has no callback")
			}
		}
	case 0x0900:
		e.cbMutex.RLock()
		cb := e.changeConnectionLinkKeyCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.changeConnectionLinkKeyCompleteEvent == nil {
				e.changeConnectionLinkKeyCompleteEvent = &ChangeConnectionLinkKeyCompleteEvent{}
			}

			if e.changeConnectionLinkKeyCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.changeConnectionLinkKeyCompleteEvent).Debug("ChangeConnectionLinkKeyCompleteEvent decoded")
				}
				e.changeConnectionLinkKeyCompleteEvent = cb(e.changeConnectionLinkKeyCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ChangeConnectionLinkKeyCompleteEvent has no callback")
			}
		}
	case 0x0A00:
		e.cbMutex.RLock()
		cb := e.masterLinkKeyCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.masterLinkKeyCompleteEvent == nil {
				e.masterLinkKeyCompleteEvent = &MasterLinkKeyCompleteEvent{}
			}

			if e.masterLinkKeyCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.masterLinkKeyCompleteEvent).Debug("MasterLinkKeyCompleteEvent decoded")
				}
				e.masterLinkKeyCompleteEvent = cb(e.masterLinkKeyCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("MasterLinkKeyCompleteEvent has no callback")
			}
		}
	case 0x0B00:
		e.cbMutex.RLock()
		cb := e.readRemoteSupportedFeaturesCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.readRemoteSupportedFeaturesCompleteEvent == nil {
				e.readRemoteSupportedFeaturesCompleteEvent = &ReadRemoteSupportedFeaturesCompleteEvent{}
			}

			if e.readRemoteSupportedFeaturesCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.readRemoteSupportedFeaturesCompleteEvent).Debug("ReadRemoteSupportedFeaturesCompleteEvent decoded")
				}
				e.readRemoteSupportedFeaturesCompleteEvent = cb(e.readRemoteSupportedFeaturesCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ReadRemoteSupportedFeaturesCompleteEvent has no callback")
			}
		}
	case 0x0C00:
		e.cbMutex.RLock()
		cb := e.readRemoteVersionInformationCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.readRemoteVersionInformationCompleteEvent == nil {
				e.readRemoteVersionInformationCompleteEvent = &ReadRemoteVersionInformationCompleteEvent{}
			}

			if e.readRemoteVersionInformationCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.readRemoteVersionInformationCompleteEvent).Debug("ReadRemoteVersionInformationCompleteEvent decoded")
				}
				e.readRemoteVersionInformationCompleteEvent = cb(e.readRemoteVersionInformationCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ReadRemoteVersionInformationCompleteEvent has no callback")
			}
		}
	case 0x0D00:
		e.cbMutex.RLock()
		cb := e.qoSSetupCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.qoSSetupCompleteEvent == nil {
				e.qoSSetupCompleteEvent = &QoSSetupCompleteEvent{}
			}

			if e.qoSSetupCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.qoSSetupCompleteEvent).Debug("QoSSetupCompleteEvent decoded")
				}
				e.qoSSetupCompleteEvent = cb(e.qoSSetupCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("QoSSetupCompleteEvent has no callback")
			}
		}
	case 0x0E00:
		e.cbMutex.RLock()
		cb := e.commandCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.commandCompleteEvent == nil {
				e.commandCompleteEvent = &CommandCompleteEvent{}
			}

			if e.commandCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
					e.logger.WithField("0data", e.commandCompleteEvent).Trace("CommandCompleteEvent decoded")
				}
				e.commandCompleteEvent = cb(e.commandCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("CommandCompleteEvent has no callback")
			}
		}
	case 0x0F00:
		e.cbMutex.RLock()
		cb := e.commandStatusEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.commandStatusEvent == nil {
				e.commandStatusEvent = &CommandStatusEvent{}
			}

			if e.commandStatusEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
					e.logger.WithField("0data", e.commandStatusEvent).Trace("CommandStatusEvent decoded")
				}
				e.commandStatusEvent = cb(e.commandStatusEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("CommandStatusEvent has no callback")
			}
		}
	case 0x1000:
		e.cbMutex.RLock()
		cb := e.hardwareErrorEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.hardwareErrorEvent == nil {
				e.hardwareErrorEvent = &HardwareErrorEvent{}
			}

			if e.hardwareErrorEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.hardwareErrorEvent).Debug("HardwareErrorEvent decoded")
				}
				e.hardwareErrorEvent = cb(e.hardwareErrorEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("HardwareErrorEvent has no callback")
			}
		}
	case 0x1100:
		e.cbMutex.RLock()
		cb := e.flushOccurredEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.flushOccurredEvent == nil {
				e.flushOccurredEvent = &FlushOccurredEvent{}
			}

			if e.flushOccurredEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.flushOccurredEvent).Debug("FlushOccurredEvent decoded")
				}
				e.flushOccurredEvent = cb(e.flushOccurredEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("FlushOccurredEvent has no callback")
			}
		}
	case 0x1200:
		e.cbMutex.RLock()
		cb := e.roleChangeEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.roleChangeEvent == nil {
				e.roleChangeEvent = &RoleChangeEvent{}
			}

			if e.roleChangeEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.roleChangeEvent).Debug("RoleChangeEvent decoded")
				}
				e.roleChangeEvent = cb(e.roleChangeEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("RoleChangeEvent has no callback")
			}
		}
	case 0x1300:
		e.cbMutex.RLock()
		cb := e.numberOfCompletedPacketsEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.numberOfCompletedPacketsEvent == nil {
				e.numberOfCompletedPacketsEvent = &NumberOfCompletedPacketsEvent{}
			}

			if e.numberOfCompletedPacketsEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
					e.logger.WithField("0data", e.numberOfCompletedPacketsEvent).Trace("NumberOfCompletedPacketsEvent decoded")
				}
				e.numberOfCompletedPacketsEvent = cb(e.numberOfCompletedPacketsEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("NumberOfCompletedPacketsEvent has no callback")
			}
		}
	case 0x1400:
		e.cbMutex.RLock()
		cb := e.modeChangeEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.modeChangeEvent == nil {
				e.modeChangeEvent = &ModeChangeEvent{}
			}

			if e.modeChangeEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.modeChangeEvent).Debug("ModeChangeEvent decoded")
				}
				e.modeChangeEvent = cb(e.modeChangeEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ModeChangeEvent has no callback")
			}
		}
	case 0x1500:
		e.cbMutex.RLock()
		cb := e.returnLinkKeysEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.returnLinkKeysEvent == nil {
				e.returnLinkKeysEvent = &ReturnLinkKeysEvent{}
			}

			if e.returnLinkKeysEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.returnLinkKeysEvent).Debug("ReturnLinkKeysEvent decoded")
				}
				e.returnLinkKeysEvent = cb(e.returnLinkKeysEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ReturnLinkKeysEvent has no callback")
			}
		}
	case 0x1600:
		e.cbMutex.RLock()
		cb := e.pINCodeRequestEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.pINCodeRequestEvent == nil {
				e.pINCodeRequestEvent = &PINCodeRequestEvent{}
			}

			if e.pINCodeRequestEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.pINCodeRequestEvent).Debug("PINCodeRequestEvent decoded")
				}
				e.pINCodeRequestEvent = cb(e.pINCodeRequestEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("PINCodeRequestEvent has no callback")
			}
		}
	case 0x1700:
		e.cbMutex.RLock()
		cb := e.linkKeyRequestEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.linkKeyRequestEvent == nil {
				e.linkKeyRequestEvent = &LinkKeyRequestEvent{}
			}

			if e.linkKeyRequestEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.linkKeyRequestEvent).Debug("LinkKeyRequestEvent decoded")
				}
				e.linkKeyRequestEvent = cb(e.linkKeyRequestEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LinkKeyRequestEvent has no callback")
			}
		}
	case 0x1800:
		e.cbMutex.RLock()
		cb := e.linkKeyNotificationEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.linkKeyNotificationEvent == nil {
				e.linkKeyNotificationEvent = &LinkKeyNotificationEvent{}
			}

			if e.linkKeyNotificationEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.linkKeyNotificationEvent).Debug("LinkKeyNotificationEvent decoded")
				}
				e.linkKeyNotificationEvent = cb(e.linkKeyNotificationEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LinkKeyNotificationEvent has no callback")
			}
		}
	case 0x1900:
		e.cbMutex.RLock()
		cb := e.loopbackCommandEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.loopbackCommandEvent == nil {
				e.loopbackCommandEvent = &LoopbackCommandEvent{}
			}

			if e.loopbackCommandEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.loopbackCommandEvent).Debug("LoopbackCommandEvent decoded")
				}
				e.loopbackCommandEvent = cb(e.loopbackCommandEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LoopbackCommandEvent has no callback")
			}
		}
	case 0x1A00:
		e.cbMutex.RLock()
		cb := e.dataBufferOverflowEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.dataBufferOverflowEvent == nil {
				e.dataBufferOverflowEvent = &DataBufferOverflowEvent{}
			}

			if e.dataBufferOverflowEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.dataBufferOverflowEvent).Debug("DataBufferOverflowEvent decoded")
				}
				e.dataBufferOverflowEvent = cb(e.dataBufferOverflowEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("DataBufferOverflowEvent has no callback")
			}
		}
	case 0x1B00:
		e.cbMutex.RLock()
		cb := e.maxSlotsChangeEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.maxSlotsChangeEvent == nil {
				e.maxSlotsChangeEvent = &MaxSlotsChangeEvent{}
			}

			if e.maxSlotsChangeEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.maxSlotsChangeEvent).Debug("MaxSlotsChangeEvent decoded")
				}
				e.maxSlotsChangeEvent = cb(e.maxSlotsChangeEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("MaxSlotsChangeEvent has no callback")
			}
		}
	case 0x1C00:
		e.cbMutex.RLock()
		cb := e.readClockOffsetCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.readClockOffsetCompleteEvent == nil {
				e.readClockOffsetCompleteEvent = &ReadClockOffsetCompleteEvent{}
			}

			if e.readClockOffsetCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.readClockOffsetCompleteEvent).Debug("ReadClockOffsetCompleteEvent decoded")
				}
				e.readClockOffsetCompleteEvent = cb(e.readClockOffsetCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ReadClockOffsetCompleteEvent has no callback")
			}
		}
	case 0x1D00:
		e.cbMutex.RLock()
		cb := e.connectionPacketTypeChangedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.connectionPacketTypeChangedEvent == nil {
				e.connectionPacketTypeChangedEvent = &ConnectionPacketTypeChangedEvent{}
			}

			if e.connectionPacketTypeChangedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.connectionPacketTypeChangedEvent).Debug("ConnectionPacketTypeChangedEvent decoded")
				}
				e.connectionPacketTypeChangedEvent = cb(e.connectionPacketTypeChangedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ConnectionPacketTypeChangedEvent has no callback")
			}
		}
	case 0x1E00:
		e.cbMutex.RLock()
		cb := e.qoSViolationEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.qoSViolationEvent == nil {
				e.qoSViolationEvent = &QoSViolationEvent{}
			}

			if e.qoSViolationEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.qoSViolationEvent).Debug("QoSViolationEvent decoded")
				}
				e.qoSViolationEvent = cb(e.qoSViolationEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("QoSViolationEvent has no callback")
			}
		}
	case 0x2000:
		e.cbMutex.RLock()
		cb := e.pageScanRepetitionModeChangeEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.pageScanRepetitionModeChangeEvent == nil {
				e.pageScanRepetitionModeChangeEvent = &PageScanRepetitionModeChangeEvent{}
			}

			if e.pageScanRepetitionModeChangeEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.pageScanRepetitionModeChangeEvent).Debug("PageScanRepetitionModeChangeEvent decoded")
				}
				e.pageScanRepetitionModeChangeEvent = cb(e.pageScanRepetitionModeChangeEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("PageScanRepetitionModeChangeEvent has no callback")
			}
		}
	case 0x2100:
		e.cbMutex.RLock()
		cb := e.flowSpecificationCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.flowSpecificationCompleteEvent == nil {
				e.flowSpecificationCompleteEvent = &FlowSpecificationCompleteEvent{}
			}

			if e.flowSpecificationCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.flowSpecificationCompleteEvent).Debug("FlowSpecificationCompleteEvent decoded")
				}
				e.flowSpecificationCompleteEvent = cb(e.flowSpecificationCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("FlowSpecificationCompleteEvent has no callback")
			}
		}
	case 0x2200:
		e.cbMutex.RLock()
		cb := e.inquiryResultwithRSSIEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.inquiryResultwithRSSIEvent == nil {
				e.inquiryResultwithRSSIEvent = &InquiryResultwithRSSIEvent{}
			}

			if e.inquiryResultwithRSSIEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.inquiryResultwithRSSIEvent).Debug("InquiryResultwithRSSIEvent decoded")
				}
				e.inquiryResultwithRSSIEvent = cb(e.inquiryResultwithRSSIEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("InquiryResultwithRSSIEvent has no callback")
			}
		}
	case 0x2300:
		e.cbMutex.RLock()
		cb := e.readRemoteExtendedFeaturesCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.readRemoteExtendedFeaturesCompleteEvent == nil {
				e.readRemoteExtendedFeaturesCompleteEvent = &ReadRemoteExtendedFeaturesCompleteEvent{}
			}

			if e.readRemoteExtendedFeaturesCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.readRemoteExtendedFeaturesCompleteEvent).Debug("ReadRemoteExtendedFeaturesCompleteEvent decoded")
				}
				e.readRemoteExtendedFeaturesCompleteEvent = cb(e.readRemoteExtendedFeaturesCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ReadRemoteExtendedFeaturesCompleteEvent has no callback")
			}
		}
	case 0x2C00:
		e.cbMutex.RLock()
		cb := e.synchronousConnectionCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.synchronousConnectionCompleteEvent == nil {
				e.synchronousConnectionCompleteEvent = &SynchronousConnectionCompleteEvent{}
			}

			if e.synchronousConnectionCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.synchronousConnectionCompleteEvent).Debug("SynchronousConnectionCompleteEvent decoded")
				}
				e.synchronousConnectionCompleteEvent = cb(e.synchronousConnectionCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("SynchronousConnectionCompleteEvent has no callback")
			}
		}
	case 0x2D00:
		e.cbMutex.RLock()
		cb := e.synchronousConnectionChangedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.synchronousConnectionChangedEvent == nil {
				e.synchronousConnectionChangedEvent = &SynchronousConnectionChangedEvent{}
			}

			if e.synchronousConnectionChangedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.synchronousConnectionChangedEvent).Debug("SynchronousConnectionChangedEvent decoded")
				}
				e.synchronousConnectionChangedEvent = cb(e.synchronousConnectionChangedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("SynchronousConnectionChangedEvent has no callback")
			}
		}
	case 0x2E00:
		e.cbMutex.RLock()
		cb := e.sniffSubratingEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.sniffSubratingEvent == nil {
				e.sniffSubratingEvent = &SniffSubratingEvent{}
			}

			if e.sniffSubratingEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.sniffSubratingEvent).Debug("SniffSubratingEvent decoded")
				}
				e.sniffSubratingEvent = cb(e.sniffSubratingEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("SniffSubratingEvent has no callback")
			}
		}
	case 0x2F00:
		e.cbMutex.RLock()
		cb := e.extendedInquiryResultEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.extendedInquiryResultEvent == nil {
				e.extendedInquiryResultEvent = &ExtendedInquiryResultEvent{}
			}

			if e.extendedInquiryResultEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.extendedInquiryResultEvent).Debug("ExtendedInquiryResultEvent decoded")
				}
				e.extendedInquiryResultEvent = cb(e.extendedInquiryResultEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ExtendedInquiryResultEvent has no callback")
			}
		}
	case 0x3000:
		e.cbMutex.RLock()
		cb := e.encryptionKeyRefreshCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.encryptionKeyRefreshCompleteEvent == nil {
				e.encryptionKeyRefreshCompleteEvent = &EncryptionKeyRefreshCompleteEvent{}
			}

			if e.encryptionKeyRefreshCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.encryptionKeyRefreshCompleteEvent).Debug("EncryptionKeyRefreshCompleteEvent decoded")
				}
				e.encryptionKeyRefreshCompleteEvent = cb(e.encryptionKeyRefreshCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("EncryptionKeyRefreshCompleteEvent has no callback")
			}
		}
	case 0x3100:
		e.cbMutex.RLock()
		cb := e.iOCapabilityRequestEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.iOCapabilityRequestEvent == nil {
				e.iOCapabilityRequestEvent = &IOCapabilityRequestEvent{}
			}

			if e.iOCapabilityRequestEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.iOCapabilityRequestEvent).Debug("IOCapabilityRequestEvent decoded")
				}
				e.iOCapabilityRequestEvent = cb(e.iOCapabilityRequestEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("IOCapabilityRequestEvent has no callback")
			}
		}
	case 0x3200:
		e.cbMutex.RLock()
		cb := e.iOCapabilityResponseEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.iOCapabilityResponseEvent == nil {
				e.iOCapabilityResponseEvent = &IOCapabilityResponseEvent{}
			}

			if e.iOCapabilityResponseEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.iOCapabilityResponseEvent).Debug("IOCapabilityResponseEvent decoded")
				}
				e.iOCapabilityResponseEvent = cb(e.iOCapabilityResponseEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("IOCapabilityResponseEvent has no callback")
			}
		}
	case 0x3300:
		e.cbMutex.RLock()
		cb := e.userConfirmationRequestEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.userConfirmationRequestEvent == nil {
				e.userConfirmationRequestEvent = &UserConfirmationRequestEvent{}
			}

			if e.userConfirmationRequestEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.userConfirmationRequestEvent).Debug("UserConfirmationRequestEvent decoded")
				}
				e.userConfirmationRequestEvent = cb(e.userConfirmationRequestEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("UserConfirmationRequestEvent has no callback")
			}
		}
	case 0x3400:
		e.cbMutex.RLock()
		cb := e.userPasskeyRequestEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.userPasskeyRequestEvent == nil {
				e.userPasskeyRequestEvent = &UserPasskeyRequestEvent{}
			}

			if e.userPasskeyRequestEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.userPasskeyRequestEvent).Debug("UserPasskeyRequestEvent decoded")
				}
				e.userPasskeyRequestEvent = cb(e.userPasskeyRequestEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("UserPasskeyRequestEvent has no callback")
			}
		}
	case 0x3500:
		e.cbMutex.RLock()
		cb := e.remoteOOBDataRequestEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.remoteOOBDataRequestEvent == nil {
				e.remoteOOBDataRequestEvent = &RemoteOOBDataRequestEvent{}
			}

			if e.remoteOOBDataRequestEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.remoteOOBDataRequestEvent).Debug("RemoteOOBDataRequestEvent decoded")
				}
				e.remoteOOBDataRequestEvent = cb(e.remoteOOBDataRequestEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("RemoteOOBDataRequestEvent has no callback")
			}
		}
	case 0x3600:
		e.cbMutex.RLock()
		cb := e.simplePairingCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.simplePairingCompleteEvent == nil {
				e.simplePairingCompleteEvent = &SimplePairingCompleteEvent{}
			}

			if e.simplePairingCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.simplePairingCompleteEvent).Debug("SimplePairingCompleteEvent decoded")
				}
				e.simplePairingCompleteEvent = cb(e.simplePairingCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("SimplePairingCompleteEvent has no callback")
			}
		}
	case 0x3800:
		e.cbMutex.RLock()
		cb := e.linkSupervisionTimeoutChangedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.linkSupervisionTimeoutChangedEvent == nil {
				e.linkSupervisionTimeoutChangedEvent = &LinkSupervisionTimeoutChangedEvent{}
			}

			if e.linkSupervisionTimeoutChangedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.linkSupervisionTimeoutChangedEvent).Debug("LinkSupervisionTimeoutChangedEvent decoded")
				}
				e.linkSupervisionTimeoutChangedEvent = cb(e.linkSupervisionTimeoutChangedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LinkSupervisionTimeoutChangedEvent has no callback")
			}
		}
	case 0x3900:
		e.cbMutex.RLock()
		cb := e.enhancedFlushCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.enhancedFlushCompleteEvent == nil {
				e.enhancedFlushCompleteEvent = &EnhancedFlushCompleteEvent{}
			}

			if e.enhancedFlushCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.enhancedFlushCompleteEvent).Debug("EnhancedFlushCompleteEvent decoded")
				}
				e.enhancedFlushCompleteEvent = cb(e.enhancedFlushCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("EnhancedFlushCompleteEvent has no callback")
			}
		}
	case 0x3B00:
		e.cbMutex.RLock()
		cb := e.userPasskeyNotificationEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.userPasskeyNotificationEvent == nil {
				e.userPasskeyNotificationEvent = &UserPasskeyNotificationEvent{}
			}

			if e.userPasskeyNotificationEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.userPasskeyNotificationEvent).Debug("UserPasskeyNotificationEvent decoded")
				}
				e.userPasskeyNotificationEvent = cb(e.userPasskeyNotificationEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("UserPasskeyNotificationEvent has no callback")
			}
		}
	case 0x3C00:
		e.cbMutex.RLock()
		cb := e.keypressNotificationEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.keypressNotificationEvent == nil {
				e.keypressNotificationEvent = &KeypressNotificationEvent{}
			}

			if e.keypressNotificationEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.keypressNotificationEvent).Debug("KeypressNotificationEvent decoded")
				}
				e.keypressNotificationEvent = cb(e.keypressNotificationEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("KeypressNotificationEvent has no callback")
			}
		}
	case 0x3D00:
		e.cbMutex.RLock()
		cb := e.remoteHostSupportedFeaturesNotificationEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.remoteHostSupportedFeaturesNotificationEvent == nil {
				e.remoteHostSupportedFeaturesNotificationEvent = &RemoteHostSupportedFeaturesNotificationEvent{}
			}

			if e.remoteHostSupportedFeaturesNotificationEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.remoteHostSupportedFeaturesNotificationEvent).Debug("RemoteHostSupportedFeaturesNotificationEvent decoded")
				}
				e.remoteHostSupportedFeaturesNotificationEvent = cb(e.remoteHostSupportedFeaturesNotificationEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("RemoteHostSupportedFeaturesNotificationEvent has no callback")
			}
		}
	case 0x4000:
		e.cbMutex.RLock()
		cb := e.physicalLinkCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.physicalLinkCompleteEvent == nil {
				e.physicalLinkCompleteEvent = &PhysicalLinkCompleteEvent{}
			}

			if e.physicalLinkCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.physicalLinkCompleteEvent).Debug("PhysicalLinkCompleteEvent decoded")
				}
				e.physicalLinkCompleteEvent = cb(e.physicalLinkCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("PhysicalLinkCompleteEvent has no callback")
			}
		}
	case 0x4100:
		e.cbMutex.RLock()
		cb := e.channelSelectedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.channelSelectedEvent == nil {
				e.channelSelectedEvent = &ChannelSelectedEvent{}
			}

			if e.channelSelectedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.channelSelectedEvent).Debug("ChannelSelectedEvent decoded")
				}
				e.channelSelectedEvent = cb(e.channelSelectedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ChannelSelectedEvent has no callback")
			}
		}
	case 0x4200:
		e.cbMutex.RLock()
		cb := e.disconnectionPhysicalLinkCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.disconnectionPhysicalLinkCompleteEvent == nil {
				e.disconnectionPhysicalLinkCompleteEvent = &DisconnectionPhysicalLinkCompleteEvent{}
			}

			if e.disconnectionPhysicalLinkCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.disconnectionPhysicalLinkCompleteEvent).Debug("DisconnectionPhysicalLinkCompleteEvent decoded")
				}
				e.disconnectionPhysicalLinkCompleteEvent = cb(e.disconnectionPhysicalLinkCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("DisconnectionPhysicalLinkCompleteEvent has no callback")
			}
		}
	case 0x4300:
		e.cbMutex.RLock()
		cb := e.physicalLinkLossEarlyWarningEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.physicalLinkLossEarlyWarningEvent == nil {
				e.physicalLinkLossEarlyWarningEvent = &PhysicalLinkLossEarlyWarningEvent{}
			}

			if e.physicalLinkLossEarlyWarningEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.physicalLinkLossEarlyWarningEvent).Debug("PhysicalLinkLossEarlyWarningEvent decoded")
				}
				e.physicalLinkLossEarlyWarningEvent = cb(e.physicalLinkLossEarlyWarningEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("PhysicalLinkLossEarlyWarningEvent has no callback")
			}
		}
	case 0x4400:
		e.cbMutex.RLock()
		cb := e.physicalLinkRecoveryEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.physicalLinkRecoveryEvent == nil {
				e.physicalLinkRecoveryEvent = &PhysicalLinkRecoveryEvent{}
			}

			if e.physicalLinkRecoveryEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.physicalLinkRecoveryEvent).Debug("PhysicalLinkRecoveryEvent decoded")
				}
				e.physicalLinkRecoveryEvent = cb(e.physicalLinkRecoveryEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("PhysicalLinkRecoveryEvent has no callback")
			}
		}
	case 0x4500:
		e.cbMutex.RLock()
		cb := e.logicalLinkCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.logicalLinkCompleteEvent == nil {
				e.logicalLinkCompleteEvent = &LogicalLinkCompleteEvent{}
			}

			if e.logicalLinkCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.logicalLinkCompleteEvent).Debug("LogicalLinkCompleteEvent decoded")
				}
				e.logicalLinkCompleteEvent = cb(e.logicalLinkCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LogicalLinkCompleteEvent has no callback")
			}
		}
	case 0x4600:
		e.cbMutex.RLock()
		cb := e.disconnectionLogicalLinkCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.disconnectionLogicalLinkCompleteEvent == nil {
				e.disconnectionLogicalLinkCompleteEvent = &DisconnectionLogicalLinkCompleteEvent{}
			}

			if e.disconnectionLogicalLinkCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.disconnectionLogicalLinkCompleteEvent).Debug("DisconnectionLogicalLinkCompleteEvent decoded")
				}
				e.disconnectionLogicalLinkCompleteEvent = cb(e.disconnectionLogicalLinkCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("DisconnectionLogicalLinkCompleteEvent has no callback")
			}
		}
	case 0x4700:
		e.cbMutex.RLock()
		cb := e.flowSpecModifyCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.flowSpecModifyCompleteEvent == nil {
				e.flowSpecModifyCompleteEvent = &FlowSpecModifyCompleteEvent{}
			}

			if e.flowSpecModifyCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.flowSpecModifyCompleteEvent).Debug("FlowSpecModifyCompleteEvent decoded")
				}
				e.flowSpecModifyCompleteEvent = cb(e.flowSpecModifyCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("FlowSpecModifyCompleteEvent has no callback")
			}
		}
	case 0x4800:
		e.cbMutex.RLock()
		cb := e.numberOfCompletedDataBlocksEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.numberOfCompletedDataBlocksEvent == nil {
				e.numberOfCompletedDataBlocksEvent = &NumberOfCompletedDataBlocksEvent{}
			}

			if e.numberOfCompletedDataBlocksEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.numberOfCompletedDataBlocksEvent).Debug("NumberOfCompletedDataBlocksEvent decoded")
				}
				e.numberOfCompletedDataBlocksEvent = cb(e.numberOfCompletedDataBlocksEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("NumberOfCompletedDataBlocksEvent has no callback")
			}
		}
	case 0x4C00:
		e.cbMutex.RLock()
		cb := e.shortRangeModeChangeCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.shortRangeModeChangeCompleteEvent == nil {
				e.shortRangeModeChangeCompleteEvent = &ShortRangeModeChangeCompleteEvent{}
			}

			if e.shortRangeModeChangeCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.shortRangeModeChangeCompleteEvent).Debug("ShortRangeModeChangeCompleteEvent decoded")
				}
				e.shortRangeModeChangeCompleteEvent = cb(e.shortRangeModeChangeCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ShortRangeModeChangeCompleteEvent has no callback")
			}
		}
	case 0x4D00:
		e.cbMutex.RLock()
		cb := e.aMPStatusChangeEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.aMPStatusChangeEvent == nil {
				e.aMPStatusChangeEvent = &AMPStatusChangeEvent{}
			}

			if e.aMPStatusChangeEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.aMPStatusChangeEvent).Debug("AMPStatusChangeEvent decoded")
				}
				e.aMPStatusChangeEvent = cb(e.aMPStatusChangeEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("AMPStatusChangeEvent has no callback")
			}
		}
	case 0x4900:
		e.cbMutex.RLock()
		cb := e.aMPStartTestEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.aMPStartTestEvent == nil {
				e.aMPStartTestEvent = &AMPStartTestEvent{}
			}

			if e.aMPStartTestEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.aMPStartTestEvent).Debug("AMPStartTestEvent decoded")
				}
				e.aMPStartTestEvent = cb(e.aMPStartTestEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("AMPStartTestEvent has no callback")
			}
		}
	case 0x4A00:
		e.cbMutex.RLock()
		cb := e.aMPTestEndEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.aMPTestEndEvent == nil {
				e.aMPTestEndEvent = &AMPTestEndEvent{}
			}

			if e.aMPTestEndEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.aMPTestEndEvent).Debug("AMPTestEndEvent decoded")
				}
				e.aMPTestEndEvent = cb(e.aMPTestEndEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("AMPTestEndEvent has no callback")
			}
		}
	case 0x4B00:
		e.cbMutex.RLock()
		cb := e.aMPReceiverReportEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.aMPReceiverReportEvent == nil {
				e.aMPReceiverReportEvent = &AMPReceiverReportEvent{}
			}

			if e.aMPReceiverReportEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.aMPReceiverReportEvent).Debug("AMPReceiverReportEvent decoded")
				}
				e.aMPReceiverReportEvent = cb(e.aMPReceiverReportEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("AMPReceiverReportEvent has no callback")
			}
		}
	case 0x3E01:
		e.cbMutex.RLock()
		cb := e.lEConnectionCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEConnectionCompleteEvent == nil {
				e.lEConnectionCompleteEvent = &LEConnectionCompleteEvent{}
			}

			if e.lEConnectionCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEConnectionCompleteEvent).Debug("LEConnectionCompleteEvent decoded")
				}
				e.lEConnectionCompleteEvent = cb(e.lEConnectionCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEConnectionCompleteEvent has no callback")
			}
		}
	case 0x3E02:
		e.cbMutex.RLock()
		cb := e.lEAdvertisingReportEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEAdvertisingReportEvent == nil {
				e.lEAdvertisingReportEvent = &LEAdvertisingReportEvent{}
			}

			if e.lEAdvertisingReportEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
					e.logger.WithField("0data", e.lEAdvertisingReportEvent).Trace("LEAdvertisingReportEvent decoded")
				}
				e.lEAdvertisingReportEvent = cb(e.lEAdvertisingReportEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEAdvertisingReportEvent has no callback")
			}
		}
	case 0x3E03:
		e.cbMutex.RLock()
		cb := e.lEConnectionUpdateCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEConnectionUpdateCompleteEvent == nil {
				e.lEConnectionUpdateCompleteEvent = &LEConnectionUpdateCompleteEvent{}
			}

			if e.lEConnectionUpdateCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEConnectionUpdateCompleteEvent).Debug("LEConnectionUpdateCompleteEvent decoded")
				}
				e.lEConnectionUpdateCompleteEvent = cb(e.lEConnectionUpdateCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEConnectionUpdateCompleteEvent has no callback")
			}
		}
	case 0x3E04:
		e.cbMutex.RLock()
		cb := e.lEReadRemoteFeaturesCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEReadRemoteFeaturesCompleteEvent == nil {
				e.lEReadRemoteFeaturesCompleteEvent = &LEReadRemoteFeaturesCompleteEvent{}
			}

			if e.lEReadRemoteFeaturesCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEReadRemoteFeaturesCompleteEvent).Debug("LEReadRemoteFeaturesCompleteEvent decoded")
				}
				e.lEReadRemoteFeaturesCompleteEvent = cb(e.lEReadRemoteFeaturesCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEReadRemoteFeaturesCompleteEvent has no callback")
			}
		}
	case 0x3E05:
		e.cbMutex.RLock()
		cb := e.lELongTermKeyRequestEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lELongTermKeyRequestEvent == nil {
				e.lELongTermKeyRequestEvent = &LELongTermKeyRequestEvent{}
			}

			if e.lELongTermKeyRequestEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lELongTermKeyRequestEvent).Debug("LELongTermKeyRequestEvent decoded")
				}
				e.lELongTermKeyRequestEvent = cb(e.lELongTermKeyRequestEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LELongTermKeyRequestEvent has no callback")
			}
		}
	case 0x3E06:
		e.cbMutex.RLock()
		cb := e.lERemoteConnectionParameterRequestEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lERemoteConnectionParameterRequestEvent == nil {
				e.lERemoteConnectionParameterRequestEvent = &LERemoteConnectionParameterRequestEvent{}
			}

			if e.lERemoteConnectionParameterRequestEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lERemoteConnectionParameterRequestEvent).Debug("LERemoteConnectionParameterRequestEvent decoded")
				}
				e.lERemoteConnectionParameterRequestEvent = cb(e.lERemoteConnectionParameterRequestEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LERemoteConnectionParameterRequestEvent has no callback")
			}
		}
	case 0x3E07:
		e.cbMutex.RLock()
		cb := e.lEDataLengthChangeEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEDataLengthChangeEvent == nil {
				e.lEDataLengthChangeEvent = &LEDataLengthChangeEvent{}
			}

			if e.lEDataLengthChangeEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEDataLengthChangeEvent).Debug("LEDataLengthChangeEvent decoded")
				}
				e.lEDataLengthChangeEvent = cb(e.lEDataLengthChangeEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEDataLengthChangeEvent has no callback")
			}
		}
	case 0x3E08:
		e.cbMutex.RLock()
		cb := e.lEReadLocalP256PublicKeyCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEReadLocalP256PublicKeyCompleteEvent == nil {
				e.lEReadLocalP256PublicKeyCompleteEvent = &LEReadLocalP256PublicKeyCompleteEvent{}
			}

			if e.lEReadLocalP256PublicKeyCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEReadLocalP256PublicKeyCompleteEvent).Debug("LEReadLocalP256PublicKeyCompleteEvent decoded")
				}
				e.lEReadLocalP256PublicKeyCompleteEvent = cb(e.lEReadLocalP256PublicKeyCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEReadLocalP256PublicKeyCompleteEvent has no callback")
			}
		}
	case 0x3E09:
		e.cbMutex.RLock()
		cb := e.lEGenerateDHKeyCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEGenerateDHKeyCompleteEvent == nil {
				e.lEGenerateDHKeyCompleteEvent = &LEGenerateDHKeyCompleteEvent{}
			}

			if e.lEGenerateDHKeyCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEGenerateDHKeyCompleteEvent).Debug("LEGenerateDHKeyCompleteEvent decoded")
				}
				e.lEGenerateDHKeyCompleteEvent = cb(e.lEGenerateDHKeyCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEGenerateDHKeyCompleteEvent has no callback")
			}
		}
	case 0x3E0A:
		e.cbMutex.RLock()
		cb := e.lEEnhancedConnectionCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEEnhancedConnectionCompleteEvent == nil {
				e.lEEnhancedConnectionCompleteEvent = &LEEnhancedConnectionCompleteEvent{}
			}

			if e.lEEnhancedConnectionCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEEnhancedConnectionCompleteEvent).Debug("LEEnhancedConnectionCompleteEvent decoded")
				}
				e.lEEnhancedConnectionCompleteEvent = cb(e.lEEnhancedConnectionCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEEnhancedConnectionCompleteEvent has no callback")
			}
		}
	case 0x3E0B:
		e.cbMutex.RLock()
		cb := e.lEDirectedAdvertisingReportEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEDirectedAdvertisingReportEvent == nil {
				e.lEDirectedAdvertisingReportEvent = &LEDirectedAdvertisingReportEvent{}
			}

			if e.lEDirectedAdvertisingReportEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEDirectedAdvertisingReportEvent).Debug("LEDirectedAdvertisingReportEvent decoded")
				}
				e.lEDirectedAdvertisingReportEvent = cb(e.lEDirectedAdvertisingReportEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEDirectedAdvertisingReportEvent has no callback")
			}
		}
	case 0x3E0C:
		e.cbMutex.RLock()
		cb := e.lEPHYUpdateCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEPHYUpdateCompleteEvent == nil {
				e.lEPHYUpdateCompleteEvent = &LEPHYUpdateCompleteEvent{}
			}

			if e.lEPHYUpdateCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEPHYUpdateCompleteEvent).Debug("LEPHYUpdateCompleteEvent decoded")
				}
				e.lEPHYUpdateCompleteEvent = cb(e.lEPHYUpdateCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEPHYUpdateCompleteEvent has no callback")
			}
		}
	case 0x3E0D:
		e.cbMutex.RLock()
		cb := e.lEExtendedAdvertisingReportEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEExtendedAdvertisingReportEvent == nil {
				e.lEExtendedAdvertisingReportEvent = &LEExtendedAdvertisingReportEvent{}
			}

			if e.lEExtendedAdvertisingReportEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEExtendedAdvertisingReportEvent).Debug("LEExtendedAdvertisingReportEvent decoded")
				}
				e.lEExtendedAdvertisingReportEvent = cb(e.lEExtendedAdvertisingReportEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEExtendedAdvertisingReportEvent has no callback")
			}
		}
	case 0x3E0E:
		e.cbMutex.RLock()
		cb := e.lEPeriodicAdvertisingSyncEstablishedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEPeriodicAdvertisingSyncEstablishedEvent == nil {
				e.lEPeriodicAdvertisingSyncEstablishedEvent = &LEPeriodicAdvertisingSyncEstablishedEvent{}
			}

			if e.lEPeriodicAdvertisingSyncEstablishedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEPeriodicAdvertisingSyncEstablishedEvent).Debug("LEPeriodicAdvertisingSyncEstablishedEvent decoded")
				}
				e.lEPeriodicAdvertisingSyncEstablishedEvent = cb(e.lEPeriodicAdvertisingSyncEstablishedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEPeriodicAdvertisingSyncEstablishedEvent has no callback")
			}
		}
	case 0x3E0F:
		e.cbMutex.RLock()
		cb := e.lEPeriodicAdvertisingReportEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEPeriodicAdvertisingReportEvent == nil {
				e.lEPeriodicAdvertisingReportEvent = &LEPeriodicAdvertisingReportEvent{}
			}

			if e.lEPeriodicAdvertisingReportEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEPeriodicAdvertisingReportEvent).Debug("LEPeriodicAdvertisingReportEvent decoded")
				}
				e.lEPeriodicAdvertisingReportEvent = cb(e.lEPeriodicAdvertisingReportEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEPeriodicAdvertisingReportEvent has no callback")
			}
		}
	case 0x3E10:
		e.cbMutex.RLock()
		cb := e.lEPeriodicAdvertisingSyncLostEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEPeriodicAdvertisingSyncLostEvent == nil {
				e.lEPeriodicAdvertisingSyncLostEvent = &LEPeriodicAdvertisingSyncLostEvent{}
			}

			if e.lEPeriodicAdvertisingSyncLostEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEPeriodicAdvertisingSyncLostEvent).Debug("LEPeriodicAdvertisingSyncLostEvent decoded")
				}
				e.lEPeriodicAdvertisingSyncLostEvent = cb(e.lEPeriodicAdvertisingSyncLostEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEPeriodicAdvertisingSyncLostEvent has no callback")
			}
		}
	case 0x3E11:
		e.cbMutex.RLock()
		cb := e.lEScanTimeoutEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEScanTimeoutEvent == nil {
				e.lEScanTimeoutEvent = &LEScanTimeoutEvent{}
			}

			if e.lEScanTimeoutEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEScanTimeoutEvent).Debug("LEScanTimeoutEvent decoded")
				}
				e.lEScanTimeoutEvent = cb(e.lEScanTimeoutEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEScanTimeoutEvent has no callback")
			}
		}
	case 0x3E12:
		e.cbMutex.RLock()
		cb := e.lEAdvertisingSetTerminatedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEAdvertisingSetTerminatedEvent == nil {
				e.lEAdvertisingSetTerminatedEvent = &LEAdvertisingSetTerminatedEvent{}
			}

			if e.lEAdvertisingSetTerminatedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEAdvertisingSetTerminatedEvent).Debug("LEAdvertisingSetTerminatedEvent decoded")
				}
				e.lEAdvertisingSetTerminatedEvent = cb(e.lEAdvertisingSetTerminatedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEAdvertisingSetTerminatedEvent has no callback")
			}
		}
	case 0x3E13:
		e.cbMutex.RLock()
		cb := e.lEScanRequestReceivedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEScanRequestReceivedEvent == nil {
				e.lEScanRequestReceivedEvent = &LEScanRequestReceivedEvent{}
			}

			if e.lEScanRequestReceivedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEScanRequestReceivedEvent).Debug("LEScanRequestReceivedEvent decoded")
				}
				e.lEScanRequestReceivedEvent = cb(e.lEScanRequestReceivedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEScanRequestReceivedEvent has no callback")
			}
		}
	case 0x3E14:
		e.cbMutex.RLock()
		cb := e.lEChannelSelectionAlgorithmEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEChannelSelectionAlgorithmEvent == nil {
				e.lEChannelSelectionAlgorithmEvent = &LEChannelSelectionAlgorithmEvent{}
			}

			if e.lEChannelSelectionAlgorithmEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEChannelSelectionAlgorithmEvent).Debug("LEChannelSelectionAlgorithmEvent decoded")
				}
				e.lEChannelSelectionAlgorithmEvent = cb(e.lEChannelSelectionAlgorithmEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEChannelSelectionAlgorithmEvent has no callback")
			}
		}
	case 0x3E15:
		e.cbMutex.RLock()
		cb := e.lEConnectionlessIQReportEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEConnectionlessIQReportEvent == nil {
				e.lEConnectionlessIQReportEvent = &LEConnectionlessIQReportEvent{}
			}

			if e.lEConnectionlessIQReportEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEConnectionlessIQReportEvent).Debug("LEConnectionlessIQReportEvent decoded")
				}
				e.lEConnectionlessIQReportEvent = cb(e.lEConnectionlessIQReportEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEConnectionlessIQReportEvent has no callback")
			}
		}
	case 0x3E16:
		e.cbMutex.RLock()
		cb := e.lEConnectionIQReportEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEConnectionIQReportEvent == nil {
				e.lEConnectionIQReportEvent = &LEConnectionIQReportEvent{}
			}

			if e.lEConnectionIQReportEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEConnectionIQReportEvent).Debug("LEConnectionIQReportEvent decoded")
				}
				e.lEConnectionIQReportEvent = cb(e.lEConnectionIQReportEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEConnectionIQReportEvent has no callback")
			}
		}
	case 0x3E17:
		e.cbMutex.RLock()
		cb := e.lECTERequestFailedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lECTERequestFailedEvent == nil {
				e.lECTERequestFailedEvent = &LECTERequestFailedEvent{}
			}

			if e.lECTERequestFailedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lECTERequestFailedEvent).Debug("LECTERequestFailedEvent decoded")
				}
				e.lECTERequestFailedEvent = cb(e.lECTERequestFailedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LECTERequestFailedEvent has no callback")
			}
		}
	case 0x3E18:
		e.cbMutex.RLock()
		cb := e.lEPeriodicAdvertisingSyncTransferReceivedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEPeriodicAdvertisingSyncTransferReceivedEvent == nil {
				e.lEPeriodicAdvertisingSyncTransferReceivedEvent = &LEPeriodicAdvertisingSyncTransferReceivedEvent{}
			}

			if e.lEPeriodicAdvertisingSyncTransferReceivedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEPeriodicAdvertisingSyncTransferReceivedEvent).Debug("LEPeriodicAdvertisingSyncTransferReceivedEvent decoded")
				}
				e.lEPeriodicAdvertisingSyncTransferReceivedEvent = cb(e.lEPeriodicAdvertisingSyncTransferReceivedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEPeriodicAdvertisingSyncTransferReceivedEvent has no callback")
			}
		}
	case 0x3E19:
		e.cbMutex.RLock()
		cb := e.lECISEstablishedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lECISEstablishedEvent == nil {
				e.lECISEstablishedEvent = &LECISEstablishedEvent{}
			}

			if e.lECISEstablishedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lECISEstablishedEvent).Debug("LECISEstablishedEvent decoded")
				}
				e.lECISEstablishedEvent = cb(e.lECISEstablishedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LECISEstablishedEvent has no callback")
			}
		}
	case 0x3E1A:
		e.cbMutex.RLock()
		cb := e.lECISRequestEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lECISRequestEvent == nil {
				e.lECISRequestEvent = &LECISRequestEvent{}
			}

			if e.lECISRequestEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lECISRequestEvent).Debug("LECISRequestEvent decoded")
				}
				e.lECISRequestEvent = cb(e.lECISRequestEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LECISRequestEvent has no callback")
			}
		}
	case 0x3E1B:
		e.cbMutex.RLock()
		cb := e.lECreateBIGCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lECreateBIGCompleteEvent == nil {
				e.lECreateBIGCompleteEvent = &LECreateBIGCompleteEvent{}
			}

			if e.lECreateBIGCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lECreateBIGCompleteEvent).Debug("LECreateBIGCompleteEvent decoded")
				}
				e.lECreateBIGCompleteEvent = cb(e.lECreateBIGCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LECreateBIGCompleteEvent has no callback")
			}
		}
	case 0x3E1C:
		e.cbMutex.RLock()
		cb := e.lETerminateBIGCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lETerminateBIGCompleteEvent == nil {
				e.lETerminateBIGCompleteEvent = &LETerminateBIGCompleteEvent{}
			}

			if e.lETerminateBIGCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lETerminateBIGCompleteEvent).Debug("LETerminateBIGCompleteEvent decoded")
				}
				e.lETerminateBIGCompleteEvent = cb(e.lETerminateBIGCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LETerminateBIGCompleteEvent has no callback")
			}
		}
	case 0x3E1D:
		e.cbMutex.RLock()
		cb := e.lEBIGSyncEstablishedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEBIGSyncEstablishedEvent == nil {
				e.lEBIGSyncEstablishedEvent = &LEBIGSyncEstablishedEvent{}
			}

			if e.lEBIGSyncEstablishedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEBIGSyncEstablishedEvent).Debug("LEBIGSyncEstablishedEvent decoded")
				}
				e.lEBIGSyncEstablishedEvent = cb(e.lEBIGSyncEstablishedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEBIGSyncEstablishedEvent has no callback")
			}
		}
	case 0x3E1E:
		e.cbMutex.RLock()
		cb := e.lEBIGSyncLostEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEBIGSyncLostEvent == nil {
				e.lEBIGSyncLostEvent = &LEBIGSyncLostEvent{}
			}

			if e.lEBIGSyncLostEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEBIGSyncLostEvent).Debug("LEBIGSyncLostEvent decoded")
				}
				e.lEBIGSyncLostEvent = cb(e.lEBIGSyncLostEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEBIGSyncLostEvent has no callback")
			}
		}
	case 0x3E1F:
		e.cbMutex.RLock()
		cb := e.lERequestPeerSCACompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lERequestPeerSCACompleteEvent == nil {
				e.lERequestPeerSCACompleteEvent = &LERequestPeerSCACompleteEvent{}
			}

			if e.lERequestPeerSCACompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lERequestPeerSCACompleteEvent).Debug("LERequestPeerSCACompleteEvent decoded")
				}
				e.lERequestPeerSCACompleteEvent = cb(e.lERequestPeerSCACompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LERequestPeerSCACompleteEvent has no callback")
			}
		}
	case 0x3E20:
		e.cbMutex.RLock()
		cb := e.lEPathLossThresholdEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEPathLossThresholdEvent == nil {
				e.lEPathLossThresholdEvent = &LEPathLossThresholdEvent{}
			}

			if e.lEPathLossThresholdEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEPathLossThresholdEvent).Debug("LEPathLossThresholdEvent decoded")
				}
				e.lEPathLossThresholdEvent = cb(e.lEPathLossThresholdEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEPathLossThresholdEvent has no callback")
			}
		}
	case 0x3E21:
		e.cbMutex.RLock()
		cb := e.lETransmitPowerReportingEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lETransmitPowerReportingEvent == nil {
				e.lETransmitPowerReportingEvent = &LETransmitPowerReportingEvent{}
			}

			if e.lETransmitPowerReportingEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lETransmitPowerReportingEvent).Debug("LETransmitPowerReportingEvent decoded")
				}
				e.lETransmitPowerReportingEvent = cb(e.lETransmitPowerReportingEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LETransmitPowerReportingEvent has no callback")
			}
		}
	case 0x3E22:
		e.cbMutex.RLock()
		cb := e.lEBIGInfoAdvertisingReportEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.lEBIGInfoAdvertisingReportEvent == nil {
				e.lEBIGInfoAdvertisingReportEvent = &LEBIGInfoAdvertisingReportEvent{}
			}

			if e.lEBIGInfoAdvertisingReportEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.lEBIGInfoAdvertisingReportEvent).Debug("LEBIGInfoAdvertisingReportEvent decoded")
				}
				e.lEBIGInfoAdvertisingReportEvent = cb(e.lEBIGInfoAdvertisingReportEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("LEBIGInfoAdvertisingReportEvent has no callback")
			}
		}
	case 0x4E00:
		e.cbMutex.RLock()
		cb := e.triggeredClockCaptureEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.triggeredClockCaptureEvent == nil {
				e.triggeredClockCaptureEvent = &TriggeredClockCaptureEvent{}
			}

			if e.triggeredClockCaptureEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.triggeredClockCaptureEvent).Debug("TriggeredClockCaptureEvent decoded")
				}
				e.triggeredClockCaptureEvent = cb(e.triggeredClockCaptureEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("TriggeredClockCaptureEvent has no callback")
			}
		}
	case 0x4F00:
		e.cbMutex.RLock()
		cb := e.synchronizationTrainCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.synchronizationTrainCompleteEvent == nil {
				e.synchronizationTrainCompleteEvent = &SynchronizationTrainCompleteEvent{}
			}

			if e.synchronizationTrainCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.synchronizationTrainCompleteEvent).Debug("SynchronizationTrainCompleteEvent decoded")
				}
				e.synchronizationTrainCompleteEvent = cb(e.synchronizationTrainCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("SynchronizationTrainCompleteEvent has no callback")
			}
		}
	case 0x5000:
		e.cbMutex.RLock()
		cb := e.synchronizationTrainReceivedEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.synchronizationTrainReceivedEvent == nil {
				e.synchronizationTrainReceivedEvent = &SynchronizationTrainReceivedEvent{}
			}

			if e.synchronizationTrainReceivedEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.synchronizationTrainReceivedEvent).Debug("SynchronizationTrainReceivedEvent decoded")
				}
				e.synchronizationTrainReceivedEvent = cb(e.synchronizationTrainReceivedEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("SynchronizationTrainReceivedEvent has no callback")
			}
		}
	case 0x5100:
		e.cbMutex.RLock()
		cb := e.connectionlessSlaveBroadcastReceiveEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.connectionlessSlaveBroadcastReceiveEvent == nil {
				e.connectionlessSlaveBroadcastReceiveEvent = &ConnectionlessSlaveBroadcastReceiveEvent{}
			}

			if e.connectionlessSlaveBroadcastReceiveEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.connectionlessSlaveBroadcastReceiveEvent).Debug("ConnectionlessSlaveBroadcastReceiveEvent decoded")
				}
				e.connectionlessSlaveBroadcastReceiveEvent = cb(e.connectionlessSlaveBroadcastReceiveEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ConnectionlessSlaveBroadcastReceiveEvent has no callback")
			}
		}
	case 0x5200:
		e.cbMutex.RLock()
		cb := e.connectionlessSlaveBroadcastTimeoutEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.connectionlessSlaveBroadcastTimeoutEvent == nil {
				e.connectionlessSlaveBroadcastTimeoutEvent = &ConnectionlessSlaveBroadcastTimeoutEvent{}
			}

			if e.connectionlessSlaveBroadcastTimeoutEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.connectionlessSlaveBroadcastTimeoutEvent).Debug("ConnectionlessSlaveBroadcastTimeoutEvent decoded")
				}
				e.connectionlessSlaveBroadcastTimeoutEvent = cb(e.connectionlessSlaveBroadcastTimeoutEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ConnectionlessSlaveBroadcastTimeoutEvent has no callback")
			}
		}
	case 0x5300:
		e.cbMutex.RLock()
		cb := e.truncatedPageCompleteEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.truncatedPageCompleteEvent == nil {
				e.truncatedPageCompleteEvent = &TruncatedPageCompleteEvent{}
			}

			if e.truncatedPageCompleteEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.truncatedPageCompleteEvent).Debug("TruncatedPageCompleteEvent decoded")
				}
				e.truncatedPageCompleteEvent = cb(e.truncatedPageCompleteEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("TruncatedPageCompleteEvent has no callback")
			}
		}
	case 0x5400:
		e.cbMutex.RLock()
		cb := e.slavePageResponseTimeoutEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.slavePageResponseTimeoutEvent == nil {
				e.slavePageResponseTimeoutEvent = &SlavePageResponseTimeoutEvent{}
			}

			if e.slavePageResponseTimeoutEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.slavePageResponseTimeoutEvent).Debug("SlavePageResponseTimeoutEvent decoded")
				}
				e.slavePageResponseTimeoutEvent = cb(e.slavePageResponseTimeoutEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("SlavePageResponseTimeoutEvent has no callback")
			}
		}
	case 0x5500:
		e.cbMutex.RLock()
		cb := e.connectionlessSlaveBroadcastChannelMapChangeEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.connectionlessSlaveBroadcastChannelMapChangeEvent == nil {
				e.connectionlessSlaveBroadcastChannelMapChangeEvent = &ConnectionlessSlaveBroadcastChannelMapChangeEvent{}
			}

			if e.connectionlessSlaveBroadcastChannelMapChangeEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.connectionlessSlaveBroadcastChannelMapChangeEvent).Debug("ConnectionlessSlaveBroadcastChannelMapChangeEvent decoded")
				}
				e.connectionlessSlaveBroadcastChannelMapChangeEvent = cb(e.connectionlessSlaveBroadcastChannelMapChangeEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("ConnectionlessSlaveBroadcastChannelMapChangeEvent has no callback")
			}
		}
	case 0x5600:
		e.cbMutex.RLock()
		cb := e.inquiryResponseNotificationEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.inquiryResponseNotificationEvent == nil {
				e.inquiryResponseNotificationEvent = &InquiryResponseNotificationEvent{}
			}

			if e.inquiryResponseNotificationEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.inquiryResponseNotificationEvent).Debug("InquiryResponseNotificationEvent decoded")
				}
				e.inquiryResponseNotificationEvent = cb(e.inquiryResponseNotificationEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("InquiryResponseNotificationEvent has no callback")
			}
		}
	case 0x5700:
		e.cbMutex.RLock()
		cb := e.authenticatedPayloadTimeoutExpiredEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.authenticatedPayloadTimeoutExpiredEvent == nil {
				e.authenticatedPayloadTimeoutExpiredEvent = &AuthenticatedPayloadTimeoutExpiredEvent{}
			}

			if e.authenticatedPayloadTimeoutExpiredEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.authenticatedPayloadTimeoutExpiredEvent).Debug("AuthenticatedPayloadTimeoutExpiredEvent decoded")
				}
				e.authenticatedPayloadTimeoutExpiredEvent = cb(e.authenticatedPayloadTimeoutExpiredEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("AuthenticatedPayloadTimeoutExpiredEvent has no callback")
			}
		}
	case 0x5800:
		e.cbMutex.RLock()
		cb := e.sAMStatusChangeEventCallback
		e.cbMutex.RUnlock()

		if cb != nil {
			if e.sAMStatusChangeEvent == nil {
				e.sAMStatusChangeEvent = &SAMStatusChangeEvent{}
			}

			if e.sAMStatusChangeEvent.decode(params) {
				if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
					e.logger.WithField("0data", e.sAMStatusChangeEvent).Debug("SAMStatusChangeEvent decoded")
				}
				e.sAMStatusChangeEvent = cb(e.sAMStatusChangeEvent)
			}
		}else{
			if e.logger != nil && e.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
				e.logger.Debug("SAMStatusChangeEvent has no callback")
			}
		}
	}
}
