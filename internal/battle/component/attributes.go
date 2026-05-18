package component

import "battle/internal/battle/config"

// Attribute 单条属性当前值与上限。
type Attribute struct {
	Current int
	Max     int
}

// Attributes 实体属性表；读写请用 [AttrCurrent]、[AttrAdd] 等包级函数，勿在组件上挂方法。
type Attributes struct {
	Base map[config.AttributeType]*Attribute
}

func (*Attributes) Component() {}
