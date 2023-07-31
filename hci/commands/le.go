package hcicommands

import (
	"encoding/binary"
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

// LESetEventMaskInput represents the input of the command specified in Section 7.8.1
type LESetEventMaskInput struct {
	LEEventMask uint64
}

func (i LESetEventMaskInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint64(w.Put(8), i.LEEventMask)
	return w.Data
}

// LESetEventMaskSync executes the command specified in Section 7.8.1 synchronously
func (c *Commands) LESetEventMaskSync (params LESetEventMaskInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetEventMask started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0001}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetEventMask completed")
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
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.LEACLDataPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.TotalNumLEACLDataPackets = uint8(r.GetOne())
	return r.Valid()
}

// LEReadBufferSizeSync executes the command specified in Section 7.8.2 synchronously
func (c *Commands) LEReadBufferSizeSync (result *LEReadBufferSizeOutput) (*LEReadBufferSizeOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadBufferSize started")
	}
	if result == nil {
		result = &LEReadBufferSizeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0002}, nil)
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
		}).Debug("LEReadBufferSize completed")
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
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.LEACLDataPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.TotalNumLEACLDataPackets = uint8(r.GetOne())
	o.ISODataPacketLength = binary.LittleEndian.Uint16(r.Get(2))
	o.TotalNumISODataPackets = uint8(r.GetOne())
	return r.Valid()
}

// LEReadBufferSizeV2Sync executes the command specified in Section 7.8.2 synchronously
func (c *Commands) LEReadBufferSizeV2Sync (result *LEReadBufferSizeV2Output) (*LEReadBufferSizeV2Output, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadBufferSizeV2 started")
	}
	if result == nil {
		result = &LEReadBufferSizeV2Output{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0060}, nil)
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
		}).Debug("LEReadBufferSizeV2 completed")
	}

	 return result, err
}
// LEReadLocalSupportedFeaturesOutput represents the output of the command specified in Section 7.8.3
type LEReadLocalSupportedFeaturesOutput struct {
	Status uint8
	LEFeatures uint64
}

func (o *LEReadLocalSupportedFeaturesOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.LEFeatures = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// LEReadLocalSupportedFeaturesSync executes the command specified in Section 7.8.3 synchronously
func (c *Commands) LEReadLocalSupportedFeaturesSync (result *LEReadLocalSupportedFeaturesOutput) (*LEReadLocalSupportedFeaturesOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadLocalSupportedFeatures started")
	}
	if result == nil {
		result = &LEReadLocalSupportedFeaturesOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0003}, nil)
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
		}).Debug("LEReadLocalSupportedFeatures completed")
	}

	 return result, err
}
// LESetRandomAddressInput represents the input of the command specified in Section 7.8.4
type LESetRandomAddressInput struct {
	RandomAddess bleutil.MacAddr
}

func (i LESetRandomAddressInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	i.RandomAddess.Encode(w.Put(6))
	return w.Data
}

// LESetRandomAddressSync executes the command specified in Section 7.8.4 synchronously
func (c *Commands) LESetRandomAddressSync (params LESetRandomAddressInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetRandomAddress started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0005}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetRandomAddress completed")
	}

	 return err
}
// LESetAdvertisingParametersInput represents the input of the command specified in Section 7.8.5
type LESetAdvertisingParametersInput struct {
	AdvertisingIntervalMin uint16
	AdvertisingIntervalMax uint16
	AdvertisingType uint8
	OwnAddressType bleutil.MacAddrType
	PeerAddressType bleutil.MacAddrType
	PeerAddress bleutil.MacAddr
	AdvertisingChannelMap uint8
	AdvertisingFilterPolicy uint8
}

func (i LESetAdvertisingParametersInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.AdvertisingIntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.AdvertisingIntervalMax)
	w.PutOne(uint8(i.AdvertisingType))
	w.PutOne(uint8(i.OwnAddressType))
	w.PutOne(uint8(i.PeerAddressType))
	i.PeerAddress.Encode(w.Put(6))
	w.PutOne(uint8(i.AdvertisingChannelMap))
	w.PutOne(uint8(i.AdvertisingFilterPolicy))
	return w.Data
}

// LESetAdvertisingParametersSync executes the command specified in Section 7.8.5 synchronously
func (c *Commands) LESetAdvertisingParametersSync (params LESetAdvertisingParametersInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetAdvertisingParameters started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0006}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "0params": params,
		}).Debug("LESetAdvertisingParameters completed")
	}

	 return err
}
// LEReadAdvertisingPhysicalChannelTxPowerOutput represents the output of the command specified in Section 7.8.6
type LEReadAdvertisingPhysicalChannelTxPowerOutput struct {
	Status uint8
	TXPowerLevel uint8
}

func (o *LEReadAdvertisingPhysicalChannelTxPowerOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.TXPowerLevel = uint8(r.GetOne())
	return r.Valid()
}

// LEReadAdvertisingPhysicalChannelTxPowerSync executes the command specified in Section 7.8.6 synchronously
func (c *Commands) LEReadAdvertisingPhysicalChannelTxPowerSync (result *LEReadAdvertisingPhysicalChannelTxPowerOutput) (*LEReadAdvertisingPhysicalChannelTxPowerOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadAdvertisingPhysicalChannelTxPower started")
	}
	if result == nil {
		result = &LEReadAdvertisingPhysicalChannelTxPowerOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0007}, nil)
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
		}).Debug("LEReadAdvertisingPhysicalChannelTxPower completed")
	}

	 return result, err
}
// LESetAdvertisingDataInput represents the input of the command specified in Section 7.8.7
type LESetAdvertisingDataInput struct {
	AdvertisingDataLength uint8
	AdvertisingData [31]byte
}

func (i LESetAdvertisingDataInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertisingDataLength))
	copy(w.Put(31), i.AdvertisingData[:])
	return w.Data
}

// LESetAdvertisingDataSync executes the command specified in Section 7.8.7 synchronously
func (c *Commands) LESetAdvertisingDataSync (params LESetAdvertisingDataInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetAdvertisingData started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0008}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "0params": params,
		}).Debug("LESetAdvertisingData completed")
	}

	 return err
}
// LESetScanResponseDataInput represents the input of the command specified in Section 7.8.8
type LESetScanResponseDataInput struct {
	ScanResponseDataLength uint8
	ScanResponseData [31]byte
}

func (i LESetScanResponseDataInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.ScanResponseDataLength))
	copy(w.Put(31), i.ScanResponseData[:])
	return w.Data
}

// LESetScanResponseDataSync executes the command specified in Section 7.8.8 synchronously
func (c *Commands) LESetScanResponseDataSync (params LESetScanResponseDataInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetScanResponseData started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0009}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "0params": params,
		}).Debug("LESetScanResponseData completed")
	}

	 return err
}
// LESetAdvertisingEnableInput represents the input of the command specified in Section 7.8.9
type LESetAdvertisingEnableInput struct {
	AdvertisingEnable uint8
}

func (i LESetAdvertisingEnableInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertisingEnable))
	return w.Data
}

// LESetAdvertisingEnableSync executes the command specified in Section 7.8.9 synchronously
func (c *Commands) LESetAdvertisingEnableSync (params LESetAdvertisingEnableInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetAdvertisingEnable started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000A}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
			 "0params": params,
		}).Debug("LESetAdvertisingEnable completed")
	}

	 return err
}
// LESetScanParametersInput represents the input of the command specified in Section 7.8.10
type LESetScanParametersInput struct {
	LEScanType uint8
	LEScanInterval uint16
	LEScanWindow uint16
	OwnAddressType bleutil.MacAddrType
	ScanningFilterPolicy uint8
}

func (i LESetScanParametersInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.LEScanType))
	binary.LittleEndian.PutUint16(w.Put(2), i.LEScanInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.LEScanWindow)
	w.PutOne(uint8(i.OwnAddressType))
	w.PutOne(uint8(i.ScanningFilterPolicy))
	return w.Data
}

// LESetScanParametersSync executes the command specified in Section 7.8.10 synchronously
func (c *Commands) LESetScanParametersSync (params LESetScanParametersInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetScanParameters started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000B}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetScanParameters completed")
	}

	 return err
}
// LESetScanEnableInput represents the input of the command specified in Section 7.8.11
type LESetScanEnableInput struct {
	LEScanEnable uint8
	FilterDuplicates uint8
}

func (i LESetScanEnableInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.LEScanEnable))
	w.PutOne(uint8(i.FilterDuplicates))
	return w.Data
}

// LESetScanEnableSync executes the command specified in Section 7.8.11 synchronously
func (c *Commands) LESetScanEnableSync (params LESetScanEnableInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetScanEnable started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000C}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetScanEnable completed")
	}

	 return err
}
// LECreateConnectionInput represents the input of the command specified in Section 7.8.12
type LECreateConnectionInput struct {
	LEScanInterval uint16
	LEScanWindow uint16
	InitiatorFilterPolicy uint8
	PeerAddressType bleutil.MacAddrType
	PeerAddress bleutil.MacAddr
	OwnAddressType bleutil.MacAddrType
	ConnectionIntervalMin uint16
	ConnectionIntervalMax uint16
	ConnectionLatency uint16
	SupervisionTimeout uint16
	MinCELength uint16
	MaxCELength uint16
}

func (i LECreateConnectionInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.LEScanInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.LEScanWindow)
	w.PutOne(uint8(i.InitiatorFilterPolicy))
	w.PutOne(uint8(i.PeerAddressType))
	i.PeerAddress.Encode(w.Put(6))
	w.PutOne(uint8(i.OwnAddressType))
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionIntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionIntervalMax)
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionLatency)
	binary.LittleEndian.PutUint16(w.Put(2), i.SupervisionTimeout)
	binary.LittleEndian.PutUint16(w.Put(2), i.MinCELength)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxCELength)
	return w.Data
}

// LECreateConnectionSync executes the command specified in Section 7.8.12 synchronously
func (c *Commands) LECreateConnectionSync (params LECreateConnectionInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LECreateConnection started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000D}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LECreateConnection completed")
	}

	 return err
}
// LECreateConnectionCancelSync executes the command specified in Section 7.8.13 synchronously
func (c *Commands) LECreateConnectionCancelSync () error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LECreateConnectionCancel started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000E}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
		}).Debug("LECreateConnectionCancel completed")
	}

	 return err
}
// LEReadWhiteListSizeOutput represents the output of the command specified in Section 7.8.14
type LEReadWhiteListSizeOutput struct {
	Status uint8
	WhiteListSize uint8
}

func (o *LEReadWhiteListSizeOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.WhiteListSize = uint8(r.GetOne())
	return r.Valid()
}

// LEReadWhiteListSizeSync executes the command specified in Section 7.8.14 synchronously
func (c *Commands) LEReadWhiteListSizeSync (result *LEReadWhiteListSizeOutput) (*LEReadWhiteListSizeOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadWhiteListSize started")
	}
	if result == nil {
		result = &LEReadWhiteListSizeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x000F}, nil)
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
		}).Debug("LEReadWhiteListSize completed")
	}

	 return result, err
}
// LEClearWhiteListSync executes the command specified in Section 7.8.15 synchronously
func (c *Commands) LEClearWhiteListSync () error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEClearWhiteList started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0010}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
		}).Debug("LEClearWhiteList completed")
	}

	 return err
}
// LEAddDeviceToWhiteListInput represents the input of the command specified in Section 7.8.16
type LEAddDeviceToWhiteListInput struct {
	AddressType bleutil.MacAddrType
	Address bleutil.MacAddr
}

func (i LEAddDeviceToWhiteListInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AddressType))
	i.Address.Encode(w.Put(6))
	return w.Data
}

