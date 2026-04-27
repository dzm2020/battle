package component

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
