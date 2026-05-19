package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/system/attrs"
)

// HealthSystem 消费 [ResolvedDamage]，从 [Attributes] 的 hp 扣减（唯一生命数据源），然后移除 [ResolvedDamage]。
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
		attrs.Sub(attr, config.AttrHp, rd.Amount)
		s.world.RemoveComponent(e, &component.ResolvedDamage{})
	})
}