// LEAddDeviceToWhiteListSync executes the command specified in Section 7.8.16 synchronously
func (c *Commands) LEAddDeviceToWhiteListSync (params LEAddDeviceToWhiteListInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEAddDeviceToWhiteList started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0011}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEAddDeviceToWhiteList completed")
	}

	 return err
}
// LERemoveDeviceFromWhiteListInput represents the input of the command specified in Section 7.8.17
type LERemoveDeviceFromWhiteListInput struct {
	AddressType bleutil.MacAddrType
	Address bleutil.MacAddr
}

func (i LERemoveDeviceFromWhiteListInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AddressType))
	i.Address.Encode(w.Put(6))
	return w.Data
}

// LERemoveDeviceFromWhiteListSync executes the command specified in Section 7.8.17 synchronously
func (c *Commands) LERemoveDeviceFromWhiteListSync (params LERemoveDeviceFromWhiteListInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LERemoveDeviceFromWhiteList started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0012}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LERemoveDeviceFromWhiteList completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionIntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionIntervalMax)
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionLatency)
	binary.LittleEndian.PutUint16(w.Put(2), i.SupervisionTimeout)
	binary.LittleEndian.PutUint16(w.Put(2), i.MinCELength)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxCELength)
	return w.Data
}

// LEConnectionUpdateSync executes the command specified in Section 7.8.18 synchronously
func (c *Commands) LEConnectionUpdateSync (params LEConnectionUpdateInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEConnectionUpdate started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0013}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEConnectionUpdate completed")
	}

	 return err
}
// LESetHostChannelClassificationInput represents the input of the command specified in Section 7.8.19
type LESetHostChannelClassificationInput struct {
	ChannelMap [5]byte
}

func (i LESetHostChannelClassificationInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	copy(w.Put(5), i.ChannelMap[:])
	return w.Data
}

// LESetHostChannelClassificationSync executes the command specified in Section 7.8.19 synchronously
func (c *Commands) LESetHostChannelClassificationSync (params LESetHostChannelClassificationInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetHostChannelClassification started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0014}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetHostChannelClassification completed")
	}

	 return err
}
// LEReadChannelMapInput represents the input of the command specified in Section 7.8.20
type LEReadChannelMapInput struct {
	ConnectionHandle uint16
}

func (i LEReadChannelMapInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data
}

// LEReadChannelMapOutput represents the output of the command specified in Section 7.8.20
type LEReadChannelMapOutput struct {
	Status uint8
	ConnectionHandle uint16
	ChannelMap [5]byte
}

func (o *LEReadChannelMapOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	copy(o.ChannelMap[:], r.Get(5))
	return r.Valid()
}

// LEReadChannelMapSync executes the command specified in Section 7.8.20 synchronously
func (c *Commands) LEReadChannelMapSync (params LEReadChannelMapInput, result *LEReadChannelMapOutput) (*LEReadChannelMapOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReadChannelMap started")
	}
	if result == nil {
		result = &LEReadChannelMapOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0015}, nil)
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
		}).Debug("LEReadChannelMap completed")
	}

	 return result, err
}
// LEReadRemoteFeaturesInput represents the input of the command specified in Section 7.8.21
type LEReadRemoteFeaturesInput struct {
	ConnectionHandle uint16
}

func (i LEReadRemoteFeaturesInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data
}

// LEReadRemoteFeaturesSync executes the command specified in Section 7.8.21 synchronously
func (c *Commands) LEReadRemoteFeaturesSync (params LEReadRemoteFeaturesInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReadRemoteFeatures started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0016}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEReadRemoteFeatures completed")
	}

	 return err
}
// LEEncryptInput represents the input of the command specified in Section 7.8.22
type LEEncryptInput struct {
	Key [16]byte
	PlaintextData [16]byte
}

func (i LEEncryptInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	copy(w.Put(16), i.Key[:])
	copy(w.Put(16), i.PlaintextData[:])
	return w.Data
}

// LEEncryptOutput represents the output of the command specified in Section 7.8.22
type LEEncryptOutput struct {
	Status uint8
	EncryptedData [16]byte
}

func (o *LEEncryptOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	copy(o.EncryptedData[:], r.Get(16))
	return r.Valid()
}

// LEEncryptSync executes the command specified in Section 7.8.22 synchronously
func (c *Commands) LEEncryptSync (params LEEncryptInput, result *LEEncryptOutput) (*LEEncryptOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEEncrypt started")
	}
	if result == nil {
		result = &LEEncryptOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0017}, nil)
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
		}).Debug("LEEncrypt completed")
	}

	 return result, err
}
// LERandOutput represents the output of the command specified in Section 7.8.23
type LERandOutput struct {
	Status uint8
	RandomNumber uint64
}

func (o *LERandOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.RandomNumber = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// LERandSync executes the command specified in Section 7.8.23 synchronously
func (c *Commands) LERandSync (result *LERandOutput) (*LERandOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LERand started")
	}
	if result == nil {
		result = &LERandOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0018}, nil)
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
		}).Debug("LERand completed")
	}

	 return result, err
}
// LEEnableEncryptionInput represents the input of the command specified in Section 7.8.24
type LEEnableEncryptionInput struct {
	ConnectionHandle uint16
	RandomNumber uint64
	EncryptedDiversifier uint16
	LongTermKey [16]byte
}

func (i LEEnableEncryptionInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint64(w.Put(8), i.RandomNumber)
	binary.LittleEndian.PutUint16(w.Put(2), i.EncryptedDiversifier)
	copy(w.Put(16), i.LongTermKey[:])
	return w.Data
}

// LEEnableEncryptionSync executes the command specified in Section 7.8.24 synchronously
func (c *Commands) LEEnableEncryptionSync (params LEEnableEncryptionInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEEnableEncryption started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0019}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEEnableEncryption completed")
	}

	 return err
}
// LELongTermKeyRequestReplyInput represents the input of the command specified in Section 7.8.25
type LELongTermKeyRequestReplyInput struct {
	ConnectionHandle uint16
	LongTermKey [16]byte
}

func (i LELongTermKeyRequestReplyInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	copy(w.Put(16), i.LongTermKey[:])
	return w.Data
}

// LELongTermKeyRequestReplyOutput represents the output of the command specified in Section 7.8.25
type LELongTermKeyRequestReplyOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LELongTermKeyRequestReplyOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LELongTermKeyRequestReplySync executes the command specified in Section 7.8.25 synchronously
func (c *Commands) LELongTermKeyRequestReplySync (params LELongTermKeyRequestReplyInput, result *LELongTermKeyRequestReplyOutput) (*LELongTermKeyRequestReplyOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LELongTermKeyRequestReply started")
	}
	if result == nil {
		result = &LELongTermKeyRequestReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001A}, nil)
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
		}).Debug("LELongTermKeyRequestReply completed")
	}

	 return result, err
}
// LELongTermKeyRequestNegativeReplyInput represents the input of the command specified in Section 7.8.26
type LELongTermKeyRequestNegativeReplyInput struct {
	ConnectionHandle uint16
}

func (i LELongTermKeyRequestNegativeReplyInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data
}

// LELongTermKeyRequestNegativeReplyOutput represents the output of the command specified in Section 7.8.26
type LELongTermKeyRequestNegativeReplyOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LELongTermKeyRequestNegativeReplyOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LELongTermKeyRequestNegativeReplySync executes the command specified in Section 7.8.26 synchronously
func (c *Commands) LELongTermKeyRequestNegativeReplySync (params LELongTermKeyRequestNegativeReplyInput, result *LELongTermKeyRequestNegativeReplyOutput) (*LELongTermKeyRequestNegativeReplyOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LELongTermKeyRequestNegativeReply started")
	}
	if result == nil {
		result = &LELongTermKeyRequestNegativeReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001B}, nil)
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
		}).Debug("LELongTermKeyRequestNegativeReply completed")
	}

	 return result, err
}
// LEReadSupportedStatesOutput represents the output of the command specified in Section 7.8.27
type LEReadSupportedStatesOutput struct {
	Status uint8
	LEStates uint64
}

func (o *LEReadSupportedStatesOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.LEStates = binary.LittleEndian.Uint64(r.Get(8))
	return r.Valid()
}

// LEReadSupportedStatesSync executes the command specified in Section 7.8.27 synchronously
func (c *Commands) LEReadSupportedStatesSync (result *LEReadSupportedStatesOutput) (*LEReadSupportedStatesOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadSupportedStates started")
	}
	if result == nil {
		result = &LEReadSupportedStatesOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001C}, nil)
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
		}).Debug("LEReadSupportedStates completed")
	}

	 return result, err
}
// LEReceiverTestInput represents the input of the command specified in Section 7.8.28
type LEReceiverTestInput struct {
	RXChannel uint8
}

func (i LEReceiverTestInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.RXChannel))
	return w.Data
}

// LEReceiverTestSync executes the command specified in Section 7.8.28 synchronously
func (c *Commands) LEReceiverTestSync (params LEReceiverTestInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReceiverTest started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001D}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEReceiverTest completed")
	}

	 return err
}
// LEReceiverTestV2Input represents the input of the command specified in Section 7.8.28
type LEReceiverTestV2Input struct {
	RXChannel uint8
	PHY uint8
}

func (i LEReceiverTestV2Input) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.RXChannel))
	w.PutOne(uint8(i.PHY))
	return w.Data
}

// LEReceiverTestV2Sync executes the command specified in Section 7.8.28 synchronously
func (c *Commands) LEReceiverTestV2Sync (params LEReceiverTestV2Input) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReceiverTestV2 started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0033}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEReceiverTestV2 completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.RXChannel))
	w.PutOne(uint8(i.PHY))
	w.PutOne(uint8(i.ModulationIndex))
	w.PutOne(uint8(i.ExpectedCTELength))
	w.PutOne(uint8(i.ExpectedCTEType))
	w.PutOne(uint8(i.SlotDurations))
	w.PutOne(uint8(i.SwitchingPatternLength))
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(uint8(m))
	}
	return w.Data
}

// LEReceiverTestV3Sync executes the command specified in Section 7.8.28 synchronously
func (c *Commands) LEReceiverTestV3Sync (params LEReceiverTestV3Input) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReceiverTestV3 started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004F}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEReceiverTestV3 completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.TXChannel))
	w.PutOne(uint8(i.TestDataLength))
	w.PutOne(uint8(i.PacketPayload))
	return w.Data
}

// LETransmitterTestSync executes the command specified in Section 7.8.29 synchronously
func (c *Commands) LETransmitterTestSync (params LETransmitterTestInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LETransmitterTest started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001E}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LETransmitterTest completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.TXChannel))
	w.PutOne(uint8(i.TestDataLength))
	w.PutOne(uint8(i.PacketPayload))
	w.PutOne(uint8(i.PHY))
	return w.Data
}

// LETransmitterTestV2Sync executes the command specified in Section 7.8.29 synchronously
func (c *Commands) LETransmitterTestV2Sync (params LETransmitterTestV2Input) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LETransmitterTestV2 started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0034}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LETransmitterTestV2 completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.TXChannel))
	w.PutOne(uint8(i.TestDataLength))
	w.PutOne(uint8(i.PacketPayload))
	w.PutOne(uint8(i.PHY))
	w.PutOne(uint8(i.CTELength))
	w.PutOne(uint8(i.CTEType))
	w.PutOne(uint8(i.SwitchingPatternLength))
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(uint8(m))
	}
	return w.Data
}

