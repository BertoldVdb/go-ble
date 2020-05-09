package hcicommands

import (
	"encoding/binary"
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

// InformationalReadLocalSupportedCommandsOutput represents the output of the command specified in Section 7.4.2
type InformationalReadLocalSupportedCommandsOutput struct {
	Status uint8
	SupportedCommands [64]byte
}

func (o *InformationalReadLocalSupportedCommandsOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	copy(o.SupportedCommands[:], r.Get(64))
	return r.Valid()
}

// InformationalReadLocalSupportedCommandsSync executes the command specified in Section 7.4.2 synchronously
func (c *Commands) InformationalReadLocalSupportedCommandsSync (result *InformationalReadLocalSupportedCommandsOutput) (*InformationalReadLocalSupportedCommandsOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("InformationalReadLocalSupportedCommands started")
	}
	if result == nil {
		result = &InformationalReadLocalSupportedCommandsOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 4, OCF: 0x0002}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "1result": result,
		}).Debug("InformationalReadLocalSupportedCommands completed")
	}

	 return result, err
}
// InformationalReadLocalSupportedFeaturesOutput represents the output of the command specified in Section 7.4.3
type InformationalReadLocalSupportedFeaturesOutput struct {
	Status uint8
	LMPFeatures uint64
}

func (o *InformationalReadLocalSupportedFeaturesOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.LMPFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// InformationalReadLocalSupportedFeaturesSync executes the command specified in Section 7.4.3 synchronously
func (c *Commands) InformationalReadLocalSupportedFeaturesSync (result *InformationalReadLocalSupportedFeaturesOutput) (*InformationalReadLocalSupportedFeaturesOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("InformationalReadLocalSupportedFeatures started")
	}
	if result == nil {
		result = &InformationalReadLocalSupportedFeaturesOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 4, OCF: 0x0003}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "1result": result,
		}).Debug("InformationalReadLocalSupportedFeatures completed")
	}

	 return result, err
}
// InformationalReadLocalExtendedFeaturesInput represents the input of the command specified in Section 7.4.4
type InformationalReadLocalExtendedFeaturesInput struct {
	PageNumber uint8
}

func (i InformationalReadLocalExtendedFeaturesInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.PageNumber))
	return w.Data
}

// InformationalReadLocalExtendedFeaturesOutput represents the output of the command specified in Section 7.4.4
type InformationalReadLocalExtendedFeaturesOutput struct {
	Status uint8
	PageNumber uint8
	MaximumPageNumber uint8
	ExtendedLMPFeatures uint64
}

func (o *InformationalReadLocalExtendedFeaturesOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.PageNumber = uint8(r.GetOne())
	o.MaximumPageNumber = uint8(r.GetOne())
	o.ExtendedLMPFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// InformationalReadLocalExtendedFeaturesSync executes the command specified in Section 7.4.4 synchronously
func (c *Commands) InformationalReadLocalExtendedFeaturesSync (params InformationalReadLocalExtendedFeaturesInput, result *InformationalReadLocalExtendedFeaturesOutput) (*InformationalReadLocalExtendedFeaturesOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("InformationalReadLocalExtendedFeatures started")
	}
	if result == nil {
		result = &InformationalReadLocalExtendedFeaturesOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 4, OCF: 0x0004}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "0params": params,
			 "1result": result,
		}).Debug("InformationalReadLocalExtendedFeatures completed")
	}

	 return result, err
}
// InformationalReadBufferSizeOutput represents the output of the command specified in Section 7.4.5
type InformationalReadBufferSizeOutput struct {
	Status uint8
	ACLDataPacketLength uint16
	SynchronousDataPacketLength uint8
	TotalNumACLDataPackets uint16
	TotalNumSynchronousDataPackets uint16
}

func (o *InformationalReadBufferSizeOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ACLDataPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.SynchronousDataPacketLength = uint8(r.GetOne())
	o.TotalNumACLDataPackets = binary.LittleEndian.Uint16(r.Get(2))
	o.TotalNumSynchronousDataPackets = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// InformationalReadBufferSizeSync executes the command specified in Section 7.4.5 synchronously
func (c *Commands) InformationalReadBufferSizeSync (result *InformationalReadBufferSizeOutput) (*InformationalReadBufferSizeOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("InformationalReadBufferSize started")
	}
	if result == nil {
		result = &InformationalReadBufferSizeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 4, OCF: 0x0005}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "1result": result,
		}).Debug("InformationalReadBufferSize completed")
	}

	 return result, err
}
// InformationalReadBDADDROutput represents the output of the command specified in Section 7.4.6
type InformationalReadBDADDROutput struct {
	Status uint8
	BDADDR bleutil.MacAddr
}

func (o *InformationalReadBDADDROutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.BDADDR.Decode(r.Get(6))
	return r.Valid()
}

