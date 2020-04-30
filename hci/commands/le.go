package hcicommands

import (
	"encoding/binary"
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
)

// LESetEventMaskInput represents the input of the command specified in Section 7.8.1
type LESetEventMaskInput struct {
	LEEventMask uint64
}

func (i LESetEventMaskInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint64(w.Put(8), i.LEEventMask)
	return w.Data()
}

// LESetEventMaskSync executes the command specified in Section 7.8.1 synchronously
func (c *Commands) LESetEventMaskSync (params LESetEventMaskInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0001}, nil)
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

// LEReadBufferSizeOutput represents the output of the command specified in Section 7.8.2
type LEReadBufferSizeOutput struct {
	Status uint8
	LEACLDataPacketLength uint16
	TotalNumLEACLDataPackets uint8
}

func (o *LEReadBufferSizeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.LEACLDataPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.TotalNumLEACLDataPackets = r.GetOne()
	return r.Valid()
}

// LEReadBufferSizeSync executes the command specified in Section 7.8.2 synchronously
func (c *Commands) LEReadBufferSizeSync (result *LEReadBufferSizeOutput) (*LEReadBufferSizeOutput, error) {
	if result == nil {
		result = &LEReadBufferSizeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0002}, nil)
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

// LEReadBufferSizeV2Output represents the output of the command specified in Section 7.8.2
type LEReadBufferSizeV2Output struct {
	Status uint8
	LEACLDataPacketLength uint16
	TotalNumLEACLDataPackets uint8
	ISODataPacketLength uint16
	TotalNumISODataPackets uint8
}

func (o *LEReadBufferSizeV2Output) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.LEACLDataPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.TotalNumLEACLDataPackets = r.GetOne()
	o.ISODataPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.TotalNumISODataPackets = r.GetOne()
	return r.Valid()
}

// LEReadBufferSizeV2Sync executes the command specified in Section 7.8.2 synchronously
func (c *Commands) LEReadBufferSizeV2Sync (result *LEReadBufferSizeV2Output) (*LEReadBufferSizeV2Output, error) {
	if result == nil {
		result = &LEReadBufferSizeV2Output{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0060}, nil)
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

// LEReadLocalSupportedFeaturesOutput represents the output of the command specified in Section 7.8.3
type LEReadLocalSupportedFeaturesOutput struct {
	Status uint8
	LEFeatures uint64
}

func (o *LEReadLocalSupportedFeaturesOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.LEFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// LEReadLocalSupportedFeaturesSync executes the command specified in Section 7.8.3 synchronously
func (c *Commands) LEReadLocalSupportedFeaturesSync (result *LEReadLocalSupportedFeaturesOutput) (*LEReadLocalSupportedFeaturesOutput, error) {
	if result == nil {
		result = &LEReadLocalSupportedFeaturesOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0003}, nil)
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

// LESetRandomAddressInput represents the input of the command specified in Section 7.8.4
type LESetRandomAddressInput struct {
	RandomAddess [6]byte
}

func (i LESetRandomAddressInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.RandomAddess[:])
	return w.Data()
}

// LESetRandomAddressSync executes the command specified in Section 7.8.4 synchronously
func (c *Commands) LESetRandomAddressSync (params LESetRandomAddressInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0005}, nil)
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

// LESetAdvertisingParametersInput represents the input of the command specified in Section 7.8.5
type LESetAdvertisingParametersInput struct {
	AdvertisingIntervalMin uint16
	AdvertisingIntervalMax uint16
	AdvertisingType uint8
	OwnAddressType uint8
	PeerAddressType uint8
	PeerAddress [6]byte
	AdvertisingChannelMap uint8
	AdvertisingFilterPolicy uint8
}

func (i LESetAdvertisingParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.AdvertisingIntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.AdvertisingIntervalMax)
	w.PutOne(i.AdvertisingType)
	w.PutOne(i.OwnAddressType)
	w.PutOne(i.PeerAddressType)
	copy(w.Put(6), i.PeerAddress[:])
	w.PutOne(i.AdvertisingChannelMap)
	w.PutOne(i.AdvertisingFilterPolicy)
	return w.Data()
}

// LESetAdvertisingParametersSync executes the command specified in Section 7.8.5 synchronously
func (c *Commands) LESetAdvertisingParametersSync (params LESetAdvertisingParametersInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0006}, nil)
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

// LEReadAdvertisingPhysicalChannelTxPowerOutput represents the output of the command specified in Section 7.8.6
type LEReadAdvertisingPhysicalChannelTxPowerOutput struct {
	Status uint8
	TXPowerLevel uint8
}

func (o *LEReadAdvertisingPhysicalChannelTxPowerOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.TXPowerLevel = r.GetOne()
	return r.Valid()
}

// LEReadAdvertisingPhysicalChannelTxPowerSync executes the command specified in Section 7.8.6 synchronously
func (c *Commands) LEReadAdvertisingPhysicalChannelTxPowerSync (result *LEReadAdvertisingPhysicalChannelTxPowerOutput) (*LEReadAdvertisingPhysicalChannelTxPowerOutput, error) {
	if result == nil {
		result = &LEReadAdvertisingPhysicalChannelTxPowerOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0007}, nil)
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

// LESetAdvertisingDataInput represents the input of the command specified in Section 7.8.7
type LESetAdvertisingDataInput struct {
	AdvertisingDataLength uint8
	AdvertisingData [31]byte
}

func (i LESetAdvertisingDataInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertisingDataLength)
	copy(w.Put(31), i.AdvertisingData[:])
	return w.Data()
}

// LESetAdvertisingDataSync executes the command specified in Section 7.8.7 synchronously
func (c *Commands) LESetAdvertisingDataSync (params LESetAdvertisingDataInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0008}, nil)
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

// LESetScanResponseDataInput represents the input of the command specified in Section 7.8.8
type LESetScanResponseDataInput struct {
	ScanResponseDataLength uint8
	ScanResponseData [31]byte
}

func (i LESetScanResponseDataInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.ScanResponseDataLength)
	copy(w.Put(31), i.ScanResponseData[:])
	return w.Data()
}

// LESetScanResponseDataSync executes the command specified in Section 7.8.8 synchronously
func (c *Commands) LESetScanResponseDataSync (params LESetScanResponseDataInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0009}, nil)
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

// LESetAdvertisingEnableInput represents the input of the command specified in Section 7.8.9
type LESetAdvertisingEnableInput struct {
	AdvertisingEnable uint8
}

func (i LESetAdvertisingEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertisingEnable)
	return w.Data()
}

// LESetAdvertisingEnableSync executes the command specified in Section 7.8.9 synchronously
func (c *Commands) LESetAdvertisingEnableSync (params LESetAdvertisingEnableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000A}, nil)
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

// LESetScanParametersInput represents the input of the command specified in Section 7.8.10
type LESetScanParametersInput struct {
	LEScanType uint8
	LEScanInterval uint16
	LEScanWindow uint16
	OwnAddressType uint8
	ScanningFilterPolicy uint8
}

func (i LESetScanParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.LEScanType)
	binary.LittleEndian.PutUint16(w.Put(2), i.LEScanInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.LEScanWindow)
	w.PutOne(i.OwnAddressType)
	w.PutOne(i.ScanningFilterPolicy)
	return w.Data()
}

// LESetScanParametersSync executes the command specified in Section 7.8.10 synchronously
func (c *Commands) LESetScanParametersSync (params LESetScanParametersInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000B}, nil)
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

// LESetScanEnableInput represents the input of the command specified in Section 7.8.11
type LESetScanEnableInput struct {
	LEScanEnable uint8
	FilterDuplicates uint8
}

func (i LESetScanEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.LEScanEnable)
	w.PutOne(i.FilterDuplicates)
	return w.Data()
}

// LESetScanEnableSync executes the command specified in Section 7.8.11 synchronously
func (c *Commands) LESetScanEnableSync (params LESetScanEnableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000C}, nil)
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

// LECreateConnectionInput represents the input of the command specified in Section 7.8.12
type LECreateConnectionInput struct {
	LEScanInterval uint16
	LEScanWindow uint16
	InitiatorFilterPolicy uint8
	PeerAddressType uint8
	PeerAddress [6]byte
	OwnAddressType uint8
	ConnectionIntervalMin uint16
	ConnectionIntervalMax uint16
	ConnectionLatency uint16
	SupervisionTimeout uint16
	MinCELength uint16
	MaxCELength uint16
}

func (i LECreateConnectionInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.LEScanInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.LEScanWindow)
	w.PutOne(i.InitiatorFilterPolicy)
	w.PutOne(i.PeerAddressType)
	copy(w.Put(6), i.PeerAddress[:])
	w.PutOne(i.OwnAddressType)
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionIntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionIntervalMax)
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionLatency)
	binary.LittleEndian.PutUint16(w.Put(2), i.SupervisionTimeout)
	binary.LittleEndian.PutUint16(w.Put(2), i.MinCELength)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxCELength)
	return w.Data()
}

// LECreateConnectionSync executes the command specified in Section 7.8.12 synchronously
func (c *Commands) LECreateConnectionSync (params LECreateConnectionInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000D}, nil)
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

// LECreateConnectionCancelSync executes the command specified in Section 7.8.13 synchronously
func (c *Commands) LECreateConnectionCancelSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000E}, nil)
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

// LEReadWhiteListSizeOutput represents the output of the command specified in Section 7.8.14
type LEReadWhiteListSizeOutput struct {
	Status uint8
	WhiteListSize uint8
}

func (o *LEReadWhiteListSizeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.WhiteListSize = r.GetOne()
	return r.Valid()
}

// LEReadWhiteListSizeSync executes the command specified in Section 7.8.14 synchronously
func (c *Commands) LEReadWhiteListSizeSync (result *LEReadWhiteListSizeOutput) (*LEReadWhiteListSizeOutput, error) {
	if result == nil {
		result = &LEReadWhiteListSizeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000F}, nil)
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

// LEClearWhiteListSync executes the command specified in Section 7.8.15 synchronously
func (c *Commands) LEClearWhiteListSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0010}, nil)
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

// LEAddDeviceToWhiteListInput represents the input of the command specified in Section 7.8.16
type LEAddDeviceToWhiteListInput struct {
	AddressType uint8
	Address [6]byte
}

func (i LEAddDeviceToWhiteListInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AddressType)
	copy(w.Put(6), i.Address[:])
	return w.Data()
}

// LEAddDeviceToWhiteListSync executes the command specified in Section 7.8.16 synchronously
func (c *Commands) LEAddDeviceToWhiteListSync (params LEAddDeviceToWhiteListInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0011}, nil)
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

// LERemoveDeviceFromWhiteListInput represents the input of the command specified in Section 7.8.17
type LERemoveDeviceFromWhiteListInput struct {
	AddressType uint8
	Address [6]byte
}

func (i LERemoveDeviceFromWhiteListInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AddressType)
	copy(w.Put(6), i.Address[:])
	return w.Data()
}

// LERemoveDeviceFromWhiteListSync executes the command specified in Section 7.8.17 synchronously
func (c *Commands) LERemoveDeviceFromWhiteListSync (params LERemoveDeviceFromWhiteListInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0012}, nil)
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

// LEConnectionUpdateInput represents the input of the command specified in Section 7.8.18
type LEConnectionUpdateInput struct {
	ConnectionHandle uint16
	ConnectionIntervalMin uint16
	ConnectionIntervalMax uint16
	ConnectionLatency uint16
	SupervisionTimeout uint16
	MinCELength uint16
	MaxCELength uint16
}

func (i LEConnectionUpdateInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionIntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionIntervalMax)
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionLatency)
	binary.LittleEndian.PutUint16(w.Put(2), i.SupervisionTimeout)
	binary.LittleEndian.PutUint16(w.Put(2), i.MinCELength)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxCELength)
	return w.Data()
}

// LEConnectionUpdateSync executes the command specified in Section 7.8.18 synchronously
func (c *Commands) LEConnectionUpdateSync (params LEConnectionUpdateInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0013}, nil)
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

// LESetHostChannelClassificationInput represents the input of the command specified in Section 7.8.19
type LESetHostChannelClassificationInput struct {
	ChannelMap [5]byte
}

func (i LESetHostChannelClassificationInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(5), i.ChannelMap[:])
	return w.Data()
}

// LESetHostChannelClassificationSync executes the command specified in Section 7.8.19 synchronously
func (c *Commands) LESetHostChannelClassificationSync (params LESetHostChannelClassificationInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0014}, nil)
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

