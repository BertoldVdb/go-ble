package bleadvertiser

type LegacyAdvertisementType uint8

const (
	LegacyAdvertisementTypeInd              LegacyAdvertisementType = 0
	LegacyAdvertisementTypeDirectInd        LegacyAdvertisementType = 1
	LegacyAdvertisementTypeScanInd          LegacyAdvertisementType = 2
	LegacyAdvertisementTypeNonConnInd       LegacyAdvertisementType = 3
	LegacyAdvertisementTypeDirectIndLowDuty LegacyAdvertisementType = 4
)

func (a LegacyAdvertisementType) String() string {
	switch a {
	case LegacyAdvertisementTypeInd:
		return "ADV_IND"
	case LegacyAdvertisementTypeDirectInd:
		return "ADV_DIRECT_IND"
	case LegacyAdvertisementTypeScanInd:
		return "ADV_SCAN_IND"
	case LegacyAdvertisementTypeNonConnInd:
		return "ADV_NONCONN_IND"
	case LegacyAdvertisementTypeDirectIndLowDuty:
		return "ADV_DIRECT_IND_LOW_DUTY"
	}

	return "Invalid"
}

func UtilPDUAddRecord(pdu []byte, rt uint8, data []byte) []byte {
	pdu = append(pdu, uint8(len(data)+1))
	pdu = append(pdu, rt)
	return append(pdu, data...)
}
