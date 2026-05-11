package buff_effect

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
)

// handleBufferEffectStatChange 属性修改效果
func handleBufferEffectStatChange(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	mods := ecs.EnsureGetComponent[*component.StatModifiers](world, e)
	if len(desc.ParamsString) < 1 || len(desc.Params) < 1 {
		log.Error("[buff] 属性变更效果：参数不足 实体=%v Buff编号=%d", e, buff.BuffId)
		return
	}
	key := desc.ParamsString[0]
	delta := int(desc.Params[0]) * buff.Stacks
	switch key {
	case config.AttrArmor:
		mods.ArmorDelta += delta
	case config.AttrMagicResist:
		mods.MRDelta += delta
	case config.AttrAttackDamage:
		mods.AttackDamageDelta += delta
	case config.AttrHitPermille:
		mods.HitDeltaPermille += delta
	case config.AttrDodgePermille:
		mods.DodgeDeltaPermille += delta
	case config.AttrCritRate:
		mods.CritRateDeltaPermille += delta
	case config.AttrCritDamage:
		mods.CritDamageDeltaPermille += delta
	default:
		log.Error("[buff] 属性变更效果：未识别的属性键 实体=%v Buff编号=%d 键=%s", e, buff.BuffId, key)
	}
	log.Info("[buff] 属性变更效果 实体=%v Buff编号=%d 键=%s mods=%+v", e, buff.BuffId, key, mods)
}
