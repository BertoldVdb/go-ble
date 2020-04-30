package hcicommands

import (
	"encoding/binary"
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
)

// LinkControlInquiryInput represents the input of the command specified in Section 7.1.1
type LinkControlInquiryInput struct {
	LAP uint32
	InquiryLength uint8
	NumResponses uint8
}

func (i LinkControlInquiryInput) encode(data []byte) []byte {
	w := writer{data: data};
	encodeUint24(w.Put(3), i.LAP)
	w.PutOne(i.InquiryLength)
	w.PutOne(i.NumResponses)
	return w.Data()
}

// LinkControlInquirySync executes the command specified in Section 7.1.1 synchronously
func (c *Commands) LinkControlInquirySync (params LinkControlInquiryInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0001}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlInquiryCancelSync executes the command specified in Section 7.1.2 synchronously
func (c *Commands) LinkControlInquiryCancelSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0002}, nil)
	if err != nil {
		return err
	}

	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlPeriodicInquiryModeInput represents the input of the command specified in Section 7.1.3
type LinkControlPeriodicInquiryModeInput struct {
	MaxPeriodLength uint16
	MinPeriodLength uint16
	LAP uint32
	InquiryLength uint8
	NumResponses uint8
}

func (i LinkControlPeriodicInquiryModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxPeriodLength)
	binary.LittleEndian.PutUint16(w.Put(2), i.MinPeriodLength)
	encodeUint24(w.Put(3), i.LAP)
	w.PutOne(i.InquiryLength)
	w.PutOne(i.NumResponses)
	return w.Data()
}

// LinkControlPeriodicInquiryModeSync executes the command specified in Section 7.1.3 synchronously
func (c *Commands) LinkControlPeriodicInquiryModeSync (params LinkControlPeriodicInquiryModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0003}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlExitPeriodicInquiryModeSync executes the command specified in Section 7.1.4 synchronously
func (c *Commands) LinkControlExitPeriodicInquiryModeSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0004}, nil)
	if err != nil {
		return err
	}

	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlCreateConnectionInput represents the input of the command specified in Section 7.1.5
type LinkControlCreateConnectionInput struct {
	BDADDR [6]byte
	PacketType uint16
	PageScanRepetitionMode uint8
	Reserved uint8
	ClockOffset uint16
	AllowRoleSwitch uint8
}

func (i LinkControlCreateConnectionInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	binary.LittleEndian.PutUint16(w.Put(2), i.PacketType)
	w.PutOne(i.PageScanRepetitionMode)
	w.PutOne(i.Reserved)
	binary.LittleEndian.PutUint16(w.Put(2), i.ClockOffset)
	w.PutOne(i.AllowRoleSwitch)
	return w.Data()
}

// LinkControlCreateConnectionSync executes the command specified in Section 7.1.5 synchronously
func (c *Commands) LinkControlCreateConnectionSync (params LinkControlCreateConnectionInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0005}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlDisconnectInput represents the input of the command specified in Section 7.1.6
type LinkControlDisconnectInput struct {
	ConnectionHandle uint16
	Reason uint8
}

func (i LinkControlDisconnectInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.Reason)
	return w.Data()
}

// LinkControlDisconnectSync executes the command specified in Section 7.1.6 synchronously
func (c *Commands) LinkControlDisconnectSync (params LinkControlDisconnectInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0006}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlCreateConnectionCancelInput represents the input of the command specified in Section 7.1.7
type LinkControlCreateConnectionCancelInput struct {
	BDADDR [6]byte
}

func (i LinkControlCreateConnectionCancelInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	return w.Data()
}

// LinkControlCreateConnectionCancelOutput represents the output of the command specified in Section 7.1.7
type LinkControlCreateConnectionCancelOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlCreateConnectionCancelOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlCreateConnectionCancelSync executes the command specified in Section 7.1.7 synchronously
func (c *Commands) LinkControlCreateConnectionCancelSync (params LinkControlCreateConnectionCancelInput, result *LinkControlCreateConnectionCancelOutput) (*LinkControlCreateConnectionCancelOutput, error) {
	if result == nil {
		result = &LinkControlCreateConnectionCancelOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0008}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlAcceptConnectionRequestInput represents the input of the command specified in Section 7.1.8
type LinkControlAcceptConnectionRequestInput struct {
	BDADDR [6]byte
	Role uint8
}

func (i LinkControlAcceptConnectionRequestInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.Role)
	return w.Data()
}

// LinkControlAcceptConnectionRequestSync executes the command specified in Section 7.1.8 synchronously
func (c *Commands) LinkControlAcceptConnectionRequestSync (params LinkControlAcceptConnectionRequestInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0009}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlRejectConnectionRequestInput represents the input of the command specified in Section 7.1.9
type LinkControlRejectConnectionRequestInput struct {
	BDADDR [6]byte
	Reason uint8
}

func (i LinkControlRejectConnectionRequestInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.Reason)
	return w.Data()
}

// LinkControlRejectConnectionRequestSync executes the command specified in Section 7.1.9 synchronously
func (c *Commands) LinkControlRejectConnectionRequestSync (params LinkControlRejectConnectionRequestInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x000A}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlLinkKeyRequestReplyInput represents the input of the command specified in Section 7.1.10
type LinkControlLinkKeyRequestReplyInput struct {
	BDADDR [6]byte
	LinkKey [16]byte
}

func (i LinkControlLinkKeyRequestReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	copy(w.Put(16), i.LinkKey[:])
	return w.Data()
}

