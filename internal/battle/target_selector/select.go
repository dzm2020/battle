package target_selector

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"battle/internal/battle/target_selector/target_fliter"
	"battle/internal/battle/target_selector/target_sort"
)

// Select 根据 [config.TargetSelectConfigByID] 选取实体列表。
func Select(w *ecs.World, caster ecs.Entity, selectDescID int32) []ecs.Entity {
	if w == nil || caster == 0 || !w.EntityExists(caster) {
		return nil
	}
	desc := config.GetTargetSelectConfigByID(selectDescID)
	if desc == nil {
		return nil
	}
	if desc.MaxCount == 0 {
		return nil
	}
	var candidates []ecs.Entity
	seen := make(map[ecs.Entity]struct{})
	w.ForEachEntity(func(e ecs.Entity) {
		if e == caster && !desc.IncludeSelf {
			return
		}
		ctx := &target_fliter.Context{World: w, Caster: caster, Target: e}
		if !target_fliter.Apply(ctx, desc.Filters...) {
			return
		}
		if _, ok := seen[e]; ok {
			return
		}
		seen[e] = struct{}{}
		candidates = append(candidates, e)
	})
	target_sort.Apply(w, caster, candidates, desc.SortType, desc.SortOrder)
	if desc.MaxCount > 0 && len(candidates) > desc.MaxCount {
		return candidates[:desc.MaxCount]
	}
	return candidates
}
