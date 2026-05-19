package attrs

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// Recompute 根据基础属性与 Buff 修正计算最终属性表。
func Recompute(base *component.Attributes, mods *component.BuffStatModifiers) map[config.AttributeType]int {
	keys := make(map[config.AttributeType]struct{})
	if base != nil && base.Base != nil {
		for k := range base.Base {
			keys[k] = struct{}{}
		}
	}
	if mods != nil && mods.Modifiers != nil {
		for k := range mods.Modifiers {
			keys[k] = struct{}{}
		}
	}
	out := make(map[config.AttributeType]int, len(keys))
	for k := range keys {
		v := Current(base, k)
		if mods != nil && mods.Modifiers != nil {
			v += int(mods.Modifiers[k])
		}
		out[k] = v
	}
	return out
}

// Final 读取 [component.FinalAttributes] 中的最终值；未登记键返回 0。
func Final(w *ecs.World, e ecs.Entity, key config.AttributeType) int {
	if w == nil || e == 0 {
		return 0
	}
	c, ok := w.GetComponent(e, &component.FinalAttributes{})
	if !ok {
		return 0
	}
	fa := c.(*component.FinalAttributes)
	if fa.Values == nil {
		return 0
	}
	return fa.Values[key]
}

// GetAttributeFinalValue 读取最终属性（优先 [component.FinalAttributes]；无组件时回退为 [Recompute]）。
func GetAttributeFinalValue(w *ecs.World, e ecs.Entity, typ config.AttributeType) int {
	if w == nil || e == 0 {
		return 0
	}
	if c, ok := w.GetComponent(e, &component.FinalAttributes{}); ok {
		fa := c.(*component.FinalAttributes)
		if fa.Values != nil {
			return fa.Values[typ]
		}
	}
	var mods *component.BuffStatModifiers
	if m, ok := w.GetComponent(e, &component.BuffStatModifiers{}); ok {
		mods = m.(*component.BuffStatModifiers)
	}
	var base *component.Attributes
	if a, ok := w.GetComponent(e, &component.Attributes{}); ok {
		base = a.(*component.Attributes)
	}
	return Recompute(base, mods)[typ]
}
