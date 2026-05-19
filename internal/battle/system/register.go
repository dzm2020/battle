package system

import (
	"battle/ecs"
)

// AddCoreCombatSystems 注册通用战斗管线（不含 [BattleInitSystem]，由 [room] 在开房时单独挂载）。
func AddCoreCombatSystems(w *ecs.World) {
	w.AddSystem(&SpawnSystem{})
	w.AddSystem(&BuffSystem{})
	w.AddSystem(&AttributeSystem{})
	w.AddSystem(&CooldownSystem{})
	w.AddSystem(&CastValidationSystem{})
	w.AddSystem(&ResourceSystem{})
	w.AddSystem(&CastStateSystem{})
	w.AddSystem(&DamageSystem{})
	w.AddSystem(&HealSystem{})
	w.AddSystem(&HealthSystem{})
	w.AddSystem(&DeathSystem{})
	w.AddSystem(&BattleEndSystem{})
}

// AddPVESystems 注册 PVE 副本战斗管线；在通用管线之后挂 [PVERulesSystem]。
func AddPVESystems(w *ecs.World) {
	AddCoreCombatSystems(w)
	w.AddSystem(&PVERulesSystem{})
}

// AddPVPSystems 注册 PVP 副本战斗管线；在通用管线之后挂 [PVPRulesSystem]。
func AddPVPSystems(w *ecs.World) {
	AddCoreCombatSystems(w)
	w.AddSystem(&PVPRulesSystem{})
}

// AddSystems 注册单测用完整管线：含 [BattleInitSystem] 与 [AddCoreCombatSystems]。
// 正式开房时 [BattleInitSystem] 仅由 [room.Create] 挂载，战斗 System 由 [room_bootstrap.Installer] 在首帧注册。
func AddSystems(w *ecs.World) {
	w.AddSystem(&BattleInitSystem{})
	AddCoreCombatSystems(w)
}

// AddCombatSystems 为 [AddSystems] 的别名（历史命名，单测沿用）。
func AddCombatSystems(w *ecs.World) {
	AddSystems(w)
}
