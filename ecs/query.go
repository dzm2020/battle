package ecs

// Query 查询器（泛型，类型安全）
type Query[T1 Component] struct {
	world   *World
	compID1 uint8
	results []Entity
}

// NewQuery 创建查询器
func NewQuery[T1 Component](w *World) *Query[T1] {
	var zero T1
	compID1 := ensureCompID(w, zero)
	return &Query[T1]{
		world:   w,
		compID1: compID1,
		results: make([]Entity, 0, 100),
	}
}

// Query2 两个组件的查询器
type Query2[T1, T2 Component] struct {
	world   *World
	compID1 uint8
	compID2 uint8
	results []Entity
}

// NewQuery2 创建双组件查询器
func NewQuery2[T1, T2 Component](w *World) *Query2[T1, T2] {
	var zero1 T1
	var zero2 T2
	compID1 := ensureCompID(w, zero1)
	compID2 := ensureCompID(w, zero2)
	return &Query2[T1, T2]{
		world:   w,
		compID1: compID1,
		compID2: compID2,
		results: make([]Entity, 0, 100),
	}
}

// Query3 三个组件的查询器
type Query3[T1, T2, T3 Component] struct {
	world   *World
	compID1 uint8
	compID2 uint8
	compID3 uint8
	results []Entity
}

// NewQuery3 创建三组件查询器
func NewQuery3[T1, T2, T3 Component](w *World) *Query3[T1, T2, T3] {
	var zero1 T1
	var zero2 T2
	var zero3 T3
	return &Query3[T1, T2, T3]{
		world:   w,
		compID1: ensureCompID(w, zero1),
		compID2: ensureCompID(w, zero2),
		compID3: ensureCompID(w, zero3),
		results: make([]Entity, 0, 100),
	}
}

// Query4 四个组件的查询器
type Query4[T1, T2, T3, T4 Component] struct {
	world   *World
	compID1 uint8
	compID2 uint8
	compID3 uint8
	compID4 uint8
	results []Entity
}

// NewQuery4 创建四组件查询器
func NewQuery4[T1, T2, T3, T4 Component](w *World) *Query4[T1, T2, T3, T4] {
	var zero1 T1
	var zero2 T2
	var zero3 T3
	var zero4 T4
	return &Query4[T1, T2, T3, T4]{
		world:   w,
		compID1: ensureCompID(w, zero1),
		compID2: ensureCompID(w, zero2),
		compID3: ensureCompID(w, zero3),
		compID4: ensureCompID(w, zero4),
		results: make([]Entity, 0, 100),
	}
}

// Query5 五个组件的查询器
type Query5[T1, T2, T3, T4, T5 Component] struct {
	world   *World
	compID1 uint8
	compID2 uint8
	compID3 uint8
	compID4 uint8
	compID5 uint8
	results []Entity
}

// NewQuery5 创建五组件查询器
func NewQuery5[T1, T2, T3, T4, T5 Component](w *World) *Query5[T1, T2, T3, T4, T5] {
	var zero1 T1
	var zero2 T2
	var zero3 T3
	var zero4 T4
	var zero5 T5
	return &Query5[T1, T2, T3, T4, T5]{
		world:   w,
		compID1: ensureCompID(w, zero1),
		compID2: ensureCompID(w, zero2),
		compID3: ensureCompID(w, zero3),
		compID4: ensureCompID(w, zero4),
		compID5: ensureCompID(w, zero5),
		results: make([]Entity, 0, 100),
	}
}

// Query6 六个组件的查询器
type Query6[T1, T2, T3, T4, T5, T6 Component] struct {
	world   *World
	compID1 uint8
	compID2 uint8
	compID3 uint8
	compID4 uint8
	compID5 uint8
	compID6 uint8
	results []Entity
}

