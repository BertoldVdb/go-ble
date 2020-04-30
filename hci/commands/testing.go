package hcicommands

import (
	"encoding/binary"
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
)

// TestingReadLoopbackModeOutput represents the output of the command specified in Section 7.6.1
type TestingReadLoopbackModeOutput struct {
	Status uint8
	LoopbackMode uint8
}

func (o *TestingReadLoopbackModeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.LoopbackMode = r.GetOne()
	return r.Valid()
}

// TestingReadLoopbackModeSync executes the command specified in Section 7.6.1 synchronously
func (c *Commands) TestingReadLoopbackModeSync (result *TestingReadLoopbackModeOutput) (*TestingReadLoopbackModeOutput, error) {
	if result == nil {
		result = &TestingReadLoopbackModeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0001}, nil)
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

// TestingWriteLoopbackModeInput represents the input of the command specified in Section 7.6.2
type TestingWriteLoopbackModeInput struct {
	LoopbackMode uint8
}

func (i TestingWriteLoopbackModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.LoopbackMode)
	return w.Data()
}

// TestingWriteLoopbackModeSync executes the command specified in Section 7.6.2 synchronously
func (c *Commands) TestingWriteLoopbackModeSync (params TestingWriteLoopbackModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0002}, nil)
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

// TestingEnableDeviceUnderTestModeSync executes the command specified in Section 7.6.3 synchronously
func (c *Commands) TestingEnableDeviceUnderTestModeSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0003}, nil)
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

// TestingWriteSimplePairingDebugModeInput represents the input of the command specified in Section 7.6.4
type TestingWriteSimplePairingDebugModeInput struct {
	SimplePairingDebugMode uint8
}

func (i TestingWriteSimplePairingDebugModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.SimplePairingDebugMode)
	return w.Data()
}

// TestingWriteSimplePairingDebugModeSync executes the command specified in Section 7.6.4 synchronously
func (c *Commands) TestingWriteSimplePairingDebugModeSync (params TestingWriteSimplePairingDebugModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0004}, nil)
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

// TestingEnableAMPReceiverReportsInput represents the input of the command specified in Section 7.6.5
type TestingEnableAMPReceiverReportsInput struct {
	Enable uint8
	Interval uint8
}

func (i TestingEnableAMPReceiverReportsInput) encode(data []byte) []byte {
	w := writer{data: data};
	w.PutOne(i.Enable)
	w.PutOne(i.Interval)
	return w.Data()
}

// TestingEnableAMPReceiverReportsSync executes the command specified in Section 7.6.5 synchronously
func (c *Commands) TestingEnableAMPReceiverReportsSync (params TestingEnableAMPReceiverReportsInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0007}, nil)
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

// TestingAMPTestEndSync executes the command specified in Section 7.6.6 synchronously
func (c *Commands) TestingAMPTestEndSync () error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x0008}, nil)
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

// TestingWriteSecureConnectionsTestModeInput represents the input of the command specified in Section 7.6.8
type TestingWriteSecureConnectionsTestModeInput struct {
	ConnectionHandle uint16
	DM1ACLUMode uint8
	eSCOLoopbackMode uint8
}

func (i TestingWriteSecureConnectionsTestModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.DM1ACLUMode)
	w.PutOne(i.eSCOLoopbackMode)
	return w.Data()
}

// TestingWriteSecureConnectionsTestModeOutput represents the output of the command specified in Section 7.6.8
type TestingWriteSecureConnectionsTestModeOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *TestingWriteSecureConnectionsTestModeOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// TestingWriteSecureConnectionsTestModeSync executes the command specified in Section 7.6.8 synchronously
func (c *Commands) TestingWriteSecureConnectionsTestModeSync (params TestingWriteSecureConnectionsTestModeInput, result *TestingWriteSecureConnectionsTestModeOutput) (*TestingWriteSecureConnectionsTestModeOutput, error) {
	if result == nil {
		result = &TestingWriteSecureConnectionsTestModeOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 6, OCF: 0x000A}, nil)
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

