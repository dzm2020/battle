package component

import "battle/ecs"

// BuffList 单实体上的 Buff 容器：所有可叠加/并存的实例都缓存在 Buffs 切片中（即运行期
// 的“状态表/缓冲表”），与 [config.BuffConfig] 表内键（即 BuffId）一一对应。
// 无 Buff 时本组件应被移除，并同时清掉 [StatModifiers]、[ControlState] 等派生数据。
type BuffList struct {
	// Buffs 按加入顺序或系统维护顺序排列；同 BuffId 可因叠层策略出现多条（Independent）或单条（Refresh/Merge）。
	Buffs []BuffInstance
}

func (*BuffList) Component() {}

// BuffInstance 单条 Buff 的运行态（“槽位”数据），不内联全部效果，只存层数、剩余时间
// 与 DoT/HoT 节拍；效果定义在 [config.BuffConfig] 中，由 [config.Tab].BuffConfigConfigByID[BuffId] 查表。
type BuffInstance struct {
	// BuffId 即 [config.BuffConfig] 在表中的 int32 主键；表内无此键时 [system.BuffSystem] 会丢弃该实例。
	BuffId uint32

	// Stacks 层数，参与属性/DoT/HoT 的每跳强度（控制类效果不按层数放大，只由是否含该 Def 决定）。
	Stacks int

	// FramesLeft 剩余持续帧；每帧末在 [system.BuffSystem] 中减 1。非负时到期则移除此实例。
	// 为负表示不受时间轴自然结束（如永久或需逻辑/驱散结束），本系统不自动递减或删除。
	FramesLeft int

	// TickCountdown 为 DoT/HoT 用：每帧先自减，<0 时结算一跳后按间隔重置，与 FramesLeft 独立。
	// 无 DoT/HoT 的 Def 在 [buff.apply] 中仍可能赋初值，但 [system.BuffSystem] 会早退不推进。
	TickCountdown int

	// 施法者实体（用于伤害来源）
	Caster ecs.Entity
}
