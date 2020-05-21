package bleutil

import (
	"math/rand"

	_ "github.com/BertoldVdb/go-misc/seed"
)

func RandomRange(min int, max int) int {
	return min + rand.Intn(max-min+1)
}
