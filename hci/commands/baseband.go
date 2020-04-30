package hcicommands

import (
	"encoding/binary"
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
)

// BasebandSetEventMaskInput represents the input of the command specified in Section 7.3.1
type BasebandSetEventMaskInput struct {
	EventMask uint64
}

func (i BasebandSetEventMaskInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint64(w.Put(8), i.EventMask)
	return w.Data()
}

// BasebandSetEventMaskSync executes the command specified in Section 7.3.1 synchronously
func (c *Commands) BasebandSetEventMaskSync (params BasebandSetEventMaskInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0001}, nil)
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

// BasebandResetSync executes the command specified in Section 7.3.2 synchronously
func (c *Commands) BasebandResetSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0003}, nil)
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

// BasebandFlushInput represents the input of the command specified in Section 7.3.4
type BasebandFlushInput struct {
	ConnectionHandle uint16
}

func (i BasebandFlushInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// BasebandFlushOutput represents the output of the command specified in Section 7.3.4
type BasebandFlushOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *BasebandFlushOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandFlushSync executes the command specified in Section 7.3.4 synchronously
func (c *Commands) BasebandFlushSync (params BasebandFlushInput, result *BasebandFlushOutput) (*BasebandFlushOutput, error) {
	if result == nil {
		result = &BasebandFlushOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0008}, nil)
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

// BasebandReadPINTypeOutput represents the output of the command specified in Section 7.3.5
type BasebandReadPINTypeOutput struct {
	Status uint8
	PINType uint8
}

func (o *BasebandReadPINTypeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.PINType = r.GetOne()
	return r.Valid()
}

// BasebandReadPINTypeSync executes the command specified in Section 7.3.5 synchronously
func (c *Commands) BasebandReadPINTypeSync (result *BasebandReadPINTypeOutput) (*BasebandReadPINTypeOutput, error) {
	if result == nil {
		result = &BasebandReadPINTypeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0009}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWritePINTypeInput represents the input of the command specified in Section 7.3.6
type BasebandWritePINTypeInput struct {
	PINType uint8
}

func (i BasebandWritePINTypeInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PINType)
	return w.Data()
}

// BasebandWritePINTypeSync executes the command specified in Section 7.3.6 synchronously
func (c *Commands) BasebandWritePINTypeSync (params BasebandWritePINTypeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x000A}, nil)
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

// BasebandReadStoredLinkKeyInput represents the input of the command specified in Section 7.3.8
type BasebandReadStoredLinkKeyInput struct {
	BDADDR [6]byte
	ReadAll uint8
}

func (i BasebandReadStoredLinkKeyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.ReadAll)
	return w.Data()
}

// BasebandReadStoredLinkKeyOutput represents the output of the command specified in Section 7.3.8
type BasebandReadStoredLinkKeyOutput struct {
	Status uint8
	MaxNumKeys uint16
	NumKeysRead uint16
}

func (o *BasebandReadStoredLinkKeyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.MaxNumKeys = binary.LittleEndian.Uint16(r.Get(2))
	o.NumKeysRead = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadStoredLinkKeySync executes the command specified in Section 7.3.8 synchronously
func (c *Commands) BasebandReadStoredLinkKeySync (params BasebandReadStoredLinkKeyInput, result *BasebandReadStoredLinkKeyOutput) (*BasebandReadStoredLinkKeyOutput, error) {
	if result == nil {
		result = &BasebandReadStoredLinkKeyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x000D}, nil)
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

// BasebandWriteStoredLinkKeyInput represents the input of the command specified in Section 7.3.9
type BasebandWriteStoredLinkKeyInput struct {
	NumKeysToWrite uint8
	BDADDR [][6]byte
	LinkKey [][16]byte
}

func (i BasebandWriteStoredLinkKeyInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.NumKeysToWrite)
	if len(i.BDADDR) != int(i.NumKeysToWrite) {
		panic("len(i.BDADDR) != int(i.NumKeysToWrite)")
	}
	for _, m := range i.BDADDR {
		copy(w.Put(6), m[:])
	}
	if len(i.LinkKey) != int(i.NumKeysToWrite) {
		panic("len(i.LinkKey) != int(i.NumKeysToWrite)")
	}
	for _, m := range i.LinkKey {
		copy(w.Put(16), m[:])
	}
	return w.Data()
}

// BasebandWriteStoredLinkKeyOutput represents the output of the command specified in Section 7.3.9
type BasebandWriteStoredLinkKeyOutput struct {
	Status uint8
	NumKeysWritten uint8
}

func (o *BasebandWriteStoredLinkKeyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.NumKeysWritten = r.GetOne()
	return r.Valid()
}

// BasebandWriteStoredLinkKeySync executes the command specified in Section 7.3.9 synchronously
func (c *Commands) BasebandWriteStoredLinkKeySync (params BasebandWriteStoredLinkKeyInput, result *BasebandWriteStoredLinkKeyOutput) (*BasebandWriteStoredLinkKeyOutput, error) {
	if result == nil {
		result = &BasebandWriteStoredLinkKeyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0011}, nil)
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

// BasebandDeleteStoredLinkKeyInput represents the input of the command specified in Section 7.3.10
type BasebandDeleteStoredLinkKeyInput struct {
	BDADDR [6]byte
	DeleteAll uint8
}

func (i BasebandDeleteStoredLinkKeyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.DeleteAll)
	return w.Data()
}

// BasebandDeleteStoredLinkKeyOutput represents the output of the command specified in Section 7.3.10
type BasebandDeleteStoredLinkKeyOutput struct {
	Status uint8
	NumKeysDeleted uint16
}

func (o *BasebandDeleteStoredLinkKeyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.NumKeysDeleted = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandDeleteStoredLinkKeySync executes the command specified in Section 7.3.10 synchronously
func (c *Commands) BasebandDeleteStoredLinkKeySync (params BasebandDeleteStoredLinkKeyInput, result *BasebandDeleteStoredLinkKeyOutput) (*BasebandDeleteStoredLinkKeyOutput, error) {
	if result == nil {
		result = &BasebandDeleteStoredLinkKeyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0012}, nil)
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

// BasebandWriteLocalNameInput represents the input of the command specified in Section 7.3.11
type BasebandWriteLocalNameInput struct {
	LocalName [248]byte
}

func (i BasebandWriteLocalNameInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(248), i.LocalName[:])
	return w.Data()
}

// BasebandWriteLocalNameSync executes the command specified in Section 7.3.11 synchronously
func (c *Commands) BasebandWriteLocalNameSync (params BasebandWriteLocalNameInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0013}, nil)
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

// BasebandReadLocalNameOutput represents the output of the command specified in Section 7.3.12
type BasebandReadLocalNameOutput struct {
	Status uint8
	LocalName [248]byte
}

func (o *BasebandReadLocalNameOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.LocalName[:], r.Get(248))
	return r.Valid()
}

// BasebandReadLocalNameSync executes the command specified in Section 7.3.12 synchronously
func (c *Commands) BasebandReadLocalNameSync (result *BasebandReadLocalNameOutput) (*BasebandReadLocalNameOutput, error) {
	if result == nil {
		result = &BasebandReadLocalNameOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0014}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandReadConnectionAcceptTimeoutOutput represents the output of the command specified in Section 7.3.13
type BasebandReadConnectionAcceptTimeoutOutput struct {
	Status uint8
	ConnectionAcceptTimeout uint16
}

func (o *BasebandReadConnectionAcceptTimeoutOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionAcceptTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadConnectionAcceptTimeoutSync executes the command specified in Section 7.3.13 synchronously
func (c *Commands) BasebandReadConnectionAcceptTimeoutSync (result *BasebandReadConnectionAcceptTimeoutOutput) (*BasebandReadConnectionAcceptTimeoutOutput, error) {
	if result == nil {
		result = &BasebandReadConnectionAcceptTimeoutOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0015}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteConnectionAcceptTimeoutInput represents the input of the command specified in Section 7.3.14
type BasebandWriteConnectionAcceptTimeoutInput struct {
	ConnectionAcceptTimeout uint16
}

func (i BasebandWriteConnectionAcceptTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionAcceptTimeout)
	return w.Data()
}

// BasebandWriteConnectionAcceptTimeoutSync executes the command specified in Section 7.3.14 synchronously
func (c *Commands) BasebandWriteConnectionAcceptTimeoutSync (params BasebandWriteConnectionAcceptTimeoutInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0016}, nil)
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

// BasebandReadPageTimeoutOutput represents the output of the command specified in Section 7.3.15
type BasebandReadPageTimeoutOutput struct {
	Status uint8
	PageTimeout uint16
}

func (o *BasebandReadPageTimeoutOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.PageTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadPageTimeoutSync executes the command specified in Section 7.3.15 synchronously
func (c *Commands) BasebandReadPageTimeoutSync (result *BasebandReadPageTimeoutOutput) (*BasebandReadPageTimeoutOutput, error) {
	if result == nil {
		result = &BasebandReadPageTimeoutOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0017}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWritePageTimeoutInput represents the input of the command specified in Section 7.3.16
type BasebandWritePageTimeoutInput struct {
	PageTimeout uint16
}

func (i BasebandWritePageTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.PageTimeout)
	return w.Data()
}

// BasebandWritePageTimeoutSync executes the command specified in Section 7.3.16 synchronously
func (c *Commands) BasebandWritePageTimeoutSync (params BasebandWritePageTimeoutInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0018}, nil)
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

// BasebandReadScanEnableOutput represents the output of the command specified in Section 7.3.17
type BasebandReadScanEnableOutput struct {
	Status uint8
	ScanEnable uint8
}

func (o *BasebandReadScanEnableOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ScanEnable = r.GetOne()
	return r.Valid()
}

// BasebandReadScanEnableSync executes the command specified in Section 7.3.17 synchronously
func (c *Commands) BasebandReadScanEnableSync (result *BasebandReadScanEnableOutput) (*BasebandReadScanEnableOutput, error) {
	if result == nil {
		result = &BasebandReadScanEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0019}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteScanEnableInput represents the input of the command specified in Section 7.3.18
type BasebandWriteScanEnableInput struct {
	ScanEnable uint8
}

func (i BasebandWriteScanEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.ScanEnable)
	return w.Data()
}

// BasebandWriteScanEnableSync executes the command specified in Section 7.3.18 synchronously
func (c *Commands) BasebandWriteScanEnableSync (params BasebandWriteScanEnableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x001A}, nil)
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

