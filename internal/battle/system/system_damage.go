package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"math"
	"math/rand/v2"
)

// 千分比默认值（Attributes 字段为 0 时使用）。
const (
	defaultHitPermille         = 950
	defaultCritChancePermille    = 50
	defaultCritDamageBonusPermille = 500 // 额外 50% → 暴伤系数 (1000+500)/1000 = 1.5
	defaultDodgePermille       = 50
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
			if int(rand.UintN(1000)) >= hitChance {
				s.world.EmitEvent(ecs.Event{
					Kind:     ecs.EventDamageMissed,
					Entity:   victim,
					Attacker: source,
				})
				s.world.RemoveComponent(victim, &component.PendingDamage{})
				return
			}

			raw = applyBlockIfAny(s.world, victim, raw)
			raw = applyCritIfAny(s.world, source, raw)
		}

		defPhys, defMag := effectiveDefense(s.world, victim)
		final := MitigatedDamage(raw, pd.Type, defPhys, defMag)

		s.world.RemoveComponent(victim, &component.PendingDamage{})
		s.world.AddComponent(victim, &component.ResolvedDamage{Amount: final, Source: source})
	})
}

func effectiveDefense(w *ecs.World, victim ecs.Entity) (phys int, mag int) {
	if c, ok := w.GetComponent(victim, &component.Attributes{}); ok {
		a := c.(*component.Attributes)
		phys = a.PhysicalArmor
		mag = a.MagicResist
	}
	if sm, ok := w.GetComponent(victim, &component.StatModifiers{}); ok {
		m := sm.(*component.StatModifiers)
		phys += m.ArmorDelta
		mag += m.MRDelta
	}
	return phys, mag
}

func combatHitChance(w *ecs.World, attacker, victim ecs.Entity) int {
	hit := defaultHitPermille
	dodge := defaultDodgePermille
	if a, ok := w.GetComponent(attacker, &component.Attributes{}); ok {
		attr := a.(*component.Attributes)
		if attr.HitPermille > 0 {
			hit = attr.HitPermille
		}
		if sm, ok := w.GetComponent(attacker, &component.StatModifiers{}); ok {
			hit += sm.(*component.StatModifiers).HitDeltaPermille
		}
	}
	if v, ok := w.GetComponent(victim, &component.Attributes{}); ok {
		attr := v.(*component.Attributes)
		if attr.DodgePermille > 0 {
			dodge = attr.DodgePermille
		}
		if sm, ok := w.GetComponent(victim, &component.StatModifiers{}); ok {
			dodge += sm.(*component.StatModifiers).DodgeDeltaPermille
		}
	}
	chance := hit - dodge + 750
	if chance < 50 {
		chance = 50
	}
	if chance > 995 {
		chance = 995
	}
	return chance
}

func applyBlockIfAny(w *ecs.World, victim ecs.Entity, raw int) int {
	v, ok := w.GetComponent(victim, &component.Attributes{})
	if !ok {
		return raw
	}
	attr := v.(*component.Attributes)
	blockChance := attr.BlockChancePermille
	if blockChance <= 0 {
		return raw
	}
	if int(rand.UintN(1000)) >= blockChance {
		return raw
	}
	amt := attr.BlockAmount
	if amt <= 0 {
		return raw
	}
	out := raw - amt
	if out < 0 {
		return 0
	}
	return out
}

func applyCritIfAny(w *ecs.World, attacker ecs.Entity, raw int) int {
	a, ok := w.GetComponent(attacker, &component.Attributes{})
	if !ok {
		return raw
	}
	attr := a.(*component.Attributes)
	crit := attr.CritChancePermille
	if crit <= 0 {
		crit = defaultCritChancePermille
	}
	if sm, ok := w.GetComponent(attacker, &component.StatModifiers{}); ok {
		crit += sm.(*component.StatModifiers).CritChanceDeltaPermille
	}
	if int(rand.UintN(1000)) >= crit {
		return raw
	}
	bonus := attr.CritDamagePermille
	if bonus <= 0 {
		bonus = defaultCritDamageBonusPermille
	}
	if sm, ok := w.GetComponent(attacker, &component.StatModifiers{}); ok {
		bonus += sm.(*component.StatModifiers).CritDamageDeltaPermille
	}
	mult := 1000 + bonus
	return int(math.Floor(float64(raw*mult) / 1000))
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
	denom := 100 + def
	if denom < 1 {
		denom = 1
	}
	out := int(math.Floor(float64(raw*100) / float64(denom)))
	if out < 0 {
		return 0
	}
	return out
}