// LinkControlLinkKeyRequestReplyOutput represents the output of the command specified in Section 7.1.10
type LinkControlLinkKeyRequestReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlLinkKeyRequestReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlLinkKeyRequestReplySync executes the command specified in Section 7.1.10 synchronously
func (c *Commands) LinkControlLinkKeyRequestReplySync (params LinkControlLinkKeyRequestReplyInput, result *LinkControlLinkKeyRequestReplyOutput) (*LinkControlLinkKeyRequestReplyOutput, error) {
	if result == nil {
		result = &LinkControlLinkKeyRequestReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x000B}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlLinkKeyRequestNegativeReplyInput represents the input of the command specified in Section 7.1.11
type LinkControlLinkKeyRequestNegativeReplyInput struct {
	BDADDR [6]byte
}

func (i LinkControlLinkKeyRequestNegativeReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	return w.Data()
}

// LinkControlLinkKeyRequestNegativeReplyOutput represents the output of the command specified in Section 7.1.11
type LinkControlLinkKeyRequestNegativeReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlLinkKeyRequestNegativeReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlLinkKeyRequestNegativeReplySync executes the command specified in Section 7.1.11 synchronously
func (c *Commands) LinkControlLinkKeyRequestNegativeReplySync (params LinkControlLinkKeyRequestNegativeReplyInput, result *LinkControlLinkKeyRequestNegativeReplyOutput) (*LinkControlLinkKeyRequestNegativeReplyOutput, error) {
	if result == nil {
		result = &LinkControlLinkKeyRequestNegativeReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x000C}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlPINCodeRequestReplyInput represents the input of the command specified in Section 7.1.12
type LinkControlPINCodeRequestReplyInput struct {
	BDADDR [6]byte
	PINCodeLength uint8
	PINCode [16]byte
}

func (i LinkControlPINCodeRequestReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.PINCodeLength)
	copy(w.Put(16), i.PINCode[:])
	return w.Data()
}

// LinkControlPINCodeRequestReplyOutput represents the output of the command specified in Section 7.1.12
type LinkControlPINCodeRequestReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlPINCodeRequestReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlPINCodeRequestReplySync executes the command specified in Section 7.1.12 synchronously
func (c *Commands) LinkControlPINCodeRequestReplySync (params LinkControlPINCodeRequestReplyInput, result *LinkControlPINCodeRequestReplyOutput) (*LinkControlPINCodeRequestReplyOutput, error) {
	if result == nil {
		result = &LinkControlPINCodeRequestReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x000D}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlPINCodeRequestNegativeReplyInput represents the input of the command specified in Section 7.1.13
type LinkControlPINCodeRequestNegativeReplyInput struct {
	BDADDR [6]byte
}

func (i LinkControlPINCodeRequestNegativeReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	return w.Data()
}

// LinkControlPINCodeRequestNegativeReplyOutput represents the output of the command specified in Section 7.1.13
type LinkControlPINCodeRequestNegativeReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlPINCodeRequestNegativeReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlPINCodeRequestNegativeReplySync executes the command specified in Section 7.1.13 synchronously
func (c *Commands) LinkControlPINCodeRequestNegativeReplySync (params LinkControlPINCodeRequestNegativeReplyInput, result *LinkControlPINCodeRequestNegativeReplyOutput) (*LinkControlPINCodeRequestNegativeReplyOutput, error) {
	if result == nil {
		result = &LinkControlPINCodeRequestNegativeReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x000E}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlChangeConnectionPacketTypeInput represents the input of the command specified in Section 7.1.14
type LinkControlChangeConnectionPacketTypeInput struct {
	ConnectionHandle uint16
	PacketType uint16
}

func (i LinkControlChangeConnectionPacketTypeInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.PacketType)
	return w.Data()
}

// LinkControlChangeConnectionPacketTypeSync executes the command specified in Section 7.1.14 synchronously
func (c *Commands) LinkControlChangeConnectionPacketTypeSync (params LinkControlChangeConnectionPacketTypeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x000F}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlAuthenticationRequestedInput represents the input of the command specified in Section 7.1.15
type LinkControlAuthenticationRequestedInput struct {
	ConnectionHandle uint16
}

func (i LinkControlAuthenticationRequestedInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LinkControlAuthenticationRequestedSync executes the command specified in Section 7.1.15 synchronously
func (c *Commands) LinkControlAuthenticationRequestedSync (params LinkControlAuthenticationRequestedInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0011}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlSetConnectionEncryptionInput represents the input of the command specified in Section 7.1.16
type LinkControlSetConnectionEncryptionInput struct {
	ConnectionHandle uint16
	EncryptionEnable uint8
}

func (i LinkControlSetConnectionEncryptionInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.EncryptionEnable)
	return w.Data()
}

// LinkControlSetConnectionEncryptionSync executes the command specified in Section 7.1.16 synchronously
func (c *Commands) LinkControlSetConnectionEncryptionSync (params LinkControlSetConnectionEncryptionInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0013}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlChangeConnectionLinkKeyInput represents the input of the command specified in Section 7.1.17
type LinkControlChangeConnectionLinkKeyInput struct {
	ConnectionHandle uint16
}

func (i LinkControlChangeConnectionLinkKeyInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LinkControlChangeConnectionLinkKeySync executes the command specified in Section 7.1.17 synchronously
func (c *Commands) LinkControlChangeConnectionLinkKeySync (params LinkControlChangeConnectionLinkKeyInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0015}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlMasterLinkKeyInput represents the input of the command specified in Section 7.1.18
type LinkControlMasterLinkKeyInput struct {
	KeyFlag uint8
}

func (i LinkControlMasterLinkKeyInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.KeyFlag)
	return w.Data()
}

// LinkControlMasterLinkKeySync executes the command specified in Section 7.1.18 synchronously
func (c *Commands) LinkControlMasterLinkKeySync (params LinkControlMasterLinkKeyInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0017}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlRemoteNameRequestInput represents the input of the command specified in Section 7.1.19
type LinkControlRemoteNameRequestInput struct {
	BDADDR [6]byte
	PageScanRepetitionMode uint8
	Reserved uint8
	ClockOffset uint16
}

func (i LinkControlRemoteNameRequestInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.PageScanRepetitionMode)
	w.PutOne(i.Reserved)
	binary.LittleEndian.PutUint16(w.Put(2), i.ClockOffset)
	return w.Data()
}

// LinkControlRemoteNameRequestSync executes the command specified in Section 7.1.19 synchronously
func (c *Commands) LinkControlRemoteNameRequestSync (params LinkControlRemoteNameRequestInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0019}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlRemoteNameRequestCancelInput represents the input of the command specified in Section 7.1.20
type LinkControlRemoteNameRequestCancelInput struct {
	BDADDR [6]byte
}

func (i LinkControlRemoteNameRequestCancelInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	return w.Data()
}

// LinkControlRemoteNameRequestCancelOutput represents the output of the command specified in Section 7.1.20
type LinkControlRemoteNameRequestCancelOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlRemoteNameRequestCancelOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlRemoteNameRequestCancelSync executes the command specified in Section 7.1.20 synchronously
func (c *Commands) LinkControlRemoteNameRequestCancelSync (params LinkControlRemoteNameRequestCancelInput, result *LinkControlRemoteNameRequestCancelOutput) (*LinkControlRemoteNameRequestCancelOutput, error) {
	if result == nil {
		result = &LinkControlRemoteNameRequestCancelOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x001A}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlReadRemoteSupportedFeaturesInput represents the input of the command specified in Section 7.1.21
type LinkControlReadRemoteSupportedFeaturesInput struct {
	ConnectionHandle uint16
}

func (i LinkControlReadRemoteSupportedFeaturesInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LinkControlReadRemoteSupportedFeaturesSync executes the command specified in Section 7.1.21 synchronously
func (c *Commands) LinkControlReadRemoteSupportedFeaturesSync (params LinkControlReadRemoteSupportedFeaturesInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x001B}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlReadRemoteExtendedFeaturesInput represents the input of the command specified in Section 7.1.22
type LinkControlReadRemoteExtendedFeaturesInput struct {
	ConnectionHandle uint16
	PageNumber uint8
}

func (i LinkControlReadRemoteExtendedFeaturesInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.PageNumber)
	return w.Data()
}

// LinkControlReadRemoteExtendedFeaturesSync executes the command specified in Section 7.1.22 synchronously
func (c *Commands) LinkControlReadRemoteExtendedFeaturesSync (params LinkControlReadRemoteExtendedFeaturesInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x001C}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlReadRemoteVersionInformationInput represents the input of the command specified in Section 7.1.23
type LinkControlReadRemoteVersionInformationInput struct {
	ConnectionHandle uint16
}

func (i LinkControlReadRemoteVersionInformationInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LinkControlReadRemoteVersionInformationSync executes the command specified in Section 7.1.23 synchronously
func (c *Commands) LinkControlReadRemoteVersionInformationSync (params LinkControlReadRemoteVersionInformationInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x001D}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlReadClockOffsetInput represents the input of the command specified in Section 7.1.24
type LinkControlReadClockOffsetInput struct {
	ConnectionHandle uint16
}

func (i LinkControlReadClockOffsetInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LinkControlReadClockOffsetSync executes the command specified in Section 7.1.24 synchronously
func (c *Commands) LinkControlReadClockOffsetSync (params LinkControlReadClockOffsetInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x001F}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlReadLMPHandleInput represents the input of the command specified in Section 7.1.25
type LinkControlReadLMPHandleInput struct {
	ConnectionHandle uint16
}

func (i LinkControlReadLMPHandleInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LinkControlReadLMPHandleOutput represents the output of the command specified in Section 7.1.25
type LinkControlReadLMPHandleOutput struct {
	Status uint8
	ConnectionHandle uint16
	LMPHandle uint8
	Reserved uint32
}

func (o *LinkControlReadLMPHandleOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LMPHandle = r.GetOne()
	o.Reserved = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// LinkControlReadLMPHandleSync executes the command specified in Section 7.1.25 synchronously
func (c *Commands) LinkControlReadLMPHandleSync (params LinkControlReadLMPHandleInput, result *LinkControlReadLMPHandleOutput) (*LinkControlReadLMPHandleOutput, error) {
	if result == nil {
		result = &LinkControlReadLMPHandleOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0020}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlSetupSynchronousConnectionInput represents the input of the command specified in Section 7.1.26
type LinkControlSetupSynchronousConnectionInput struct {
	ConnectionHandle uint16
	TransmitBandwidth uint32
	ReceiveBandwidth uint32
	MaxLatency uint16
	VoiceSetting uint16
	RetransmissionEffort uint8
	PacketType uint16
}

func (i LinkControlSetupSynchronousConnectionInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint32(w.Put(4), i.TransmitBandwidth)
	binary.LittleEndian.PutUint32(w.Put(4), i.ReceiveBandwidth)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxLatency)
	binary.LittleEndian.PutUint16(w.Put(2), i.VoiceSetting)
	w.PutOne(i.RetransmissionEffort)
	binary.LittleEndian.PutUint16(w.Put(2), i.PacketType)
	return w.Data()
}

// LinkControlSetupSynchronousConnectionSync executes the command specified in Section 7.1.26 synchronously
func (c *Commands) LinkControlSetupSynchronousConnectionSync (params LinkControlSetupSynchronousConnectionInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0028}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlAcceptSynchronousConnectionRequestInput represents the input of the command specified in Section 7.1.27
type LinkControlAcceptSynchronousConnectionRequestInput struct {
	BDADDR [6]byte
	TransmitBandwidth uint32
	ReceiveBandwidth uint32
	MaxLatency uint16
	VoiceSettings uint16
	RetransmissionEffort uint8
	PacketType uint16
}

func (i LinkControlAcceptSynchronousConnectionRequestInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	binary.LittleEndian.PutUint32(w.Put(4), i.TransmitBandwidth)
	binary.LittleEndian.PutUint32(w.Put(4), i.ReceiveBandwidth)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxLatency)
	binary.LittleEndian.PutUint16(w.Put(2), i.VoiceSettings)
	w.PutOne(i.RetransmissionEffort)
	binary.LittleEndian.PutUint16(w.Put(2), i.PacketType)
	return w.Data()
}

// LinkControlAcceptSynchronousConnectionRequestSync executes the command specified in Section 7.1.27 synchronously
func (c *Commands) LinkControlAcceptSynchronousConnectionRequestSync (params LinkControlAcceptSynchronousConnectionRequestInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0029}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlRejectSynchronousConnectionRequestInput represents the input of the command specified in Section 7.1.28
type LinkControlRejectSynchronousConnectionRequestInput struct {
	BDADDR [6]byte
	Reason uint8
}

func (i LinkControlRejectSynchronousConnectionRequestInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.Reason)
	return w.Data()
}

// LinkControlRejectSynchronousConnectionRequestSync executes the command specified in Section 7.1.28 synchronously
func (c *Commands) LinkControlRejectSynchronousConnectionRequestSync (params LinkControlRejectSynchronousConnectionRequestInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x002A}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlIOCapabilityRequestReplyInput represents the input of the command specified in Section 7.1.29
type LinkControlIOCapabilityRequestReplyInput struct {
	BDADDR [6]byte
	IOCapability uint8
	OOBDataPresent uint8
	AuthenticationRequirements uint8
}

func (i LinkControlIOCapabilityRequestReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.IOCapability)
	w.PutOne(i.OOBDataPresent)
	w.PutOne(i.AuthenticationRequirements)
	return w.Data()
}

// LinkControlIOCapabilityRequestReplyOutput represents the output of the command specified in Section 7.1.29
type LinkControlIOCapabilityRequestReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlIOCapabilityRequestReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlIOCapabilityRequestReplySync executes the command specified in Section 7.1.29 synchronously
func (c *Commands) LinkControlIOCapabilityRequestReplySync (params LinkControlIOCapabilityRequestReplyInput, result *LinkControlIOCapabilityRequestReplyOutput) (*LinkControlIOCapabilityRequestReplyOutput, error) {
	if result == nil {
		result = &LinkControlIOCapabilityRequestReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x002B}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlUserConfirmationRequestReplyInput represents the input of the command specified in Section 7.1.30
type LinkControlUserConfirmationRequestReplyInput struct {
	BDADDR [6]byte
}

func (i LinkControlUserConfirmationRequestReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	return w.Data()
}

// LinkControlUserConfirmationRequestReplyOutput represents the output of the command specified in Section 7.1.30
type LinkControlUserConfirmationRequestReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlUserConfirmationRequestReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlUserConfirmationRequestReplySync executes the command specified in Section 7.1.30 synchronously
func (c *Commands) LinkControlUserConfirmationRequestReplySync (params LinkControlUserConfirmationRequestReplyInput, result *LinkControlUserConfirmationRequestReplyOutput) (*LinkControlUserConfirmationRequestReplyOutput, error) {
	if result == nil {
		result = &LinkControlUserConfirmationRequestReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x002C}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlUserConfirmationRequestNegativeReplyInput represents the input of the command specified in Section 7.1.31
type LinkControlUserConfirmationRequestNegativeReplyInput struct {
	BDADDR [6]byte
}

func (i LinkControlUserConfirmationRequestNegativeReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	return w.Data()
}

// LinkControlUserConfirmationRequestNegativeReplyOutput represents the output of the command specified in Section 7.1.31
type LinkControlUserConfirmationRequestNegativeReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlUserConfirmationRequestNegativeReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlUserConfirmationRequestNegativeReplySync executes the command specified in Section 7.1.31 synchronously
func (c *Commands) LinkControlUserConfirmationRequestNegativeReplySync (params LinkControlUserConfirmationRequestNegativeReplyInput, result *LinkControlUserConfirmationRequestNegativeReplyOutput) (*LinkControlUserConfirmationRequestNegativeReplyOutput, error) {
	if result == nil {
		result = &LinkControlUserConfirmationRequestNegativeReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x002D}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlUserPasskeyRequestReplyInput represents the input of the command specified in Section 7.1.32
type LinkControlUserPasskeyRequestReplyInput struct {
	BDADDR [6]byte
	NumericValue uint32
}

func (i LinkControlUserPasskeyRequestReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	binary.LittleEndian.PutUint32(w.Put(4), i.NumericValue)
	return w.Data()
}

// LinkControlUserPasskeyRequestReplyOutput represents the output of the command specified in Section 7.1.32
type LinkControlUserPasskeyRequestReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlUserPasskeyRequestReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlUserPasskeyRequestReplySync executes the command specified in Section 7.1.32 synchronously
func (c *Commands) LinkControlUserPasskeyRequestReplySync (params LinkControlUserPasskeyRequestReplyInput, result *LinkControlUserPasskeyRequestReplyOutput) (*LinkControlUserPasskeyRequestReplyOutput, error) {
	if result == nil {
		result = &LinkControlUserPasskeyRequestReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x002E}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlUserPasskeyRequestNegativeReplyInput represents the input of the command specified in Section 7.1.33
type LinkControlUserPasskeyRequestNegativeReplyInput struct {
	BDADDR [6]byte
}

func (i LinkControlUserPasskeyRequestNegativeReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	return w.Data()
}

// LinkControlUserPasskeyRequestNegativeReplyOutput represents the output of the command specified in Section 7.1.33
type LinkControlUserPasskeyRequestNegativeReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlUserPasskeyRequestNegativeReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlUserPasskeyRequestNegativeReplySync executes the command specified in Section 7.1.33 synchronously
func (c *Commands) LinkControlUserPasskeyRequestNegativeReplySync (params LinkControlUserPasskeyRequestNegativeReplyInput, result *LinkControlUserPasskeyRequestNegativeReplyOutput) (*LinkControlUserPasskeyRequestNegativeReplyOutput, error) {
	if result == nil {
		result = &LinkControlUserPasskeyRequestNegativeReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x002F}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlRemoteOOBDataRequestReplyInput represents the input of the command specified in Section 7.1.34
type LinkControlRemoteOOBDataRequestReplyInput struct {
	BDADDR [6]byte
	C [16]byte
	R [16]byte
}

func (i LinkControlRemoteOOBDataRequestReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	copy(w.Put(16), i.C[:])
	copy(w.Put(16), i.R[:])
	return w.Data()
}

// LinkControlRemoteOOBDataRequestReplyOutput represents the output of the command specified in Section 7.1.34
type LinkControlRemoteOOBDataRequestReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlRemoteOOBDataRequestReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlRemoteOOBDataRequestReplySync executes the command specified in Section 7.1.34 synchronously
func (c *Commands) LinkControlRemoteOOBDataRequestReplySync (params LinkControlRemoteOOBDataRequestReplyInput, result *LinkControlRemoteOOBDataRequestReplyOutput) (*LinkControlRemoteOOBDataRequestReplyOutput, error) {
	if result == nil {
		result = &LinkControlRemoteOOBDataRequestReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0030}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlRemoteOOBDataRequestNegativeReplyInput represents the input of the command specified in Section 7.1.35
type LinkControlRemoteOOBDataRequestNegativeReplyInput struct {
	BDADDR [6]byte
}

func (i LinkControlRemoteOOBDataRequestNegativeReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	return w.Data()
}

// LinkControlRemoteOOBDataRequestNegativeReplyOutput represents the output of the command specified in Section 7.1.35
type LinkControlRemoteOOBDataRequestNegativeReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlRemoteOOBDataRequestNegativeReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlRemoteOOBDataRequestNegativeReplySync executes the command specified in Section 7.1.35 synchronously
func (c *Commands) LinkControlRemoteOOBDataRequestNegativeReplySync (params LinkControlRemoteOOBDataRequestNegativeReplyInput, result *LinkControlRemoteOOBDataRequestNegativeReplyOutput) (*LinkControlRemoteOOBDataRequestNegativeReplyOutput, error) {
	if result == nil {
		result = &LinkControlRemoteOOBDataRequestNegativeReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0033}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlIOCapabilityRequestNegativeReplyInput represents the input of the command specified in Section 7.1.36
type LinkControlIOCapabilityRequestNegativeReplyInput struct {
	BDADDR [6]byte
	Reason uint8
}

func (i LinkControlIOCapabilityRequestNegativeReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.Reason)
	return w.Data()
}

// LinkControlIOCapabilityRequestNegativeReplyOutput represents the output of the command specified in Section 7.1.36
type LinkControlIOCapabilityRequestNegativeReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlIOCapabilityRequestNegativeReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlIOCapabilityRequestNegativeReplySync executes the command specified in Section 7.1.36 synchronously
func (c *Commands) LinkControlIOCapabilityRequestNegativeReplySync (params LinkControlIOCapabilityRequestNegativeReplyInput, result *LinkControlIOCapabilityRequestNegativeReplyOutput) (*LinkControlIOCapabilityRequestNegativeReplyOutput, error) {
	if result == nil {
		result = &LinkControlIOCapabilityRequestNegativeReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0034}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlCreatePhysicalLinkInput represents the input of the command specified in Section 7.1.37
type LinkControlCreatePhysicalLinkInput struct {
	PhysicalLinkHandle uint8
	DedicatedAMPKeyLength uint8
	DedicatedAMPKeyType uint8
	DedicatedAMPKey []byte
}

func (i LinkControlCreatePhysicalLinkInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PhysicalLinkHandle)
	w.PutOne(i.DedicatedAMPKeyLength)
	w.PutOne(i.DedicatedAMPKeyType)
	w.PutSlice(i.DedicatedAMPKey)
	return w.Data()
}

// LinkControlCreatePhysicalLinkSync executes the command specified in Section 7.1.37 synchronously
func (c *Commands) LinkControlCreatePhysicalLinkSync (params LinkControlCreatePhysicalLinkInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0035}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlAcceptPhysicalLinkInput represents the input of the command specified in Section 7.1.38
type LinkControlAcceptPhysicalLinkInput struct {
	PhysicalLinkHandle uint8
	DedicatedAMPKeyLength uint8
	DedicatedAMPKeyType uint8
	DedicatedAMPKey []byte
}

func (i LinkControlAcceptPhysicalLinkInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PhysicalLinkHandle)
	w.PutOne(i.DedicatedAMPKeyLength)
	w.PutOne(i.DedicatedAMPKeyType)
	w.PutSlice(i.DedicatedAMPKey)
	return w.Data()
}

// LinkControlAcceptPhysicalLinkSync executes the command specified in Section 7.1.38 synchronously
func (c *Commands) LinkControlAcceptPhysicalLinkSync (params LinkControlAcceptPhysicalLinkInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0036}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlDisconnectPhysicalLinkInput represents the input of the command specified in Section 7.1.39
type LinkControlDisconnectPhysicalLinkInput struct {
	PhysicalLinkHandle uint8
	Reason uint8
}

func (i LinkControlDisconnectPhysicalLinkInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PhysicalLinkHandle)
	w.PutOne(i.Reason)
	return w.Data()
}

// LinkControlDisconnectPhysicalLinkSync executes the command specified in Section 7.1.39 synchronously
func (c *Commands) LinkControlDisconnectPhysicalLinkSync (params LinkControlDisconnectPhysicalLinkInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0037}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlCreateLogicalLinkInput represents the input of the command specified in Section 7.1.40
type LinkControlCreateLogicalLinkInput struct {
	PhysicalLinkHandle uint8
	TXFlowSpec [16]byte
	RXFlowSpec [16]byte
}

func (i LinkControlCreateLogicalLinkInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PhysicalLinkHandle)
	copy(w.Put(16), i.TXFlowSpec[:])
	copy(w.Put(16), i.RXFlowSpec[:])
	return w.Data()
}

// LinkControlCreateLogicalLinkSync executes the command specified in Section 7.1.40 synchronously
func (c *Commands) LinkControlCreateLogicalLinkSync (params LinkControlCreateLogicalLinkInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0038}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlAcceptLogicalLinkInput represents the input of the command specified in Section 7.1.41
type LinkControlAcceptLogicalLinkInput struct {
	PhysicalLinkHandle uint8
	TXFlowSpec [16]byte
	RXFlowSpec [16]byte
}

func (i LinkControlAcceptLogicalLinkInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PhysicalLinkHandle)
	copy(w.Put(16), i.TXFlowSpec[:])
	copy(w.Put(16), i.RXFlowSpec[:])
	return w.Data()
}

// LinkControlAcceptLogicalLinkSync executes the command specified in Section 7.1.41 synchronously
func (c *Commands) LinkControlAcceptLogicalLinkSync (params LinkControlAcceptLogicalLinkInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0039}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlDisconnectLogicalLinkInput represents the input of the command specified in Section 7.1.42
type LinkControlDisconnectLogicalLinkInput struct {
	LogicalLinkHandle uint16
}

func (i LinkControlDisconnectLogicalLinkInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.LogicalLinkHandle)
	return w.Data()
}

// LinkControlDisconnectLogicalLinkSync executes the command specified in Section 7.1.42 synchronously
func (c *Commands) LinkControlDisconnectLogicalLinkSync (params LinkControlDisconnectLogicalLinkInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x003A}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlLogicalLinkCancelInput represents the input of the command specified in Section 7.1.43
type LinkControlLogicalLinkCancelInput struct {
	PhysicalLinkHandle uint8
	TXFlowSpecID uint8
}

func (i LinkControlLogicalLinkCancelInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PhysicalLinkHandle)
	w.PutOne(i.TXFlowSpecID)
	return w.Data()
}

// LinkControlLogicalLinkCancelOutput represents the output of the command specified in Section 7.1.43
type LinkControlLogicalLinkCancelOutput struct {
	Status uint8
	PhysicalLinkHandle uint8
	TXFlowSpecID uint8
}

func (o *LinkControlLogicalLinkCancelOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.PhysicalLinkHandle = r.GetOne()
	o.TXFlowSpecID = r.GetOne()
	return r.Valid()
}

// LinkControlLogicalLinkCancelSync executes the command specified in Section 7.1.43 synchronously
func (c *Commands) LinkControlLogicalLinkCancelSync (params LinkControlLogicalLinkCancelInput, result *LinkControlLogicalLinkCancelOutput) (*LinkControlLogicalLinkCancelOutput, error) {
	if result == nil {
		result = &LinkControlLogicalLinkCancelOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x003B}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlFlowSpecModifyInput represents the input of the command specified in Section 7.1.44
type LinkControlFlowSpecModifyInput struct {
	Handle uint16
	TXFlowSpec [16]byte
	RXFlowSpec [16]byte
}

func (i LinkControlFlowSpecModifyInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.Handle)
	copy(w.Put(16), i.TXFlowSpec[:])
	copy(w.Put(16), i.RXFlowSpec[:])
	return w.Data()
}

// LinkControlFlowSpecModifySync executes the command specified in Section 7.1.44 synchronously
func (c *Commands) LinkControlFlowSpecModifySync (params LinkControlFlowSpecModifyInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x003C}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlEnhancedSetupSynchronousConnectionInput represents the input of the command specified in Section 7.1.45
type LinkControlEnhancedSetupSynchronousConnectionInput struct {
	ConnectionHandle uint16
	TransmitBandwidth uint32
	ReceiveBandwidth uint32
	TransmitCodingFormat [5]byte
	ReceiveCodingFormat [5]byte
	TransmitCodecFrameSize uint16
	ReceiveCodecFrameSize uint16
	InputBandwidth uint32
	OutputBandwidth uint32
	InputCodingFormat [5]byte
	OutputCodingFormat [5]byte
	InputCodedDataSize uint16
	OutputCodedDataSize uint16
	InputPCMDataFormat uint8
	OutputPCMDataFormat uint8
	InputPCMSamplePayloadMSBPosition uint8
	OutputPCMSamplePayloadMSBPosition uint8
	InputDataPath uint8
	OutputDataPath uint8
	InputTransportUnitSize uint8
	OutputTransportUnitSize uint8
	MaxLatency uint16
	PacketType uint16
	RetransmissionEffort uint8
}

func (i LinkControlEnhancedSetupSynchronousConnectionInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint32(w.Put(4), i.TransmitBandwidth)
	binary.LittleEndian.PutUint32(w.Put(4), i.ReceiveBandwidth)
	copy(w.Put(5), i.TransmitCodingFormat[:])
	copy(w.Put(5), i.ReceiveCodingFormat[:])
	binary.LittleEndian.PutUint16(w.Put(2), i.TransmitCodecFrameSize)
	binary.LittleEndian.PutUint16(w.Put(2), i.ReceiveCodecFrameSize)
	binary.LittleEndian.PutUint32(w.Put(4), i.InputBandwidth)
	binary.LittleEndian.PutUint32(w.Put(4), i.OutputBandwidth)
	copy(w.Put(5), i.InputCodingFormat[:])
	copy(w.Put(5), i.OutputCodingFormat[:])
	binary.LittleEndian.PutUint16(w.Put(2), i.InputCodedDataSize)
	binary.LittleEndian.PutUint16(w.Put(2), i.OutputCodedDataSize)
	w.PutOne(i.InputPCMDataFormat)
	w.PutOne(i.OutputPCMDataFormat)
	w.PutOne(i.InputPCMSamplePayloadMSBPosition)
	w.PutOne(i.OutputPCMSamplePayloadMSBPosition)
	w.PutOne(i.InputDataPath)
	w.PutOne(i.OutputDataPath)
	w.PutOne(i.InputTransportUnitSize)
	w.PutOne(i.OutputTransportUnitSize)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxLatency)
	binary.LittleEndian.PutUint16(w.Put(2), i.PacketType)
	w.PutOne(i.RetransmissionEffort)
	return w.Data()
}

// LinkControlEnhancedSetupSynchronousConnectionSync executes the command specified in Section 7.1.45 synchronously
func (c *Commands) LinkControlEnhancedSetupSynchronousConnectionSync (params LinkControlEnhancedSetupSynchronousConnectionInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x003D}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlEnhancedAcceptSynchronousConnectionRequestInput represents the input of the command specified in Section 7.1.46
type LinkControlEnhancedAcceptSynchronousConnectionRequestInput struct {
	BDADDR [6]byte
	TransmitBandwidth uint32
	ReceiveBandwidth uint32
	TransmitCodingFormat [5]byte
	ReceiveCodingFormat [5]byte
	TransmitCodecFrameSize uint16
	ReceiveCodecFrameSize uint16
	InputBandwidth uint32
	OutputBandwidth uint32
	InputCodingFormat [5]byte
	OutputCodingFormat [5]byte
	InputCodedDataSize uint16
	OutputCodedDataSize uint16
	InputPCMDataFormat uint8
	OutputPCMDataFormat uint8
	InputPCMSamplePayloadMSBPosition uint8
	OutputPCMSamplePayloadMSBPosition uint8
	InputDataPath uint8
	OutputDataPath uint8
	InputTransportUnitSize uint8
	OutputTransportUnitSize uint8
	MaxLatency uint16
	PacketType uint16
	RetransmissionEffort uint8
}

func (i LinkControlEnhancedAcceptSynchronousConnectionRequestInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	binary.LittleEndian.PutUint32(w.Put(4), i.TransmitBandwidth)
	binary.LittleEndian.PutUint32(w.Put(4), i.ReceiveBandwidth)
	copy(w.Put(5), i.TransmitCodingFormat[:])
	copy(w.Put(5), i.ReceiveCodingFormat[:])
	binary.LittleEndian.PutUint16(w.Put(2), i.TransmitCodecFrameSize)
	binary.LittleEndian.PutUint16(w.Put(2), i.ReceiveCodecFrameSize)
	binary.LittleEndian.PutUint32(w.Put(4), i.InputBandwidth)
	binary.LittleEndian.PutUint32(w.Put(4), i.OutputBandwidth)
	copy(w.Put(5), i.InputCodingFormat[:])
	copy(w.Put(5), i.OutputCodingFormat[:])
	binary.LittleEndian.PutUint16(w.Put(2), i.InputCodedDataSize)
	binary.LittleEndian.PutUint16(w.Put(2), i.OutputCodedDataSize)
	w.PutOne(i.InputPCMDataFormat)
	w.PutOne(i.OutputPCMDataFormat)
	w.PutOne(i.InputPCMSamplePayloadMSBPosition)
	w.PutOne(i.OutputPCMSamplePayloadMSBPosition)
	w.PutOne(i.InputDataPath)
	w.PutOne(i.OutputDataPath)
	w.PutOne(i.InputTransportUnitSize)
	w.PutOne(i.OutputTransportUnitSize)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxLatency)
	binary.LittleEndian.PutUint16(w.Put(2), i.PacketType)
	w.PutOne(i.RetransmissionEffort)
	return w.Data()
}

// LinkControlEnhancedAcceptSynchronousConnectionRequestSync executes the command specified in Section 7.1.46 synchronously
func (c *Commands) LinkControlEnhancedAcceptSynchronousConnectionRequestSync (params LinkControlEnhancedAcceptSynchronousConnectionRequestInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x003E}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlTruncatedPageInput represents the input of the command specified in Section 7.1.47
type LinkControlTruncatedPageInput struct {
	BDADDR [6]byte
	PageScanRepetitionMode uint8
	ClockOffset uint16
}

func (i LinkControlTruncatedPageInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.PageScanRepetitionMode)
	binary.LittleEndian.PutUint16(w.Put(2), i.ClockOffset)
	return w.Data()
}

