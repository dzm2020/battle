package component

import "battle/internal/battle/config"

// Attribute 单条属性当前值与上限。
type Attribute struct {
	Current int
	Max     int
}

// Attributes 实体属性表（纯数据）。运行时读写见 [battle/internal/battle/system/attrs]。
type Attributes struct {
	Base map[config.AttributeType]*Attribute
}

func (*Attributes) Component() {}
