package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

const threatPerDamagePoint = 1 // 可按游戏改为伤害 * 系数

// ThreatSystem 订阅伤害事件，向受害者 [ThreatBook] 累加来自攻击源的威胁值。
// 需在战斗单位上挂载 ThreatBook（通常仅 NPC/需要 AI 的单位）。
type ThreatSystem struct {
	world *ecs.World
}

func (s *ThreatSystem) Initialize(w *ecs.World) {
	s.world = w
	w.Subscribe(ecs.EventDamageApplied, func(ev ecs.Event) {
		if ev.Kind != ecs.EventDamageApplied || ev.Attacker == 0 || ev.IntPayload <= 0 {
			return
		}
		addThreat(s.world, ev.Entity, ev.Attacker, ev.IntPayload*threatPerDamagePoint)
	})
}

func addThreat(w *ecs.World, victim, source ecs.Entity, amt int) {
	c, ok := w.GetComponent(victim, &component.ThreatBook{})
	var tb *component.ThreatBook
	if ok {
		tb = c.(*component.ThreatBook)
	} else {
		tb = &component.ThreatBook{}
		w.AddComponent(victim, tb)
	}
	for i := range tb.Entries {
		if tb.Entries[i].Source == source {
			tb.Entries[i].Amount += amt
			return
		}
	}
	tb.Entries = append(tb.Entries, component.ThreatEntry{Source: source, Amount: amt})
}

func (s *ThreatSystem) Update(dt float64) {}