// LETransmitterTestV3Sync executes the command specified in Section 7.8.29 synchronously
func (c *Commands) LETransmitterTestV3Sync (params LETransmitterTestV3Input) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LETransmitterTestV3 started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0050}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LETransmitterTestV3 completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.TXChannel))
	w.PutOne(uint8(i.TestDataLength))
	w.PutOne(uint8(i.PacketPayload))
	w.PutOne(uint8(i.PHY))
	w.PutOne(uint8(i.CTELength))
	w.PutOne(uint8(i.CTEType))
	w.PutOne(uint8(i.SwitchingPatternLength))
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(uint8(m))
	}
	w.PutOne(uint8(i.TransmitPowerLevel))
	return w.Data
}

// LETransmitterTestV4Sync executes the command specified in Section 7.8.29 synchronously
func (c *Commands) LETransmitterTestV4Sync (params LETransmitterTestV4Input) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LETransmitterTestV4 started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x007B}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LETransmitterTestV4 completed")
	}

	 return err
}
// LETestEndOutput represents the output of the command specified in Section 7.8.30
type LETestEndOutput struct {
	Status uint8
	NumPackets uint16
}

func (o *LETestEndOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.NumPackets = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LETestEndSync executes the command specified in Section 7.8.30 synchronously
func (c *Commands) LETestEndSync (result *LETestEndOutput) (*LETestEndOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LETestEnd started")
	}
	if result == nil {
		result = &LETestEndOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x001F}, nil)
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
		}).Debug("LETestEnd completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.IntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.IntervalMax)
	binary.LittleEndian.PutUint16(w.Put(2), i.Latency)
	binary.LittleEndian.PutUint16(w.Put(2), i.Timeout)
	binary.LittleEndian.PutUint16(w.Put(2), i.MinCELength)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxCELength)
	return w.Data
}

// LERemoteConnectionParameterRequestReplyOutput represents the output of the command specified in Section 7.8.31
type LERemoteConnectionParameterRequestReplyOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LERemoteConnectionParameterRequestReplyOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LERemoteConnectionParameterRequestReplySync executes the command specified in Section 7.8.31 synchronously
func (c *Commands) LERemoteConnectionParameterRequestReplySync (params LERemoteConnectionParameterRequestReplyInput, result *LERemoteConnectionParameterRequestReplyOutput) (*LERemoteConnectionParameterRequestReplyOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LERemoteConnectionParameterRequestReply started")
	}
	if result == nil {
		result = &LERemoteConnectionParameterRequestReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0020}, nil)
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
		}).Debug("LERemoteConnectionParameterRequestReply completed")
	}

	 return result, err
}
// LERemoteConnectionParameterRequestNegativeReplyInput represents the input of the command specified in Section 7.8.32
type LERemoteConnectionParameterRequestNegativeReplyInput struct {
	ConnectionHandle uint16
	Reason uint8
}

func (i LERemoteConnectionParameterRequestNegativeReplyInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.Reason))
	return w.Data
}

// LERemoteConnectionParameterRequestNegativeReplyOutput represents the output of the command specified in Section 7.8.32
type LERemoteConnectionParameterRequestNegativeReplyOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LERemoteConnectionParameterRequestNegativeReplyOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LERemoteConnectionParameterRequestNegativeReplySync executes the command specified in Section 7.8.32 synchronously
func (c *Commands) LERemoteConnectionParameterRequestNegativeReplySync (params LERemoteConnectionParameterRequestNegativeReplyInput, result *LERemoteConnectionParameterRequestNegativeReplyOutput) (*LERemoteConnectionParameterRequestNegativeReplyOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LERemoteConnectionParameterRequestNegativeReply started")
	}
	if result == nil {
		result = &LERemoteConnectionParameterRequestNegativeReplyOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0021}, nil)
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
		}).Debug("LERemoteConnectionParameterRequestNegativeReply completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.TXOctets)
	binary.LittleEndian.PutUint16(w.Put(2), i.TXTime)
	return w.Data
}

// LESetDataLengthOutput represents the output of the command specified in Section 7.8.33
type LESetDataLengthOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetDataLengthOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetDataLengthSync executes the command specified in Section 7.8.33 synchronously
func (c *Commands) LESetDataLengthSync (params LESetDataLengthInput, result *LESetDataLengthOutput) (*LESetDataLengthOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetDataLength started")
	}
	if result == nil {
		result = &LESetDataLengthOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0022}, nil)
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
		}).Debug("LESetDataLength completed")
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
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.SuggestedMaxTXOctets = binary.LittleEndian.Uint16(r.Get(2))
	o.SuggestedMaxTXTime = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEReadSuggestedDefaultDataLengthSync executes the command specified in Section 7.8.34 synchronously
func (c *Commands) LEReadSuggestedDefaultDataLengthSync (result *LEReadSuggestedDefaultDataLengthOutput) (*LEReadSuggestedDefaultDataLengthOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadSuggestedDefaultDataLength started")
	}
	if result == nil {
		result = &LEReadSuggestedDefaultDataLengthOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0023}, nil)
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
		}).Debug("LEReadSuggestedDefaultDataLength completed")
	}

	 return result, err
}
// LEWriteSuggestedDefaultDataLengthInput represents the input of the command specified in Section 7.8.35
type LEWriteSuggestedDefaultDataLengthInput struct {
	SuggestedMaxTXOctets uint16
	SuggestedMaxTXTime uint16
}

func (i LEWriteSuggestedDefaultDataLengthInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.SuggestedMaxTXOctets)
	binary.LittleEndian.PutUint16(w.Put(2), i.SuggestedMaxTXTime)
	return w.Data
}

// LEWriteSuggestedDefaultDataLengthSync executes the command specified in Section 7.8.35 synchronously
func (c *Commands) LEWriteSuggestedDefaultDataLengthSync (params LEWriteSuggestedDefaultDataLengthInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEWriteSuggestedDefaultDataLength started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0024}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEWriteSuggestedDefaultDataLength completed")
	}

	 return err
}
// LEReadLocalP256PublicKeyInput represents the input of the command specified in Section 7.8.36
type LEReadLocalP256PublicKeyInput struct {
	RemoteP256PublicKey [64]byte
}

func (i LEReadLocalP256PublicKeyInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	copy(w.Put(64), i.RemoteP256PublicKey[:])
	return w.Data
}

// LEReadLocalP256PublicKeySync executes the command specified in Section 7.8.36 synchronously
func (c *Commands) LEReadLocalP256PublicKeySync (params LEReadLocalP256PublicKeyInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReadLocalP256PublicKey started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0026}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEReadLocalP256PublicKey completed")
	}

	 return err
}
// LEGenerateDHKeyV2Input represents the input of the command specified in Section 7.8.37
type LEGenerateDHKeyV2Input struct {
	RemoteP256PublicKey [64]byte
	KeyType uint8
}

func (i LEGenerateDHKeyV2Input) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	copy(w.Put(64), i.RemoteP256PublicKey[:])
	w.PutOne(uint8(i.KeyType))
	return w.Data
}

// LEGenerateDHKeyV2Sync executes the command specified in Section 7.8.37 synchronously
func (c *Commands) LEGenerateDHKeyV2Sync (params LEGenerateDHKeyV2Input) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEGenerateDHKeyV2 started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x005E}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEGenerateDHKeyV2 completed")
	}

	 return err
}
// LEAddDeviceToResolvingListInput represents the input of the command specified in Section 7.8.38
type LEAddDeviceToResolvingListInput struct {
	PeerIdentityAddressType bleutil.MacAddrType
	PeerIdentityAddress bleutil.MacAddr
	PeerIRK [16]byte
	LocalIRK [16]byte
}

func (i LEAddDeviceToResolvingListInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.PeerIdentityAddressType))
	i.PeerIdentityAddress.Encode(w.Put(6))
	copy(w.Put(16), i.PeerIRK[:])
	copy(w.Put(16), i.LocalIRK[:])
	return w.Data
}

// LEAddDeviceToResolvingListSync executes the command specified in Section 7.8.38 synchronously
func (c *Commands) LEAddDeviceToResolvingListSync (params LEAddDeviceToResolvingListInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEAddDeviceToResolvingList started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0027}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEAddDeviceToResolvingList completed")
	}

	 return err
}
// LERemoveDeviceFromResolvingListInput represents the input of the command specified in Section 7.8.39
type LERemoveDeviceFromResolvingListInput struct {
	PeerIdentityAddressType bleutil.MacAddrType
	PeerDeviceAddress bleutil.MacAddr
}

func (i LERemoveDeviceFromResolvingListInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.PeerIdentityAddressType))
	i.PeerDeviceAddress.Encode(w.Put(6))
	return w.Data
}

// LERemoveDeviceFromResolvingListSync executes the command specified in Section 7.8.39 synchronously
func (c *Commands) LERemoveDeviceFromResolvingListSync (params LERemoveDeviceFromResolvingListInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LERemoveDeviceFromResolvingList started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0028}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LERemoveDeviceFromResolvingList completed")
	}

	 return err
}
// LEClearResolvingListSync executes the command specified in Section 7.8.40 synchronously
func (c *Commands) LEClearResolvingListSync () error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEClearResolvingList started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0029}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
		}).Debug("LEClearResolvingList completed")
	}

	 return err
}
// LEReadResolvingListSizeOutput represents the output of the command specified in Section 7.8.41
type LEReadResolvingListSizeOutput struct {
	Status uint8
	ResolvingListSize uint8
}

func (o *LEReadResolvingListSizeOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ResolvingListSize = uint8(r.GetOne())
	return r.Valid()
}

// LEReadResolvingListSizeSync executes the command specified in Section 7.8.41 synchronously
func (c *Commands) LEReadResolvingListSizeSync (result *LEReadResolvingListSizeOutput) (*LEReadResolvingListSizeOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadResolvingListSize started")
	}
	if result == nil {
		result = &LEReadResolvingListSizeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002A}, nil)
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
		}).Debug("LEReadResolvingListSize completed")
	}

	 return result, err
}
// LEReadPeerResolvableAddressInput represents the input of the command specified in Section 7.8.42
type LEReadPeerResolvableAddressInput struct {
	PeerIdentityAddressType bleutil.MacAddrType
	PeerIdentityAddress bleutil.MacAddr
}

func (i LEReadPeerResolvableAddressInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.PeerIdentityAddressType))
	i.PeerIdentityAddress.Encode(w.Put(6))
	return w.Data
}

// LEReadPeerResolvableAddressOutput represents the output of the command specified in Section 7.8.42
type LEReadPeerResolvableAddressOutput struct {
	Status uint8
	PeerResolvableAddress bleutil.MacAddr
}

func (o *LEReadPeerResolvableAddressOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.PeerResolvableAddress.Decode(r.Get(6))
	return r.Valid()
}

// LEReadPeerResolvableAddressSync executes the command specified in Section 7.8.42 synchronously
func (c *Commands) LEReadPeerResolvableAddressSync (params LEReadPeerResolvableAddressInput, result *LEReadPeerResolvableAddressOutput) (*LEReadPeerResolvableAddressOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReadPeerResolvableAddress started")
	}
	if result == nil {
		result = &LEReadPeerResolvableAddressOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002B}, nil)
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
		}).Debug("LEReadPeerResolvableAddress completed")
	}

	 return result, err
}
// LEReadLocalResolvableAddressInput represents the input of the command specified in Section 7.8.43
type LEReadLocalResolvableAddressInput struct {
	PeerIdentityAddressType bleutil.MacAddrType
	PeerIdentityAddress bleutil.MacAddr
}

func (i LEReadLocalResolvableAddressInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.PeerIdentityAddressType))
	i.PeerIdentityAddress.Encode(w.Put(6))
	return w.Data
}

