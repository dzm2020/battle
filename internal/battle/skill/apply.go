package skill

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"battle/internal/battle/target_selector"
)

// ApplySkillEffects 按 [config.SkillBaseConfig.EffectIDs] 顺序执行技能效果（伤害、加 Buff 等）。
// caster：施法者；mainTarget：施法请求中的主目标（选目标 ID 为 0 时使用）。
func ApplySkillEffects(w *ecs.World, caster, mainTarget ecs.Entity, skillID int) {
	if w == nil || caster == 0 || !w.EntityExists(caster) {
		return
	}

	desc := config.GetSkillConfigByID(int32(skillID))

	if desc == nil {
		return
	}
	for _, eid := range desc.EffectIDs {
		effectDesc := config.GetSkillEffectConfigByID(int32(eid))
		if effectDesc == nil {
			continue
		}
		//  选取目标
		targets := target_selector.Select(w, caster, int32(effectDesc.TargetSelectID))
		//  执行效果
		for _, t := range targets {
			applySkillEffect(w, caster, t, effectDesc)
		}
	}
}
