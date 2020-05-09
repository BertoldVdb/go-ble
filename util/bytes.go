package bleutil

func ReverseSlice(in []byte) {
	l := len(in)
	for i := 0; i < len(in)/2; i++ {
		in[i], in[l-1-i] = in[l-1-i], in[i]
	}
}
