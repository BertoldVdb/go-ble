package hcicommands

type writer struct {
	data []byte
}

func (w *writer) PutOne(v uint8) {
	w.data = append(w.data, byte(v))
}

func (w *writer) Put(cnt int) []byte {
	for i := 0; i < cnt; i++ {
		w.data = append(w.data, 0)
	}

	return w.data[len(w.data)-cnt:]
}

func (w *writer) PutSlice(sl []byte) {
	w.data = append(w.data, sl...)
}

func (w *writer) Data() []byte {
	return w.data
}

type reader struct {
	data  []byte
	index int
	fail  bool
}

var dummyBuf [512]byte

func (r *reader) GetOne() uint8 {
	if r.index+1 > len(r.data) {
		r.fail = true
		return 0
	}
	tmp := r.data[r.index]
	r.index++
	return tmp
}

func (r *reader) Get(cnt int) []byte {
	if r.index+cnt > len(r.data) {
		r.fail = true
		return dummyBuf[0:cnt]
	}
	tmp := r.data[r.index : r.index+cnt]
	r.index += cnt
	return tmp
}

func (r *reader) GetRemainder() []byte {
	if r.index > len(r.data) {
		r.fail = true
		return nil
	}

	return r.data[r.index:]
}

func (r *reader) Valid() bool {
	return !r.fail
}

func encodeUint24(data []byte, value uint32) {
	data[0] = byte(value >> 0)
	data[1] = byte(value >> 8)
	data[2] = byte(value >> 16)
	if value>>24 > 0 {
		panic("Value does not fit in 24-bit number")
	}
}

func decodeUint24(data []byte) uint32 {
	result := uint32(0)
	result |= uint32(data[0]) << 0
	result |= uint32(data[1]) << 8
	result |= uint32(data[2]) << 16
	return result
}

func countSetBits(input uint64) int {
	result := 0
	for i := 0; i < 64; i++ {
		if input&1 > 0 {
			result++
		}
		input >>= 1
	}
	return result
}
