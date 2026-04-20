package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// BattleEndPayloadDraw 全员阵亡 / 同归于尽时 [ecs.EventBattleEnd].IntPayload 取值。
const BattleEndPayloadDraw = -1

// BattleEndSystem 根据仍存活的阵营数判定战斗是否结束，并派发 [ecs.EventBattleEnd]。
// 仅统计同时挂载 [component.Team] 与 [component.Health] 且 Current > 0 的实体（参战单位）。
// IntPayload：获胜方 [component.Team].Side（0–255）；平局为 [BattleEndPayloadDraw]。
//
// 在 [Initialize] 时快照「开局」存活阵营数（须在本系统随 [AddCombatSystems] 注册前已放入参战实体，否则 openingSides 恒为 0）；
// 首帧仅用其建立基线，避免首帧击杀时无法产生「上一帧仍为多阵营」的过渡。
//
// 须注册在 [DeathSystem] 之后（见 [AddCombatSystems]），以便死亡实体已从世界移除后再统计。
type BattleEndSystem struct {
	world        *ecs.World
	q            *ecs.Query2[*component.Team, *component.Health]
	done         bool
	prevSides    int // -1 未建立基线；之后为上一帧结算后的存活阵营数
	openingSides int // Initialize 时的存活阵营数（世界已含实体时应先加人再 Register 本系统）
}

func (s *BattleEndSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery2[*component.Team, *component.Health](w)
	s.prevSides = -1
	s.openingSides, _ = countAliveSides(s.q)
}

func (s *BattleEndSystem) Update(dt float64) {
	if s.done {
		return
	}

	currSides, winningSide := countAliveSides(s.q)

	if s.prevSides < 0 {
		if s.openingSides >= 2 {
			s.prevSides = s.openingSides
		} else {
			s.prevSides = currSides
		}
		return
	}

	switch {
	case currSides == 0 && s.prevSides > 0:
		s.finish(BattleEndPayloadDraw)
	case s.prevSides >= 2 && currSides == 1:
		s.finish(int(winningSide))
	default:
		s.prevSides = currSides
	}
}

func (s *BattleEndSystem) finish(winnerPayload int) {
	s.done = true
	s.world.EmitEvent(ecs.Event{
		Kind:       ecs.EventBattleEnd,
		IntPayload: winnerPayload,
	})
}

// countAliveSides 返回有存活单位的阵营个数；若恰好 1 个阵营则返回该 Side。
func countAliveSides(q *ecs.Query2[*component.Team, *component.Health]) (count int, soleSide uint8) {
	seen := make(map[uint8]bool)
	q.ForEach(func(_ ecs.Entity, t *component.Team, h *component.Health) {
		if h.Current <= 0 {
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
