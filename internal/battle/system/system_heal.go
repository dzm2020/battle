package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// HealSystem 消费 [PendingHeal]，直接增加 [Health].Current 并派发 [EventHealApplied]；在 [DamageSystem] 之后、
// [HealthSystem] 之前运行，使同帧先扣血再治疗时治疗仍有效（与项目系统顺序一致时先伤后疗由注册顺序决定）。
type HealSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.PendingHeal]
}

func (s *HealSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.PendingHeal](w)
}

func (s *HealSystem) Update(dt float64) {
	s.q.ForEach(func(e ecs.Entity, ph *component.PendingHeal) {
		if ph.Amount <= 0 {
			s.world.RemoveComponent(e, &component.PendingHeal{})
			return
		}
		h, ok := s.world.GetComponent(e, &component.Health{})
		if !ok {
			s.world.RemoveComponent(e, &component.PendingHeal{})
			return
		}
		hp := h.(*component.Health)
		hp.Current += ph.Amount
		if hp.Current > hp.Max {
			hp.Current = hp.Max
		}
		if a, ok := s.world.GetComponent(e, &component.Attributes{}); ok {
			a.(*component.Attributes).Set(config.AttrHp, hp.Current)
		}
		src := ph.Source
		s.world.EmitEvent(ecs.Event{
			Kind:       ecs.EventHealApplied,
			Entity:     e,
			Attacker:   src,
			IntPayload: ph.Amount,
		})
		s.world.RemoveComponent(e, &component.PendingHeal{})
	})
}
