// Package transform 提供 [component.Transform2D] 的 World 级读取。
package transform

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// XY 读取实体平面坐标。
func XY(w *ecs.World, e ecs.Entity) (int, int, bool) {
	t, ok := w.GetComponent(e, &component.Transform2D{})
	if !ok {
		return 0, 0, false
	}
	tr := t.(*component.Transform2D)
	return tr.X, tr.Y, true
}
