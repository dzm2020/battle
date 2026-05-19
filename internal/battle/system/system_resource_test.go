package system

import (
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/system/attrs"
	"battle/internal/battle/system/skill"
)

func TestResourceSystem_AppliesConsumeQueue(t *testing.T) {
	w := newCombatWorld(t)
	e := spawnCombatEntity(w, 100, 100)
	attrs.EnqueueConsume(w, e, config.AttrMana, 30)

	res := &ResourceSystem{}
	res.Initialize(w)
	res.Update(0)

	if got := attrs.Current(getAttrs(t, w, e), config.AttrMana); got != 70 {
		t.Fatalf("mana after consume: got %d want 70", got)
	}
}

func TestResourceSystem_ManaRegen(t *testing.T) {
	w := newCombatWorld(t)
	e := spawnCombatEntity(w, 100, 0)
	a := getAttrs(t, w, e)
	attrs.SetRange(a, config.AttrMana, 40, 50)
	w.AddComponent(e, &component.ResourceRegen{
		PerFrame: map[config.AttributeType]int{config.AttrMana: 5},
	})

	res := &ResourceSystem{}
	res.Initialize(w)
	res.Update(0)

	if got := attrs.Current(getAttrs(t, w, e), config.AttrMana); got != 45 {
		t.Fatalf("mana after regen: got %d want 45", got)
	}
}

func TestResourceSystem_DefaultManaRegenWithoutComponent(t *testing.T) {
	w := newCombatWorld(t)
	e := spawnCombatEntity(w, 100, 0)
	a := getAttrs(t, w, e)
	attrs.SetRange(a, config.AttrMana, 40, 50)

	res := &ResourceSystem{}
	res.Initialize(w)
	res.Update(0)

	want := 40 + attrs.DefaultManaRegenPerFrame
	if got := attrs.Current(getAttrs(t, w, e), config.AttrMana); got != want {
		t.Fatalf("default regen: got %d want %d", got, want)
	}
}

func TestCastValidation_EnqueueConsumeProcessedByResourceSystem(t *testing.T) {
	mustLoadBattleConfig(t)
	w := newCombatWorld(t)
	caster := spawnCombatEntity(w, 100, 100)
	target := spawnCombatEntity(w, 100, 0)
	if err := skill.Add(w, caster, 3); err != nil {
		t.Fatal(err)
	}
	w.AddComponent(caster, &component.SkillCastRequest{SkillID: 3, TargetEntity: target})

	cv := &CastValidationSystem{}
	cv.Initialize(w)
	cv.Update(0)

	res := &ResourceSystem{}
	res.Initialize(w)
	res.Update(0)

	if got := attrs.Current(getAttrs(t, w, caster), config.AttrMana); got != 50 {
		t.Fatalf("mana after skill 3 (cost 50): got %d want 50", got)
	}
}

func getAttrs(t *testing.T, w *ecs.World, e ecs.Entity) *component.Attributes {
	t.Helper()
	c, ok := w.GetComponent(e, &component.Attributes{})
	if !ok {
		t.Fatal("no Attributes")
	}
	return c.(*component.Attributes)
}
