package component

import "battle/ecs"

// DamageType 伤害类型，影响与哪项减伤属性结合；真实伤害忽略防御相关项。
type DamageType uint8

const (
	DamagePhysical DamageType = iota
	DamageMagical
	DamageTrue
)

type PendingDamageBuff struct {
	Amount int
}

func (*PendingDamageBuff) Component() {}

// PendingDamage 待结算的原始伤害（由技能/DoT 等写入；[DamageSystem] 消费后移除）。
// Source 为伤害来源实体，0 表示无来源或环境/DoT 不参入命中与暴击（仅走减免）。
type PendingDamage struct {
	Amount int
	Type   DamageType
	Source ecs.Entity
}

func (*PendingDamage) Component() {}

// ResolvedDamage 已按类型与护甲结算、待 [HealthSystem] 从生命上扣减的数值（[HealthSystem] 消费后移除）。
// Source 为伤害来源，供事件、日志与仇恨使用。
type ResolvedDamage struct {
	Amount int
	Source ecs.Entity
}

func (*ResolvedDamage) Component() {}

// Health 生命；单位一般与 [Attributes] 同加于战斗实体。
type Health struct {
	Current int
	Max     int
}

func (*Health) Component() {}
