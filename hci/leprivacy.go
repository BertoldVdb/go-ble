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
	return 0

	if c.Info.RandomAddr == 0 {
		return 0
	}

	if usage == LEAddrUsageConnect {
		//TODO: add a system for using quirks depending on the host controller
		//Cambridge Silicon Radio dongles don't seem to connect properly when using the random address.
		return 0
	}

	return 1
}

func (c *Controller) setLERandomAddress() error {
	/* Read 48 random bits */
	var rnd [8]byte
	_, err := crypto_rand.Read(rnd[2:])
	if err != nil {
		return err
	}
	value := binary.BigEndian.Uint64(rnd[:])

	/* Randomize last x bits */
	c.Info.RandomAddr = c.Info.BdAddr.BDADDR
	c.Info.RandomAddr ^= bleutil.MacAddr(value >> (48 - c.config.LERandomAddrBits))

	/* Not all devices support the random address so failure is not fatal */
	err = c.Cmds.LESetRandomAddressSync(hcicommands.LESetRandomAddressInput{
		RandomAddess: c.Info.RandomAddr,
	})
	c.logger.WithError(err).WithField("0addr", c.Info.RandomAddr).Info("Set LE random address")
	if err != nil {
		c.Info.RandomAddr = 0
	}
	return nil
}
