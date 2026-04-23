package buff

import (
	"strings"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/control"
)

func init() {
	registerEffectHandler(config.BufferEffectStatChange, handleBufferEffectStatChange)
	registerEffectHandler(config.BufferEffectControl, handleBufferEffectControl)
	registerEffectHandler(config.BufferEffectDamage, handleBufferEffectDamage)
	registerEffectHandler(config.BufferEffectHeal, handleBufferEffectHeal)
}

// effectHandlerFn 按模板中某一种 [BufferEffectType] 聚合处理：DoT/HoT 仍在 [TickPeriodicBuffEffects]。
type effectHandlerFn func(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig)

var effectHandlerDict = make(map[config.BufferEffectType]effectHandlerFn)

func registerEffectHandler(typ config.BufferEffectType, fn effectHandlerFn) {
	effectHandlerDict[typ] = fn
}

func applyEffect(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	if fn := effectHandlerDict[desc.EffectType]; fn != nil {
		fn(world, e, buff, desc)
	} else {

	}
}

// handleBufferEffectStatChange 属性修改效果
func handleBufferEffectStatChange(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	mods := ecs.EnsureGetComponent[*component.StatModifiers](world, e)
	if len(desc.ParamsString) < 1 || len(desc.Params) < 1 {
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
	}
}

// handleBufferEffectControl 控制效果
func handleBufferEffectControl(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	if desc == nil || len(desc.ParamsString) < 1 {
		return
	}
	ctrl := ecs.EnsureGetComponent[*component.ControlState](world, e)
	switch strings.ToLower(strings.TrimSpace(desc.ParamsString[0])) {
	case "stun", "stunned":
		ctrl.Flags |= control.FlagStunned
	case "silence", "silenced":
		ctrl.Flags |= control.FlagSilenced
	case "root", "rooted":
		ctrl.Flags |= control.FlagRooted
	}
}

// handleBufferEffectDamage 伤害效果（周期触发时的结算也可走同一 dict；具体合并见 [Manager.applyPeriodicEffects]）。
func handleBufferEffectDamage(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	if desc == nil || len(desc.Params) < 1 {
		return
	}
	st := max(buff.Stacks, 1)
	amt := int(desc.Params[0]) * st
	if amt <= 0 {
		return
	}
	com := ecs.EnsureGetComponent[*component.PendingDamageBuff](world, e)
	com.Amount += amt
}

// handleBufferEffectHeal 治疗效果
func handleBufferEffectHeal(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	if desc == nil || len(desc.Params) < 1 {
		return
	}
	st := max(buff.Stacks, 1)
	heal := int(desc.Params[0]) * st
	if heal <= 0 {
		return
	}
	com := ecs.EnsureGetComponent[*component.PendingHealBuff](world, e)
	com.Amount += heal
}
