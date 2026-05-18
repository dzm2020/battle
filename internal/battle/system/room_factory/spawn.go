package room_factory

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/pb"
	"battle/internal/battle/resource"
	"errors"
	"fmt"
)

// spawnMonstersForDesc 按副本配置刷怪；无空位时返回错误。
func spawnMonstersForDesc(w *ecs.World, spec *resource.RoomSpec, side component.SideType) error {
	desc := config.GetDungeonConfigByID(spec.DungeonId)
	if desc == nil {
		return fmt.Errorf("dungeon config not found: %d", spec.DungeonId)
	}
	grid, ok := resource.Grid(w)
	if !ok || grid == nil {
		return errors.New("grid not initialized")
	}

	for _, monsterID := range desc.Monster {
		cellX, cellZ, ok := grid.PickFreeCell(side == component.SideTypeRed)
		if !ok {
			return errors.New("failed to pick cell")
		}
		if err := resource.EnqueueSpawn(w, &resource.SpawnRequest{
			UnitID: monsterID,
			Side:   side,
			CellX:  cellX,
			CellY:  cellZ,
		}); err != nil {
			return err
		}
	}
	return nil
}

// spawnPlayersOnGridWithTeam 将单个 [pb.Player] 的单位入队刷怪请求。
func spawnPlayersOnGridWithTeam(w *ecs.World, player *pb.Player, side component.SideType) error {
	if player == nil {
		return nil
	}
	grid, ok := resource.Grid(w)
	if !ok || grid == nil {
		return errors.New("grid not initialized")
	}

	for _, u := range player.Units {
		if u == nil {
			return errors.New("nil unit")
		}
		cellX, cellZ, ok := grid.PickFreeCell(side == component.SideTypeRed)
		if !ok {
			return errors.New("failed to pick cell")
		}
		if err := resource.EnqueueSpawn(w, &resource.SpawnRequest{
			UnitID: int32(u.ID),
			Side:   side,
			CellX:  cellX,
			CellY:  cellZ,
			Data:   u,
		}); err != nil {
			return err
		}
	}
	return nil
}
