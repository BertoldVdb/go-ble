package hci

import (
	"encoding/binary"

	crypto_rand "crypto/rand"

	hcicommands "github.com/BertoldVdb/go-ble/hci/commands"
	bleutil "github.com/BertoldVdb/go-ble/util"
)

type LEAddrUsage int

const (
	LEAddrUsageScan      LEAddrUsage = 0
	LEAddrUsageConnect   LEAddrUsage = 1
	LEAddrUsageAdvertise LEAddrUsage = 2
)

func (c *Controller) GetLERecommenedOwnAddrType(usage LEAddrUsage) bleutil.MacAddrType {
	if c.Info.RandomAddr == 0 {
		return 0
	}

	if c.config.PrivacyScan && usage == LEAddrUsageScan {
		return 1
	}
	if c.config.PrivacyConnect && usage == LEAddrUsageConnect {
		return 1
	}
	if c.config.PrivacyAdvertise && usage == LEAddrUsageAdvertise {
		return 1
	}

	return 0
}

func (c *Controller) GetOwnAddress(t bleutil.MacAddrType) bleutil.BLEAddr {
	result := bleutil.BLEAddr{
		MacAddrType: t,
	}

	switch t {
	case bleutil.MacAddrPublic:
		result.MacAddr = c.Info.BdAddr.BDADDR
	case bleutil.MacAddrRandom:
		result.MacAddr = c.Info.RandomAddr
	default:
		panic("Unsupported address type")
	}

	return result
}

func (c *Controller) setLERandomAddress() error {
	/* Read 48 random bits */
	var value uint64

	for {
		var rnd [8]byte
		_, err := crypto_rand.Read(rnd[2:])
		if err != nil {
			return err
		}
		value = binary.BigEndian.Uint64(rnd[:])

		/* Very unlikely, may not be zero or all one */
		if value != 0 && value != 0xFFFFFFFFFFFF {
			break
		}
	}

	/* Randomize last x bits */
	c.Info.RandomAddr = c.Info.BdAddr.BDADDR
	c.Info.RandomAddr ^= bleutil.MacAddr(value >> (48 - c.config.LERandomAddrBits))

	/* Two MSB must be one */
	c.Info.RandomAddr |= 0x3 << (48 - 2)

	/* Not all devices support the random address so failure is not fatal */
	err := c.Cmds.LESetRandomAddressSync(hcicommands.LESetRandomAddressInput{
		RandomAddess: c.Info.RandomAddr,
	})
	c.logger.WithError(err).WithField("0addr", c.Info.RandomAddr).Info("Set LE random address")
	if err != nil {
		c.Info.RandomAddr = 0
	}
	return nil
}
