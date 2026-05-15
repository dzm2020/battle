package unit

import (
	"battle/ecs"
	"battle/internal/battle/pb"
)

func CreateByUnit(w *ecs.World, unit *pb.PlayerUnit, components ...ecs.Component) (ecs.Entity, error) {
	return spawnUnitFromPBUnit(w, unit, components...)
}
