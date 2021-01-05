package hcicommands

import (
	"encoding/binary"
	hcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

// VendorCypressSetTxPwrInput represents the input of the command specified in Section 7.63.121
type VendorCypressSetTxPwrInput struct {
	ConnectionHandle uint16
	TXPower uint8
}

func (i VendorCypressSetTxPwrInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	binary.LittleEndian.PutUint16(w.Put(2), i.ConnectionHandle)
	w.PutOne(uint8(i.TXPower))
	return w.Data
}

// VendorCypressSetTxPwrSync executes the command specified in Section 7.63.121 synchronously
func (c *Commands) VendorCypressSetTxPwrSync (params VendorCypressSetTxPwrInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("VendorCypressSetTxPwr started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 63, OCF: 0x0026}, nil)
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
		}).Debug("VendorCypressSetTxPwr completed")
	}

	 return err
}
// VendorCypressSetBDAddrInput represents the input of the command specified in Section 7.63.121
type VendorCypressSetBDAddrInput struct {
	StaticAddress bleutil.MacAddr
}

func (i VendorCypressSetBDAddrInput) encode(data []byte) []byte {
	w := bleutil.Writer{Data: data};
	i.StaticAddress.Encode(w.Put(6))
	return w.Data
}

// VendorCypressSetBDAddrSync executes the command specified in Section 7.63.121 synchronously
func (c *Commands) VendorCypressSetBDAddrSync (params VendorCypressSetBDAddrInput) error {
	var err2 error
	var response []byte
	if c.logger != nil && c.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		c.logger.WithFields(logrus.Fields{
			 "0params": params,
		}).Trace("VendorCypressSetBDAddr started")
	}
	buffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: 63, OCF: 0x0001}, nil)
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
		}).Debug("VendorCypressSetBDAddr completed")
	}

	 return err
}
