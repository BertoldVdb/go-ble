package bleutil

import (
	"bytes"
	"testing"
)

// --- macaddr ---

func TestMacAddrEncodeDecode(t *testing.T) {
	cases := []struct {
		name string
		mac  MacAddr
		want []byte
	}{
		{"zero", 0, []byte{0, 0, 0, 0, 0, 0}},
		{"low byte", 0x12, []byte{0x12, 0, 0, 0, 0, 0}},
		{"all bytes", 0x563412af00ee, []byte{0xee, 0x00, 0xaf, 0x12, 0x34, 0x56}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var enc [6]byte
			c.mac.Encode(enc[:])
			if !bytes.Equal(enc[:], c.want) {
				t.Fatalf("Encode: got %x want %x", enc, c.want)
			}
			var got MacAddr
			got.Decode(c.want)
			if got != c.mac {
				t.Fatalf("Decode: got %x want %x", uint64(got), uint64(c.mac))
			}
		})
	}
}

func TestMacAddrFromString(t *testing.T) {
	cases := []struct {
		in      string
		want    MacAddr
		wantErr bool
	}{
		{"00:11:22:33:44:55", 0x001122334455, false},
		{"AA-BB-CC-DD-EE-FF", 0xAABBCCDDEEFF, false},
		{"aabbccddeeff", 0xAABBCCDDEEFF, false},
		{"too:short", 0, true},
		{"not-hex-zz", 0, true},
		{"00:11:22:33:44", 0, true},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			got, err := MacAddrFromString(c.in)
			if (err != nil) != c.wantErr {
				t.Fatalf("err: got %v wantErr=%v", err, c.wantErr)
			}
			if err == nil && got != c.want {
				t.Fatalf("got %x want %x", uint64(got), uint64(c.want))
			}
		})
	}
}

func TestMacAddrString(t *testing.T) {
	m := MacAddr(0x010203040506)
	if got, want := m.String(), "01:02:03:04:05:06"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestMacAddrFromStringPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on invalid input")
		}
	}()
	MacAddrFromStringPanic("not-a-mac")
}

func TestBLEAddrIsLessAndUint64(t *testing.T) {
	a := BLEAddr{MacAddr: 1, MacAddrType: MacAddrPublic}
	b := BLEAddr{MacAddr: 2, MacAddrType: MacAddrPublic}
	c := BLEAddr{MacAddr: 1, MacAddrType: MacAddrRandom}

	if !a.IsLess(b) {
		t.Error("a < b expected")
	}
	if b.IsLess(a) {
		t.Error("!(b < a)")
	}
	if !a.IsLess(c) {
		t.Error("a (public) < a (random) expected since type is the tiebreaker")
	}
	if a.GetUint64() == c.GetUint64() {
		t.Error("BLEAddr with different types should hash differently")
	}
	if a.Network() != "BLE" || a.MacAddrType.String() == "" {
		t.Error("trivial getters")
	}
}

func TestMacAddrTypeString(t *testing.T) {
	cases := []struct {
		t    MacAddrType
		want string
	}{
		{MacAddrPublic, "Public"},
		{MacAddrRandom, "Random"},
		{MacAddrPublicIdentity, "Public Identity"},
		{MacAddrStaticIdentity, "Static Identity"},
		{MacAddrType(99), "Unknown"},
	}
	for _, c := range cases {
		if got := c.t.String(); got != c.want {
			t.Errorf("type %d: got %q want %q", c.t, got, c.want)
		}
	}
}

// --- uuid ---

func TestUUIDFromBytesValid(t *testing.T) {
	if _, ok := UUIDFromBytesValid(nil); ok {
		t.Error("nil should be invalid")
	}
	if _, ok := UUIDFromBytesValid(make([]byte, 3)); ok {
		t.Error("3-byte should be invalid")
	}
	if _, ok := UUIDFromBytesValid(make([]byte, 8)); ok {
		t.Error("8-byte should be invalid")
	}
	if _, ok := UUIDFromBytesValid(make([]byte, 2)); !ok {
		t.Error("2-byte should be valid")
	}
	if _, ok := UUIDFromBytesValid(make([]byte, 4)); !ok {
		t.Error("4-byte should be valid")
	}
	if _, ok := UUIDFromBytesValid(make([]byte, 16)); !ok {
		t.Error("16-byte should be valid")
	}
}

func TestUUIDFromBytesPanicOnBadInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on invalid UUID length")
		}
	}()
	UUIDFromBytes(make([]byte, 7))
}

func TestUUIDFromString(t *testing.T) {
	u, err := UUIDFromString("180a")
	if err != nil {
		t.Fatal(err)
	}
	if u.GetLength() != 2 {
		t.Errorf("expected 2-byte length, got %d", u.GetLength())
	}
	if got := u.String(); got != "180a" {
		t.Errorf("String: got %q want %q", got, "180a")
	}

	u128, err := UUIDFromString("0000180a-0000-1000-8000-00805f9b34fb")
	if err != nil {
		t.Fatal(err)
	}
	// This is the 16-bit-base form — it should report length 2.
	if u128.GetLength() != 2 {
		t.Errorf("base UUID length: got %d want 2", u128.GetLength())
	}
}

func TestUUIDStringFormats(t *testing.T) {
	short := UUIDFromStringPanic("180a")
	if got := short.String(); got != "180a" {
		t.Errorf("16-bit UUID string: got %q want %q", got, "180a")
	}

	full := UUIDFromStringPanic("01020304-0506-0708-090a-0b0c0d0e0f10")
	got := full.String()
	if len(got) != 36 {
		t.Errorf("full UUID string length: got %d want 36", len(got))
	}
}

func TestUUIDFromStringInvalid(t *testing.T) {
	if _, err := UUIDFromString("not-hex"); err == nil {
		t.Error("expected error on non-hex")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("UUIDFromStringPanic expected panic")
		}
	}()
	UUIDFromStringPanic("zzz")
}

func TestUUIDIsZero(t *testing.T) {
	var z UUID
	if !z.IsZero() {
		t.Error("zero UUID should be zero")
	}
	z[5] = 1
	if z.IsZero() {
		t.Error("non-zero UUID should not be zero")
	}
}

func TestUUIDToBytes(t *testing.T) {
	short := UUIDFromStringPanic("1234")
	b := short.UUIDToBytes()
	if len(b) != 2 {
		t.Errorf("16-bit UUID bytes length: got %d want 2", len(b))
	}

	full := UUIDFromStringPanic("01020304-0506-0708-090a-0b0c0d0e0f10")
	b = full.UUIDToBytes()
	if len(b) != 16 {
		t.Errorf("full UUID bytes length: got %d want 16", len(b))
	}
}

func TestUUIDCreateVariant(t *testing.T) {
	full := UUIDFromStringPanic("01020304-0506-0708-090a-0b0c0d0e0f10")
	v := full.CreateVariant(0x42)
	if v == full {
		t.Error("variant should differ from base")
	}
	v2 := full.CreateVariantAlt(0x42)
	if v2 == full {
		t.Error("alt variant should differ from base")
	}
	if v == v2 {
		t.Error("variant and altvariant should differ")
	}
}

func TestUUIDCreateVariantPanicOn16Bit(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when creating variant of 16-bit UUID")
		}
	}()
	UUIDFromStringPanic("180a").CreateVariant(1)
}

// --- bytes ---

func TestCopySlice(t *testing.T) {
	in := []byte{1, 2, 3}
	out := CopySlice(in)
	if !bytes.Equal(in, out) {
		t.Error("CopySlice differs")
	}
	out[0] = 99
	if in[0] == 99 {
		t.Error("CopySlice did not produce an independent copy")
	}
}

func TestReverseSlice(t *testing.T) {
	cases := []struct{ in, want []byte }{
		{nil, nil},
		{[]byte{}, []byte{}},
		{[]byte{1}, []byte{1}},
		{[]byte{1, 2}, []byte{2, 1}},
		{[]byte{1, 2, 3, 4, 5}, []byte{5, 4, 3, 2, 1}},
	}
	for _, c := range cases {
		cp := append([]byte(nil), c.in...)
		ReverseSlice(cp)
		if !bytes.Equal(cp, c.want) {
			t.Errorf("ReverseSlice(%v): got %v want %v", c.in, cp, c.want)
		}
	}
}

