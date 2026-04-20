package buff

import (
	"battle/internal/battle/component"
	"battle/internal/battle/control"
)

// StackPolicy 同名 Buff（同一 [DescriptorConfig].ID）再次施加时的叠层策略。
type StackPolicy uint8

const (
	// StackIndependent 始终追加一条新的 [component.BuffInstance]，各自 FramesLeft 独立递减。
	StackIndependent StackPolicy = iota
	// StackRefresh 若已存在同 DefID 实例：仅将 FramesLeft 重置为 DurationFrames，不改变 Stacks；
	// 若不存在则新建一条 1 层实例（实现见 [ApplyBuff]）。
	StackRefresh
	// StackMerge 若已存在同 DefID：Stacks++（封顶 MaxStacks）并重置 FramesLeft；否则新建。
	StackMerge
)

// EffectKind 描述一条 [EffectConfig] 的语义；同一 [DescriptorConfig] 可组合多种 Kind。
type EffectKind uint8

const (
	EffectStatMod EffectKind = iota // 属性增量，写入 [component.StatModifiers]
	EffectDoT                       // 周期性伤害，经 [component.MergePendingDamage] 写入结算链
	EffectHoT                       // 周期性治疗，直接修改 [component.Health].Current
	EffectControl                   // 控制位，OR 入 [component.ControlState].Flags
)

// EffectConfig 静态配置的一条子效果；同属一条 [DescriptorConfig] 的多条 EffectConfig 在同一 Buff 实例上并存。
// 数值类效果在 [BuffSystem] 中与实例 Stacks 相乘；控制类不按层翻倍，仅为按位或到 [component.ControlState]。
type EffectConfig struct {
	// Kind 决定本条如何参与结算；必须与下方非零字段语义一致（见各字段说明）。
	Kind EffectKind `json:"kind"`

	// ArmorDeltaPerStack 每层对物甲的增量；仅 EffectStatMod 使用，汇总至 [component.StatModifiers].ArmorDelta。
	ArmorDeltaPerStack int `json:"armorDeltaPerStack,omitempty"`
	// MRDeltaPerStack 每层对魔抗的增量；仅 EffectStatMod，汇总至 StatModifiers.MRDelta。
	MRDeltaPerStack int `json:"mrDeltaPerStack,omitempty"`
	// PowerDeltaPerStack 每层对 [component.Attributes].PhysicalPower 的等价增量；仅 EffectStatMod。
	PowerDeltaPerStack int `json:"powerDeltaPerStack,omitempty"`

	// DamagePerTick 单次 DoT 一跳的基础伤害（再乘 Stacks）；仅 EffectDoT。
	DamagePerTick int `json:"damagePerTick,omitempty"`
	// DamageType DoT 的伤害类型（物/魔/真）；仅 EffectDoT，写入 [component.MergePendingDamage]。
	DamageType component.DamageType `json:"damageType,omitempty"`
	// TickIntervalFrames DoT/HoT 相邻两跳之间的帧间隔，>=1（实现会钳制为至少 1）；
	// DoT/HoT 共用 [component.BuffInstance].TickCountdown；同一 DescriptorConfig 取首个 DoT/HoT 出现的间隔为准。
	TickIntervalFrames int `json:"tickIntervalFrames,omitempty"`

	// HealPerTick 单次 HoT 一跳的基础治疗量（再乘 Stacks）；仅 EffectHoT，直接加 [component.Health].Current。
	HealPerTick int `json:"healPerTick,omitempty"`

	// Control 控制位掩码；仅 EffectControl，与实例上其它控制效果按位或写入 ControlState.Flags。
	Control control.Flags `json:"control,omitempty"`
}

// DescriptorConfig 单种 Buff 的模板：与运行时 [component.BuffInstance] 通过 ID 关联，可多效果组合。
type DescriptorConfig struct {
	ID             uint32      `json:"id"`             // 全表唯一，与 BuffInstance.DefID 一致
	MaxStacks      int         `json:"maxStacks"`      // Merge/Refresh 时 Stacks 上限，至少按 1 处理
	Policy         StackPolicy `json:"policy"`         // 再次施加同名 Buff 时的策略
	DurationFrames int         `json:"durationFrames"` // >=0 时走时间到期；=-1 为不限天然时长（FramesLeft 为负）
	Effects        []EffectConfig `json:"effects"`     // 本 Buff 实例激活时生效的全部子效果
}
