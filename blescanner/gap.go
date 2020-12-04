package blescanner

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	bleutil "github.com/BertoldVdb/go-ble/util"
	"github.com/sirupsen/logrus"
)

const (
	GAPTypeFlags                = 0x01
	GAPTypeLocalNameShort       = 0x08
	GAPTypeLocalNameComplete    = 0x09
	GAPTypeTXPower              = 0x0A
	GAPTypeManufacturerSpecific = 0xFF
)

type GAPRecord struct {
	Type      uint8
	EventType EventType
	Data      []byte
}

func (g *GAPRecord) copyTo(n *GAPRecord) *GAPRecord {
	if n == nil {
		n = &GAPRecord{}
	}

	n.Type = g.Type
	n.EventType = g.EventType
	n.Data = append(n.Data[:0], g.Data...)

	return n
}

func (g GAPRecord) String() string {
	return fmt.Sprintf("%-15s %s", g.EventType, hex.EncodeToString(g.Data))
}

func (d *BLEDevice) handleRecord(event EventType, record []byte) {
	if len(record) == 0 {
		return
	}

	gapType := record[0]

	gap := d.gapSingle
	if d.gapFields != nil {
		gap = d.gapFields[gapType]
		if gap == nil {
			gap = &GAPRecord{
				Data: make([]byte, 0, len(record)),
			}
		}
	}

	gap.Type = gapType
	gap.EventType = event
	gap.Data = append(gap.Data[:0], record[1:]...)

	if d.scanner.logger != nil && d.scanner.logger.Logger.IsLevelEnabled(logrus.TraceLevel) {
		d.scanner.logger.WithFields(logrus.Fields{
			"0addr":  d.addr,
			"1event": gap.EventType,
			"2type":  gap.Type,
			"3data":  hex.EncodeToString(gap.Data),
		}).Trace("Processing GAP record")
	}

	if d.gapFields != nil {
		d.gapFields[gapType] = gap
	}

	switch gapType {
	case GAPTypeManufacturerSpecific:
		d.handleManufacturerSpecific(gap)

	case GAPTypeLocalNameShort:
		fallthrough
	case GAPTypeLocalNameComplete:
		d.handleName(gap)

	case GAPTypeFlags:
		d.handleFlags(gap)

	case GAPTypeTXPower:
		d.handleTXPower(gap)
	}

	if gapType >= 0x2 && gapType <= 0x7 {
		d.handleUUID(gap)
	}
}

func (d *BLEDevice) GetGAPTypes(result []int) []int {
	result = result[:0]

	if d.gapFields != nil {
		for key := range d.gapFields {
			result = append(result, int(key))
		}
	}

	return result
}

func (d *BLEDevice) GetGAPRecord(gapType int, buf *GAPRecord) *GAPRecord {
	if d.gapFields == nil {
		return nil
	}

	internal, ok := d.gapFields[uint8(gapType)]
	if !ok || internal == nil {
		return nil
	}

	return internal.copyTo(buf)
}

func (d *BLEDevice) handleManufacturerSpecific(gap *GAPRecord) {
	if len(gap.Data) < 2 {
		return
	}
	manufacturerID := binary.LittleEndian.Uint16(gap.Data)

	d.scanner.Lock()
	cb, ok := d.scanner.manufacturerSpecificCallback[manufacturerID]
	d.scanner.Unlock()
	if !ok || cb == nil {
		return
	}

	if d.scanner.logger != nil && d.scanner.logger.Logger.IsLevelEnabled(logrus.DebugLevel) {
		d.scanner.logger.WithFields(logrus.Fields{
			"0id":    fmt.Sprintf("0x%04x", manufacturerID),
			"1event": gap.EventType,
			"2data":  hex.EncodeToString(gap.Data),
		}).Debug("Manufacturer data callback triggered")
	}

	cb(d, gap)
}

func (d *BLEDevice) handleName(gap *GAPRecord) {
	oldState := d.nameState
	if gap.Type == GAPTypeLocalNameComplete {
		d.nameState = 2
	} else {
		/* Don't replace a full name with only a partial one... */
		if d.nameState > 1 {
			return
		}
		d.nameState = 1
	}

	if d.nameState != oldState {
		d.name = string(gap.Data)

		if d.scanner.logger != nil {
			d.scanner.logger.WithFields(logrus.Fields{
				"0addr": d.addr,
				"1name": d.name,
			}).Info("Found device name")
		}
	}
}

func (d *BLEDevice) GetNameType() (string, int) {
	return d.name, d.nameState
}

func (d *BLEDevice) GetName() string {
	name, _ := d.GetNameType()
	return name
}

func (d *BLEDevice) handleFlags(gap *GAPRecord) {
	if len(gap.Data) >= 1 {
		d.flags = gap.Data[0]
	}
}

func (d *BLEDevice) GetFlags() uint8 {
	return d.flags
}

func (d *BLEDevice) handleTXPower(gap *GAPRecord) {
	if len(gap.Data) >= 1 {
		d.txPower = int8(gap.Data[0])
	}
}

func (d *BLEDevice) GetTXPower() int8 {
	return d.txPower
}

func (d *BLEDevice) handleUUID(gap *GAPRecord) {
	complete := false
	if gap.Type%2 == 1 {
		complete = true
	}

	l := 0
	tt := gap.Type>>1 - 1
	switch tt {
	case 0:
		l = 2
	case 1:
		l = 4
	case 2:
		l = 16
	default:
		panic("Unsupported UUID type supplied")
	}

	if complete {
		d.services[tt] = d.services[tt][:0]
	}

	data := gap.Data
	for {
		if len(data) < l {
			if len(data) != 0 {
				if d.scanner.logger != nil {
					d.scanner.logger.WithFields(logrus.Fields{
						"0addr": d.addr,
						"1name": d.name,
					}).Debug("UUID beacon incomplete")
				}
			}
			return
		}
		value := data[:l]
		data = data[l:]

		uuid := bleutil.UUIDFromBytes(value)

		found := false
		for _, m := range d.services[tt] {
			if m == uuid {
				found = true
				break
			}
		}

		if !found {
			if len(d.services[tt]) < 32 {
				d.services[tt] = append(d.services[tt], uuid)
			}
		}
	}
}

func (d *BLEDevice) GetServices(serviceType int, in []bleutil.UUID) []bleutil.UUID {
	in = in[:0]

	if serviceType < 0 {
		in = append(in, d.services[0]...)
		in = append(in, d.services[1]...)
		in = append(in, d.services[2]...)
	} else {
		in = append(in, d.services[serviceType]...)
	}

	return in
}
