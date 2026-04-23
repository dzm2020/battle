package component

import "battle/ecs"

// MergePendingDamage 合并待结算伤害；source 为攻击来源实体，无意图来源时传 0（仅走护甲、不参与命中暴击）。
// 同一帧多段伤害合并时：Amount 累加；Type 不一致则退化为真实伤害。
// Source 仅当全程为同一来源时保留；若与已缓存来源冲突则置 0。
func MergePendingDamage(w *ecs.World, e ecs.Entity, amount int, dt DamageType, source ecs.Entity) {
	pd := ecs.EnsureGetComponent[*PendingDamage](w, e)
	pd.Amount += amount
	if source != 0 {
		if pd.Source == 0 || pd.Source == source {
			pd.Source = source
		} else {
			pd.Source = 0
		}
	}
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
