package system

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/skill"
)

// AddCombatSystems 注册完整战斗管线（帧顺序）：
// Buff → 冷却 → 技能吟唱 → 技能意图 → 伤害（命中/格挡/暴击/减免）→ 治疗 → 战斗日志订阅 → 仇恨订阅 → 扣血 → 死亡。
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
	w.AddSystem(NewSkillChannelSystem(skillConfig, buffConfig))
	w.AddSystem(NewSkillIntentSystem(skillConfig, buffConfig))
	w.AddSystem(&DamageSystem{})
	w.AddSystem(&HealSystem{})
	w.AddSystem(NewCombatLogSystem(512))
	w.AddSystem(&ThreatSystem{})
	w.AddSystem(&HealthSystem{})
	w.AddSystem(&DeathSystem{})
}
