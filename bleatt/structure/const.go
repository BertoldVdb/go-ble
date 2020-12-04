package attstructure

import (
	bleutil "github.com/BertoldVdb/go-ble/util"
)

var (
	UUIDPrimaryService   = bleutil.UUIDFromStringPanic("2800")
	UUIDSecondaryService = bleutil.UUIDFromStringPanic("2801")
	UUIDIncludedService  = bleutil.UUIDFromStringPanic("2802")
	UUIDCharacteristic   = bleutil.UUIDFromStringPanic("2803")

	UUIDCharacteristicExtendedProperties  = bleutil.UUIDFromStringPanic("2900")
	UUIDCharacteristicUserDescription     = bleutil.UUIDFromStringPanic("2901")
	UUIDCharacteristicClientConfiguration = bleutil.UUIDFromStringPanic("2902")
	UUIDCharacteristicServerConfiguration = bleutil.UUIDFromStringPanic("2903")
	UUIDCharacteristicFormat              = bleutil.UUIDFromStringPanic("2904")
	UUIDCharacteristicAggregateFormat     = bleutil.UUIDFromStringPanic("2905")

	UUIDDatabaseHash = bleutil.UUIDFromStringPanic("2b2a")
)

type CharacteristicFlag uint16

const (
	CharacteristicBroadcast   CharacteristicFlag = 0x1
	CharacteristicRead        CharacteristicFlag = 0x2
	CharacteristicWriteNoAck  CharacteristicFlag = 0x4
	CharacteristicWriteAck    CharacteristicFlag = 0x8
	CharacteristicNotify      CharacteristicFlag = 0x10
	CharacteristicIndicate    CharacteristicFlag = 0x20
	CharacteristicSignedWrite CharacteristicFlag = 0x40
	CharacteristicExtended    CharacteristicFlag = 0x80
)
