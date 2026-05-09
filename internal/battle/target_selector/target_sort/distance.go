package target_sort

import (
	"battle/ecs"
	"battle/internal/battle/utils"
)

func compareDistanceSquared(w *ecs.World, ref, a, b ecs.Entity) int {
	da := utils.DistanceSquaredFromRef(w, ref, a)
	db := utils.DistanceSquaredFromRef(w, ref, b)
	switch {
	case da < db:
		return -1
	case da > db:
		return 1
	default:
		return 0
	}
}
