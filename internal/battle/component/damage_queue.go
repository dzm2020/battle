package component

import "battle/ecs"

type DamageType int

const (
	DamagePhysical DamageType = iota
	DamageMagic
	DamageTrue // 真实伤害无视减免
)

// PendingDamage 是目标实体上待结算的单次伤害条目。
type PendingDamage struct {
	Source    ecs.Entity
	RawDamage float64
	Type      DamageType
}

// DamageQueue 存放待 [system.DamageSystem] 结算的伤害条目；追加请用 [DamageQueueAppend]。
type DamageQueue struct {
	Entries []*PendingDamage
}

func (*DamageQueue) Component() {}

// DamageQueueAppend 向队列追加一条待结算伤害。
func DamageQueueAppend(q *DamageQueue, dmg *PendingDamage) {
	if q == nil || dmg == nil {
		return
	}
	q.Entries = append(q.Entries, dmg)
}
