package utils

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/land"
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
	dx := float64(x - rx)
	dy := float64(y - ry)
	return math.Sqrt(dx*dx + dy*dy)
}

func GetLandFreeCell(grid *land.Grid, camp component.SideType) (cellX, cellZ int, ok bool) {
	if grid == nil {
		return 0, 0, false
	}
	return grid.PickFreeCell(camp == component.SideTypeRed)
}
