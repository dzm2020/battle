package skill

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

func init() {
	registerSkillEffect(config.EffectDamage, handleSkillEffectDamage)
	registerSkillEffect(config.EffectAddBuff, handleSkillEffectAddBuff)
}

type skillEffectFn func(w *ecs.World, caster, target ecs.Entity, eff *config.SkillEffectConfig)

var skillEffectDict = make(map[config.EffectType]skillEffectFn)

func registerSkillEffect(typ config.EffectType, fn skillEffectFn) {
	skillEffectDict[typ] = fn
}

func applySkillEffect(w *ecs.World, caster, target ecs.Entity, eff *config.SkillEffectConfig) {
	if w == nil || eff == nil || target == 0 || !w.EntityExists(target) {
		return
	}
	if fn := skillEffectDict[eff.EffectType]; fn != nil {
		fn(w, caster, target, eff)
	}
}

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

// handleSkillEffectAddBuff 添加 Buff：走 [buff.AddBuff]，模板来自全局 [config.Tab.BuffConfigConfigByID]。
// IntParams[0]：Buff 模板 ID（uint32）。
func handleSkillEffectAddBuff(w *ecs.World, caster, target ecs.Entity, eff *config.SkillEffectConfig) {
	if len(eff.IntParams) < 1 || eff.IntParams[0] <= 0 {
		return
	}
	buffID := uint32(eff.IntParams[0])
	_ = buff.AddBuff(w, caster, target, buffID)
}