// LEReadChannelMapInput represents the input of the command specified in Section 7.8.20
type LEReadChannelMapInput struct {
	ConnectionHandle uint16
}

func (i LEReadChannelMapInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LEReadChannelMapOutput represents the output of the command specified in Section 7.8.20
type LEReadChannelMapOutput struct {
	Status uint8
	ConnectionHandle uint16
	ChannelMap [5]byte
}

func (o *LEReadChannelMapOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	copy(o.ChannelMap[:], r.Get(5))
	return r.Valid()
}

// LEReadChannelMapSync executes the command specified in Section 7.8.20 synchronously
func (c *Commands) LEReadChannelMapSync (params LEReadChannelMapInput, result *LEReadChannelMapOutput) (*LEReadChannelMapOutput, error) {
	if result == nil {
		result = &LEReadChannelMapOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0015}, nil)
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

// LEReadRemoteFeaturesInput represents the input of the command specified in Section 7.8.21
type LEReadRemoteFeaturesInput struct {
	ConnectionHandle uint16
}

func (i LEReadRemoteFeaturesInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LEReadRemoteFeaturesSync executes the command specified in Section 7.8.21 synchronously
func (c *Commands) LEReadRemoteFeaturesSync (params LEReadRemoteFeaturesInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0016}, nil)
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

// LEEncryptInput represents the input of the command specified in Section 7.8.22
type LEEncryptInput struct {
	Key [16]byte
	PlaintextData [16]byte
}

func (i LEEncryptInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(16), i.Key[:])
	copy(w.Put(16), i.PlaintextData[:])
	return w.Data()
}

// LEEncryptOutput represents the output of the command specified in Section 7.8.22
type LEEncryptOutput struct {
	Status uint8
	EncryptedData [16]byte
}

func (o *LEEncryptOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.EncryptedData[:], r.Get(16))
	return r.Valid()
}

// LEEncryptSync executes the command specified in Section 7.8.22 synchronously
func (c *Commands) LEEncryptSync (params LEEncryptInput, result *LEEncryptOutput) (*LEEncryptOutput, error) {
	if result == nil {
		result = &LEEncryptOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0017}, nil)
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

// LERandOutput represents the output of the command specified in Section 7.8.23
type LERandOutput struct {
	Status uint8
	RandomNumber uint64
}

func (o *LERandOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.RandomNumber = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// LERandSync executes the command specified in Section 7.8.23 synchronously
func (c *Commands) LERandSync (result *LERandOutput) (*LERandOutput, error) {
	if result == nil {
		result = &LERandOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0018}, nil)
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

// LEEnableEncryptionInput represents the input of the command specified in Section 7.8.24
type LEEnableEncryptionInput struct {
	ConnectionHandle uint16
	RandomNumber uint64
	EncryptedDiversifier uint16
}

func (i LEEnableEncryptionInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint64(w.Put(8), i.RandomNumber)
	binary.LittleEndian.PutUint16(w.Put(2), i.EncryptedDiversifier)
	return w.Data()
}

// LEEnableEncryptionSync executes the command specified in Section 7.8.24 synchronously
func (c *Commands) LEEnableEncryptionSync (params LEEnableEncryptionInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0019}, nil)
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

// LELongTermKeyRequestReplyInput represents the input of the command specified in Section 7.8.25
type LELongTermKeyRequestReplyInput struct {
	ConnectionHandle uint16
	LongTermKey [16]byte
}

func (i LELongTermKeyRequestReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	copy(w.Put(16), i.LongTermKey[:])
	return w.Data()
}

// LELongTermKeyRequestReplyOutput represents the output of the command specified in Section 7.8.25
type LELongTermKeyRequestReplyOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LELongTermKeyRequestReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LELongTermKeyRequestReplySync executes the command specified in Section 7.8.25 synchronously
func (c *Commands) LELongTermKeyRequestReplySync (params LELongTermKeyRequestReplyInput, result *LELongTermKeyRequestReplyOutput) (*LELongTermKeyRequestReplyOutput, error) {
	if result == nil {
		result = &LELongTermKeyRequestReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001A}, nil)
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

// LELongTermKeyRequestNegativeReplyInput represents the input of the command specified in Section 7.8.26
type LELongTermKeyRequestNegativeReplyInput struct {
	ConnectionHandle uint16
}

func (i LELongTermKeyRequestNegativeReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LELongTermKeyRequestNegativeReplyOutput represents the output of the command specified in Section 7.8.26
type LELongTermKeyRequestNegativeReplyOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LELongTermKeyRequestNegativeReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LELongTermKeyRequestNegativeReplySync executes the command specified in Section 7.8.26 synchronously
func (c *Commands) LELongTermKeyRequestNegativeReplySync (params LELongTermKeyRequestNegativeReplyInput, result *LELongTermKeyRequestNegativeReplyOutput) (*LELongTermKeyRequestNegativeReplyOutput, error) {
	if result == nil {
		result = &LELongTermKeyRequestNegativeReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001B}, nil)
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

// LEReadSupportedStatesOutput represents the output of the command specified in Section 7.8.27
type LEReadSupportedStatesOutput struct {
	Status uint8
	LEStates uint64
}

func (o *LEReadSupportedStatesOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.LEStates = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// LEReadSupportedStatesSync executes the command specified in Section 7.8.27 synchronously
func (c *Commands) LEReadSupportedStatesSync (result *LEReadSupportedStatesOutput) (*LEReadSupportedStatesOutput, error) {
	if result == nil {
		result = &LEReadSupportedStatesOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001C}, nil)
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

// LEReceiverTestInput represents the input of the command specified in Section 7.8.28
type LEReceiverTestInput struct {
	RXChannel uint8
}

func (i LEReceiverTestInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.RXChannel)
	return w.Data()
}

// LEReceiverTestSync executes the command specified in Section 7.8.28 synchronously
func (c *Commands) LEReceiverTestSync (params LEReceiverTestInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001D}, nil)
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

// LEReceiverTestV2Input represents the input of the command specified in Section 7.8.28
type LEReceiverTestV2Input struct {
	RXChannel uint8
	PHY uint8
}

func (i LEReceiverTestV2Input) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.RXChannel)
	w.PutOne(i.PHY)
	return w.Data()
}

// LEReceiverTestV2Sync executes the command specified in Section 7.8.28 synchronously
func (c *Commands) LEReceiverTestV2Sync (params LEReceiverTestV2Input) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0033}, nil)
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

// LEReceiverTestV3Input represents the input of the command specified in Section 7.8.28
type LEReceiverTestV3Input struct {
	RXChannel uint8
	PHY uint8
	ModulationIndex uint8
	ExpectedCTELength uint8
	ExpectedCTEType uint8
	SlotDurations uint8
	SwitchingPatternLength uint8
	AntennaIDs []uint8
}

func (i LEReceiverTestV3Input) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.RXChannel)
	w.PutOne(i.PHY)
	w.PutOne(i.ModulationIndex)
	w.PutOne(i.ExpectedCTELength)
	w.PutOne(i.ExpectedCTEType)
	w.PutOne(i.SlotDurations)
	w.PutOne(i.SwitchingPatternLength)
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(m)
	}
	return w.Data()
}

// LEReceiverTestV3Sync executes the command specified in Section 7.8.28 synchronously
func (c *Commands) LEReceiverTestV3Sync (params LEReceiverTestV3Input) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004F}, nil)
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

// LETransmitterTestInput represents the input of the command specified in Section 7.8.29
type LETransmitterTestInput struct {
	TXChannel uint8
	TestDataLength uint8
	PacketPayload uint8
}

func (i LETransmitterTestInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.TXChannel)
	w.PutOne(i.TestDataLength)
	w.PutOne(i.PacketPayload)
	return w.Data()
}

// LETransmitterTestSync executes the command specified in Section 7.8.29 synchronously
func (c *Commands) LETransmitterTestSync (params LETransmitterTestInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001E}, nil)
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

// LETransmitterTestV2Input represents the input of the command specified in Section 7.8.29
type LETransmitterTestV2Input struct {
	TXChannel uint8
	TestDataLength uint8
	PacketPayload uint8
	PHY uint8
}

func (i LETransmitterTestV2Input) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.TXChannel)
	w.PutOne(i.TestDataLength)
	w.PutOne(i.PacketPayload)
	w.PutOne(i.PHY)
	return w.Data()
}

// LETransmitterTestV2Sync executes the command specified in Section 7.8.29 synchronously
func (c *Commands) LETransmitterTestV2Sync (params LETransmitterTestV2Input) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0034}, nil)
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

// LETransmitterTestV3Input represents the input of the command specified in Section 7.8.29
type LETransmitterTestV3Input struct {
	TXChannel uint8
	TestDataLength uint8
	PacketPayload uint8
	PHY uint8
	CTELength uint8
	CTEType uint8
	SwitchingPatternLength uint8
	AntennaIDs []uint8
}

func (i LETransmitterTestV3Input) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.TXChannel)
	w.PutOne(i.TestDataLength)
	w.PutOne(i.PacketPayload)
	w.PutOne(i.PHY)
	w.PutOne(i.CTELength)
	w.PutOne(i.CTEType)
	w.PutOne(i.SwitchingPatternLength)
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(m)
	}
	return w.Data()
}

// LETransmitterTestV3Sync executes the command specified in Section 7.8.29 synchronously
func (c *Commands) LETransmitterTestV3Sync (params LETransmitterTestV3Input) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0050}, nil)
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

// LETransmitterTestV4Input represents the input of the command specified in Section 7.8.29
type LETransmitterTestV4Input struct {
	TXChannel uint8
	TestDataLength uint8
	PacketPayload uint8
	PHY uint8
	CTELength uint8
	CTEType uint8
	SwitchingPatternLength uint8
	AntennaIDs []uint8
	TransmitPowerLevel uint8
}

func (i LETransmitterTestV4Input) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.TXChannel)
	w.PutOne(i.TestDataLength)
	w.PutOne(i.PacketPayload)
	w.PutOne(i.PHY)
	w.PutOne(i.CTELength)
	w.PutOne(i.CTEType)
	w.PutOne(i.SwitchingPatternLength)
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(m)
	}
	w.PutOne(i.TransmitPowerLevel)
	return w.Data()
}

// LETransmitterTestV4Sync executes the command specified in Section 7.8.29 synchronously
func (c *Commands) LETransmitterTestV4Sync (params LETransmitterTestV4Input) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x007B}, nil)
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

// LETestEndOutput represents the output of the command specified in Section 7.8.30
type LETestEndOutput struct {
	Status uint8
	NumPackets uint16
}

func (o *LETestEndOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.NumPackets = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LETestEndSync executes the command specified in Section 7.8.30 synchronously
func (c *Commands) LETestEndSync (result *LETestEndOutput) (*LETestEndOutput, error) {
	if result == nil {
		result = &LETestEndOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001F}, nil)
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

// LERemoteConnectionParameterRequestReplyInput represents the input of the command specified in Section 7.8.31
type LERemoteConnectionParameterRequestReplyInput struct {
	ConnectionHandle uint16
	IntervalMin uint16
	IntervalMax uint16
	Latency uint16
	Timeout uint16
	MinCELength uint16
	MaxCELength uint16
}

func (i LERemoteConnectionParameterRequestReplyInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.IntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.IntervalMax)
	binary.LittleEndian.PutUint16(w.Put(2), i.Latency)
	binary.LittleEndian.PutUint16(w.Put(2), i.Timeout)
	binary.LittleEndian.PutUint16(w.Put(2), i.MinCELength)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxCELength)
	return w.Data()
}

// LERemoteConnectionParameterRequestReplyOutput represents the output of the command specified in Section 7.8.31
type LERemoteConnectionParameterRequestReplyOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LERemoteConnectionParameterRequestReplyOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LERemoteConnectionParameterRequestReplySync executes the command specified in Section 7.8.31 synchronously
func (c *Commands) LERemoteConnectionParameterRequestReplySync (params LERemoteConnectionParameterRequestReplyInput, result *LERemoteConnectionParameterRequestReplyOutput) (*LERemoteConnectionParameterRequestReplyOutput, error) {
	if result == nil {
		result = &LERemoteConnectionParameterRequestReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0020}, nil)
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

// LESetDataLengthInput represents the input of the command specified in Section 7.8.33
type LESetDataLengthInput struct {
	ConnectionHandle uint16
	TXOctets uint16
	TXTime uint16
}

