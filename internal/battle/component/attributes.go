package component

type Attribute struct {
	Current int
	Max     int
}

type Attributes struct {
	Values map[string]*Attribute
}

func (*Attributes) Component() {}

// Get 读取 Current；nil、空 map 或未登记键返回 0。
func (a *Attributes) Get(key string) int {
	if a == nil || a.Values == nil {
		return 0
	}
	if x := a.Values[key]; x != nil {
		return x.Current
	}
	return 0
}

// Set 写入 Current。键已存在时保留 Max，并在 Max>0 时将 Current 夹到不超过 Max；
// 新键则 Current、Max 均为 v。
func (a *Attributes) Set(key string, v int) {
	if a.Values == nil {
		a.Values = make(map[string]*Attribute)
	}
	if existing := a.Values[key]; existing != nil {
		existing.Current = v
		if existing.Max > 0 && existing.Current > existing.Max {
			existing.Current = existing.Max
		}
		return
	}
	a.Values[key] = &Attribute{Current: v, Max: v}
}

// SetRange 同时写入 Current 与 Max；若 max < current 则将 max 抬到 current。
func (a *Attributes) SetRange(key string, current, max int) {
	if max < current {
		max = current
	}
	if a.Values == nil {
		a.Values = make(map[string]*Attribute)
	}
	a.Values[key] = &Attribute{Current: current, Max: max}
}

// Add 在给定键上累加 Current（常用于临时修正；持久 Buff 多用 [StatModifiers]）。
func (a *Attributes) Add(key string, delta int) {
	a.Set(key, a.Get(key)+delta)
}
