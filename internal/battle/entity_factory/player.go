package entity_factory

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
	"battle/internal/battle/pb"
	"battle/internal/battle/skill"
	"fmt"
)

func CreateByPlayer(w *ecs.World, p *pb.Player, components ...ecs.Component) (ecs.Entity, error) {
	pe := w.CreateEntity()

	playerComponents := &component.Player{
		ID:    p.ID,
		Base:  p.Base,
		Units: make(map[uint32]ecs.Entity),
	}
	w.AddComponent(pe, playerComponents)

	for u, unit := range p.Units {
		e := w.CreateEntity()

		initAttrComponentFromStats(w, e, unit.Stats)

		if err := skill.Add(w, e, unit.Ability...); err != nil {
			w.RemoveEntity(e)
			return 0, fmt.Errorf("unit: 初始技能  挂载失败: %w", err)
		}

		if err := buff.Add(w, e, e, unit.BuffDefIDs...); err != nil {
			w.RemoveEntity(e)
			return 0, fmt.Errorf("unit: 初始 Buff 挂载失败: %w", err)
		}
		for _, com := range components {
			w.AddComponent(e, com)
		}
		playerComponents.Units[u] = e
		return e, nil
	}
	return pe, nil
}

func initAttrComponentFromStats(w *ecs.World, e ecs.Entity, stats []pb.Attribute) {
	if len(stats) == 0 {
		return
	}
	attrs := ecs.EnsureGetComponent[*component.Attributes](w, e)
	for _, attr := range stats {
		attrs.SetRange(attr.Type, attr.InitValue, attr.MaxValue)
	}
	w.AddComponent(e, attrs)
}
