package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/event"
	"battle/internal/battle/resource"
	"battle/internal/battle/system/attrs"
)

// BattleEndPayloadDraw 全员阵亡 / 同归于尽时 battle end 事件的 [event.Payload].IntPayload 取值。
const BattleEndPayloadDraw = -1

// BattleEndSystem 订阅 [event.KindDeath]，统计仍存活（[component.Team] + hp>0）的阵营数；
// 若 <= 1 则战斗结束并设置 [resource.PhaseSettled]。
//
// 须在 [DeathSystem] 之后注册（死亡事件在移除实体前派发，本 System 在回调中统计剩余单位）。
type BattleEndSystem struct {
	world      *ecs.World
	q          *ecs.Query2[*component.Team, *component.Attributes]
	done       bool
	deathUnsub func()
}

func (s *BattleEndSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery2[*component.Team, *component.Attributes](w)
	s.deathUnsub = w.Subscribe(event.KindDeath, func(ecs.Event) {
		s.onDeath()
	})
}

func (s *BattleEndSystem) Update(_ float64) {}

func (s *BattleEndSystem) onDeath() {
	if s.done || s.world == nil || s.q == nil {
		return
	}
	currSides, winningSide := countAliveSides(s.q)
	if currSides > 1 {
		return
	}
	if currSides == 0 {
		s.finish(BattleEndPayloadDraw)
		return
	}
	s.finish(sideToBattleEndPayload(winningSide))
}

func (s *BattleEndSystem) finish(winnerPayload int) {
	if s.done {
		return
	}
	s.done = true
	resource.SetPhase(s.world, resource.PhaseSettled)
	s.world.EmitEvent(ecs.Event{
		Kind:    event.KindBattleEnd,
		Payload: event.Payload{IntPayload: winnerPayload},
	})
}

// countAliveSides 返回有存活单位的阵营个数；若恰好 1 个阵营则返回该 Side。
func countAliveSides(q *ecs.Query2[*component.Team, *component.Attributes]) (count int, soleSide component.SideType) {
	seen := make(map[component.SideType]bool)
	q.ForEach(func(_ ecs.Entity, t *component.Team, h *component.Attributes) {
		if attrs.HP(h) <= 0 {
			return
		}
		seen[t.Side] = true
	})
	count = len(seen)
	if count == 1 {
		for side := range seen {
			soleSide = side
			break
		}
	}
	return count, soleSide
}

func sideToBattleEndPayload(side component.SideType) int {
	switch side {
	case component.SideTypeRed:
		return 0
	case component.SideTypeBlue:
		return 1
	default:
		return BattleEndPayloadDraw
	}
}