// BasebandReadPageScanActivityOutput represents the output of the command specified in Section 7.3.19
type BasebandReadPageScanActivityOutput struct {
	Status uint8
	PageScanInterval uint16
	PageScanWindow uint16
}

func (o *BasebandReadPageScanActivityOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.PageScanInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.PageScanWindow = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadPageScanActivitySync executes the command specified in Section 7.3.19 synchronously
func (c *Commands) BasebandReadPageScanActivitySync (result *BasebandReadPageScanActivityOutput) (*BasebandReadPageScanActivityOutput, error) {
	if result == nil {
		result = &BasebandReadPageScanActivityOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x001B}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWritePageScanActivityInput represents the input of the command specified in Section 7.3.20
type BasebandWritePageScanActivityInput struct {
	PageScanInterval uint16
	PageScanWindow uint16
}

func (i BasebandWritePageScanActivityInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.PageScanInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.PageScanWindow)
	return w.Data()
}

// BasebandWritePageScanActivitySync executes the command specified in Section 7.3.20 synchronously
func (c *Commands) BasebandWritePageScanActivitySync (params BasebandWritePageScanActivityInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x001C}, nil)
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

// BasebandReadInquiryScanActivityOutput represents the output of the command specified in Section 7.3.21
type BasebandReadInquiryScanActivityOutput struct {
	Status uint8
	InquiryScanInterval uint16
	InquiryScanWindow uint16
}

func (o *BasebandReadInquiryScanActivityOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.InquiryScanInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.InquiryScanWindow = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadInquiryScanActivitySync executes the command specified in Section 7.3.21 synchronously
func (c *Commands) BasebandReadInquiryScanActivitySync (result *BasebandReadInquiryScanActivityOutput) (*BasebandReadInquiryScanActivityOutput, error) {
	if result == nil {
		result = &BasebandReadInquiryScanActivityOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x001D}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteInquiryScanActivityInput represents the input of the command specified in Section 7.3.22
type BasebandWriteInquiryScanActivityInput struct {
	InquiryScanInterval uint16
	InquiryScanWindow uint16
}

func (i BasebandWriteInquiryScanActivityInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.InquiryScanInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.InquiryScanWindow)
	return w.Data()
}

// BasebandWriteInquiryScanActivitySync executes the command specified in Section 7.3.22 synchronously
func (c *Commands) BasebandWriteInquiryScanActivitySync (params BasebandWriteInquiryScanActivityInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x001E}, nil)
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

// BasebandReadAuthenticationEnableOutput represents the output of the command specified in Section 7.3.23
type BasebandReadAuthenticationEnableOutput struct {
	Status uint8
	AuthenticationEnable uint8
}

func (o *BasebandReadAuthenticationEnableOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.AuthenticationEnable = r.GetOne()
	return r.Valid()
}

// BasebandReadAuthenticationEnableSync executes the command specified in Section 7.3.23 synchronously
func (c *Commands) BasebandReadAuthenticationEnableSync (result *BasebandReadAuthenticationEnableOutput) (*BasebandReadAuthenticationEnableOutput, error) {
	if result == nil {
		result = &BasebandReadAuthenticationEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x001F}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteAuthenticationEnableInput represents the input of the command specified in Section 7.3.24
type BasebandWriteAuthenticationEnableInput struct {
	AuthenticationEnable uint8
}

func (i BasebandWriteAuthenticationEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AuthenticationEnable)
	return w.Data()
}

// BasebandWriteAuthenticationEnableSync executes the command specified in Section 7.3.24 synchronously
func (c *Commands) BasebandWriteAuthenticationEnableSync (params BasebandWriteAuthenticationEnableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0020}, nil)
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

// BasebandReadClassofDeviceOutput represents the output of the command specified in Section 7.3.25
type BasebandReadClassofDeviceOutput struct {
	Status uint8
	ClassOfDevice uint32
}

func (o *BasebandReadClassofDeviceOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ClassOfDevice = decodeUint24(r.Get(3))
	return r.Valid()
}

// BasebandReadClassofDeviceSync executes the command specified in Section 7.3.25 synchronously
func (c *Commands) BasebandReadClassofDeviceSync (result *BasebandReadClassofDeviceOutput) (*BasebandReadClassofDeviceOutput, error) {
	if result == nil {
		result = &BasebandReadClassofDeviceOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0023}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteClassofDeviceInput represents the input of the command specified in Section 7.3.26
type BasebandWriteClassofDeviceInput struct {
	ClassOfDevice uint32
}

func (i BasebandWriteClassofDeviceInput) encode(data []byte) []byte {
	w := writer{data: data};
	encodeUint24(w.Put(3), i.ClassOfDevice)
	return w.Data()
}

// BasebandWriteClassofDeviceSync executes the command specified in Section 7.3.26 synchronously
func (c *Commands) BasebandWriteClassofDeviceSync (params BasebandWriteClassofDeviceInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0024}, nil)
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

// BasebandReadVoiceSettingOutput represents the output of the command specified in Section 7.3.27
type BasebandReadVoiceSettingOutput struct {
	Status uint8
	VoiceSetting uint16
}

func (o *BasebandReadVoiceSettingOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.VoiceSetting = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadVoiceSettingSync executes the command specified in Section 7.3.27 synchronously
func (c *Commands) BasebandReadVoiceSettingSync (result *BasebandReadVoiceSettingOutput) (*BasebandReadVoiceSettingOutput, error) {
	if result == nil {
		result = &BasebandReadVoiceSettingOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0025}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteVoiceSettingInput represents the input of the command specified in Section 7.3.28
type BasebandWriteVoiceSettingInput struct {
	VoiceSetting uint16
}

func (i BasebandWriteVoiceSettingInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.VoiceSetting)
	return w.Data()
}

// BasebandWriteVoiceSettingSync executes the command specified in Section 7.3.28 synchronously
func (c *Commands) BasebandWriteVoiceSettingSync (params BasebandWriteVoiceSettingInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0026}, nil)
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

// BasebandReadAutomaticFlushTimeoutInput represents the input of the command specified in Section 7.3.29
type BasebandReadAutomaticFlushTimeoutInput struct {
	ConnectionHandle uint16
}

func (i BasebandReadAutomaticFlushTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// BasebandReadAutomaticFlushTimeoutOutput represents the output of the command specified in Section 7.3.29
type BasebandReadAutomaticFlushTimeoutOutput struct {
	Status uint8
	ConnectionHandle uint16
	FlushTimeout uint16
}

func (o *BasebandReadAutomaticFlushTimeoutOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.FlushTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadAutomaticFlushTimeoutSync executes the command specified in Section 7.3.29 synchronously
func (c *Commands) BasebandReadAutomaticFlushTimeoutSync (params BasebandReadAutomaticFlushTimeoutInput, result *BasebandReadAutomaticFlushTimeoutOutput) (*BasebandReadAutomaticFlushTimeoutOutput, error) {
	if result == nil {
		result = &BasebandReadAutomaticFlushTimeoutOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0027}, nil)
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

// BasebandWriteAutomaticFlushTimeoutInput represents the input of the command specified in Section 7.3.30
type BasebandWriteAutomaticFlushTimeoutInput struct {
	ConnectionHandle uint16
	FlushTimeout uint16
}

func (i BasebandWriteAutomaticFlushTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.FlushTimeout)
	return w.Data()
}

// BasebandWriteAutomaticFlushTimeoutOutput represents the output of the command specified in Section 7.3.30
type BasebandWriteAutomaticFlushTimeoutOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *BasebandWriteAutomaticFlushTimeoutOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandWriteAutomaticFlushTimeoutSync executes the command specified in Section 7.3.30 synchronously
func (c *Commands) BasebandWriteAutomaticFlushTimeoutSync (params BasebandWriteAutomaticFlushTimeoutInput, result *BasebandWriteAutomaticFlushTimeoutOutput) (*BasebandWriteAutomaticFlushTimeoutOutput, error) {
	if result == nil {
		result = &BasebandWriteAutomaticFlushTimeoutOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0028}, nil)
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

// BasebandReadNumBroadcastRetransmissionsOutput represents the output of the command specified in Section 7.3.31
type BasebandReadNumBroadcastRetransmissionsOutput struct {
	Status uint8
	NumBroadcastRetransmissions uint8
}

func (o *BasebandReadNumBroadcastRetransmissionsOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.NumBroadcastRetransmissions = r.GetOne()
	return r.Valid()
}

// BasebandReadNumBroadcastRetransmissionsSync executes the command specified in Section 7.3.31 synchronously
func (c *Commands) BasebandReadNumBroadcastRetransmissionsSync (result *BasebandReadNumBroadcastRetransmissionsOutput) (*BasebandReadNumBroadcastRetransmissionsOutput, error) {
	if result == nil {
		result = &BasebandReadNumBroadcastRetransmissionsOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0029}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteNumBroadcastRetransmissionsInput represents the input of the command specified in Section 7.3.32
type BasebandWriteNumBroadcastRetransmissionsInput struct {
	NumBroadcastRetransmissions uint8
}

func (i BasebandWriteNumBroadcastRetransmissionsInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.NumBroadcastRetransmissions)
	return w.Data()
}

// BasebandWriteNumBroadcastRetransmissionsSync executes the command specified in Section 7.3.32 synchronously
func (c *Commands) BasebandWriteNumBroadcastRetransmissionsSync (params BasebandWriteNumBroadcastRetransmissionsInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x002A}, nil)
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

