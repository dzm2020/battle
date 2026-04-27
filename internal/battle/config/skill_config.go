package config

// ======================= 枚举定义 =======================

// SkillType 技能类型
type SkillType int

const (
	SkillTypeActive  SkillType = 0 // 主动技能
	SkillTypePassive SkillType = 1 // 被动技能（预留）
)

// AttributeType 消耗属性类型
type AttributeType int

const (
	AttrNone   AttributeType = 0 // 无消耗
	AttrHP     AttributeType = 1 // 生命值
	AttrMP     AttributeType = 2 // 法力值
	AttrEnergy AttributeType = 3 // 能量
	AttrRage   AttributeType = 4 // 怒气
)

// EffectType 技能效果类型
type EffectType int

const (
	EffectDamage  EffectType = 1 // 造成伤害
	EffectAddBuff EffectType = 2 // 添加Buff
	EffectSummon  EffectType = 3 // 召唤单位
	EffectBlink   EffectType = 4 // 闪现/瞬移
	EffectResetCD EffectType = 5 // 重置冷却
)

// CompareOp 比较运算符
type CompareOp int

const (
	OpEqual        CompareOp = 1 // ==
	OpNotEqual     CompareOp = 2 // !=
	OpGreater      CompareOp = 3 // >
	OpLess         CompareOp = 4 // <
	OpGreaterEqual CompareOp = 5 // >=
	OpLessEqual    CompareOp = 6 // <=
)

// ConditionType 被动条件类型（用于被动条件配置表）
type ConditionType int

const (
	CondOnSkillCast   ConditionType = 1 // 当释放某个技能时
	CondOnTakeDamage  ConditionType = 2 // 当受到伤害时
	CondOnDealDamage  ConditionType = 3 // 当造成伤害时
	CondOnBuffApplied ConditionType = 4 // 当获得Buff时
	CondOnHealthBelow ConditionType = 5 // 当生命值低于某比例时
)

// ======================= 配置结构体 =======================

// SkillBaseConfig 技能基础配置
type SkillBaseConfig struct {
	ID              int           `json:"id"`               // 技能ID
	SkillType       SkillType     `json:"skill_type"`       // 技能类型（主动/被动）
	ConsumeType     AttributeType `json:"consume_type"`     // 消耗属性类型
	ConsumeValue    int           `json:"consume_value"`    // 消耗值
	PreCastFrames   int           `json:"pre_cast_frames"`  // 前摇帧数
	AfterCastFrames int           `json:"post_cast_frames"` // 后摇帧数
	CooldownFrames  int           `json:"cooldown_frames"`  // 冷却帧数
	EffectIDs       []int         `json:"effect_ids"`       // 技能效果ID列表（按顺序执行）
}

// SkillEffectConfig 技能效果配置
type SkillEffectConfig struct {
	EffectID       int        `json:"effect_id"`        // 效果ID
	EffectType     EffectType `json:"effect_type"`      // 效果类型
	IntParams      []int      `json:"int_params"`       // 整数参数列表
	StringParams   []string   `json:"string_params"`    // 字符串参数列表
	TargetSelectID int        `json:"target_select_id"` // 选取目标配置ID（0表示不需要目标）
}

// ======================= 被动技能相关 =======================

// PassiveConditionConfig 被动条件配置表
// 描述被动技能的触发条件和触发后的效果（通常触发后执行一组效果ID）
type PassiveConditionConfig struct {
	ID               int           `json:"id"`                 // 条件配置ID
	ConditionType    ConditionType `json:"condition_type"`     // 条件类型
	IntParams        []int         `json:"int_params"`         // 条件参数（整数）
	StringParams     []string      `json:"string_params"`      // 条件参数（字符串）
	TriggerEffectIDs []int         `json:"trigger_effect_ids"` // 触发时执行的效果ID列表
}
