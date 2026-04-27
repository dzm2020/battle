package buff

import (
	"strings"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/control"
	"battle/internal/battle/log"
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
		log.Debug("[buff] 未注册的效果类型 实体=%v Buff编号=%d 效果类型=%d", e, buff.BuffId, desc.EffectType)
	}
}

// handleBufferEffectStatChange 属性修改效果
func handleBufferEffectStatChange(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	mods := ecs.EnsureGetComponent[*component.StatModifiers](world, e)
	if len(desc.ParamsString) < 1 || len(desc.Params) < 1 {
		log.Debug("[buff] 属性变更效果：参数不足 实体=%v Buff编号=%d", e, buff.BuffId)
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
		log.Debug("[buff] 属性变更效果：未识别的属性键 实体=%v Buff编号=%d 键=%s", e, buff.BuffId, key)
	}
}

// handleBufferEffectControl 控制效果
func handleBufferEffectControl(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	if desc == nil || len(desc.ParamsString) < 1 {
		log.Debug("[buff] 控制效果：参数不足 实体=%v Buff编号=%d", e, buff.BuffId)
		return
	}
	ctrl := ecs.EnsureGetComponent[*component.ControlState](world, e)
	tag := strings.ToLower(strings.TrimSpace(desc.ParamsString[0]))
	switch tag {
	case "stun", "stunned":
		ctrl.Flags |= control.FlagStunned
		log.Debug("[buff] 控制效果：眩晕 实体=%v Buff编号=%d", e, buff.BuffId)
	case "silence", "silenced":
		ctrl.Flags |= control.FlagSilenced
		log.Debug("[buff] 控制效果：沉默 实体=%v Buff编号=%d", e, buff.BuffId)
	case "root", "rooted":
		ctrl.Flags |= control.FlagRooted
		log.Debug("[buff] 控制效果：禁锢 实体=%v Buff编号=%d", e, buff.BuffId)
	default:
		log.Debug("[buff] 控制效果：未识别的标签 实体=%v Buff编号=%d 标签=%s", e, buff.BuffId, tag)
	}
}

// handleBufferEffectDamage 伤害效果（周期触发时的结算也可走同一 dict；列表合并见 [tickApplyBuffEffect]）。
func handleBufferEffectDamage(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	if desc == nil || len(desc.Params) < 1 {
		log.Debug("[buff] 伤害效果：参数不足 实体=%v Buff编号=%d", e, buff.BuffId)
		return
	}
	st := max(buff.Stacks, 1)
	amt := int(desc.Params[0]) * st
	if amt <= 0 {
		log.Debug("[buff] 伤害效果：跳过非正数额 实体=%v Buff编号=%d 数额=%d", e, buff.BuffId, amt)
		return
	}
	com := ecs.EnsureGetComponent[*component.PendingDamageBuff](world, e)
	com.Amount += amt
	log.Debug("[buff] 持续伤害：累计写入 实体=%v Buff编号=%d 本帧增量=%d 当前合计=%d", e, buff.BuffId, amt, com.Amount)
}

// handleBufferEffectHeal 治疗效果
func handleBufferEffectHeal(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	if desc == nil || len(desc.Params) < 1 {
		log.Debug("[buff] 治疗效果：参数不足 实体=%v Buff编号=%d", e, buff.BuffId)
		return
	}
	st := max(buff.Stacks, 1)
	heal := int(desc.Params[0]) * st
	if heal <= 0 {
		log.Debug("[buff] 治疗效果：跳过非正数额 实体=%v Buff编号=%d 数额=%d", e, buff.BuffId, heal)
		return
	}
	com := ecs.EnsureGetComponent[*component.PendingHealBuff](world, e)
	com.Amount += heal
	log.Debug("[buff] 持续治疗：累计写入 实体=%v Buff编号=%d 本帧增量=%d 当前合计=%d", e, buff.BuffId, heal, com.Amount)
}