// BasebandReadHoldModeActivityOutput represents the output of the command specified in Section 7.3.33
type BasebandReadHoldModeActivityOutput struct {
	Status uint8
	HoldModeActivity uint8
}

func (o *BasebandReadHoldModeActivityOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.HoldModeActivity = r.GetOne()
	return r.Valid()
}

// BasebandReadHoldModeActivitySync executes the command specified in Section 7.3.33 synchronously
func (c *Commands) BasebandReadHoldModeActivitySync (result *BasebandReadHoldModeActivityOutput) (*BasebandReadHoldModeActivityOutput, error) {
	if result == nil {
		result = &BasebandReadHoldModeActivityOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x002B}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteHoldModeActivityInput represents the input of the command specified in Section 7.3.34
type BasebandWriteHoldModeActivityInput struct {
	HoldModeActivity uint8
}

func (i BasebandWriteHoldModeActivityInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.HoldModeActivity)
	return w.Data()
}

// BasebandWriteHoldModeActivitySync executes the command specified in Section 7.3.34 synchronously
func (c *Commands) BasebandWriteHoldModeActivitySync (params BasebandWriteHoldModeActivityInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x002C}, nil)
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

// BasebandReadTransmitPowerLevelInput represents the input of the command specified in Section 7.3.35
type BasebandReadTransmitPowerLevelInput struct {
	ConnectionHandle uint16
	Type uint8
}

func (i BasebandReadTransmitPowerLevelInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.Type)
	return w.Data()
}

// BasebandReadTransmitPowerLevelOutput represents the output of the command specified in Section 7.3.35
type BasebandReadTransmitPowerLevelOutput struct {
	Status uint8
	ConnectionHandle uint16
	TXPowerLevel uint8
}

func (o *BasebandReadTransmitPowerLevelOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.TXPowerLevel = r.GetOne()
	return r.Valid()
}

// BasebandReadTransmitPowerLevelSync executes the command specified in Section 7.3.35 synchronously
func (c *Commands) BasebandReadTransmitPowerLevelSync (params BasebandReadTransmitPowerLevelInput, result *BasebandReadTransmitPowerLevelOutput) (*BasebandReadTransmitPowerLevelOutput, error) {
	if result == nil {
		result = &BasebandReadTransmitPowerLevelOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x002D}, nil)
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

// BasebandReadSynchronousFlowControlEnableOutput represents the output of the command specified in Section 7.3.36
type BasebandReadSynchronousFlowControlEnableOutput struct {
	Status uint8
	SynchronousFlowControlEnable uint8
}

func (o *BasebandReadSynchronousFlowControlEnableOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.SynchronousFlowControlEnable = r.GetOne()
	return r.Valid()
}

// BasebandReadSynchronousFlowControlEnableSync executes the command specified in Section 7.3.36 synchronously
func (c *Commands) BasebandReadSynchronousFlowControlEnableSync (result *BasebandReadSynchronousFlowControlEnableOutput) (*BasebandReadSynchronousFlowControlEnableOutput, error) {
	if result == nil {
		result = &BasebandReadSynchronousFlowControlEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x002E}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteSynchronousFlowControlEnableInput represents the input of the command specified in Section 7.3.37
type BasebandWriteSynchronousFlowControlEnableInput struct {
	SynchronousFlowControlEnable uint8
}

func (i BasebandWriteSynchronousFlowControlEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.SynchronousFlowControlEnable)
	return w.Data()
}

// BasebandWriteSynchronousFlowControlEnableSync executes the command specified in Section 7.3.37 synchronously
func (c *Commands) BasebandWriteSynchronousFlowControlEnableSync (params BasebandWriteSynchronousFlowControlEnableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x002F}, nil)
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

// BasebandSetControllerToHostFlowControlInput represents the input of the command specified in Section 7.3.38
type BasebandSetControllerToHostFlowControlInput struct {
	FlowControlEnable uint8
}

func (i BasebandSetControllerToHostFlowControlInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.FlowControlEnable)
	return w.Data()
}

// BasebandSetControllerToHostFlowControlSync executes the command specified in Section 7.3.38 synchronously
func (c *Commands) BasebandSetControllerToHostFlowControlSync (params BasebandSetControllerToHostFlowControlInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0031}, nil)
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

// BasebandHostBufferSizeInput represents the input of the command specified in Section 7.3.39
type BasebandHostBufferSizeInput struct {
	HostACLDataPacketLength uint16
	HostSynchronousDataPacketLength uint8
	HostTotalNumACLDataPackets uint16
	HostTotalNumSynchronousDataPackets uint16
}

func (i BasebandHostBufferSizeInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.HostACLDataPacketLength)
	w.PutOne(i.HostSynchronousDataPacketLength)
	binary.LittleEndian.PutUint16(w.Put(2), i.HostTotalNumACLDataPackets)
	binary.LittleEndian.PutUint16(w.Put(2), i.HostTotalNumSynchronousDataPackets)
	return w.Data()
}

// BasebandHostBufferSizeSync executes the command specified in Section 7.3.39 synchronously
func (c *Commands) BasebandHostBufferSizeSync (params BasebandHostBufferSizeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0033}, nil)
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

// BasebandHostNumberOfCompletedPacketsInput represents the input of the command specified in Section 7.3.40
type BasebandHostNumberOfCompletedPacketsInput struct {
	NumHandles uint8
	ConnectionHandle []uint16
	HostNumCompletedPackets []uint16
}

func (i BasebandHostNumberOfCompletedPacketsInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.NumHandles)
	if len(i.ConnectionHandle) != int(i.NumHandles) {
		panic("len(i.ConnectionHandle) != int(i.NumHandles)")
	}
	for _, m := range i.ConnectionHandle {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.HostNumCompletedPackets) != int(i.NumHandles) {
		panic("len(i.HostNumCompletedPackets) != int(i.NumHandles)")
	}
	for _, m := range i.HostNumCompletedPackets {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	return w.Data()
}

// BasebandHostNumberOfCompletedPacketsSync executes the command specified in Section 7.3.40 synchronously
func (c *Commands) BasebandHostNumberOfCompletedPacketsSync (params BasebandHostNumberOfCompletedPacketsInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0035}, nil)
	if err != nil {
		return err
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	_, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		return err
	}


	err2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

	return err
}

// BasebandReadLinkSupervisionTimeoutInput represents the input of the command specified in Section 7.3.41
type BasebandReadLinkSupervisionTimeoutInput struct {
	Handle uint16
}

func (i BasebandReadLinkSupervisionTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.Handle)
	return w.Data()
}

// BasebandReadLinkSupervisionTimeoutOutput represents the output of the command specified in Section 7.3.41
type BasebandReadLinkSupervisionTimeoutOutput struct {
	Status uint8
	Handle uint16
	LinkSupervisionTimeout uint16
}

func (o *BasebandReadLinkSupervisionTimeoutOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	o.LinkSupervisionTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadLinkSupervisionTimeoutSync executes the command specified in Section 7.3.41 synchronously
func (c *Commands) BasebandReadLinkSupervisionTimeoutSync (params BasebandReadLinkSupervisionTimeoutInput, result *BasebandReadLinkSupervisionTimeoutOutput) (*BasebandReadLinkSupervisionTimeoutOutput, error) {
	if result == nil {
		result = &BasebandReadLinkSupervisionTimeoutOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0036}, nil)
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

// BasebandWriteLinkSupervisionTimeoutInput represents the input of the command specified in Section 7.3.42
type BasebandWriteLinkSupervisionTimeoutInput struct {
	Handle uint16
	LinkSupervisionTimeout uint16
}

func (i BasebandWriteLinkSupervisionTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.Handle)
	binary.LittleEndian.PutUint16(w.Put(2), i.LinkSupervisionTimeout)
	return w.Data()
}

// BasebandWriteLinkSupervisionTimeoutOutput represents the output of the command specified in Section 7.3.42
type BasebandWriteLinkSupervisionTimeoutOutput struct {
	Status uint8
	Handle uint16
}

