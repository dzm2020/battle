package system

import (
	"testing"

	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/system/attrs"
	"battle/internal/battle/system/skill"
)

func TestCastValidationSystem_AcceptsValidRequest(t *testing.T) {
	mustLoadBattleConfig(t)
	w := newCombatWorld(t)
	caster := spawnCombatEntity(w, 100, 100)
	target := spawnCombatEntity(w, 100, 0)
	if err := skill.Add(w, caster, 1); err != nil {
		t.Fatal(err)
	}
	w.AddComponent(caster, &component.SkillCastRequest{SkillID: 1, TargetEntity: target})

	sys := &CastValidationSystem{}
	sys.Initialize(w)
	sys.Update(0)

	if _, ok := w.GetComponent(caster, &component.SkillCastRequest{}); ok {
		t.Fatal("request should be consumed")
	}
	state, ok := w.GetComponent(caster, &component.SkillCastState{})
	if !ok || state == nil {
		t.Fatal("expected SkillCastState")
	}
	st := state.(*component.SkillCastState)
	if !st.IsCasting || st.SkillId != 1 {
		t.Fatalf("unexpected cast state: %+v", st)
	}
}

func TestCastValidationSystem_RejectsStun(t *testing.T) {
	mustLoadBattleConfig(t)
	w := newCombatWorld(t)
	caster := spawnCombatEntity(w, 100, 100)
	target := spawnCombatEntity(w, 100, 0)
	w.AddComponent(caster, &component.BuffControlState{Flags: component.FlagStunned})
	if err := skill.Add(w, caster, 1); err != nil {
		t.Fatal(err)
	}
	w.AddComponent(caster, &component.SkillCastRequest{SkillID: 1, TargetEntity: target})

	sys := &CastValidationSystem{}
	sys.Initialize(w)
	sys.Update(0)

	if _, ok := w.GetComponent(caster, &component.SkillCastState{}); ok {
		t.Fatal("stunned caster should not enter cast state")
	}
}

func TestCastValidationSystem_RejectsInsufficientMana(t *testing.T) {
	mustLoadBattleConfig(t)
	w := newCombatWorld(t)
	caster := spawnCombatEntity(w, 100, 0)
	target := spawnCombatEntity(w, 100, 0)
	if err := skill.Add(w, caster, 3); err != nil {
		t.Fatal(err)
	}
	attr, _ := w.GetComponent(caster, &component.Attributes{})
	attrs.SetCurrent(attr.(*component.Attributes), config.AttrMana, 0)

	w.AddComponent(caster, &component.SkillCastRequest{SkillID: 3, TargetEntity: target})

	sys := &CastValidationSystem{}
	sys.Initialize(w)
	sys.Update(0)

	if _, ok := w.GetComponent(caster, &component.SkillCastState{}); ok {
		t.Fatal("expected reject when mana insufficient")
	}
}
