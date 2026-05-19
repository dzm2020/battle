package system

import (
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
)

func TestDamageSystem_TrueDamage(t *testing.T) {
	w := newCombatWorld(t)
	target := spawnCombatEntity(w, 100, 0)
	source := spawnCombatEntity(w, 100, 0)

	dq := ecs.EnsureGetComponent[*component.DamageQueue](w, target)
	dq.Entries = []*component.PendingDamage{{
		Source:    source,
		RawDamage: 42,
		Type:      component.DamageTrue,
	}}

	sys := &DamageSystem{}
	sys.Initialize(w)
	sys.Update(0)

	resolved, ok := w.GetComponent(target, &component.ResolvedDamage{})
	if !ok {
		t.Fatal("expected ResolvedDamage")
	}
	if resolved.(*component.ResolvedDamage).Amount != 42 {
		t.Fatalf("ResolvedDamage: got %d want 42", resolved.(*component.ResolvedDamage).Amount)
	}
}

func TestDamageSystem_EmptyQueue(t *testing.T) {
	w := newCombatWorld(t)
	target := spawnCombatEntity(w, 100, 0)
	ecs.EnsureGetComponent[*component.DamageQueue](w, target)

	sys := &DamageSystem{}
	sys.Initialize(w)
	sys.Update(0)

	if _, ok := w.GetComponent(target, &component.ResolvedDamage{}); ok {
		t.Fatal("expected no ResolvedDamage for empty queue")
	}
}
