package hcicommands

import (
	"encoding/binary"
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

// TestingReadLoopbackModeOutput represents the output of the command specified in Section 7.6.1
type TestingReadLoopbackModeOutput struct {
	Status uint8
	LoopbackMode uint8
}

func (o *TestingReadLoopbackModeOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.LoopbackMode = uint8(r.GetOne())
	return r.Valid()
}

// TestingReadLoopbackModeSync executes the command specified in Section 7.6.1 synchronously
func (c *Commands) TestingReadLoopbackModeSync (result *TestingReadLoopbackModeOutput) (*TestingReadLoopbackModeOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("TestingReadLoopbackMode started")
	}
	if result == nil {
		result = &TestingReadLoopbackModeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0001}, nil)
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
		}).Debug("TestingReadLoopbackMode completed")
	}

	 return result, err
}
// TestingWriteLoopbackModeInput represents the input of the command specified in Section 7.6.2
type TestingWriteLoopbackModeInput struct {
	LoopbackMode uint8
}

func (i TestingWriteLoopbackModeInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.LoopbackMode))
	return w.Data
}

// TestingWriteLoopbackModeSync executes the command specified in Section 7.6.2 synchronously
func (c *Commands) TestingWriteLoopbackModeSync (params TestingWriteLoopbackModeInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("TestingWriteLoopbackMode started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0002}, nil)
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
		}).Debug("TestingWriteLoopbackMode completed")
	}

	 return err
}
// TestingEnableDeviceUnderTestModeSync executes the command specified in Section 7.6.3 synchronously
func (c *Commands) TestingEnableDeviceUnderTestModeSync () error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("TestingEnableDeviceUnderTestMode started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0003}, nil)
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
		}).Debug("TestingEnableDeviceUnderTestMode completed")
	}

	 return err
}
// TestingWriteSimplePairingDebugModeInput represents the input of the command specified in Section 7.6.4
type TestingWriteSimplePairingDebugModeInput struct {
	SimplePairingDebugMode uint8
}

func (i TestingWriteSimplePairingDebugModeInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.SimplePairingDebugMode))
	return w.Data
}

// TestingWriteSimplePairingDebugModeSync executes the command specified in Section 7.6.4 synchronously
func (c *Commands) TestingWriteSimplePairingDebugModeSync (params TestingWriteSimplePairingDebugModeInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("TestingWriteSimplePairingDebugMode started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0004}, nil)
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
		}).Debug("TestingWriteSimplePairingDebugMode completed")
	}

	 return err
}
// TestingEnableAMPReceiverReportsInput represents the input of the command specified in Section 7.6.5
type TestingEnableAMPReceiverReportsInput struct {
	Enable uint8
	Interval uint8
}

func (i TestingEnableAMPReceiverReportsInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	w.PutOne(uint8(i.Enable))
	w.PutOne(uint8(i.Interval))
	return w.Data
}

// TestingEnableAMPReceiverReportsSync executes the command specified in Section 7.6.5 synchronously
func (c *Commands) TestingEnableAMPReceiverReportsSync (params TestingEnableAMPReceiverReportsInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("TestingEnableAMPReceiverReports started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0007}, nil)
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
		}).Debug("TestingEnableAMPReceiverReports completed")
	}

	 return err
}
// TestingAMPTestEndSync executes the command specified in Section 7.6.6 synchronously
func (c *Commands) TestingAMPTestEndSync () error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
		}).Trace("TestingAMPTestEnd started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0008}, nil)
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
		}).Debug("TestingAMPTestEnd completed")
	}

	 return err
}
// TestingWriteSecureConnectionsTestModeInput represents the input of the command specified in Section 7.6.8
type TestingWriteSecureConnectionsTestModeInput struct {
	ConnectionHandle uint16
	DM1ACLUMode uint8
	eSCOLoopbackMode uint8
}

func (i TestingWriteSecureConnectionsTestModeInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.DM1ACLUMode))
	w.PutOne(uint8(i.eSCOLoopbackMode))
	return w.Data
}

// TestingWriteSecureConnectionsTestModeOutput represents the output of the command specified in Section 7.6.8
type TestingWriteSecureConnectionsTestModeOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *TestingWriteSecureConnectionsTestModeOutput) decode(data []byte) bool {
	r := bleutil.Reader{Data: data};
	o.Status = uint8(r.GetOne())
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// TestingWriteSecureConnectionsTestModeSync executes the command specified in Section 7.6.8 synchronously
func (c *Commands) TestingWriteSecureConnectionsTestModeSync (params TestingWriteSecureConnectionsTestModeInput, result *TestingWriteSecureConnectionsTestModeOutput) (*TestingWriteSecureConnectionsTestModeOutput, error) {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("TestingWriteSecureConnectionsTestMode started")
	}
	if result == nil {
		result = &TestingWriteSecureConnectionsTestModeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x000A}, nil)
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
		}).Debug("TestingWriteSecureConnectionsTestMode completed")
	}

	 return result, err
}
