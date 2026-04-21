package config

type Attribute = string

const (
	AttrHp            Attribute = "hp"             // 生命值
	AttrMana          Attribute = "mana"           // 法力值
	AttrAttackDamage  Attribute = "attack_damage"  // 攻击力
	AttrAbilityPower  Attribute = "ability_power"  // 法术强度
	AttrArmor         Attribute = "armor"          // 护甲
	AttrMagicResist   Attribute = "magic_resist"   // 魔法抗性
	AttrAttackSpeed   Attribute = "attack_speed"   // 攻击速度
	AttrAttackRange   Attribute = "attack_range"   // 攻击距离
	AttrCritRate      Attribute = "crit_rate"      // 暴击率
	AttrCritDamage    Attribute = "crit_damage"    // 暴击伤害
	AttrHitPermille   Attribute = "hit_permille"   // 命中率
	AttrDodgePermille Attribute = "dodge_permille" // 闪避率
)

// AttributeConfig
// @Description: 属性配置表
type AttributeConfig struct {
	ID        int32     // 配置项唯一ID
	Type      Attribute // 属性类型
	InitValue int32     // 属性初始值
	MaxValue  int32     // 属性最大值
}