func (i LESetDataLengthInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.TXOctets)
	binary.LittleEndian.PutUint16(w.Put(2), i.TXTime)
	return w.Data()
}

// LESetDataLengthOutput represents the output of the command specified in Section 7.8.33
type LESetDataLengthOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetDataLengthOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetDataLengthSync executes the command specified in Section 7.8.33 synchronously
func (c *Commands) LESetDataLengthSync (params LESetDataLengthInput, result *LESetDataLengthOutput) (*LESetDataLengthOutput, error) {
	if result == nil {
		result = &LESetDataLengthOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0022}, nil)
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

// LEReadSuggestedDefaultDataLengthOutput represents the output of the command specified in Section 7.8.34
type LEReadSuggestedDefaultDataLengthOutput struct {
	Status uint8
	SuggestedMaxTXOctets uint16
	SuggestedMaxTXTime uint16
}

func (o *LEReadSuggestedDefaultDataLengthOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.SuggestedMaxTXOctets = binary.LittleEndian.Uint16(r.Get(2))
	o.SuggestedMaxTXTime = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEReadSuggestedDefaultDataLengthSync executes the command specified in Section 7.8.34 synchronously
func (c *Commands) LEReadSuggestedDefaultDataLengthSync (result *LEReadSuggestedDefaultDataLengthOutput) (*LEReadSuggestedDefaultDataLengthOutput, error) {
	if result == nil {
		result = &LEReadSuggestedDefaultDataLengthOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0023}, nil)
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

// LEWriteSuggestedDefaultDataLengthInput represents the input of the command specified in Section 7.8.35
type LEWriteSuggestedDefaultDataLengthInput struct {
	SuggestedMaxTXOctets uint16
	SuggestedMaxTXTime uint16
}

func (i LEWriteSuggestedDefaultDataLengthInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.SuggestedMaxTXOctets)
	binary.LittleEndian.PutUint16(w.Put(2), i.SuggestedMaxTXTime)
	return w.Data()
}

// LEWriteSuggestedDefaultDataLengthSync executes the command specified in Section 7.8.35 synchronously
func (c *Commands) LEWriteSuggestedDefaultDataLengthSync (params LEWriteSuggestedDefaultDataLengthInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0024}, nil)
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

// LEReadLocalP256PublicKeyInput represents the input of the command specified in Section 7.8.36
type LEReadLocalP256PublicKeyInput struct {
	RemoteP256PublicKey [64]byte
}

func (i LEReadLocalP256PublicKeyInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(64), i.RemoteP256PublicKey[:])
	return w.Data()
}

// LEReadLocalP256PublicKeySync executes the command specified in Section 7.8.36 synchronously
func (c *Commands) LEReadLocalP256PublicKeySync (params LEReadLocalP256PublicKeyInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0026}, nil)
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

// LEGenerateDHKeyV2Input represents the input of the command specified in Section 7.8.37
type LEGenerateDHKeyV2Input struct {
	RemoteP256PublicKey [64]byte
	KeyType uint8
}

func (i LEGenerateDHKeyV2Input) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(64), i.RemoteP256PublicKey[:])
	w.PutOne(i.KeyType)
	return w.Data()
}

// LEGenerateDHKeyV2Sync executes the command specified in Section 7.8.37 synchronously
func (c *Commands) LEGenerateDHKeyV2Sync (params LEGenerateDHKeyV2Input) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x005E}, nil)
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

// LEAddDeviceToResolvingListInput represents the input of the command specified in Section 7.8.38
type LEAddDeviceToResolvingListInput struct {
	PeerIdentityAddressType uint8
	PeerIdentityAddress [6]byte
	PeerIRK [16]byte
	LocalIRK [16]byte
}

func (i LEAddDeviceToResolvingListInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PeerIdentityAddressType)
	copy(w.Put(6), i.PeerIdentityAddress[:])
	copy(w.Put(16), i.PeerIRK[:])
	copy(w.Put(16), i.LocalIRK[:])
	return w.Data()
}

// LEAddDeviceToResolvingListSync executes the command specified in Section 7.8.38 synchronously
func (c *Commands) LEAddDeviceToResolvingListSync (params LEAddDeviceToResolvingListInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0027}, nil)
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

// LERemoveDeviceFromResolvingListInput represents the input of the command specified in Section 7.8.39
type LERemoveDeviceFromResolvingListInput struct {
	PeerIdentityAddressType uint8
	PeerDeviceAddress [6]byte
}

func (i LERemoveDeviceFromResolvingListInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PeerIdentityAddressType)
	copy(w.Put(6), i.PeerDeviceAddress[:])
	return w.Data()
}

// LERemoveDeviceFromResolvingListSync executes the command specified in Section 7.8.39 synchronously
func (c *Commands) LERemoveDeviceFromResolvingListSync (params LERemoveDeviceFromResolvingListInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0028}, nil)
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

// LEClearResolvingListSync executes the command specified in Section 7.8.40 synchronously
func (c *Commands) LEClearResolvingListSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0029}, nil)
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

// LEReadResolvingListSizeOutput represents the output of the command specified in Section 7.8.41
type LEReadResolvingListSizeOutput struct {
	Status uint8
	ResolvingListSize uint8
}

func (o *LEReadResolvingListSizeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ResolvingListSize = r.GetOne()
	return r.Valid()
}

// LEReadResolvingListSizeSync executes the command specified in Section 7.8.41 synchronously
func (c *Commands) LEReadResolvingListSizeSync (result *LEReadResolvingListSizeOutput) (*LEReadResolvingListSizeOutput, error) {
	if result == nil {
		result = &LEReadResolvingListSizeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002A}, nil)
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

// LEReadPeerResolvableAddressInput represents the input of the command specified in Section 7.8.42
type LEReadPeerResolvableAddressInput struct {
	PeerIdentityAddressType uint8
	PeerIdentityAddress [6]byte
}

func (i LEReadPeerResolvableAddressInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PeerIdentityAddressType)
	copy(w.Put(6), i.PeerIdentityAddress[:])
	return w.Data()
}

// LEReadPeerResolvableAddressOutput represents the output of the command specified in Section 7.8.42
type LEReadPeerResolvableAddressOutput struct {
	Status uint8
	PeerResolvableAddress [6]byte
}

func (o *LEReadPeerResolvableAddressOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.PeerResolvableAddress[:], r.Get(6))
	return r.Valid()
}

// LEReadPeerResolvableAddressSync executes the command specified in Section 7.8.42 synchronously
func (c *Commands) LEReadPeerResolvableAddressSync (params LEReadPeerResolvableAddressInput, result *LEReadPeerResolvableAddressOutput) (*LEReadPeerResolvableAddressOutput, error) {
	if result == nil {
		result = &LEReadPeerResolvableAddressOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002B}, nil)
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

// LEReadLocalResolvableAddressInput represents the input of the command specified in Section 7.8.43
type LEReadLocalResolvableAddressInput struct {
	PeerIdentityAddressType uint8
	PeerIdentityAddress [6]byte
}

func (i LEReadLocalResolvableAddressInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PeerIdentityAddressType)
	copy(w.Put(6), i.PeerIdentityAddress[:])
	return w.Data()
}

// LEReadLocalResolvableAddressOutput represents the output of the command specified in Section 7.8.43
type LEReadLocalResolvableAddressOutput struct {
	Status uint8
	LocalResolvableAddress [6]byte
}

func (o *LEReadLocalResolvableAddressOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	copy(o.LocalResolvableAddress[:], r.Get(6))
	return r.Valid()
}

// LEReadLocalResolvableAddressSync executes the command specified in Section 7.8.43 synchronously
func (c *Commands) LEReadLocalResolvableAddressSync (params LEReadLocalResolvableAddressInput, result *LEReadLocalResolvableAddressOutput) (*LEReadLocalResolvableAddressOutput, error) {
	if result == nil {
		result = &LEReadLocalResolvableAddressOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002C}, nil)
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

// LESetAddressResolutionEnableInput represents the input of the command specified in Section 7.8.44
type LESetAddressResolutionEnableInput struct {
	AddressResolutionEnable uint8
}

func (i LESetAddressResolutionEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AddressResolutionEnable)
	return w.Data()
}

// LESetAddressResolutionEnableSync executes the command specified in Section 7.8.44 synchronously
func (c *Commands) LESetAddressResolutionEnableSync (params LESetAddressResolutionEnableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002D}, nil)
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

// LESetResolvablePrivateAddressTimeoutInput represents the input of the command specified in Section 7.8.45
type LESetResolvablePrivateAddressTimeoutInput struct {
	RPATimeout uint16
}

func (i LESetResolvablePrivateAddressTimeoutInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.RPATimeout)
	return w.Data()
}

// LESetResolvablePrivateAddressTimeoutSync executes the command specified in Section 7.8.45 synchronously
func (c *Commands) LESetResolvablePrivateAddressTimeoutSync (params LESetResolvablePrivateAddressTimeoutInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002E}, nil)
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

// LEReadMaximumDataLengthOutput represents the output of the command specified in Section 7.8.46
type LEReadMaximumDataLengthOutput struct {
	Status uint8
	SupportedMaxTXOctets uint16
	SupportedMaxTXTime uint16
	SupportedMaxRXOctets uint16
	SupportedMaxRXTime uint16
}

func (o *LEReadMaximumDataLengthOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.SupportedMaxTXOctets = binary.LittleEndian.Uint16(r.Get(2))
	o.SupportedMaxTXTime = binary.LittleEndian.Uint16(r.Get(2))
	o.SupportedMaxRXOctets = binary.LittleEndian.Uint16(r.Get(2))
	o.SupportedMaxRXTime = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEReadMaximumDataLengthSync executes the command specified in Section 7.8.46 synchronously
func (c *Commands) LEReadMaximumDataLengthSync (result *LEReadMaximumDataLengthOutput) (*LEReadMaximumDataLengthOutput, error) {
	if result == nil {
		result = &LEReadMaximumDataLengthOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002F}, nil)
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

// LEReadPHYInput represents the input of the command specified in Section 7.8.47
type LEReadPHYInput struct {
	ConnectionHandle uint16
}

func (i LEReadPHYInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LEReadPHYOutput represents the output of the command specified in Section 7.8.47
type LEReadPHYOutput struct {
	Status uint8
	ConnectionHandle uint16
	TXPHY uint8
	RXPHY uint8
}

func (o *LEReadPHYOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.TXPHY = r.GetOne()
	o.RXPHY = r.GetOne()
	return r.Valid()
}

// LEReadPHYSync executes the command specified in Section 7.8.47 synchronously
func (c *Commands) LEReadPHYSync (params LEReadPHYInput, result *LEReadPHYOutput) (*LEReadPHYOutput, error) {
	if result == nil {
		result = &LEReadPHYOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0030}, nil)
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

// LESetDefaultPHYInput represents the input of the command specified in Section 7.8.48
type LESetDefaultPHYInput struct {
	AllPHYs uint8
	TXPHYs uint8
	RXPHYs uint8
}

func (i LESetDefaultPHYInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AllPHYs)
	w.PutOne(i.TXPHYs)
	w.PutOne(i.RXPHYs)
	return w.Data()
}

// LESetDefaultPHYSync executes the command specified in Section 7.8.48 synchronously
func (c *Commands) LESetDefaultPHYSync (params LESetDefaultPHYInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0031}, nil)
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

// LESetPHYInput represents the input of the command specified in Section 7.8.49
type LESetPHYInput struct {
	ConnectionHandle uint16
	AllPHYs uint8
	TXPHYs uint8
	RXPHYs uint8
	PHYOptions uint16
}

func (i LESetPHYInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.AllPHYs)
	w.PutOne(i.TXPHYs)
	w.PutOne(i.RXPHYs)
	binary.LittleEndian.PutUint16(w.Put(2), i.PHYOptions)
	return w.Data()
}

// LESetPHYSync executes the command specified in Section 7.8.49 synchronously
func (c *Commands) LESetPHYSync (params LESetPHYInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0032}, nil)
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

// LESetAdvertisingSetRandomAddressInput represents the input of the command specified in Section 7.8.52
type LESetAdvertisingSetRandomAddressInput struct {
	AdvertisingHandle uint8
	AdvertisingRandomAddress [6]byte
}

func (i LESetAdvertisingSetRandomAddressInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertisingHandle)
	copy(w.Put(6), i.AdvertisingRandomAddress[:])
	return w.Data()
}

