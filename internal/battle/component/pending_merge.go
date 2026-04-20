package component

import "battle/ecs"

// MergePendingDamage 向实体“合并”待结算伤害，供 [system.BuffSystem] 的 DoT 与技能等写入同一条
// [PendingDamage]：同实体上该组件只应存在一份，因此多源同帧伤害需累加 Amount。若新加类型与
// 既有 Type 不一致，则将 Type 置为 [DamageTrue]，避免静默覆盖导致丢伤害（见 docs/buff-design.md）。
func MergePendingDamage(w *ecs.World, e ecs.Entity, amount int, dt DamageType) {
	if amount <= 0 {
		return
	}
	c, ok := w.GetComponent(e, &PendingDamage{})
	if !ok {
		w.AddComponent(e, &PendingDamage{Amount: amount, Type: dt})
		return
	}
	pd := c.(*PendingDamage)
	pd.Amount += amount
	if pd.Type != dt {
		pd.Type = DamageTrue
	}
}
