package target_sort

import (
	"battle/ecs"
	"battle/internal/battle/system/distance"
)

func compareDistanceSquared(w *ecs.World, ref, a, b ecs.Entity) int {
	da := distance.SquaredFromRef(w, ref, a)
	db := distance.SquaredFromRef(w, ref, b)
	switch {
	case da < db:
		return -1
	case da > db:
		return 1
	default:
		return 0
	}
}