// LEReadLocalResolvableAddressOutput represents the output of the command specified in Section 7.8.43
type LEReadLocalResolvableAddressOutput struct {
	Status uint8
	LocalResolvableAddress bleutil.MacAddr
}

func (o *LEReadLocalResolvableAddressOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.LocalResolvableAddress.Decode(r.Get(6))
	return r.Valid()
}

// LEReadLocalResolvableAddressSync executes the command specified in Section 7.8.43 synchronously
func (c *Commands) LEReadLocalResolvableAddressSync (params LEReadLocalResolvableAddressInput, result *LEReadLocalResolvableAddressOutput) (*LEReadLocalResolvableAddressOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReadLocalResolvableAddress started")
	}
	if result == nil {
		result = &LEReadLocalResolvableAddressOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002C}, nil)
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
		}).Debug("LEReadLocalResolvableAddress completed")
	}

	 return result, err
}
// LESetAddressResolutionEnableInput represents the input of the command specified in Section 7.8.44
type LESetAddressResolutionEnableInput struct {
	AddressResolutionEnable uint8
}

func (i LESetAddressResolutionEnableInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AddressResolutionEnable))
	return w.Data
}

// LESetAddressResolutionEnableSync executes the command specified in Section 7.8.44 synchronously
func (c *Commands) LESetAddressResolutionEnableSync (params LESetAddressResolutionEnableInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetAddressResolutionEnable started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002D}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetAddressResolutionEnable completed")
	}

	 return err
}
// LESetResolvablePrivateAddressTimeoutInput represents the input of the command specified in Section 7.8.45
type LESetResolvablePrivateAddressTimeoutInput struct {
	RPATimeout uint16
}

func (i LESetResolvablePrivateAddressTimeoutInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.RPATimeout)
	return w.Data
}

// LESetResolvablePrivateAddressTimeoutSync executes the command specified in Section 7.8.45 synchronously
func (c *Commands) LESetResolvablePrivateAddressTimeoutSync (params LESetResolvablePrivateAddressTimeoutInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetResolvablePrivateAddressTimeout started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002E}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetResolvablePrivateAddressTimeout completed")
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
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.SupportedMaxTXOctets = binary.LittleEndian.Uint16(r.Get(2))
	o.SupportedMaxTXTime = binary.LittleEndian.Uint16(r.Get(2))
	o.SupportedMaxRXOctets = binary.LittleEndian.Uint16(r.Get(2))
	o.SupportedMaxRXTime = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEReadMaximumDataLengthSync executes the command specified in Section 7.8.46 synchronously
func (c *Commands) LEReadMaximumDataLengthSync (result *LEReadMaximumDataLengthOutput) (*LEReadMaximumDataLengthOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadMaximumDataLength started")
	}
	if result == nil {
		result = &LEReadMaximumDataLengthOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x002F}, nil)
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
		}).Debug("LEReadMaximumDataLength completed")
	}

	 return result, err
}
// LEReadPHYInput represents the input of the command specified in Section 7.8.47
type LEReadPHYInput struct {
	ConnectionHandle uint16
}

func (i LEReadPHYInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data
}

// LEReadPHYOutput represents the output of the command specified in Section 7.8.47
type LEReadPHYOutput struct {
	Status uint8
	ConnectionHandle uint16
	TXPHY uint8
	RXPHY uint8
}

func (o *LEReadPHYOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.TXPHY = uint8(r.GetOne())
	o.RXPHY = uint8(r.GetOne())
	return r.Valid()
}

// LEReadPHYSync executes the command specified in Section 7.8.47 synchronously
func (c *Commands) LEReadPHYSync (params LEReadPHYInput, result *LEReadPHYOutput) (*LEReadPHYOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReadPHY started")
	}
	if result == nil {
		result = &LEReadPHYOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0030}, nil)
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
		}).Debug("LEReadPHY completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AllPHYs))
	w.PutOne(uint8(i.TXPHYs))
	w.PutOne(uint8(i.RXPHYs))
	return w.Data
}

// LESetDefaultPHYSync executes the command specified in Section 7.8.48 synchronously
func (c *Commands) LESetDefaultPHYSync (params LESetDefaultPHYInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetDefaultPHY started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0031}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetDefaultPHY completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.AllPHYs))
	w.PutOne(uint8(i.TXPHYs))
	w.PutOne(uint8(i.RXPHYs))
	binary.LittleEndian.PutUint16(w.Put(2), i.PHYOptions)
	return w.Data
}

// LESetPHYSync executes the command specified in Section 7.8.49 synchronously
func (c *Commands) LESetPHYSync (params LESetPHYInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetPHY started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0032}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetPHY completed")
	}

	 return err
}
// LESetAdvertisingSetRandomAddressInput represents the input of the command specified in Section 7.8.52
type LESetAdvertisingSetRandomAddressInput struct {
	AdvertisingHandle uint8
	AdvertisingRandomAddress bleutil.MacAddr
}

func (i LESetAdvertisingSetRandomAddressInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertisingHandle))
	i.AdvertisingRandomAddress.Encode(w.Put(6))
	return w.Data
}

// LESetAdvertisingSetRandomAddressSync executes the command specified in Section 7.8.52 synchronously
func (c *Commands) LESetAdvertisingSetRandomAddressSync (params LESetAdvertisingSetRandomAddressInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetAdvertisingSetRandomAddress started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0035}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetAdvertisingSetRandomAddress completed")
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
	OwnAddressType bleutil.MacAddrType
	PeerAddressType bleutil.MacAddrType
	PeerAddress bleutil.MacAddr
	AdvertisingFilterPolicy uint8
	AdvertisingTXPower uint8
	PrimaryAdvertisingPHY uint8
	SecondaryAdvertisingMaxSkip uint8
	SecondaryAdvertisingPHY uint8
	AdvertisingSID uint8
	ScanRequestNotificationEnable uint8
}

func (i LESetExtendedAdvertisingParametersInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertisingHandle))
	binary.LittleEndian.PutUint16(w.Put(2), i.AdvertisingEventProperties)
	bleutil.EncodeUint24(w.Put(3), i.PrimaryAdvertisingIntervalMin)
	bleutil.EncodeUint24(w.Put(3), i.PrimaryAdvertisingIntervalMax)
	w.PutOne(uint8(i.PrimaryAdvertisingChannelMap))
	w.PutOne(uint8(i.OwnAddressType))
	w.PutOne(uint8(i.PeerAddressType))
	i.PeerAddress.Encode(w.Put(6))
	w.PutOne(uint8(i.AdvertisingFilterPolicy))
	w.PutOne(uint8(i.AdvertisingTXPower))
	w.PutOne(uint8(i.PrimaryAdvertisingPHY))
	w.PutOne(uint8(i.SecondaryAdvertisingMaxSkip))
	w.PutOne(uint8(i.SecondaryAdvertisingPHY))
	w.PutOne(uint8(i.AdvertisingSID))
	w.PutOne(uint8(i.ScanRequestNotificationEnable))
	return w.Data
}

// LESetExtendedAdvertisingParametersOutput represents the output of the command specified in Section 7.8.53
type LESetExtendedAdvertisingParametersOutput struct {
	Status uint8
	SelectedTXPower uint8
}

func (o *LESetExtendedAdvertisingParametersOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.SelectedTXPower = uint8(r.GetOne())
	return r.Valid()
}

// LESetExtendedAdvertisingParametersSync executes the command specified in Section 7.8.53 synchronously
func (c *Commands) LESetExtendedAdvertisingParametersSync (params LESetExtendedAdvertisingParametersInput, result *LESetExtendedAdvertisingParametersOutput) (*LESetExtendedAdvertisingParametersOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetExtendedAdvertisingParameters started")
	}
	if result == nil {
		result = &LESetExtendedAdvertisingParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0036}, nil)
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
		}).Debug("LESetExtendedAdvertisingParameters completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertisingHandle))
	w.PutOne(uint8(i.Operation))
	w.PutOne(uint8(i.FragmentPreference))
	w.PutOne(uint8(i.AdvertisingDataLength))
	w.PutSlice(i.AdvertisingData)
	return w.Data
}

// LESetExtendedAdvertisingDataSync executes the command specified in Section 7.8.54 synchronously
func (c *Commands) LESetExtendedAdvertisingDataSync (params LESetExtendedAdvertisingDataInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetExtendedAdvertisingData started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0037}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetExtendedAdvertisingData completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertisingHandle))
	w.PutOne(uint8(i.Operation))
	w.PutOne(uint8(i.FragmentPreference))
	w.PutOne(uint8(i.ScanResponseDataLength))
	w.PutSlice(i.ScanResponseData)
	return w.Data
}

// LESetExtendedScanResponseDataSync executes the command specified in Section 7.8.55 synchronously
func (c *Commands) LESetExtendedScanResponseDataSync (params LESetExtendedScanResponseDataInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetExtendedScanResponseData started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0038}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetExtendedScanResponseData completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.Enable))
	w.PutOne(uint8(i.NumSets))
	if len(i.AdvertisingHandle) != int(i.NumSets) {
		panic("len(i.AdvertisingHandle) != int(i.NumSets)")
	}
	for _, m := range i.AdvertisingHandle {
		w.PutOne(uint8(m))
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
		w.PutOne(uint8(m))
	}
	return w.Data
}

// LESetExtendedAdvertisingEnableSync executes the command specified in Section 7.8.56 synchronously
func (c *Commands) LESetExtendedAdvertisingEnableSync (params LESetExtendedAdvertisingEnableInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetExtendedAdvertisingEnable started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0039}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetExtendedAdvertisingEnable completed")
	}

	 return err
}
// LEReadMaximumAdvertisingDataLengthOutput represents the output of the command specified in Section 7.8.57
type LEReadMaximumAdvertisingDataLengthOutput struct {
	Status uint8
	MaxAdvertisingDataLength uint16
}

func (o *LEReadMaximumAdvertisingDataLengthOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.MaxAdvertisingDataLength = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEReadMaximumAdvertisingDataLengthSync executes the command specified in Section 7.8.57 synchronously
func (c *Commands) LEReadMaximumAdvertisingDataLengthSync (result *LEReadMaximumAdvertisingDataLengthOutput) (*LEReadMaximumAdvertisingDataLengthOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadMaximumAdvertisingDataLength started")
	}
	if result == nil {
		result = &LEReadMaximumAdvertisingDataLengthOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003A}, nil)
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
		}).Debug("LEReadMaximumAdvertisingDataLength completed")
	}

	 return result, err
}
// LEReadNumberofSupportedAdvertisingSetsOutput represents the output of the command specified in Section 7.8.58
type LEReadNumberofSupportedAdvertisingSetsOutput struct {
	Status uint8
	NumSupportedAdvertisingSets uint8
}

func (o *LEReadNumberofSupportedAdvertisingSetsOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.NumSupportedAdvertisingSets = uint8(r.GetOne())
	return r.Valid()
}

// LEReadNumberofSupportedAdvertisingSetsSync executes the command specified in Section 7.8.58 synchronously
func (c *Commands) LEReadNumberofSupportedAdvertisingSetsSync (result *LEReadNumberofSupportedAdvertisingSetsOutput) (*LEReadNumberofSupportedAdvertisingSetsOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadNumberofSupportedAdvertisingSets started")
	}
	if result == nil {
		result = &LEReadNumberofSupportedAdvertisingSetsOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003B}, nil)
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
		}).Debug("LEReadNumberofSupportedAdvertisingSets completed")
	}

	 return result, err
}
// LERemoveAdvertisingSetInput represents the input of the command specified in Section 7.8.59
type LERemoveAdvertisingSetInput struct {
	AdvertisingHandle uint8
}

func (i LERemoveAdvertisingSetInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertisingHandle))
	return w.Data
}

