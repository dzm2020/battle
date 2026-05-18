package system

import (
	"battle/ecs"
)

// AddCoreCombatSystems 注册通用战斗管线（不含 [BattleInitSystem]，由 [room] 在开房时单独挂载）。
func AddCoreCombatSystems(w *ecs.World) {
	w.AddSystem(&SpawnSystem{})
	w.AddSystem(&BuffSystem{})
	w.AddSystem(&CooldownSystem{})
	w.AddSystem(&CastValidationSystem{})
	w.AddSystem(&CastStateSystem{})
	w.AddSystem(&DamageSystem{})
	w.AddSystem(&HealSystem{})
	w.AddSystem(&HealthSystem{})
	w.AddSystem(&DeathSystem{})
	w.AddSystem(&BattleEndSystem{})
}

// AddPVESystems 注册 PVE 副本战斗管线；可在通用管线之上追加 PVE 专用 System。
func AddPVESystems(w *ecs.World) {
	AddCoreCombatSystems(w)
}

// AddPVPSystems 注册 PVP 副本战斗管线；可在通用管线之上追加 PVP 专用 System。
func AddPVPSystems(w *ecs.World) {
	AddCoreCombatSystems(w)
}
