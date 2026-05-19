package room_bootstrap

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/pb"
	"battle/internal/battle/resource"
	"errors"
	"fmt"
)

// spawnDungeonMonsters 按副本配置将怪物入队；无空位时返回错误。
func spawnDungeonMonsters(w *ecs.World, spec *resource.RoomSpec, side component.SideType) error {
	desc := config.GetDungeonConfigByID(spec.DungeonId)
	if desc == nil {
		return fmt.Errorf("room_bootstrap: dungeon config not found: %d", spec.DungeonId)
	}
	grid, ok := resource.Grid(w)
	if !ok || grid == nil {
		return errors.New("room_bootstrap: grid not initialized")
	}

	for _, unitID := range desc.Monster {
		cellX, cellY, ok := grid.PickFreeCell(side == component.SideTypeRed)
		if !ok {
			return errors.New("room_bootstrap: failed to pick cell for monster")
		}
		if err := resource.EnqueueSpawn(w, &resource.SpawnRequest{
			UnitID: unitID,
			Side:   side,
			CellX:  cellX,
			CellY:  cellY,
		}); err != nil {
			return err
		}
	}
	return nil
}

// spawnPlayerUnits 将单个 [pb.Player] 的战斗单位入队；player 为 nil 时跳过。
// 先创建带 [component.Player] 的编队实体，刷怪时写入 Player.Units 映射。
func spawnPlayerUnits(w *ecs.World, player *pb.Player, side component.SideType) error {
	if player == nil {
		return nil
	}
	grid, ok := resource.Grid(w)
	if !ok || grid == nil {
		return errors.New("room_bootstrap: grid not initialized")
	}

	teamEntity := w.CreateEntity()
	w.AddComponent(teamEntity, &component.Player{
		ID:    player.ID,
		Base:  player.Base,
		Units: make(map[uint32]ecs.Entity),
	})

	for _, unit := range player.Units {
		if unit == nil {
			return errors.New("room_bootstrap: nil player unit")
		}
		cellX, cellY, ok := grid.PickFreeCell(side == component.SideTypeRed)
		if !ok {
			return errors.New("room_bootstrap: failed to pick cell for player unit")
		}
		if err := resource.EnqueueSpawn(w, &resource.SpawnRequest{
			UnitID:     int32(unit.ID),
			Side:       side,
			CellX:      cellX,
			CellY:      cellY,
			TeamEntity: teamEntity,
			Data:       unit,
		}); err != nil {
			return err
		}
	}
	return nil
}