// LERemoveAdvertisingSetSync executes the command specified in Section 7.8.59 synchronously
func (c *Commands) LERemoveAdvertisingSetSync (params LERemoveAdvertisingSetInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LERemoveAdvertisingSet started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003C}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LERemoveAdvertisingSet completed")
	}

	 return err
}
// LEClearAdvertisingSetsSync executes the command specified in Section 7.8.60 synchronously
func (c *Commands) LEClearAdvertisingSetsSync () error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEClearAdvertisingSets started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003D}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
		}).Debug("LEClearAdvertisingSets completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertisingHandle))
	binary.LittleEndian.PutUint16(w.Put(2), i.PeriodicAdvertisingIntervalMin)
	binary.LittleEndian.PutUint16(w.Put(2), i.PeriodicAdvertisingIntervalMax)
	binary.LittleEndian.PutUint16(w.Put(2), i.PeriodicAdvertisingProperties)
	return w.Data
}

// LESetPeriodicAdvertisingParametersSync executes the command specified in Section 7.8.61 synchronously
func (c *Commands) LESetPeriodicAdvertisingParametersSync (params LESetPeriodicAdvertisingParametersInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetPeriodicAdvertisingParameters started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003E}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetPeriodicAdvertisingParameters completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertisingHandle))
	w.PutOne(uint8(i.Operation))
	w.PutOne(uint8(i.AdvertisingDataLength))
	w.PutSlice(i.AdvertisingData)
	return w.Data
}

// LESetPeriodicAdvertisingDataSync executes the command specified in Section 7.8.62 synchronously
func (c *Commands) LESetPeriodicAdvertisingDataSync (params LESetPeriodicAdvertisingDataInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetPeriodicAdvertisingData started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x003F}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetPeriodicAdvertisingData completed")
	}

	 return err
}
// LESetPeriodicAdvertisingEnableInput represents the input of the command specified in Section 7.8.63
type LESetPeriodicAdvertisingEnableInput struct {
	Enable uint8
	AdvertisingHandle uint8
}

func (i LESetPeriodicAdvertisingEnableInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.Enable))
	w.PutOne(uint8(i.AdvertisingHandle))
	return w.Data
}

// LESetPeriodicAdvertisingEnableSync executes the command specified in Section 7.8.63 synchronously
func (c *Commands) LESetPeriodicAdvertisingEnableSync (params LESetPeriodicAdvertisingEnableInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetPeriodicAdvertisingEnable started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0040}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetPeriodicAdvertisingEnable completed")
	}

	 return err
}
// LESetExtendedScanParametersInput represents the input of the command specified in Section 7.8.64
type LESetExtendedScanParametersInput struct {
	OwnAddressType bleutil.MacAddrType
	ScanningFilterPolicy uint8
	ScanningPHYs uint8
	ScanType []uint8
	ScanInterval []uint16
	ScanWindow []uint16
}

func (i LESetExtendedScanParametersInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.OwnAddressType))
	w.PutOne(uint8(i.ScanningFilterPolicy))
	w.PutOne(uint8(i.ScanningPHYs))
	var0 := bleutil.CountSetBits(uint64(i.ScanningPHYs))
	if len(i.ScanType) != var0 {
		panic("len(i.ScanType) != var0")
	}
	for _, m := range i.ScanType {
		w.PutOne(uint8(m))
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
	return w.Data
}

// LESetExtendedScanParametersSync executes the command specified in Section 7.8.64 synchronously
func (c *Commands) LESetExtendedScanParametersSync (params LESetExtendedScanParametersInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetExtendedScanParameters started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0041}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetExtendedScanParameters completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.Enable))
	w.PutOne(uint8(i.FilterDuplicates))
	binary.LittleEndian.PutUint16(w.Put(2), i.Duration)
	binary.LittleEndian.PutUint16(w.Put(2), i.Period)
	return w.Data
}

// LESetExtendedScanEnableSync executes the command specified in Section 7.8.65 synchronously
func (c *Commands) LESetExtendedScanEnableSync (params LESetExtendedScanEnableInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetExtendedScanEnable started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0042}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetExtendedScanEnable completed")
	}

	 return err
}
// LEExtendedCreateConnectionInput represents the input of the command specified in Section 7.8.66
type LEExtendedCreateConnectionInput struct {
	InitiatingFilterPolicy uint8
	OwnAddressType bleutil.MacAddrType
	PeerAddressType bleutil.MacAddrType
	PeerAddress bleutil.MacAddr
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.InitiatingFilterPolicy))
	w.PutOne(uint8(i.OwnAddressType))
	w.PutOne(uint8(i.PeerAddressType))
	i.PeerAddress.Encode(w.Put(6))
	w.PutOne(uint8(i.InitiatingPHYs))
	var1 := bleutil.CountSetBits(uint64(i.InitiatingPHYs))
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
	return w.Data
}

// LEExtendedCreateConnectionSync executes the command specified in Section 7.8.66 synchronously
func (c *Commands) LEExtendedCreateConnectionSync (params LEExtendedCreateConnectionInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEExtendedCreateConnection started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0043}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEExtendedCreateConnection completed")
	}

	 return err
}
// LEPeriodicAdvertisingCreateSyncInput represents the input of the command specified in Section 7.8.67
type LEPeriodicAdvertisingCreateSyncInput struct {
	Options uint8
	AdvertisingSID uint8
	AdvertiserAddressType bleutil.MacAddrType
	AdvertiserAddress bleutil.MacAddr
	Skip uint16
	SyncTimeout uint16
	SyncCTEType uint8
}

func (i LEPeriodicAdvertisingCreateSyncInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.Options))
	w.PutOne(uint8(i.AdvertisingSID))
	w.PutOne(uint8(i.AdvertiserAddressType))
	i.AdvertiserAddress.Encode(w.Put(6))
	binary.LittleEndian.PutUint16(w.Put(2), i.Skip)
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncTimeout)
	w.PutOne(uint8(i.SyncCTEType))
	return w.Data
}

// LEPeriodicAdvertisingCreateSyncSync executes the command specified in Section 7.8.67 synchronously
func (c *Commands) LEPeriodicAdvertisingCreateSyncSync (params LEPeriodicAdvertisingCreateSyncInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEPeriodicAdvertisingCreateSync started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0044}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEPeriodicAdvertisingCreateSync completed")
	}

	 return err
}
// LEPeriodicAdvertisingCreateSyncCancelSync executes the command specified in Section 7.8.68 synchronously
func (c *Commands) LEPeriodicAdvertisingCreateSyncCancelSync () error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEPeriodicAdvertisingCreateSyncCancel started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0045}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
		}).Debug("LEPeriodicAdvertisingCreateSyncCancel completed")
	}

	 return err
}
// LEPeriodicAdvertisingTerminateSyncInput represents the input of the command specified in Section 7.8.69
type LEPeriodicAdvertisingTerminateSyncInput struct {
	SyncHandle uint16
}

func (i LEPeriodicAdvertisingTerminateSyncInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncHandle)
	return w.Data
}

// LEPeriodicAdvertisingTerminateSyncSync executes the command specified in Section 7.8.69 synchronously
func (c *Commands) LEPeriodicAdvertisingTerminateSyncSync (params LEPeriodicAdvertisingTerminateSyncInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEPeriodicAdvertisingTerminateSync started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0046}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEPeriodicAdvertisingTerminateSync completed")
	}

	 return err
}
// LEAddDeviceToPeriodicAdvertiserListInput represents the input of the command specified in Section 7.8.70
type LEAddDeviceToPeriodicAdvertiserListInput struct {
	AdvertiserAddressType bleutil.MacAddrType
	AdvertiserAddress bleutil.MacAddr
	AdvertisingSID uint8
}

func (i LEAddDeviceToPeriodicAdvertiserListInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertiserAddressType))
	i.AdvertiserAddress.Encode(w.Put(6))
	w.PutOne(uint8(i.AdvertisingSID))
	return w.Data
}

// LEAddDeviceToPeriodicAdvertiserListSync executes the command specified in Section 7.8.70 synchronously
func (c *Commands) LEAddDeviceToPeriodicAdvertiserListSync (params LEAddDeviceToPeriodicAdvertiserListInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEAddDeviceToPeriodicAdvertiserList started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0047}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEAddDeviceToPeriodicAdvertiserList completed")
	}

	 return err
}
// LERemoveDeviceFromPeriodicAdvertiserListInput represents the input of the command specified in Section 7.8.71
type LERemoveDeviceFromPeriodicAdvertiserListInput struct {
	AdvertiserAddressType bleutil.MacAddrType
	AdvertiserAddress bleutil.MacAddr
	AdvertisingSID uint8
}

func (i LERemoveDeviceFromPeriodicAdvertiserListInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertiserAddressType))
	i.AdvertiserAddress.Encode(w.Put(6))
	w.PutOne(uint8(i.AdvertisingSID))
	return w.Data
}

// LERemoveDeviceFromPeriodicAdvertiserListSync executes the command specified in Section 7.8.71 synchronously
func (c *Commands) LERemoveDeviceFromPeriodicAdvertiserListSync (params LERemoveDeviceFromPeriodicAdvertiserListInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LERemoveDeviceFromPeriodicAdvertiserList started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0048}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LERemoveDeviceFromPeriodicAdvertiserList completed")
	}

	 return err
}
// LEClearPeriodicAdvertiserListSync executes the command specified in Section 7.8.72 synchronously
func (c *Commands) LEClearPeriodicAdvertiserListSync () error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEClearPeriodicAdvertiserList started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0049}, nil)
	if err != nil {
		goto log
	}

	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
	}

	err = HciErrorToGo(response, err)

	err2 = c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
	if err2 != nil {
		err = err2
	}

log:
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		c.logger.WithError(err).WithFields(logrus.Fields{
		}).Debug("LEClearPeriodicAdvertiserList completed")
	}

	 return err
}
// LEReadPeriodicAdvertiserListSizeOutput represents the output of the command specified in Section 7.8.73
type LEReadPeriodicAdvertiserListSizeOutput struct {
	Status uint8
	PeriodicAdvertiserListSize uint8
}

func (o *LEReadPeriodicAdvertiserListSizeOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.PeriodicAdvertiserListSize = uint8(r.GetOne())
	return r.Valid()
}

// LEReadPeriodicAdvertiserListSizeSync executes the command specified in Section 7.8.73 synchronously
func (c *Commands) LEReadPeriodicAdvertiserListSizeSync (result *LEReadPeriodicAdvertiserListSizeOutput) (*LEReadPeriodicAdvertiserListSizeOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadPeriodicAdvertiserListSize started")
	}
	if result == nil {
		result = &LEReadPeriodicAdvertiserListSizeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004A}, nil)
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
		}).Debug("LEReadPeriodicAdvertiserListSize completed")
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
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.MinTXPower = uint8(r.GetOne())
	o.MaxTXPower = uint8(r.GetOne())
	return r.Valid()
}

// LEReadTransmitPowerSync executes the command specified in Section 7.8.74 synchronously
func (c *Commands) LEReadTransmitPowerSync (result *LEReadTransmitPowerOutput) (*LEReadTransmitPowerOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadTransmitPower started")
	}
	if result == nil {
		result = &LEReadTransmitPowerOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004B}, nil)
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
		}).Debug("LEReadTransmitPower completed")
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
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.RFTXPathCompensationValue = binary.LittleEndian.Uint16(r.Get(2))
	o.RFRXPathCompensationValue = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEReadRFPathCompensationSync executes the command specified in Section 7.8.75 synchronously
