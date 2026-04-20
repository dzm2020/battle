package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"math"
)

// DamageSystem 读取 [PendingDamage]，将有效物甲/魔抗视为 [component.Attributes] + [component.StatModifiers]
//（后者由 [BuffSystem] 每帧刷新），再写入 [ResolvedDamage] 并移除 Pending。
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
		phys := 0
		mag := 0
		if c, ok := s.world.GetComponent(e, &component.Attributes{}); ok {
			a := c.(*component.Attributes)
			phys = a.PhysicalArmor
			mag = a.MagicResist
		}
		if sm, ok := s.world.GetComponent(e, &component.StatModifiers{}); ok {
			m := sm.(*component.StatModifiers)
			phys += m.ArmorDelta
			mag += m.MRDelta
		}
		final := MitigatedDamage(pd.Amount, pd.Type, phys, mag)
		s.world.RemoveComponent(e, &component.PendingDamage{})
		s.world.AddComponent(e, &component.ResolvedDamage{Amount: final})
	})
}

// MitigatedDamage 根据类型与护甲/魔抗计算最终伤害（已含 Buff 叠加后的有效防御）。
// True 不参与减免；物理使用 physicalArmor，魔法使用 magicResist。
// 公式：damage * 100 / max(100 + defense, 1)。
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
