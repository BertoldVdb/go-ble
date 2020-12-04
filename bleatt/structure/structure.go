package attstructure

import (
	"encoding/binary"
	"sync"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

type Structure struct {
	services []*PrimaryService
	exported *ExportedStructure
}

type PrimaryService struct {
	parent          *Structure
	uuid            bleutil.UUID
	characteristics []*Characteristic
}

type Characteristic struct {
	parent       *PrimaryService
	uuid         bleutil.UUID
	flags        CharacteristicFlag
	initialValue []byte

	valueConfig ValueConfig

	ValueHandle *GATTHandle
}

type ValueConfig struct {
	ValueBeforeReadCb func(h *GATTHandle, offset int) error
	ValueAfterReadCb  func(h *GATTHandle, offset int, bytes int) error
	ValueWriteCb      func(h *GATTHandle) error
	LengthFixed       bool
	LengthMax         uint16
}

func NewStructure() *Structure {
	return &Structure{}
}

func (s *Structure) AddPrimaryService(uuid bleutil.UUID) *PrimaryService {
	p := &PrimaryService{
		parent: s,
		uuid:   uuid,
	}
	s.services = append(s.services, p)
	return p
}

func (p *PrimaryService) AddCharacteristic(uuid bleutil.UUID, flags CharacteristicFlag, valueConfig ValueConfig) *Characteristic {
	c := &Characteristic{
		parent: p,
		uuid:   uuid,
		flags:  flags,

		valueConfig: valueConfig,
	}

	p.characteristics = append(p.characteristics, c)

	return c
}

func (p *PrimaryService) AddCharacteristicReadOnly(uuid bleutil.UUID, value []byte) *Characteristic {
	c := p.AddCharacteristic(uuid, CharacteristicRead, ValueConfig{})
	c.initialValue = value
	return c
}

type HandleInfo struct {
	Handle         uint16
	GroupEndHandle uint16
	UUIDWidth      int
	UUID           bleutil.UUID

	Flags CharacteristicFlag
}

type GATTHandle struct {
	Info HandleInfo

	Value       []byte
	ValueConfig ValueConfig

	CCCHandle *GATTHandle
}

type ExportedStructure struct {
	sync.Mutex
	idx     uint16
	Handles []*GATTHandle

	HandleSet func(*Characteristic, []byte) (int, error)
}

func (c *Characteristic) SetValue(new []byte) (int, error) {
	e := c.parent.parent.exported
	e.Lock()

	if c.valueConfig.LengthFixed {
		copy(c.ValueHandle.Value, new)
	} else {
		c.ValueHandle.Value = append(c.ValueHandle.Value[:0], new...)
	}

	e.Unlock()

	return e.HandleSet(c, new)

}

func (c *Characteristic) GetValue(buf []byte) []byte {
	e := c.parent.parent.exported
	e.Lock()
	defer e.Unlock()

	return append(buf[:0], c.ValueHandle.Value...)
}

func (result *ExportedStructure) Append(s *Structure) {
	result.Lock()
	defer result.Unlock()

	s.exported = result

	for _, p := range s.services {
		result.idx++

		serviceDescr := &GATTHandle{
			Info: HandleInfo{Handle: result.idx,
				UUIDWidth: UUIDPrimaryService.GetLength(),
				UUID:      UUIDPrimaryService,
				Flags:     CharacteristicRead,
			},
			Value: p.uuid.UUIDToBytes(),
		}

		result.Handles = append(result.Handles, serviceDescr)

		for _, c := range p.characteristics {
			result.idx++
			charDescrValue := []byte{byte(c.flags), 0, 0}
			charDescrValue = append(charDescrValue, c.uuid.UUIDToBytes()...)
			binary.LittleEndian.PutUint16(charDescrValue[1:], result.idx+1)

			charDescr := &GATTHandle{
				Info: HandleInfo{Handle: result.idx,
					UUIDWidth: UUIDCharacteristic.GetLength(),
					UUID:      UUIDCharacteristic,
					Flags:     CharacteristicRead,
				},
				Value: charDescrValue,
			}

			result.Handles = append(result.Handles, charDescr)

			result.idx++
			var valueCopy []byte
			if c.valueConfig.LengthFixed {
				valueCopy = make([]byte, c.valueConfig.LengthMax)
			} else {
				valueCopy = make([]byte, len(c.initialValue))
			}
			copy(valueCopy, c.initialValue)

			charValue := &GATTHandle{
				Info: HandleInfo{Handle: result.idx,
					UUIDWidth: c.uuid.GetLength(),
					UUID:      c.uuid,
					Flags:     c.flags,
				},
				Value:       valueCopy,
				ValueConfig: c.valueConfig,
			}
			c.ValueHandle = charValue
			result.Handles = append(result.Handles, charValue)

			/* Does it need a CCC? */
			if c.flags&(CharacteristicIndicate|CharacteristicIndicate) > 0 {
				result.idx++
				charValue.CCCHandle = &GATTHandle{
					Info: HandleInfo{Handle: result.idx,
						UUIDWidth: UUIDCharacteristicClientConfiguration.GetLength(),
						UUID:      UUIDCharacteristicClientConfiguration,
						Flags:     CharacteristicRead | CharacteristicWriteAck,
					},
					Value: []byte{0, 0},
				}

				result.Handles = append(result.Handles, charValue.CCCHandle)
			}
		}

		serviceDescr.Info.GroupEndHandle = result.idx
	}
}
