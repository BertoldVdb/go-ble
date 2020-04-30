package hcicommands

import (
	"encoding/binary"
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
)

// StatusResetFailedContactCounterInput represents the input of the command specified in Section 7.5.2
type StatusResetFailedContactCounterInput struct {
	Handle uint16
}

func (i StatusResetFailedContactCounterInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.Handle)
	return w.Data()
}

// StatusResetFailedContactCounterOutput represents the output of the command specified in Section 7.5.2
type StatusResetFailedContactCounterOutput struct {
	Status uint8
	Handle uint16
}

func (o *StatusResetFailedContactCounterOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// StatusResetFailedContactCounterSync executes the command specified in Section 7.5.2 synchronously
func (c *Commands) StatusResetFailedContactCounterSync (params StatusResetFailedContactCounterInput, result *StatusResetFailedContactCounterOutput) (*StatusResetFailedContactCounterOutput, error) {
	if result == nil {
		result = &StatusResetFailedContactCounterOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 5, OCF: 0x0002}, nil)
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

// StatusReadLinkQualityInput represents the input of the command specified in Section 7.5.3
type StatusReadLinkQualityInput struct {
	Handle uint16
}

func (i StatusReadLinkQualityInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.Handle)
	return w.Data()
}

// StatusReadLinkQualityOutput represents the output of the command specified in Section 7.5.3
type StatusReadLinkQualityOutput struct {
	Status uint8
	Handle uint16
	LinkQuality uint8
}

func (o *StatusReadLinkQualityOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	o.LinkQuality = r.GetOne()
	return r.Valid()
}

// StatusReadLinkQualitySync executes the command specified in Section 7.5.3 synchronously
func (c *Commands) StatusReadLinkQualitySync (params StatusReadLinkQualityInput, result *StatusReadLinkQualityOutput) (*StatusReadLinkQualityOutput, error) {
	if result == nil {
		result = &StatusReadLinkQualityOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 5, OCF: 0x0003}, nil)
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

// StatusReadRSSIInput represents the input of the command specified in Section 7.5.4
type StatusReadRSSIInput struct {
	Handle uint16
}

func (i StatusReadRSSIInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.Handle)
	return w.Data()
}

// StatusReadRSSIOutput represents the output of the command specified in Section 7.5.4
type StatusReadRSSIOutput struct {
	Status uint8
	Handle uint16
	RSSI uint8
}

func (o *StatusReadRSSIOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	o.RSSI = r.GetOne()
	return r.Valid()
}

// StatusReadRSSISync executes the command specified in Section 7.5.4 synchronously
func (c *Commands) StatusReadRSSISync (params StatusReadRSSIInput, result *StatusReadRSSIOutput) (*StatusReadRSSIOutput, error) {
	if result == nil {
		result = &StatusReadRSSIOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 5, OCF: 0x0005}, nil)
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

// StatusReadAFHChannelMapInput represents the input of the command specified in Section 7.5.5
type StatusReadAFHChannelMapInput struct {
	ConnectionHandle uint16
}

func (i StatusReadAFHChannelMapInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// StatusReadAFHChannelMapOutput represents the output of the command specified in Section 7.5.5
type StatusReadAFHChannelMapOutput struct {
	Status uint8
	ConnectionHandle uint16
	AFHMode uint8
	AFHChannelMap [10]byte
}

func (o *StatusReadAFHChannelMapOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.AFHMode = r.GetOne()
	copy(o.AFHChannelMap[:], r.Get(10))
	return r.Valid()
}

// StatusReadAFHChannelMapSync executes the command specified in Section 7.5.5 synchronously
func (c *Commands) StatusReadAFHChannelMapSync (params StatusReadAFHChannelMapInput, result *StatusReadAFHChannelMapOutput) (*StatusReadAFHChannelMapOutput, error) {
	if result == nil {
		result = &StatusReadAFHChannelMapOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 5, OCF: 0x0006}, nil)
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

// StatusReadClockInput represents the input of the command specified in Section 7.5.6
type StatusReadClockInput struct {
	ConnectionHandle uint16
	WhichClock []byte
}

func (i StatusReadClockInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutSlice(i.WhichClock)
	return w.Data()
}

// StatusReadClockOutput represents the output of the command specified in Section 7.5.6
type StatusReadClockOutput struct {
	Status uint8
	ConnectionHandle uint16
	Clock uint32
	Accuracy uint16
}

func (o *StatusReadClockOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.Clock = binary.LittleEndian.Uint32(r.Get(4))
	o.Accuracy = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// StatusReadClockSync executes the command specified in Section 7.5.6 synchronously
func (c *Commands) StatusReadClockSync (params StatusReadClockInput, result *StatusReadClockOutput) (*StatusReadClockOutput, error) {
	if result == nil {
		result = &StatusReadClockOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 5, OCF: 0x0007}, nil)
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

// StatusReadEncryptionKeySizeInput represents the input of the command specified in Section 7.5.7
type StatusReadEncryptionKeySizeInput struct {
	ConnectionHandle uint16
}

func (i StatusReadEncryptionKeySizeInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// StatusReadEncryptionKeySizeOutput represents the output of the command specified in Section 7.5.7
type StatusReadEncryptionKeySizeOutput struct {
	Status uint8
	ConnectionHandle uint16
	KeySize uint8
}

func (o *StatusReadEncryptionKeySizeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.KeySize = r.GetOne()
	return r.Valid()
}

// StatusReadEncryptionKeySizeSync executes the command specified in Section 7.5.7 synchronously
func (c *Commands) StatusReadEncryptionKeySizeSync (params StatusReadEncryptionKeySizeInput, result *StatusReadEncryptionKeySizeOutput) (*StatusReadEncryptionKeySizeOutput, error) {
	if result == nil {
		result = &StatusReadEncryptionKeySizeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 5, OCF: 0x0008}, nil)
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

// StatusReadLocalAMPInfoOutput represents the output of the command specified in Section 7.5.8
type StatusReadLocalAMPInfoOutput struct {
	Status uint8
	AMPStatus uint8
	TotalBandwidth uint32
	MaxGuaranteedBandwidth uint32
	MinLatency uint32
	MaxPDUSize uint16
	ControllerType uint8
	PALCapabilities uint16
	MaxAMPAssocLength uint16
	MaxFlushTimeout uint32
	BestEffortFlushTimeout uint32
}

func (o *StatusReadLocalAMPInfoOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.AMPStatus = r.GetOne()
	o.TotalBandwidth = binary.LittleEndian.Uint32(r.Get(4))
	o.MaxGuaranteedBandwidth = binary.LittleEndian.Uint32(r.Get(4))
	o.MinLatency = binary.LittleEndian.Uint32(r.Get(4))
	o.MaxPDUSize = binary.LittleEndian.Uint16(r.Get(2))
	o.ControllerType = r.GetOne()
	o.PALCapabilities = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxAMPAssocLength = binary.LittleEndian.Uint16(r.Get(2))
	o.MaxFlushTimeout = binary.LittleEndian.Uint32(r.Get(4))
	o.BestEffortFlushTimeout = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// StatusReadLocalAMPInfoSync executes the command specified in Section 7.5.8 synchronously
func (c *Commands) StatusReadLocalAMPInfoSync (result *StatusReadLocalAMPInfoOutput) (*StatusReadLocalAMPInfoOutput, error) {
	if result == nil {
		result = &StatusReadLocalAMPInfoOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 5, OCF: 0x0009}, nil)
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

// StatusReadLocalAMPASSOCInput represents the input of the command specified in Section 7.5.9
type StatusReadLocalAMPASSOCInput struct {
	PhysicalLinkHandle uint8
	LengthSoFar uint16
	AMPAssocLength uint16
}

func (i StatusReadLocalAMPASSOCInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PhysicalLinkHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.LengthSoFar)
	binary.LittleEndian.PutUint16(w.Put(2), i.AMPAssocLength)
	return w.Data()
}

// StatusReadLocalAMPASSOCOutput represents the output of the command specified in Section 7.5.9
type StatusReadLocalAMPASSOCOutput struct {
	Status uint8
	PhysicalLinkHandle uint8
	AMPAssocRemainingLength uint16
	AMPAssocFragment []byte
}

func (o *StatusReadLocalAMPASSOCOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.PhysicalLinkHandle = r.GetOne()
	o.AMPAssocRemainingLength = binary.LittleEndian.Uint16(r.Get(2))
	o.AMPAssocFragment = append(o.AMPAssocFragment[:0], r.GetRemainder()...)
	return r.Valid()
}

// StatusReadLocalAMPASSOCSync executes the command specified in Section 7.5.9 synchronously
func (c *Commands) StatusReadLocalAMPASSOCSync (params StatusReadLocalAMPASSOCInput, result *StatusReadLocalAMPASSOCOutput) (*StatusReadLocalAMPASSOCOutput, error) {
	if result == nil {
		result = &StatusReadLocalAMPASSOCOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 5, OCF: 0x000A}, nil)
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

// StatusWriteRemoteAMPASSOCInput represents the input of the command specified in Section 7.5.10
type StatusWriteRemoteAMPASSOCInput struct {
	PhysicalLinkHandle uint8
	LengthSoFar uint16
	AMPAssocRemainingLength uint16
	AMPAssocFragment []byte
}

func (i StatusWriteRemoteAMPASSOCInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.PhysicalLinkHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.LengthSoFar)
	binary.LittleEndian.PutUint16(w.Put(2), i.AMPAssocRemainingLength)
	w.PutSlice(i.AMPAssocFragment)
	return w.Data()
}

// StatusWriteRemoteAMPASSOCOutput represents the output of the command specified in Section 7.5.10
type StatusWriteRemoteAMPASSOCOutput struct {
	Status uint8
	PhysicalLinkHandle uint8
}

func (o *StatusWriteRemoteAMPASSOCOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.PhysicalLinkHandle = r.GetOne()
	return r.Valid()
}

// StatusWriteRemoteAMPASSOCSync executes the command specified in Section 7.5.10 synchronously
func (c *Commands) StatusWriteRemoteAMPASSOCSync (params StatusWriteRemoteAMPASSOCInput, result *StatusWriteRemoteAMPASSOCOutput) (*StatusWriteRemoteAMPASSOCOutput, error) {
	if result == nil {
		result = &StatusWriteRemoteAMPASSOCOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 5, OCF: 0x000B}, nil)
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

// StatusGetMWSTransportLayerConfigurationOutput represents the output of the command specified in Section 7.5.11
type StatusGetMWSTransportLayerConfigurationOutput struct {
	Status uint8
	NumTransports uint8
	TransportLayer []uint8
	NumBaudRates []uint8
	ToMWSBaudRate []uint32
	FromMWSBaudRate []uint32
}

func (o *StatusGetMWSTransportLayerConfigurationOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.NumTransports = r.GetOne()
	if cap(o.TransportLayer) < int(o.NumTransports) {
		o.TransportLayer = make([]uint8, 0, int(o.NumTransports))
	}
	o.TransportLayer = o.TransportLayer[:int(o.NumTransports)]
	for j:=0; j<int(o.NumTransports); j++ {
		o.TransportLayer[j] = r.GetOne()
	}
	if cap(o.NumBaudRates) < int(o.NumTransports) {
		o.NumBaudRates = make([]uint8, 0, int(o.NumTransports))
	}
	o.NumBaudRates = o.NumBaudRates[:int(o.NumTransports)]
	for j:=0; j<int(o.NumTransports); j++ {
		o.NumBaudRates[j] = r.GetOne()
	}
	var0 := 0
	for _, m := range o.NumBaudRates {
		var0 += int(m)
	}
	if cap(o.ToMWSBaudRate) < var0 {
		o.ToMWSBaudRate = make([]uint32, 0, var0)
	}
	o.ToMWSBaudRate = o.ToMWSBaudRate[:var0]
	for j:=0; j<var0; j++ {
		o.ToMWSBaudRate[j] = binary.LittleEndian.Uint32(r.Get(4))
	}
	if cap(o.FromMWSBaudRate) < var0 {
		o.FromMWSBaudRate = make([]uint32, 0, var0)
	}
	o.FromMWSBaudRate = o.FromMWSBaudRate[:var0]
	for j:=0; j<var0; j++ {
		o.FromMWSBaudRate[j] = binary.LittleEndian.Uint32(r.Get(4))
	}
	return r.Valid()
}

// StatusGetMWSTransportLayerConfigurationSync executes the command specified in Section 7.5.11 synchronously
func (c *Commands) StatusGetMWSTransportLayerConfigurationSync (result *StatusGetMWSTransportLayerConfigurationOutput) (*StatusGetMWSTransportLayerConfigurationOutput, error) {
	if result == nil {
		result = &StatusGetMWSTransportLayerConfigurationOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 5, OCF: 0x000C}, nil)
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

// StatusSetTriggeredClockCaptureInput represents the input of the command specified in Section 7.5.12
type StatusSetTriggeredClockCaptureInput struct {
	ConnectionHandle uint16
	Enable uint8
	WhichClock uint8
	LPOAllowed uint8
	NumClockCapturesToFilter uint8
}

func (i StatusSetTriggeredClockCaptureInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.Enable)
	w.PutOne(i.WhichClock)
	w.PutOne(i.LPOAllowed)
	w.PutOne(i.NumClockCapturesToFilter)
	return w.Data()
}

// StatusSetTriggeredClockCaptureSync executes the command specified in Section 7.5.12 synchronously
func (c *Commands) StatusSetTriggeredClockCaptureSync (params StatusSetTriggeredClockCaptureInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 5, OCF: 0x000D}, nil)
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

