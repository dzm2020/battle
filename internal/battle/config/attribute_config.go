package config

type AttributeType = string

const (
	AttrHp            AttributeType = "hp"             // 生命值（由 HealthSystem 处理，不参与 ResourceSystem 恢复）
	AttrMana          AttributeType = "mana"           // 法力值
	AttrRage          AttributeType = "rage"           // 怒气
	AttrEnergy        AttributeType = "energy"         // 能量
	AttrAttackDamage  AttributeType = "attack_damage"  // 攻击力
	AttrAbilityPower  AttributeType = "ability_power"  // 法术强度
	AttrArmor         AttributeType = "armor"          // 护甲
	AttrMagicResist   AttributeType = "magic_resist"   // 魔法抗性
	AttrAttackSpeed   AttributeType = "attack_speed"   // 攻击速度
	AttrAttackRange   AttributeType = "attack_range"   // 攻击距离
	AttrCritRate      AttributeType = "crit_rate"      // 暴击率
	AttrCritDamage    AttributeType = "crit_damage"    // 暴击伤害
	AttrHitPermille   AttributeType = "hit_permille"   // 命中率
	AttrDodgePermille AttributeType = "dodge_permille" // 闪避率
)

// IsCombatResource 是否为战斗资源（法力/怒气/能量等，由 [system.ResourceSystem] 负责消耗与自然恢复）。
func IsCombatResource(typ AttributeType) bool {
	switch typ {
	case AttrMana, AttrRage, AttrEnergy:
		return true
	default:
		return false
	}
}

// AttributeConfig
// @Description: 属性配置表
type AttributeConfig struct {
	ID        int32         // 配置项唯一ID
	Type      AttributeType // 属性类型
	InitValue int32         // 属性初始值
	MaxValue  int32         // 属性最大值
}
