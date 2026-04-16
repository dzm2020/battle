package geom

import "math"

// Vec2 战斗平面坐标（第 9 天 AOI 可复用）。仅做几何计算，不参与业务判断。
type Vec2 struct {
	X float64
	Y float64
}

// DistSq 平方距离，避免开方比较时可用 DistSq <= r*r。
func DistSq(a, b Vec2) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return dx*dx + dy*dy
}

// Dist 欧氏距离。
func Dist(a, b Vec2) float64 {
	return math.Sqrt(DistSq(a, b))
}