func (c *Commands) LEReadRFPathCompensationSync (result *LEReadRFPathCompensationOutput) (*LEReadRFPathCompensationOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadRFPathCompensation started")
	}
	if result == nil {
		result = &LEReadRFPathCompensationOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004C}, nil)
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
		}).Debug("LEReadRFPathCompensation completed")
	}

	 return result, err
}
// LEWriteRFPathCompensationInput represents the input of the command specified in Section 7.8.76
type LEWriteRFPathCompensationInput struct {
	RFTXPathCompensationValue uint16
	RFRXPathCompensationValue uint16
}

func (i LEWriteRFPathCompensationInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.RFTXPathCompensationValue)
	binary.LittleEndian.PutUint16(w.Put(2), i.RFRXPathCompensationValue)
	return w.Data
}

// LEWriteRFPathCompensationSync executes the command specified in Section 7.8.76 synchronously
func (c *Commands) LEWriteRFPathCompensationSync (params LEWriteRFPathCompensationInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEWriteRFPathCompensation started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004D}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEWriteRFPathCompensation completed")
	}

	 return err
}
// LESetPrivacyModeInput represents the input of the command specified in Section 7.8.77
type LESetPrivacyModeInput struct {
	PeerIdentityAddressType bleutil.MacAddrType
	PeerIdentityAddress bleutil.MacAddr
	PrivacyMode uint8
}

func (i LESetPrivacyModeInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.PeerIdentityAddressType))
	i.PeerIdentityAddress.Encode(w.Put(6))
	w.PutOne(uint8(i.PrivacyMode))
	return w.Data
}

// LESetPrivacyModeSync executes the command specified in Section 7.8.77 synchronously
func (c *Commands) LESetPrivacyModeSync (params LESetPrivacyModeInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetPrivacyMode started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x004E}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetPrivacyMode completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertisingHandle))
	w.PutOne(uint8(i.CTELength))
	w.PutOne(uint8(i.CTEType))
	w.PutOne(uint8(i.CTECount))
	w.PutOne(uint8(i.SwitchingPatternLength))
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(uint8(m))
	}
	return w.Data
}

// LESetConnectionlessCTETransmitParametersSync executes the command specified in Section 7.8.80 synchronously
func (c *Commands) LESetConnectionlessCTETransmitParametersSync (params LESetConnectionlessCTETransmitParametersInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetConnectionlessCTETransmitParameters started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0051}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetConnectionlessCTETransmitParameters completed")
	}

	 return err
}
// LESetConnectionlessCTETransmitEnableInput represents the input of the command specified in Section 7.8.81
type LESetConnectionlessCTETransmitEnableInput struct {
	AdvertisingHandle uint8
	CTEEnable uint8
}

func (i LESetConnectionlessCTETransmitEnableInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.AdvertisingHandle))
	w.PutOne(uint8(i.CTEEnable))
	return w.Data
}

// LESetConnectionlessCTETransmitEnableSync executes the command specified in Section 7.8.81 synchronously
func (c *Commands) LESetConnectionlessCTETransmitEnableSync (params LESetConnectionlessCTETransmitEnableInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetConnectionlessCTETransmitEnable started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0052}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetConnectionlessCTETransmitEnable completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncHandle)
	w.PutOne(uint8(i.SamplingEnable))
	w.PutOne(uint8(i.SlotDurations))
	w.PutOne(uint8(i.MaxSampledCTEs))
	w.PutOne(uint8(i.SwitchingPatternLength))
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(uint8(m))
	}
	return w.Data
}

// LESetConnectionlessIQSamplingEnableOutput represents the output of the command specified in Section 7.8.82
type LESetConnectionlessIQSamplingEnableOutput struct {
	Status uint8
	SyncHandle uint16
}

func (o *LESetConnectionlessIQSamplingEnableOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.SyncHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetConnectionlessIQSamplingEnableSync executes the command specified in Section 7.8.82 synchronously
func (c *Commands) LESetConnectionlessIQSamplingEnableSync (params LESetConnectionlessIQSamplingEnableInput, result *LESetConnectionlessIQSamplingEnableOutput) (*LESetConnectionlessIQSamplingEnableOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetConnectionlessIQSamplingEnable started")
	}
	if result == nil {
		result = &LESetConnectionlessIQSamplingEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0053}, nil)
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
		}).Debug("LESetConnectionlessIQSamplingEnable completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.SamplingEnable))
	w.PutOne(uint8(i.SlotDurations))
	w.PutOne(uint8(i.SwitchingPatternLength))
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(uint8(m))
	}
	return w.Data
}

// LESetConnectionCTEReceiveParametersOutput represents the output of the command specified in Section 7.8.83
type LESetConnectionCTEReceiveParametersOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetConnectionCTEReceiveParametersOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetConnectionCTEReceiveParametersSync executes the command specified in Section 7.8.83 synchronously
func (c *Commands) LESetConnectionCTEReceiveParametersSync (params LESetConnectionCTEReceiveParametersInput, result *LESetConnectionCTEReceiveParametersOutput) (*LESetConnectionCTEReceiveParametersOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetConnectionCTEReceiveParameters started")
	}
	if result == nil {
		result = &LESetConnectionCTEReceiveParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0054}, nil)
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
		}).Debug("LESetConnectionCTEReceiveParameters completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.CTETypes))
	w.PutOne(uint8(i.SwitchingPatternLength))
	if len(i.AntennaIDs) != int(i.SwitchingPatternLength) {
		panic("len(i.AntennaIDs) != int(i.SwitchingPatternLength)")
	}
	for _, m := range i.AntennaIDs {
		w.PutOne(uint8(m))
	}
	return w.Data
}

// LESetConnectionCTETransmitParametersOutput represents the output of the command specified in Section 7.8.84
type LESetConnectionCTETransmitParametersOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetConnectionCTETransmitParametersOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetConnectionCTETransmitParametersSync executes the command specified in Section 7.8.84 synchronously
func (c *Commands) LESetConnectionCTETransmitParametersSync (params LESetConnectionCTETransmitParametersInput, result *LESetConnectionCTETransmitParametersOutput) (*LESetConnectionCTETransmitParametersOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetConnectionCTETransmitParameters started")
	}
	if result == nil {
		result = &LESetConnectionCTETransmitParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0055}, nil)
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
		}).Debug("LESetConnectionCTETransmitParameters completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.Enable))
	binary.LittleEndian.PutUint16(w.Put(2), i.CTERequestInterval)
	w.PutOne(uint8(i.RequestedCTELength))
	w.PutOne(uint8(i.RequestedCTEType))
	return w.Data
}

// LEConnectionCTERequestEnableOutput represents the output of the command specified in Section 7.8.85
type LEConnectionCTERequestEnableOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEConnectionCTERequestEnableOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEConnectionCTERequestEnableSync executes the command specified in Section 7.8.85 synchronously
func (c *Commands) LEConnectionCTERequestEnableSync (params LEConnectionCTERequestEnableInput, result *LEConnectionCTERequestEnableOutput) (*LEConnectionCTERequestEnableOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEConnectionCTERequestEnable started")
	}
	if result == nil {
		result = &LEConnectionCTERequestEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0056}, nil)
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
		}).Debug("LEConnectionCTERequestEnable completed")
	}

	 return result, err
}
// LEConnectionCTEResponseEnableInput represents the input of the command specified in Section 7.8.86
type LEConnectionCTEResponseEnableInput struct {
	ConnectionHandle uint16
	Enable uint8
}

func (i LEConnectionCTEResponseEnableInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.Enable))
	return w.Data
}

// LEConnectionCTEResponseEnableOutput represents the output of the command specified in Section 7.8.86
type LEConnectionCTEResponseEnableOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEConnectionCTEResponseEnableOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEConnectionCTEResponseEnableSync executes the command specified in Section 7.8.86 synchronously
func (c *Commands) LEConnectionCTEResponseEnableSync (params LEConnectionCTEResponseEnableInput, result *LEConnectionCTEResponseEnableOutput) (*LEConnectionCTEResponseEnableOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEConnectionCTEResponseEnable started")
	}
	if result == nil {
		result = &LEConnectionCTEResponseEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0057}, nil)
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
		}).Debug("LEConnectionCTEResponseEnable completed")
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
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.SupportedSwitchingSamplingRates = uint8(r.GetOne())
	o.NumAntennae = uint8(r.GetOne())
	o.MaxSwitchingPatternLength = uint8(r.GetOne())
	o.MaxCTELength = uint8(r.GetOne())
	return r.Valid()
}

// LEReadAntennaInformationSync executes the command specified in Section 7.8.87 synchronously
func (c *Commands) LEReadAntennaInformationSync (result *LEReadAntennaInformationOutput) (*LEReadAntennaInformationOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("LEReadAntennaInformation started")
	}
	if result == nil {
		result = &LEReadAntennaInformationOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0058}, nil)
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
		}).Debug("LEReadAntennaInformation completed")
	}

	 return result, err
}
// LESetPeriodicAdvertisingReceiveEnableInput represents the input of the command specified in Section 7.8.88
type LESetPeriodicAdvertisingReceiveEnableInput struct {
	SyncHandle uint16
	Enable uint8
}

func (i LESetPeriodicAdvertisingReceiveEnableInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncHandle)
	w.PutOne(uint8(i.Enable))
	return w.Data
}

// LESetPeriodicAdvertisingReceiveEnableSync executes the command specified in Section 7.8.88 synchronously
func (c *Commands) LESetPeriodicAdvertisingReceiveEnableSync (params LESetPeriodicAdvertisingReceiveEnableInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetPeriodicAdvertisingReceiveEnable started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0059}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetPeriodicAdvertisingReceiveEnable completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.ServiceData)
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncHandle)
	return w.Data
}

// LEPeriodicAdvertisingSyncTransferOutput represents the output of the command specified in Section 7.8.89
type LEPeriodicAdvertisingSyncTransferOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEPeriodicAdvertisingSyncTransferOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEPeriodicAdvertisingSyncTransferSync executes the command specified in Section 7.8.89 synchronously
func (c *Commands) LEPeriodicAdvertisingSyncTransferSync (params LEPeriodicAdvertisingSyncTransferInput, result *LEPeriodicAdvertisingSyncTransferOutput) (*LEPeriodicAdvertisingSyncTransferOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEPeriodicAdvertisingSyncTransfer started")
	}
	if result == nil {
		result = &LEPeriodicAdvertisingSyncTransferOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x005A}, nil)
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
		}).Debug("LEPeriodicAdvertisingSyncTransfer completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.ServiceData)
	w.PutOne(uint8(i.AdvertisingHandle))
	return w.Data
}

// LEPeriodicAdvertisingSetInfoTransferOutput represents the output of the command specified in Section 7.8.90
type LEPeriodicAdvertisingSetInfoTransferOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEPeriodicAdvertisingSetInfoTransferOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEPeriodicAdvertisingSetInfoTransferSync executes the command specified in Section 7.8.90 synchronously
func (c *Commands) LEPeriodicAdvertisingSetInfoTransferSync (params LEPeriodicAdvertisingSetInfoTransferInput, result *LEPeriodicAdvertisingSetInfoTransferOutput) (*LEPeriodicAdvertisingSetInfoTransferOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEPeriodicAdvertisingSetInfoTransfer started")
	}
	if result == nil {
		result = &LEPeriodicAdvertisingSetInfoTransferOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x005B}, nil)
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
		}).Debug("LEPeriodicAdvertisingSetInfoTransfer completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.Mode))
	binary.LittleEndian.PutUint16(w.Put(2), i.Skip)
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncTimeout)
	w.PutOne(uint8(i.CTEType))
	return w.Data
}

