package skill

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
)

// ResolveTargets 根据技能 [SkillConfig.Target] 从世界中解析本帧生效的实体列表（不含校验施法合法性）。
func ResolveTargets(w *ecs.World, caster ecs.Entity, primary ecs.Entity, sk SkillConfig) []ecs.Entity {
	switch sk.Target {
	case TargetSelf:
		return []ecs.Entity{caster}
	case TargetSingleEnemy:
		return []ecs.Entity{primary}
	case TargetAllEnemySides:
		return collectEnemySides(w, caster)
	default:
		return nil
	}
}

func collectEnemySides(w *ecs.World, caster ecs.Entity) []ecs.Entity {
	var casterSide uint8
	if c, ok := w.GetComponent(caster, &component.Team{}); ok {
		casterSide = c.(*component.Team).Side
	}
	var out []ecs.Entity
	q := ecs.NewQuery2[*component.Team, *component.Health](w)
	q.ForEach(func(e ecs.Entity, tm *component.Team, hp *component.Health) {
		if e == caster {
			return
		}
		if tm.Side != casterSide {
			out = append(out, e)
		}
	})
	return out
}

// ExecuteEffects 对已选目标顺序执行 [SkillConfig].Effects（需先完成资源与冷却逻辑由调用方保证）。
func ExecuteEffects(w *ecs.World, targets []ecs.Entity, sk SkillConfig, buffConfig *buff.DefinitionConfig) {
	for i := range sk.Effects {
		eff := &sk.Effects[i]
		for _, te := range targets {
			switch eff.Kind {
			case EffectDamage:
				if eff.Amount > 0 {
					component.MergePendingDamage(w, te, eff.Amount, eff.DamageType)
				}
			case EffectHeal:
				if eff.Amount <= 0 {
					continue
				}
				if h, ok := w.GetComponent(te, &component.Health{}); ok {
					hp := h.(*component.Health)
					hp.Current += eff.Amount
					if hp.Current > hp.Max {
						hp.Current = hp.Max
					}
				}
			case EffectApplyBuff:
				if eff.BuffDefID != 0 && buffConfig != nil {
					buff.ApplyBuff(w, buffConfig, te, eff.BuffDefID)
				}
			}
		}
	}
}
