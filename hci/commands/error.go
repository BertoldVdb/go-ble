package hcicommands

import "errors"

var hciErrors [0x46]error

type HCIError struct {
	code int
	text string
}

func (h HCIError) Error() string {
	return h.text
}

func (h HCIError) Code() int {
	return h.code
}

var (
	ErrorMalformed = errors.New("Parsing reply failed")
)

func init() {
	makeError := func(code int, text string) {
		hciErrors[code] = HCIError{
			code: code,
			text: text,
		}
	}

	makeError(0x01, "Unknown HCI Command")
	makeError(0x02, "Unknown Connection Identifier")
	makeError(0x03, "Hardware Failure")
	makeError(0x04, "Page Timeout")
	makeError(0x05, "Authentication Failure")
	makeError(0x06, "PIN or Key Missing")
	makeError(0x07, "Memory Capacity Exceeded")
	makeError(0x08, "Connection Timeout")
	makeError(0x09, "Connection Limit Exceeded")
	makeError(0x0A, "Synchronous Connection Limit To A Device Exceeded")
	makeError(0x0B, "Connection Already Exists")
	makeError(0x0C, "Command Disallowed")
	makeError(0x0D, "Connection Rejected due to Limited Resources")
	makeError(0x0E, "Connection Rejected Due To Security Reasons")
	makeError(0x0F, "Connection Rejected due to Unacceptable BD_ADDR")
	makeError(0x10, "Connection Accept Timeout Exceeded")
	makeError(0x11, "Unsupported Feature or Parameter Value")
	makeError(0x12, "Invalid HCI Command Parameters")
	makeError(0x13, "Remote User Terminated Connection")
	makeError(0x14, "Remote Device Terminated Connection due to Low Resources")
	makeError(0x15, "Remote Device Terminated Connection due to Power Off")
	makeError(0x16, "Connection Terminated By Local Host")
	makeError(0x17, "Repeated Attempts")
	makeError(0x18, "Pairing Not Allowed")
	makeError(0x19, "Unknown LMP PDU")
	makeError(0x1A, "Unsupported Remote Feature / Unsupported LMP Feature")
	makeError(0x1B, "SCO Offset Rejected")
	makeError(0x1C, "SCO Interval Rejected")
	makeError(0x1D, "SCO Air Mode Rejected")
	makeError(0x1E, "Invalid LMP Parameters / Invalid LL Parameters")
	makeError(0x1F, "Unspecified Error")
	makeError(0x20, "Unsupported LMP Parameter Value / Unsupported LL Parameter Value")
	makeError(0x21, "Role Change Not Allowed")
	makeError(0x22, "LMP Response Timeout / LL Response Timeout")
	makeError(0x23, "LMP Error Transaction Collision / LL Procedure Collision")
	makeError(0x24, "LMP PDU Not Allowed")
	makeError(0x25, "Encryption Mode Not Acceptable")
	makeError(0x26, "Link Key cannot be Changed")
	makeError(0x27, "Requested QoS Not Supported")
	makeError(0x28, "Instant Passed")
	makeError(0x29, "Pairing With Unit Key Not Supported")
	makeError(0x2A, "Different Transaction Collision")
	makeError(0x2B, "Reserved for future use")
	makeError(0x2C, "QoS Unacceptable Parameter")
	makeError(0x2D, "QoS Rejected")
	makeError(0x2E, "Channel Classification Not Supported")
	makeError(0x2F, "Insufficient Security")
	makeError(0x30, "Parameter Out Of Mandatory Range")
	makeError(0x31, "Reserved for future use")
	makeError(0x32, "Role Switch Pending")
	makeError(0x33, "Reserved for future use")
	makeError(0x34, "Reserved Slot Violation")
	makeError(0x35, "Role Switch Failed")
	makeError(0x36, "Extended Inquiry Response Too Large")
	makeError(0x37, "Secure Simple Pairing Not Supported By Host")
	makeError(0x38, "Host Busy - Pairing")
	makeError(0x39, "Connection Rejected due to No Suitable Channel Found")
	makeError(0x3A, "Controller Busy")
	makeError(0x3B, "Unacceptable Connection Parameters")
	makeError(0x3C, "Advertising Timeout")
	makeError(0x3D, "Connection Terminated due to MIC Failure")
	makeError(0x3E, "Connection Failed to be Established / Synchronization Timeout")
	makeError(0x3F, "MAC Connection Failed")
	makeError(0x40, "Coarse Clock Adjustment Rejected but Will Try to Adjust Using Clock Dragging")
	makeError(0x41, "Type 0 Submap Not Defined")
	makeError(0x42, "Unknown Advertising Identifier")
	makeError(0x43, "Limit Reached")
	makeError(0x44, "Operation Cancelled by Host")
	makeError(0x45, "Packet Too Long")
}

func HciErrorToGo(result []byte, err error) error {
	if len(result) == 0 {
		return ErrorMalformed
	}

	value := int(result[0])
	if value == 0 {
		return err
	}

	if value >= len(hciErrors) {
		return ErrorMalformed
	}
	return hciErrors[value]
}