// LESetAdvertisingSetRandomAddressSync executes the command specified in Section 7.8.52 synchronously
func (c *Commands) LESetAdvertisingSetRandomAddressSync (params LESetAdvertisingSetRandomAddressInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0035}, nil)
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

// LESetExtendedAdvertisingParametersInput represents the input of the command specified in Section 7.8.53
type LESetExtendedAdvertisingParametersInput struct {
	AdvertisingHandle uint8
	AdvertisingEventProperties uint16
	PrimaryAdvertisingIntervalMin uint32
	PrimaryAdvertisingIntervalMax uint32
	PrimaryAdvertisingChannelMap uint8
	OwnAddressType uint8
	PeerAddressType uint8
	PeerAddress [6]byte
	AdvertisingFilterPolicy uint8
	AdvertisingTXPower uint8
	PrimaryAdvertisingPHY uint8
	SecondaryAdvertisingMaxSkip uint8
	SecondaryAdvertisingPHY uint8
	AdvertisingSID uint8
	ScanRequestNotificationEnable uint8
}

func (i LESetExtendedAdvertisingParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertisingHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.AdvertisingEventProperties)
	encodeUint24(w.Put(3), i.PrimaryAdvertisingIntervalMin)
	encodeUint24(w.Put(3), i.PrimaryAdvertisingIntervalMax)
	w.PutOne(i.PrimaryAdvertisingChannelMap)
	w.PutOne(i.OwnAddressType)
	w.PutOne(i.PeerAddressType)
	copy(w.Put(6), i.PeerAddress[:])
	w.PutOne(i.AdvertisingFilterPolicy)
	w.PutOne(i.AdvertisingTXPower)
	w.PutOne(i.PrimaryAdvertisingPHY)
	w.PutOne(i.SecondaryAdvertisingMaxSkip)
	w.PutOne(i.SecondaryAdvertisingPHY)
	w.PutOne(i.AdvertisingSID)
	w.PutOne(i.ScanRequestNotificationEnable)
	return w.Data()
}

// LESetExtendedAdvertisingParametersOutput represents the output of the command specified in Section 7.8.53
type LESetExtendedAdvertisingParametersOutput struct {
	Status uint8
	SelectedTXPower uint8
}

func (o *LESetExtendedAdvertisingParametersOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.SelectedTXPower = r.GetOne()
	return r.Valid()
}

// LESetExtendedAdvertisingParametersSync executes the command specified in Section 7.8.53 synchronously
func (c *Commands) LESetExtendedAdvertisingParametersSync (params LESetExtendedAdvertisingParametersInput, result *LESetExtendedAdvertisingParametersOutput) (*LESetExtendedAdvertisingParametersOutput, error) {
	if result == nil {
		result = &LESetExtendedAdvertisingParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0036}, nil)
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

// LESetExtendedAdvertisingDataInput represents the input of the command specified in Section 7.8.54
type LESetExtendedAdvertisingDataInput struct {
	AdvertisingHandle uint8
	Operation uint8
	FragmentPreference uint8
	AdvertisingDataLength uint8
	AdvertisingData []byte
}

func (i LESetExtendedAdvertisingDataInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertisingHandle)
	w.PutOne(i.Operation)
	w.PutOne(i.FragmentPreference)
	w.PutOne(i.AdvertisingDataLength)
	w.PutSlice(i.AdvertisingData)
	return w.Data()
}

// LESetExtendedAdvertisingDataSync executes the command specified in Section 7.8.54 synchronously
func (c *Commands) LESetExtendedAdvertisingDataSync (params LESetExtendedAdvertisingDataInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0037}, nil)
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

// LESetExtendedScanResponseDataInput represents the input of the command specified in Section 7.8.55
type LESetExtendedScanResponseDataInput struct {
	AdvertisingHandle uint8
	Operation uint8
	FragmentPreference uint8
	ScanResponseDataLength uint8
	ScanResponseData []byte
}

func (i LESetExtendedScanResponseDataInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertisingHandle)
	w.PutOne(i.Operation)
	w.PutOne(i.FragmentPreference)
	w.PutOne(i.ScanResponseDataLength)
	w.PutSlice(i.ScanResponseData)
	return w.Data()
}

// LESetExtendedScanResponseDataSync executes the command specified in Section 7.8.55 synchronously
func (c *Commands) LESetExtendedScanResponseDataSync (params LESetExtendedScanResponseDataInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0038}, nil)
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

// LESetExtendedAdvertisingEnableInput represents the input of the command specified in Section 7.8.56
type LESetExtendedAdvertisingEnableInput struct {
	Enable uint8
	NumSets uint8
	AdvertisingHandle []uint8
	Duration []uint16
	MaxExtendedAdvertisingEvents []uint8
}

func (i LESetExtendedAdvertisingEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.Enable)
	w.PutOne(i.NumSets)
	if len(i.AdvertisingHandle) != int(i.NumSets) {
		panic("len(i.AdvertisingHandle) != int(i.NumSets)")
	}
	for _, m := range i.AdvertisingHandle {
		w.PutOne(m)
	}
	if len(i.Duration) != int(i.NumSets) {
		panic("len(i.Duration) != int(i.NumSets)")
	}
	for _, m := range i.Duration {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.MaxExtendedAdvertisingEvents) != int(i.NumSets) {
		panic("len(i.MaxExtendedAdvertisingEvents) != int(i.NumSets)")
	}
	for _, m := range i.MaxExtendedAdvertisingEvents {
		w.PutOne(m)
	}
	return w.Data()
}

// LESetExtendedAdvertisingEnableSync executes the command specified in Section 7.8.56 synchronously
func (c *Commands) LESetExtendedAdvertisingEnableSync (params LESetExtendedAdvertisingEnableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0039}, nil)
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

// LEReadMaximumAdvertisingDataLengthOutput represents the output of the command specified in Section 7.8.57
type LEReadMaximumAdvertisingDataLengthOutput struct {
	Status uint8
	MaxAdvertisingDataLength uint16
}

func (o *LEReadMaximumAdvertisingDataLengthOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.MaxAdvertisingDataLength = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEReadMaximumAdvertisingDataLengthSync executes the command specified in Section 7.8.57 synchronously
func (c *Commands) LEReadMaximumAdvertisingDataLengthSync (result *LEReadMaximumAdvertisingDataLengthOutput) (*LEReadMaximumAdvertisingDataLengthOutput, error) {
	if result == nil {
		result = &LEReadMaximumAdvertisingDataLengthOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003A}, nil)
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

// LEReadNumberofSupportedAdvertisingSetsOutput represents the output of the command specified in Section 7.8.58
type LEReadNumberofSupportedAdvertisingSetsOutput struct {
	Status uint8
	NumSupportedAdvertisingSets uint8
}

func (o *LEReadNumberofSupportedAdvertisingSetsOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.NumSupportedAdvertisingSets = r.GetOne()
	return r.Valid()
}

// LEReadNumberofSupportedAdvertisingSetsSync executes the command specified in Section 7.8.58 synchronously
func (c *Commands) LEReadNumberofSupportedAdvertisingSetsSync (result *LEReadNumberofSupportedAdvertisingSetsOutput) (*LEReadNumberofSupportedAdvertisingSetsOutput, error) {
	if result == nil {
		result = &LEReadNumberofSupportedAdvertisingSetsOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003B}, nil)
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

// LERemoveAdvertisingSetInput represents the input of the command specified in Section 7.8.59
type LERemoveAdvertisingSetInput struct {
	AdvertisingHandle uint8
}

func (i LERemoveAdvertisingSetInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertisingHandle)
	return w.Data()
}

// LERemoveAdvertisingSetSync executes the command specified in Section 7.8.59 synchronously
func (c *Commands) LERemoveAdvertisingSetSync (params LERemoveAdvertisingSetInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003C}, nil)
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

// LEClearAdvertisingSetsSync executes the command specified in Section 7.8.60 synchronously
func (c *Commands) LEClearAdvertisingSetsSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003D}, nil)
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

// LESetPeriodicAdvertisingParametersInput represents the input of the command specified in Section 7.8.61
type LESetPeriodicAdvertisingParametersInput struct {
	AdvertisingHandle uint8
	PeriodicAdvertisingIntervalMin uint16
	PeriodicAdvertisingIntervalMax uint16
	PeriodicAdvertisingProperties uint16
}

func (i LESetPeriodicAdvertisingParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertisingHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.PeriodicAdvertisingIntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.PeriodicAdvertisingIntervalMax)
	binary.LittleEndian.PutUint16(w.Put(2), i.PeriodicAdvertisingProperties)
	return w.Data()
}

// LESetPeriodicAdvertisingParametersSync executes the command specified in Section 7.8.61 synchronously
func (c *Commands) LESetPeriodicAdvertisingParametersSync (params LESetPeriodicAdvertisingParametersInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003E}, nil)
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

// LESetPeriodicAdvertisingDataInput represents the input of the command specified in Section 7.8.62
type LESetPeriodicAdvertisingDataInput struct {
	AdvertisingHandle uint8
	Operation uint8
	AdvertisingDataLength uint8
	AdvertisingData []byte
}

func (i LESetPeriodicAdvertisingDataInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertisingHandle)
	w.PutOne(i.Operation)
	w.PutOne(i.AdvertisingDataLength)
	w.PutSlice(i.AdvertisingData)
	return w.Data()
}

// LESetPeriodicAdvertisingDataSync executes the command specified in Section 7.8.62 synchronously
func (c *Commands) LESetPeriodicAdvertisingDataSync (params LESetPeriodicAdvertisingDataInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003F}, nil)
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

// LESetPeriodicAdvertisingEnableInput represents the input of the command specified in Section 7.8.63
type LESetPeriodicAdvertisingEnableInput struct {
	Enable uint8
	AdvertisingHandle uint8
}

func (i LESetPeriodicAdvertisingEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.Enable)
	w.PutOne(i.AdvertisingHandle)
	return w.Data()
}

// LESetPeriodicAdvertisingEnableSync executes the command specified in Section 7.8.63 synchronously
func (c *Commands) LESetPeriodicAdvertisingEnableSync (params LESetPeriodicAdvertisingEnableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0040}, nil)
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

// LESetExtendedScanParametersInput represents the input of the command specified in Section 7.8.64
type LESetExtendedScanParametersInput struct {
	OwnAddressType uint8
	ScanningFilterPolicy uint8
	ScanningPHYs uint8
	ScanType []uint8
	ScanInterval []uint16
	ScanWindow []uint16
}

func (i LESetExtendedScanParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.OwnAddressType)
	w.PutOne(i.ScanningFilterPolicy)
	w.PutOne(i.ScanningPHYs)
	var0 := countSetBits(uint64(i.ScanningPHYs))
	if len(i.ScanType) != var0 {
		panic("len(i.ScanType) != var0")
	}
	for _, m := range i.ScanType {
		w.PutOne(m)
	}
	if len(i.ScanInterval) != var0 {
		panic("len(i.ScanInterval) != var0")
	}
	for _, m := range i.ScanInterval {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.ScanWindow) != var0 {
		panic("len(i.ScanWindow) != var0")
	}
	for _, m := range i.ScanWindow {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	return w.Data()
}

// LESetExtendedScanParametersSync executes the command specified in Section 7.8.64 synchronously
func (c *Commands) LESetExtendedScanParametersSync (params LESetExtendedScanParametersInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0041}, nil)
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

// LESetExtendedScanEnableInput represents the input of the command specified in Section 7.8.65
type LESetExtendedScanEnableInput struct {
	Enable uint8
	FilterDuplicates uint8
	Duration uint16
	Period uint16
}

func (i LESetExtendedScanEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.Enable)
	w.PutOne(i.FilterDuplicates)
	binary.LittleEndian.PutUint16(w.Put(2), i.Duration)
	binary.LittleEndian.PutUint16(w.Put(2), i.Period)
	return w.Data()
}

// LESetExtendedScanEnableSync executes the command specified in Section 7.8.65 synchronously
func (c *Commands) LESetExtendedScanEnableSync (params LESetExtendedScanEnableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0042}, nil)
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

// LEExtendedCreateConnectionInput represents the input of the command specified in Section 7.8.66
type LEExtendedCreateConnectionInput struct {
	InitiatingFilterPolicy uint8
	OwnAddressType uint8
	PeerAddressType uint8
	PeerAddress [6]byte
	InitiatingPHYs uint8
	ScanInterval []uint16
	ScanWindow []uint16
	ConnectionIntervalMin []uint16
	ConnectionIntervalMax []uint16
	ConnectionLatency []uint16
	SupervisionTimeout []uint16
	MinCELength []uint16
	MaxCELength []uint16
}

func (i LEExtendedCreateConnectionInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.InitiatingFilterPolicy)
	w.PutOne(i.OwnAddressType)
	w.PutOne(i.PeerAddressType)
	copy(w.Put(6), i.PeerAddress[:])
	w.PutOne(i.InitiatingPHYs)
	var1 := countSetBits(uint64(i.InitiatingPHYs))
	if len(i.ScanInterval) != var1 {
		panic("len(i.ScanInterval) != var1")
	}
	for _, m := range i.ScanInterval {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.ScanWindow) != var1 {
		panic("len(i.ScanWindow) != var1")
	}
	for _, m := range i.ScanWindow {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.ConnectionIntervalMin) != var1 {
		panic("len(i.ConnectionIntervalMin) != var1")
	}
	for _, m := range i.ConnectionIntervalMin {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.ConnectionIntervalMax) != var1 {
		panic("len(i.ConnectionIntervalMax) != var1")
	}
	for _, m := range i.ConnectionIntervalMax {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.ConnectionLatency) != var1 {
		panic("len(i.ConnectionLatency) != var1")
	}
	for _, m := range i.ConnectionLatency {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.SupervisionTimeout) != var1 {
		panic("len(i.SupervisionTimeout) != var1")
	}
	for _, m := range i.SupervisionTimeout {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.MinCELength) != var1 {
		panic("len(i.MinCELength) != var1")
	}
	for _, m := range i.MinCELength {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.MaxCELength) != var1 {
		panic("len(i.MaxCELength) != var1")
	}
	for _, m := range i.MaxCELength {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	return w.Data()
}

// LEExtendedCreateConnectionSync executes the command specified in Section 7.8.66 synchronously
func (c *Commands) LEExtendedCreateConnectionSync (params LEExtendedCreateConnectionInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0043}, nil)
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

// LEPeriodicAdvertisingCreateSyncInput represents the input of the command specified in Section 7.8.67
type LEPeriodicAdvertisingCreateSyncInput struct {
	Options uint8
	AdvertisingSID uint8
	AdvertiserAddressType uint8
	AdvertiserAddress [6]byte
	Skip uint16
	SyncTimeout uint16
	SyncCTEType uint8
}

func (i LEPeriodicAdvertisingCreateSyncInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.Options)
	w.PutOne(i.AdvertisingSID)
	w.PutOne(i.AdvertiserAddressType)
	copy(w.Put(6), i.AdvertiserAddress[:])
	binary.LittleEndian.PutUint16(w.Put(2), i.Skip)
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncTimeout)
	w.PutOne(i.SyncCTEType)
	return w.Data()
}

// LEPeriodicAdvertisingCreateSyncSync executes the command specified in Section 7.8.67 synchronously
func (c *Commands) LEPeriodicAdvertisingCreateSyncSync (params LEPeriodicAdvertisingCreateSyncInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0044}, nil)
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

// LEPeriodicAdvertisingCreateSyncCancelSync executes the command specified in Section 7.8.68 synchronously
func (c *Commands) LEPeriodicAdvertisingCreateSyncCancelSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0045}, nil)
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

// LEPeriodicAdvertisingTerminateSyncInput represents the input of the command specified in Section 7.8.69
type LEPeriodicAdvertisingTerminateSyncInput struct {
	SyncHandle uint16
}

func (i LEPeriodicAdvertisingTerminateSyncInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncHandle)
	return w.Data()
}

// LEPeriodicAdvertisingTerminateSyncSync executes the command specified in Section 7.8.69 synchronously
func (c *Commands) LEPeriodicAdvertisingTerminateSyncSync (params LEPeriodicAdvertisingTerminateSyncInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0046}, nil)
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

// LEAddDeviceToPeriodicAdvertiserListInput represents the input of the command specified in Section 7.8.70
type LEAddDeviceToPeriodicAdvertiserListInput struct {
	AdvertiserAddressType uint8
	AdvertiserAddress [6]byte
	AdvertisingSID uint8
}

func (i LEAddDeviceToPeriodicAdvertiserListInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertiserAddressType)
	copy(w.Put(6), i.AdvertiserAddress[:])
	w.PutOne(i.AdvertisingSID)
	return w.Data()
}

// LEAddDeviceToPeriodicAdvertiserListSync executes the command specified in Section 7.8.70 synchronously
func (c *Commands) LEAddDeviceToPeriodicAdvertiserListSync (params LEAddDeviceToPeriodicAdvertiserListInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0047}, nil)
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

// LERemoveDeviceFromPeriodicAdvertiserListInput represents the input of the command specified in Section 7.8.71
type LERemoveDeviceFromPeriodicAdvertiserListInput struct {
	AdvertiserAddressType uint8
	AdvertiserAddress [6]byte
	AdvertisingSID uint8
}

func (i LERemoveDeviceFromPeriodicAdvertiserListInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertiserAddressType)
	copy(w.Put(6), i.AdvertiserAddress[:])
	w.PutOne(i.AdvertisingSID)
	return w.Data()
}

// LERemoveDeviceFromPeriodicAdvertiserListSync executes the command specified in Section 7.8.71 synchronously
func (c *Commands) LERemoveDeviceFromPeriodicAdvertiserListSync (params LERemoveDeviceFromPeriodicAdvertiserListInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0048}, nil)
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

// LEClearPeriodicAdvertiserListSync executes the command specified in Section 7.8.72 synchronously
func (c *Commands) LEClearPeriodicAdvertiserListSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0049}, nil)
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

// LEReadPeriodicAdvertiserListSizeOutput represents the output of the command specified in Section 7.8.73
type LEReadPeriodicAdvertiserListSizeOutput struct {
	Status uint8
	PeriodicAdvertiserListSize uint8
}

func (o *LEReadPeriodicAdvertiserListSizeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.PeriodicAdvertiserListSize = r.GetOne()
	return r.Valid()
}

// LEReadPeriodicAdvertiserListSizeSync executes the command specified in Section 7.8.73 synchronously
func (c *Commands) LEReadPeriodicAdvertiserListSizeSync (result *LEReadPeriodicAdvertiserListSizeOutput) (*LEReadPeriodicAdvertiserListSizeOutput, error) {
	if result == nil {
		result = &LEReadPeriodicAdvertiserListSizeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004A}, nil)
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

// LEReadTransmitPowerOutput represents the output of the command specified in Section 7.8.74
type LEReadTransmitPowerOutput struct {
	Status uint8
	MinTXPower uint8
	MaxTXPower uint8
}

func (o *LEReadTransmitPowerOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.MinTXPower = r.GetOne()
	o.MaxTXPower = r.GetOne()
	return r.Valid()
}

// LEReadTransmitPowerSync executes the command specified in Section 7.8.74 synchronously
func (c *Commands) LEReadTransmitPowerSync (result *LEReadTransmitPowerOutput) (*LEReadTransmitPowerOutput, error) {
	if result == nil {
		result = &LEReadTransmitPowerOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004B}, nil)
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

// LEReadRFPathCompensationOutput represents the output of the command specified in Section 7.8.75
type LEReadRFPathCompensationOutput struct {
	Status uint8
	RFTXPathCompensationValue uint16
	RFRXPathCompensationValue uint16
}

func (o *LEReadRFPathCompensationOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.RFTXPathCompensationValue = binary.LittleEndian.Uint16(r.Get(2))
	o.RFRXPathCompensationValue = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEReadRFPathCompensationSync executes the command specified in Section 7.8.75 synchronously
func (c *Commands) LEReadRFPathCompensationSync (result *LEReadRFPathCompensationOutput) (*LEReadRFPathCompensationOutput, error) {
	if result == nil {
		result = &LEReadRFPathCompensationOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004C}, nil)
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

// LEWriteRFPathCompensationInput represents the input of the command specified in Section 7.8.76
type LEWriteRFPathCompensationInput struct {
	RFTXPathCompensationValue uint16
	RFRXPathCompensationValue uint16
}

func (i LEWriteRFPathCompensationInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.RFTXPathCompensationValue)
	binary.LittleEndian.PutUint16(w.Put(2), i.RFRXPathCompensationValue)
	return w.Data()
}

// LEWriteRFPathCompensationSync executes the command specified in Section 7.8.76 synchronously
func (c *Commands) LEWriteRFPathCompensationSync (params LEWriteRFPathCompensationInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004D}, nil)
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

// LESetPrivacyModeInput represents the input of the command specified in Section 7.8.77
type LESetPrivacyModeInput struct {
	PeerIdentityAddressType uint8
	PeerIdentityAddress [6]byte
	PrivacyMode uint8
}

func (i LESetPrivacyModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PeerIdentityAddressType)
	copy(w.Put(6), i.PeerIdentityAddress[:])
	w.PutOne(i.PrivacyMode)
	return w.Data()
}

// LESetPrivacyModeSync executes the command specified in Section 7.8.77 synchronously
func (c *Commands) LESetPrivacyModeSync (params LESetPrivacyModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004E}, nil)
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

// LESetConnectionlessCTETransmitParametersInput represents the input of the command specified in Section 7.8.80
type LESetConnectionlessCTETransmitParametersInput struct {
	AdvertisingHandle uint8
	CTELength uint8
	CTEType uint8
	CTECount uint8
	SwitchingPatternLength uint8
	AntennaIDs []uint8
}

func (i LESetConnectionlessCTETransmitParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertisingHandle)
	w.PutOne(i.CTELength)
	w.PutOne(i.CTEType)
	w.PutOne(i.CTECount)
	w.PutOne(i.SwitchingPatternLength)
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(m)
	}
	return w.Data()
}

// LESetConnectionlessCTETransmitParametersSync executes the command specified in Section 7.8.80 synchronously
func (c *Commands) LESetConnectionlessCTETransmitParametersSync (params LESetConnectionlessCTETransmitParametersInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0051}, nil)
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

// LESetConnectionlessCTETransmitEnableInput represents the input of the command specified in Section 7.8.81
type LESetConnectionlessCTETransmitEnableInput struct {
	AdvertisingHandle uint8
	CTEEnable uint8
}

func (i LESetConnectionlessCTETransmitEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.AdvertisingHandle)
	w.PutOne(i.CTEEnable)
	return w.Data()
}

// LESetConnectionlessCTETransmitEnableSync executes the command specified in Section 7.8.81 synchronously
func (c *Commands) LESetConnectionlessCTETransmitEnableSync (params LESetConnectionlessCTETransmitEnableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0052}, nil)
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

// LESetConnectionlessIQSamplingEnableInput represents the input of the command specified in Section 7.8.82
type LESetConnectionlessIQSamplingEnableInput struct {
	SyncHandle uint16
	SamplingEnable uint8
	SlotDurations uint8
	MaxSampledCTEs uint8
	SwitchingPatternLength uint8
	AntennaIDs []uint8
}

func (i LESetConnectionlessIQSamplingEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncHandle)
	w.PutOne(i.SamplingEnable)
	w.PutOne(i.SlotDurations)
	w.PutOne(i.MaxSampledCTEs)
	w.PutOne(i.SwitchingPatternLength)
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(m)
	}
	return w.Data()
}

// LESetConnectionlessIQSamplingEnableOutput represents the output of the command specified in Section 7.8.82
type LESetConnectionlessIQSamplingEnableOutput struct {
	Status uint8
	SyncHandle uint16
}

func (o *LESetConnectionlessIQSamplingEnableOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetConnectionlessIQSamplingEnableSync executes the command specified in Section 7.8.82 synchronously
func (c *Commands) LESetConnectionlessIQSamplingEnableSync (params LESetConnectionlessIQSamplingEnableInput, result *LESetConnectionlessIQSamplingEnableOutput) (*LESetConnectionlessIQSamplingEnableOutput, error) {
	if result == nil {
		result = &LESetConnectionlessIQSamplingEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0053}, nil)
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

// LESetConnectionCTEReceiveParametersInput represents the input of the command specified in Section 7.8.83
type LESetConnectionCTEReceiveParametersInput struct {
	ConnectionHandle uint16
	SamplingEnable uint8
	SlotDurations uint8
	SwitchingPatternLength uint8
	AntennaIDs []uint8
}

func (i LESetConnectionCTEReceiveParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.SamplingEnable)
	w.PutOne(i.SlotDurations)
	w.PutOne(i.SwitchingPatternLength)
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(m)
	}
	return w.Data()
}

