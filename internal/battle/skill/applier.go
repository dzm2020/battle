package skill

import (
	"battle/internal/battle/buff"
	"battle/internal/battle/calc"
	"battle/internal/battle/entity"
)

// ApplyContext 技能成功结算上下文。
type ApplyContext struct {
	Frame  uint64
	Caster *entity.Entity
	Target *entity.Entity
	Config SkillConfig
}

// EffectApplier 技能效果钩子：伤害、治疗、Buff。
type EffectApplier interface {
	OnSkillApplied(ctx ApplyContext)
}

// DefaultApplier 空实现，便于测试只校验不落效果。
type DefaultApplier struct{}

func (DefaultApplier) OnSkillApplied(ApplyContext) {}

// BattleApplier 第 6～7 天默认战斗效果：物理伤害、治疗、挂 Buff。
type BattleApplier struct{}

func (BattleApplier) OnSkillApplied(ctx ApplyContext) {
	if ctx.Caster == nil {
		return
	}
	eff := ctx.Target
	if eff != nil && ctx.Config.Damage > 0 && !eff.IsDead() {
		mul := ctx.Caster.OutgoingDamageMul
		dmg := calc.PhysicalHit(ctx.Caster.Derived, eff.Derived, mul)
		dmg += ctx.Config.Damage
		if dmg < 1 {
			dmg = 1
		}
		eff.Runtime.CurHP -= dmg
		if eff.Runtime.CurHP < 0 {
			eff.Runtime.CurHP = 0
		}
	}
	if eff != nil && ctx.Config.Heal > 0 && !eff.IsDead() {
		eff.Runtime.CurHP += ctx.Config.Heal
		buff.ClampHPToMax(&eff.Runtime, eff.Derived)
	}
	if eff != nil && len(ctx.Config.TargetBuffIDs) > 0 && !eff.IsDead() {
		for _, bid := range ctx.Config.TargetBuffIDs {
			eff.AddBuff(ctx.Frame, bid)
		}
	}
}
