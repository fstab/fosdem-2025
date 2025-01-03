package util

import (
	"math"
	"math/rand"
	"time"
)

func Sleep() {
	time.Sleep(time.Duration(math.Abs(rand.NormFloat64()*50+50)) * time.Millisecond)
}
