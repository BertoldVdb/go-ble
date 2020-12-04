package bleutil

import (
	"github.com/BertoldVdb/go-misc/bufferfifo"
	pdu "github.com/BertoldVdb/go-misc/pdubuf"
)

var rxtxFreeBuffers = bufferfifo.New(64)

const headerCap = 64

func GetBufferCap(length int, capacity int) *pdu.PDU {
	buf := rxtxFreeBuffers.Pop()

	if buf == nil {
		buf = pdu.Alloc(headerCap, length, capacity)
	} else {
		buf.Realloc(headerCap, length, capacity)
	}

	buf.SetState(0, 1)

	return buf
}

func GetBuffer(length int) *pdu.PDU {
	return GetBufferCap(length, 2*length)
}

func CopyBufferFromSlice(slice []byte) *pdu.PDU {
	pdu := GetBufferCap(len(slice), cap(slice))
	copy(pdu.Buf(), slice)
	return pdu
}

func CopySliceFromBuffer(buf *pdu.PDU) []byte {
	sl := make([]byte, buf.Len())
	copy(sl, buf.Buf())

	return sl
}

func ReleaseBuffer(buf *pdu.PDU) int {
	if buf == nil {
		return 0
	}
	buf.SetState(1, 0)
	return rxtxFreeBuffers.Push(buf)
}
