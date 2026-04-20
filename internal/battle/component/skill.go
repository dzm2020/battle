package component

import "battle/ecs"

// SkillUser 战斗单位可施放技能的运行时状态：资源池、已学会技能列表与冷却表。
// 静态数值（消耗、冷却时长、效果）来自 [skill.CatalogConfig]，不在组件内重复存储。
type SkillUser struct {
	// Mana / Rage / Energy 三种资源槽；具体技能消耗哪种由 [skill.SkillConfig.Resource] 决定。
	Mana   int
	Rage   int
	Energy int

	// GrantedSkillIDs 当前允许施放的技能模板 ID 列表（须在 [skill.CatalogConfig] 中存在）。
	GrantedSkillIDs []uint32

	// CooldownRemaining 记录技能 ID → 剩余冷却帧数；<=0 的条目应在 [system.CooldownSystem] 中删除。
	// nil 表示从未进入过冷却，等同于“无条目”。
	CooldownRemaining map[uint32]int
}

func (*SkillUser) Component() {}

// CastIntent 由外部玩法层写入，表示“本实体希望施放某技能”；[system.SkillIntentSystem] 消费后应移除组件。
// 同一实体同一帧至多处理一次意图（后者覆盖前者由玩法层避免）。
type CastIntent struct {
	// SkillID 对应 [skill.SkillConfig].ID。
	SkillID uint32
	// Target 主目标；[skill.TargetSelf] 时忽略；单体敌方、部分 AOE 需要有效实体 ID。
	Target ecs.Entity
}

func (*CastIntent) Component() {}

// SkillCastState 吟唱/引导中的状态；FramesLeft 每帧递减，归零当帧结算效果并进入冷却。
// 与 CastIntent 互斥：吟唱期间不应再写入新的 CastIntent（应由玩法层阻止）。
type SkillCastState struct {
	SkillID uint32
	// PrimaryTarget 施放开始时锁定；[skill.TargetSelf] 时可为 0。
	PrimaryTarget ecs.Entity
	// FramesLeft 剩余吟唱帧；初始化为 [skill.SkillConfig].CastFrames（>0）。
	FramesLeft int
}

func (*SkillCastState) Component() {}
