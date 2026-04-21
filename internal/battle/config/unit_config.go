package config

type UnitConfig struct {
	ID         string   `json:"id" yaml:"id"`                               // 唯一标识
	Name       string   `json:"name" yaml:"name"`                           // 显示名称，如 "亚托克斯"
	Stats      []int32  `json:"stats" yaml:"stats"`                         // 基础属性值 AttributeConfig配置表ID
	Ability    []int32  `json:"ability,omitempty" yaml:"ability,omitempty"` // 技能配置（可选） 技能配置表ID
	BuffDefIDs []uint32 `json:"spawnBuffDefIds,omitempty"`                  // 初始Buff  BuffConfig配置表ID
}
