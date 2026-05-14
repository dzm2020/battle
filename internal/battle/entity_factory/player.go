package entity_factory

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/pb"
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
		e, err := spawnUnitFromPBUnit(w, unit, components...)
		if err != nil {
			return 0, err
		}
		playerComponents.Units[u] = e
	}
	return pe, nil
}
