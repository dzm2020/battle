package config

// CampRelation 表示候选实体与施法者的阵营关系；与 TargetScopeSelf 组合时逻辑上可忽略。
type CampRelation uint8

const (
	CampEnemy           CampRelation = iota // 0：敌对阵营（与施法者 Team.Side 不同且含 Health）
	CampAllyIncludeSelf                     // 1：友方含施法者自身
	CampAllyExcludeSelf                     // 2：友方不含施法者
	CampEveryone                            // 3：不限阵营（常配合 FullScreen）
	CampSpecificSide                        // 4：仅指定 Side，配合 [SkillConfig.CampSide]
)

// TargetScope 表示技能在空间上的作用方式（单体、群体几何近似、链式等）。
// JSON 数值从 1 起；**0 保留为非法**，表示未正确配置 scope。
type TargetScope uint8

const (
	_ TargetScope = iota // 0：非法；禁止在 JSON 中省略 scope 或使用 0

	TargetScopeSelf       // 1：仅施法者自身，忽略 CastIntent.Target 语义上的「指向」
	TargetScopeSingle     // 2：单体；主目标为 CastIntent.Target，合法性由 Camp 约束
	TargetScopeCone       // 3：扇形（当前与圆共用：先按 Camp 取候选，再以 aoeRadius + Transform2D 裁剪）
	TargetScopeCircle     // 4：圆形；锚点优先为主目标实体坐标，无主目标则用施法者
	TargetScopeLineRect   // 5：直线/矩形（当前实现同 Circle 的球筛选）
	TargetScopeMulti      // 6：场上符合 Camp 的全员多目标（文档中「全体敌方/友方」的基础形态）
	TargetScopeFullScreen // 7：除施法者外全场带 Team+Health 的单位，再按 Camp 缩小
	TargetScopeChain      // 8：链式；首跳为 CastIntent.Target，后续从同 Camp 候选中扩展
	TargetScopeRandom     // 9：从候选中随机取若干（个数见 MaxTargets，默认 3）
)

// PickRule 在 Buff 过滤之后、MaxTargets 截断之前，对候选列表排序（或保持无序）。
// 最近/最远依赖施法者与目标上的 [component.Transform2D]。
type PickRule uint8

const (
	PickNone          PickRule = iota // 0：不排序，顺序由 ECS 遍历与收集顺序决定
	PickNearest                       // 1：距施法者距离升序（缺坐标的目标视为无穷远）
	PickFarthest                      // 2：距施法者距离降序
	PickHPCurrentAsc                  // 3：当前生命值升序（斩杀/补刀倾向）
	PickHPPercentAsc                  // 4：当前生命/最大生命升序（群体治疗常用）
	PickAttackHighest                 // 5：Attributes.Values["attack_damage"] 降序；无 Attributes 视为 0
)

type EffectType string

const (
	// 伤害类
	EffectDamageDirect   EffectType = "damage_direct"    // 直接伤害
	EffectDamageOverTime EffectType = "damage_over_time" // 持续伤害

	// 治疗与护盾
	EffectHealDirect   EffectType = "heal_direct"    // 直接治疗
	EffectHealOverTime EffectType = "heal_over_time" // 持续治疗
	EffectShield       EffectType = "shield"         // 添加护盾

	// 属性修改（Buff/Debuff）
	EffectModifyStat EffectType = "modify_stat" // 修改一个或多个属性（攻击力、护甲等）
	EffectModifyMana EffectType = "modify_mana" // 增加/减少法力值

	// 控制效果
	EffectStun      EffectType = "stun"      // 眩晕
	EffectRoot      EffectType = "root"      // 定身（不能移动，可攻击）
	EffectSilence   EffectType = "silence"   // 沉默（不能施法）
	EffectTaunt     EffectType = "taunt"     // 嘲讽
	EffectKnockback EffectType = "knockback" // 击退
	EffectFear      EffectType = "fear"      // 恐惧

	// 特殊机制
	EffectSummon   EffectType = "summon"   // 召唤单位
	EffectTeleport EffectType = "teleport" // 传送
	EffectRevive   EffectType = "revive"   // 复活
	EffectCleanse  EffectType = "cleanse"  // 净化所有负面效果

	// 添加Buff
	EffectApplyBuff EffectType = "apply_buff" // 添加Buff
)

type EffectConfig struct {
	Type   EffectType             `json:"type" yaml:"type"`
	Params map[string]interface{} `json:"params" yaml:"params"` // key:参数名 value:参数值
}

// SkillConfig
// @Description: 技能配置
type SkillConfig struct {
	ID uint32 `json:"id"` // 全局唯一
	// 消耗
	Resource Attribute `json:"resource"` // 消耗哪种属性
	Cost     int       `json:"cost"`     // 单次施放扣除量；与 ResourceNone 组合时须为 0

	// --- 目标选取三维度（必填 scope、camp；pickRule 可选）---
	Camp             CampRelation `json:"camp"`                       // 阵营：敌方/友方/全体等；敌方为 JSON 数值 0
	Scope            TargetScope  `json:"scope"`                      // 作用范围：单体/群体/链式等；JSON 值为 0 表示非法配置
	PickRule         PickRule     `json:"pickRule,omitempty"`         // 选取规则：最近、血量排序等；0 表示不排序
	CampSide         uint8        `json:"campSide,omitempty"`         // 仅 CampSpecificSide：指定 Team.Side
	AOERadius        float64      `json:"aoeRadius,omitempty"`        // 球半径；与 Cone/Circle/Multi 等组合时筛选距离；0 表示不做距离裁剪
	RequireBuffDefID uint32       `json:"requireBuffDefId,omitempty"` // 候选目标必须携带该 Buff 模板 ID
	ForbidBuffDefID  uint32       `json:"forbidBuffDefId,omitempty"`  // 候选目标不得携带该 Buff 模板 ID

	MaxTargets int `json:"maxTargets,omitempty"` // 选取目标数 排序后至多保留多少个目标；随机模式缺省时内部另有默认次数

	// 前摇  后摇
	PreCastDelay  float64 `json:"pre_cast_delay" yaml:"pre_cast_delay"`   // 前摇时间
	PostCastDelay float64 `json:"post_cast_delay" yaml:"post_cast_delay"` // 后摇时间
	Interruptible bool    `json:"interruptible" yaml:"interruptible"`     // 可选：是否可被打断（前摇期间）
	// 效果
	Effects []EffectConfig `json:"effects"` // 命中目标集后依次执行的效果链
}