// LinkControlTruncatedPageSync executes the command specified in Section 7.1.47 synchronously
func (c *Commands) LinkControlTruncatedPageSync (params LinkControlTruncatedPageInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x003F}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlTruncatedPageCancelInput represents the input of the command specified in Section 7.1.48
type LinkControlTruncatedPageCancelInput struct {
	BDADDR [6]byte
}

func (i LinkControlTruncatedPageCancelInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	return w.Data()
}

// LinkControlTruncatedPageCancelOutput represents the output of the command specified in Section 7.1.48
type LinkControlTruncatedPageCancelOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlTruncatedPageCancelOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlTruncatedPageCancelSync executes the command specified in Section 7.1.48 synchronously
func (c *Commands) LinkControlTruncatedPageCancelSync (params LinkControlTruncatedPageCancelInput, result *LinkControlTruncatedPageCancelOutput) (*LinkControlTruncatedPageCancelOutput, error) {
	if result == nil {
		result = &LinkControlTruncatedPageCancelOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0040}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlSetConnectionlessSlaveBroadcastInput represents the input of the command specified in Section 7.1.49
type LinkControlSetConnectionlessSlaveBroadcastInput struct {
	Enable uint8
	LTADDR uint8
	LPOAllowed uint8
	PacketType uint16
	IntervalMin uint16
	IntervalMax uint16
	SupervisionTimeout uint16
}

func (i LinkControlSetConnectionlessSlaveBroadcastInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.Enable)
	w.PutOne(i.LTADDR)
	w.PutOne(i.LPOAllowed)
	binary.LittleEndian.PutUint16(w.Put(2), i.PacketType)
	binary.LittleEndian.PutUint16(w.Put(2), i.IntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.IntervalMax)
	binary.LittleEndian.PutUint16(w.Put(2), i.SupervisionTimeout)
	return w.Data()
}

