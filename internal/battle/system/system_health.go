package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// HealthSystem 应用 [ResolvedDamage]，派发 [ecs.EventDamageApplied]，再移除 ResolvedDamage。
type HealthSystem struct {
	world *ecs.World
	q     *ecs.Query2[*component.ResolvedDamage, *component.Health]
}

func (s *HealthSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery2[*component.ResolvedDamage, *component.Health](w)
}

func (s *HealthSystem) Update(dt float64) {
	s.q.ForEach(func(e ecs.Entity, rd *component.ResolvedDamage, h *component.Health) {
		if rd.Amount <= 0 {
			s.world.RemoveComponent(e, &component.ResolvedDamage{})
			return
		}
		if rd.Amount > h.Current {
			h.Current = 0
		} else {
			h.Current -= rd.Amount
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
