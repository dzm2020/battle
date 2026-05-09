package event

import "battle/ecs"

// NewCombatWorld 创建战斗用 ECS 世界；生命周期事件 Payload 使用本包的 [Payload]。
func NewCombatWorld(initEntityNum int32) *ecs.World {
	return ecs.NewWorld(initEntityNum, ecs.LifecycleHooks{
		EntityCreated: func(e ecs.Entity) any {
			return Payload{Entity: e}
		},
		EntityDestroyed: func(e ecs.Entity) any {
			return Payload{Entity: e}
		},
		ComponentAdded: func(e ecs.Entity, id uint8, c ecs.Component) any {
			return Payload{Entity: e, ComponentID: id, Component: c}
		},
		ComponentRemoved: func(e ecs.Entity, id uint8, c ecs.Component) any {
			return Payload{Entity: e, ComponentID: id, Component: c}
		},
	})
}
