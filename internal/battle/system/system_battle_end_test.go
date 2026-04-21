package system

import (
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/skill"
)

func TestBattleEnd_TwoSides_Victory(t *testing.T) {
	skillConfig := skill.NewCatalogConfig()
	skillConfig.Register(skill.SkillConfig{
		ID:         1,
		Resource:   skill.ResourceNone,
		Cost:       0,
		Scope:      skill.TargetScopeSingle,
		Camp:       skill.CampEnemy,
		CastFrames: 0,
		Effects: []skill.EffectConfig{
			// 真伤：避免命中判定的随机性导致多帧未击杀、本测得不到结束事件。
			{Kind: skill.EffectDamage, Amount: 200, DamageType: component.DamageTrue},
		},
	})
	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)

	caster := w.CreateEntity()
	foe := w.CreateEntity()
	w.AddComponent(caster, &component.Team{Side: 1})
	w.AddComponent(foe, &component.Team{Side: 2})
	w.AddComponent(caster, &component.Health{Current: 100, Max: 100})
	w.AddComponent(caster, &component.SkillUser{GrantedSkillIDs: []uint32{1}})
	w.AddComponent(foe, &component.Health{Current: 30, Max: 30})
	w.AddComponent(foe, &component.Attributes{})
	w.AddComponent(caster, &component.CastIntent{SkillID: 1, Target: foe})

	AddCombatSystems(w, skillConfig)

	var winner int
	var battleEndCount int
	cancel := w.Subscribe(ecs.EventBattleEnd, func(ev ecs.Event) {
		battleEndCount++
		winner = ev.IntPayload
	})
	defer cancel()

	w.Update(0) // 建立基线 prevSides=2
	w.Update(0) // foe 阵亡 → 仅剩阵营 1

	if battleEndCount != 1 {
		t.Fatalf("battle end events want 1 got %d", battleEndCount)
	}
	if winner != 1 {
		t.Fatalf("winner side want 1 got %d", winner)
	}
}

func TestBattleEnd_SingleSide_NoPrematureVictory(t *testing.T) {
	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)

	a := w.CreateEntity()
	b := w.CreateEntity()
	w.AddComponent(a, &component.Team{Side: 1})
	w.AddComponent(b, &component.Team{Side: 1})
	w.AddComponent(a, &component.Health{Current: 10, Max: 10})
	w.AddComponent(b, &component.Health{Current: 10, Max: 10})

	AddCombatSystems(w, skill.NewCatalogConfig())

	var n int
	cancel := w.Subscribe(ecs.EventBattleEnd, func(ev ecs.Event) { n++ })
	defer cancel()

	for i := 0; i < 8; i++ {
		w.Update(0)
	}
	if n != 0 {
		t.Fatalf("single faction should not emit victory, got %d events", n)
	}
}

func TestBattleEnd_TwoSides_AllDead_Draw(t *testing.T) {
	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)

	e1 := w.CreateEntity()
	e2 := w.CreateEntity()
	w.AddComponent(e1, &component.Team{Side: 1})
	w.AddComponent(e2, &component.Team{Side: 2})
	w.AddComponent(e1, &component.Health{Current: 10, Max: 10})
	w.AddComponent(e2, &component.Health{Current: 10, Max: 10})

	AddCombatSystems(w, skill.NewCatalogConfig())

	var payload int
	var battleEndCount int
	cancel := w.Subscribe(ecs.EventBattleEnd, func(ev ecs.Event) {
		battleEndCount++
		payload = ev.IntPayload
	})
	defer cancel()

	w.Update(0) // prevSides → 2

	h1, _ := w.GetComponent(e1, &component.Health{})
	h1.(*component.Health).Current = 0
	h2, _ := w.GetComponent(e2, &component.Health{})
	h2.(*component.Health).Current = 0

	w.Update(0)

	if battleEndCount != 1 {
		t.Fatalf("battle end want 1 got %d", battleEndCount)
	}
	if payload != BattleEndPayloadDraw {
		t.Fatalf("draw payload want %d got %d", BattleEndPayloadDraw, payload)
	}
}
