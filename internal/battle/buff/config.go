package buff

import "battle/internal/battle/control"

// Kind Buff 语义类型（同一 Kind 可有不同配表实例）。
type Kind uint8

const (
	KindNone Kind = iota
	// KindInstantHeal 瞬时治疗，不入持续列表。
	KindInstantHeal
	// KindInstantDamage 瞬时伤害（演示用）。
	KindInstantDamage
	// KindDot 周期伤害，需 DurationFrames + TickIntervalFrames。
	KindDot
	// KindStatATK 持续期间增加攻击力（按层数叠加平铺值）。
	KindStatATK
	// KindStun 持续眩晕。
	KindStun
	// KindSlow 持续减速（乘算移速系数）。
	KindSlow
	// KindDamageAmp 持续提高造成的伤害（乘在 outgoing 上）。
	KindDamageAmp
)

// StackPolicy 同 ID 再次添加时的规则。
type StackPolicy uint8

const (
	// StackRefresh 刷新持续时间；层数不变（瞬时类忽略）。
	StackRefresh StackPolicy = iota
	// StackLayer 层数 +1，不超过 MaxStacks；持续时间取 max(旧剩余, 新持续)。
	StackLayer
	// StackReplace 总是移除旧实例并添加新实例（同 ID 单实例）。
	StackReplace
)

// BuffConfig 配表项：由 Registry 持有，运行时只读。
type BuffConfig struct {
	ID   string
	Name string
	Kind Kind

	DurationFrames     uint64 // 持续帧数；瞬时类可为 0
	TickIntervalFrames uint64 // DoT/HoT 间隔；非周期类为 0
	TickDeltaHP        int64  // 每跳修改 CurHP（负数为伤害）

	// StatATKFlat 每层为宿主增加的攻击力（KindStatATK）。
	StatATKFlat int64
	MaxStacks   int32

	SlowMoveMul float64 // KindSlow：移速乘子，如 0.6 表示 60% 移速

	// OutDamageMul 每层对「造成伤害」的乘子，最终 outgoing *= pow(OutDamageMul, stacks)（KindDamageAmp）。
	OutDamageMul float64

	Control control.Flags // KindStun 等写控制位

	StackPolicy StackPolicy
	Priority    int32 // 预留：同组覆盖时比较；当前 StackReplace 外未启用
}
