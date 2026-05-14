package room

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// PlaceOnGrid 将已存在于本房间 [Room.World] 的实体占住 [Room.Grid] 上按 placementSide 策略找到的第一个空闲格，并写入 [component.Transform2D]。
// 允许阶段：Lobby、PreBattle、Fighting（便于战中刷怪）；Settled/Closed 返回 [ErrWrongPhase]。
//
// placementSide 决定扫描顺序：红方 [component.SideTypeRed] 使用网格升序首个空闲格，
// 其它取值（含蓝方、空字符串）使用降序首个空闲格，与 PVP 两侧落点习惯一致；纯 PVE 刷怪可用蓝或任意非红值。
func (r *Room) PlaceOnGrid(e ecs.Entity, placementSide component.SideType) error {
	if r.phaseIs(PhaseClosed) || r.phaseIs(PhaseSettled) {
		return ErrWrongPhase
	}
	if !r.phaseIs(PhaseLobby) && !r.phaseIs(PhasePreBattle) && !r.phaseIs(PhaseFighting) {
		return ErrWrongPhase
	}
	grid := r.grid
	if grid == nil {
		return ErrNoSpatialGrid
	}
	w := r.world

	cellX, cellZ, ok := grid.PickFreeCell(placementSide == component.SideTypeRed)
	if !ok {
		return ErrGridFull
	}
	_ = grid.AddUnit(uint64(e), cellX, cellZ)
	w.AddComponent(e, &component.Transform2D{X: cellX, Y: cellZ})
	return nil
}
