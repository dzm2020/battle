package buff

import (
	"testing"

	"battle/internal/battle/attr"
)

type mockHost struct {
	base attr.Base
	rt   *attr.Runtime
}

func (m *mockHost) AttrBase() attr.Base { return m.base }

func (m *mockHost) AttrRuntime() *attr.Runtime { return m.rt }

func (m *mockHost) IsDead() bool {
	return m.rt != nil && m.rt.CurHP <= 0
}

func TestManager_StunSlowAmpMods(t *testing.T) {
	rt := &attr.Runtime{CurHP: 100, CurMP: 50}
	h := &mockHost{base: attr.Base{Level: 1}, rt: rt}
	reg := DemoRegistry()
	mgr := NewManager(reg)

	mgr.Add(1, "demo_stun", h)
	mods := mgr.Tick(1, h)
	if !mods.Control.HasStun() {
		t.Fatalf("stun expected, got %#v", mods.Control)
	}

	mgr.Reset()
	mgr.Add(1, "demo_slow", h)
	mods = mgr.Tick(1, h)
	if mods.MoveSpeedMul >= 1 {
		t.Fatalf("slow expected, mul=%v", mods.MoveSpeedMul)
	}

	mgr.Reset()
	mgr.Add(1, "demo_amp", h)
	mgr.Add(1, "demo_amp", h)
	mods = mgr.Tick(1, h)
	if mods.OutgoingDamageMul < 1.25*1.25-1e-6 {
		t.Fatalf("amp stacks: got %v", mods.OutgoingDamageMul)
	}
}

func TestManager_DoT(t *testing.T) {
	rt := &attr.Runtime{CurHP: 200, CurMP: 50}
	h := &mockHost{base: attr.Base{Level: 1}, rt: rt}
	mgr := NewManager(DemoRegistry())
	hp0 := rt.CurHP
	mgr.Add(1, "demo_poison", h)
	for f := uint64(1); f < 61; f++ {
		_ = mgr.Tick(f, h)
	}
	if rt.CurHP != hp0 {
		t.Fatalf("before first tick hp should stay %d, got %d", hp0, rt.CurHP)
	}
	_ = mgr.Tick(61, h)
	if rt.CurHP >= hp0 {
		t.Fatalf("expected dot damage, hp %d -> %d", hp0, rt.CurHP)
	}
}

func TestManager_InstantHeal(t *testing.T) {
	rt := &attr.Runtime{CurHP: 10, CurMP: 50, Shield: 0}
	h := &mockHost{base: attr.Base{Level: 1}, rt: rt}
	mgr := NewManager(DemoRegistry())
	mgr.Add(5, "demo_instant_heal", h)
	if rt.CurHP != 40 {
		t.Fatalf("heal got %d", rt.CurHP)
	}
}

func TestManager_StunExpires(t *testing.T) {
	rt := &attr.Runtime{CurHP: 50, CurMP: 50}
	h := &mockHost{base: attr.Base{Level: 1}, rt: rt}
	mgr := NewManager(DemoRegistry())
	mgr.Add(1, "demo_stun", h)
	mods := mgr.Tick(120, h)
	if !mods.Control.HasStun() {
		t.Fatal("still active at 120")
	}
	mods = mgr.Tick(121, h)
	if mods.Control.HasStun() {
		t.Fatal("stun should clear after expire")
	}
	if mods.Control != 0 {
		t.Fatalf("control residue %#v", mods.Control)
	}
}
