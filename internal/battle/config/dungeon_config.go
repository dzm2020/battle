package config

// DungeonConfig 副本 / 关卡配置。
type DungeonConfig struct {
	ID      int32   `json:"id" yaml:"id"`           // 唯一标识（与 Dungeon 表 JSON 顶层键对应）
	Type    int32   `json:"type" yaml:"type"`       // 副本类型
	Monster []int32 `json:"monster" yaml:"monster"` // 怪物单位模板 ID 列表（与 [Unit.json] 中顶层键一致的数值键，见 [Tables.UnitConfigByID]）
}
