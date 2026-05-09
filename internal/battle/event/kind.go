package event

import "battle/ecs"

// 战斗等业务自定义事件类别（勿与 [ecs] 框架内置 Kind 冲突）。
const (
	KindDamageApplied ecs.EventKind = ecs.FirstUserEventKind + iota
	KindDamageMissed
	KindHealApplied
	KindDeath
	KindBattleEnd
	KindAddBuffRequest
	KindRemoveBuffRequest
)
