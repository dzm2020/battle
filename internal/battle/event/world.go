package event

import "battle/ecs"

// NewCombatWorld 创建空战斗世界；单测须再调用 [component.Register]，并按需 ecs.AddResource 注入网格与 [resource.SpawnRequestQueue]。
func NewCombatWorld(initEntityNum int32) *ecs.World {
	return ecs.NewWorld(initEntityNum)
}
