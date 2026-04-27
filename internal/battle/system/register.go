package system

import (
	"battle/ecs"
	"battle/internal/battle/buff"
)

// AddCombatSystems 注册完整战斗管线（帧顺序）：
// Buff → 冷却 → 技能吟唱 → 技能意图 → 伤害（命中/格挡/暴击/减免）→ 治疗 → 战斗日志订阅 → 仇恨订阅 → 扣血 → 死亡 → 战斗结束判定。
// skillConfig 为 nil 时使用空技能表；Buff 模板来自全局 [config.Tab.BuffConfigConfigByID]。
func AddCombatSystems(w *ecs.World) {
	w.AddSystem(buff.NewBuffSystem())
	w.AddSystem(NewCooldownSystem())
	w.AddSystem(&DamageSystem{})
	w.AddSystem(&HealSystem{})
	//w.AddSystem(NewCombatLogSystem(512))
	w.AddSystem(&HealthSystem{})
	w.AddSystem(&DeathSystem{})
	w.AddSystem(&BattleEndSystem{})
}
