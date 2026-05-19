package component

// ResolvedDamage 已按类型与护甲结算、待 [HealthSystem] 从 [Attributes] 的 hp 扣减的数值（消费后移除）。
// Source 为伤害来源，供事件、日志与仇恨使用。
type ResolvedDamage struct {
	Amount int
}

func (*ResolvedDamage) Component() {}
