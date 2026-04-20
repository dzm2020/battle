package component

import (
	"battle/ecs"
)

// MergePendingDamage 合并待结算伤害；source 为攻击来源实体，无意图来源时传 0（仅走护甲、不参与命中暴击）。
func MergePendingDamage(w *ecs.World, target ecs.Entity, amount int, dt DamageType, source ecs.Entity) {
	if amount <= 0 {
		return
	}
	c, ok := w.GetComponent(target, &PendingDamage{})
	if !ok {
		w.AddComponent(target, &PendingDamage{Amount: amount, Type: dt, Source: source})
		return
	}
	pd := c.(*PendingDamage)
	pd.Amount += amount
	pd.Source = source
	if pd.Type != dt {
		pd.Type = DamageTrue
	}
}

// MergePendingHeal 合并待结算治疗；由 [HealSystem] 消费。
func MergePendingHeal(w *ecs.World, target ecs.Entity, amount int, source ecs.Entity) {
	if amount <= 0 {
		return
	}
	c, ok := w.GetComponent(target, &PendingHeal{})
	if !ok {
		w.AddComponent(target, &PendingHeal{Amount: amount, Source: source})
		return
	}
	ph := c.(*PendingHeal)
	ph.Amount += amount
	if source != 0 {
		ph.Source = source
	}
}
