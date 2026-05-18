package room_factory

import (
	"battle/internal/battle/component"
	"battle/internal/battle/pb"
	"battle/internal/battle/system/runtime"
	"errors"
)

// spawnMonstersForDesc 按副本配置刷怪；无空位时返回错误。
func spawnMonstersForDesc(spec *Spec, side component.SideType) error {
	w := spec.World
	desc := spec.Desc
	ctx, err := runtime.MustGet(w)
	if err != nil {
		return err
	}
	if ctx.Grid == nil {
		return runtime.ErrNilGrid
	}
	for _, monsterID := range desc.Monster {
		cellX, cellZ, ok := ctx.Grid.PickFreeCell(side == component.SideTypeRed)
		if !ok {
			return errors.New("failed to pick cell")
		}
		if err := runtime.EnqueueSpawn(w, &component.SpawnRequest{
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
func spawnPlayersOnGridWithTeam(spec *Spec, player *pb.Player, side component.SideType) error {
	if player == nil {
		return nil
	}
	w := spec.World
	ctx, err := runtime.MustGet(w)
	if err != nil {
		return err
	}
	if ctx.Grid == nil {
		return runtime.ErrNilGrid
	}

	for _, u := range player.Units {
		if u == nil {
			return errors.New("nil unit")
		}
		cellX, cellZ, ok := ctx.Grid.PickFreeCell(side == component.SideTypeRed)
		if !ok {
			return errors.New("failed to pick cell")
		}
		if err := runtime.EnqueueSpawn(w, &component.SpawnRequest{
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
