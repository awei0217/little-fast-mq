package utils

import (
	"math/rand"
	"time"
)

func init() {

	rand.Seed(time.Now().Unix())
}

func GetRandUint32() uint32 {

	return rand.Uint32()
}
