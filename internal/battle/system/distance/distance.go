// Package distance 提供实体间平面距离计算（基于 [component.Transform2D]）。
package distance

import (
	"math"

	"battle/ecs"
	"battle/internal/battle/system/transform"
)

// SquaredFromRef 计算 ref 与 e 之间欧氏距离的平方（无 Sqrt，适合排序比较）。
func SquaredFromRef(w *ecs.World, ref, e ecs.Entity) float64 {
	rx, ry, rok := transform.XY(w, ref)
	if !rok {
		rx, ry = 0, 0
	}
	x, y, ok := transform.XY(w, e)
	if !ok {
		return math.MaxFloat64
	}
	dx := float64(x - rx)
	dy := float64(y - ry)
	return dx*dx + dy*dy
}
