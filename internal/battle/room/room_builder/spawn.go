package room_builder

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/land"
	"battle/internal/battle/pb"
	"errors"
)

// spawnMonstersForDesc 按副本配置刷怪；无空位时跳过该只怪（与原先 pve 循环 continue 一致）。
func spawnMonstersForDesc(spec *Spec, side component.SideType) error {
	w := spec.World
	desc := spec.Desc
	queue, _ := ecs.GetResource[*component.SpawnRequestQueue](w)
	grid, _ := ecs.GetResource[*land.Grid](w)
	for _, monsterID := range desc.Monster {
		cellX, cellZ, ok := grid.PickFreeCell(side == component.SideTypeRed)
		if !ok {
			return errors.New("failed to pick cell")
		}

		queue.Queue = append(queue.Queue, &component.SpawnRequest{
			UnitID: monsterID,
			Side:   side,
			CellX:  cellX,
			CellY:  cellZ,
		})
	}
	return nil
}

// spawnPlayersOnGridWithTeam 将单个 [pb.Player] 的单位放入网格；单位带 [component.Team]（含队长实体引用）。
func spawnPlayersOnGridWithTeam(spec *Spec, player *pb.Player, side component.SideType) error {
	w := spec.World

	queue, _ := ecs.GetResource[*component.SpawnRequestQueue](w)
	grid, _ := ecs.GetResource[*land.Grid](w)

	for _, unit := range player.Units {
		if unit == nil {
			return errors.New("nil unit")
		}
		cellX, cellZ, ok := grid.PickFreeCell(side == component.SideTypeRed)
		if !ok {
			return errors.New("failed to pick cell")
		}
		queue.Queue = append(queue.Queue, &component.SpawnRequest{
			UnitID: int32(unit.ID),
			Side:   side,
			CellX:  cellX,
			CellY:  cellZ,
			Data:   unit,
		})
	}
	return nil
}