// LinkControlSetConnectionlessSlaveBroadcastOutput represents the output of the command specified in Section 7.1.49
type LinkControlSetConnectionlessSlaveBroadcastOutput struct {
	Status uint8
	LTADDR uint8
	Interval uint16
}

func (o *LinkControlSetConnectionlessSlaveBroadcastOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.LTADDR = r.GetOne()
	o.Interval = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LinkControlSetConnectionlessSlaveBroadcastSync executes the command specified in Section 7.1.49 synchronously
func (c *Commands) LinkControlSetConnectionlessSlaveBroadcastSync (params LinkControlSetConnectionlessSlaveBroadcastInput, result *LinkControlSetConnectionlessSlaveBroadcastOutput) (*LinkControlSetConnectionlessSlaveBroadcastOutput, error) {
	if result == nil {
		result = &LinkControlSetConnectionlessSlaveBroadcastOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0041}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlSetConnectionlessSlaveBroadcastReceiveInput represents the input of the command specified in Section 7.1.50
type LinkControlSetConnectionlessSlaveBroadcastReceiveInput struct {
	Enable uint8
	BDADDR [6]byte
	LTADDR uint8
	Interval uint16
	ClockOffset uint32
	CSBsupervisionTO uint16
	RemoteTimingAccuracy uint8
	Skip uint8
	PacketType uint16
	AFHChannelMap [10]byte
}

func (i LinkControlSetConnectionlessSlaveBroadcastReceiveInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.Enable)
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.LTADDR)
	binary.LittleEndian.PutUint16(w.Put(2), i.Interval)
	binary.LittleEndian.PutUint32(w.Put(4), i.ClockOffset)
	binary.LittleEndian.PutUint16(w.Put(2), i.CSBsupervisionTO)
	w.PutOne(i.RemoteTimingAccuracy)
	w.PutOne(i.Skip)
	binary.LittleEndian.PutUint16(w.Put(2), i.PacketType)
	copy(w.Put(10), i.AFHChannelMap[:])
	return w.Data()
}

