package skill

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
	"cmp"
	"math"
	"math/rand/v2"
	"slices"
)

// ResolveTargets 根据 [SkillConfig] 解析本帧生效的目标实体 ID 列表。
// primary 为 [component.CastIntent].Target；自身、随机、最低血群体等玩法下可为 0。
func ResolveTargets(w *ecs.World, caster ecs.Entity, primary ecs.Entity, sk SkillConfig) []ecs.Entity {
	spec := effectiveTargetSpec(sk)

	switch spec.Scope {
	case 0: // JSON 缺 scope 或未写；非法配置
		return nil
	case TargetScopeSelf:
		raw := []ecs.Entity{caster}
		return filterByBuffDefs(w, raw, sk)
	case TargetScopeChain:
		return resolveChainAfterBuffFilter(w, caster, primary, sk, spec)
	case TargetScopeRandom:
		raw := collectByCamp(w, caster, spec)
		raw = filterSpatial(w, caster, primary, raw, spec)
		raw = filterByBuffDefs(w, raw, sk)
		return pickRandomEnemySubset(raw, sk)
	default:
		raw := collectCandidates(w, caster, primary, spec)
		raw = filterByBuffDefs(w, raw, sk)

		raw = sortByPickRule(w, caster, raw, spec.Pick)
		return applyMaxTargetsGeneric(raw, spec.MaxTargets)
	}
}

// collectCandidates 根据 TargetScope 产生原始候选实体列表（不含 Buff 与排序）。
func collectCandidates(w *ecs.World, caster, primary ecs.Entity, spec targetSpec) []ecs.Entity {
	switch spec.Scope {
	case TargetScopeSingle:
		if !validPrimaryForCamp(w, caster, primary, spec) {
			return nil
		}
		out := []ecs.Entity{primary}
		if spec.AOERadius > 0 {
			out = filterSpatial(w, caster, primary, out, spec)
		}
		return out

	case TargetScopeCone, TargetScopeCircle, TargetScopeLineRect:
		pool := collectByCamp(w, caster, spec)
		pool = filterSpatial(w, caster, primary, pool, spec)
		return pool

	case TargetScopeMulti:
		pool := collectByCamp(w, caster, spec)
		pool = filterSpatial(w, caster, primary, pool, spec)
		return pool

	case TargetScopeFullScreen:
		return collectFullScreenForCamp(w, caster, spec)

	default:
		return nil
	}
}

// collectByCamp 根据 targetSpec 中的 Camp / CampSide 委派到 collectByCampRelation。
func collectByCamp(w *ecs.World, caster ecs.Entity, spec targetSpec) []ecs.Entity {
	return collectByCampRelation(w, caster, spec.Camp, spec.CampSide)
}

// collectByCampRelation 按阵营关系枚举拉取世界中带 Team+Health 的实体（语义见各 collect* 函数）。
func collectByCampRelation(w *ecs.World, caster ecs.Entity, camp CampRelation, side uint8) []ecs.Entity {
	switch camp {
	case CampEnemy:
		return collectEnemySides(w, caster)
	case CampAllyIncludeSelf:
		return collectAllySides(w, caster, true)
	case CampAllyExcludeSelf:
		return collectAllySides(w, caster, false)
	case CampEveryone:
		return collectFullScreen(w, caster)
	case CampSpecificSide:
		return collectTeamSide(w, side)
	default:
		return nil
	}
}

// collectTeamSide 取出指定 Side 且含生命值的所有实体（包含施法者若在队内）。
func collectTeamSide(w *ecs.World, wantSide uint8) []ecs.Entity {
	var out []ecs.Entity
	q := ecs.NewQuery2[*component.Team, *component.Health](w)
	q.ForEach(func(e ecs.Entity, tm *component.Team, hp *component.Health) {
		if tm.Side == wantSide {
			out = append(out, e)
		}
	})
	return out
}

// collectFullScreenForCamp 在 TargetScopeFullScreen 下按 Camp 选取池子，并可套一层球范围过滤。
func collectFullScreenForCamp(w *ecs.World, caster ecs.Entity, spec targetSpec) []ecs.Entity {
	switch spec.Camp {
	case CampEveryone:
		out := collectFullScreen(w, caster)
		return filterSpatial(w, caster, 0, out, spec)
	case CampEnemy:
		out := collectEnemySides(w, caster)
		return filterSpatial(w, caster, 0, out, spec)
	case CampAllyIncludeSelf:
		out := collectAllySides(w, caster, true)
		return filterSpatial(w, caster, 0, out, spec)
	case CampAllyExcludeSelf:
		out := collectAllySides(w, caster, false)
		return filterSpatial(w, caster, 0, out, spec)
	case CampSpecificSide:
		out := collectTeamSide(w, spec.CampSide)
		return filterSpatial(w, caster, 0, out, spec)
	default:
		return nil
	}
}

