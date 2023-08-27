package attstructure

import (
	"encoding/binary"
	"errors"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

func ImportStructure(gattHandles []*GATTHandle, rCb ClientReadHandler, wCb ClientWriteHandler) (*Structure, error) {
	s := &Structure{
		isClient: true,

		clientRead:      rCb,
		clientWrite:     wCb,
		clientNotifyMap: make(map[uint16](ClientNotifyHandler)),
	}

	var currentService *Service
	var currentCharacteristic *Characteristic

	finishCharacteristic := func() {
		if currentCharacteristic != nil && currentService != nil {
			if currentCharacteristic.ValueHandle != nil {
				currentService.characteristics = append(currentService.characteristics, currentCharacteristic)
			}
			currentCharacteristic = nil
		}
	}
	finishService := func() {
		if currentService != nil {
			finishCharacteristic()

			s.services = append(s.services, currentService)
			currentCharacteristic = nil
		}
	}

	for _, m := range gattHandles {
		if m.Info.UUID == UUIDIncludedService {
			return nil, errors.New("including services is not supported (TODO)")
		}

		if m.Info.UUID == UUIDPrimaryService || m.Info.UUID == UUIDSecondaryService {
			finishService()

			uuid, ok := bleutil.UUIDFromBytesValid(m.Value)
			if !ok {
				return nil, errors.New("Service UUID cannot be parsed")
			}

			/* It is theoretically allowed to split a service definition,
			   I wonder when you would use it... */
			currentService = nil
			for _, k := range s.services {
				if k.uuid == uuid {
					currentService = k
					break
				}
			}
			if currentService == nil {
				currentService = &Service{
					isPrimary: m.Info.UUID == UUIDPrimaryService,
					parent:    s,
					uuid:      uuid,
				}
			}
		} else {
			if currentService == nil {
				return nil, errors.New("this item must be in a service definition")
			}

			/* The standard says you can have multiple characteristic definitions with
			   the same UUID. I don't know what that means in practice (different permissions?).
			   We just add them all, the code that uses this can decide which one to access */
			if m.Info.UUID == UUIDCharacteristic {
				finishCharacteristic()

				if len(m.Value) < 3 {
					return nil, errors.New("Characteristic definition has an invalid length")
				}

				flags := CharacteristicFlag(m.Value[0])
				handle := binary.LittleEndian.Uint16(m.Value[1:])
				uuid, ok := bleutil.UUIDFromBytesValid(m.Value[3:])
				if !ok {
					return nil, errors.New("Characteristic UUID cannot be parsed")
				}

				currentCharacteristic = &Characteristic{
					parent:      currentService,
					uuid:        uuid,
					flags:       flags,
					valueIsNext: true,
					ValueHandle: &GATTHandle{
						Info: HandleInfo{
							Handle:    handle,
							UUID:      uuid,
							UUIDWidth: uuid.GetLength(),
							Flags:     flags,
						},
					},
				}

				continue
			}

			/* The Characteristic Value declaration contains the value of the characteristic. It
			 * is the first Attribute after the characteristic declaration. All characteristic
			 * definitions shall have a Characteristic Value declaration.
			 *
			 * I have seen a (TELink based) device that puts garbage information in the value handle and UUID above, so
			 * we need to parse it again here. Hoever, I have also seen some devices that did not value declarations in the
			 * full scan, so use both mechanisms.
			 */
			if currentCharacteristic != nil {
				if currentCharacteristic.valueIsNext {
					currentCharacteristic.valueIsNext = false
					currentCharacteristic.uuid = m.Info.UUID
					currentCharacteristic.ValueHandle.Info.UUID = m.Info.UUID
					currentCharacteristic.ValueHandle.Info.UUIDWidth = m.Info.UUIDWidth
					currentCharacteristic.ValueHandle.Info.Handle = m.Info.Handle
				}

				if m.Info.UUID == UUIDCharacteristicClientConfiguration {
					currentCharacteristic.ValueHandle.CCCHandle = m
				}
			}
		}
	}

	finishService()

	return s, nil
}
