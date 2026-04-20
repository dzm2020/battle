// Package skill 提供技能静态模板（[SkillConfig]、[EffectConfig]）、配置表（[CatalogConfig]）以及
// JSON/YAML 加载（[LoadCatalogConfigFromJSON]、[LoadCatalogConfigFromYAML]）、目标解析与效果执行
// （[ResolveTargets]、[ExecuteEffects]）。施法由 [component.CastIntent]、[system.SkillSystem] 等驱动；
// Buff 联动通过 EffectApplyBuff 使用 [buff.DefinitionConfig]。详见 docs/skill-design.md。
package skill
