package system

import (
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/system/attrs"
)

func TestAttributeSystem_RecomputesFinalFromBaseAndBuff(t *testing.T) {
	w := newCombatWorld(t)
	e := spawnCombatEntity(w, 100, 50)
	base, _ := w.GetComponent(e, &component.Attributes{})
	attrs.SetRange(base.(*component.Attributes), config.AttrArmor, 10, 10)

	mods := ecs.EnsureGetComponent[*component.BuffStatModifiers](w, e)
	mods.Modifiers = map[config.AttributeType]int32{config.AttrArmor: 5}

	sys := &AttributeSystem{}
	sys.Initialize(w)
	sys.Update(0)

	fa, ok := w.GetComponent(e, &component.FinalAttributes{})
	if !ok {
		t.Fatal("expected FinalAttributes")
	}
	if got := fa.(*component.FinalAttributes).Values[config.AttrArmor]; got != 15 {
		t.Fatalf("final armor: got %d want 15", got)
	}
	if got := attrs.Final(w, e, config.AttrArmor); got != 15 {
		t.Fatalf("attrs.Final: got %d want 15", got)
	}
}

func TestAttributeSystem_ClearsBuffOnlyKeysWhenModifiersRemoved(t *testing.T) {
	w := newCombatWorld(t)
	e := spawnCombatEntity(w, 100, 50)

	sys := &AttributeSystem{}
	sys.Initialize(w)

	mods := ecs.EnsureGetComponent[*component.BuffStatModifiers](w, e)
	mods.Modifiers = map[config.AttributeType]int32{config.AttrCritRate: 100}
	sys.Update(0)
	if attrs.Final(w, e, config.AttrCritRate) != 100 {
		t.Fatalf("with buff: got %d", attrs.Final(w, e, config.AttrCritRate))
	}

	w.RemoveComponent(e, &component.BuffStatModifiers{})
	sys.Update(0)
	if attrs.Final(w, e, config.AttrCritRate) != 0 {
		t.Fatalf("after remove buff mods: got %d want 0", attrs.Final(w, e, config.AttrCritRate))
	}
}
