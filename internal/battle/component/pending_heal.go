package component

import "battle/ecs"

// PendingHeal 待结算的直接治疗量（技能或 HoT 汇总写入）；[HealSystem] 消费后移除。
type PendingHeal struct {
	Amount int
	Source ecs.Entity
}

func (*PendingHeal) Component() {}

type PendingHealBuff struct {
	Amount int
}

func (*PendingHealBuff) Component() {}
