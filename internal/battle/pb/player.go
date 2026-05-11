package pb

import "battle/internal/battle/config"

type Player struct {
	ID    uint32                 `json:"id" yaml:"id"`       // 玩家ID
	Base  *PlayerBase            `json:"base" yaml:"base"`   // 玩家基础数据
	Units map[uint32]*PlayerUnit `json:"units" yaml:"units"` // 战斗单位
}

type PlayerBase struct {
	Name  string // 玩家名字
	Level uint32 // 等级
}

type PlayerUnit struct {
	ID         uint32      `json:"id" yaml:"id"`                               // 唯一标识（与 Unit 表顶层键一致）
	Stats      []Attribute `json:"stats" yaml:"stats"`                         // 基础属性值 AttributeConfig配置表ID
	Ability    []int32     `json:"ability,omitempty" yaml:"ability,omitempty"` // 技能配置（可选） 技能配置表ID
	BuffDefIDs []uint32    `json:"spawnBuffDefIds,omitempty"`                  // 初始Buff  BuffConfig配置表ID
}

type Attribute struct {
	ID        int32            // 配置项唯一ID
	Type      config.Attribute // 属性类型
	InitValue int32            // 属性初始值
	MaxValue  int32            // 属性最大值
}
