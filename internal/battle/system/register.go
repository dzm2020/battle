package system

import (
	"battle/ecs"
)

// AddCombatSystems 注册完整战斗管线（单帧内执行顺序）：
// Buff → 技能冷却递减 → 技能施法校验（消耗/控制/CD，生成 [component.SkillCastState]）→ 技能阶段推进（前摇/生效/后摇）→
// 伤害结算 → 治疗结算 → 扣血 → 死亡移除 → 战斗结束判定。
// Buff 模板见 [config.Tab.BuffConfigConfigByID]，技能配置见 [config.Tab.SkillConfigByID]。
func AddCombatSystems(w *ecs.World) {
	w.AddSystem(&BufferSystem{})
	w.AddSystem(&SkillCooldownSystem{})
	w.AddSystem(&SkillCastValidationSystem{})
	w.AddSystem(&SkillCastStateSystem{})
	w.AddSystem(&DamageSystem{})
	w.AddSystem(&HealSystem{})
	w.AddSystem(&HealthSystem{})
	w.AddSystem(&DeathSystem{})
	w.AddSystem(&BattleEndSystem{})
}
