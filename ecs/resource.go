package ecs

import "reflect"

// Resource 提供对 World 资源的类型安全访问；由 [NewResource] 创建。
type Resource[T any] struct {
	world *World
	id    ResID
}

// New 在未初始化的 [Resource] 上创建 mapper（避免重复写类型参数）。
func (Resource[T]) New(world *World) Resource[T] {
	return NewResource[T](world)
}

// NewResource 创建资源 mapper（不向 World 添加资源，仅用于后续访问）。
//
// 参见 [World.Resources]。
func NewResource[T any](w *World) Resource[T] {
	return Resource[T]{
		id:    ResourceID[T](w),
		world: w,
	}
}

// Add 向 World 添加资源。
//
// 若该类型资源已存在则 panic。
func (g *Resource[T]) Add(res *T) {
	g.world.Resources().Add(g.id, res)
}

// Remove 从 World 移除资源。
//
// 若该类型资源不存在则 panic。
func (g *Resource[T]) Remove() {
	g.world.Resources().Remove(g.id)
}

// Get 获取资源指针；不存在时返回 nil。
func (g *Resource[T]) Get() *T {
	res := g.world.Resources().Get(g.id)
	if res == nil {
		return nil
	}
	return res.(*T)
}

// Has 是否已添加该类型资源。
func (g *Resource[T]) Has() bool {
	return g.world.Resources().Has(g.id)
}

// ResourceID 返回资源类型的 [ResID]；若未注册则自动注册。
func ResourceID[T any](w *World) ResID {
	return w.resourceID(reflect.TypeFor[T]())
}

// ResourceIDs 返回已注册的全部资源 ID。
func ResourceIDs(w *World) []ResID {
	intIDs := w.resources.registry.ids
	ids := make([]ResID, len(intIDs))
	for i, iid := range intIDs {
		ids[i] = ResID{id: iid}
	}
	return ids
}

// ResourceTypeID 按 reflect.Type 返回 [ResID]；若未注册则自动注册。
func ResourceTypeID(w *World, tp reflect.Type) ResID {
	return w.resourceID(tp)
}

// ResourceType 按 [ResID] 返回 reflect.Type。
func ResourceType(w *World, id ResID) (reflect.Type, bool) {
	return w.resources.registry.resourceType(id.id)
}

// GetResource 按类型获取资源指针；不存在时返回 nil。
//
// 重复访问建议缓存 [Resource] mapper。
func GetResource[T any](w *World) *T {
	res := w.resources.Get(ResourceID[T](w))
	if res == nil {
		return nil
	}
	return res.(*T)
}

// AddResource 向 World 添加资源并返回其 [ResID]。
//
// 若该类型资源已存在则 panic。
func AddResource[T any](w *World, res *T) ResID {
	id := ResourceID[T](w)
	w.resources.Add(id, res)
	return id
}
