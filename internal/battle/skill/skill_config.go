package skill

import "battle/internal/battle/component"

// ResourceType 技能消耗的资源种类；与 [component.SkillUser] 中的 Mana/Rage/Energy 字段对应。
type ResourceType uint8

const (
	// ResourceNone 无资源消耗，仍受冷却约束。
	ResourceNone ResourceType = iota
	ResourceMana
	ResourceRage
	ResourceEnergy
)

// TargetKind 技能主目标选取方式；群体类仍可能依赖 [component.CastIntent.Target] 作为指向目标（可选）。
type TargetKind uint8

const (
	// TargetSelf 仅作用于施法者自身（忽略 CastIntent.Target）。
	TargetSelf TargetKind = iota
	// TargetSingleEnemy 单体敌方：需通过 [component.CastIntent.Target] 指定实体，且阵营与施法者不同（见 [component.Team]）。
	TargetSingleEnemy
	// TargetAllEnemySides 场上除施法者外，所有具有 [component.Team] 且 Side 与施法者不同的、
	// 同时含 [component.Health] 的实体（用于范围伤害/群体挂 Buff）。
	TargetAllEnemySides
)

// EffectKind 单条技能效果条目类型，一条 [SkillConfig] 可由多个 [EffectConfig] 顺序执行。
type EffectKind uint8

const (
	// EffectDamage 对选取目标列表分别 [component.MergePendingDamage]。
	EffectDamage EffectKind = iota
	// EffectHeal 对选取目标列表分别增加 [component.Health].Current（不超过 Max）。
	EffectHeal
	// EffectApplyBuff 对选取目标列表分别 [buff.ApplyBuff]；DefID 为 Buff 模板 ID。
	EffectApplyBuff
)

// EffectConfig 技能在“结算目标集”上执行的一条原子效果；与 [SkillConfig].Effects 顺序一致执行。
type EffectConfig struct {
	// Kind 决定本条目使用哪些可选字段。
	Kind EffectKind `json:"kind"`

	// Amount 对 EffectDamage 为伤害量（再经伤害系统减免）；对 EffectHeal 为治疗量；其它 Kind 忽略。
	Amount int `json:"amount,omitempty"`
	// DamageType 仅 EffectDamage 使用。
	DamageType component.DamageType `json:"damageType,omitempty"`
	// BuffDefID 仅 EffectApplyBuff，对应 [buff.DefinitionConfig] 中的 DescriptorConfig.ID。
	BuffDefID uint32 `json:"buffDefId,omitempty"`
}

// SkillConfig 技能的静态模板（可由 JSON/YAML 加载）；运行时施放状态见 [component.SkillUser]、[component.SkillCastState]。
type SkillConfig struct {
	// ID 全表唯一，对应施放请求中的技能 ID。
	ID uint32 `json:"id"`

	// Resource 消耗种类；ResourceNone 时 Cost 必须为 0，否则 [SkillIntentSystem] 拒绝施放。
	Resource ResourceType `json:"resource"`
	// Cost 单次施放扣除的资源数值（非法或不足时施放失败）。
	Cost int `json:"cost"`

	// CooldownFrames 完整冷却帧数；从“效果结算完毕”当帧起算（瞬发结算后立刻进入冷却；吟唱在吟唱结束结算后）。
	CooldownFrames int `json:"cooldownFrames"`

	// Target 目标选取策略；决定 [ResolveTargets] 如何从世界收集实体列表。
	Target TargetKind `json:"target"`

	// CastFrames 吟唱/引导帧数；0 表示瞬发。大于 0 时在发起当帧扣除资源并写入 [component.SkillCastState]，
	// 经过 CastFrames 次 [SkillChannelSystem] 更新后结算效果并进入冷却（资源不在结算帧二次扣除）。
	CastFrames int `json:"castFrames"`

	// Effects 命中目标后顺序执行的效果列表（伤害、治疗、挂 Buff 可混排）。
	Effects []EffectConfig `json:"effects"`
}
