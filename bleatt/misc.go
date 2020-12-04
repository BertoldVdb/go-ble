package bleatt

import (
	attstructure "github.com/BertoldVdb/go-ble/bleatt/structure"
	bleutil "github.com/BertoldVdb/go-ble/util"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
)

func getOpcode(buf *pdu.PDU) (bool, ATTCommand, bool, bool) {
	if buf.Len() < 1 {
		return false, 0, false, false
	}

	opcode := buf.Buf()[0]
	isAuthenticated := (opcode >> 7) == 1
	method := ATTCommand(opcode & 0x3F)
	isForServer := opcode&1 == 0

	if method == ATTHandleValueCNF {
		isForServer = false
	}

	return true, method, isAuthenticated, isForServer
}

func isPartOfGATTDatabase(uuid bleutil.UUID) int {
	if uuid == attstructure.UUIDPrimaryService || uuid == attstructure.UUIDSecondaryService || uuid == attstructure.UUIDIncludedService ||
		uuid == attstructure.UUIDCharacteristic || uuid == attstructure.UUIDCharacteristicExtendedProperties {
		return 2 /* Include value in hash */
	}
	if uuid == attstructure.UUIDCharacteristicUserDescription || uuid == attstructure.UUIDCharacteristicClientConfiguration ||
		uuid == attstructure.UUIDCharacteristicServerConfiguration || uuid == attstructure.UUIDCharacteristicFormat || uuid == attstructure.UUIDCharacteristicAggregateFormat {
		return 1
	}

	return 0
}
