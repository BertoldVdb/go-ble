package attstructure

import (
	"testing"

	bleutil "github.com/BertoldVdb/go-ble/util"
)

func TestStructureBuildAndExport(t *testing.T) {
	s := NewStructure()
	if got := s.GetServices(); len(got) != 0 {
		t.Errorf("empty structure: got %d services", len(got))
	}

	svc := s.AddPrimaryService(bleutil.UUIDFromStringPanic("180a"))
	c := svc.AddCharacteristic(
		bleutil.UUIDFromStringPanic("2a29"),
		CharacteristicRead,
		ValueConfig{},
	)
	if got := c.GetFlags(); got != CharacteristicRead {
		t.Errorf("flags: got %#x want %#x", got, CharacteristicRead)
	}
	if got := s.GetService(bleutil.UUIDFromStringPanic("180a")); got != svc {
		t.Error("GetService did not find primary service")
	}
	if got := s.GetService(bleutil.UUIDFromStringPanic("180b")); got != nil {
		t.Error("GetService returned non-nil for unknown UUID")
	}

	if got := svc.GetCharacteristic(bleutil.UUIDFromStringPanic("2a29")); got != c {
		t.Error("GetCharacteristic did not find by UUID")
	}
	if got := svc.GetCharacteristic(bleutil.UUIDFromStringPanic("2a2a")); got != nil {
		t.Error("GetCharacteristic returned non-nil for unknown UUID")
	}

	exp := &ExportedStructure{}
	exp.Append(s)
	if len(exp.Handles) == 0 {
		t.Fatal("ExportedStructure has no handles")
	}

	// Verify the structure prints something nontrivial via its String().
	if got := s.String(); len(got) == 0 {
		t.Error("Structure.String() empty")
	}
}

func TestExportedStructureMultipleServices(t *testing.T) {
	s := NewStructure()
	svc1 := s.AddPrimaryService(bleutil.UUIDFromStringPanic("1800"))
	svc1.AddCharacteristicReadOnly(bleutil.UUIDFromStringPanic("2a00"), []byte("Name"))
	svc2 := s.AddPrimaryService(bleutil.UUIDFromStringPanic("180a"))
	svc2.AddCharacteristic(
		bleutil.UUIDFromStringPanic("2a29"),
		CharacteristicRead|CharacteristicNotify,
		ValueConfig{},
	)

	exp := &ExportedStructure{}
	exp.Append(s)

	// Each service produces 1 service descriptor; each characteristic
	// produces 2 handles (descriptor + value); a Notify characteristic
	// produces an additional CCCD.
	wantHandles := 1 + 2 + 1 + 2 + 1
	if got := len(exp.Handles); got != wantHandles {
		t.Errorf("handles: got %d want %d", got, wantHandles)
	}
}

func TestImportStructureInvalid(t *testing.T) {
	// Empty handles should still produce a Structure (or error gracefully).
	_, err := ImportStructure(nil, nil, nil)
	_ = err // either result is acceptable; main goal is it doesn't panic
}

func TestServiceAddCharacteristicReadOnly(t *testing.T) {
	s := NewStructure()
	svc := s.AddPrimaryService(bleutil.UUIDFromStringPanic("180a"))
	c := svc.AddCharacteristicReadOnly(bleutil.UUIDFromStringPanic("2a29"), []byte{1, 2, 3})

	if c.GetFlags() != CharacteristicRead {
		t.Errorf("flags: got %#x want %#x", c.GetFlags(), CharacteristicRead)
	}

	chars := svc.GetCharacteristics()
	if len(chars) != 1 || chars[0] != c {
		t.Errorf("GetCharacteristics: %+v", chars)
	}
}
