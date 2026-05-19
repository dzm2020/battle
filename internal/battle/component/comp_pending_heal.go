package component

// PendingHeal 待结算的直接治疗量（技能或 HoT 汇总写入）；[HealSystem] 消费后移除。
type PendingHeal struct {
	Amount int
}

func (*PendingHeal) Component() {}
