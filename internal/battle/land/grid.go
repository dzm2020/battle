// Package land 提供均匀网格的空间分区（水平 XZ），用于 AOI / 邻近单位查询。
package land

import (
	"battle/internal/battle/config"
	"errors"
	"fmt"
)

// ErrInvalidGridConfig 边界或 cellSize 非法。
var ErrInvalidGridConfig = errors.New("spatial: invalid bounds or cellSize")

// Unit 网格中的可移动对象；ID 可与 [battle/ecs.Entity]（uint64）对应。
// 逻辑位置由所在格子索引 cellX、cellZ 表示。
type Unit struct {
	ID    uint64
	cellX int
	cellZ int
}

// GridCell 单个网格单元内的单位集合。
type GridCell struct {
	Units map[uint64]*Unit
}

// Empty 格子内无任何单位。
func (c *GridCell) Empty() bool {
	return c == nil || len(c.Units) == 0
}

// Grid 二维均匀网格空间分区（XZ 平面）。
type Grid struct {
	cells    [][]*GridCell
	cellSize float64
	width    int
	height   int
	minX     float64
	minZ     float64
}

func CreateGridByID(mapID int32) (*Grid, error) {
	mapDesc := config.GetMapConfigByID(mapID)
	if mapDesc == nil {
		return nil, fmt.Errorf("map id %v not found", mapID)
	}
	return NewSpatialGrid(mapDesc.MinX, mapDesc.MinZ, mapDesc.MaxX, mapDesc.MaxZ, mapDesc.CellSize)
}

// NewSpatialGrid 创建覆盖 [minX,maxX)×[minZ,maxZ) 的网格；要求 maxX>minX、maxZ>minZ、cellSize>0。
func NewSpatialGrid(minX, minZ, maxX, maxZ float64, cellSize float64) (*Grid, error) {
	if cellSize <= 0 || maxX <= minX || maxZ <= minZ {
		return nil, ErrInvalidGridConfig
	}
	width := int((maxX-minX)/cellSize) + 1
	height := int((maxZ-minZ)/cellSize) + 1
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	cells := make([][]*GridCell, width)
	for i := range cells {
		cells[i] = make([]*GridCell, height)
		for j := range cells[i] {
			cells[i][j] = &GridCell{
				Units: make(map[uint64]*Unit),
			}
		}
	}

	return &Grid{
		cells:    cells,
		cellSize: cellSize,
		width:    width,
		height:   height,
		minX:     minX,
		minZ:     minZ,
	}, nil
}

// CellSize 单元格边长。
func (g *Grid) CellSize() float64 { return g.cellSize }

// Width / Height 单元格数量。
func (g *Grid) Width() int  { return g.width }
func (g *Grid) Height() int { return g.height }

func (g *Grid) cellIndex(x, z float64) (int, int) {
	cx := int((x - g.minX) / g.cellSize)
	cz := int((z - g.minZ) / g.cellSize)
	if cx < 0 {
		cx = 0
	}
	if cx >= g.width {
		cx = g.width - 1
	}
	if cz < 0 {
		cz = 0
	}
	if cz >= g.height {
		cz = g.height - 1
	}
	return cx, cz
}

// AddUnit 将单位放入当前坐标所在格子，并写入 unit 内缓存的格子索引。必须先 Add 再 Update。
func (g *Grid) AddUnit(unitID uint64, cx int, cz int) error {
	unit := &Unit{ID: unitID}

	cell := g.cells[cx][cz]
	cell.Units[unit.ID] = unit

	unit.cellX, unit.cellZ = cx, cz
	return nil
}

// RemoveUnit 从网格移除单位。未加入过网格时无副作用。
func (g *Grid) RemoveUnit(unit *Unit) {
	if unit == nil {
		return
	}
	cx, cz := unit.cellX, unit.cellZ
	if cx < 0 || cx >= g.width || cz < 0 || cz >= g.height {
		return
	}
	cell := g.cells[cx][cz]
	delete(cell.Units, unit.ID)
}

// RemoveUnitByID 按单位 ID 在网格中查找并移除；未找到则返回 false。
func (g *Grid) RemoveUnitByID(unitID uint64) bool {
	if g == nil || unitID == 0 {
		return false
	}
	for cx := 0; cx < g.width; cx++ {
		for cz := 0; cz < g.height; cz++ {
			cell := g.cells[cx][cz]
			if u, ok := cell.Units[unitID]; ok {
				g.RemoveUnit(u)
				return true
			}
		}
	}
	return false
}

