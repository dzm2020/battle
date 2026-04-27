package system

//  用来计算每帧造成的伤害，将计算结果保存到ResolvedDamage组件中

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/define"
	"math"
	"math/rand/v2"
)

// DamageSystem 读取 [PendingDamage]，可选执行 **命中→格挡→暴击**，再经护甲/魔抗减免写入 [ResolvedDamage]。
type DamageSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.PendingDamage]
}

func (s *DamageSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.PendingDamage](w)
}

func (s *DamageSystem) Update(dt float64) {
	s.q.ForEach(func(victim ecs.Entity, pd *component.PendingDamage) {
		raw := pd.Amount
		if raw <= 0 {
			s.world.RemoveComponent(victim, &component.PendingDamage{})
			return
		}

		source := pd.Source

		if source != 0 && pd.Type != component.DamageTrue {
			hitChance := combatHitChance(s.world, source, victim)
			if int(rand.UintN(define.Thousand)) >= hitChance {
				s.world.EmitEvent(ecs.Event{
					Kind:     ecs.EventDamageMissed,
					Entity:   victim,
					Attacker: source,
				})
				s.world.RemoveComponent(victim, &component.PendingDamage{})
				return
			}
			//  计算暴击加成
			raw = applyCritIfAny(s.world, source, raw)
		}
		//  计算防御减免
		defPhys, defMag := effectiveDefense(s.world, victim)
		final := MitigatedDamage(raw, pd.Type, defPhys, defMag)

		s.world.RemoveComponent(victim, &component.PendingDamage{})
		s.world.AddComponent(victim, &component.ResolvedDamage{Amount: final})
	})
}

// 获取防御值
func effectiveDefense(w *ecs.World, victim ecs.Entity) (phys int, mag int) {
	//  基础物防 魔防
	if c, ok := w.GetComponent(victim, &component.Attributes{}); ok {
		a := c.(*component.Attributes)
		phys = a.Get(config.AttrArmor)
		mag = a.Get(config.AttrMagicResist)
	}
	//  buff提供的物防 魔防
	if sm, ok := w.GetComponent(victim, &component.StatModifiers{}); ok {
		m := sm.(*component.StatModifiers)
		phys += m.ArmorDelta
		mag += m.MRDelta
	}
	return phys, mag
}

// 命中率计算
func combatHitChance(w *ecs.World, attacker, victim ecs.Entity) int {
	hit := define.DefaultHitPermille
	dodge := define.DefaultDodgePermille
	if a, ok := w.GetComponent(attacker, &component.Attributes{}); ok {
		attr := a.(*component.Attributes)
		if attr.Get(config.AttrHitPermille) > 0 {
			hit = attr.Get(config.AttrHitPermille)
		}
		if sm, ok := w.GetComponent(attacker, &component.StatModifiers{}); ok {
			hit += sm.(*component.StatModifiers).HitDeltaPermille
		}
	}
	if v, ok := w.GetComponent(victim, &component.Attributes{}); ok {
		attr := v.(*component.Attributes)
		if attr.Get(config.AttrDodgePermille) > 0 {
			dodge = attr.Get(config.AttrDodgePermille)
		}
		if sm, ok := w.GetComponent(victim, &component.StatModifiers{}); ok {
			dodge += sm.(*component.StatModifiers).DodgeDeltaPermille
		}
	}
	chance := hit - dodge
	return chance
}

// 计算暴击
func applyCritIfAny(w *ecs.World, attacker ecs.Entity, raw int) int {
	a, ok := w.GetComponent(attacker, &component.Attributes{})
	if !ok {
		return raw
	}
	attr := a.(*component.Attributes)
	crit := attr.Get(config.AttrCritRate)
	if crit <= 0 {
		crit = define.DefaultCritRatePermille
	}
	if sm, ok := w.GetComponent(attacker, &component.StatModifiers{}); ok {
		crit += sm.(*component.StatModifiers).CritRateDeltaPermille
	}
	if int(rand.UintN(define.Thousand)) >= crit {
		return raw
	}
	bonus := attr.Get(config.AttrCritDamage)
	if bonus <= 0 {
		bonus = define.DefaultCritDamageBonusPermille
	}
	if sm, ok := w.GetComponent(attacker, &component.StatModifiers{}); ok {
		bonus += sm.(*component.StatModifiers).CritDamageDeltaPermille
	}
	mult := define.Thousand + bonus
	return int(math.Floor(float64(raw*mult) / define.Thousand))
}

// MitigatedDamage 根据类型与护甲/魔抗计算最终伤害。
func MitigatedDamage(raw int, t component.DamageType, physicalArmor, magicResist int) int {
	if raw <= 0 {
		return 0
	}
	if t == component.DamageTrue {
		return raw
	}
	def := 0
	switch t {
	case component.DamagePhysical:
		def = physicalArmor
	case component.DamageMagical:
		def = magicResist
	default:
		def = 0
	}
	denom := define.Hundred + def
	if denom < 1 {
		denom = 1
	}
	out := int(math.Floor(float64(raw*define.Hundred) / float64(denom)))
	if out < 0 {
		return 0
	}
	return out
}
