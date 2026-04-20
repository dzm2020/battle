package skill

// 本文件定义与仓库根目录 skill_record.md 一致的三维目标模型：作用范围 × 阵营关系 × 选取规则。
// JSON 配置字段名为 scope、camp、pickRule（见 [SkillConfig]）。

// TargetScope 表示技能在空间上的作用方式（单体、群体几何近似、链式等）。
// JSON 数值从 1 起；**0 保留为非法**，表示未正确配置 scope。
type TargetScope uint8

const (
	_ TargetScope = iota // 0：非法；禁止在 JSON 中省略 scope 或使用 0

	TargetScopeSelf      // 1：仅施法者自身，忽略 CastIntent.Target 语义上的「指向」
	TargetScopeSingle    // 2：单体；主目标为 CastIntent.Target，合法性由 Camp 约束
	TargetScopeCone      // 3：扇形（当前与圆共用：先按 Camp 取候选，再以 aoeRadius + Transform2D 裁剪）
	TargetScopeCircle    // 4：圆形；锚点优先为主目标实体坐标，无主目标则用施法者
	TargetScopeLineRect  // 5：直线/矩形（当前实现同 Circle 的球筛选）
	TargetScopeMulti     // 6：场上符合 Camp 的全员多目标（文档中「全体敌方/友方」的基础形态）
	TargetScopeFullScreen // 7：除施法者外全场带 Team+Health 的单位，再按 Camp 缩小
	TargetScopeChain     // 8：链式；首跳为 CastIntent.Target，后续从同 Camp 候选中扩展
	TargetScopeRandom    // 9：从候选中随机取若干（个数见 MaxTargets，默认 3）
)

// CampRelation 表示候选实体与施法者的阵营关系；与 TargetScopeSelf 组合时逻辑上可忽略。
type CampRelation uint8

const (
	CampEnemy CampRelation = iota // 0：敌对阵营（与施法者 Team.Side 不同且含 Health）
	CampAllyIncludeSelf           // 1：友方含施法者自身
	CampAllyExcludeSelf           // 2：友方不含施法者
	CampEveryone                  // 3：不限阵营（常配合 FullScreen）
	CampSpecificSide              // 4：仅指定 Side，配合 [SkillConfig.CampSide]
)

// PickRule 在 Buff 过滤之后、MaxTargets 截断之前，对候选列表排序（或保持无序）。
// 最近/最远依赖施法者与目标上的 [component.Transform2D]。
type PickRule uint8

const (
	PickNone PickRule = iota // 0：不排序，顺序由 ECS 遍历与收集顺序决定
	PickNearest              // 1：距施法者距离升序（缺坐标的目标视为无穷远）
	PickFarthest             // 2：距施法者距离降序
	PickHPCurrentAsc         // 3：当前生命值升序（斩杀/补刀倾向）
	PickHPPercentAsc         // 4：当前生命/最大生命升序（群体治疗常用）
	PickAttackHighest        // 5：Attributes.PhysicalPower 降序；无 Attributes 视为 0
)

// targetSpec 将 [SkillConfig] 展开为 [ResolveTargets] 使用的运行期快照（无业务逻辑，仅字段拷贝）。
type targetSpec struct {
	Scope            TargetScope
	Camp             CampRelation
	CampSide         uint8
	Pick             PickRule
	MaxTargets       int
	ChainJumps       int
	RequireBuffDefID uint32
	ForbidBuffDefID  uint32
	AOERadius        float64
}

// effectiveTargetSpec 从配置构造 targetSpec；若日后增加默认值或校验可集中在此处。
func effectiveTargetSpec(sk SkillConfig) targetSpec {
	return targetSpec{
		Scope:            sk.Scope,
		Camp:             sk.Camp,
		Pick:             sk.PickRule,
		MaxTargets:       sk.MaxTargets,
		ChainJumps:       sk.ChainJumps,
		RequireBuffDefID: sk.RequireBuffDefID,
		ForbidBuffDefID:  sk.ForbidBuffDefID,
		AOERadius:        sk.AOERadius,
		CampSide:         sk.CampSide,
	}
}
