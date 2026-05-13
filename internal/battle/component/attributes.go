package component

import "battle/internal/battle/config"

type Attribute struct {
	Current int
	Max     int
}

type Attributes struct {
	Base map[config.AttributeType]*Attribute
}

func (*Attributes) Component() {}

func (a *Attributes) Init(key string, current, max int) {
	if max < current {
		max = current
	}
	if a.Base == nil {
		a.Base = make(map[string]*Attribute)
	}
	a.Base[key] = &Attribute{Current: current, Max: max}
}

// Get 读取 Current；nil、空 map 或未登记键返回 0。
func (a *Attributes) Get(key string) int {
	if a == nil || a.Base == nil {
		return 0
	}
	if x := a.Base[key]; x != nil {
		return x.Current
	}
	return 0
}

// Set 写入 Current。键已存在时保留 Max，并在 Max>0 时将 Current 夹到不超过 Max；
// 新键则 Current、Max 均为 v。
func (a *Attributes) Set(key string, v int) {
	if a.Base == nil {
		a.Base = make(map[string]*Attribute)
	}
	if existing := a.Base[key]; existing != nil {
		existing.Current = v
		existing.Current = min(v, existing.Current, existing.Max)
		existing.Current = max(v, existing.Current, 0)
		return
	}
	a.Base[key] = &Attribute{Current: v, Max: v}
}

// Add 在给定键上累加 Current（常用于临时修正；持久 Buff 多用 [StatModifiers]）。
func (a *Attributes) Add(key string, delta int) {
	a.Set(key, a.Get(key)+delta)
}

func (a *Attributes) Sub(key string, delta int) {
	a.Set(key, a.Get(key)-delta)
}

// SetRange 同时设置 Current 与 Max（用于单位初始化与读表生成属性）。
func (a *Attributes) SetRange(key config.AttributeType, current, max int) {
	if max < current {
		max = current
	}
	if a.Base == nil {
		a.Base = make(map[config.AttributeType]*Attribute)
	}
	a.Base[key] = &Attribute{Current: current, Max: max}
}

func (a *Attributes) GetHP() int {
	return a.Get(config.AttrHp)
}
