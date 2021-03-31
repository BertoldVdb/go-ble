package servicedeviceinformation

import (
	attperipheral "github.com/BertoldVdb/go-ble/bleatt/helpers/peripheral"
	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	hciconnmgr "github.com/BertoldVdb/go-ble/hci/connmgr"
	bleutil "github.com/BertoldVdb/go-ble/util"
)

type DeviceInformationConfig struct {
	ManufacturerName   string
	ModelNumber        string
	SerialNumber       string
	HardwareRevision   string
	FirmwareRevision   string
	SoftwareRevision   string
	SystemID           string
	RegulatoryDataList []byte
	PnPID              []byte
}

func DefaultConfig() *DeviceInformationConfig {
	return &DeviceInformationConfig{}
}

type DeviceInformation struct {
	config *DeviceInformationConfig
}

func (s *DeviceInformation) CreateStructure(structure *attstructure.Structure) error {
	pdi := structure.AddPrimaryService(bleutil.UUIDFromStringPanic("180A"))

	register := func(uuid string, data []byte) {
		if len(data) > 0 {
			pdi.AddCharacteristicReadOnly(bleutil.UUIDFromStringPanic(uuid), data)
		}
	}

	register("2A29", []byte(s.config.ManufacturerName))
	register("2A24", []byte(s.config.ModelNumber))
	register("2A25", []byte(s.config.SerialNumber))
	register("2A27", []byte(s.config.HardwareRevision))
	register("2A26", []byte(s.config.FirmwareRevision))
	register("2A28", []byte(s.config.SoftwareRevision))
	register("2A23", []byte(s.config.SystemID))
	register("2A2A", s.config.RegulatoryDataList)
	register("2A50", s.config.PnPID)

	return nil
}

func (s *DeviceInformation) Disconnected() {
}

func (s *DeviceInformation) Connected(conn hciconnmgr.BufferConn) error {
	return nil
}

func CreateService(config *DeviceInformationConfig) func() attperipheral.PeripheralImplementation {
	return func() attperipheral.PeripheralImplementation {
		return &DeviceInformation{
			config: config,
		}
	}
}
