package skill

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
)

// ExecuteEffects 对目标列表依次执行技能效果链（伤害 / 治疗 / 挂 Buff）。
func ExecuteEffects(w *ecs.World, caster ecs.Entity, targets []ecs.Entity, sk SkillConfig) {
	for _, ef := range sk.Effects {
		for _, t := range targets {
			switch ef.Kind {
			case EffectDamage:
				component.MergePendingDamage(w, t, ef.Amount, ef.DamageType, caster)
			case EffectHeal:
				component.MergePendingHeal(w, t, ef.Amount, caster)
			case EffectApplyBuff:
				buff.ApplyBuff(w, t, ef.BuffDefID)
			}
		}
	}
}
