package component

import "battle/ecs"

// PendingDamage 待结算的原始伤害（由技能/DoT 等写入；[DamageSystem] 消费后移除）。
// Source 为伤害来源实体，0 表示无来源或环境/DoT 不参入命中与暴击（仅走减免）。
type PendingDamage struct {
	Amount int
	Type   DamageType
	Source ecs.Entity
}

func (*PendingDamage) Component() {}
