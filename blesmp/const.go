package blesmp

type psmOpcode byte

const (
	opcodePairingRequest  psmOpcode = 0x1
	opcodePairingResponse psmOpcode = 0x2
	opcodePairingConfirm  psmOpcode = 0x3
	opcodePairingRandom   psmOpcode = 0x4
	opcodePairingFailed   psmOpcode = 0x5

	opcodeKDEncryptionInformation      psmOpcode = 0x6
	opcodeKDInitiatorIdentification    psmOpcode = 0x7
	opcodeKDIdentityInformation        psmOpcode = 0x8
	opcodeKDIdentityAddressInformation psmOpcode = 0x9
	opcodeKDSigningInformation         psmOpcode = 0xA

	opcodeKDSecurityRequest psmOpcode = 0xB

	opcodePairingPublicKey     psmOpcode = 0xC
	opcodePairingDHKeyCheck    psmOpcode = 0xD
	opcodeKeypressNotification psmOpcode = 0xE
)

type psmFailedReason byte

const (
	failedConfirmValueFailed  psmFailedReason = 0x4
	failedEncryptionKeySize   psmFailedReason = 0x6
	failedCommandNotSupported psmFailedReason = 0x7
	failedUnspecifiedReason   psmFailedReason = 0x8
)