// InformationalReadBDADDRSync executes the command specified in Section 7.4.6 synchronously
func (c *Commands) InformationalReadBDADDRSync (result *InformationalReadBDADDROutput) (*InformationalReadBDADDROutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("InformationalReadBDADDR started")
	}
	if result == nil {
		result = &InformationalReadBDADDROutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 4, OCF: 0x0009}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "1result": result,
		}).Debug("InformationalReadBDADDR completed")
	}

	 return result, err
}
// InformationalReadDataBlockSizeOutput represents the output of the command specified in Section 7.4.7
type InformationalReadDataBlockSizeOutput struct {
	Status uint8
	MaxACLDataPacketLength uint16
	DataBlockLength uint16
	TotalNumDataBlocks uint16
}

func (o *InformationalReadDataBlockSizeOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.MaxACLDataPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.DataBlockLength = binary.LittleEndian.Uint16(r.Get(2))
	o.TotalNumDataBlocks = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// InformationalReadDataBlockSizeSync executes the command specified in Section 7.4.7 synchronously
func (c *Commands) InformationalReadDataBlockSizeSync (result *InformationalReadDataBlockSizeOutput) (*InformationalReadDataBlockSizeOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("InformationalReadDataBlockSize started")
	}
	if result == nil {
		result = &InformationalReadDataBlockSizeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 4, OCF: 0x000A}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "1result": result,
		}).Debug("InformationalReadDataBlockSize completed")
	}

	 return result, err
}
// InformationalReadLocalSupportedCodecsOutput represents the output of the command specified in Section 7.4.8
type InformationalReadLocalSupportedCodecsOutput struct {
	Status uint8
	NumSupportedStandardCodecs uint8
	StandardCodecID []uint8
	NumSupportedVendorSpecificCodecs uint8
	VendorSpecificCodecID []uint32
}

func (o *InformationalReadLocalSupportedCodecsOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.NumSupportedStandardCodecs = uint8(r.GetOne())
	if cap(o.StandardCodecID) < int(o.NumSupportedStandardCodecs) {
		o.StandardCodecID = make([]uint8, 0, int(o.NumSupportedStandardCodecs))
	}
	o.StandardCodecID = o.StandardCodecID[:int(o.NumSupportedStandardCodecs)]
	for j:=0; j<int(o.NumSupportedStandardCodecs); j++ {
		o.StandardCodecID[j] = uint8(r.GetOne())
	}
	o.NumSupportedVendorSpecificCodecs = uint8(r.GetOne())
	if cap(o.VendorSpecificCodecID) < int(o.NumSupportedVendorSpecificCodecs) {
		o.VendorSpecificCodecID = make([]uint32, 0, int(o.NumSupportedVendorSpecificCodecs))
	}
	o.VendorSpecificCodecID = o.VendorSpecificCodecID[:int(o.NumSupportedVendorSpecificCodecs)]
	for j:=0; j<int(o.NumSupportedVendorSpecificCodecs); j++ {
		o.VendorSpecificCodecID[j] = binary.LittleEndian.Uint32(r.Get(4))
	}
	return r.Valid()
}

// InformationalReadLocalSupportedCodecsSync executes the command specified in Section 7.4.8 synchronously
func (c *Commands) InformationalReadLocalSupportedCodecsSync (result *InformationalReadLocalSupportedCodecsOutput) (*InformationalReadLocalSupportedCodecsOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("InformationalReadLocalSupportedCodecs started")
	}
	if result == nil {
		result = &InformationalReadLocalSupportedCodecsOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 4, OCF: 0x000B}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "1result": result,
		}).Debug("InformationalReadLocalSupportedCodecs completed")
	}

	 return result, err
}
// InformationalReadLocalSupportedCodecsV2Output represents the output of the command specified in Section 7.4.8
type InformationalReadLocalSupportedCodecsV2Output struct {
	Status uint8
	NumSupportedStandardCodecs uint8
	StandardCodecID []uint8
	StandardCodecTransport []uint8
	NumSupportedVendorSpecificCodecs uint8
	VendorSpecificCodecID []uint32
	VendorSpecificCodecTransport []uint8
}

