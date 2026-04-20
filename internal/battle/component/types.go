package component

// DamageType 伤害类型，影响与哪项减伤属性结合；真实伤害忽略防御相关项。
type DamageType uint8

const (
	DamagePhysical DamageType = iota
	DamageMagical
	DamageTrue
)

// PendingDamage 待结算的原始伤害（由技能/DoT 等写入；[DamageSystem] 消费后移除）。
type PendingDamage struct {
	Amount int
	Type   DamageType
}

func (*PendingDamage) Component() {}

// ResolvedDamage 已按类型与护甲结算、待 [HealthSystem] 从生命上扣减的数值（[HealthSystem] 消费后移除）。
type ResolvedDamage struct {
	Amount int
}

func (*ResolvedDamage) Component() {}

// Health 生命；单位一般与 [Attributes] 同加于战斗实体。
type Health struct {
	Current int
	Max     int
}

func (*Health) Component() {}

// Attributes 用于伤害减免；缺失时结算侧按全 0 防御处理。
type Attributes struct {
	PhysicalArmor int
	MagicResist   int
}

func (*Attributes) Component() {}
