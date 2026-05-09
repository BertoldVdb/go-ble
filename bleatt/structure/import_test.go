package attstructure

import (
	"encoding/binary"
	"testing"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

func makeCharValue(flags byte, valueHandle uint16, charUUID bleutil.UUID) []byte {
	out := []byte{flags, 0, 0}
	binary.LittleEndian.PutUint16(out[1:], valueHandle)
	out = append(out, charUUID.UUIDToBytes()...)
	return out
}

func TestImportStructureRoundTrip(t *testing.T) {
	svcUUID := bleutil.UUIDFromStringPanic("180a")
	charUUID := bleutil.UUIDFromStringPanic("2a29")

	handles := []*GATTHandle{
		{Info: HandleInfo{Handle: 1, UUID: UUIDPrimaryService}, Value: svcUUID.UUIDToBytes()},
		{Info: HandleInfo{Handle: 2, UUID: UUIDCharacteristic}, Value: makeCharValue(0x02, 3, charUUID)},
		{Info: HandleInfo{Handle: 3, UUID: charUUID}, Value: []byte("manufacturer")},
	}

	s, err := ImportStructure(handles, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	svc := s.GetService(svcUUID)
	if svc == nil {
		t.Fatal("service not imported")
	}
	c := svc.GetCharacteristic(charUUID)
	if c == nil {
		t.Fatal("characteristic not imported")
	}
	if c.ValueHandle == nil {
		t.Fatal("value handle not set")
	}
}

func TestImportStructureRejectsIncludedService(t *testing.T) {
	handles := []*GATTHandle{
		{Info: HandleInfo{Handle: 1, UUID: UUIDIncludedService}, Value: nil},
	}
	if _, err := ImportStructure(handles, nil, nil); err == nil {
		t.Error("expected error on included-service handle")
	}
}

func TestImportStructureRejectsBareCharacteristic(t *testing.T) {
	// Characteristic without preceding service definition should error.
	charUUID := bleutil.UUIDFromStringPanic("2a29")
	handles := []*GATTHandle{
		{Info: HandleInfo{Handle: 1, UUID: UUIDCharacteristic}, Value: makeCharValue(0x02, 2, charUUID)},
	}
	if _, err := ImportStructure(handles, nil, nil); err == nil {
		t.Error("expected error on bare characteristic")
	}
}

func TestImportStructureRejectsShortCharValue(t *testing.T) {
	svcUUID := bleutil.UUIDFromStringPanic("180a")
	handles := []*GATTHandle{
		{Info: HandleInfo{Handle: 1, UUID: UUIDPrimaryService}, Value: svcUUID.UUIDToBytes()},
		{Info: HandleInfo{Handle: 2, UUID: UUIDCharacteristic}, Value: []byte{0x02}}, // <3 bytes
	}
	if _, err := ImportStructure(handles, nil, nil); err == nil {
		t.Error("expected error on short characteristic value")
	}
}

func TestImportStructureCCCD(t *testing.T) {
	svcUUID := bleutil.UUIDFromStringPanic("180a")
	charUUID := bleutil.UUIDFromStringPanic("2a29")

	handles := []*GATTHandle{
		{Info: HandleInfo{Handle: 1, UUID: UUIDPrimaryService}, Value: svcUUID.UUIDToBytes()},
		{Info: HandleInfo{Handle: 2, UUID: UUIDCharacteristic}, Value: makeCharValue(0x10, 3, charUUID)},
		{Info: HandleInfo{Handle: 3, UUID: charUUID}, Value: nil},
		{Info: HandleInfo{Handle: 4, UUID: UUIDCharacteristicClientConfiguration}, Value: []byte{0, 0}},
	}
	s, err := ImportStructure(handles, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	c := s.GetService(svcUUID).GetCharacteristic(charUUID)
	if c.ValueHandle.CCCHandle == nil {
		t.Fatal("CCCHandle not linked to characteristic")
	}
	if c.ValueHandle.CCCHandle.Info.Handle != 4 {
		t.Errorf("CCCHandle handle: got %d want 4", c.ValueHandle.CCCHandle.Info.Handle)
	}
}
