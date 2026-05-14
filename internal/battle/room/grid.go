package room

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/land"
)

// Grid 将 [land.Grid] 与房间 [ecs.World] 绑定，负责同步空间占用与 [component.Transform2D]。
type Grid struct {
	*land.Grid
	world *ecs.World
}

// NewGrid 创建与给定 ECS 世界绑定的房间网格包装；base 或 w 为 nil 时返回 nil。
func NewGrid(w *ecs.World, base *land.Grid) *Grid {
	if w == nil || base == nil {
		return nil
	}
	return &Grid{Grid: base, world: w}
}

// FreeCell 按 placementSide 选择扫描顺序，返回第一个空闲格索引（与 [Grid.Add] 规则一致）。
func (g *Grid) FreeCell(placementSide component.SideType) (cellX, cellZ int, ok bool) {
	if g == nil || g.Grid == nil {
		return 0, 0, false
	}
	return g.PickFreeCell(placementSide == component.SideTypeRed)
}

// Add 在网格上为实体占第一个空闲格并写入 [component.Transform2D]。
func (g *Grid) Add(e ecs.Entity, placementSide component.SideType) error {
	if g == nil || g.Grid == nil || g.world == nil {
		return ErrNoSpatialGrid
	}
	w := g.world
	cellX, cellZ, ok := g.FreeCell(placementSide)
	if !ok {
		return ErrGridFull
	}
	_ = g.AddUnit(uint64(e), cellX, cellZ)
	w.AddComponent(e, &component.Transform2D{X: cellX, Y: cellZ})
	return nil
}

// Remove 从网格移除该实体的占用，并移除其 [component.Transform2D]（若存在）；不会 [RemoveEntity]。
func (g *Grid) Remove(e ecs.Entity) error {
	if g == nil || g.Grid == nil || g.world == nil {
		return ErrNoSpatialGrid
	}
	if e == 0 {
		return ErrInvalidEntity
	}
	if !g.RemoveUnitByID(uint64(e)) {
		return ErrEntityNotOnGrid
	}
	g.world.RemoveComponent(e, (*component.Transform2D)(nil))
	return nil
}
