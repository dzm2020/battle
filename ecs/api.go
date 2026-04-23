package ecs

import "reflect"

// ========== 系统接口 ==========

// System 系统接口
type System interface {
	Initialize(w *World)
	Update(dt float64)
}

// ========== 查询器 ==========

// EnsureGetComponent 若实体上已有类型 T 的组件则返回；否则新建零值实例 [World.AddComponent] 后返回。
// T 必须为指针组件类型（如 *MyComp），并实现 [Component]。
func EnsureGetComponent[T Component](w *World, e Entity) T {
	var proto T
	if c, ok := w.GetComponent(e, proto); ok {
		return c.(T)
	}
	typ := reflect.TypeOf(proto)
	if typ.Kind() != reflect.Ptr {
		panic("ecs: EnsureGetComponent[T] expects T to be a pointer component type (*Foo)")
	}
	newComp := reflect.New(typ.Elem()).Interface().(T)
	w.AddComponent(e, newComp)
	return newComp
}
