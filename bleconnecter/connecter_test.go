package bleconnecter

import (
	"context"
	"testing"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

// Connect must reject empty/nil peer lists. Previously a nil list was
// silently treated as "match any peer", letting any peer satisfy a
// pending Connect — a fundamental authorization gap, particularly in
// peripheral mode where the LL allowlist is disabled.
func TestConnectRejectsNilPeers(t *testing.T) {
	c := &BLEConnecter{}
	_, _, err := c.Connect(context.Background(), true, nil, BLEConnectionParametersRequested{})
	if err != ErrorNoPeers {
		t.Fatalf("nil peers: got %v want %v", err, ErrorNoPeers)
	}
}

func TestConnectRejectsEmptyPeers(t *testing.T) {
	c := &BLEConnecter{}
	_, _, err := c.Connect(context.Background(), false, []bleutil.BLEAddr{}, BLEConnectionParametersRequested{})
	if err != ErrorNoPeers {
		t.Fatalf("empty peers: got %v want %v", err, ErrorNoPeers)
	}
}

func TestMakeValidClampsIntervalToFloor(t *testing.T) {
	r := BLEConnectionParametersRequested{
		ConnectionIntervalMin: 1,
		ConnectionIntervalMax: 1,
	}
	r.makeValid()
	if r.ConnectionIntervalMin != 0x6 || r.ConnectionIntervalMax != 0x6 {
		t.Errorf("interval not clamped to spec floor: min=%d max=%d", r.ConnectionIntervalMin, r.ConnectionIntervalMax)
	}
}

func TestMakeValidClampsIntervalToCeiling(t *testing.T) {
	r := BLEConnectionParametersRequested{
		ConnectionIntervalMin: 0xFFFF,
		ConnectionIntervalMax: 0xFFFF,
	}
	r.makeValid()
	if r.ConnectionIntervalMin != 0xC80 || r.ConnectionIntervalMax != 0xC80 {
		t.Errorf("interval not clamped to spec ceiling: min=%d max=%d", r.ConnectionIntervalMin, r.ConnectionIntervalMax)
	}
}

// If user supplies min > max, makeValid should normalise.
func TestMakeValidNormalisesMinGtMax(t *testing.T) {
	r := BLEConnectionParametersRequested{
		ConnectionIntervalMin: 0x100,
		ConnectionIntervalMax: 0x40,
	}
	r.makeValid()
	if r.ConnectionIntervalMin > r.ConnectionIntervalMax {
		t.Errorf("min > max after makeValid: min=%d max=%d", r.ConnectionIntervalMin, r.ConnectionIntervalMax)
	}
}

func TestMakeValidLatencyClamped(t *testing.T) {
	r := BLEConnectionParametersRequested{
		ConnectionLatency: 0xFFFF,
	}
	r.makeValid()
	if r.ConnectionLatency != 0x1F3 {
		t.Errorf("latency not clamped: got %d want 0x1F3", r.ConnectionLatency)
	}
}

// SupervisionTimeout must be at least the spec minimum (10) and at
// least (1 + Latency) * Interval * 2.
func TestMakeValidSupervisionTimeoutFloor(t *testing.T) {
	r := BLEConnectionParametersRequested{
		ConnectionIntervalMin: 0x10,
		ConnectionIntervalMax: 0x10,
		ConnectionLatency:     0,
		SupervisionTimeout:    0,
	}
	r.makeValid()
	if r.SupervisionTimeout < 10 {
		t.Errorf("supervision timeout below floor: %d", r.SupervisionTimeout)
	}
}

func TestMakeValidSupervisionTimeoutCeiling(t *testing.T) {
	r := BLEConnectionParametersRequested{
		ConnectionIntervalMin: 0x10,
		ConnectionIntervalMax: 0x10,
		SupervisionTimeout:    0xFFFF,
	}
	r.makeValid()
	if r.SupervisionTimeout != 0xC80 {
		t.Errorf("supervision timeout above ceiling: got %d want 0xC80", r.SupervisionTimeout)
	}
}

// BLEConnection trivial getters
func TestBLEConnectionIsCentral(t *testing.T) {
	c := &BLEConnection{isCentral: true}
	if !c.IsCentral() {
		t.Error("IsCentral should return true")
	}
	c.isCentral = false
	if c.IsCentral() {
		t.Error("IsCentral should return false")
	}
}

func TestBLEConnectionRemoteAddr(t *testing.T) {
	c := &BLEConnection{
		peerAddr: bleutil.BLEAddr{
			MacAddr:     0x010203040506,
			MacAddrType: bleutil.MacAddrPublic,
		},
	}
	got := c.RemoteAddr()
	if got == nil {
		t.Fatal("RemoteAddr returned nil")
	}
	if got.Network() != "BLE" {
		t.Errorf("Network: got %q want BLE", got.Network())
	}
}
