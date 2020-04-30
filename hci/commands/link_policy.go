package hcicommands

import (
	"encoding/binary"
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
)

// LinkPolicyHoldModeInput represents the input of the command specified in Section 7.2.1
type LinkPolicyHoldModeInput struct {
	ConnectionHandle uint16
	HoldModeMaxInterval uint16
	HoldModeMinInterval uint16
}

func (i LinkPolicyHoldModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.HoldModeMaxInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.HoldModeMinInterval)
	return w.Data()
}

// LinkPolicyHoldModeSync executes the command specified in Section 7.2.1 synchronously
func (c *Commands) LinkPolicyHoldModeSync (params LinkPolicyHoldModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x0001}, nil)
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

// LinkPolicySniffModeInput represents the input of the command specified in Section 7.2.2
type LinkPolicySniffModeInput struct {
	ConnectionHandle uint16
	SniffMaxInterval uint16
	SniffMinInterval uint16
	SniffAttempt uint16
	SniffTimeout uint16
}

func (i LinkPolicySniffModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.SniffMaxInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.SniffMinInterval)
	binary.LittleEndian.PutUint16(w.Put(2), i.SniffAttempt)
	binary.LittleEndian.PutUint16(w.Put(2), i.SniffTimeout)
	return w.Data()
}

// LinkPolicySniffModeSync executes the command specified in Section 7.2.2 synchronously
func (c *Commands) LinkPolicySniffModeSync (params LinkPolicySniffModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x0003}, nil)
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

// LinkPolicyExitSniffModeInput represents the input of the command specified in Section 7.2.3
type LinkPolicyExitSniffModeInput struct {
	ConnectionHandle uint16
}

func (i LinkPolicyExitSniffModeInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LinkPolicyExitSniffModeSync executes the command specified in Section 7.2.3 synchronously
func (c *Commands) LinkPolicyExitSniffModeSync (params LinkPolicyExitSniffModeInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x0004}, nil)
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

// LinkPolicyQoSSetupInput represents the input of the command specified in Section 7.2.6
type LinkPolicyQoSSetupInput struct {
	ConnectionHandle uint16
	Unused uint8
	ServiceType uint8
	TokenRate uint32
	PeakBandwidth uint32
	Latency uint32
	DelayVariation uint32
}

func (i LinkPolicyQoSSetupInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.Unused)
	w.PutOne(i.ServiceType)
	binary.LittleEndian.PutUint32(w.Put(4), i.TokenRate)
	binary.LittleEndian.PutUint32(w.Put(4), i.PeakBandwidth)
	binary.LittleEndian.PutUint32(w.Put(4), i.Latency)
	binary.LittleEndian.PutUint32(w.Put(4), i.DelayVariation)
	return w.Data()
}

// LinkPolicyQoSSetupSync executes the command specified in Section 7.2.6 synchronously
func (c *Commands) LinkPolicyQoSSetupSync (params LinkPolicyQoSSetupInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x0007}, nil)
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

// LinkPolicyRoleDiscoveryInput represents the input of the command specified in Section 7.2.7
type LinkPolicyRoleDiscoveryInput struct {
	ConnectionHandle uint16
}

func (i LinkPolicyRoleDiscoveryInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LinkPolicyRoleDiscoveryOutput represents the output of the command specified in Section 7.2.7
type LinkPolicyRoleDiscoveryOutput struct {
	Status uint8
	ConnectionHandle uint16
	CurrentRole uint8
}

func (o *LinkPolicyRoleDiscoveryOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.CurrentRole = r.GetOne()
	return r.Valid()
}

// LinkPolicyRoleDiscoverySync executes the command specified in Section 7.2.7 synchronously
func (c *Commands) LinkPolicyRoleDiscoverySync (params LinkPolicyRoleDiscoveryInput, result *LinkPolicyRoleDiscoveryOutput) (*LinkPolicyRoleDiscoveryOutput, error) {
	if result == nil {
		result = &LinkPolicyRoleDiscoveryOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x0009}, nil)
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

// LinkPolicySwitchRoleInput represents the input of the command specified in Section 7.2.8
type LinkPolicySwitchRoleInput struct {
	BDADDR [6]byte
	Role uint8
}

func (i LinkPolicySwitchRoleInput) encode(data []byte) []byte {
	w := writer{data: data};
	copy(w.Put(6), i.BDADDR[:])
	w.PutOne(i.Role)
	return w.Data()
}

// LinkPolicySwitchRoleSync executes the command specified in Section 7.2.8 synchronously
func (c *Commands) LinkPolicySwitchRoleSync (params LinkPolicySwitchRoleInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x000B}, nil)
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

// LinkPolicyReadLinkPolicySettingsInput represents the input of the command specified in Section 7.2.9
type LinkPolicyReadLinkPolicySettingsInput struct {
	ConnectionHandle uint16
}

func (i LinkPolicyReadLinkPolicySettingsInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	return w.Data()
}

// LinkPolicyReadLinkPolicySettingsOutput represents the output of the command specified in Section 7.2.9
type LinkPolicyReadLinkPolicySettingsOutput struct {
	Status uint8
	ConnectionHandle uint16
	LinkPolicySettings uint16
}

