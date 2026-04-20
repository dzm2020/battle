// Package buff 提供 Buff 静态定义（[DescriptorConfig]、[EffectConfig]）、配置容器（[DefinitionConfig]）、
// 施加逻辑（[ApplyBuff]）与 JSON 加载（[LoadDefinitionConfigFromJSON]）。运行时实例存放在 [component.BuffList].Buffs，
// 由 [system.BuffSystem] 每帧驱动；详见 docs/buff-design.md。
package buff
