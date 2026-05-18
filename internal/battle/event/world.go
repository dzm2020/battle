package event

import "battle/ecs"

// NewCombatWorld 创建空战斗世界（须再调用 [component.RegisterCombatTypesWorld]；
// 若测试涉及刷怪/格子，另调 [runtime.Install]）。
func NewCombatWorld(initEntityNum int32) *ecs.World {
	return ecs.NewWorld(initEntityNum)
}
