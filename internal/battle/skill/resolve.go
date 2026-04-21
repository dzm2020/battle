package skill

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// ValidCastTargets 意图阶段校验：非法 scope、缺失主目标（单体）、或解析不到任何目标时失败。
func ValidCastTargets(w *ecs.World, caster, primary ecs.Entity, sk SkillConfig) bool {
	if sk.Scope == TargetScopeIllegal {
		return false
	}
	if !w.EntityExists(caster) {
		return false
	}
	switch sk.Scope {
	case TargetScopeSelf:
		return true
	case TargetScopeSingle:
		if primary == 0 || !w.EntityExists(primary) {
			return false
		}
		cs, ok := casterSide(w, caster)
		if !ok {
			return false
		}
		return campMatch(w, caster, primary, sk, cs)
	case TargetScopeMulti, TargetScopeCircle, TargetScopeCone, TargetScopeFullScreen, TargetScopeChain, TargetScopeRandom:
		targets := ResolveTargets(w, caster, primary, sk)
		return len(targets) > 0
	default:
		return false
	}
}

// ResolveTargets 解析技能命中实体列表（scope × camp；PickRule / MaxTargets / AOE 简化为 MaxTargets 截断）。
func ResolveTargets(w *ecs.World, caster, primary ecs.Entity, sk SkillConfig) []ecs.Entity {
	cs, ok := casterSide(w, caster)
	if !ok {
		return nil
	}
	switch sk.Scope {
	case TargetScopeSelf:
		return []ecs.Entity{caster}
	case TargetScopeSingle:
		if primary == 0 || !w.EntityExists(primary) {
			return nil
		}
		if !campMatch(w, caster, primary, sk, cs) {
			return nil
		}
		return []ecs.Entity{primary}
	case TargetScopeMulti:
		out := collectByCamp(w, caster, sk, cs)
		return trimMaxTargets(out, sk.MaxTargets)
	default:
		return nil
	}
}

func casterSide(w *ecs.World, caster ecs.Entity) (uint8, bool) {
	t, ok := w.GetComponent(caster, &component.Team{})
	if !ok {
		return 0, false
	}
	return t.(*component.Team).Side, true
}

func campMatch(w *ecs.World, caster, target ecs.Entity, sk SkillConfig, casterSideVal uint8) bool {
	t, ok := w.GetComponent(target, &component.Team{})
	if !ok {
		return false
	}
	ts := t.(*component.Team).Side
	switch sk.Camp {
	case CampEnemy:
		return ts != casterSideVal
	case CampAllyIncludeSelf:
		return ts == casterSideVal
	case CampAllyExcludeSelf:
		return ts == casterSideVal && target != caster
	case CampEveryone:
		return true
	case CampSpecificSide:
		return ts == sk.CampSide
	default:
		return false
	}
}

func collectByCamp(w *ecs.World, caster ecs.Entity, sk SkillConfig, casterSideVal uint8) []ecs.Entity {
	q := ecs.NewQuery2[*component.Team, *component.Health](w)
	var out []ecs.Entity
	q.ForEach(func(e ecs.Entity, tm *component.Team, _ *component.Health) {
		if !entityPassesCamp(caster, e, sk, casterSideVal, tm.Side) {
			return
		}
		out = append(out, e)
	})
	return out
}

func entityPassesCamp(caster, e ecs.Entity, sk SkillConfig, casterSideVal, targetSide uint8) bool {
	switch sk.Camp {
	case CampEnemy:
		return targetSide != casterSideVal
	case CampAllyIncludeSelf:
		return targetSide == casterSideVal
	case CampAllyExcludeSelf:
		return targetSide == casterSideVal && e != caster
	case CampEveryone:
		return true
	case CampSpecificSide:
		return targetSide == sk.CampSide
	default:
		return false
	}
}

func trimMaxTargets(targets []ecs.Entity, max int) []ecs.Entity {
	if max <= 0 || len(targets) <= max {
		return targets
	}
	return targets[:max]
}
