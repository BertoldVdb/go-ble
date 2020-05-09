package blescanner

type EventType uint8

const (
	EventTypeInd        EventType = 0
	EventTypeDirectInd  EventType = 1
	EventTypeScanInd    EventType = 2
	EventTypeNonConnInd EventType = 3
	EventTypeScanRsp    EventType = 4
)

func (a EventType) String() string {
	switch a {
	case EventTypeInd:
		return "ADV_IND"
	case EventTypeDirectInd:
		return "ADV_DIRECT_IND"
	case EventTypeScanInd:
		return "ADV_SCAN_IND"
	case EventTypeNonConnInd:
		return "ADV_NONCONN_IND"
	case EventTypeScanRsp:
		return "SCAN_RSP"
	}

	return "Invalid"
}
