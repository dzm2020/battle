package skill

import "battle/internal/battle/component"

// ResourceType 技能消耗的资源种类；与 [component.SkillUser] 中的 Mana/Rage/Energy 字段对应。
type ResourceType uint8

const (
	ResourceNone ResourceType = iota // 无消耗（此时 Cost 必须为 0）
	ResourceMana                     // 法力
	ResourceRage                     // 怒气
	ResourceEnergy                   // 能量
)

// EffectKind 单条技能效果条目类型；一条技能可有多个 [EffectConfig]，按数组顺序依次执行。
type EffectKind uint8

const (
	EffectDamage EffectKind = iota // 写入 [component.PendingDamage]，经伤害系统减免后扣血
	EffectHeal                     // 直接增加 [component.Health].Current（不超过 Max）
	EffectApplyBuff                // 调用 [buff.ApplyBuff] 挂载 Buff
)

// EffectConfig 在 [ResolveTargets] 得到的目标列表上，对每个目标执行一条原子效果。
type EffectConfig struct {
	Kind       EffectKind           `json:"kind"`                 // 效果类型，决定下列哪些字段生效
	Amount     int                  `json:"amount,omitempty"`     // 伤害/治疗量；其它 Kind 可忽略
	DamageType component.DamageType `json:"damageType,omitempty"` // 仅 EffectDamage
	BuffDefID  uint32               `json:"buffDefId,omitempty"`  // 仅 EffectApplyBuff，对应 Buff 模板 ID
}

// SkillConfig 技能的静态模板（JSON/YAML）。
//
// 目标选取必须由作用范围 × 阵营 × 选取规则组合描述（见仓库根目录 skill_record.md）：
// [TargetScope]、[CampRelation]、[PickRule]，详见 skill_target_spec.go。
type SkillConfig struct {
	ID uint32 `json:"id"` // 全局唯一；与 CastIntent.SkillID、SkillUser.GrantedSkillIDs 一致

	Resource ResourceType `json:"resource"` // 消耗哪种资源
	Cost     int          `json:"cost"`     // 单次施放扣除量；与 ResourceNone 组合时须为 0

	CooldownFrames int `json:"cooldownFrames"` // 冷却帧数；从「效果结算完毕」当帧写入 SkillUser.CooldownRemaining

	// --- 目标三维度（必填 scope、camp；pickRule 可选）---
	Scope     TargetScope   `json:"scope"`               // 作用范围：单体/群体/链式等；JSON 值为 0 表示非法配置
	Camp      CampRelation  `json:"camp"`                // 阵营：敌方/友方/全体等；敌方为 JSON 数值 0
	PickRule  PickRule      `json:"pickRule,omitempty"`  // 选取规则：最近、血量排序等；0 表示不排序
	CampSide  uint8         `json:"campSide,omitempty"`  // 仅 CampSpecificSide：指定 Team.Side
	AOERadius float64       `json:"aoeRadius,omitempty"` // 球半径；与 Cone/Circle/Multi 等组合时筛选距离；0 表示不做距离裁剪

	MaxTargets       int    `json:"maxTargets,omitempty"`       // 排序后至多保留多少个目标；随机模式缺省时内部另有默认次数
	ChainJumps       int    `json:"chainJumps,omitempty"`       // 链式：首目标之外的额外跳跃次数；≤0 时默认额外 2 跳
	RequireBuffDefID uint32 `json:"requireBuffDefId,omitempty"` // 候选目标必须携带该 Buff 模板 ID
	ForbidBuffDefID  uint32 `json:"forbidBuffDefId,omitempty"`  // 候选目标不得携带该 Buff 模板 ID

	CastFrames int `json:"castFrames"` // 吟唱帧数；0 瞬发且当帧结算；>0 则在发起当帧扣费并进入 SkillCastState

	Effects []EffectConfig `json:"effects"` // 命中目标集后依次执行的效果链
}
