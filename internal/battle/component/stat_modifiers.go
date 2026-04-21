package component

// StatModifiers 本帧由 [system.BuffSystem] 从 [BuffList] 全量重算，仅作“相对基础 [Attributes.Values] 的
// 临时修正量”，不修改 [Attributes] 持久值。 [system.DamageSystem] 在计算有效物甲/魔抗时
// 会加上 ArmorDelta、MRDelta；攻击强度增量对应键 attack_damage。
type StatModifiers struct {
	AttackDamageDelta int // 相对 [Attributes.Values]["attack_damage"] 的 Buff 增量（物理强度）
	ArmorDelta        int // 与基础物甲相加得有效护甲（参与物伤减免）
	MRDelta           int // 与基础魔抗相加得有效魔抗（参与法伤减免）

	// 千分比修正，与 [Attributes.Values] 中命中/闪避/暴击等键叠加后参与 [DamageSystem] 判定
	HitDeltaPermille         int
	DodgeDeltaPermille       int
	CritRateDeltaPermille    int
	CritDamageDeltaPermille  int
}

func (*StatModifiers) Component() {}
