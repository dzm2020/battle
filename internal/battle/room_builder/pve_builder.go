package room_builder

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/entity_factory"
	"battle/internal/battle/room"
	"battle/internal/battle/system"
	"battle/internal/battle/utils"
)

func pveBuilder(r *room.Room, desc *config.DungeonConfig, options *Options) error {
	if options == nil {
		return nil
	}

	grid := r.Grid()
	for _, monsterID := range desc.Monster {
		cellX, cellZ, ok := utils.GetLandFreeCell(grid, component.SideTypeBlue)
		if !ok {
			continue
		}
		components := []ecs.Component{
			&component.Team{Side: component.SideTypeBlue},
			&component.Transform2D{X: cellX, Y: cellZ},
		}
		e, err := entity_factory.CreateByID(r.World(), monsterID, components...)
		if err != nil {
			return err
		}
		if err = grid.AddUnit(uint64(e), cellX, cellZ); err != nil {
			return err
		}
	}

	if err := spawnPlayersOnGridWithTeam(r, options.Self, component.SideTypeRed); err != nil {
		return err
	}

	system.AddCombatSystems(r.World())

	return nil
}
