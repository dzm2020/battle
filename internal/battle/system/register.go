package system

import (
	"battle/ecs"
	"battle/internal/battle/buff"
)

// AddCombatSystems 注册战斗管线：必须先 [buff.BuffSystem] 再伤害结算，以便 DoT 写入的
// [component.PendingDamage] 与当帧 [component.StatModifiers] 在 [DamageSystem] 中生效。
// buffDefs 为 nil 时使用空表（仅适合测试占位）。
func AddCombatSystems(w *ecs.World, buffDefs *buff.DefinitionRegistry) {
	if buffDefs == nil {
		buffDefs = buff.NewRegistry()
	}
	w.AddSystem(NewBuffSystem(buffDefs))
	w.AddSystem(&DamageSystem{})
	w.AddSystem(&HealthSystem{})
	w.AddSystem(&DeathSystem{})
}