func (o *BasebandWriteLinkSupervisionTimeoutOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandWriteLinkSupervisionTimeoutSync executes the command specified in Section 7.3.42 synchronously
func (c *Commands) BasebandWriteLinkSupervisionTimeoutSync (params BasebandWriteLinkSupervisionTimeoutInput, result *BasebandWriteLinkSupervisionTimeoutOutput) (*BasebandWriteLinkSupervisionTimeoutOutput, error) {
	if result == nil {
		result = &BasebandWriteLinkSupervisionTimeoutOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0037}, nil)
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

// BasebandReadNumberOfSupportedIACOutput represents the output of the command specified in Section 7.3.43
type BasebandReadNumberOfSupportedIACOutput struct {
	Status uint8
	NumSupportedIAC uint8
}

func (o *BasebandReadNumberOfSupportedIACOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.NumSupportedIAC = r.GetOne()
	return r.Valid()
}

// BasebandReadNumberOfSupportedIACSync executes the command specified in Section 7.3.43 synchronously
func (c *Commands) BasebandReadNumberOfSupportedIACSync (result *BasebandReadNumberOfSupportedIACOutput) (*BasebandReadNumberOfSupportedIACOutput, error) {
	if result == nil {
		result = &BasebandReadNumberOfSupportedIACOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0038}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandReadCurrentIACLAPOutput represents the output of the command specified in Section 7.3.44
type BasebandReadCurrentIACLAPOutput struct {
	Status uint8
	NumCurrentIAC uint8
	IACLAP []uint32
}

func (o *BasebandReadCurrentIACLAPOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.NumCurrentIAC = r.GetOne()
	if cap(o.IACLAP) < int(o.NumCurrentIAC) {
		o.IACLAP = make([]uint32, 0, int(o.NumCurrentIAC))
	}
	o.IACLAP = o.IACLAP[:int(o.NumCurrentIAC)]
	for j:=0; j<int(o.NumCurrentIAC); j++ {
		o.IACLAP[j] = decodeUint24(r.Get(3))
	}
	return r.Valid()
}

// BasebandReadCurrentIACLAPSync executes the command specified in Section 7.3.44 synchronously
func (c *Commands) BasebandReadCurrentIACLAPSync (result *BasebandReadCurrentIACLAPOutput) (*BasebandReadCurrentIACLAPOutput, error) {
	if result == nil {
		result = &BasebandReadCurrentIACLAPOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0039}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteCurrentIACLAPInput represents the input of the command specified in Section 7.3.45
type BasebandWriteCurrentIACLAPInput struct {
	NumCurrentIAC uint8
	IACLAP []uint32
}

func (i BasebandWriteCurrentIACLAPInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.NumCurrentIAC)
	if len(i.IACLAP) != int(i.NumCurrentIAC) {
		panic("len(i.IACLAP) != int(i.NumCurrentIAC)")
	}
	for _, m := range i.IACLAP {
		encodeUint24(w.Put(3), m)
	}
	return w.Data()
}

// BasebandWriteCurrentIACLAPSync executes the command specified in Section 7.3.45 synchronously
func (c *Commands) BasebandWriteCurrentIACLAPSync (params BasebandWriteCurrentIACLAPInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x003A}, nil)
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

// BasebandSetAFHHostChannelClassificationInput represents the input of the command specified in Section 7.3.46
type BasebandSetAFHHostChannelClassificationInput struct {
	AFHHostChannelClassification [10]byte
}

func (i BasebandSetAFHHostChannelClassificationInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(10), i.AFHHostChannelClassification[:])
	return w.Data()
}

// BasebandSetAFHHostChannelClassificationSync executes the command specified in Section 7.3.46 synchronously
func (c *Commands) BasebandSetAFHHostChannelClassificationSync (params BasebandSetAFHHostChannelClassificationInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x003F}, nil)
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

// BasebandReadInquiryScanTypeOutput represents the output of the command specified in Section 7.3.47
type BasebandReadInquiryScanTypeOutput struct {
	Status uint8
	InquiryScanType uint8
}

func (o *BasebandReadInquiryScanTypeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.InquiryScanType = r.GetOne()
	return r.Valid()
}

// BasebandReadInquiryScanTypeSync executes the command specified in Section 7.3.47 synchronously
func (c *Commands) BasebandReadInquiryScanTypeSync (result *BasebandReadInquiryScanTypeOutput) (*BasebandReadInquiryScanTypeOutput, error) {
	if result == nil {
		result = &BasebandReadInquiryScanTypeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0042}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteInquiryScanTypeInput represents the input of the command specified in Section 7.3.48
type BasebandWriteInquiryScanTypeInput struct {
	ScanType uint8
}

func (i BasebandWriteInquiryScanTypeInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.ScanType)
	return w.Data()
}

// BasebandWriteInquiryScanTypeSync executes the command specified in Section 7.3.48 synchronously
func (c *Commands) BasebandWriteInquiryScanTypeSync (params BasebandWriteInquiryScanTypeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0043}, nil)
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

// BasebandReadInquiryModeOutput represents the output of the command specified in Section 7.3.49
type BasebandReadInquiryModeOutput struct {
	Status uint8
	InquiryMode uint8
}

func (o *BasebandReadInquiryModeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.InquiryMode = r.GetOne()
	return r.Valid()
}

// BasebandReadInquiryModeSync executes the command specified in Section 7.3.49 synchronously
func (c *Commands) BasebandReadInquiryModeSync (result *BasebandReadInquiryModeOutput) (*BasebandReadInquiryModeOutput, error) {
	if result == nil {
		result = &BasebandReadInquiryModeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0044}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteInquiryModeInput represents the input of the command specified in Section 7.3.50
type BasebandWriteInquiryModeInput struct {
	InquiryMode uint8
}

func (i BasebandWriteInquiryModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.InquiryMode)
	return w.Data()
}

// BasebandWriteInquiryModeSync executes the command specified in Section 7.3.50 synchronously
func (c *Commands) BasebandWriteInquiryModeSync (params BasebandWriteInquiryModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0045}, nil)
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

// BasebandReadPageScanTypeOutput represents the output of the command specified in Section 7.3.51
type BasebandReadPageScanTypeOutput struct {
	Status uint8
	PageScanType uint8
}

func (o *BasebandReadPageScanTypeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.PageScanType = r.GetOne()
	return r.Valid()
}

// BasebandReadPageScanTypeSync executes the command specified in Section 7.3.51 synchronously
func (c *Commands) BasebandReadPageScanTypeSync (result *BasebandReadPageScanTypeOutput) (*BasebandReadPageScanTypeOutput, error) {
	if result == nil {
		result = &BasebandReadPageScanTypeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0046}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWritePageScanTypeInput represents the input of the command specified in Section 7.3.52
type BasebandWritePageScanTypeInput struct {
	PageScanType uint8
}

func (i BasebandWritePageScanTypeInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PageScanType)
	return w.Data()
}

// BasebandWritePageScanTypeSync executes the command specified in Section 7.3.52 synchronously
func (c *Commands) BasebandWritePageScanTypeSync (params BasebandWritePageScanTypeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0047}, nil)
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

// BasebandReadAFHChannelAssessmentModeOutput represents the output of the command specified in Section 7.3.53
type BasebandReadAFHChannelAssessmentModeOutput struct {
	Status uint8
	AFHChannelAssessmentMode uint8
}

func (o *BasebandReadAFHChannelAssessmentModeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.AFHChannelAssessmentMode = r.GetOne()
	return r.Valid()
}

// BasebandReadAFHChannelAssessmentModeSync executes the command specified in Section 7.3.53 synchronously
func (c *Commands) BasebandReadAFHChannelAssessmentModeSync (result *BasebandReadAFHChannelAssessmentModeOutput) (*BasebandReadAFHChannelAssessmentModeOutput, error) {
	if result == nil {
		result = &BasebandReadAFHChannelAssessmentModeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0048}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteAFHChannelAssessmentModeInput represents the input of the command specified in Section 7.3.54
type BasebandWriteAFHChannelAssessmentModeInput struct {
	AFHChannelAssessmentMode uint8
}

func (i BasebandWriteAFHChannelAssessmentModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AFHChannelAssessmentMode)
	return w.Data()
}

// BasebandWriteAFHChannelAssessmentModeSync executes the command specified in Section 7.3.54 synchronously
func (c *Commands) BasebandWriteAFHChannelAssessmentModeSync (params BasebandWriteAFHChannelAssessmentModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0049}, nil)
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

// BasebandReadExtendedInquiryResponseOutput represents the output of the command specified in Section 7.3.55
type BasebandReadExtendedInquiryResponseOutput struct {
	Status uint8
	FECRequired uint8
	ExtendedInquiryResponse [240]byte
}

func (o *BasebandReadExtendedInquiryResponseOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.FECRequired = r.GetOne()
	copy(o.ExtendedInquiryResponse[:], r.Get(240))
	return r.Valid()
}

// BasebandReadExtendedInquiryResponseSync executes the command specified in Section 7.3.55 synchronously
func (c *Commands) BasebandReadExtendedInquiryResponseSync (result *BasebandReadExtendedInquiryResponseOutput) (*BasebandReadExtendedInquiryResponseOutput, error) {
	if result == nil {
		result = &BasebandReadExtendedInquiryResponseOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0051}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteExtendedInquiryResponseInput represents the input of the command specified in Section 7.3.56
type BasebandWriteExtendedInquiryResponseInput struct {
	FECRequired uint8
	ExtendedInquiryResponse [240]byte
}

func (i BasebandWriteExtendedInquiryResponseInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.FECRequired)
	copy(w.Put(240), i.ExtendedInquiryResponse[:])
	return w.Data()
}

// BasebandWriteExtendedInquiryResponseSync executes the command specified in Section 7.3.56 synchronously
func (c *Commands) BasebandWriteExtendedInquiryResponseSync (params BasebandWriteExtendedInquiryResponseInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0052}, nil)
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

// BasebandRefreshEncryptionKeyInput represents the input of the command specified in Section 7.3.57
type BasebandRefreshEncryptionKeyInput struct {
	ConnectionHandle uint16
}

func (i BasebandRefreshEncryptionKeyInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// BasebandRefreshEncryptionKeySync executes the command specified in Section 7.3.57 synchronously
func (c *Commands) BasebandRefreshEncryptionKeySync (params BasebandRefreshEncryptionKeyInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0053}, nil)
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

