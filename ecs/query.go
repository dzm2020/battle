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
	compID1, _ := w.Registry().ID(zero)
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
	compID1, _ := w.Registry().ID(zero1)
	compID2, _ := w.Registry().ID(zero2)
	return &Query2[T1, T2]{
		world:   w,
		compID1: compID1,
		compID2: compID2,
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
