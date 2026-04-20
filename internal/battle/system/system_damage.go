package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"math"
)

// DamageSystem 读取 [PendingDamage]，结合 [Attributes]（可无）计算减免，写入 [ResolvedDamage] 并移除 Pending。
type DamageSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.PendingDamage]
}

func (s *DamageSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.PendingDamage](w)
}

func (s *DamageSystem) Update(dt float64) {
	s.q.ForEach(func(e ecs.Entity, pd *component.PendingDamage) {
		var attrs *component.Attributes
		if c, ok := s.world.GetComponent(e, &component.Attributes{}); ok {
			attrs = c.(*component.Attributes)
		}
		final := s.MitigatedDamage(pd.Amount, pd.Type, attrs)
		s.world.RemoveComponent(e, &component.PendingDamage{})
		s.world.AddComponent(e, &component.ResolvedDamage{Amount: final})
	})
}

// MitigatedDamage 根据类型与护甲/魔抗计算最终伤害。
// True 不参与减免；物理使用 PhysicalArmor，魔法使用 MagicResist。
// 公式：damage * 100 / max(100 + defense, 1)，防御为负时仍可放大承受伤害。
func (s *DamageSystem) MitigatedDamage(raw int, t component.DamageType, attrs *component.Attributes) int {
	if raw <= 0 {
		return 0
	}
	if t == component.DamageTrue {
		return raw
	}
	def := 0
	switch t {
	case component.DamagePhysical:
		def = attrs.PhysicalArmor
	case component.DamageMagical:
		def = attrs.MagicResist
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
