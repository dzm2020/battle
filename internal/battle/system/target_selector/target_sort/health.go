package target_sort

import (
	"battle/ecs"
	"battle/internal/battle/system/attrs"
)

func compareHealthCurrent(w *ecs.World, _ ecs.Entity, a, b ecs.Entity) int {
	ha := attrs.HealthCurrent(w, a)
	hb := attrs.HealthCurrent(w, b)
	switch {
	case ha < hb:
		return -1
	case ha > hb:
		return 1
	default:
		return 0
	}
}
