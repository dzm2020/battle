package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// DeathSystem 对生命已耗尽的单位派发 [ecs.EventDeath] 并移除实体。
// 须在 [DamageSystem]、[HealthSystem] 之后更新（见 [AddCombatSystems] 顺序）。
type DeathSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.Attributes]
}

func (s *DeathSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.Attributes](w)
}

func (s *DeathSystem) Update(dt float64) {
	s.q.ForEach(func(e ecs.Entity, h *component.Attributes) {
		hp := h.Get(config.AttrHp)
		if hc, ok := s.world.GetComponent(e, &component.Health{}); ok {
			cur := hc.(*component.Health).Current
			if hp <= 0 && cur > 0 {
				hp = cur
			}
		}
		if hp > 0 {
			return
		}
		s.world.EmitEvent(ecs.Event{Kind: ecs.EventDeath, Entity: e})
		s.world.RemoveEntity(e)
	})
}