// LESetConnectionCTEReceiveParametersOutput represents the output of the command specified in Section 7.8.83
type LESetConnectionCTEReceiveParametersOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetConnectionCTEReceiveParametersOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetConnectionCTEReceiveParametersSync executes the command specified in Section 7.8.83 synchronously
func (c *Commands) LESetConnectionCTEReceiveParametersSync (params LESetConnectionCTEReceiveParametersInput, result *LESetConnectionCTEReceiveParametersOutput) (*LESetConnectionCTEReceiveParametersOutput, error) {
	if result == nil {
		result = &LESetConnectionCTEReceiveParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0054}, nil)
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

// LESetConnectionCTETransmitParametersInput represents the input of the command specified in Section 7.8.84
type LESetConnectionCTETransmitParametersInput struct {
	ConnectionHandle uint16
	CTETypes uint8
	SwitchingPatternLength uint8
	AntennaIDs []uint8
}

func (i LESetConnectionCTETransmitParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.CTETypes)
	w.PutOne(i.SwitchingPatternLength)
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(m)
	}
	return w.Data()
}

// LESetConnectionCTETransmitParametersOutput represents the output of the command specified in Section 7.8.84
type LESetConnectionCTETransmitParametersOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetConnectionCTETransmitParametersOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetConnectionCTETransmitParametersSync executes the command specified in Section 7.8.84 synchronously
func (c *Commands) LESetConnectionCTETransmitParametersSync (params LESetConnectionCTETransmitParametersInput, result *LESetConnectionCTETransmitParametersOutput) (*LESetConnectionCTETransmitParametersOutput, error) {
	if result == nil {
		result = &LESetConnectionCTETransmitParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0055}, nil)
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

// LEConnectionCTERequestEnableInput represents the input of the command specified in Section 7.8.85
type LEConnectionCTERequestEnableInput struct {
	ConnectionHandle uint16
	Enable uint8
	CTERequestInterval uint16
	RequestedCTELength uint8
	RequestedCTEType uint8
}

func (i LEConnectionCTERequestEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.Enable)
	binary.LittleEndian.PutUint16(w.Put(2), i.CTERequestInterval)
	w.PutOne(i.RequestedCTELength)
	w.PutOne(i.RequestedCTEType)
	return w.Data()
}

// LEConnectionCTERequestEnableOutput represents the output of the command specified in Section 7.8.85
type LEConnectionCTERequestEnableOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEConnectionCTERequestEnableOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEConnectionCTERequestEnableSync executes the command specified in Section 7.8.85 synchronously
func (c *Commands) LEConnectionCTERequestEnableSync (params LEConnectionCTERequestEnableInput, result *LEConnectionCTERequestEnableOutput) (*LEConnectionCTERequestEnableOutput, error) {
	if result == nil {
		result = &LEConnectionCTERequestEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0056}, nil)
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

// LEConnectionCTEResponseEnableInput represents the input of the command specified in Section 7.8.86
type LEConnectionCTEResponseEnableInput struct {
	ConnectionHandle uint16
	Enable uint8
}

func (i LEConnectionCTEResponseEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.Enable)
	return w.Data()
}

// LEConnectionCTEResponseEnableOutput represents the output of the command specified in Section 7.8.86
type LEConnectionCTEResponseEnableOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEConnectionCTEResponseEnableOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEConnectionCTEResponseEnableSync executes the command specified in Section 7.8.86 synchronously
func (c *Commands) LEConnectionCTEResponseEnableSync (params LEConnectionCTEResponseEnableInput, result *LEConnectionCTEResponseEnableOutput) (*LEConnectionCTEResponseEnableOutput, error) {
	if result == nil {
		result = &LEConnectionCTEResponseEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0057}, nil)
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

// LEReadAntennaInformationOutput represents the output of the command specified in Section 7.8.87
type LEReadAntennaInformationOutput struct {
	Status uint8
	SupportedSwitchingSamplingRates uint8
	NumAntennae uint8
	MaxSwitchingPatternLength uint8
	MaxCTELength uint8
}

func (o *LEReadAntennaInformationOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.SupportedSwitchingSamplingRates = r.GetOne()
	o.NumAntennae = r.GetOne()
	o.MaxSwitchingPatternLength = r.GetOne()
	o.MaxCTELength = r.GetOne()
	return r.Valid()
}

// LEReadAntennaInformationSync executes the command specified in Section 7.8.87 synchronously
func (c *Commands) LEReadAntennaInformationSync (result *LEReadAntennaInformationOutput) (*LEReadAntennaInformationOutput, error) {
	if result == nil {
		result = &LEReadAntennaInformationOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0058}, nil)
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

// LESetPeriodicAdvertisingReceiveEnableInput represents the input of the command specified in Section 7.8.88
type LESetPeriodicAdvertisingReceiveEnableInput struct {
	SyncHandle uint16
	Enable uint8
}

func (i LESetPeriodicAdvertisingReceiveEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncHandle)
	w.PutOne(i.Enable)
	return w.Data()
}

// LESetPeriodicAdvertisingReceiveEnableSync executes the command specified in Section 7.8.88 synchronously
func (c *Commands) LESetPeriodicAdvertisingReceiveEnableSync (params LESetPeriodicAdvertisingReceiveEnableInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0059}, nil)
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

// LEPeriodicAdvertisingSyncTransferInput represents the input of the command specified in Section 7.8.89
type LEPeriodicAdvertisingSyncTransferInput struct {
	ConnectionHandle uint16
	ServiceData uint16
	SyncHandle uint16
}

func (i LEPeriodicAdvertisingSyncTransferInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.ServiceData)
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncHandle)
	return w.Data()
}

// LEPeriodicAdvertisingSyncTransferOutput represents the output of the command specified in Section 7.8.89
type LEPeriodicAdvertisingSyncTransferOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEPeriodicAdvertisingSyncTransferOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEPeriodicAdvertisingSyncTransferSync executes the command specified in Section 7.8.89 synchronously
func (c *Commands) LEPeriodicAdvertisingSyncTransferSync (params LEPeriodicAdvertisingSyncTransferInput, result *LEPeriodicAdvertisingSyncTransferOutput) (*LEPeriodicAdvertisingSyncTransferOutput, error) {
	if result == nil {
		result = &LEPeriodicAdvertisingSyncTransferOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x005A}, nil)
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

// LEPeriodicAdvertisingSetInfoTransferInput represents the input of the command specified in Section 7.8.90
type LEPeriodicAdvertisingSetInfoTransferInput struct {
	ConnectionHandle uint16
	ServiceData uint16
	AdvertisingHandle uint8
}

func (i LEPeriodicAdvertisingSetInfoTransferInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.ServiceData)
	w.PutOne(i.AdvertisingHandle)
	return w.Data()
}

// LEPeriodicAdvertisingSetInfoTransferOutput represents the output of the command specified in Section 7.8.90
type LEPeriodicAdvertisingSetInfoTransferOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEPeriodicAdvertisingSetInfoTransferOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEPeriodicAdvertisingSetInfoTransferSync executes the command specified in Section 7.8.90 synchronously
func (c *Commands) LEPeriodicAdvertisingSetInfoTransferSync (params LEPeriodicAdvertisingSetInfoTransferInput, result *LEPeriodicAdvertisingSetInfoTransferOutput) (*LEPeriodicAdvertisingSetInfoTransferOutput, error) {
	if result == nil {
		result = &LEPeriodicAdvertisingSetInfoTransferOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x005B}, nil)
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

// LESetPeriodicAdvertisingSyncTransferParametersInput represents the input of the command specified in Section 7.8.91
type LESetPeriodicAdvertisingSyncTransferParametersInput struct {
	ConnectionHandle uint16
	Mode uint8
	Skip uint16
	SyncTimeout uint16
	CTEType uint8
}

func (i LESetPeriodicAdvertisingSyncTransferParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.Mode)
	binary.LittleEndian.PutUint16(w.Put(2), i.Skip)
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncTimeout)
	w.PutOne(i.CTEType)
	return w.Data()
}

// LESetPeriodicAdvertisingSyncTransferParametersOutput represents the output of the command specified in Section 7.8.91
type LESetPeriodicAdvertisingSyncTransferParametersOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetPeriodicAdvertisingSyncTransferParametersOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetPeriodicAdvertisingSyncTransferParametersSync executes the command specified in Section 7.8.91 synchronously
func (c *Commands) LESetPeriodicAdvertisingSyncTransferParametersSync (params LESetPeriodicAdvertisingSyncTransferParametersInput, result *LESetPeriodicAdvertisingSyncTransferParametersOutput) (*LESetPeriodicAdvertisingSyncTransferParametersOutput, error) {
	if result == nil {
		result = &LESetPeriodicAdvertisingSyncTransferParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x005C}, nil)
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

// LEModifySleepClockAccuracyInput represents the input of the command specified in Section 7.8.94
type LEModifySleepClockAccuracyInput struct {
	Action uint8
}

func (i LEModifySleepClockAccuracyInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.Action)
	return w.Data()
}

// LEModifySleepClockAccuracySync executes the command specified in Section 7.8.94 synchronously
func (c *Commands) LEModifySleepClockAccuracySync (params LEModifySleepClockAccuracyInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x005F}, nil)
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

// LEReadISOTXSyncInput represents the input of the command specified in Section 7.8.96
type LEReadISOTXSyncInput struct {
	ConnectionHandle uint16
}

func (i LEReadISOTXSyncInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LEReadISOTXSyncOutput represents the output of the command specified in Section 7.8.96
type LEReadISOTXSyncOutput struct {
	Status uint8
	ConnectionHandle uint16
	PacketSequenceNumber uint16
	TimeStamp uint32
	TimeOffset uint32
}

func (o *LEReadISOTXSyncOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PacketSequenceNumber = binary.LittleEndian.Uint16(r.Get(2))
	o.TimeStamp = binary.LittleEndian.Uint32(r.Get(4))
	o.TimeOffset = decodeUint24(r.Get(3))
	return r.Valid()
}

// LEReadISOTXSyncSync executes the command specified in Section 7.8.96 synchronously
func (c *Commands) LEReadISOTXSyncSync (params LEReadISOTXSyncInput, result *LEReadISOTXSyncOutput) (*LEReadISOTXSyncOutput, error) {
	if result == nil {
		result = &LEReadISOTXSyncOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0061}, nil)
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

// LESetCIGParametersInput represents the input of the command specified in Section 7.8.97
type LESetCIGParametersInput struct {
	CIGID uint8
	SDUIntervalMToS uint32
	SDUIntervalSToM uint32
	SlavesClockAccuracy uint8
	Packing uint8
	Framing uint8
	MaxTransportLatencyMToS uint16
	MaxTransportLatencySToM uint16
	CISCount uint8
	CISID []uint8
	MaxSDUMToS []uint16
	MaxSDUSToM []uint16
	PHYMToS []uint8
	PHYSToM []uint8
	RTNMToS []uint8
	RTNSToM []uint8
}

func (i LESetCIGParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.CIGID)
	encodeUint24(w.Put(3), i.SDUIntervalMToS)
	encodeUint24(w.Put(3), i.SDUIntervalSToM)
	w.PutOne(i.SlavesClockAccuracy)
	w.PutOne(i.Packing)
	w.PutOne(i.Framing)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxTransportLatencyMToS)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxTransportLatencySToM)
	w.PutOne(i.CISCount)
	if len(i.CISID) != int(i.CISCount) {
		panic("len(i.CISID) != int(i.CISCount)")
	}
	for _, m := range i.CISID {
		w.PutOne(m)
	}
	if len(i.MaxSDUMToS) != int(i.CISCount) {
		panic("len(i.MaxSDUMToS) != int(i.CISCount)")
	}
	for _, m := range i.MaxSDUMToS {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.MaxSDUSToM) != int(i.CISCount) {
		panic("len(i.MaxSDUSToM) != int(i.CISCount)")
	}
	for _, m := range i.MaxSDUSToM {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.PHYMToS) != int(i.CISCount) {
		panic("len(i.PHYMToS) != int(i.CISCount)")
	}
	for _, m := range i.PHYMToS {
		w.PutOne(m)
	}
	if len(i.PHYSToM) != int(i.CISCount) {
		panic("len(i.PHYSToM) != int(i.CISCount)")
	}
	for _, m := range i.PHYSToM {
		w.PutOne(m)
	}
	if len(i.RTNMToS) != int(i.CISCount) {
		panic("len(i.RTNMToS) != int(i.CISCount)")
	}
	for _, m := range i.RTNMToS {
		w.PutOne(m)
	}
	if len(i.RTNSToM) != int(i.CISCount) {
		panic("len(i.RTNSToM) != int(i.CISCount)")
	}
	for _, m := range i.RTNSToM {
		w.PutOne(m)
	}
	return w.Data()
}

