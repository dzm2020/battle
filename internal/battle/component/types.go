package component

import "battle/ecs"

// DamageType 伤害类型，影响与哪项减伤属性结合；真实伤害忽略防御相关项。
type DamageType uint8

const (
	DamagePhysical DamageType = iota
	DamageMagical
	DamageTrue
)

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

// Attributes 基础战斗属性。部分「千分比」字段为 0 时，[DamageSystem] 使用包内默认常数（如命中 95%）。
// 仅当 [PendingDamage].Source 非 0 时，对目标才做命中/闪避/格挡/暴击判定；Source 为 0 时只走护甲减免（如部分 DoT）。
type Attributes struct {
	PhysicalPower int
	PhysicalArmor int
	MagicResist   int

	// --- 第五阶段：命中/暴击/格挡/闪避（千分比 0–1000，0=用系统默认）---
	HitPermille         int // 命中，默认 950
	CritChancePermille  int // 暴击率，默认 50
	CritDamagePermille  int // 暴伤倍率，在 1000=1.0 倍基础上额外加算，如 500=+0.5 即总 1.5 倍
	DodgePermille       int // 闪避，默认 50
	BlockChancePermille int // 格挡触发率
	BlockAmount         int // 格挡时固定格挡伤害量（在暴伤后、减免前）
}

func (*Attributes) Component() {}
