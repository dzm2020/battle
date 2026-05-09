package utils

import (
	"battle/ecs"
	"math"
)

// DistanceSquaredFromRef
//
//	@Description: 计算两个对象之间距离
//	@param w
//	@param ref
//	@param e
//	@return float64
func DistanceSquaredFromRef(w *ecs.World, ref, e ecs.Entity) float64 {
	rx, ry, rok := TransformXY(w, ref)
	if !rok {
		rx, ry = 0, 0
	}
	x, y, ok := TransformXY(w, e)
	if !ok {
		return math.MaxFloat64
	}
	dx := x - rx
	dy := y - ry
	return math.Sqrt(dx*dx + dy*dy)
}
