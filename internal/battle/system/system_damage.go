package system

//  用来计算每帧造成的伤害，将计算结果保存到ResolvedDamage组件中

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/event"
	"battle/internal/battle/system/attrs"
	"battle/internal/battle/system/combatmath"
	"math"
	"math/rand/v2"
)

// DamageSystem 读取 [DamageQueue] 中各条 [PendingDamage]，可选执行 **命中→暴击**，再经护甲/魔抗减免，合并写入 [ResolvedDamage]。
type DamageSystem struct {
	world *ecs.World
	q     *ecs.Query2[*component.DamageQueue, *component.Attributes]
}

func (s *DamageSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery2[*component.DamageQueue, *component.Attributes](w)
}

func (s *DamageSystem) Update(dt float64) {
	s.q.ForEach(func(target ecs.Entity, dq *component.DamageQueue, _ *component.Attributes) {
		if len(dq.Entries) == 0 {
			return
		}
		entries := dq.Entries
		dq.Entries = nil

		total := 0
		for _, pd := range entries {
			if pd == nil {
				continue
			}
			raw := int(math.Floor(pd.RawDamage))
			if raw <= 0 {
				continue
			}
			final := s.calDamage(s.world, target, pd)
			total += final
		}

		if total <= 0 {
			return
		}
		resolved := ecs.EnsureGetComponent[*component.ResolvedDamage](s.world, target)
		resolved.Amount += total
	})
}

func (s *DamageSystem) calDamage(w *ecs.World, target ecs.Entity, entry *component.PendingDamage) int {
	damage := int(math.Floor(entry.RawDamage))
	if damage <= 0 {
		return 0
	}
	hit := attrs.GetAttributeFinalValue(w, entry.Source, config.AttrHitPermille)
	dodge := attrs.GetAttributeFinalValue(w, target, config.AttrDodgePermille)
	chance := hit - dodge
	if chance < 0 {
		chance = 0
	}
	// 未配置命中/闪避时视为必中（测试与缺省单位）
	if hit == 0 && dodge == 0 {
		chance = combatmath.Thousand
	}
	//  miss
	if chance < combatmath.Thousand && int(rand.UintN(combatmath.Thousand)) > chance {
		s.world.EmitEvent(ecs.Event{
			Kind: event.KindDamageMissed,
			Payload: event.Payload{
				Entity:   target,
				Attacker: entry.Source,
			},
		})
		return 0
	}
	//  暴击伤害
	crit := attrs.GetAttributeFinalValue(w, entry.Source, config.AttrCritRate)
	if int(rand.UintN(combatmath.Thousand)) < crit {
		bonus := attrs.GetAttributeFinalValue(w, entry.Source, config.AttrCritDamage)
		mult := combatmath.Thousand + bonus
		damage = int(math.Floor(float64(damage*mult) / combatmath.Thousand))
	}

	if entry.Type == component.DamageTrue {
		return damage
	}
	//  抗性减少伤害
	def := 0
	switch entry.Type {
	case component.DamagePhysical:
		def = attrs.GetAttributeFinalValue(w, target, config.AttrArmor)
	case component.DamageMagic:
		def = attrs.GetAttributeFinalValue(w, target, config.AttrMagicResist)
	default:
		def = 0
	}
	denom := combatmath.Hundred + def
	if denom < 1 {
		denom = 1
	}
	damage = int(math.Floor(float64(damage*combatmath.Hundred) / float64(denom)))
	return damage
}
