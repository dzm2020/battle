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
	tab := config.Tab

	desc := tab.SkillConfigByID[int32(skillID)]

	if desc == nil {
		return
	}
	for _, eid := range desc.EffectIDs {
		eff := lookupSkillEffect(tab, int32(eid))
		if eff == nil {
			continue
		}
		targets := target_selector.SelectTargets(w, caster, mainTarget, int32(eff.TargetSelectID))
		for _, t := range targets {
			applySkillEffect(w, caster, t, eff)
		}
	}
}

func lookupSkillEffect(tab *config.Tables, id int32) *config.SkillEffectConfig {
	if tab == nil || tab.SkillEffectConfigByID == nil {
		return nil
	}
	return tab.SkillEffectConfigByID[id]
}
