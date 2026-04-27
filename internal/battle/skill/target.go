package skill

import (
	"battle/ecs"
	"battle/internal/battle/target_selector"
)

// ResolveSkillEffectTargets 根据 [config.SkillEffectConfig.TargetSelectID] 解析效果目标列表。
// selectID==0：仅使用显式主目标（若无效则返回空）。
func ResolveSkillEffectTargets(w *ecs.World, caster, mainTarget ecs.Entity, selectID int) []ecs.Entity {
	return target_selector.SelectTargets(w, caster, mainTarget, int32(selectID))
}
