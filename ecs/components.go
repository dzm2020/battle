package ecs

import "reflect"

// Component 组件接口（所有组件必须实现）
type Component interface {
	// Component 是标记接口，仅用于类型约束
	Component()
}

// componentType 组件类型信息
type componentType struct {
	id   uint8
	name string
	typ  reflect.Type
}

// ComponentRegistry 组件类型注册表
type ComponentRegistry struct {
	types    []componentType
	typeToID map[reflect.Type]uint8
	nextID   uint8
}

// NewComponentRegistry 创建组件注册表
func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		types:    make([]componentType, 0, 32),
		typeToID: make(map[reflect.Type]uint8),
		nextID:   0,
	}
}

// Register 注册组件类型
func (r *ComponentRegistry) Register(c Component) uint8 {
	typ := reflect.TypeOf(c)
	if id, exists := r.typeToID[typ]; exists {
		return id
	}
	id := r.nextID
	r.nextID++
	r.typeToID[typ] = id
	r.types = append(r.types, componentType{
		id:   id,
		name: typ.Name(),
		typ:  typ,
	})
	return id
}

// ID 获取组件类型 ID
func (r *ComponentRegistry) ID(c Component) (uint8, bool) {
	id, ok := r.typeToID[reflect.TypeOf(c)]
	return id, ok
}