// LESetPeriodicAdvertisingSyncTransferParametersOutput represents the output of the command specified in Section 7.8.91
type LESetPeriodicAdvertisingSyncTransferParametersOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetPeriodicAdvertisingSyncTransferParametersOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetPeriodicAdvertisingSyncTransferParametersSync executes the command specified in Section 7.8.91 synchronously
func (c *Commands) LESetPeriodicAdvertisingSyncTransferParametersSync (params LESetPeriodicAdvertisingSyncTransferParametersInput, result *LESetPeriodicAdvertisingSyncTransferParametersOutput) (*LESetPeriodicAdvertisingSyncTransferParametersOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetPeriodicAdvertisingSyncTransferParameters started")
	}
	if result == nil {
		result = &LESetPeriodicAdvertisingSyncTransferParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x005C}, nil)
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
		}).Debug("LESetPeriodicAdvertisingSyncTransferParameters completed")
	}

	 return result, err
}
// LEModifySleepClockAccuracyInput represents the input of the command specified in Section 7.8.94
type LEModifySleepClockAccuracyInput struct {
	Action uint8
}

func (i LEModifySleepClockAccuracyInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.Action))
	return w.Data
}

// LEModifySleepClockAccuracySync executes the command specified in Section 7.8.94 synchronously
func (c *Commands) LEModifySleepClockAccuracySync (params LEModifySleepClockAccuracyInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEModifySleepClockAccuracy started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x005F}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEModifySleepClockAccuracy completed")
	}

	 return err
}
// LEReadISOTXSyncInput represents the input of the command specified in Section 7.8.96
type LEReadISOTXSyncInput struct {
	ConnectionHandle uint16
}

func (i LEReadISOTXSyncInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data
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
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PacketSequenceNumber = binary.LittleEndian.Uint16(r.Get(2))
	o.TimeStamp = binary.LittleEndian.Uint32(r.Get(4))
	o.TimeOffset = bleutil.DecodeUint24(r.Get(3))
	return r.Valid()
}

// LEReadISOTXSyncSync executes the command specified in Section 7.8.96 synchronously
func (c *Commands) LEReadISOTXSyncSync (params LEReadISOTXSyncInput, result *LEReadISOTXSyncOutput) (*LEReadISOTXSyncOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReadISOTXSync started")
	}
	if result == nil {
		result = &LEReadISOTXSyncOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0061}, nil)
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
		}).Debug("LEReadISOTXSync completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.CIGID))
	bleutil.EncodeUint24(w.Put(3), i.SDUIntervalMToS)
	bleutil.EncodeUint24(w.Put(3), i.SDUIntervalSToM)
	w.PutOne(uint8(i.SlavesClockAccuracy))
	w.PutOne(uint8(i.Packing))
	w.PutOne(uint8(i.Framing))
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxTransportLatencyMToS)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxTransportLatencySToM)
	w.PutOne(uint8(i.CISCount))
	if len(i.CISID) != int(i.CISCount) {
		panic("len(i.CISID) != int(i.CISCount)")
	}
	for _, m := range i.CISID {
		w.PutOne(uint8(m))
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
		w.PutOne(uint8(m))
	}
	if len(i.PHYSToM) != int(i.CISCount) {
		panic("len(i.PHYSToM) != int(i.CISCount)")
	}
	for _, m := range i.PHYSToM {
		w.PutOne(uint8(m))
	}
	if len(i.RTNMToS) != int(i.CISCount) {
		panic("len(i.RTNMToS) != int(i.CISCount)")
	}
	for _, m := range i.RTNMToS {
		w.PutOne(uint8(m))
	}
	if len(i.RTNSToM) != int(i.CISCount) {
		panic("len(i.RTNSToM) != int(i.CISCount)")
	}
	for _, m := range i.RTNSToM {
		w.PutOne(uint8(m))
	}
	return w.Data
}

// LESetCIGParametersOutput represents the output of the command specified in Section 7.8.97
type LESetCIGParametersOutput struct {
	Status uint8
	CIGID uint8
	CISCount uint8
	ConnectionHandle []uint16
}

func (o *LESetCIGParametersOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.CIGID = uint8(r.GetOne())
	o.CISCount = uint8(r.GetOne())
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
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetCIGParameters started")
	}
	if result == nil {
		result = &LESetCIGParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0062}, nil)
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
		}).Debug("LESetCIGParameters completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.CIGID))
	bleutil.EncodeUint24(w.Put(3), i.SDUIntervalMToS)
	bleutil.EncodeUint24(w.Put(3), i.SDUIntervalSToM)
	w.PutOne(uint8(i.FTMToS))
	w.PutOne(uint8(i.FTSToM))
	binary.LittleEndian.PutUint16(w.Put(2), i.ISOInterval)
	w.PutOne(uint8(i.SlavesClockAccuracy))
	w.PutOne(uint8(i.Packing))
	w.PutOne(uint8(i.Framing))
	w.PutOne(uint8(i.CISCount))
	if len(i.CISID) != int(i.CISCount) {
		panic("len(i.CISID) != int(i.CISCount)")
	}
	for _, m := range i.CISID {
		w.PutOne(uint8(m))
	}
	if len(i.NSE) != int(i.CISCount) {
		panic("len(i.NSE) != int(i.CISCount)")
	}
	for _, m := range i.NSE {
		w.PutOne(uint8(m))
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
		w.PutOne(uint8(m))
	}
	if len(i.PHYSToM) != int(i.CISCount) {
		panic("len(i.PHYSToM) != int(i.CISCount)")
	}
	for _, m := range i.PHYSToM {
		w.PutOne(uint8(m))
	}
	if len(i.BNMToS) != int(i.CISCount) {
		panic("len(i.BNMToS) != int(i.CISCount)")
	}
	for _, m := range i.BNMToS {
		w.PutOne(uint8(m))
	}
	if len(i.BNSToM) != int(i.CISCount) {
		panic("len(i.BNSToM) != int(i.CISCount)")
	}
	for _, m := range i.BNSToM {
		w.PutOne(uint8(m))
	}
	return w.Data
}

// LESetCIGParametersTestOutput represents the output of the command specified in Section 7.8.98
type LESetCIGParametersTestOutput struct {
	Status uint8
	CIGID uint8
	CISCount uint8
	ConnectionHandle []uint16
}

func (o *LESetCIGParametersTestOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.CIGID = uint8(r.GetOne())
	o.CISCount = uint8(r.GetOne())
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
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetCIGParametersTest started")
	}
	if result == nil {
		result = &LESetCIGParametersTestOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0063}, nil)
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
		}).Debug("LESetCIGParametersTest completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.CISCount))
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
	return w.Data
}

// LECreateCISSync executes the command specified in Section 7.8.99 synchronously
func (c *Commands) LECreateCISSync (params LECreateCISInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LECreateCIS started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0064}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LECreateCIS completed")
	}

	 return err
}
// LERemoveCIGInput represents the input of the command specified in Section 7.8.100
type LERemoveCIGInput struct {
	CIGID uint8
}

func (i LERemoveCIGInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.CIGID))
	return w.Data
}

// LERemoveCIGOutput represents the output of the command specified in Section 7.8.100
type LERemoveCIGOutput struct {
	Status uint8
	CIGID uint8
}

func (o *LERemoveCIGOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.CIGID = uint8(r.GetOne())
	return r.Valid()
}

// LERemoveCIGSync executes the command specified in Section 7.8.100 synchronously
func (c *Commands) LERemoveCIGSync (params LERemoveCIGInput, result *LERemoveCIGOutput) (*LERemoveCIGOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LERemoveCIG started")
	}
	if result == nil {
		result = &LERemoveCIGOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0065}, nil)
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
		}).Debug("LERemoveCIG completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.BIGHandle))
	w.PutOne(uint8(i.AdvertisingHandle))
	w.PutOne(uint8(i.NumBIS))
	bleutil.EncodeUint24(w.Put(3), i.SDUInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxSDU)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxTransportLatency)
	w.PutOne(uint8(i.RTN))
	w.PutOne(uint8(i.PHY))
	w.PutOne(uint8(i.Packing))
	w.PutOne(uint8(i.Framing))
	w.PutOne(uint8(i.Encryption))
	copy(w.Put(16), i.BroadcastCode[:])
	return w.Data
}

// LEAcceptCISRequestSync executes the command specified in Section 7.8.101 synchronously
func (c *Commands) LEAcceptCISRequestSync (params LEAcceptCISRequestInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEAcceptCISRequest started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0066}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEAcceptCISRequest completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.BIGHandle))
	w.PutOne(uint8(i.AdvertisingHandle))
	w.PutOne(uint8(i.NumBIS))
	bleutil.EncodeUint24(w.Put(3), i.SDUInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.ISOInterval)
	w.PutOne(uint8(i.NSE))
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxSDU)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxPDU)
	w.PutOne(uint8(i.PHY))
	w.PutOne(uint8(i.Packing))
	w.PutOne(uint8(i.Framing))
	w.PutOne(uint8(i.BN))
	w.PutOne(uint8(i.IRC))
	w.PutOne(uint8(i.PTO))
	w.PutOne(uint8(i.Encryption))
	copy(w.Put(16), i.BroadcastCode[:])
	return w.Data
}

// LECreateBIGTestSync executes the command specified in Section 7.8.104 synchronously
func (c *Commands) LECreateBIGTestSync (params LECreateBIGTestInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LECreateBIGTest started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0069}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LECreateBIGTest completed")
	}

	 return err
}
// LETerminateBIGInput represents the input of the command specified in Section 7.8.105
type LETerminateBIGInput struct {
	BIGHandle uint8
	Reason uint8
}

func (i LETerminateBIGInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.BIGHandle))
	w.PutOne(uint8(i.Reason))
	return w.Data
}

// LETerminateBIGSync executes the command specified in Section 7.8.105 synchronously
func (c *Commands) LETerminateBIGSync (params LETerminateBIGInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LETerminateBIG started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006A}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LETerminateBIG completed")
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
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.BIGHandle))
	binary.LittleEndian.PutUint16(w.Put(2), i.SyncHandle)
	w.PutOne(uint8(i.Encryption))
	copy(w.Put(16), i.BroadcastCode[:])
	w.PutOne(uint8(i.MSE))
	binary.LittleEndian.PutUint16(w.Put(2), i.BIGSyncTimeout)
	w.PutOne(uint8(i.NumBIS))
	if len(i.BIS) != int(i.NumBIS) {
		panic("len(i.BIS) != int(i.NumBIS)")
	}
	for _, m := range i.BIS {
		w.PutOne(uint8(m))
	}
	return w.Data
}

// LEBIGCreateSyncSync executes the command specified in Section 7.8.106 synchronously
func (c *Commands) LEBIGCreateSyncSync (params LEBIGCreateSyncInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEBIGCreateSync started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006B}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEBIGCreateSync completed")
	}

	 return err
}
// LEBIGTerminateSyncInput represents the input of the command specified in Section 7.8.107
type LEBIGTerminateSyncInput struct {
	BIGHandle uint8
}

func (i LEBIGTerminateSyncInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.BIGHandle))
	return w.Data
}

// LEBIGTerminateSyncOutput represents the output of the command specified in Section 7.8.107
type LEBIGTerminateSyncOutput struct {
	Status uint8
	BIGHandle uint8
}

func (o *LEBIGTerminateSyncOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.BIGHandle = uint8(r.GetOne())
	return r.Valid()
}

