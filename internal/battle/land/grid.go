// Package land 提供均匀网格的空间分区（水平 XZ），用于 AOI / 邻近单位查询。
package land

import (
	"errors"
)

// ErrInvalidGridConfig 边界或 cellSize 非法。
var ErrInvalidGridConfig = errors.New("spatial: invalid bounds or cellSize")

// Unit 网格中的可移动对象；ID 可与 [battle/ecs.Entity]（uint64）对应。
// Pos 使用水平面 X、Z（深度）；若世界使用 XY 平面，可将 component.Transform2D 的 Y 映射到 Pos.Z。
type Unit struct {
	ID  uint64
	Pos struct {
		X, Z float64
	}
	cellX int
	cellZ int
}

// GridCell 单个网格单元内的单位集合。
type GridCell struct {
	Units map[uint64]*Unit
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
func (g *Grid) AddUnit(unit *Unit) {
	if unit == nil {
		return
	}
	cx, cz := g.cellIndex(unit.Pos.X, unit.Pos.Z)
	cell := g.cells[cx][cz]
	cell.Units[unit.ID] = unit

	unit.cellX, unit.cellZ = cx, cz
}

// RemoveUnit 从网格移除单位（不改变 Pos）。未加入过网格时无副作用。
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

// UpdateUnit 更新坐标；若跨格则迁移。同格内只改坐标。
func (g *Grid) UpdateUnit(unit *Unit, newX, newZ float64) {
	if unit == nil {
		return
	}
	newCX, newCZ := g.cellIndex(newX, newZ)

	if newCX == unit.cellX && newCZ == unit.cellZ {
		unit.Pos.X, unit.Pos.Z = newX, newZ
		return
	}

	oldCX, oldCZ := unit.cellX, unit.cellZ
	if oldCX >= 0 && oldCX < g.width && oldCZ >= 0 && oldCZ < g.height {
		oldCell := g.cells[oldCX][oldCZ]

		delete(oldCell.Units, unit.ID)

	}

	unit.Pos.X, unit.Pos.Z = newX, newZ
	unit.cellX, unit.cellZ = newCX, newCZ

	newCell := g.cells[newCX][newCZ]
	newCell.Units[unit.ID] = unit

}

// GetNearbyUnits 返回与 (centerX,centerZ) 平面距离不超过 radius 的单位（圆形筛选）。
func (g *Grid) GetNearbyUnits(centerX, centerZ, radius float64) []*Unit {
	if radius < 0 {
		radius = 0
	}
	r2 := radius * radius

	minX := centerX - radius
	maxX := centerX + radius
	minZ := centerZ - radius
	maxZ := centerZ + radius

	minCX, minCZ := g.cellIndex(minX, minZ)
	maxCX, maxCZ := g.cellIndex(maxX, maxZ)

	out := make([]*Unit, 0, 16)
	for cx := minCX; cx <= maxCX; cx++ {
		for cz := minCZ; cz <= maxCZ; cz++ {
			if cx < 0 || cx >= g.width || cz < 0 || cz >= g.height {
				continue
			}
			cell := g.cells[cx][cz]
			for _, u := range cell.Units {
				dx := u.Pos.X - centerX
				dz := u.Pos.Z - centerZ
				if dx*dx+dz*dz <= r2 {
					out = append(out, u)
				}
			}
		}
	}
	return out
}