// NewQuery6 创建六组件查询器
func NewQuery6[T1, T2, T3, T4, T5, T6 Component](w *World) *Query6[T1, T2, T3, T4, T5, T6] {
	var zero1 T1
	var zero2 T2
	var zero3 T3
	var zero4 T4
	var zero5 T5
	var zero6 T6
	return &Query6[T1, T2, T3, T4, T5, T6]{
		world:   w,
		compID1: ensureCompID(w, zero1),
		compID2: ensureCompID(w, zero2),
		compID3: ensureCompID(w, zero3),
		compID4: ensureCompID(w, zero4),
		compID5: ensureCompID(w, zero5),
		compID6: ensureCompID(w, zero6),
		results: make([]Entity, 0, 100),
	}
}

// Query7 七个组件的查询器
type Query7[T1, T2, T3, T4, T5, T6, T7 Component] struct {
	world   *World
	compID1 uint8
	compID2 uint8
	compID3 uint8
	compID4 uint8
	compID5 uint8
	compID6 uint8
	compID7 uint8
	results []Entity
}

// NewQuery7 创建七组件查询器
func NewQuery7[T1, T2, T3, T4, T5, T6, T7 Component](w *World) *Query7[T1, T2, T3, T4, T5, T6, T7] {
	var zero1 T1
	var zero2 T2
	var zero3 T3
	var zero4 T4
	var zero5 T5
	var zero6 T6
	var zero7 T7
	return &Query7[T1, T2, T3, T4, T5, T6, T7]{
		world:   w,
		compID1: ensureCompID(w, zero1),
		compID2: ensureCompID(w, zero2),
		compID3: ensureCompID(w, zero3),
		compID4: ensureCompID(w, zero4),
		compID5: ensureCompID(w, zero5),
		compID6: ensureCompID(w, zero6),
		compID7: ensureCompID(w, zero7),
		results: make([]Entity, 0, 100),
	}
}

// Collect 收集所有匹配的实体
func (q *Query[T1]) Collect() []Entity {
	q.results = q.results[:0]
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) {
			q.results = append(q.results, e)
		}
	}
	return q.results
}

// ForEach 遍历所有匹配的实体
func (q *Query[T1]) ForEach(fn func(Entity, T1)) {
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) {
			comp, _ := ec.Get(q.compID1)
			fn(e, comp.(T1))
		}
	}
}

// Collect 收集 Query2 的结果
func (q *Query2[T1, T2]) Collect() []Entity {
	q.results = q.results[:0]
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) {
			q.results = append(q.results, e)
		}
	}
	return q.results
}

// ForEach 遍历 Query2 的结果
func (q *Query2[T1, T2]) ForEach(fn func(Entity, T1, T2)) {
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) {
			comp1, _ := ec.Get(q.compID1)
			comp2, _ := ec.Get(q.compID2)
			fn(e, comp1.(T1), comp2.(T2))
		}
	}
}

// Collect 收集 Query3 的结果
func (q *Query3[T1, T2, T3]) Collect() []Entity {
	q.results = q.results[:0]
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) && ec.Has(q.compID3) {
			q.results = append(q.results, e)
		}
	}
	return q.results
}

// ForEach 遍历 Query3 的结果
func (q *Query3[T1, T2, T3]) ForEach(fn func(Entity, T1, T2, T3)) {
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) && ec.Has(q.compID3) {
			comp1, _ := ec.Get(q.compID1)
			comp2, _ := ec.Get(q.compID2)
			comp3, _ := ec.Get(q.compID3)
			fn(e, comp1.(T1), comp2.(T2), comp3.(T3))
		}
	}
}

// Collect 收集 Query4 的结果
func (q *Query4[T1, T2, T3, T4]) Collect() []Entity {
	q.results = q.results[:0]
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) && ec.Has(q.compID3) && ec.Has(q.compID4) {
			q.results = append(q.results, e)
		}
	}
	return q.results
}

// ForEach 遍历 Query4 的结果
func (q *Query4[T1, T2, T3, T4]) ForEach(fn func(Entity, T1, T2, T3, T4)) {
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) && ec.Has(q.compID3) && ec.Has(q.compID4) {
			comp1, _ := ec.Get(q.compID1)
			comp2, _ := ec.Get(q.compID2)
			comp3, _ := ec.Get(q.compID3)
			comp4, _ := ec.Get(q.compID4)
			fn(e, comp1.(T1), comp2.(T2), comp3.(T3), comp4.(T4))
		}
	}
}

