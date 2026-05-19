// Package distance 提供实体间平面距离计算（基于 [component.Transform2D]）。
package distance

import (
	"math"

	"battle/ecs"
	"battle/internal/battle/system/attrs"
)

// FromRef 计算 ref 与 e 之间的欧氏距离。
func FromRef(w *ecs.World, ref, e ecs.Entity) float64 {
	rx, ry, rok := attrs.TransformXY(w, ref)
	if !rok {
		rx, ry = 0, 0
	}
	x, y, ok := attrs.TransformXY(w, e)
	if !ok {
		return math.MaxFloat64
	}
	dx := float64(x - rx)
	dy := float64(y - ry)
	return math.Sqrt(dx*dx + dy*dy)
}
