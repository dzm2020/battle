package ecs

// ========== 世界 ==========

// World ECS 世界
type World struct {
	registry *ComponentRegistry
	entities map[Entity]*EntityComponents
	systems  []System
}

// NewWorld 创建世界
func NewWorld(initEntityNum int32) *World {
	return &World{
		registry: NewComponentRegistry(),
		entities: make(map[Entity]*EntityComponents, initEntityNum),
		systems:  make([]System, 0, 16),
	}
}

// Registry 获取组件注册表
func (w *World) Registry() *ComponentRegistry {
	return w.registry
}

// CreateEntity 创建实体
func (w *World) CreateEntity() Entity {
	e := NewEntity()
	w.entities[e] = NewEntityComponents()
	return e
}

// RemoveEntity 移除实体
func (w *World) RemoveEntity(e Entity) {
	delete(w.entities, e)
}

// AddComponent 为实体添加组件
func (w *World) AddComponent(e Entity, comp Component) {
	compID, ok := w.registry.ID(comp)
	if !ok {
		compID = w.registry.Register(comp)
	}
	if ec, exists := w.entities[e]; exists {
		ec.Add(compID, comp)
	}
}

// RemoveComponent 从实体移除组件
func (w *World) RemoveComponent(e Entity, comp Component) {
	if compID, ok := w.registry.ID(comp); ok {
		if ec, exists := w.entities[e]; exists {
			ec.Remove(compID)
		}
	}
}

// GetComponent 获取实体的组件
func (w *World) GetComponent(e Entity, comp Component) (Component, bool) {
	if compID, ok := w.registry.ID(comp); ok {
		if ec, exists := w.entities[e]; exists {
			return ec.Get(compID)
		}
	}
	return nil, false
}

// AddSystem 添加系统
func (w *World) AddSystem(s System) {
	s.Initialize(w)
	w.systems = append(w.systems, s)
}

// Update 更新所有系统
func (w *World) Update(dt float64) {
	for _, s := range w.systems {
		s.Update(dt)
	}
}
