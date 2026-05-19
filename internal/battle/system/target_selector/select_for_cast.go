package target_selector

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"battle/internal/battle/system/target_selector/target_fliter"
)

// SelectForCast 为技能效果解析目标列表。
//
// 当 castTarget 非 0、仍存在于世界且通过 [TargetSelectConfig] 筛选时：
//   - MaxCount == 1：直接采用点选目标（不再走表排序）；
//   - MaxCount > 1：在表驱动 [Select] 结果基础上保证点选目标排在首位且计入数量上限。
//
// castTarget 为 0、不存在或未通过筛选时，与 [Select] 行为一致。
func SelectForCast(w *ecs.World, caster, castTarget ecs.Entity, selectDescID int32) []ecs.Entity {
	if selectDescID == 0 {
		return nil
	}
	desc := config.GetTargetSelectConfigByID(selectDescID)
	if desc == nil || desc.MaxCount == 0 {
		return nil
	}
	if w == nil || caster == 0 || !w.EntityExists(caster) {
		return nil
	}

	if castTarget != 0 && w.EntityExists(castTarget) && passesFilters(w, caster, castTarget, desc) {
		if desc.MaxCount == 1 {
			return []ecs.Entity{castTarget}
		}
		return ensurePrimaryFirst(castTarget, Select(w, caster, selectDescID), desc.MaxCount)
	}

	return Select(w, caster, selectDescID)
}

func passesFilters(w *ecs.World, caster, target ecs.Entity, desc *config.TargetSelectConfig) bool {
	if desc == nil {
		return false
	}
	if target == caster && !desc.IncludeSelf {
		return false
	}
	ctx := &target_fliter.Context{World: w, Caster: caster, Target: target}
	return target_fliter.Apply(ctx, desc.Filters...)
}

func ensurePrimaryFirst(primary ecs.Entity, selected []ecs.Entity, maxCount int) []ecs.Entity {
	out := make([]ecs.Entity, 0, len(selected)+1)
	seen := make(map[ecs.Entity]struct{}, len(selected)+1)
	out = append(out, primary)
	seen[primary] = struct{}{}
	for _, e := range selected {
		if _, ok := seen[e]; ok {
			continue
		}
		out = append(out, e)
		seen[e] = struct{}{}
		if maxCount > 0 && len(out) >= maxCount {
			break
		}
	}
	if maxCount > 0 && len(out) > maxCount {
		return out[:maxCount]
	}
	return out
}
