package system

import (
	"battle/ecs"
	"fmt"
)

// CombatLogEntry 一条战斗日志（便于回放或 UI）。
type CombatLogEntry struct {
	Kind    string       // miss / damage / heal / death
	Target  ecs.Entity
	Source  ecs.Entity
	Value   int
	Message string
}

// CombatLogSystem 订阅战斗事件并写入环形缓冲；MaxEntries 满后丢弃最旧条目。
type CombatLogSystem struct {
	world      *ecs.World
	Entries    []CombatLogEntry
	MaxEntries int
	cancelFns  []func()
}

func NewCombatLogSystem(maxEntries int) *CombatLogSystem {
	if maxEntries <= 0 {
		maxEntries = 512
	}
	return &CombatLogSystem{MaxEntries: maxEntries}
}

func (s *CombatLogSystem) Initialize(w *ecs.World) {
	s.world = w
	s.cancelFns = append(s.cancelFns, w.Subscribe(ecs.EventDamageMissed, func(ev ecs.Event) {
		s.push(CombatLogEntry{
			Kind: "miss", Target: ev.Entity, Source: ev.Attacker,
			Message: fmt.Sprintf("miss victim=%v attacker=%v", ev.Entity, ev.Attacker),
		})
	}))
	s.cancelFns = append(s.cancelFns, w.Subscribe(ecs.EventDamageApplied, func(ev ecs.Event) {
		s.push(CombatLogEntry{
			Kind: "damage", Target: ev.Entity, Source: ev.Attacker, Value: ev.IntPayload,
			Message: fmt.Sprintf("damage %d victim=%v attacker=%v", ev.IntPayload, ev.Entity, ev.Attacker),
		})
	}))
	s.cancelFns = append(s.cancelFns, w.Subscribe(ecs.EventHealApplied, func(ev ecs.Event) {
		s.push(CombatLogEntry{
			Kind: "heal", Target: ev.Entity, Source: ev.Attacker, Value: ev.IntPayload,
			Message: fmt.Sprintf("heal %d target=%v source=%v", ev.IntPayload, ev.Entity, ev.Attacker),
		})
	}))
	s.cancelFns = append(s.cancelFns, w.Subscribe(ecs.EventDeath, func(ev ecs.Event) {
		s.push(CombatLogEntry{
			Kind: "death", Target: ev.Entity,
			Message: fmt.Sprintf("death entity=%v", ev.Entity),
		})
	}))
	s.cancelFns = append(s.cancelFns, w.Subscribe(ecs.EventBattleEnd, func(ev ecs.Event) {
		msg := fmt.Sprintf("battle_end winner_payload=%d", ev.IntPayload)
		if ev.IntPayload == BattleEndPayloadDraw {
			msg = "battle_end draw"
		}
		s.push(CombatLogEntry{
			Kind: "battle_end",
			Message: msg,
		})
	}))
}

func (s *CombatLogSystem) push(e CombatLogEntry) {
	if len(s.Entries) >= s.MaxEntries {
		s.Entries = s.Entries[1:]
	}
	s.Entries = append(s.Entries, e)
}

func (s *CombatLogSystem) Update(dt float64) {}
