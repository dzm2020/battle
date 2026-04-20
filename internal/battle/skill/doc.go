// Package skill 提供技能静态模板（[SkillConfig]、[EffectConfig]）、配置表（[CatalogConfig]）以及
// JSON/YAML 加载（[LoadCatalogConfigFromJSON]、[LoadCatalogConfigFromYAML]）、目标解析与效果执行
// （[ResolveTargets]、[ExecuteEffects]、[ValidCastTargets]）。
//
// 目标选取仅通过 [TargetScope]、[CampRelation]、[PickRule] 三维组合配置（见仓库根目录 skill_record.md）。
//
// [ResolveTargets] 管线（Self：仅 Buff 过滤；Chain/Random：分支处理；其余）：
// 按 Camp 收集候选 → 可选 [filterSpatial]（aoeRadius 与 Transform2D）→ [filterByBuffDefs]
// → [sortByPickRule] → [applyMaxTargetsGeneric]。不含资源与冷却校验；控制状态由 [action.CanAct] 等在施法系统中处理。
//
// 施法由 [component.CastIntent]、[system.SkillIntentSystem] 与 [system.SkillChannelSystem] 驱动；
// Buff 联动通过 EffectApplyBuff 使用 [buff.DefinitionConfig]。详见 docs/skill-design.md。
package skill
