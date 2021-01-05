package deviceinfo

import (
	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	bleutil "github.com/BertoldVdb/go-ble/util"
)

type ControllerInfo struct {
	SupportedCommands       *hcicommands.InformationalReadLocalSupportedCommandsOutput
	SupportedFeatures       *hcicommands.InformationalReadLocalSupportedFeaturesOutput
	LESupportedFeatures     *hcicommands.LEReadLocalSupportedFeaturesOutput
	BdAddr                  *hcicommands.InformationalReadBDADDROutput
	RandomAddr              bleutil.MacAddr
	LocalVersionInformation *hcicommands.InformationalReadLocalVersionInformationOutput
}

func (c *ControllerInfo) Read(cmds *hcicommands.Commands) error {
	var err error

	c.BdAddr, err = cmds.InformationalReadBDADDRSync(nil)
	if err != nil {
		return err
	}

	c.SupportedCommands, err = cmds.InformationalReadLocalSupportedCommandsSync(nil)
	if err != nil {
		return err
	}

	c.SupportedFeatures, err = cmds.InformationalReadLocalSupportedFeaturesSync(nil)
	if err != nil {
		return err
	}

	c.LocalVersionInformation, err = cmds.InformationalReadLocalVersionInformationSync(nil)
	if err != nil {
		return err
	}

	c.LESupportedFeatures, err = cmds.LEReadLocalSupportedFeaturesSync(nil)
	return err
}
