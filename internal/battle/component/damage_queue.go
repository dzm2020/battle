package component

import "battle/ecs"

type DamageType int

const (
	DamagePhysical DamageType = iota
	DamageMagic
	DamageTrue // 真实伤害无视减免
)

// PendingDamage 是目标实体上待结算的单次伤害条目
type PendingDamage struct {
	Source    ecs.Entity
	RawDamage float64
	Type      DamageType
}

// DamageQueue 组件 – 存放多个待处理伤害
type DamageQueue struct {
	Entries []*PendingDamage
}

func (dq *DamageQueue) Component() {}

func (dq *DamageQueue) Add(dmg *PendingDamage) {
	dq.Entries = append(dq.Entries, dmg)
}
