package bleutil

import (
	"math/rand"

	_ "github.com/BertoldVdb/go-misc/seed"
)

func RandomRange(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func ClampUint16(value uint16, min uint16, max uint16) uint16 {
	if value > max {
		return max
	}
	if value < min {
		return min
	}
	return value
}
