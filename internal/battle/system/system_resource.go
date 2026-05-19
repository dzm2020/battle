package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/system/attrs"
)

// ResourceSystem 处理战斗资源（法力、怒气、能量等）的消耗与自然恢复。
//
// 流程（每帧）：
//  1. 自然恢复（不超过 Max）
//  2. 消费 [component.ResourceConsumeQueue]（由 [CastValidationSystem] 施法校验通过后入队）
//
// 须注册在 [CastValidationSystem] 之后、[CastStateSystem] 之前。
type ResourceSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.Attributes]
}

func (s *ResourceSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.Attributes](w)
}

func (s *ResourceSystem) Update(_ float64) {
	if s.world == nil || s.q == nil {
		return
	}
	s.q.ForEach(func(e ecs.Entity, attr *component.Attributes) {
		if attr == nil {
			return
		}
		s.applyRegen(e, attr)
		s.applyConsumeQueue(e, attr)
	})
}

func (s *ResourceSystem) applyConsumeQueue(e ecs.Entity, attr *component.Attributes) {
	q, ok := s.world.GetComponent(e, &component.ResourceConsumeQueue{})
	if !ok {
		return
	}
	cq := q.(*component.ResourceConsumeQueue)
	if len(cq.Entries) == 0 {
		return
	}
	for _, entry := range cq.Entries {
		if entry.Amount <= 0 || entry.Type == "" {
			continue
		}
		if !attrs.CanAfford(attr, entry.Type, entry.Amount) {
			continue
		}
		attrs.ApplyConsume(attr, entry.Type, entry.Amount)
	}
	cq.Entries = nil
}

func (s *ResourceSystem) applyRegen(e ecs.Entity, attr *component.Attributes) {
	var regen *component.ResourceRegen
	if r, ok := s.world.GetComponent(e, &component.ResourceRegen{}); ok {
		regen = r.(*component.ResourceRegen)
	}
	rates := attrs.RegenRatesFor(attr, regen)
	attrs.ApplyRegen(attr, rates)
}