// BasebandReadSimplePairingModeOutput represents the output of the command specified in Section 7.3.58
type BasebandReadSimplePairingModeOutput struct {
	Status uint8
	SimplePairingMode uint8
}

func (o *BasebandReadSimplePairingModeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.SimplePairingMode = r.GetOne()
	return r.Valid()
}

// BasebandReadSimplePairingModeSync executes the command specified in Section 7.3.58 synchronously
func (c *Commands) BasebandReadSimplePairingModeSync (result *BasebandReadSimplePairingModeOutput) (*BasebandReadSimplePairingModeOutput, error) {
	if result == nil {
		result = &BasebandReadSimplePairingModeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0055}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteSimplePairingModeInput represents the input of the command specified in Section 7.3.59
type BasebandWriteSimplePairingModeInput struct {
	SimplePairingMode uint8
}

func (i BasebandWriteSimplePairingModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.SimplePairingMode)
	return w.Data()
}

// BasebandWriteSimplePairingModeSync executes the command specified in Section 7.3.59 synchronously
func (c *Commands) BasebandWriteSimplePairingModeSync (params BasebandWriteSimplePairingModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0056}, nil)
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

// BasebandReadLocalOOBDataOutput represents the output of the command specified in Section 7.3.60
type BasebandReadLocalOOBDataOutput struct {
	Status uint8
	C [16]byte
	R [16]byte
}

func (o *BasebandReadLocalOOBDataOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.C[:], r.Get(16))
	copy(o.R[:], r.Get(16))
	return r.Valid()
}

// BasebandReadLocalOOBDataSync executes the command specified in Section 7.3.60 synchronously
func (c *Commands) BasebandReadLocalOOBDataSync (result *BasebandReadLocalOOBDataOutput) (*BasebandReadLocalOOBDataOutput, error) {
	if result == nil {
		result = &BasebandReadLocalOOBDataOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0057}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandReadInquiryResponseTransmitPowerLevelOutput represents the output of the command specified in Section 7.3.61
type BasebandReadInquiryResponseTransmitPowerLevelOutput struct {
	Status uint8
	TXPower uint8
}

func (o *BasebandReadInquiryResponseTransmitPowerLevelOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.TXPower = r.GetOne()
	return r.Valid()
}

// BasebandReadInquiryResponseTransmitPowerLevelSync executes the command specified in Section 7.3.61 synchronously
func (c *Commands) BasebandReadInquiryResponseTransmitPowerLevelSync (result *BasebandReadInquiryResponseTransmitPowerLevelOutput) (*BasebandReadInquiryResponseTransmitPowerLevelOutput, error) {
	if result == nil {
		result = &BasebandReadInquiryResponseTransmitPowerLevelOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0058}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteInquiryTransmitPowerLevelInput represents the input of the command specified in Section 7.3.62
type BasebandWriteInquiryTransmitPowerLevelInput struct {
	TXPower uint8
}

func (i BasebandWriteInquiryTransmitPowerLevelInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.TXPower)
	return w.Data()
}

// BasebandWriteInquiryTransmitPowerLevelSync executes the command specified in Section 7.3.62 synchronously
func (c *Commands) BasebandWriteInquiryTransmitPowerLevelSync (params BasebandWriteInquiryTransmitPowerLevelInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0059}, nil)
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

// BasebandSendKeypressNotificationInput represents the input of the command specified in Section 7.3.63
type BasebandSendKeypressNotificationInput struct {
	BDADDR [6]byte
	NotificationType uint8
}

func (i BasebandSendKeypressNotificationInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.NotificationType)
	return w.Data()
}

// BasebandSendKeypressNotificationOutput represents the output of the command specified in Section 7.3.63
type BasebandSendKeypressNotificationOutput struct {
	Status uint8
	BDADDR [6]byte
}

func (o *BasebandSendKeypressNotificationOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.BDADDR[:], r.Get(6))
	return r.Valid()
}

// BasebandSendKeypressNotificationSync executes the command specified in Section 7.3.63 synchronously
func (c *Commands) BasebandSendKeypressNotificationSync (params BasebandSendKeypressNotificationInput, result *BasebandSendKeypressNotificationOutput) (*BasebandSendKeypressNotificationOutput, error) {
	if result == nil {
		result = &BasebandSendKeypressNotificationOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0060}, nil)
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

// BasebandReadDefaultErroneousDataReportingOutput represents the output of the command specified in Section 7.3.64
type BasebandReadDefaultErroneousDataReportingOutput struct {
	Status uint8
	ErroneousDataReporting uint8
}

func (o *BasebandReadDefaultErroneousDataReportingOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ErroneousDataReporting = r.GetOne()
	return r.Valid()
}

// BasebandReadDefaultErroneousDataReportingSync executes the command specified in Section 7.3.64 synchronously
func (c *Commands) BasebandReadDefaultErroneousDataReportingSync (result *BasebandReadDefaultErroneousDataReportingOutput) (*BasebandReadDefaultErroneousDataReportingOutput, error) {
	if result == nil {
		result = &BasebandReadDefaultErroneousDataReportingOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x005A}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteDefaultErroneousDataReportingInput represents the input of the command specified in Section 7.3.65
type BasebandWriteDefaultErroneousDataReportingInput struct {
	ErroneousDataReporting uint8
}

func (i BasebandWriteDefaultErroneousDataReportingInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.ErroneousDataReporting)
	return w.Data()
}

// BasebandWriteDefaultErroneousDataReportingSync executes the command specified in Section 7.3.65 synchronously
func (c *Commands) BasebandWriteDefaultErroneousDataReportingSync (params BasebandWriteDefaultErroneousDataReportingInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x005B}, nil)
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

// BasebandEnhancedFlushInput represents the input of the command specified in Section 7.3.66
type BasebandEnhancedFlushInput struct {
	Handle uint16
	PacketType uint8
}

func (i BasebandEnhancedFlushInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.Handle)
	w.PutOne(i.PacketType)
	return w.Data()
}

// BasebandEnhancedFlushSync executes the command specified in Section 7.3.66 synchronously
func (c *Commands) BasebandEnhancedFlushSync (params BasebandEnhancedFlushInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x005F}, nil)
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

// BasebandReadLogicalLinkAcceptTimeoutOutput represents the output of the command specified in Section 7.3.67
type BasebandReadLogicalLinkAcceptTimeoutOutput struct {
	Status uint8
	LogicalLinkAcceptTimeout uint16
}

func (o *BasebandReadLogicalLinkAcceptTimeoutOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.LogicalLinkAcceptTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadLogicalLinkAcceptTimeoutSync executes the command specified in Section 7.3.67 synchronously
func (c *Commands) BasebandReadLogicalLinkAcceptTimeoutSync (result *BasebandReadLogicalLinkAcceptTimeoutOutput) (*BasebandReadLogicalLinkAcceptTimeoutOutput, error) {
	if result == nil {
		result = &BasebandReadLogicalLinkAcceptTimeoutOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0061}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteLogicalLinkAcceptTimeoutInput represents the input of the command specified in Section 7.3.68
type BasebandWriteLogicalLinkAcceptTimeoutInput struct {
	LogicalLinkAcceptTimeout uint16
}

func (i BasebandWriteLogicalLinkAcceptTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.LogicalLinkAcceptTimeout)
	return w.Data()
}

// BasebandWriteLogicalLinkAcceptTimeoutSync executes the command specified in Section 7.3.68 synchronously
func (c *Commands) BasebandWriteLogicalLinkAcceptTimeoutSync (params BasebandWriteLogicalLinkAcceptTimeoutInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0062}, nil)
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

// BasebandSetEventMaskPage2Input represents the input of the command specified in Section 7.3.69
type BasebandSetEventMaskPage2Input struct {
	EventMaskPage2 uint64
}

func (i BasebandSetEventMaskPage2Input) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint64(w.Put(8), i.EventMaskPage2)
	return w.Data()
}

// BasebandSetEventMaskPage2Sync executes the command specified in Section 7.3.69 synchronously
func (c *Commands) BasebandSetEventMaskPage2Sync (params BasebandSetEventMaskPage2Input) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0063}, nil)
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

// BasebandReadLocationDataSync executes the command specified in Section 7.3.70 synchronously
func (c *Commands) BasebandReadLocationDataSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0064}, nil)
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

// BasebandWriteLocationDataSync executes the command specified in Section 7.3.71 synchronously
func (c *Commands) BasebandWriteLocationDataSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0065}, nil)
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

// BasebandReadFlowControlModeSync executes the command specified in Section 7.3.72 synchronously
func (c *Commands) BasebandReadFlowControlModeSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0066}, nil)
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

// BasebandWriteFlowControlModeInput represents the input of the command specified in Section 7.3.73
type BasebandWriteFlowControlModeInput struct {
	FlowControlMode uint8
}

func (i BasebandWriteFlowControlModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.FlowControlMode)
	return w.Data()
}

// BasebandWriteFlowControlModeSync executes the command specified in Section 7.3.73 synchronously
func (c *Commands) BasebandWriteFlowControlModeSync (params BasebandWriteFlowControlModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0067}, nil)
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

// BasebandReadEnhancedTransmitPowerLevelInput represents the input of the command specified in Section 7.3.74
type BasebandReadEnhancedTransmitPowerLevelInput struct {
	ConnectionHandle uint16
	Type uint8
}

func (i BasebandReadEnhancedTransmitPowerLevelInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.Type)
	return w.Data()
}

