package blescanner

import "testing"

func TestEventTypeString(t *testing.T) {
	cases := []struct {
		e    EventType
		want string
	}{
		{EventTypeInd, "ADV_IND"},
		{EventTypeDirectInd, "ADV_DIRECT_IND"},
		{EventTypeScanInd, "ADV_SCAN_IND"},
		{EventTypeNonConnInd, "ADV_NONCONN_IND"},
		{EventTypeScanRsp, "SCAN_RSP"},
		{EventType(99), "Invalid"},
	}
	for _, c := range cases {
		if got := c.e.String(); got != c.want {
			t.Errorf("type %d: got %q want %q", c.e, got, c.want)
		}
	}
}
