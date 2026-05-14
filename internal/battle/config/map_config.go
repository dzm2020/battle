package config

// MapConfig 战斗地图空间边界，与 [land.NewSpatialGrid] 参数一致（XZ 平面、均匀 cell）。
type MapConfig struct {
	ID       int32   `json:"id" yaml:"id"`
	MinX     float64 `json:"min_x" yaml:"min_x"`
	MinZ     float64 `json:"min_z" yaml:"min_z"`
	MaxX     float64 `json:"max_x" yaml:"max_x"`
	MaxZ     float64 `json:"max_z" yaml:"max_z"`
	CellSize float64 `json:"cell_size" yaml:"cell_size"`
}