// BasebandReadEnhancedTransmitPowerLevelOutput represents the output of the command specified in Section 7.3.74
type BasebandReadEnhancedTransmitPowerLevelOutput struct {
	Status uint8
	ConnectionHandle uint16
	TXPowerLevelGFSK uint8
	TXPowerLevelDQPSK uint8
	TXPowerLevel8DPSK uint8
}

func (o *BasebandReadEnhancedTransmitPowerLevelOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.TXPowerLevelGFSK = r.GetOne()
	o.TXPowerLevelDQPSK = r.GetOne()
	o.TXPowerLevel8DPSK = r.GetOne()
	return r.Valid()
}

// BasebandReadEnhancedTransmitPowerLevelSync executes the command specified in Section 7.3.74 synchronously
func (c *Commands) BasebandReadEnhancedTransmitPowerLevelSync (params BasebandReadEnhancedTransmitPowerLevelInput, result *BasebandReadEnhancedTransmitPowerLevelOutput) (*BasebandReadEnhancedTransmitPowerLevelOutput, error) {
	if result == nil {
		result = &BasebandReadEnhancedTransmitPowerLevelOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0068}, nil)
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

// BasebandReadBestEffortFlushTimeoutInput represents the input of the command specified in Section 7.3.75
type BasebandReadBestEffortFlushTimeoutInput struct {
	LogicalLinkHandle uint16
}

func (i BasebandReadBestEffortFlushTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.LogicalLinkHandle)
	return w.Data()
}

// BasebandReadBestEffortFlushTimeoutOutput represents the output of the command specified in Section 7.3.75
type BasebandReadBestEffortFlushTimeoutOutput struct {
	Status uint8
	BestEffortFlushTimeout uint32
}

func (o *BasebandReadBestEffortFlushTimeoutOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.BestEffortFlushTimeout = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// BasebandReadBestEffortFlushTimeoutSync executes the command specified in Section 7.3.75 synchronously
func (c *Commands) BasebandReadBestEffortFlushTimeoutSync (params BasebandReadBestEffortFlushTimeoutInput, result *BasebandReadBestEffortFlushTimeoutOutput) (*BasebandReadBestEffortFlushTimeoutOutput, error) {
	if result == nil {
		result = &BasebandReadBestEffortFlushTimeoutOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0069}, nil)
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

// BasebandWriteBestEffortFlushTimeoutInput represents the input of the command specified in Section 7.3.76
type BasebandWriteBestEffortFlushTimeoutInput struct {
	LogicalLinkHandle uint16
	BestEffortFlushTimeout uint32
}

func (i BasebandWriteBestEffortFlushTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.LogicalLinkHandle)
	binary.LittleEndian.PutUint32(w.Put(4), i.BestEffortFlushTimeout)
	return w.Data()
}

// BasebandWriteBestEffortFlushTimeoutSync executes the command specified in Section 7.3.76 synchronously
func (c *Commands) BasebandWriteBestEffortFlushTimeoutSync (params BasebandWriteBestEffortFlushTimeoutInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x006A}, nil)
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

// BasebandShortRangeModeInput represents the input of the command specified in Section 7.3.77
type BasebandShortRangeModeInput struct {
	PhysicalLinkHandle uint8
	ShortRangeMode uint8
}

func (i BasebandShortRangeModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PhysicalLinkHandle)
	w.PutOne(i.ShortRangeMode)
	return w.Data()
}

// BasebandShortRangeModeSync executes the command specified in Section 7.3.77 synchronously
func (c *Commands) BasebandShortRangeModeSync (params BasebandShortRangeModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x006B}, nil)
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

// BasebandReadLEHostSupportOutput represents the output of the command specified in Section 7.3.78
type BasebandReadLEHostSupportOutput struct {
	Status uint8
	LESupportedHost uint8
	SimultaneousLEHost uint8
}

func (o *BasebandReadLEHostSupportOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.LESupportedHost = r.GetOne()
	o.SimultaneousLEHost = r.GetOne()
	return r.Valid()
}

// BasebandReadLEHostSupportSync executes the command specified in Section 7.3.78 synchronously
func (c *Commands) BasebandReadLEHostSupportSync (result *BasebandReadLEHostSupportOutput) (*BasebandReadLEHostSupportOutput, error) {
	if result == nil {
		result = &BasebandReadLEHostSupportOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x006C}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteLEHostSupportInput represents the input of the command specified in Section 7.3.79
type BasebandWriteLEHostSupportInput struct {
	LESupportedHost uint8
	SimultaneousLEHost uint8
}

func (i BasebandWriteLEHostSupportInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.LESupportedHost)
	w.PutOne(i.SimultaneousLEHost)
	return w.Data()
}

// BasebandWriteLEHostSupportSync executes the command specified in Section 7.3.79 synchronously
func (c *Commands) BasebandWriteLEHostSupportSync (params BasebandWriteLEHostSupportInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x006}, nil)
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

// BasebandSetMWSChannelParametersInput represents the input of the command specified in Section 7.3.80
type BasebandSetMWSChannelParametersInput struct {
	MWSChannelEnable uint8
	MWSRXCenterFrequency uint16
	MWSTXCenterFrequency uint16
	MWSRXChannelBandwidth uint16
	MWSTXChannelBandwidth uint16
	MWSChannelType uint8
}

func (i BasebandSetMWSChannelParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.MWSChannelEnable)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSRXCenterFrequency)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSTXCenterFrequency)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSRXChannelBandwidth)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSTXChannelBandwidth)
	w.PutOne(i.MWSChannelType)
	return w.Data()
}

// BasebandSetMWSChannelParametersSync executes the command specified in Section 7.3.80 synchronously
func (c *Commands) BasebandSetMWSChannelParametersSync (params BasebandSetMWSChannelParametersInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x006E}, nil)
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

// BasebandSetExternalFrameConfigurationInput represents the input of the command specified in Section 7.3.81
type BasebandSetExternalFrameConfigurationInput struct {
	MWSFrameDuration uint16
	MWSFrameSyncAssertOffset uint16
	MWSFrameSyncAssertJitter uint16
	MWSNumPeriods uint8
	PeriodDuration []uint16
	PeriodType []uint8
}

func (i BasebandSetExternalFrameConfigurationInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSFrameDuration)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSFrameSyncAssertOffset)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSFrameSyncAssertJitter)
	w.PutOne(i.MWSNumPeriods)
	if len(i.PeriodDuration) != int(i.MWSNumPeriods) {
		panic("len(i.PeriodDuration) != int(i.MWSNumPeriods)")
	}
	for _, m := range i.PeriodDuration {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.PeriodType) != int(i.MWSNumPeriods) {
		panic("len(i.PeriodType) != int(i.MWSNumPeriods)")
	}
	for _, m := range i.PeriodType {
		w.PutOne(m)
	}
	return w.Data()
}

// BasebandSetExternalFrameConfigurationSync executes the command specified in Section 7.3.81 synchronously
func (c *Commands) BasebandSetExternalFrameConfigurationSync (params BasebandSetExternalFrameConfigurationInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x006F}, nil)
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

// BasebandSetMWSSignalingInput represents the input of the command specified in Section 7.3.82
type BasebandSetMWSSignalingInput struct {
	MWSRXAssertOffset uint16
	MWSRXAssertJitter uint16
	MWSRXDeassertOffset uint16
	MWSRXDeassertJitter uint16
	MWSTXAssertOffset uint16
	MWSTXAssertJitter uint16
	MWSTXDeassertOffset uint16
	MWSTXDeassertJitter uint16
	MWSPatternAssertOffset uint16
	MWSPatternAssertJitter uint16
	MWSInactivityDurationAssertOffset uint16
	MWSInactivityDurationAssertJitter uint16
	MWSScanFrequencyAssertOffset uint16
	MWSScanFrequencyAssertJitter uint16
	MWSPriorityAssertOffsetRequest uint16
}

func (i BasebandSetMWSSignalingInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSRXAssertOffset)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSRXAssertJitter)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSRXDeassertOffset)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSRXDeassertJitter)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSTXAssertOffset)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSTXAssertJitter)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSTXDeassertOffset)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSTXDeassertJitter)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSPatternAssertOffset)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSPatternAssertJitter)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSInactivityDurationAssertOffset)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSInactivityDurationAssertJitter)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSScanFrequencyAssertOffset)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSScanFrequencyAssertJitter)
	binary.LittleEndian.PutUint16(w.Put(2), i.MWSPriorityAssertOffsetRequest)
	return w.Data()
}

// BasebandSetMWSSignalingOutput represents the output of the command specified in Section 7.3.82
type BasebandSetMWSSignalingOutput struct {
	Status uint8
	BluetoothRXPriorityAssertOffset uint16
	BluetoothRXPriorityAssertJitter uint16
	BluetoothRXPriorityDeassertOffset uint16
	BluetoothRXPriorityDeassertJitter uint16
	I802RXPriorityAssertOffset uint16
	I802RXPriorityAssertJitter uint16
	I802RXPriorityDeassertOffset uint16
	I802RXPriorityDeassertJitter uint16
	BluetoothTXOnAssertOffset uint16
	BluetoothTXOnAssertJitter uint16
	BluetoothTXOnDeassertOffset uint16
	BluetoothTXOnDeassertJitter uint16
	I802TXOnAssertOffset uint16
	I802TXOnAssertJitter uint16
	I802TXOnDeassertOffset uint16
	I802TXOnDeassertJitter uint16
}

func (o *BasebandSetMWSSignalingOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.BluetoothRXPriorityAssertOffset = binary.LittleEndian.Uint16(r.Get(2))
	o.BluetoothRXPriorityAssertJitter = binary.LittleEndian.Uint16(r.Get(2))
	o.BluetoothRXPriorityDeassertOffset = binary.LittleEndian.Uint16(r.Get(2))
	o.BluetoothRXPriorityDeassertJitter = binary.LittleEndian.Uint16(r.Get(2))
	o.I802RXPriorityAssertOffset = binary.LittleEndian.Uint16(r.Get(2))
	o.I802RXPriorityAssertJitter = binary.LittleEndian.Uint16(r.Get(2))
	o.I802RXPriorityDeassertOffset = binary.LittleEndian.Uint16(r.Get(2))
	o.I802RXPriorityDeassertJitter = binary.LittleEndian.Uint16(r.Get(2))
	o.BluetoothTXOnAssertOffset = binary.LittleEndian.Uint16(r.Get(2))
	o.BluetoothTXOnAssertJitter = binary.LittleEndian.Uint16(r.Get(2))
	o.BluetoothTXOnDeassertOffset = binary.LittleEndian.Uint16(r.Get(2))
	o.BluetoothTXOnDeassertJitter = binary.LittleEndian.Uint16(r.Get(2))
	o.I802TXOnAssertOffset = binary.LittleEndian.Uint16(r.Get(2))
	o.I802TXOnAssertJitter = binary.LittleEndian.Uint16(r.Get(2))
	o.I802TXOnDeassertOffset = binary.LittleEndian.Uint16(r.Get(2))
	o.I802TXOnDeassertJitter = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandSetMWSSignalingSync executes the command specified in Section 7.3.82 synchronously
func (c *Commands) BasebandSetMWSSignalingSync (params BasebandSetMWSSignalingInput, result *BasebandSetMWSSignalingOutput) (*BasebandSetMWSSignalingOutput, error) {
	if result == nil {
		result = &BasebandSetMWSSignalingOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0070}, nil)
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

// BasebandSetMWSTransportLayerInput represents the input of the command specified in Section 7.3.83
type BasebandSetMWSTransportLayerInput struct {
	TransportLayer uint8
	ToMWSBaudRate uint32
	FromMWSBaudRate uint32
}

func (i BasebandSetMWSTransportLayerInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.TransportLayer)
	binary.LittleEndian.PutUint32(w.Put(4), i.ToMWSBaudRate)
	binary.LittleEndian.PutUint32(w.Put(4), i.FromMWSBaudRate)
	return w.Data()
}

// BasebandSetMWSTransportLayerSync executes the command specified in Section 7.3.83 synchronously
func (c *Commands) BasebandSetMWSTransportLayerSync (params BasebandSetMWSTransportLayerInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0071}, nil)
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

// BasebandSetMWSScanFrequencyTableInput represents the input of the command specified in Section 7.3.84
type BasebandSetMWSScanFrequencyTableInput struct {
	NumScanFrequencies uint8
	ScanFrequencyLow []uint16
	ScanFrequencyHigh []uint16
}

func (i BasebandSetMWSScanFrequencyTableInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.NumScanFrequencies)
	if len(i.ScanFrequencyLow) != int(i.NumScanFrequencies) {
		panic("len(i.ScanFrequencyLow) != int(i.NumScanFrequencies)")
	}
	for _, m := range i.ScanFrequencyLow {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.ScanFrequencyHigh) != int(i.NumScanFrequencies) {
		panic("len(i.ScanFrequencyHigh) != int(i.NumScanFrequencies)")
	}
	for _, m := range i.ScanFrequencyHigh {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	return w.Data()
}

// BasebandSetMWSScanFrequencyTableSync executes the command specified in Section 7.3.84 synchronously
func (c *Commands) BasebandSetMWSScanFrequencyTableSync (params BasebandSetMWSScanFrequencyTableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0072}, nil)
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

// BasebandSetMWSPATTERNConfigurationInput represents the input of the command specified in Section 7.3.85
type BasebandSetMWSPATTERNConfigurationInput struct {
	MWSPatternIndex uint8
	MWSPatternNumIntervals uint8
	MWSPatternIntervalType []uint8
}

func (i BasebandSetMWSPATTERNConfigurationInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.MWSPatternIndex)
	w.PutOne(i.MWSPatternNumIntervals)
	if len(i.MWSPatternIntervalType) != int(i.MWSPatternNumIntervals) {
		panic("len(i.MWSPatternIntervalType) != int(i.MWSPatternNumIntervals)")
	}
	for _, m := range i.MWSPatternIntervalType {
		w.PutOne(m)
	}
	return w.Data()
}

// BasebandSetMWSPATTERNConfigurationSync executes the command specified in Section 7.3.85 synchronously
func (c *Commands) BasebandSetMWSPATTERNConfigurationSync (params BasebandSetMWSPATTERNConfigurationInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0073}, nil)
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

// BasebandSetReservedLTADDRInput represents the input of the command specified in Section 7.3.86
type BasebandSetReservedLTADDRInput struct {
	LTADDR uint8
}

func (i BasebandSetReservedLTADDRInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.LTADDR)
	return w.Data()
}

// BasebandSetReservedLTADDROutput represents the output of the command specified in Section 7.3.86
type BasebandSetReservedLTADDROutput struct {
	Status uint8
	LTADDR uint8
}

func (o *BasebandSetReservedLTADDROutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.LTADDR = r.GetOne()
	return r.Valid()
}

// BasebandSetReservedLTADDRSync executes the command specified in Section 7.3.86 synchronously
func (c *Commands) BasebandSetReservedLTADDRSync (params BasebandSetReservedLTADDRInput, result *BasebandSetReservedLTADDROutput) (*BasebandSetReservedLTADDROutput, error) {
	if result == nil {
		result = &BasebandSetReservedLTADDROutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0074}, nil)
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

// BasebandDeleteReservedLTADDRInput represents the input of the command specified in Section 7.3.87
type BasebandDeleteReservedLTADDRInput struct {
	LTADDR uint8
}

func (i BasebandDeleteReservedLTADDRInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.LTADDR)
	return w.Data()
}

// BasebandDeleteReservedLTADDROutput represents the output of the command specified in Section 7.3.87
type BasebandDeleteReservedLTADDROutput struct {
	Status uint8
	LTADDR uint8
}

func (o *BasebandDeleteReservedLTADDROutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.LTADDR = r.GetOne()
	return r.Valid()
}

// BasebandDeleteReservedLTADDRSync executes the command specified in Section 7.3.87 synchronously
func (c *Commands) BasebandDeleteReservedLTADDRSync (params BasebandDeleteReservedLTADDRInput, result *BasebandDeleteReservedLTADDROutput) (*BasebandDeleteReservedLTADDROutput, error) {
	if result == nil {
		result = &BasebandDeleteReservedLTADDROutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0075}, nil)
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

// BasebandSetConnectionlessSlaveBroadcastDataInput represents the input of the command specified in Section 7.3.88
type BasebandSetConnectionlessSlaveBroadcastDataInput struct {
	LTADDR uint8
	Fragment uint8
	DataLength uint8
	Data []byte
}

func (i BasebandSetConnectionlessSlaveBroadcastDataInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.LTADDR)
	w.PutOne(i.Fragment)
	w.PutOne(i.DataLength)
	w.PutSlice(i.Data)
	return w.Data()
}

// BasebandSetConnectionlessSlaveBroadcastDataOutput represents the output of the command specified in Section 7.3.88
type BasebandSetConnectionlessSlaveBroadcastDataOutput struct {
	Status uint8
	LTADDR uint8
}

func (o *BasebandSetConnectionlessSlaveBroadcastDataOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.LTADDR = r.GetOne()
	return r.Valid()
}

// BasebandSetConnectionlessSlaveBroadcastDataSync executes the command specified in Section 7.3.88 synchronously
func (c *Commands) BasebandSetConnectionlessSlaveBroadcastDataSync (params BasebandSetConnectionlessSlaveBroadcastDataInput, result *BasebandSetConnectionlessSlaveBroadcastDataOutput) (*BasebandSetConnectionlessSlaveBroadcastDataOutput, error) {
	if result == nil {
		result = &BasebandSetConnectionlessSlaveBroadcastDataOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0076}, nil)
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

// BasebandReadSynchronizationTrainParametersOutput represents the output of the command specified in Section 7.3.89
type BasebandReadSynchronizationTrainParametersOutput struct {
	Status uint8
	SyncTrainInterval uint16
	synchronizationtrainTO uint32
	ServiceData uint8
}

func (o *BasebandReadSynchronizationTrainParametersOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.SyncTrainInterval = binary.LittleEndian.Uint16(r.Get(2))
	o.synchronizationtrainTO = binary.LittleEndian.Uint32(r.Get(4))
	o.ServiceData = r.GetOne()
	return r.Valid()
}

// BasebandReadSynchronizationTrainParametersSync executes the command specified in Section 7.3.89 synchronously
func (c *Commands) BasebandReadSynchronizationTrainParametersSync (result *BasebandReadSynchronizationTrainParametersOutput) (*BasebandReadSynchronizationTrainParametersOutput, error) {
	if result == nil {
		result = &BasebandReadSynchronizationTrainParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0077}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteSynchronizationTrainParametersInput represents the input of the command specified in Section 7.3.90
type BasebandWriteSynchronizationTrainParametersInput struct {
	IntervalMin uint16
	IntervalMax uint16
	synchronizationtrainTO uint32
	ServiceData uint8
}

func (i BasebandWriteSynchronizationTrainParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.IntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.IntervalMax)
	binary.LittleEndian.PutUint32(w.Put(4), i.synchronizationtrainTO)
	w.PutOne(i.ServiceData)
	return w.Data()
}

