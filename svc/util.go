package svc

import (
	"math"
)

// IsPositive checks if a given float64 value is positive, excluding NaN and infinity.
func IsPositive(amt float64) bool {
	return amt > 0 && !math.IsNaN(amt) && !math.IsInf(amt, 0)
}
