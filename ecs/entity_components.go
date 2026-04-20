package ecs

// ========== 实体组件存储 ==========

// EntityComponents 实体的组件集合
type EntityComponents struct {
	components map[uint8]Component
	mask       uint64 // 位掩码，快速判断组件组合
}

// NewEntityComponents 创建实体组件集合
func NewEntityComponents() *EntityComponents {
	return &EntityComponents{
		components: make(map[uint8]Component, 4),
		mask:       0,
	}
}

// Add 添加组件
func (ec *EntityComponents) Add(compID uint8, comp Component) {
	if _, exists := ec.components[compID]; !exists {
		ec.components[compID] = comp
		ec.mask |= (1 << compID)
	}
}

// Remove 移除组件
func (ec *EntityComponents) Remove(compID uint8) {
	if _, exists := ec.components[compID]; exists {
		delete(ec.components, compID)
		ec.mask &^= (1 << compID)
	}
}

// Has 检查是否拥有指定组件
func (ec *EntityComponents) Has(compID uint8) bool {
	return (ec.mask & (1 << compID)) != 0
}

// Get 获取组件
func (ec *EntityComponents) Get(compID uint8) (Component, bool) {
	comp, ok := ec.components[compID]
	return comp, ok
}