// spatialAnchorEntity 圆形/范围的圆心锚点：有主目标用主目标，否则用施法者。
func spatialAnchorEntity(caster, primary ecs.Entity) ecs.Entity {
	if primary != 0 {
		return primary
	}
	return caster
}

// filterSpatial 保留与锚点距离不超过 aoeRadius 的实体；无 Transform2D 的候选被丢弃；半径≤0 时返回原列表。
func filterSpatial(w *ecs.World, caster, primary ecs.Entity, pool []ecs.Entity, spec targetSpec) []ecs.Entity {
	if spec.AOERadius <= 0 || len(pool) == 0 {
		return pool
	}
	anchor := spatialAnchorEntity(caster, primary)
	at, ok := w.GetComponent(anchor, &component.Transform2D{})
	if !ok {
		return pool
	}
	ap := at.(*component.Transform2D)
	r2 := spec.AOERadius * spec.AOERadius

	out := make([]ecs.Entity, 0, len(pool))
	for _, e := range pool {
		t, ok := w.GetComponent(e, &component.Transform2D{})
		if !ok {
			continue
		}
		tp := t.(*component.Transform2D)
		dx := ap.X - tp.X
		dy := ap.Y - tp.Y
		if dx*dx+dy*dy <= r2 {
			out = append(out, e)
		}
	}
	return out
}

// validPrimaryForCamp 校验 CastIntent.Target 是否满足当前 Camp（单体/链首跳）。
func validPrimaryForCamp(w *ecs.World, caster, primary ecs.Entity, spec targetSpec) bool {
	if primary == 0 {
		return false
	}
	switch spec.Camp {
	case CampEnemy:
		return ValidEnemyPair(w, caster, primary)
	case CampAllyIncludeSelf, CampAllyExcludeSelf:
		if spec.Camp == CampAllyExcludeSelf && primary == caster {
			return false
		}
		return ValidAllyPair(w, caster, primary)
	case CampSpecificSide:
		tm, ok := w.GetComponent(primary, &component.Team{})
		return ok && tm.(*component.Team).Side == spec.CampSide
	case CampEveryone:
		_, ok1 := w.GetComponent(primary, &component.Team{})
		_, ok2 := w.GetComponent(primary, &component.Health{})
		return ok1 && ok2
	default:
		return false
	}
}

// ValidEnemyPair 施法者与主目标为不同阵营且均含 [component.Team]。
func ValidEnemyPair(w *ecs.World, caster, target ecs.Entity) bool {
	if target == 0 || caster == target {
		return false
	}
	ca, ok1 := w.GetComponent(caster, &component.Team{})
	tg, ok2 := w.GetComponent(target, &component.Team{})
	if !ok1 || !ok2 {
		return false
	}
	return ca.(*component.Team).Side != tg.(*component.Team).Side
}

// ValidAllyPair 主目标与施法者同阵营（均含 Team）；允许 target == caster。
func ValidAllyPair(w *ecs.World, caster, target ecs.Entity) bool {
	if target == 0 {
		return false
	}
	ca, ok1 := w.GetComponent(caster, &component.Team{})
	tg, ok2 := w.GetComponent(target, &component.Team{})
	if !ok1 || !ok2 {
		return false
	}
	return ca.(*component.Team).Side == tg.(*component.Team).Side
}

// ValidCastTargets 供施法意图系统调用：scope 非法、Buff Requirement 不满足、半结构非法时返回 false。
func ValidCastTargets(w *ecs.World, caster, primary ecs.Entity, sk SkillConfig) bool {
	spec := effectiveTargetSpec(sk)

	switch spec.Scope {
	case 0:
		return false
	case TargetScopeSelf:
		return true
	case TargetScopeSingle:
		return validPrimaryForCamp(w, caster, primary, spec) && len(ResolveTargets(w, caster, primary, sk)) > 0
	case TargetScopeChain:
		switch spec.Camp {
		case CampEnemy:
			return ValidEnemyPair(w, caster, primary) && len(ResolveTargets(w, caster, primary, sk)) > 0
		case CampAllyIncludeSelf, CampAllyExcludeSelf:
			if spec.Camp == CampAllyExcludeSelf && primary == caster {
				return false
			}
			return ValidAllyPair(w, caster, primary) && len(ResolveTargets(w, caster, primary, sk)) > 0
		default:
			return len(ResolveTargets(w, caster, primary, sk)) > 0
		}
	default:
		return len(ResolveTargets(w, caster, primary, sk)) > 0
	}
}

