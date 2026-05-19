package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/system/attrs"
)

// AttributeSystem 在 [BuffSystem] 之后根据 [component.Attributes] 与 [component.BuffStatModifiers]
// 重算最终属性并写入 [component.FinalAttributes]，供 [DamageSystem] 等消费。
type AttributeSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.Attributes]
}

func (s *AttributeSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.Attributes](w)
}

func (s *AttributeSystem) Update(_ float64) {
	if s.world == nil || s.q == nil {
		return
	}
	s.q.ForEach(func(e ecs.Entity, base *component.Attributes) {
		if base == nil {
			return
		}
		var mods *component.BuffStatModifiers
		if m, ok := s.world.GetComponent(e, &component.BuffStatModifiers{}); ok {
			mods = m.(*component.BuffStatModifiers)
		}
		fa := ecs.EnsureGetComponent[*component.FinalAttributes](s.world, e)
		fa.Values = attrs.Recompute(base, mods)
	})
}
