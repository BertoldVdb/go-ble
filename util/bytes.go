package bleutil

import (
	"reflect"
	"unsafe"
)

func CopySlice(in []byte) []byte {
	out := make([]byte, len(in))

	copy(out, in)

	return out
}

func ReverseSlice(in []byte) {
	l := len(in)
	for i := 0; i < len(in)/2; i++ {
		in[i], in[l-1-i] = in[l-1-i], in[i]
	}
}

func SameSlice(s1, s2 []byte) bool {
	h1 := (*reflect.SliceHeader)(unsafe.Pointer(&s1))
	h2 := (*reflect.SliceHeader)(unsafe.Pointer(&s2))

	return h1.Data == h2.Data && h1.Len == h2.Len
}
