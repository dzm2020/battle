// Package combatmath 提供战斗数值比较与千分比/百分比换算常量。
package combatmath

import (
	"math"
	"strings"
)

const (
	Thousand = 1000
	Hundred  = 100
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
