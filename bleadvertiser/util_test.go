package bleadvertiser

import (
	"bytes"
	"testing"
)

func TestUtilPDUAddRecord(t *testing.T) {
	out := UtilPDUAddRecord(nil, 0x01, []byte{0x06})
	if !bytes.Equal(out, []byte{0x02, 0x01, 0x06}) {
		t.Errorf("flags record: %x", out)
	}

	out = UtilPDUAddRecord(out, 0x09, []byte("test"))
	want := []byte{0x02, 0x01, 0x06, 0x05, 0x09, 't', 'e', 's', 't'}
	if !bytes.Equal(out, want) {
		t.Errorf("name record: %x want %x", out, want)
	}

	// Empty payload: still emits length=1 + type byte.
	out = UtilPDUAddRecord(nil, 0xff, nil)
	if !bytes.Equal(out, []byte{0x01, 0xff}) {
		t.Errorf("empty payload: %x", out)
	}
}

func TestLegacyAdvertisementTypeString(t *testing.T) {
	cases := []struct {
		t    LegacyAdvertisementType
		want string
	}{
		{LegacyAdvertisementTypeInd, "ADV_IND"},
		{LegacyAdvertisementTypeDirectInd, "ADV_DIRECT_IND"},
		{LegacyAdvertisementTypeScanInd, "ADV_SCAN_IND"},
		{LegacyAdvertisementTypeNonConnInd, "ADV_NONCONN_IND"},
		{LegacyAdvertisementTypeDirectIndLowDuty, "ADV_DIRECT_IND_LOW_DUTY"},
		{LegacyAdvertisementType(99), "Invalid"},
	}
	for _, c := range cases {
		if got := c.t.String(); got != c.want {
			t.Errorf("type %d: got %q want %q", c.t, got, c.want)
		}
	}
}
