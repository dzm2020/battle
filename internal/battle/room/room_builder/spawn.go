package room_builder

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/pb"
	"battle/internal/battle/room"
	"battle/internal/battle/unit"
)

// spawnMonstersForDesc 按副本配置刷怪；无空位时跳过该只怪（与原先 pve 循环 continue 一致）。
func spawnMonstersForDesc(r *room.Room, desc *config.DungeonConfig, side component.SideType) error {
	for _, monsterID := range desc.Monster {
		if err := spawnMonsterAtFreeCell(r, monsterID, side); err != nil {
			return err
		}
	}
	return nil
}

func spawnMonsterAtFreeCell(r *room.Room, monsterID int32, side component.SideType) error {
	grid := r.Grid()
	if grid == nil {
		return room.ErrNoSpatialGrid
	}
	cellX, cellZ, ok := grid.PickFreeCell(side == component.SideTypeRed)
	if !ok {
		return nil
	}
	components := []ecs.Component{
		&component.Team{Side: side},
		&component.Transform2D{X: cellX, Y: cellZ},
	}
	e, err := unit.CreateByID(r.World(), monsterID, components...)
	if err != nil {
		return err
	}
	return grid.AddUnit(uint64(e), cellX, cellZ)
}

// spawnPlayersOnGridWithTeam 将单个 [pb.Player] 的单位放入网格；单位带 [component.Team]（含队长实体引用）。
func spawnPlayersOnGridWithTeam(r *room.Room, player *pb.Player, side component.SideType) error {
	if player == nil {
		return nil
	}
	w := r.World()
	grid := r.Grid()
	if grid == nil {
		return room.ErrNoSpatialGrid
	}

	teamEntity := w.CreateEntity()
	pc := &component.Player{
		ID:    player.ID,
		Base:  player.Base,
		Units: make(map[uint32]ecs.Entity),
	}
	w.AddComponent(teamEntity, pc)

	for _, unit := range player.Units {
		if unit == nil {
			continue
		}
		cellX, cellZ, ok := grid.PickFreeCell(side == component.SideTypeRed)
		if !ok {
			continue
		}
		components := []ecs.Component{
			&component.Team{Side: side, Entity: teamEntity},
			&component.Transform2D{X: cellX, Y: cellZ},
		}
		e, err := unit.CreateByUnit(w, unit, components...)
		if err != nil {
			return err
		}
		if err = grid.AddUnit(uint64(e), cellX, cellZ); err != nil {
			return err
		}
		pc.Units[unit.ID] = e
	}
	return nil
}
