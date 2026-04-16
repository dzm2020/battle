package skill

import "battle/internal/battle/entity"

// ApplyContext 技能「生效」瞬间的上下文：第 6 天在此接入伤害、Buff 挂载等。
type ApplyContext struct {
	Frame   uint64
	Caster  *entity.Entity
	Target  *entity.Entity
	Config  SkillConfig
}

// EffectApplier 效果层抽象：System 在完成全部校验后调用，避免把伤害公式写进 System。
type EffectApplier interface {
	// OnSkillApplied 实现扣蓝之外的战斗效果；默认实现可为空操作 + 日志。
	OnSkillApplied(ctx ApplyContext)
}

// DefaultApplier 第 5～6 天之间的占位：不写伤害，仅预留钩子。
type DefaultApplier struct{}

func (DefaultApplier) OnSkillApplied(ctx ApplyContext) {
	_ = ctx
}
