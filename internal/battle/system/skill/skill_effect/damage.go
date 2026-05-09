package skill_effect

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// handleSkillEffectDamage 造成伤害：写入 [component.PendingDamage]，由 [system.DamageSystem] 结算。
// IntParams[0]：伤害量；IntParams[1] 可选：[component.DamageType]（0 物 / 1 法 / 2 真），缺省为物理。
func handleSkillEffectDamage(w *ecs.World, caster, target ecs.Entity, eff *config.SkillEffectConfig) {
	if len(eff.IntParams) < 1 {
		return
	}
	amt := eff.IntParams[0]
	if amt <= 0 {
		return
	}
	dt := component.DamagePhysical
	if len(eff.IntParams) >= 2 {
		switch eff.IntParams[1] {
		case int(component.DamageMagical):
			dt = component.DamageMagical
		case int(component.DamageTrue):
			dt = component.DamageTrue
		default:
			dt = component.DamagePhysical
		}
	}
	component.MergePendingDamage(w, target, amt, dt, caster)
}