func (o *LinkPolicyReadLinkPolicySettingsOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	o.LinkPolicySettings = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LinkPolicyReadLinkPolicySettingsSync executes the command specified in Section 7.2.9 synchronously
func (c *Commands) LinkPolicyReadLinkPolicySettingsSync (params LinkPolicyReadLinkPolicySettingsInput, result *LinkPolicyReadLinkPolicySettingsOutput) (*LinkPolicyReadLinkPolicySettingsOutput, error) {
	if result == nil {
		result = &LinkPolicyReadLinkPolicySettingsOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x000C}, nil)
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

// LinkPolicyWriteLinkPolicySettingsInput represents the input of the command specified in Section 7.2.10
type LinkPolicyWriteLinkPolicySettingsInput struct {
	ConnectionHandle uint16
	LinkPolicySettings uint16
}

func (i LinkPolicyWriteLinkPolicySettingsInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.LinkPolicySettings)
	return w.Data()
}

// LinkPolicyWriteLinkPolicySettingsOutput represents the output of the command specified in Section 7.2.10
type LinkPolicyWriteLinkPolicySettingsOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LinkPolicyWriteLinkPolicySettingsOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LinkPolicyWriteLinkPolicySettingsSync executes the command specified in Section 7.2.10 synchronously
func (c *Commands) LinkPolicyWriteLinkPolicySettingsSync (params LinkPolicyWriteLinkPolicySettingsInput, result *LinkPolicyWriteLinkPolicySettingsOutput) (*LinkPolicyWriteLinkPolicySettingsOutput, error) {
	if result == nil {
		result = &LinkPolicyWriteLinkPolicySettingsOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x000D}, nil)
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

// LinkPolicyReadDefaultLinkPolicySettingsOutput represents the output of the command specified in Section 7.2.11
type LinkPolicyReadDefaultLinkPolicySettingsOutput struct {
	Status uint8
	DefaultLinkPolicySettings uint16
}

func (o *LinkPolicyReadDefaultLinkPolicySettingsOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.DefaultLinkPolicySettings = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LinkPolicyReadDefaultLinkPolicySettingsSync executes the command specified in Section 7.2.11 synchronously
func (c *Commands) LinkPolicyReadDefaultLinkPolicySettingsSync (result *LinkPolicyReadDefaultLinkPolicySettingsOutput) (*LinkPolicyReadDefaultLinkPolicySettingsOutput, error) {
	if result == nil {
		result = &LinkPolicyReadDefaultLinkPolicySettingsOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x000E}, nil)
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

// LinkPolicyWriteDefaultLinkPolicySettingsInput represents the input of the command specified in Section 7.2.12
type LinkPolicyWriteDefaultLinkPolicySettingsInput struct {
	DefaultLinkPolicySettings uint16
}

func (i LinkPolicyWriteDefaultLinkPolicySettingsInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.DefaultLinkPolicySettings)
	return w.Data()
}

// LinkPolicyWriteDefaultLinkPolicySettingsSync executes the command specified in Section 7.2.12 synchronously
func (c *Commands) LinkPolicyWriteDefaultLinkPolicySettingsSync (params LinkPolicyWriteDefaultLinkPolicySettingsInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x000F}, nil)
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

// LinkPolicyFlowSpecificationInput represents the input of the command specified in Section 7.2.13
type LinkPolicyFlowSpecificationInput struct {
	ConnectionHandle uint16
	Unused uint8
	FlowDirection uint8
	ServiceType uint8
	TokenRate uint32
	TokenBucketSize uint32
	PeakBandwidth uint32
	AccessLatency uint32
}

func (i LinkPolicyFlowSpecificationInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(i.Unused)
	w.PutOne(i.FlowDirection)
	w.PutOne(i.ServiceType)
	binary.LittleEndian.PutUint32(w.Put(4), i.TokenRate)
	binary.LittleEndian.PutUint32(w.Put(4), i.TokenBucketSize)
	binary.LittleEndian.PutUint32(w.Put(4), i.PeakBandwidth)
	binary.LittleEndian.PutUint32(w.Put(4), i.AccessLatency)
	return w.Data()
}

// LinkPolicyFlowSpecificationSync executes the command specified in Section 7.2.13 synchronously
func (c *Commands) LinkPolicyFlowSpecificationSync (params LinkPolicyFlowSpecificationInput) error {
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x0010}, nil)
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

// LinkPolicySniffSubratingInput represents the input of the command specified in Section 7.2.14
type LinkPolicySniffSubratingInput struct {
	ConnectionHandle uint16
	MaxLatency uint16
	MinRemoteTimeout uint16
	MinLocalTimeout uint16
}

func (i LinkPolicySniffSubratingInput) encode(data []byte) []byte {
	w := writer{data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	binary.LittleEndian.PutUint16(w.Put(2), i.MaxLatency)
	binary.LittleEndian.PutUint16(w.Put(2), i.MinRemoteTimeout)
	binary.LittleEndian.PutUint16(w.Put(2), i.MinLocalTimeout)
	return w.Data()
}

// LinkPolicySniffSubratingOutput represents the output of the command specified in Section 7.2.14
type LinkPolicySniffSubratingOutput struct {
	Status uint8
	ConnectionHandle uint16
}

func (o *LinkPolicySniffSubratingOutput) decode(data []byte) bool {
	r := reader{data: data};
	o.Status = r.GetOne()
	o.ConnectionHandle = binary.LittleEndian.Uint16(r.Get(2))
	return r.Valid()
}

// LinkPolicySniffSubratingSync executes the command specified in Section 7.2.14 synchronously
func (c *Commands) LinkPolicySniffSubratingSync (params LinkPolicySniffSubratingInput, result *LinkPolicySniffSubratingOutput) (*LinkPolicySniffSubratingOutput, error) {
	if result == nil {
		result = &LinkPolicySniffSubratingOutput{}
	}

	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 2, OCF: 0x0011}, nil)
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

