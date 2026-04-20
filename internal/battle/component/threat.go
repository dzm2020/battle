package component

import "battle/ecs"

// ThreatBook 目标实体上记录的「来自各攻击源的威胁值」，用于 RPG 仇恨与 AI 选敌。
// 仅在需要仇恨的单位上挂载；同一 Source 合并累加 Amount。
type ThreatBook struct {
	Entries []ThreatEntry
}

// ThreatEntry 单条仇恨记录。
type ThreatEntry struct {
	Source ecs.Entity
	Amount int
}

func (*ThreatBook) Component() {}

// ThreatTopSource 返回威胁值最高的攻击来源；空表或无条目返回 0。
func ThreatTopSource(tb *ThreatBook) ecs.Entity {
	if tb == nil || len(tb.Entries) == 0 {
		return 0
	}
	best := tb.Entries[0].Source
	max := tb.Entries[0].Amount
	for i := 1; i < len(tb.Entries); i++ {
		if tb.Entries[i].Amount > max {
			max = tb.Entries[i].Amount
			best = tb.Entries[i].Source
		}
	}
	return best
}
