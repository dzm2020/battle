package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// HealthSystem 消费 [ResolvedDamage]，从 [Attributes] 的 hp 扣减并同步 [Health]（若存在），然后移除 [ResolvedDamage]。
type HealthSystem struct {
	world *ecs.World
	q     *ecs.Query2[*component.ResolvedDamage, *component.Attributes]
}

func (s *HealthSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery2[*component.ResolvedDamage, *component.Attributes](w)
}

func (s *HealthSystem) Update(dt float64) {
	s.q.ForEach(func(e ecs.Entity, rd *component.ResolvedDamage, attr *component.Attributes) {
		if rd.Amount <= 0 {
			s.world.RemoveComponent(e, &component.ResolvedDamage{})
			return
		}
		attr.Sub(config.AttrHp, rd.Amount)
		if h, ok := s.world.GetComponent(e, &component.Health{}); ok {
			hc := h.(*component.Health)
			hc.Current = attr.Get(config.AttrHp)
			if hc.Current < 0 {
				hc.Current = 0
			}
			if hc.Max > 0 && hc.Current > hc.Max {
				hc.Current = hc.Max
			}
		}
		s.world.RemoveComponent(e, &component.ResolvedDamage{})
	})
}
