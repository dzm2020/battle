package utils

import (
	"math"
	"strings"
)

func CompareFloat64(cur float64, op string, val float64) bool {
	switch strings.TrimSpace(op) {
	case ">":
		return cur > val
	case "<":
		return cur < val
	case "==", "=":
		return math.Abs(cur-val) < 1e-6
	case "!=":
		return math.Abs(cur-val) >= 1e-6
	case ">=":
		return cur >= val
	case "<=":
		return cur <= val
	default:
		return false
	}
}
