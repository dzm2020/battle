package skill_effect

import (
	"battle/ecs"
	"battle/internal/battle/target_selector"

	"battle/internal/battle/config"
)

type skillEffectFn func(w *ecs.World, caster, target ecs.Entity, eff *config.SkillEffectConfig)

var skillEffectDict = make(map[config.EffectType]skillEffectFn)

func registerSkillEffect(typ config.EffectType, fn skillEffectFn) {
	skillEffectDict[typ] = fn
}

// Apply 按 [config.SkillBaseConfig.EffectIDs] 顺序执行技能效果（伤害、加 Buff 等）。
// caster：施法者；skillID：技能配置 ID。
func Apply(w *ecs.World, caster ecs.Entity, skillID int) {
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

func applySkillEffect(w *ecs.World, caster, target ecs.Entity, eff *config.SkillEffectConfig) {
	if w == nil || eff == nil || target == 0 || !w.EntityExists(target) {
		return
	}
	if fn := skillEffectDict[eff.EffectType]; fn != nil {
		fn(w, caster, target, eff)
	}
}

func init() {
	registerSkillEffect(config.EffectDamage, handleSkillEffectDamage)
	registerSkillEffect(config.EffectAddBuff, handleSkillEffectAddBuff)
}