func TestSameSlice(t *testing.T) {
	a := make([]byte, 4)
	b := a[:]
	if !SameSlice(a, b) {
		t.Error("identical slices reported different")
	}
	c := make([]byte, 4)
	if SameSlice(a, c) {
		t.Error("different backing arrays reported same")
	}
}

// --- math ---

func TestRandomRange(t *testing.T) {
	for i := 0; i < 100; i++ {
		v := RandomRange(5, 10)
		if v < 5 || v > 10 {
			t.Fatalf("out of range: %d", v)
		}
	}
	// Single-value range
	if v := RandomRange(7, 7); v != 7 {
		t.Errorf("singleton: got %d want 7", v)
	}
}

func TestClampUint16(t *testing.T) {
	cases := []struct {
		v, lo, hi, want uint16
	}{
		{0, 5, 10, 5},
		{5, 5, 10, 5},
		{7, 5, 10, 7},
		{10, 5, 10, 10},
		{100, 5, 10, 10},
	}
	for _, c := range cases {
		if got := ClampUint16(c.v, c.lo, c.hi); got != c.want {
			t.Errorf("ClampUint16(%d,%d,%d) = %d want %d", c.v, c.lo, c.hi, got, c.want)
		}
	}
}

// --- readwrite ---

func TestReaderHappyPath(t *testing.T) {
	r := Reader{Data: []byte{1, 2, 3, 4, 5}}
	if r.GetOne() != 1 {
		t.Error("GetOne first byte")
	}
	got := r.Get(2)
	if !bytes.Equal(got, []byte{2, 3}) {
		t.Errorf("Get(2): %x", got)
	}
	rest := r.GetRemainder()
	if !bytes.Equal(rest, []byte{4, 5}) {
		t.Errorf("GetRemainder: %x", rest)
	}
	if !r.Valid() {
		t.Error("Reader should still be valid")
	}
}

func TestReaderOverread(t *testing.T) {
	r := Reader{Data: []byte{1}}
	r.GetOne()                // ok
	zero := r.GetOne()        // overread
	if zero != 0 || r.Valid() {
		t.Error("Overread should return zero and mark invalid")
	}
}

func TestReaderGetTooLong(t *testing.T) {
	r := Reader{Data: []byte{1, 2, 3}}
	got := r.Get(4)
	if r.Valid() {
		t.Error("Get past end should mark invalid")
	}
	if len(got) != 4 {
		t.Errorf("Get(4) len: got %d want 4 (dummy buf)", len(got))
	}
}

func TestWriter(t *testing.T) {
	w := Writer{}
	w.PutOne(1)
	w.PutSlice([]byte{2, 3})
	more := w.Put(2)
	more[0] = 4
	more[1] = 5
	if !bytes.Equal(w.Data, []byte{1, 2, 3, 4, 5}) {
		t.Errorf("Writer: %x", w.Data)
	}
}

func TestEncodeDecodeUint24(t *testing.T) {
	for _, v := range []uint32{0, 1, 0xFF, 0xFFFF, 0xFFFFFF} {
		var b [3]byte
		EncodeUint24(b[:], v)
		got := DecodeUint24(b[:])
		if got != v {
			t.Errorf("u24 round-trip %d: got %d", v, got)
		}
	}
}

func TestEncodeUint24Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on >24-bit value")
		}
	}()
	var b [3]byte
	EncodeUint24(b[:], 0x01000000)
}

func TestCountSetBits(t *testing.T) {
	cases := []struct {
		v    uint64
		want int
	}{
		{0, 0},
		{1, 1},
		{0xFF, 8},
		{0xFFFFFFFFFFFFFFFF, 64},
		{0x01010101, 4},
	}
	for _, c := range cases {
		if got := CountSetBits(c.v); got != c.want {
			t.Errorf("CountSetBits(%#x) = %d want %d", c.v, got, c.want)
		}
	}
}

// --- log/Assert ---

func TestAssertOK(t *testing.T) {
	Assert(true, "should not panic") // smoke
}

func TestAssertPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Assert(false) should panic")
		}
	}()
	Assert(false, "boom")
}
