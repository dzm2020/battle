package config

import "encoding/json"

// Camp 阵营
type Camp int

const (
	CampAll     Camp = 0 // 全部
	CampFriend  Camp = 1 // 友方
	CampEnemy   Camp = 2 // 敌方
	CampNeutral Camp = 3 // 中立
)

// UnitType 兵种（支持位掩码，如需组合可用 []UnitType）
type UnitType int

const (
	UnitAll      UnitType = 0      // 全部
	UnitHero     UnitType = 1 << 0 // 英雄
	UnitMinion   UnitType = 1 << 1 // 小兵
	UnitSummoned UnitType = 1 << 2 // 召唤物
	UnitBoss     UnitType = 1 << 3 // BOSS
)

// TargetSortType 目标排序类型
type TargetSortType int

const (
	SortNone     TargetSortType = iota // 不排序
	SortHealth                         // 生命值
	SortPosition                       // 位置
)

// SortOrder 排序顺序
type SortOrder int

const (
	OrderAsc  SortOrder = 0 // 升序
	OrderDesc SortOrder = 1 // 降序
)

type TargetSelectConfig struct {
	ID          int            `json:"id"`
	IncludeSelf bool           `json:"include_self"` // 是否包含自己（独立字段，因为太常用）
	MaxCount    int            `json:"max_count"`    // -1 不限制，0 不选任何目标
	SortType    TargetSortType `json:"sort_type"`
	SortOrder   SortOrder      `json:"sort_order"`

	Filters []Filter `json:"filters,omitempty"` // 这里暂时所有条件都是and关系，后续可以扩展成用表达式来处理
}

// FilterType 筛选条件类型（与 JSON 中 `type` 字符串一一对应，反序列化时按 string 读入）。
type FilterType string

const (
	FilterCamp       FilterType = "camp"
	FilterStatusMask FilterType = "status_mask"
	FilterProperty   FilterType = "property"
)

// Filter 筛选条件表达式（支持嵌套）
type Filter struct {
	Type   FilterType      `json:"type"`
	Params json.RawMessage `json:"params,omitempty"` // 不同 type 对应的参数 CampFilter,UnitTypeFilter ....
}

// CampFilter 阵营筛选
type CampFilter struct {
	AllowedCamps []Camp `json:"allowed_camps"` // 允许的阵营列表（空表示无限制）
}

type PropertyFilter struct {
	Property string  `json:"property"` // hp, mp, attack, speed ...
	Op       string  `json:"op"`       // ">", "<", "==", ">=", "<="
	Value    float64 `json:"value"`
}

type StatusFilter struct {
	StatusMask int `json:"status_mask"` // 状态位掩码（例如眩晕、沉默，0表示不限）
}

// UnitTypeFilter 兵种位筛选；Mask==0 表示不限。当前 ECS 若无兵种组件，仅 Mask==0 时通过。
type UnitTypeFilter struct {
	Mask int `json:"mask"`
}

// LifeStateFilter 存活筛选：alive（默认）| dead | any。
type LifeStateFilter struct {
	Mode string `json:"mode"`
}
