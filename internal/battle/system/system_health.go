package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// HealthSystem 应用 [ResolvedDamage]，派发 [ecs.EventDamageApplied]，再移除 ResolvedDamage。
type HealthSystem struct {
	world *ecs.World
	q     *ecs.Query2[*component.ResolvedDamage, *component.Attributes]
}

func (s *HealthSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery2[*component.ResolvedDamage, *component.Attributes](w)
}

func (s *HealthSystem) Update(dt float64) {
	s.q.ForEach(func(e ecs.Entity, rd *component.ResolvedDamage, h *component.Attributes) {
		if rd.Amount <= 0 {
			s.world.RemoveComponent(e, &component.ResolvedDamage{})
			return
		}

		hp := h.Get(config.AttrHp)
		if hc, ok := s.world.GetComponent(e, &component.Health{}); ok {
			cur := hc.(*component.Health).Current
			if hp <= 0 && cur > 0 {
				hp = cur
			}
		}

		if rd.Amount > hp {
			hp = 0
		} else {
			hp -= rd.Amount
		}

		h.Set(config.AttrHp, hp)

		if hc, ok := s.world.GetComponent(e, &component.Health{}); ok {
			cur := hc.(*component.Health)
			cur.Current = hp
			if cur.Max > 0 && cur.Current > cur.Max {
				cur.Current = cur.Max
			}
		}

		s.world.EmitEvent(ecs.Event{
			Kind:       ecs.EventDamageApplied,
			Entity:     e,
			Attacker:   rd.Source,
			IntPayload: rd.Amount,
		})
		s.world.RemoveComponent(e, &component.ResolvedDamage{})
	})
}
