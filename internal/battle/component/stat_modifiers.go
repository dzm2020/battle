package component

// StatModifiers 本帧由 [system.BuffSystem] 从 [BuffList] 全量重算，仅作“相对基础 [Attributes] 的
// 临时修正量”，不修改 [Attributes] 持久值。 [system.DamageSystem] 在计算有效物甲/魔抗时
// 会加上 ArmorDelta、MRDelta；物理强度增量可给技能公式或后续系统使用。
type StatModifiers struct {
	PhysicalPowerDelta int // 相对 [Attributes].PhysicalPower 的 Buff 增量
	ArmorDelta         int // 与基础物甲相加得有效护甲（参与物伤减免）
	MRDelta            int // 与基础魔抗相加得有效魔抗（参与法伤减免）
}

func (*StatModifiers) Component() {}
