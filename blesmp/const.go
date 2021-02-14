package blesmp

type smpOpcode byte

const (
	opcodePairingRequest  smpOpcode = 0x1
	opcodePairingResponse smpOpcode = 0x2
	opcodePairingConfirm  smpOpcode = 0x3
	opcodePairingRandom   smpOpcode = 0x4
	opcodePairingFailed   smpOpcode = 0x5

	opcodeKDEncryptionInformation      smpOpcode = 0x6
	opcodeKDIdentification             smpOpcode = 0x7
	opcodeKDIdentityInformation        smpOpcode = 0x8
	opcodeKDIdentityAddressInformation smpOpcode = 0x9
	opcodeKDSigningInformation         smpOpcode = 0xA

	opcodeKDSecurityRequest smpOpcode = 0xB

	opcodePairingPublicKey     smpOpcode = 0xC
	opcodePairingDHKeyCheck    smpOpcode = 0xD
	opcodeKeypressNotification smpOpcode = 0xE
)

type smpFailedReason byte

const (
	failedPasskeyEntryFailed  smpFailedReason = 0x01
	failedConfirmValueFailed  smpFailedReason = 0x4
	failedPairingNotSupported smpFailedReason = 0x5
	failedEncryptionKeySize   smpFailedReason = 0x6
	failedCommandNotSupported smpFailedReason = 0x7
	failedUnspecifiedReason   smpFailedReason = 0x8
)

type SMPState int

const (
	StateInsecure          SMPState = 0
	StateFailed            SMPState = 1
	StateSecure            SMPState = 2
	StateBusy              SMPState = 4
	StatePermanentlyFailed SMPState = 5
)

type smpIOCapability byte

const (
	cIODisplayOnly     smpIOCapability = 0
	cIODisplayYesNo    smpIOCapability = 1
	cIOKeyboardOnly    smpIOCapability = 2
	cIONoInputNoOutput smpIOCapability = 3
	cIOKeyboardDisplay smpIOCapability = 4
)
