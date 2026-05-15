package ecs

import "reflect"

// InsertResource 每局战斗插入一次（如开房时）
func InsertResource[T any](w *World, res T) {
	var zero T
	w.resources[reflect.TypeOf(zero)] = res
}

// GetResource 系统在 Update 里取用
func GetResource[T any](w *World) (T, bool) {
	var zero T
	v, ok := w.resources[reflect.TypeOf(zero)]
	if !ok {
		return zero, false
	}
	return v.(T), true
}
