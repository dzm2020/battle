package component

import "battle/internal/battle/config"

func NewAttributesFromConfigIDs(statIDs []int32) *Attributes {
	if len(statIDs) == 0 || config.Tab.AttributeConfigByID == nil {
		return &Attributes{}
	}
	m := make(map[string]int)
	for _, id := range statIDs {
		row, ok := config.Tab.AttributeConfigByID[id]
		if !ok || row.Type == "" {
			continue
		}
		m[string(row.Type)] = int(row.InitValue)
	}
	if len(m) == 0 {
		return &Attributes{}
	}
	return &Attributes{Values: m}
}

type Attributes struct {
	Values map[string]int
}

func (*Attributes) Component() {}

// Get 读取属性；nil、空 map 或未登记键返回 0。
func (a *Attributes) Get(key string) int {
	if a == nil || a.Values == nil {
		return 0
	}
	return a.Values[key]
}

// Set 写入单键（懒分配 map）。
func (a *Attributes) Set(key string, v int) {
	if a.Values == nil {
		a.Values = make(map[string]int)
	}
	a.Values[key] = v
}

// Add 在给定键上累加（常用于临时修正；持久 Buff 多用 [StatModifiers]）。
func (a *Attributes) Add(key string, delta int) {
	a.Set(key, a.Get(key)+delta)
}

// Clone 返回浅拷贝副本（拷贝 map）。
func (a *Attributes) Clone() *Attributes {
	if a == nil || a.Values == nil {
		return &Attributes{}
	}
	m := make(map[string]int, len(a.Values))
	for k, v := range a.Values {
		m[k] = v
	}
	return &Attributes{Values: m}
}
