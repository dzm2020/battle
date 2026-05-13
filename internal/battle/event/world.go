package event

import "battle/ecs"

// NewCombatWorld 创建空战斗世界（须再调用 [component.RegisterCombatTypesWorld] 注册组件类型）。
func NewCombatWorld(initEntityNum int32) *ecs.World {
	return ecs.NewWorld(initEntityNum)
}
