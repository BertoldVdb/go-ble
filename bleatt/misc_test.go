package bleatt

import (
	"testing"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

func TestGetOpcodePreservesCommandOpcodes(t *testing.T) {
	tests := []struct {
		name        string
		opcode      byte
		method      ATTCommand
		auth        bool
		isForServer bool
	}{
		{
			name:        "write request",
			opcode:      byte(ATTWriteReq),
			method:      ATTWriteReq,
			isForServer: true,
		},
		{
			name:        "write command",
			opcode:      byte(ATTWriteCMD),
			method:      ATTWriteCMD,
			isForServer: true,
		},
		{
			name:        "signed write command",
			opcode:      byte(ATTSignedWriteCMD),
			method:      ATTSignedWriteCMD,
			auth:        true,
			isForServer: true,
		},
		{
			name:        "handle value confirmation",
			opcode:      byte(ATTHandleValueCNF),
			method:      ATTHandleValueCNF,
			isForServer: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bleutil.GetBuffer(1)
			defer bleutil.ReleaseBuffer(buf)
			buf.Buf()[0] = tt.opcode

			valid, method, auth, isForServer := getOpcode(buf)
			if !valid {
				t.Fatal("opcode was not valid")
			}
			if method != tt.method || auth != tt.auth || isForServer != tt.isForServer {
				t.Fatalf("got method=%#x auth=%v isForServer=%v", method, auth, isForServer)
			}
		})
	}
}

func TestNormalizeATTMTU(t *testing.T) {
	if got := normalizeATTMTU(0); got != 0xFFFF {
		t.Fatalf("zero MTU should mean default maximum, got %d", got)
	}
	if got := normalizeATTMTU(1); got != 23 {
		t.Fatalf("short MTU should clamp to 23, got %d", got)
	}
	if got := normalizeATTMTU(128); got != 128 {
		t.Fatalf("valid MTU changed to %d", got)
	}
}

func TestAddPayloadHandlesTooSmallMTU(t *testing.T) {
	buf := bleutil.GetBuffer(3)
	defer bleutil.ReleaseBuffer(buf)

	conn := &gattDeviceConn{mtu: 1}
	full, copied := (&attServer{}).addPayload(conn, buf, []byte{1, 2, 3})
	if !full || copied != 0 || buf.Len() != 3 {
		t.Fatalf("got full=%v copied=%d len=%d", full, copied, buf.Len())
	}
}