// LinkControlSetConnectionlessSlaveBroadcastReceiveOutput represents the output of the command specified in Section 7.1.50
type LinkControlSetConnectionlessSlaveBroadcastReceiveOutput struct {
	Status uint8
	BDADDR [6]byte
	LTADDR uint8
}

func (o *LinkControlSetConnectionlessSlaveBroadcastReceiveOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	o.LTADDR = r.GetOne()
	return r.Valid()
}

// LinkControlSetConnectionlessSlaveBroadcastReceiveSync executes the command specified in Section 7.1.50 synchronously
func (c *Commands) LinkControlSetConnectionlessSlaveBroadcastReceiveSync (params LinkControlSetConnectionlessSlaveBroadcastReceiveInput, result *LinkControlSetConnectionlessSlaveBroadcastReceiveOutput) (*LinkControlSetConnectionlessSlaveBroadcastReceiveOutput, error) {
	if result == nil {
		result = &LinkControlSetConnectionlessSlaveBroadcastReceiveOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0042}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

// LinkControlStartSynchronizationTrainInput represents the input of the command specified in Section 7.1.51
type LinkControlStartSynchronizationTrainInput struct {
	BDADDR [6]byte
	SyncScanTimeout uint16
	SyncScanWindow uint16
	SyncScanInterval uint16
}

func (i LinkControlStartSynchronizationTrainInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncScanTimeout)
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncScanWindow)
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncScanInterval)
	return w.Data()
}

