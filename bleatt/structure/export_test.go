package attstructure

import (
	"testing"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

// CCCD permission derivation must reflect the characteristic's overall
// confidentiality, not only its read-encryption flag. A notify-only,
// write-encrypted characteristic must produce a CCCD whose write
// permission requires encryption — otherwise an unencrypted peer can
// subscribe to a sensitive notification stream.
func TestExportCCCDInheritsWriteEncryption(t *testing.T) {
	s := NewStructure()
	svc := s.AddPrimaryService(bleutil.UUIDFromStringPanic("180a"))
	svc.AddCharacteristic(
		bleutil.UUIDFromStringPanic("2a00"),
		CharacteristicNotify|CharacteristicWriteAck|CharacteristicWriteNeedsEncryption,
		ValueConfig{},
	)

	exp := &ExportedStructure{}
	exp.Append(s)

	var ccc *GATTHandle
	for _, h := range exp.Handles {
		if h.Info.UUID == UUIDCharacteristicClientConfiguration {
			ccc = h
			break
		}
	}
	if ccc == nil {
		t.Fatal("no CCCD generated")
	}
	if ccc.Info.Flags&CharacteristicWriteNeedsEncryption == 0 {
		t.Errorf("CCCD write must require encryption when characteristic is write-encrypted; got flags=%#x", ccc.Info.Flags)
	}
	if ccc.Info.Flags&CharacteristicReadNeedsEncryption == 0 {
		t.Errorf("CCCD read must require encryption when characteristic is encrypted; got flags=%#x", ccc.Info.Flags)
	}
}

// Conversely, a plain (no encryption flags) notify characteristic should
// produce a plain CCCD.
func TestExportCCCDPlainCharacteristic(t *testing.T) {
	s := NewStructure()
	svc := s.AddPrimaryService(bleutil.UUIDFromStringPanic("180a"))
	svc.AddCharacteristic(
		bleutil.UUIDFromStringPanic("2a00"),
		CharacteristicNotify|CharacteristicWriteAck,
		ValueConfig{},
	)

	exp := &ExportedStructure{}
	exp.Append(s)

	var ccc *GATTHandle
	for _, h := range exp.Handles {
		if h.Info.UUID == UUIDCharacteristicClientConfiguration {
			ccc = h
			break
		}
	}
	if ccc == nil {
		t.Fatal("no CCCD generated")
	}
	if ccc.Info.Flags&(CharacteristicReadNeedsEncryption|CharacteristicWriteNeedsEncryption) != 0 {
		t.Errorf("CCCD on plain characteristic must not require encryption; got flags=%#x", ccc.Info.Flags)
	}
}
