package system

import (
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/resource"
	"battle/internal/battle/system/attrs"
)

func TestBattleEndSystem_OnDeath_SetsPhaseSettled(t *testing.T) {
	w := newCombatWorld(t)
	ecs.AddResource(w, &resource.RoomPhase{Phase: resource.PhaseFighting})

	red := w.CreateEntity()
	w.AddComponent(red, &component.Team{Side: component.SideTypeRed})
	aRed := ecs.EnsureGetComponent[*component.Attributes](w, red)
	attrs.SetRange(aRed, config.AttrHp, 100, 100)

	blue := w.CreateEntity()
	w.AddComponent(blue, &component.Team{Side: component.SideTypeBlue})
	aBlue := ecs.EnsureGetComponent[*component.Attributes](w, blue)
	attrs.SetRange(aBlue, config.AttrHp, 100, 100)

	endSys := &BattleEndSystem{}
	endSys.Initialize(w)
	deathSys := &DeathSystem{}
	deathSys.Initialize(w)

	attrs.SetCurrent(aBlue, config.AttrHp, 0)
	deathSys.Update(0)

	phase := ecs.GetResource[resource.RoomPhase](w)
	if phase == nil || phase.Phase != resource.PhaseSettled {
		t.Fatalf("phase: got %v want Settled", phase)
	}
}

func TestBattleEndSystem_OnDeath_AllDead_Draw(t *testing.T) {
	w := newCombatWorld(t)
	ecs.AddResource(w, &resource.RoomPhase{Phase: resource.PhaseFighting})

	red := w.CreateEntity()
	w.AddComponent(red, &component.Team{Side: component.SideTypeRed})
	aRed := ecs.EnsureGetComponent[*component.Attributes](w, red)
	attrs.SetRange(aRed, config.AttrHp, 0, 100)

	blue := w.CreateEntity()
	w.AddComponent(blue, &component.Team{Side: component.SideTypeBlue})
	aBlue := ecs.EnsureGetComponent[*component.Attributes](w, blue)
	attrs.SetRange(aBlue, config.AttrHp, 0, 100)

	endSys := &BattleEndSystem{}
	endSys.Initialize(w)
	deathSys := &DeathSystem{}
	deathSys.Initialize(w)

	deathSys.Update(0)

	phase := ecs.GetResource[resource.RoomPhase](w)
	if phase == nil || phase.Phase != resource.PhaseSettled {
		t.Fatalf("phase: got %v want Settled", phase)
	}
}

func TestBattleEndSystem_NoEndWithoutDeath(t *testing.T) {
	w := newCombatWorld(t)
	ecs.AddResource(w, &resource.RoomPhase{Phase: resource.PhaseFighting})

	red := w.CreateEntity()
	w.AddComponent(red, &component.Team{Side: component.SideTypeRed})
	attrs.SetRange(ecs.EnsureGetComponent[*component.Attributes](w, red), config.AttrHp, 100, 100)

	blue := w.CreateEntity()
	w.AddComponent(blue, &component.Team{Side: component.SideTypeBlue})
	attrs.SetRange(ecs.EnsureGetComponent[*component.Attributes](w, blue), config.AttrHp, 100, 100)

	endSys := &BattleEndSystem{}
	endSys.Initialize(w)

	phase := ecs.GetResource[resource.RoomPhase](w)
	if phase.Phase != resource.PhaseFighting {
		t.Fatalf("without death event: got phase %v want Fighting", phase.Phase)
	}
}
