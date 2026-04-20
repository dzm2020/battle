// Package buff 提供 Buff 静态定义（[Descriptor]）、注册表（[DefinitionRegistry]）、
// 施加逻辑（[ApplyBuff]）与 JSON 加载。运行时实体上的多条实例存放在 [component.BuffList].Buffs
// 切片（缓冲表），由 [system.BuffSystem] 每帧驱动；详见 docs/buff-design.md。
package buff
