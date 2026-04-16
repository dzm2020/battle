package attr

// Derived 战斗属性（衍生）：由 Base 经固定公式计算，供技能与伤害模块只读使用。
type Derived struct {
	MaxHP int64
	MaxMP int64
	ATK   int64
	DEF   int64

	// CritRate 暴击率，范围 [0,1]。
	CritRate float64
	// CritDamage 暴击伤害倍率，例如 1.5 表示额外按 150% 结算（与项目最终公式对齐时再调）。
	CritDamage float64
	// PhysMitigation 物理承伤减免系数 [0,1)，实际物理承伤比例约为 1 - PhysMitigation。
	PhysMitigation float64
}