// LESetCIGParametersOutput represents the output of the command specified in Section 7.8.97
type LESetCIGParametersOutput struct {
	Status uint8
	CIGID uint8
	CISCount uint8
	ConnectionHandle []uint16
}

func (o *LESetCIGParametersOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.CIGID = r.GetOne()
	o.CISCount = r.GetOne()
	if cap(o.ConnectionHandle) < int(o.CISCount) {
		o.ConnectionHandle = make([]uint16, 0, int(o.CISCount))
	}
	o.ConnectionHandle = o.ConnectionHandle[:int(o.CISCount)]
	for j:=0; j<int(o.CISCount); j++ {
		o.ConnectionHandle[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// LESetCIGParametersSync executes the command specified in Section 7.8.97 synchronously
func (c *Commands) LESetCIGParametersSync (params LESetCIGParametersInput, result *LESetCIGParametersOutput) (*LESetCIGParametersOutput, error) {
	if result == nil {
		result = &LESetCIGParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0062}, nil)
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

// LESetCIGParametersTestInput represents the input of the command specified in Section 7.8.98
type LESetCIGParametersTestInput struct {
	CIGID uint8
	SDUIntervalMToS uint32
	SDUIntervalSToM uint32
	FTMToS uint8
	FTSToM uint8
	ISOInterval uint16
	SlavesClockAccuracy uint8
	Packing uint8
	Framing uint8
	CISCount uint8
	CISID []uint8
	NSE []uint8
	MaxSDUMToS []uint16
	MaxSDUSToM []uint16
	MaxPDUMToS []uint16
	MaxPDUSToM []uint16
	PHYMToS []uint8
	PHYSToM []uint8
	BNMToS []uint8
	BNSToM []uint8
}

func (i LESetCIGParametersTestInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.CIGID)
	encodeUint24(w.Put(3), i.SDUIntervalMToS)
	encodeUint24(w.Put(3), i.SDUIntervalSToM)
	w.PutOne(i.FTMToS)
	w.PutOne(i.FTSToM)
	binary.LittleEndian.PutUint16(w.Put(2), i.ISOInterval)
	w.PutOne(i.SlavesClockAccuracy)
	w.PutOne(i.Packing)
	w.PutOne(i.Framing)
	w.PutOne(i.CISCount)
	if len(i.CISID) != int(i.CISCount) {
		panic("len(i.CISID) != int(i.CISCount)")
	}
	for _, m := range i.CISID {
		w.PutOne(m)
	}
	if len(i.NSE) != int(i.CISCount) {
		panic("len(i.NSE) != int(i.CISCount)")
	}
	for _, m := range i.NSE {
		w.PutOne(m)
	}
	if len(i.MaxSDUMToS) != int(i.CISCount) {
		panic("len(i.MaxSDUMToS) != int(i.CISCount)")
	}
	for _, m := range i.MaxSDUMToS {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.MaxSDUSToM) != int(i.CISCount) {
		panic("len(i.MaxSDUSToM) != int(i.CISCount)")
	}
	for _, m := range i.MaxSDUSToM {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.MaxPDUMToS) != int(i.CISCount) {
		panic("len(i.MaxPDUMToS) != int(i.CISCount)")
	}
	for _, m := range i.MaxPDUMToS {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.MaxPDUSToM) != int(i.CISCount) {
		panic("len(i.MaxPDUSToM) != int(i.CISCount)")
	}
	for _, m := range i.MaxPDUSToM {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.PHYMToS) != int(i.CISCount) {
		panic("len(i.PHYMToS) != int(i.CISCount)")
	}
	for _, m := range i.PHYMToS {
		w.PutOne(m)
	}
	if len(i.PHYSToM) != int(i.CISCount) {
		panic("len(i.PHYSToM) != int(i.CISCount)")
	}
	for _, m := range i.PHYSToM {
		w.PutOne(m)
	}
	if len(i.BNMToS) != int(i.CISCount) {
		panic("len(i.BNMToS) != int(i.CISCount)")
	}
	for _, m := range i.BNMToS {
		w.PutOne(m)
	}
	if len(i.BNSToM) != int(i.CISCount) {
		panic("len(i.BNSToM) != int(i.CISCount)")
	}
	for _, m := range i.BNSToM {
		w.PutOne(m)
	}
	return w.Data()
}

// LESetCIGParametersTestOutput represents the output of the command specified in Section 7.8.98
type LESetCIGParametersTestOutput struct {
	Status uint8
	CIGID uint8
	CISCount uint8
	ConnectionHandle []uint16
}

func (o *LESetCIGParametersTestOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.CIGID = r.GetOne()
	o.CISCount = r.GetOne()
	if cap(o.ConnectionHandle) < int(o.CISCount) {
		o.ConnectionHandle = make([]uint16, 0, int(o.CISCount))
	}
	o.ConnectionHandle = o.ConnectionHandle[:int(o.CISCount)]
	for j:=0; j<int(o.CISCount); j++ {
		o.ConnectionHandle[j] = binary.LittleEndian.Uint16(r.Get(2))
	}
	return r.Valid()
}

// LESetCIGParametersTestSync executes the command specified in Section 7.8.98 synchronously
func (c *Commands) LESetCIGParametersTestSync (params LESetCIGParametersTestInput, result *LESetCIGParametersTestOutput) (*LESetCIGParametersTestOutput, error) {
	if result == nil {
		result = &LESetCIGParametersTestOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0063}, nil)
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

// LECreateCISInput represents the input of the command specified in Section 7.8.99
type LECreateCISInput struct {
	CISCount uint8
	CISConnectionHandle []uint16
	ACLConnectionHandle []uint16
}

func (i LECreateCISInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.CISCount)
	if len(i.CISConnectionHandle) != int(i.CISCount) {
		panic("len(i.CISConnectionHandle) != int(i.CISCount)")
	}
	for _, m := range i.CISConnectionHandle {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	if len(i.ACLConnectionHandle) != int(i.CISCount) {
		panic("len(i.ACLConnectionHandle) != int(i.CISCount)")
	}
	for _, m := range i.ACLConnectionHandle {
		binary.LittleEndian.PutUint16(w.Put(2), m)
	}
	return w.Data()
}

// LECreateCISSync executes the command specified in Section 7.8.99 synchronously
func (c *Commands) LECreateCISSync (params LECreateCISInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0064}, nil)
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

// LERemoveCIGInput represents the input of the command specified in Section 7.8.100
type LERemoveCIGInput struct {
	CIGID uint8
}

func (i LERemoveCIGInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.CIGID)
	return w.Data()
}

// LERemoveCIGOutput represents the output of the command specified in Section 7.8.100
type LERemoveCIGOutput struct {
	Status uint8
	CIGID uint8
}

func (o *LERemoveCIGOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.CIGID = r.GetOne()
	return r.Valid()
}

// LERemoveCIGSync executes the command specified in Section 7.8.100 synchronously
func (c *Commands) LERemoveCIGSync (params LERemoveCIGInput, result *LERemoveCIGOutput) (*LERemoveCIGOutput, error) {
	if result == nil {
		result = &LERemoveCIGOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0065}, nil)
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

// LEAcceptCISRequestInput represents the input of the command specified in Section 7.8.101
type LEAcceptCISRequestInput struct {
	BIGHandle uint8
	AdvertisingHandle uint8
	NumBIS uint8
	SDUInterval uint32
	MaxSDU uint16
	MaxTransportLatency uint16
	RTN uint8
	PHY uint8
	Packing uint8
	Framing uint8
	Encryption uint8
	BroadcastCode [16]byte
}

func (i LEAcceptCISRequestInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.BIGHandle)
	w.PutOne(i.AdvertisingHandle)
	w.PutOne(i.NumBIS)
	encodeUint24(w.Put(3), i.SDUInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxSDU)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxTransportLatency)
	w.PutOne(i.RTN)
	w.PutOne(i.PHY)
	w.PutOne(i.Packing)
	w.PutOne(i.Framing)
	w.PutOne(i.Encryption)
	copy(w.Put(16), i.BroadcastCode[:])
	return w.Data()
}

// LEAcceptCISRequestSync executes the command specified in Section 7.8.101 synchronously
func (c *Commands) LEAcceptCISRequestSync (params LEAcceptCISRequestInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0066}, nil)
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

// LECreateBIGTestInput represents the input of the command specified in Section 7.8.104
type LECreateBIGTestInput struct {
	BIGHandle uint8
	AdvertisingHandle uint8
	NumBIS uint8
	SDUInterval uint32
	ISOInterval uint16
	NSE uint8
	MaxSDU uint16
	MaxPDU uint16
	PHY uint8
	Packing uint8
	Framing uint8
	BN uint8
	IRC uint8
	PTO uint8
	Encryption uint8
	BroadcastCode [16]byte
}

func (i LECreateBIGTestInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.BIGHandle)
	w.PutOne(i.AdvertisingHandle)
	w.PutOne(i.NumBIS)
	encodeUint24(w.Put(3), i.SDUInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.ISOInterval)
	w.PutOne(i.NSE)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxSDU)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxPDU)
	w.PutOne(i.PHY)
	w.PutOne(i.Packing)
	w.PutOne(i.Framing)
	w.PutOne(i.BN)
	w.PutOne(i.IRC)
	w.PutOne(i.PTO)
	w.PutOne(i.Encryption)
	copy(w.Put(16), i.BroadcastCode[:])
	return w.Data()
}

// LECreateBIGTestSync executes the command specified in Section 7.8.104 synchronously
func (c *Commands) LECreateBIGTestSync (params LECreateBIGTestInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0069}, nil)
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

// LETerminateBIGInput represents the input of the command specified in Section 7.8.105
type LETerminateBIGInput struct {
	BIGHandle uint8
	Reason uint8
}

func (i LETerminateBIGInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.BIGHandle)
	w.PutOne(i.Reason)
	return w.Data()
}

// LETerminateBIGSync executes the command specified in Section 7.8.105 synchronously
func (c *Commands) LETerminateBIGSync (params LETerminateBIGInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006A}, nil)
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

// LEBIGCreateSyncInput represents the input of the command specified in Section 7.8.106
type LEBIGCreateSyncInput struct {
	BIGHandle uint8
	SyncHandle uint16
	Encryption uint8
	BroadcastCode [16]byte
	MSE uint8
	BIGSyncTimeout uint16
	NumBIS uint8
	BIS []uint8
}

func (i LEBIGCreateSyncInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.BIGHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncHandle)
	w.PutOne(i.Encryption)
	copy(w.Put(16), i.BroadcastCode[:])
	w.PutOne(i.MSE)
	binary.LittleEndian.PutUint16(w.Put(2), i.BIGSyncTimeout)
	w.PutOne(i.NumBIS)
	if len(i.BIS) != int(i.NumBIS) {
		panic("len(i.BIS) != int(i.NumBIS)")
	}
	for _, m := range i.BIS {
		w.PutOne(m)
	}
	return w.Data()
}

// LEBIGCreateSyncSync executes the command specified in Section 7.8.106 synchronously
func (c *Commands) LEBIGCreateSyncSync (params LEBIGCreateSyncInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006B}, nil)
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

// LEBIGTerminateSyncInput represents the input of the command specified in Section 7.8.107
type LEBIGTerminateSyncInput struct {
	BIGHandle uint8
}

func (i LEBIGTerminateSyncInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.BIGHandle)
	return w.Data()
}

// LEBIGTerminateSyncOutput represents the output of the command specified in Section 7.8.107
type LEBIGTerminateSyncOutput struct {
	Status uint8
	BIGHandle uint8
}

func (o *LEBIGTerminateSyncOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.BIGHandle = r.GetOne()
	return r.Valid()
}

// LEBIGTerminateSyncSync executes the command specified in Section 7.8.107 synchronously
func (c *Commands) LEBIGTerminateSyncSync (params LEBIGTerminateSyncInput, result *LEBIGTerminateSyncOutput) (*LEBIGTerminateSyncOutput, error) {
	if result == nil {
		result = &LEBIGTerminateSyncOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006C}, nil)
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

