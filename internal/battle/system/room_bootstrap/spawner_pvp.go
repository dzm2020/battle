package room_bootstrap

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/resource"
)

func spawnerPVP(w *ecs.World, spec *resource.RoomSpec) error {
	if err := spawnPlayerUnits(w, spec.Enemy, component.SideTypeBlue); err != nil {
		return err
	}
	return spawnPlayerUnits(w, spec.Self, component.SideTypeRed)
}