// Collect 收集 Query5 的结果
func (q *Query5[T1, T2, T3, T4, T5]) Collect() []Entity {
	q.results = q.results[:0]
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) && ec.Has(q.compID3) && ec.Has(q.compID4) && ec.Has(q.compID5) {
			q.results = append(q.results, e)
		}
	}
	return q.results
}

// ForEach 遍历 Query5 的结果
func (q *Query5[T1, T2, T3, T4, T5]) ForEach(fn func(Entity, T1, T2, T3, T4, T5)) {
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) && ec.Has(q.compID3) && ec.Has(q.compID4) && ec.Has(q.compID5) {
			comp1, _ := ec.Get(q.compID1)
			comp2, _ := ec.Get(q.compID2)
			comp3, _ := ec.Get(q.compID3)
			comp4, _ := ec.Get(q.compID4)
			comp5, _ := ec.Get(q.compID5)
			fn(e, comp1.(T1), comp2.(T2), comp3.(T3), comp4.(T4), comp5.(T5))
		}
	}
}

// Collect 收集 Query6 的结果
func (q *Query6[T1, T2, T3, T4, T5, T6]) Collect() []Entity {
	q.results = q.results[:0]
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) && ec.Has(q.compID3) && ec.Has(q.compID4) && ec.Has(q.compID5) && ec.Has(q.compID6) {
			q.results = append(q.results, e)
		}
	}
	return q.results
}

// ForEach 遍历 Query6 的结果
func (q *Query6[T1, T2, T3, T4, T5, T6]) ForEach(fn func(Entity, T1, T2, T3, T4, T5, T6)) {
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) && ec.Has(q.compID3) && ec.Has(q.compID4) && ec.Has(q.compID5) && ec.Has(q.compID6) {
			comp1, _ := ec.Get(q.compID1)
			comp2, _ := ec.Get(q.compID2)
			comp3, _ := ec.Get(q.compID3)
			comp4, _ := ec.Get(q.compID4)
			comp5, _ := ec.Get(q.compID5)
			comp6, _ := ec.Get(q.compID6)
			fn(e, comp1.(T1), comp2.(T2), comp3.(T3), comp4.(T4), comp5.(T5), comp6.(T6))
		}
	}
}

// Collect 收集 Query7 的结果
func (q *Query7[T1, T2, T3, T4, T5, T6, T7]) Collect() []Entity {
	q.results = q.results[:0]
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) && ec.Has(q.compID3) && ec.Has(q.compID4) && ec.Has(q.compID5) && ec.Has(q.compID6) && ec.Has(q.compID7) {
			q.results = append(q.results, e)
		}
	}
	return q.results
}

// ForEach 遍历 Query7 的结果
func (q *Query7[T1, T2, T3, T4, T5, T6, T7]) ForEach(fn func(Entity, T1, T2, T3, T4, T5, T6, T7)) {
	for e, ec := range q.world.entities {
		if ec.Has(q.compID1) && ec.Has(q.compID2) && ec.Has(q.compID3) && ec.Has(q.compID4) && ec.Has(q.compID5) && ec.Has(q.compID6) && ec.Has(q.compID7) {
			comp1, _ := ec.Get(q.compID1)
			comp2, _ := ec.Get(q.compID2)
			comp3, _ := ec.Get(q.compID3)
			comp4, _ := ec.Get(q.compID4)
			comp5, _ := ec.Get(q.compID5)
			comp6, _ := ec.Get(q.compID6)
			comp7, _ := ec.Get(q.compID7)
			fn(e, comp1.(T1), comp2.(T2), comp3.(T3), comp4.(T4), comp5.(T5), comp6.(T6), comp7.(T7))
		}
	}
}

func ensureCompID[T Component](w *World, proto T) uint8 {
	if id, ok := w.Registry().ID(proto); ok {
		return id
	}
	return w.Registry().Register(proto)
}
