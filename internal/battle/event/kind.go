package event

import "battle/ecs"

// 战斗等业务自定义事件类别（勿与 [ecs] 框架内置 Kind 冲突）。
const (
	KindDamageApplied ecs.EventKind = ecs.FirstUserEventKind + iota
	KindDamageMissed
	KindHealApplied
	KindDeath
	KindBattleStart
	KindBattleEnd
)

// Payload 业务事件负载（订阅方按需读取字段）。
type Payload struct {
	IntPayload int
	Entity     ecs.Entity
	Attacker   ecs.Entity
}