func (o *InformationalReadLocalSupportedCodecsV2Output) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.NumSupportedStandardCodecs = uint8(r.GetOne())
	if cap(o.StandardCodecID) < int(o.NumSupportedStandardCodecs) {
		o.StandardCodecID = make([]uint8, 0, int(o.NumSupportedStandardCodecs))
	}
	o.StandardCodecID = o.StandardCodecID[:int(o.NumSupportedStandardCodecs)]
	for j:=0; j<int(o.NumSupportedStandardCodecs); j++ {
		o.StandardCodecID[j] = uint8(r.GetOne())
	}
	if cap(o.StandardCodecTransport) < int(o.NumSupportedStandardCodecs) {
		o.StandardCodecTransport = make([]uint8, 0, int(o.NumSupportedStandardCodecs))
	}
	o.StandardCodecTransport = o.StandardCodecTransport[:int(o.NumSupportedStandardCodecs)]
	for j:=0; j<int(o.NumSupportedStandardCodecs); j++ {
		o.StandardCodecTransport[j] = uint8(r.GetOne())
	}
	o.NumSupportedVendorSpecificCodecs = uint8(r.GetOne())
	if cap(o.VendorSpecificCodecID) < int(o.NumSupportedVendorSpecificCodecs) {
		o.VendorSpecificCodecID = make([]uint32, 0, int(o.NumSupportedVendorSpecificCodecs))
	}
	o.VendorSpecificCodecID = o.VendorSpecificCodecID[:int(o.NumSupportedVendorSpecificCodecs)]
	for j:=0; j<int(o.NumSupportedVendorSpecificCodecs); j++ {
		o.VendorSpecificCodecID[j] = binary.LittleEndian.Uint32(r.Get(4))
	}
	if cap(o.VendorSpecificCodecTransport) < int(o.NumSupportedVendorSpecificCodecs) {
		o.VendorSpecificCodecTransport = make([]uint8, 0, int(o.NumSupportedVendorSpecificCodecs))
	}
	o.VendorSpecificCodecTransport = o.VendorSpecificCodecTransport[:int(o.NumSupportedVendorSpecificCodecs)]
	for j:=0; j<int(o.NumSupportedVendorSpecificCodecs); j++ {
		o.VendorSpecificCodecTransport[j] = uint8(r.GetOne())
	}
	return r.Valid()
}

// InformationalReadLocalSupportedCodecsV2Sync executes the command specified in Section 7.4.8 synchronously
func (c *Commands) InformationalReadLocalSupportedCodecsV2Sync (result *InformationalReadLocalSupportedCodecsV2Output) (*InformationalReadLocalSupportedCodecsV2Output, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("InformationalReadLocalSupportedCodecsV2 started")
	}
	if result == nil {
		result = &InformationalReadLocalSupportedCodecsV2Output{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 4, OCF: 0x000D}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "1result": result,
		}).Debug("InformationalReadLocalSupportedCodecsV2 completed")
	}

	 return result, err
}
// InformationalReadLocalSimplePairingOptionsOutput represents the output of the command specified in Section 7.4.9
type InformationalReadLocalSimplePairingOptionsOutput struct {
	Status uint8
	SimplePairingOptions uint8
	MaxEncryptionKeySize uint8
}

func (o *InformationalReadLocalSimplePairingOptionsOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.SimplePairingOptions = uint8(r.GetOne())
	o.MaxEncryptionKeySize = uint8(r.GetOne())
	return r.Valid()
}

// InformationalReadLocalSimplePairingOptionsSync executes the command specified in Section 7.4.9 synchronously
func (c *Commands) InformationalReadLocalSimplePairingOptionsSync (result *InformationalReadLocalSimplePairingOptionsOutput) (*InformationalReadLocalSimplePairingOptionsOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("InformationalReadLocalSimplePairingOptions started")
	}
	if result == nil {
		result = &InformationalReadLocalSimplePairingOptionsOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 4, OCF: 0x000C}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "1result": result,
		}).Debug("InformationalReadLocalSimplePairingOptions completed")
	}

	 return result, err
}
// InformationalReadLocalSupportedCodecCapabilitiesInput represents the input of the command specified in Section 7.4.10
type InformationalReadLocalSupportedCodecCapabilitiesInput struct {
	Handle uint16
}

func (i InformationalReadLocalSupportedCodecCapabilitiesInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.Handle)
	return w.Data
}

// InformationalReadLocalSupportedCodecCapabilitiesOutput represents the output of the command specified in Section 7.4.10
type InformationalReadLocalSupportedCodecCapabilitiesOutput struct {
	Status uint8
	Handle uint16
	FailedContactCounter uint16
}

func (o *InformationalReadLocalSupportedCodecCapabilitiesOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.Handle = binary.LittleEndian.Uint16(r.Get(2))
	o.FailedContactCounter = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// InformationalReadLocalSupportedCodecCapabilitiesSync executes the command specified in Section 7.4.10 synchronously
func (c *Commands) InformationalReadLocalSupportedCodecCapabilitiesSync (params InformationalReadLocalSupportedCodecCapabilitiesInput, result *InformationalReadLocalSupportedCodecCapabilitiesOutput) (*InformationalReadLocalSupportedCodecCapabilitiesOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("InformationalReadLocalSupportedCodecCapabilities started")
	}
	if result == nil {
		result = &InformationalReadLocalSupportedCodecCapabilitiesOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 4, OCF: 0x000E}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	if !result.decode(response) {
		err = ErrorMalformed
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "0params": params,
			 "1result": result,
		}).Debug("InformationalReadLocalSupportedCodecCapabilities completed")
	}

	 return result, err
}
