package target_sort

import (
	"battle/ecs"
	"battle/internal/battle/utils"
)

func compareHealthCurrent(w *ecs.World, _ ecs.Entity, a, b ecs.Entity) int {
	ha := utils.HealthCurrent(w, a)
	hb := utils.HealthCurrent(w, b)
	switch {
	case ha < hb:
		return -1
	case ha > hb:
		return 1
	default:
		return 0
	}
}
