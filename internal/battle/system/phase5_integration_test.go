package system

import (
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
)

func TestThreatTopSource_AfterDamage(t *testing.T) {
	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)
	w.AddSystem(&DamageSystem{})
	w.AddSystem(&HealSystem{})
	w.AddSystem(&ThreatSystem{})
	w.AddSystem(&HealthSystem{})

	atk := w.CreateEntity()
	vic := w.CreateEntity()
	w.AddComponent(atk, &component.Team{Side: 1})
	w.AddComponent(vic, &component.Team{Side: 2})
	w.AddComponent(vic, &component.Health{Current: 100, Max: 100})
	w.AddComponent(vic, &component.ThreatBook{})
	w.AddComponent(atk, &component.Attributes{Values: map[string]int{
		component.AttrHitPermille: 1000,
		component.AttrCritRate:    0,
	}})

	component.MergePendingDamage(w, vic, 25, component.DamagePhysical, atk)
	w.Update(0)

	tb, ok := w.GetComponent(vic, &component.ThreatBook{})
	if !ok {
		t.Fatal("ThreatBook")
	}
	top := component.ThreatTopSource(tb.(*component.ThreatBook))
	if top != atk {
		t.Fatalf("want attacker as top threat got %v", top)
	}
}

func TestClearCombatEntities_RemovesTeamUnits(t *testing.T) {
	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)
	e := w.CreateEntity()
	w.AddComponent(e, &component.Team{Side: 1})
	ClearCombatEntities(w)
	if _, ok := w.GetComponent(e, &component.Team{}); ok {
		t.Fatal("entity should be removed")
	}
}