// LEBIGTerminateSyncSync executes the command specified in Section 7.8.107 synchronously
func (c *Commands) LEBIGTerminateSyncSync (params LEBIGTerminateSyncInput, result *LEBIGTerminateSyncOutput) (*LEBIGTerminateSyncOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEBIGTerminateSync started")
	}
	if result == nil {
		result = &LEBIGTerminateSyncOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006C}, nil)
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
		}).Debug("LEBIGTerminateSync completed")
	}

	 return result, err
}
// LERequestPeerSCAInput represents the input of the command specified in Section 7.8.108
type LERequestPeerSCAInput struct {
	ConnectionHandle uint16
}

func (i LERequestPeerSCAInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data
}

// LERequestPeerSCASync executes the command specified in Section 7.8.108 synchronously
func (c *Commands) LERequestPeerSCASync (params LERequestPeerSCAInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LERequestPeerSCA started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006D}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LERequestPeerSCA completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.DataPathDirection))
	w.PutOne(uint8(i.DataPathID))
	copy(w.Put(5), i.CodecID[:])
	bleutil.EncodeUint24(w.Put(3), i.ControllerDelay)
	w.PutOne(uint8(i.CodecConfigurationLength))
	w.PutSlice(i.CodecConfiguration)
	return w.Data
}

// LESetupISODataPathOutput represents the output of the command specified in Section 7.8.109
type LESetupISODataPathOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetupISODataPathOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetupISODataPathSync executes the command specified in Section 7.8.109 synchronously
func (c *Commands) LESetupISODataPathSync (params LESetupISODataPathInput, result *LESetupISODataPathOutput) (*LESetupISODataPathOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetupISODataPath started")
	}
	if result == nil {
		result = &LESetupISODataPathOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006E}, nil)
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
		}).Debug("LESetupISODataPath completed")
	}

	 return result, err
}
// LERemoveISODataPathInput represents the input of the command specified in Section 7.8.110
type LERemoveISODataPathInput struct {
	ConnectionHandle uint16
	DataPathDirection uint8
}

func (i LERemoveISODataPathInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.DataPathDirection))
	return w.Data
}

// LERemoveISODataPathOutput represents the output of the command specified in Section 7.8.110
type LERemoveISODataPathOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LERemoveISODataPathOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LERemoveISODataPathSync executes the command specified in Section 7.8.110 synchronously
func (c *Commands) LERemoveISODataPathSync (params LERemoveISODataPathInput, result *LERemoveISODataPathOutput) (*LERemoveISODataPathOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LERemoveISODataPath started")
	}
	if result == nil {
		result = &LERemoveISODataPathOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x006F}, nil)
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
		}).Debug("LERemoveISODataPath completed")
	}

	 return result, err
}
// LEISOTransmitTestInput represents the input of the command specified in Section 7.8.111
type LEISOTransmitTestInput struct {
	ConnectionHandle uint16
	PayloadType uint8
}

func (i LEISOTransmitTestInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.PayloadType))
	return w.Data
}

// LEISOTransmitTestOutput represents the output of the command specified in Section 7.8.111
type LEISOTransmitTestOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEISOTransmitTestOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEISOTransmitTestSync executes the command specified in Section 7.8.111 synchronously
func (c *Commands) LEISOTransmitTestSync (params LEISOTransmitTestInput, result *LEISOTransmitTestOutput) (*LEISOTransmitTestOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEISOTransmitTest started")
	}
	if result == nil {
		result = &LEISOTransmitTestOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0070}, nil)
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
		}).Debug("LEISOTransmitTest completed")
	}

	 return result, err
}
// LEISOReceiveTestInput represents the input of the command specified in Section 7.8.112
type LEISOReceiveTestInput struct {
	ConnectionHandle uint16
	PayloadType uint8
}

func (i LEISOReceiveTestInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.PayloadType))
	return w.Data
}

// LEISOReceiveTestOutput represents the output of the command specified in Section 7.8.112
type LEISOReceiveTestOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LEISOReceiveTestOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LEISOReceiveTestSync executes the command specified in Section 7.8.112 synchronously
func (c *Commands) LEISOReceiveTestSync (params LEISOReceiveTestInput, result *LEISOReceiveTestOutput) (*LEISOReceiveTestOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEISOReceiveTest started")
	}
	if result == nil {
		result = &LEISOReceiveTestOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0071}, nil)
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
		}).Debug("LEISOReceiveTest completed")
	}

	 return result, err
}
// LEISOReadTestCountersInput represents the input of the command specified in Section 7.8.113
type LEISOReadTestCountersInput struct {
	ConnectionHandle uint16
}

func (i LEISOReadTestCountersInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data
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
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ReceivedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	o.MissedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	o.FailedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// LEISOReadTestCountersSync executes the command specified in Section 7.8.113 synchronously
func (c *Commands) LEISOReadTestCountersSync (params LEISOReadTestCountersInput, result *LEISOReadTestCountersOutput) (*LEISOReadTestCountersOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEISOReadTestCounters started")
	}
	if result == nil {
		result = &LEISOReadTestCountersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0072}, nil)
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
		}).Debug("LEISOReadTestCounters completed")
	}

	 return result, err
}
// LEISOTestEndInput represents the input of the command specified in Section 7.8.114
type LEISOTestEndInput struct {
	ConnectionHandle uint16
}

func (i LEISOTestEndInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data
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
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.ReceivedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	o.MissedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	o.FailedPacketCount = binary.LittleEndian.Uint32(r.Get(4))
	return r.Valid()
}

// LEISOTestEndSync executes the command specified in Section 7.8.114 synchronously
func (c *Commands) LEISOTestEndSync (params LEISOTestEndInput, result *LEISOTestEndOutput) (*LEISOTestEndOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEISOTestEnd started")
	}
	if result == nil {
		result = &LEISOTestEndOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0073}, nil)
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
		}).Debug("LEISOTestEnd completed")
	}

	 return result, err
}
// LESetHostFeatureInput represents the input of the command specified in Section 7.8.115
type LESetHostFeatureInput struct {
	BitNumber uint8
	BitValue uint8
}

func (i LESetHostFeatureInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.BitNumber))
	w.PutOne(uint8(i.BitValue))
	return w.Data
}

// LESetHostFeatureSync executes the command specified in Section 7.8.115 synchronously
func (c *Commands) LESetHostFeatureSync (params LESetHostFeatureInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetHostFeature started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0074}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LESetHostFeature completed")
	}

	 return err
}
// LEReadISOLinkQualityInput represents the input of the command specified in Section 7.8.116
type LEReadISOLinkQualityInput struct {
	ConnectionHandle uint16
}

func (i LEReadISOLinkQualityInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data
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
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
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
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReadISOLinkQuality started")
	}
	if result == nil {
		result = &LEReadISOLinkQualityOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0075}, nil)
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
		}).Debug("LEReadISOLinkQuality completed")
	}

	 return result, err
}
// LEEnhancedReadTransmitPowerLevelInput represents the input of the command specified in Section 7.8.117
type LEEnhancedReadTransmitPowerLevelInput struct {
	ConnectionHandle uint16
	PHY uint8
}

func (i LEEnhancedReadTransmitPowerLevelInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.PHY))
	return w.Data
}

// LEEnhancedReadTransmitPowerLevelOutput represents the output of the command specified in Section 7.8.117
type LEEnhancedReadTransmitPowerLevelOutput struct {
	Status uint8
	ConnectionHandle uint16
	PHY uint8
	CurrentTransmitPowerLevel uint8
}

func (o *LEEnhancedReadTransmitPowerLevelOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.PHY = uint8(r.GetOne())
	o.CurrentTransmitPowerLevel = uint8(r.GetOne())
	return r.Valid()
}

// LEEnhancedReadTransmitPowerLevelSync executes the command specified in Section 7.8.117 synchronously
func (c *Commands) LEEnhancedReadTransmitPowerLevelSync (params LEEnhancedReadTransmitPowerLevelInput, result *LEEnhancedReadTransmitPowerLevelOutput) (*LEEnhancedReadTransmitPowerLevelOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEEnhancedReadTransmitPowerLevel started")
	}
	if result == nil {
		result = &LEEnhancedReadTransmitPowerLevelOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0076}, nil)
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
		}).Debug("LEEnhancedReadTransmitPowerLevel completed")
	}

	 return result, err
}
// LEReadRemoteTransmitPowerLevelInput represents the input of the command specified in Section 7.8.118
type LEReadRemoteTransmitPowerLevelInput struct {
	ConnectionHandle uint16
	PHY uint8
}

func (i LEReadRemoteTransmitPowerLevelInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.PHY))
	return w.Data
}

// LEReadRemoteTransmitPowerLevelSync executes the command specified in Section 7.8.118 synchronously
func (c *Commands) LEReadRemoteTransmitPowerLevelSync (params LEReadRemoteTransmitPowerLevelInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LEReadRemoteTransmitPowerLevel started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0077}, nil)
	if err != nil {
		goto log
	}

	buffer.Buffer = params.encode(buffer.Buffer)
	response, err = c.hcicmdmgr.CommandRunPutBuffer(buffer)
	if err != nil {
		goto log
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
		}).Debug("LEReadRemoteTransmitPowerLevel completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.HighThreshold))
	w.PutOne(uint8(i.HighHysteresis))
	w.PutOne(uint8(i.LowThreshold))
	w.PutOne(uint8(i.LowHysteresis))
	binary.LittleEndian.PutUint16(w.Put(2), i.MinTimeSpent)
	return w.Data
}

// LESetPathLossReportingParametersOutput represents the output of the command specified in Section 7.8.119
type LESetPathLossReportingParametersOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetPathLossReportingParametersOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetPathLossReportingParametersSync executes the command specified in Section 7.8.119 synchronously
func (c *Commands) LESetPathLossReportingParametersSync (params LESetPathLossReportingParametersInput, result *LESetPathLossReportingParametersOutput) (*LESetPathLossReportingParametersOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetPathLossReportingParameters started")
	}
	if result == nil {
		result = &LESetPathLossReportingParametersOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0078}, nil)
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
		}).Debug("LESetPathLossReportingParameters completed")
	}

	 return result, err
}
// LESetPathLossReportingEnableInput represents the input of the command specified in Section 7.8.120
type LESetPathLossReportingEnableInput struct {
	ConnectionHandle uint16
	Enable uint8
}

func (i LESetPathLossReportingEnableInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.Enable))
	return w.Data
}

// LESetPathLossReportingEnableOutput represents the output of the command specified in Section 7.8.120
type LESetPathLossReportingEnableOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetPathLossReportingEnableOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetPathLossReportingEnableSync executes the command specified in Section 7.8.120 synchronously
func (c *Commands) LESetPathLossReportingEnableSync (params LESetPathLossReportingEnableInput, result *LESetPathLossReportingEnableOutput) (*LESetPathLossReportingEnableOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetPathLossReportingEnable started")
	}
	if result == nil {
		result = &LESetPathLossReportingEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x0079}, nil)
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
		}).Debug("LESetPathLossReportingEnable completed")
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
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.LocalEnable))
	w.PutOne(uint8(i.RemoteEnable))
	return w.Data
}

// LESetTransmitPowerReportingEnableOutput represents the output of the command specified in Section 7.8.121
type LESetTransmitPowerReportingEnableOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LESetTransmitPowerReportingEnableOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LESetTransmitPowerReportingEnableSync executes the command specified in Section 7.8.121 synchronously
func (c *Commands) LESetTransmitPowerReportingEnableSync (params LESetTransmitPowerReportingEnableInput, result *LESetTransmitPowerReportingEnableOutput) (*LESetTransmitPowerReportingEnableOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("LESetTransmitPowerReportingEnable started")
	}
	if result == nil {
		result = &LESetTransmitPowerReportingEnableOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 8, OCF: 0x007A}, nil)
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
		}).Debug("LESetTransmitPowerReportingEnable completed")
	}

	 return result, err
}
