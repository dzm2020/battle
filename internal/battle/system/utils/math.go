// Package utils 提供战斗 System 共用的数值比较与换算常量。
package utils

import (
	"math"
	"strings"
)

// 战斗数值换算：千分比（permille）、百分比分母等。
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
