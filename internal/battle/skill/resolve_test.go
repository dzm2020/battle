package skill

import (
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
)

func TestResolveTargets_LowestHPEnemy(t *testing.T) {
	w := ecs.NewWorld(16)
	component.RegisterCombatTypesWorld(w)

	caster := w.CreateEntity()
	a := w.CreateEntity()
	b := w.CreateEntity()
	w.AddComponent(caster, &component.Team{Side: 1})
	w.AddComponent(a, &component.Team{Side: 2})
	w.AddComponent(b, &component.Team{Side: 2})
	w.AddComponent(a, &component.Health{Current: 80, Max: 100})
	w.AddComponent(b, &component.Health{Current: 10, Max: 100})

	sk := SkillConfig{
		Scope:      TargetScopeMulti,
		Camp:       CampEnemy,
		PickRule:   PickHPCurrentAsc,
		MaxTargets: 1,
	}
	got := ResolveTargets(w, caster, 0, sk)
	if len(got) != 1 || got[0] != b {
		t.Fatalf("want lowest hp entity %v got %v", b, got)
	}
}

func TestResolveTargets_AllEnemyMaxTwoByHPPercent(t *testing.T) {
	w := ecs.NewWorld(16)
	component.RegisterCombatTypesWorld(w)

	caster := w.CreateEntity()
	low := w.CreateEntity()
	high := w.CreateEntity()
	w.AddComponent(caster, &component.Team{Side: 1})
	for _, e := range []ecs.Entity{low, high} {
		w.AddComponent(e, &component.Team{Side: 2})
	}
	w.AddComponent(low, &component.Health{Current: 10, Max: 100})  // 10%
	w.AddComponent(high, &component.Health{Current: 50, Max: 100}) // 50%

	sk := SkillConfig{
		Scope:      TargetScopeMulti,
		Camp:       CampEnemy,
		PickRule:   PickHPPercentAsc,
		MaxTargets: 1,
	}
	got := ResolveTargets(w, caster, 0, sk)
	if len(got) != 1 || got[0] != low {
		t.Fatalf("want %v got %v", low, got)
	}
}

func TestResolveTargets_ChainBuffFilterExcludesPrimary(t *testing.T) {
	w := ecs.NewWorld(16)
	component.RegisterCombatTypesWorld(w)

	caster := w.CreateEntity()
	primary := w.CreateEntity()
	other := w.CreateEntity()
	w.AddComponent(caster, &component.Team{Side: 1})
	w.AddComponent(primary, &component.Team{Side: 2})
	w.AddComponent(other, &component.Team{Side: 2})
	w.AddComponent(primary, &component.Health{Current: 50, Max: 50})
	w.AddComponent(other, &component.Health{Current: 50, Max: 50})
	w.AddComponent(primary, &component.BuffList{Buffs: []component.BuffInstance{{DefID: 7, Stacks: 1, FramesLeft: 1}}})

	sk := SkillConfig{
		Scope:            TargetScopeChain,
		Camp:             CampEnemy,
		ChainJumps:       1,
		RequireBuffDefID: 99,
	}
	if len(ResolveTargets(w, caster, primary, sk)) != 0 {
		t.Fatal("primary without required buff should yield no chain")
	}
}

// 对应 skill_record.md：群体治疗 ≈ 范围(Multi/Circle) + 友方 + 血量百分比最低 N 个。
func TestResolveTargets_Composite_GroupHealLowestPctN(t *testing.T) {
	w := ecs.NewWorld(16)
	component.RegisterCombatTypesWorld(w)

	caster := w.CreateEntity()
	low := w.CreateEntity()
	high := w.CreateEntity()
	w.AddComponent(caster, &component.Team{Side: 1})
	w.AddComponent(low, &component.Team{Side: 1})
	w.AddComponent(high, &component.Team{Side: 1})
	w.AddComponent(low, &component.Health{Current: 10, Max: 100})
	w.AddComponent(high, &component.Health{Current: 80, Max: 100})

	sk := SkillConfig{
		Scope:      TargetScopeMulti,
		Camp:       CampAllyIncludeSelf,
		PickRule:   PickHPPercentAsc,
		MaxTargets: 2,
	}
	got := ResolveTargets(w, caster, 0, sk)
	if len(got) != 2 {
		t.Fatalf("want 2 targets got %d", len(got))
	}
	if got[0] != low {
		t.Fatalf("first should be lowest pct ally")
	}
}

func TestResolveTargets_DistanceSort(t *testing.T) {
	w := ecs.NewWorld(16)
	component.RegisterCombatTypesWorld(w)

	caster := w.CreateEntity()
	near := w.CreateEntity()
	far := w.CreateEntity()
	w.AddComponent(caster, &component.Team{Side: 1})
	w.AddComponent(caster, &component.Transform2D{X: 0, Y: 0})
	for _, e := range []ecs.Entity{near, far} {
		w.AddComponent(e, &component.Team{Side: 2})
		w.AddComponent(e, &component.Health{Current: 50, Max: 50})
	}
	w.AddComponent(near, &component.Transform2D{X: 1, Y: 0})
	w.AddComponent(far, &component.Transform2D{X: 100, Y: 0})

	sk := SkillConfig{Scope: TargetScopeMulti, Camp: CampEnemy, PickRule: PickNearest, MaxTargets: 1}
	got := ResolveTargets(w, caster, 0, sk)
	if len(got) != 1 || got[0] != near {
		t.Fatalf("want nearest %v got %v", near, got)
	}
}
