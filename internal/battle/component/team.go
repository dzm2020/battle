package component

// Team 战斗阵营/队伍标识，用于技能 [skill.TargetAllEnemySides] 等筛选“敌方”。
// 未挂载 Team 的实体不参与阵营差集逻辑（单体技能仍可通过显式目标攻击）。
type Team struct {
	// Side 阵营编号；仅在同局内比较有区分意义，相同 Side 互为友方。
	Side uint8
}

func (*Team) Component() {}