// resolveChainAfterBuffFilter 链式：首目标须在 Buff 过滤后的同 Camp 候选内，再沿池顺序补足后续受击单位。
func resolveChainAfterBuffFilter(w *ecs.World, caster, primary ecs.Entity, sk SkillConfig, spec targetSpec) []ecs.Entity {
	switch spec.Camp {
	case CampEnemy:
		if !ValidEnemyPair(w, caster, primary) {
			return nil
		}
	case CampAllyIncludeSelf, CampAllyExcludeSelf:
		if spec.Camp == CampAllyExcludeSelf && primary == caster {
			return nil
		}
		if !ValidAllyPair(w, caster, primary) {
			return nil
		}
	default:
		if primary == 0 {
			return nil
		}
	}

	pool := collectByCamp(w, caster, spec)
	pool = filterByBuffDefs(w, pool, sk)
	if !entityInSlice(primary, pool) {
		return nil
	}

	additional := spec.ChainJumps
	if additional <= 0 {
		additional = 2
	}
	want := 1 + additional
	if spec.MaxTargets > 0 && spec.MaxTargets < want {
		want = spec.MaxTargets
	}
	seen := make(map[ecs.Entity]struct{}, want)
	seen[caster] = struct{}{}
	out := make([]ecs.Entity, 0, want)
	out = append(out, primary)
	seen[primary] = struct{}{}

	for _, e := range pool {
		if len(out) >= want {
			break
		}
		if _, ok := seen[e]; ok {
			continue
		}
		out = append(out, e)
		seen[e] = struct{}{}
	}
	out = sortByPickRule(w, caster, out, spec.Pick)
	return out
}

func entityInSlice(target ecs.Entity, ents []ecs.Entity) bool {
	for _, e := range ents {
		if e == target {
			return true
		}
	}
	return false
}

// pickRandomEnemySubset 打乱候选后取前 n 条；n 取自 MaxTargets，≤0 时默认为 3。
func pickRandomEnemySubset(pool []ecs.Entity, sk SkillConfig) []ecs.Entity {
	if len(pool) == 0 {
		return nil
	}
	n := sk.MaxTargets
	if n <= 0 {
		n = 3
	}
	p := append([]ecs.Entity(nil), pool...)
	rand.Shuffle(len(p), func(i, j int) { p[i], p[j] = p[j], p[i] })
	if len(p) < n {
		return p
	}
	return p[:n]
}

// collectEnemySides 除施法者外，异阵营且有 Health 的实体。
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

// collectAllySides 与施法者同 Side、含 Health；includeSelf 决定是否包含施法者。
func collectAllySides(w *ecs.World, caster ecs.Entity, includeSelf bool) []ecs.Entity {
	var casterSide uint8
	if c, ok := w.GetComponent(caster, &component.Team{}); ok {
		casterSide = c.(*component.Team).Side
	}
	var out []ecs.Entity
	q := ecs.NewQuery2[*component.Team, *component.Health](w)
	q.ForEach(func(e ecs.Entity, tm *component.Team, hp *component.Health) {
		if tm.Side != casterSide {
			return
		}
		if !includeSelf && e == caster {
			return
		}
		out = append(out, e)
	})
	return out
}

// collectFullScreen 施法者以外、凡带 Team+Health 的单位（不限阵营）。
func collectFullScreen(w *ecs.World, caster ecs.Entity) []ecs.Entity {
	var out []ecs.Entity
	q := ecs.NewQuery2[*component.Team, *component.Health](w)
	q.ForEach(func(e ecs.Entity, tm *component.Team, hp *component.Health) {
		if e == caster {
			return
		}
		out = append(out, e)
	})
	return out
}

// filterByBuffDefs 根据 RequireBuffDefID / ForbidBuffDefID 过滤（技能配置中的「有/无某 Buff」）。
func filterByBuffDefs(w *ecs.World, ents []ecs.Entity, sk SkillConfig) []ecs.Entity {
	req := sk.RequireBuffDefID
	forbid := sk.ForbidBuffDefID
	if req == 0 && forbid == 0 {
		return ents
	}
	out := make([]ecs.Entity, 0, len(ents))
	for _, e := range ents {
		bl, ok := w.GetComponent(e, &component.BuffList{})
		if req != 0 && !buffListHasDef(ok, bl, req) {
			continue
		}
		if forbid != 0 && buffListHasDef(ok, bl, forbid) {
			continue
		}
		out = append(out, e)
	}
	return out
}

func buffListHasDef(ok bool, bl interface{}, defID uint32) bool {
	if !ok || bl == nil {
		return false
	}
	for _, b := range bl.(*component.BuffList).Buffs {
		if b.DefID == defID {
			return true
		}
	}
	return false
}

