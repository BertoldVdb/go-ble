package bleutil

type Writer struct {
	Data []byte
}

func (w *Writer) PutOne(v uint8) {
	w.Data = append(w.Data, byte(v))
}

func (w *Writer) Put(cnt int) []byte {
	for i := 0; i < cnt; i++ {
		w.Data = append(w.Data, 0)
	}

	return w.Data[len(w.Data)-cnt:]
}

func (w *Writer) PutSlice(sl []byte) {
	w.Data = append(w.Data, sl...)
}

type Reader struct {
	Data  []byte
	index int
	fail  bool
}

var dummyBuf [512]byte

func (r *Reader) GetOne() uint8 {
	if r.index+1 > len(r.Data) {
		r.fail = true
		return 0
	}
	tmp := r.Data[r.index]
	r.index++
	return tmp
}

func (r *Reader) Get(cnt int) []byte {
	if r.index+cnt > len(r.Data) {
		r.fail = true
		return dummyBuf[0:cnt]
	}
	tmp := r.Data[r.index : r.index+cnt]
	r.index += cnt
	return tmp
}

func (r *Reader) GetRemainder() []byte {
	if r.index > len(r.Data) {
		r.fail = true
		return nil
	}

	return r.Data[r.index:]
}

func (r *Reader) Valid() bool {
	return !r.fail
}

func EncodeUint24(data []byte, value uint32) {
	data[0] = byte(value >> 0)
	data[1] = byte(value >> 8)
	data[2] = byte(value >> 16)
	if value>>24 > 0 {
		panic("Value does not fit in 24-bit number")
	}
}

func DecodeUint24(data []byte) uint32 {
	result := uint32(0)
	result |= uint32(data[0]) << 0
	result |= uint32(data[1]) << 8
	result |= uint32(data[2]) << 16
	return result
}

func CountSetBits(input uint64) int {
	result := 0
	for i := 0; i < 64; i++ {
		if input&1 > 0 {
			result++
		}
		input >>= 1
	}
	return result
}
