package resource

import (
	"battle/ecs"
	"battle/internal/battle/land"
)

// Grid 返回空间网格。
func Grid(w *ecs.World) (*land.Grid, bool) {
	g := ecs.GetResource[land.Grid](w)
	return g, g != nil
}
