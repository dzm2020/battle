package component

import "battle/internal/battle/config"

// AttrCurrent 读取 Current；nil、空 map 或未登记键返回 0。
func AttrCurrent(a *Attributes, key config.AttributeType) int {
	if a == nil || a.Base == nil {
		return 0
	}
	if x := a.Base[key]; x != nil {
		return x.Current
	}
	return 0
}

// AttrSetCurrent 写入 Current。键已存在时保留 Max，并在 Max>0 时将 Current 夹在 [0, Max]；
// 新键则 Current、Max 均为 v。
func AttrSetCurrent(a *Attributes, key config.AttributeType, v int) {
	if a == nil {
		return
	}
	if a.Base == nil {
		a.Base = make(map[config.AttributeType]*Attribute)
	}
	if existing := a.Base[key]; existing != nil {
		existing.Current = v
		if existing.Max > 0 && existing.Current > existing.Max {
			existing.Current = existing.Max
		}
		if existing.Current < 0 {
			existing.Current = 0
		}
		return
	}
	a.Base[key] = &Attribute{Current: v, Max: v}
}

// AttrAdd 在给定键上累加 Current。
func AttrAdd(a *Attributes, key config.AttributeType, delta int) {
	AttrSetCurrent(a, key, AttrCurrent(a, key)+delta)
}

// AttrSub 在给定键上减少 Current。
func AttrSub(a *Attributes, key config.AttributeType, delta int) {
	AttrAdd(a, key, -delta)
}

// AttrSetRange 同时设置 Current 与 Max（用于单位初始化与读表生成属性）。
func AttrSetRange(a *Attributes, key config.AttributeType, current, max int) {
	if max < current {
		max = current
	}
	if a == nil {
		return
	}
	if a.Base == nil {
		a.Base = make(map[config.AttributeType]*Attribute)
	}
	a.Base[key] = &Attribute{Current: current, Max: max}
}

// AttrMax 读取 Max；未登记键返回 0。
func AttrMax(a *Attributes, key config.AttributeType) int {
	if a == nil || a.Base == nil {
		return 0
	}
	if x := a.Base[key]; x != nil {
		return x.Max
	}
	return 0
}

// AttrHP 读取生命值 Current。
func AttrHP(a *Attributes) int {
	return AttrCurrent(a, config.AttrHp)
}

// AttrHPMax 读取生命值上限。
func AttrHPMax(a *Attributes) int {
	return AttrMax(a, config.AttrHp)
}
