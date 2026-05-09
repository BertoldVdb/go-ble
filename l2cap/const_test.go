package blel2cap

import (
	"errors"
	"testing"
)

func TestPSMTypeValues(t *testing.T) {
	// Pin the wire-level constants — accidentally changing these
	// would silently break interop with peers.
	if PSMTypeATT != 0x7 {
		t.Errorf("PSMTypeATT: got %#x want 0x7", PSMTypeATT)
	}
	if PSMTypeSecurityManager != 0x10000 {
		t.Errorf("PSMTypeSecurityManager: got %#x want 0x10000", PSMTypeSecurityManager)
	}
}

func TestSignallingErrorToUint16(t *testing.T) {
	if got := signallingErrorToUint16(nil, 0, 1); got != 0 {
		t.Errorf("nil error → success code; got %d", got)
	}
	if got := signallingErrorToUint16(errors.New("nope"), 0, 1); got != 1 {
		t.Errorf("non-nil error → failure code; got %d", got)
	}
}

func TestSignallingTokenSignalCompleteSafe(t *testing.T) {
	// signalComplete must be non-blocking and safe across multiple
	// invocations + Cleanup (uses sync.Once + select-with-default).
	tok := &signallingCommandToken{
		completeChan: make(chan struct{}, 1),
	}
	tok.signalComplete()
	tok.signalComplete() // second call: select-default, no panic
	tok.Cleanup()
	tok.Cleanup() // second close: sync.Once short-circuits
	tok.signalComplete() // post-close: select-default avoids send-on-closed
}
