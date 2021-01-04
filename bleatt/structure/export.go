package attstructure

import (
	"encoding/binary"
	"sync"
)

type ExportedStructure struct {
	sync.Mutex
	idx     uint16
	Handles []*GATTHandle

	HandleSet func(*Characteristic, []byte) (int, error)
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