// sortByPickRule 按选取规则排序；PickNone 或不足两个实体时不变序；平局用 Entity ID 打破。
func sortByPickRule(w *ecs.World, caster ecs.Entity, ents []ecs.Entity, rule PickRule) []ecs.Entity {
	if len(ents) <= 1 || rule == PickNone {
		return ents
	}
	switch rule {
	case PickHPCurrentAsc, PickHPPercentAsc:
		type row struct {
			ent ecs.Entity
			v   float64
		}
		rows := make([]row, 0, len(ents))
		for _, e := range ents {
			h, ok := w.GetComponent(e, &component.Health{})
			if !ok {
				continue
			}
			hp := h.(*component.Health)
			v := float64(hp.Current)
			if rule == PickHPPercentAsc {
				v = float64(hp.Current) / float64(max(hp.Max, 1))
			}
			rows = append(rows, row{ent: e, v: v})
		}
		slices.SortFunc(rows, func(a, b row) int {
			if c := cmp.Compare(a.v, b.v); c != 0 {
				return c
			}
			return cmp.Compare(uint64(a.ent), uint64(b.ent))
		})
		out := make([]ecs.Entity, len(rows))
		for i := range rows {
			out[i] = rows[i].ent
		}
		return out
	case PickNearest, PickFarthest:
		ct, cok := w.GetComponent(caster, &component.Transform2D{})
		if !cok {
			return ents
		}
		cpos := ct.(*component.Transform2D)
		type row struct {
			ent ecs.Entity
			d2  float64
		}
		rows := make([]row, 0, len(ents))
		for _, e := range ents {
			t, ok := w.GetComponent(e, &component.Transform2D{})
			if !ok {
				rows = append(rows, row{ent: e, d2: math.Inf(1)})
				continue
			}
			tp := t.(*component.Transform2D)
			dx := cpos.X - tp.X
			dy := cpos.Y - tp.Y
			rows = append(rows, row{ent: e, d2: dx*dx + dy*dy})
		}
		slices.SortFunc(rows, func(a, b row) int {
			if rule == PickFarthest {
				if c := cmp.Compare(b.d2, a.d2); c != 0 {
					return c
				}
			} else {
				if c := cmp.Compare(a.d2, b.d2); c != 0 {
					return c
				}
			}
			return cmp.Compare(uint64(a.ent), uint64(b.ent))
		})
		out := make([]ecs.Entity, len(rows))
		for i := range rows {
			out[i] = rows[i].ent
		}
		return out
	case PickAttackHighest:
		type row struct {
			ent ecs.Entity
			atk int
		}
		rows := make([]row, 0, len(ents))
		for _, e := range ents {
			atk := 0
			if a, ok := w.GetComponent(e, &component.Attributes{}); ok {
				atk = a.(*component.Attributes).PhysicalPower
			}
			rows = append(rows, row{ent: e, atk: atk})
		}
		slices.SortFunc(rows, func(a, b row) int {
			if c := cmp.Compare(b.atk, a.atk); c != 0 {
				return c
			}
			return cmp.Compare(uint64(a.ent), uint64(b.ent))
		})
		out := make([]ecs.Entity, len(rows))
		for i := range rows {
			out[i] = rows[i].ent
		}
		return out
	default:
		return ents
	}
}

// applyMaxTargetsGeneric maxN≤0 表示不截断。
func applyMaxTargetsGeneric(ents []ecs.Entity, maxN int) []ecs.Entity {
	if len(ents) == 0 || maxN <= 0 {
		return ents
	}
	if len(ents) <= maxN {
		return ents
	}
	return ents[:maxN]
}

// ExecuteEffects 对已选目标列表逐效果、逐目标执行；caster 作为伤害来源写入 [PendingDamage].Source。
func ExecuteEffects(w *ecs.World, caster ecs.Entity, targets []ecs.Entity, sk SkillConfig, buffConfig *buff.DefinitionConfig) {
	for i := range sk.Effects {
		eff := &sk.Effects[i]
		for _, te := range targets {
			switch eff.Kind {
			case EffectDamage:
				if eff.Amount > 0 {
					component.MergePendingDamage(w, te, eff.Amount, eff.DamageType, caster)
				}
			case EffectHeal:
				if eff.Amount <= 0 {
					continue
				}
				component.MergePendingHeal(w, te, eff.Amount, caster)
			case EffectApplyBuff:
				if eff.BuffDefID != 0 && buffConfig != nil {
					buff.ApplyBuff(w, buffConfig, te, eff.BuffDefID)
				}
			}
		}
	}
}
