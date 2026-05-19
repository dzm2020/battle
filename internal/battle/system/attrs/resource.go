package attrs

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// DefaultManaRegenPerFrame 实体有法力上限但未挂 [component.ResourceRegen] 时的默认每帧恢复。
const DefaultManaRegenPerFrame = 2

// CanAfford 当前基础资源是否足够支付（读 [component.Attributes]，不读 FinalAttributes）。
func CanAfford(attr *component.Attributes, typ config.AttributeType, amount int) bool {
	if amount <= 0 || typ == "" {
		return true
	}
	return Current(attr, typ) >= amount
}

// EnqueueConsume 将消耗请求写入 [component.ResourceConsumeQueue]。
func EnqueueConsume(w *ecs.World, e ecs.Entity, typ config.AttributeType, amount int) {
	if w == nil || e == 0 || amount <= 0 || typ == "" {
		return
	}
	q := ecs.EnsureGetComponent[*component.ResourceConsumeQueue](w, e)
	q.Entries = append(q.Entries, component.ResourceConsumeEntry{Type: typ, Amount: amount})
}

// ApplyConsume 从 [component.Attributes] 扣减资源（不低于 0）。
func ApplyConsume(attr *component.Attributes, typ config.AttributeType, amount int) {
	if attr == nil || amount <= 0 || typ == "" {
		return
	}
	Sub(attr, typ, amount)
}

// ApplyRegen 按配置对战斗资源执行自然恢复（不超过 Max）。
func ApplyRegen(attr *component.Attributes, rates map[config.AttributeType]int) {
	if attr == nil || len(rates) == 0 {
		return
	}
	for typ, perFrame := range rates {
		if !config.IsCombatResource(typ) || perFrame <= 0 {
			continue
		}
		cur := Current(attr, typ)
		max := Max(attr, typ)
		if max <= 0 {
			continue
		}
		next := cur + perFrame
		if next > max {
			next = max
		}
		SetCurrent(attr, typ, next)
	}
}

// RegenRatesFor 合并实体 [component.ResourceRegen] 与缺省法力恢复。
func RegenRatesFor(attr *component.Attributes, regen *component.ResourceRegen) map[config.AttributeType]int {
	out := make(map[config.AttributeType]int)
	if regen != nil && regen.PerFrame != nil {
		for k, v := range regen.PerFrame {
			if config.IsCombatResource(k) && v > 0 {
				out[k] = v
			}
		}
	}
	if _, ok := out[config.AttrMana]; !ok && attr != nil && Max(attr, config.AttrMana) > 0 {
		out[config.AttrMana] = DefaultManaRegenPerFrame
	}
	return out
}
