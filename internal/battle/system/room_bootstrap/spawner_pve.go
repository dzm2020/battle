package room_bootstrap

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/resource"
)

func spawnerPVE(w *ecs.World, spec *resource.RoomSpec) error {
	if err := spawnDungeonMonsters(w, spec, component.SideTypeBlue); err != nil {
		return err
	}
	return spawnPlayerUnits(w, spec.Self, component.SideTypeRed)
}