// LinkControlStartSynchronizationTrainSync executes the command specified in Section 7.1.51 synchronously
func (c *Commands) LinkControlStartSynchronizationTrainSync (params LinkControlStartSynchronizationTrainInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0044}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// LinkControlRemoteOOBExtendedDataRequestReplyInput represents the input of the command specified in Section 7.1.53
type LinkControlRemoteOOBExtendedDataRequestReplyInput struct {
	BDADDR [6]byte
	C192 [16]byte
	R192 [16]byte
	C256 [16]byte
	R256 [16]byte
}

func (i LinkControlRemoteOOBExtendedDataRequestReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	copy(w.Put(16), i.C192[:])
	copy(w.Put(16), i.R192[:])
	copy(w.Put(16), i.C256[:])
	copy(w.Put(16), i.R256[:])
	return w.Data()
}

// LinkControlRemoteOOBExtendedDataRequestReplyOutput represents the output of the command specified in Section 7.1.53
type LinkControlRemoteOOBExtendedDataRequestReplyOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *LinkControlRemoteOOBExtendedDataRequestReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// LinkControlRemoteOOBExtendedDataRequestReplySync executes the command specified in Section 7.1.53 synchronously
func (c *Commands) LinkControlRemoteOOBExtendedDataRequestReplySync (params LinkControlRemoteOOBExtendedDataRequestReplyInput, result *LinkControlRemoteOOBExtendedDataRequestReplyOutput) (*LinkControlRemoteOOBExtendedDataRequestReplyOutput, error) {
	if result == nil {
		result = &LinkControlRemoteOOBExtendedDataRequestReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 1, OCF: 0x0045}, nil)
	if err != nil {
		return result, err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err := c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return result, err
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return result, err
}

