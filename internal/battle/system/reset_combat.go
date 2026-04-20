package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// ClearCombatEntities 移除所有带 [Team] 的实体（战斗结束清理）。
func ClearCombatEntities(w *ecs.World) {
	q := ecs.NewQuery[*component.Team](w)
	var ents []ecs.Entity
	q.ForEach(func(e ecs.Entity, _ *component.Team) {
		ents = append(ents, e)
	})
	for _, e := range ents {
		w.RemoveEntity(e)
	}
}
