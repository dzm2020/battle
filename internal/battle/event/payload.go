package event

import "battle/ecs"

// Payload 战斗等业务层事件载荷；与 [ecs.EventKind] 组合为完整事件，放入 [ecs.Event.Payload]。
type Payload struct {
	Entity      ecs.Entity
	Attacker    ecs.Entity
	ComponentID uint8
	Component   ecs.Component
	IntPayload  int
}
