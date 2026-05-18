package room_builder

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/room"
	"battle/internal/battle/system"
)

func pveBuilder(spec *Spec) error {
	if options == nil {
		return nil
	}

	if err := spawnMonstersForDesc(r, desc, component.SideTypeBlue); err != nil {
		return err
	}

	if err := spawnPlayersOnGridWithTeam(r, options.Self, component.SideTypeRed); err != nil {
		return err
	}

	system.AddCombatSystems(r.World())
	return nil
}