// LERequestPeerSCAInput represents the input of the command specified in Section 7.8.108
type LERequestPeerSCAInput struct {
	ConnectionHandle uint16
}

func (i LERequestPeerSCAInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LERequestPeerSCASync executes the command specified in Section 7.8.108 synchronously
func (c *Commands) LERequestPeerSCASync (params LERequestPeerSCAInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006D}, nil)
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

// LESetupISODataPathInput represents the input of the command specified in Section 7.8.109
type LESetupISODataPathInput struct {
	ConnectionHandle uint16
	DataPathDirection uint8
	DataPathID uint8
	CodecID [5]byte
	ControllerDelay uint32
	CodecConfigurationLength uint8
	CodecConfiguration []byte
}

func (i LESetupISODataPathInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.DataPathDirection)
	w.PutOne(i.DataPathID)
	copy(w.Put(5), i.CodecID[:])
	encodeUint24(w.Put(3), i.ControllerDelay)
	w.PutOne(i.CodecConfigurationLength)
	w.PutSlice(i.CodecConfiguration)
	return w.Data()
}

// LESetupISODataPathOutput represents the output of the command specified in Section 7.8.109
type LESetupISODataPathOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetupISODataPathOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetupISODataPathSync executes the command specified in Section 7.8.109 synchronously
func (c *Commands) LESetupISODataPathSync (params LESetupISODataPathInput, result *LESetupISODataPathOutput) (*LESetupISODataPathOutput, error) {
	if result == nil {
		result = &LESetupISODataPathOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006E}, nil)
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

// LERemoveISODataPathInput represents the input of the command specified in Section 7.8.110
type LERemoveISODataPathInput struct {
	ConnectionHandle uint16
	DataPathDirection uint8
}

func (i LERemoveISODataPathInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.DataPathDirection)
	return w.Data()
}

// LERemoveISODataPathOutput represents the output of the command specified in Section 7.8.110
type LERemoveISODataPathOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LERemoveISODataPathOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LERemoveISODataPathSync executes the command specified in Section 7.8.110 synchronously
func (c *Commands) LERemoveISODataPathSync (params LERemoveISODataPathInput, result *LERemoveISODataPathOutput) (*LERemoveISODataPathOutput, error) {
	if result == nil {
		result = &LERemoveISODataPathOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006F}, nil)
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

// LEISOTransmitTestInput represents the input of the command specified in Section 7.8.111
type LEISOTransmitTestInput struct {
	ConnectionHandle uint16
	PayloadType uint8
}

func (i LEISOTransmitTestInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.PayloadType)
	return w.Data()
}

// LEISOTransmitTestOutput represents the output of the command specified in Section 7.8.111
type LEISOTransmitTestOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEISOTransmitTestOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEISOTransmitTestSync executes the command specified in Section 7.8.111 synchronously
func (c *Commands) LEISOTransmitTestSync (params LEISOTransmitTestInput, result *LEISOTransmitTestOutput) (*LEISOTransmitTestOutput, error) {
	if result == nil {
		result = &LEISOTransmitTestOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0070}, nil)
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

// LEISOReceiveTestInput represents the input of the command specified in Section 7.8.112
type LEISOReceiveTestInput struct {
	ConnectionHandle uint16
	PayloadType uint8
}

func (i LEISOReceiveTestInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.PayloadType)
	return w.Data()
}

// LEISOReceiveTestOutput represents the output of the command specified in Section 7.8.112
type LEISOReceiveTestOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEISOReceiveTestOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEISOReceiveTestSync executes the command specified in Section 7.8.112 synchronously
func (c *Commands) LEISOReceiveTestSync (params LEISOReceiveTestInput, result *LEISOReceiveTestOutput) (*LEISOReceiveTestOutput, error) {
	if result == nil {
		result = &LEISOReceiveTestOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0071}, nil)
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

// LEISOReadTestCountersInput represents the input of the command specified in Section 7.8.113
type LEISOReadTestCountersInput struct {
	ConnectionHandle uint16
}

func (i LEISOReadTestCountersInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LEISOReadTestCountersOutput represents the output of the command specified in Section 7.8.113
type LEISOReadTestCountersOutput struct {
	Status uint8
	ConnectionHandle uint16
	ReceivedPacketCount uint32
	MissedPacketCount uint32
	FailedPacketCount uint32
}

func (o *LEISOReadTestCountersOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ReceivedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	o.MissedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	o.FailedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// LEISOReadTestCountersSync executes the command specified in Section 7.8.113 synchronously
func (c *Commands) LEISOReadTestCountersSync (params LEISOReadTestCountersInput, result *LEISOReadTestCountersOutput) (*LEISOReadTestCountersOutput, error) {
	if result == nil {
		result = &LEISOReadTestCountersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0072}, nil)
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

// LEISOTestEndInput represents the input of the command specified in Section 7.8.114
type LEISOTestEndInput struct {
	ConnectionHandle uint16
}

func (i LEISOTestEndInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LEISOTestEndOutput represents the output of the command specified in Section 7.8.114
type LEISOTestEndOutput struct {
	Status uint8
	ConnectionHandle uint16
	ReceivedPacketCount uint32
	MissedPacketCount uint32
	FailedPacketCount uint32
}

func (o *LEISOTestEndOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ReceivedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	o.MissedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	o.FailedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// LEISOTestEndSync executes the command specified in Section 7.8.114 synchronously
func (c *Commands) LEISOTestEndSync (params LEISOTestEndInput, result *LEISOTestEndOutput) (*LEISOTestEndOutput, error) {
	if result == nil {
		result = &LEISOTestEndOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0073}, nil)
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

// LESetHostFeatureInput represents the input of the command specified in Section 7.8.115
type LESetHostFeatureInput struct {
	BitNumber uint8
	BitValue uint8
}

func (i LESetHostFeatureInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.BitNumber)
	w.PutOne(i.BitValue)
	return w.Data()
}

// LESetHostFeatureSync executes the command specified in Section 7.8.115 synchronously
func (c *Commands) LESetHostFeatureSync (params LESetHostFeatureInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0074}, nil)
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

// LEReadISOLinkQualityInput represents the input of the command specified in Section 7.8.116
type LEReadISOLinkQualityInput struct {
	ConnectionHandle uint16
}

func (i LEReadISOLinkQualityInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LEReadISOLinkQualityOutput represents the output of the command specified in Section 7.8.116
type LEReadISOLinkQualityOutput struct {
	Status uint8
	ConnectionHandle uint16
	TxUnACKedPackets uint32
	TxFlushedPackets uint32
	TxLastSubeventPackets uint32
	RetransmittedPackets uint32
	CRCErrorPackets uint32
	RxUnreceivedPackets uint32
	DuplicatePackets uint32
}

func (o *LEReadISOLinkQualityOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.TxUnACKedPackets = binary.LittleEndian.Uint32(r.Get(4))
	o.TxFlushedPackets = binary.LittleEndian.Uint32(r.Get(4))
	o.TxLastSubeventPackets = binary.LittleEndian.Uint32(r.Get(4))
	o.RetransmittedPackets = binary.LittleEndian.Uint32(r.Get(4))
	o.CRCErrorPackets = binary.LittleEndian.Uint32(r.Get(4))
	o.RxUnreceivedPackets = binary.LittleEndian.Uint32(r.Get(4))
	o.DuplicatePackets = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// LEReadISOLinkQualitySync executes the command specified in Section 7.8.116 synchronously
func (c *Commands) LEReadISOLinkQualitySync (params LEReadISOLinkQualityInput, result *LEReadISOLinkQualityOutput) (*LEReadISOLinkQualityOutput, error) {
	if result == nil {
		result = &LEReadISOLinkQualityOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0075}, nil)
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

// LEEnhancedReadTransmitPowerLevelInput represents the input of the command specified in Section 7.8.117
type LEEnhancedReadTransmitPowerLevelInput struct {
	ConnectionHandle uint16
	PHY uint8
}

func (i LEEnhancedReadTransmitPowerLevelInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.PHY)
	return w.Data()
}

// LEEnhancedReadTransmitPowerLevelOutput represents the output of the command specified in Section 7.8.117
type LEEnhancedReadTransmitPowerLevelOutput struct {
	Status uint8
	ConnectionHandle uint16
	PHY uint8
	CurrentTransmitPowerLevel uint8
}

func (o *LEEnhancedReadTransmitPowerLevelOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PHY = r.GetOne()
	o.CurrentTransmitPowerLevel = r.GetOne()
	return r.Valid()
}

// LEEnhancedReadTransmitPowerLevelSync executes the command specified in Section 7.8.117 synchronously
func (c *Commands) LEEnhancedReadTransmitPowerLevelSync (params LEEnhancedReadTransmitPowerLevelInput, result *LEEnhancedReadTransmitPowerLevelOutput) (*LEEnhancedReadTransmitPowerLevelOutput, error) {
	if result == nil {
		result = &LEEnhancedReadTransmitPowerLevelOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0076}, nil)
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

// LEReadRemoteTransmitPowerLevelInput represents the input of the command specified in Section 7.8.118
type LEReadRemoteTransmitPowerLevelInput struct {
	ConnectionHandle uint16
	PHY uint8
}

func (i LEReadRemoteTransmitPowerLevelInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.PHY)
	return w.Data()
}

// LEReadRemoteTransmitPowerLevelSync executes the command specified in Section 7.8.118 synchronously
func (c *Commands) LEReadRemoteTransmitPowerLevelSync (params LEReadRemoteTransmitPowerLevelInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0077}, nil)
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

// LESetPathLossReportingParametersInput represents the input of the command specified in Section 7.8.119
type LESetPathLossReportingParametersInput struct {
	ConnectionHandle uint16
	HighThreshold uint8
	HighHysteresis uint8
	LowThreshold uint8
	LowHysteresis uint8
	MinTimeSpent uint16
}

func (i LESetPathLossReportingParametersInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.HighThreshold)
	w.PutOne(i.HighHysteresis)
	w.PutOne(i.LowThreshold)
	w.PutOne(i.LowHysteresis)
	binary.LittleEndian.PutUint16(w.Put(2), i.MinTimeSpent)
	return w.Data()
}

// LESetPathLossReportingParametersOutput represents the output of the command specified in Section 7.8.119
type LESetPathLossReportingParametersOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetPathLossReportingParametersOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetPathLossReportingParametersSync executes the command specified in Section 7.8.119 synchronously
func (c *Commands) LESetPathLossReportingParametersSync (params LESetPathLossReportingParametersInput, result *LESetPathLossReportingParametersOutput) (*LESetPathLossReportingParametersOutput, error) {
	if result == nil {
		result = &LESetPathLossReportingParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0078}, nil)
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

// LESetPathLossReportingEnableInput represents the input of the command specified in Section 7.8.120
type LESetPathLossReportingEnableInput struct {
	ConnectionHandle uint16
	Enable uint8
}

func (i LESetPathLossReportingEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.Enable)
	return w.Data()
}

// LESetPathLossReportingEnableOutput represents the output of the command specified in Section 7.8.120
type LESetPathLossReportingEnableOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetPathLossReportingEnableOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetPathLossReportingEnableSync executes the command specified in Section 7.8.120 synchronously
func (c *Commands) LESetPathLossReportingEnableSync (params LESetPathLossReportingEnableInput, result *LESetPathLossReportingEnableOutput) (*LESetPathLossReportingEnableOutput, error) {
	if result == nil {
		result = &LESetPathLossReportingEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0079}, nil)
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

// LESetTransmitPowerReportingEnableInput represents the input of the command specified in Section 7.8.121
type LESetTransmitPowerReportingEnableInput struct {
	ConnectionHandle uint16
	LocalEnable uint8
	RemoteEnable uint8
}

func (i LESetTransmitPowerReportingEnableInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.LocalEnable)
	w.PutOne(i.RemoteEnable)
	return w.Data()
}

// LESetTransmitPowerReportingEnableOutput represents the output of the command specified in Section 7.8.121
type LESetTransmitPowerReportingEnableOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetTransmitPowerReportingEnableOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetTransmitPowerReportingEnableSync executes the command specified in Section 7.8.121 synchronously
func (c *Commands) LESetTransmitPowerReportingEnableSync (params LESetTransmitPowerReportingEnableInput, result *LESetTransmitPowerReportingEnableOutput) (*LESetTransmitPowerReportingEnableOutput, error) {
	if result == nil {
		result = &LESetTransmitPowerReportingEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x007A}, nil)
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

