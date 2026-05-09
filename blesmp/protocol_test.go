package blesmp

import "testing"

// getLegacyAlgorithmType picks JustWorks (0,0) when neither side requests
// MITM, and an authenticated method (passkey) when MITM is requested and
// IO capabilities allow it. Verify that NoInputNoOutput on both sides
// always falls back to JustWorks even when MITM is requested — that
// fallback is what the audit flagged as a silent downgrade; we keep the
// behaviour but document it via tests so future changes are deliberate.
func TestGetLegacyAlgorithmTypeNoMITM(t *testing.T) {
	a, b := getLegacyAlgorithmType(cIODisplayOnly, cIOKeyboardOnly, false, false)
	if a != 0 || b != 0 {
		t.Fatalf("no MITM should be JustWorks (0,0); got (%d,%d)", a, b)
	}
}

func TestGetLegacyAlgorithmTypeKeyboardDisplay(t *testing.T) {
	// Initiator has keyboard, responder has display — initiator inputs
	// the passkey shown on the responder.
	a, b := getLegacyAlgorithmType(cIOKeyboardOnly, cIODisplayOnly, true, false)
	if a != 1 || b != 2 {
		t.Fatalf("keyboard/display: got (%d,%d) want (1,2)", a, b)
	}
}

func TestGetLegacyAlgorithmTypeDisplayKeyboard(t *testing.T) {
	a, b := getLegacyAlgorithmType(cIODisplayOnly, cIOKeyboardOnly, false, true)
	if a != 2 || b != 1 {
		t.Fatalf("display/keyboard: got (%d,%d) want (2,1)", a, b)
	}
}

func TestGetLegacyAlgorithmTypeKeyboardKeyboard(t *testing.T) {
	a, b := getLegacyAlgorithmType(cIOKeyboardOnly, cIOKeyboardOnly, true, true)
	if a != 1 || b != 1 {
		t.Fatalf("keyboard/keyboard: got (%d,%d) want (1,1)", a, b)
	}
}

func TestGetLegacyAlgorithmTypeNoIOFallback(t *testing.T) {
	// Both sides NoInputNoOutput, even with MITM requested → JustWorks.
	// Documented as a current intentional fallback.
	a, b := getLegacyAlgorithmType(cIONoInputNoOutput, cIONoInputNoOutput, true, true)
	if a != 0 || b != 0 {
		t.Fatalf("NoIO/NoIO with MITM: got (%d,%d) want (0,0)", a, b)
	}
}

// All defined SMP failure reason constants must have distinct values.
// The audit noted that several (AuthenticationRequirements, OOBNotAvailable,
// DHKeyCheckFailed, NumericComparisonFailed) were missing and have now
// been added — this test pins the values.
func TestSMPFailureReasonsDistinct(t *testing.T) {
	reasons := map[smpFailedReason]string{
		failedPasskeyEntryFailed:         "PasskeyEntryFailed",
		failedOOBNotAvailable:            "OOBNotAvailable",
		failedAuthenticationRequirements: "AuthenticationRequirements",
		failedConfirmValueFailed:         "ConfirmValueFailed",
		failedPairingNotSupported:        "PairingNotSupported",
		failedEncryptionKeySize:          "EncryptionKeySize",
		failedCommandNotSupported:        "CommandNotSupported",
		failedUnspecifiedReason:          "UnspecifiedReason",
		failedRepeatedAttempts:           "RepeatedAttempts",
		failedInvalidParameters:          "InvalidParameters",
		failedDHKeyCheckFailed:           "DHKeyCheckFailed",
		failedNumericComparisonFailed:    "NumericComparisonFailed",
	}
	if len(reasons) != 12 {
		t.Fatalf("duplicate failure-reason values: %v", reasons)
	}
	// Spot-check a few well-known values per Core Spec.
	if failedAuthenticationRequirements != 0x03 {
		t.Errorf("AuthenticationRequirements: got %#x want 0x03", failedAuthenticationRequirements)
	}
	if failedDHKeyCheckFailed != 0x0B {
		t.Errorf("DHKeyCheckFailed: got %#x want 0x0B", failedDHKeyCheckFailed)
	}
}

// TestHandleKeyDistributionAcceptsExtraOpcodes verifies that valid IRK,
// IdentityAddress and CSRK PDUs from the peer are accepted (returning
// true) rather than triggering a `failedCommandNotSupported`. This is
// the key fix for the audit finding "post-pairing failure on legitimate
// IRK/CSRK distribution".
func TestHandleKeyDistributionAcceptsExtraOpcodes(t *testing.T) {
	c := &SMPConn{}

	// IRK is 16 bytes
	if !c.handleKeyDistribution(opcodeKDIdentityInformation, make([]byte, 16)) {
		t.Error("IRK distribution should be accepted")
	}
	// IdentityAddress is 7 bytes (1 type + 6 addr)
	if !c.handleKeyDistribution(opcodeKDIdentityAddressInformation, make([]byte, 7)) {
		t.Error("IdentityAddress distribution should be accepted")
	}
	// CSRK is 16 bytes
	if !c.handleKeyDistribution(opcodeKDSigningInformation, make([]byte, 16)) {
		t.Error("CSRK distribution should be accepted")
	}
	// Length mismatch must be rejected.
	if c.handleKeyDistribution(opcodeKDIdentityInformation, make([]byte, 15)) {
		t.Error("malformed IRK length should be rejected")
	}
}
