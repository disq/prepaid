package svc

import (
	"math"
)

func IsPositive(amt float64) bool {
	return amt > 0 && !math.IsNaN(amt) && !math.IsInf(amt, 0)
}
