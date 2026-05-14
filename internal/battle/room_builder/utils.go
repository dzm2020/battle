package room_builder

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/entity_factory"
	"battle/internal/battle/land"
	"battle/internal/battle/pb"
	"battle/internal/battle/room"
	"battle/internal/battle/utils"
)

func createLandByID(mapID int32) (*land.Grid, error) {
	mapDesc := config.GetMapConfigByID(mapID)
	if mapDesc == nil {
		return nil, ErrNoMapConfig
	}
	return land.NewSpatialGrid(mapDesc.MinX, mapDesc.MinZ, mapDesc.MaxX, mapDesc.MaxZ, mapDesc.CellSize)
}

// spawnPlayersOnGridWithTeam 将玩家单位放入网格；assignTeam 为 true 时为每个单位挂载 [component.Team]。
func spawnPlayersOnGridWithTeam(r *room.Room, player *pb.Player, side component.SideType) error {
	word := r.World()
	grid := r.Grid()

	pe := word.CreateEntity()

	playerComponents := &component.Player{
		ID:    player.ID,
		Base:  player.Base,
		Units: make(map[uint32]ecs.Entity),
	}
	word.AddComponent(pe, playerComponents)

	for _, unit := range player.Units {
		cellX, cellZ, ok := utils.GetLandFreeCell(grid, side)
		if !ok {
			continue
		}
		components := []ecs.Component{
			&component.Team{Side: side},
			&component.Transform2D{X: cellX, Y: cellZ},
		}
		e, err := entity_factory.CreateByUnit(word, unit, components...)
		if err != nil {
			return err
		}
		if err = grid.AddUnit(uint64(e), cellX, cellZ); err != nil {
			return err
		}
		playerComponents.Units[unit.ID] = e
	}

	return nil
}
