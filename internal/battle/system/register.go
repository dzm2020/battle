package system

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/skill"
)

// AddCombatSystems 注册完整战斗管线（帧顺序）：
// Buff → 冷却递减 → 技能（吟唱/瞬发产生 PendingDamage 与 Buff）→ 伤害减免 → 扣血 → 死亡。
// buffConfig / skillConfig 为 nil 时使用空配置表；正式战斗应注入完整 [buff.DefinitionConfig] 与 [skill.CatalogConfig]。
func AddCombatSystems(w *ecs.World, buffConfig *buff.DefinitionConfig, skillConfig *skill.CatalogConfig) {
	if buffConfig == nil {
		buffConfig = buff.NewDefinitionConfig()
	}
	if skillConfig == nil {
		skillConfig = skill.NewCatalogConfig()
	}
	w.AddSystem(NewBuffSystem(buffConfig))
	w.AddSystem(NewCooldownSystem())
	w.AddSystem(NewSkillSystem(skillConfig, buffConfig))
	w.AddSystem(&DamageSystem{})
	w.AddSystem(&HealthSystem{})
	w.AddSystem(&DeathSystem{})
}
