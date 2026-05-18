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
	q     *ecs.Query2[*component.PendingHeal, *component.Attributes]
}

func (s *HealSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery2[*component.PendingHeal, *component.Attributes](w)
}

func (s *HealSystem) Update(dt float64) {
	s.q.ForEach(func(e ecs.Entity, ph *component.PendingHeal, attr *component.Attributes) {
		if ph.Amount <= 0 {
			return
		}
		hp := component.AttrCurrent(attr, config.AttrHp)
		//  死亡就别治疗了
		if hp <= 0 {
			return
		}
		component.AttrAdd(attr, config.AttrHp, ph.Amount)

		s.world.RemoveComponent(e, &component.PendingHeal{})
	})
}