// BasebandWriteSynchronizationTrainParametersOutput represents the output of the command specified in Section 7.3.90
type BasebandWriteSynchronizationTrainParametersOutput struct {
	Status uint8
	SyncTrainInterval uint16
}

func (o *BasebandWriteSynchronizationTrainParametersOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.SyncTrainInterval = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandWriteSynchronizationTrainParametersSync executes the command specified in Section 7.3.90 synchronously
func (c *Commands) BasebandWriteSynchronizationTrainParametersSync (params BasebandWriteSynchronizationTrainParametersInput, result *BasebandWriteSynchronizationTrainParametersOutput) (*BasebandWriteSynchronizationTrainParametersOutput, error) {
	if result == nil {
		result = &BasebandWriteSynchronizationTrainParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0078}, nil)
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

// BasebandReadSecureConnectionsHostSupportOutput represents the output of the command specified in Section 7.3.91
type BasebandReadSecureConnectionsHostSupportOutput struct {
	Status uint8
	SecureConnectionsHostSupport uint8
}

func (o *BasebandReadSecureConnectionsHostSupportOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.SecureConnectionsHostSupport = r.GetOne()
	return r.Valid()
}

// BasebandReadSecureConnectionsHostSupportSync executes the command specified in Section 7.3.91 synchronously
func (c *Commands) BasebandReadSecureConnectionsHostSupportSync (result *BasebandReadSecureConnectionsHostSupportOutput) (*BasebandReadSecureConnectionsHostSupportOutput, error) {
	if result == nil {
		result = &BasebandReadSecureConnectionsHostSupportOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0079}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteSecureConnectionsHostSupportInput represents the input of the command specified in Section 7.3.92
type BasebandWriteSecureConnectionsHostSupportInput struct {
	SecureConnectionsHostSupport uint8
}

func (i BasebandWriteSecureConnectionsHostSupportInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.SecureConnectionsHostSupport)
	return w.Data()
}

// BasebandWriteSecureConnectionsHostSupportSync executes the command specified in Section 7.3.92 synchronously
func (c *Commands) BasebandWriteSecureConnectionsHostSupportSync (params BasebandWriteSecureConnectionsHostSupportInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x007A}, nil)
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

// BasebandReadAuthenticatedPayloadTimeoutInput represents the input of the command specified in Section 7.3.93
type BasebandReadAuthenticatedPayloadTimeoutInput struct {
	ConnectionHandle uint16
}

func (i BasebandReadAuthenticatedPayloadTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// BasebandReadAuthenticatedPayloadTimeoutOutput represents the output of the command specified in Section 7.3.93
type BasebandReadAuthenticatedPayloadTimeoutOutput struct {
	Status uint8
	ConnectionHandle uint16
	AuthenticatedPayloadTimeout uint16
}

func (o *BasebandReadAuthenticatedPayloadTimeoutOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.AuthenticatedPayloadTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadAuthenticatedPayloadTimeoutSync executes the command specified in Section 7.3.93 synchronously
func (c *Commands) BasebandReadAuthenticatedPayloadTimeoutSync (params BasebandReadAuthenticatedPayloadTimeoutInput, result *BasebandReadAuthenticatedPayloadTimeoutOutput) (*BasebandReadAuthenticatedPayloadTimeoutOutput, error) {
	if result == nil {
		result = &BasebandReadAuthenticatedPayloadTimeoutOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x007B}, nil)
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

// BasebandWriteAuthenticatedPayloadTimeoutInput represents the input of the command specified in Section 7.3.94
type BasebandWriteAuthenticatedPayloadTimeoutInput struct {
	ConnectionHandle uint16
	AuthenticatedPayloadTimeout uint16
}

func (i BasebandWriteAuthenticatedPayloadTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.AuthenticatedPayloadTimeout)
	return w.Data()
}

// BasebandWriteAuthenticatedPayloadTimeoutOutput represents the output of the command specified in Section 7.3.94
type BasebandWriteAuthenticatedPayloadTimeoutOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *BasebandWriteAuthenticatedPayloadTimeoutOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandWriteAuthenticatedPayloadTimeoutSync executes the command specified in Section 7.3.94 synchronously
func (c *Commands) BasebandWriteAuthenticatedPayloadTimeoutSync (params BasebandWriteAuthenticatedPayloadTimeoutInput, result *BasebandWriteAuthenticatedPayloadTimeoutOutput) (*BasebandWriteAuthenticatedPayloadTimeoutOutput, error) {
	if result == nil {
		result = &BasebandWriteAuthenticatedPayloadTimeoutOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x007C}, nil)
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

// BasebandReadLocalOOBExtendedDataOutput represents the output of the command specified in Section 7.3.95
type BasebandReadLocalOOBExtendedDataOutput struct {
	Status uint8
	C192 [16]byte
	R192 [16]byte
	C256 [16]byte
	R256 [16]byte
}

func (o *BasebandReadLocalOOBExtendedDataOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.C192[:], r.Get(16))
	copy(o.R192[:], r.Get(16))
	copy(o.C256[:], r.Get(16))
	copy(o.R256[:], r.Get(16))
	return r.Valid()
}

// BasebandReadLocalOOBExtendedDataSync executes the command specified in Section 7.3.95 synchronously
func (c *Commands) BasebandReadLocalOOBExtendedDataSync (result *BasebandReadLocalOOBExtendedDataOutput) (*BasebandReadLocalOOBExtendedDataOutput, error) {
	if result == nil {
		result = &BasebandReadLocalOOBExtendedDataOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x007D}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandReadExtendedPageTimeoutOutput represents the output of the command specified in Section 7.3.96
type BasebandReadExtendedPageTimeoutOutput struct {
	Status uint8
	ExtendedPageTimeout uint16
}

func (o *BasebandReadExtendedPageTimeoutOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ExtendedPageTimeout = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadExtendedPageTimeoutSync executes the command specified in Section 7.3.96 synchronously
func (c *Commands) BasebandReadExtendedPageTimeoutSync (result *BasebandReadExtendedPageTimeoutOutput) (*BasebandReadExtendedPageTimeoutOutput, error) {
	if result == nil {
		result = &BasebandReadExtendedPageTimeoutOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x007E}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteExtendedPageTimeoutInput represents the input of the command specified in Section 7.3.97
type BasebandWriteExtendedPageTimeoutInput struct {
	ExtendedPageTimeout uint16
}

func (i BasebandWriteExtendedPageTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ExtendedPageTimeout)
	return w.Data()
}

// BasebandWriteExtendedPageTimeoutSync executes the command specified in Section 7.3.97 synchronously
func (c *Commands) BasebandWriteExtendedPageTimeoutSync (params BasebandWriteExtendedPageTimeoutInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x007F}, nil)
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

// BasebandReadExtendedInquiryLengthOutput represents the output of the command specified in Section 7.3.98
type BasebandReadExtendedInquiryLengthOutput struct {
	Status uint8
	ExtendedInquiryLength uint16
}

func (o *BasebandReadExtendedInquiryLengthOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ExtendedInquiryLength = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandReadExtendedInquiryLengthSync executes the command specified in Section 7.3.98 synchronously
func (c *Commands) BasebandReadExtendedInquiryLengthSync (result *BasebandReadExtendedInquiryLengthOutput) (*BasebandReadExtendedInquiryLengthOutput, error) {
	if result == nil {
		result = &BasebandReadExtendedInquiryLengthOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0080}, nil)
	if err != nil {
		return result, err
	}

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

// BasebandWriteExtendedInquiryLengthInput represents the input of the command specified in Section 7.3.99
type BasebandWriteExtendedInquiryLengthInput struct {
	ExtendedInquiryLength uint16
}

func (i BasebandWriteExtendedInquiryLengthInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ExtendedInquiryLength)
	return w.Data()
}

// BasebandWriteExtendedInquiryLengthSync executes the command specified in Section 7.3.99 synchronously
func (c *Commands) BasebandWriteExtendedInquiryLengthSync (params BasebandWriteExtendedInquiryLengthInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0081}, nil)
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

// BasebandSetEcosystemBaseIntervalInput represents the input of the command specified in Section 7.3.100
type BasebandSetEcosystemBaseIntervalInput struct {
	Interval uint16
}

func (i BasebandSetEcosystemBaseIntervalInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.Interval)
	return w.Data()
}

// BasebandSetEcosystemBaseIntervalSync executes the command specified in Section 7.3.100 synchronously
func (c *Commands) BasebandSetEcosystemBaseIntervalSync (params BasebandSetEcosystemBaseIntervalInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0082}, nil)
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

// BasebandConfigureDataPathOutput represents the output of the command specified in Section 7.3.101
type BasebandConfigureDataPathOutput struct {
	Status uint8
	HCIVersion uint8
	HCIRevision uint16
	LMPPALVersion uint8
	ManufacturerName uint16
	LMPPALSubversion uint16
}

func (o *BasebandConfigureDataPathOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.HCIVersion = r.GetOne()
	o.HCIRevision = binary.LittleEndian.Uint16(r.Get(2))
	o.LMPPALVersion = r.GetOne()
	o.ManufacturerName = binary.LittleEndian.Uint16(r.Get(2))
	o.LMPPALSubversion = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// BasebandConfigureDataPathSync executes the command specified in Section 7.3.101 synchronously
func (c *Commands) BasebandConfigureDataPathSync (result *BasebandConfigureDataPathOutput) (*BasebandConfigureDataPathOutput, error) {
	if result == nil {
		result = &BasebandConfigureDataPathOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 3, OCF: 0x0083}, nil)
	if err != nil {
		return result, err
	}

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