// UpdateUnit 更新坐标；若跨格则迁移。同格内只改坐标。
func (g *Grid) UpdateUnit(unit *Unit, newCX, newCZ int) {
	if unit == nil {
		return
	}

	if newCX == unit.cellX && newCZ == unit.cellZ {
		return
	}

	oldCX, oldCZ := unit.cellX, unit.cellZ
	if oldCX >= 0 && oldCX < g.width && oldCZ >= 0 && oldCZ < g.height {
		oldCell := g.cells[oldCX][oldCZ]

		delete(oldCell.Units, unit.ID)

	}

	unit.cellX, unit.cellZ = newCX, newCZ

	newCell := g.cells[newCX][newCZ]
	newCell.Units[unit.ID] = unit

}

// GetNearbyUnits 返回与中心格 (centerCX, centerCZ) 在格子坐标系下欧氏距离不超过 radius 格的单位。
// centerCX/centerCZ 为格子索引（非世界坐标）；radius 为格子半径（与 [Unit.cellX]/cellZ 同一单位）。
func (g *Grid) GetNearbyUnits(centerCX, centerCZ, radius int) []*Unit {
	if g == nil {
		return nil
	}
	if radius < 0 {
		radius = 0
	}
	r2 := radius * radius

	minCX := centerCX - radius
	maxCX := centerCX + radius
	minCZ := centerCZ - radius
	maxCZ := centerCZ + radius

	if maxCX < 0 || minCX >= g.width || maxCZ < 0 || minCZ >= g.height {
		return nil
	}
	if minCX < 0 {
		minCX = 0
	}
	if maxCX >= g.width {
		maxCX = g.width - 1
	}
	if minCZ < 0 {
		minCZ = 0
	}
	if maxCZ >= g.height {
		maxCZ = g.height - 1
	}

	out := make([]*Unit, 0, 16)
	for cx := minCX; cx <= maxCX; cx++ {
		for cz := minCZ; cz <= maxCZ; cz++ {
			cell := g.cells[cx][cz]
			for _, u := range cell.Units {
				dx := u.cellX - centerCX
				dz := u.cellZ - centerCZ
				if dx*dx+dz*dz <= r2 {
					out = append(out, u)
				}
			}
		}
	}
	return out
}

// ForEachCellAsc 正序遍历格子：cellX ∈ [0,width)，cellZ ∈ [0,height)，先增 X 再增 Z（列优先可理解为先扫一行 Z）。
func (g *Grid) ForEachCellAsc(fn func(cellX, cellZ int, cell *GridCell)) {
	if g == nil || fn == nil {
		return
	}
	for cx := 0; cx < g.width; cx++ {
		for cz := 0; cz < g.height; cz++ {
			fn(cx, cz, g.cells[cx][cz])
		}
	}
}

// ForEachCellDesc 倒序遍历格子：从 (width-1, height-1) 递减至 (0,0)，与 [ForEachCellAsc] 顺序相反。
func (g *Grid) ForEachCellDesc(fn func(cellX, cellZ int, cell *GridCell)) {
	if g == nil || fn == nil {
		return
	}
	for cx := g.width - 1; cx >= 0; cx-- {
		for cz := g.height - 1; cz >= 0; cz-- {
			fn(cx, cz, g.cells[cx][cz])
		}
	}
}

// FirstFreeCellAsc 按与 [ForEachCellAsc] 相同顺序找到第一个空闲格；若全部占用则 ok=false。
func (g *Grid) FirstFreeCellAsc() (cellX, cellZ int, ok bool) {
	if g == nil {
		return 0, 0, false
	}
	for cx := 0; cx < g.width; cx++ {
		for cz := 0; cz < g.height; cz++ {
			if g.cells[cx][cz].Empty() {
				return cx, cz, true
			}
		}
	}
	return 0, 0, false
}

// FirstFreeCellDesc 按与 [ForEachCellDesc] 相同顺序找到第一个空闲格；若全部占用则 ok=false。
func (g *Grid) FirstFreeCellDesc() (cellX, cellZ int, ok bool) {
	if g == nil {
		return 0, 0, false
	}
	for cx := g.width - 1; cx >= 0; cx-- {
		for cz := g.height - 1; cz >= 0; cz-- {
			if g.cells[cx][cz].Empty() {
				return cx, cz, true
			}
		}
	}
	return 0, 0, false
}

// PickFreeCell 返回第一个空闲格：useAsc 为 true 时同 [FirstFreeCellAsc]，否则同 [FirstFreeCellDesc]（常用于红/蓝两侧落点策略）。
func (g *Grid) PickFreeCell(useAsc bool) (cellX, cellZ int, ok bool) {
	if g == nil {
		return 0, 0, false
	}
	if useAsc {
		return g.FirstFreeCellAsc()
	}
	return g.FirstFreeCellDesc()
}
