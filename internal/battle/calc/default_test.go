package calc

import (
	"testing"

	"battle/internal/battle/attr"
)

func TestDefaultCalculator_DerivedFromBase(t *testing.T) {
	c := DefaultCalculator{}
	b := attr.Base{Level: 10, STR: 20, AGI: 15, INT: 12, VIT: 18}
	d := c.DerivedFromBase(b)

	if d.MaxHP != 100+200+15*18 {
		t.Fatalf("MaxHP got %d", d.MaxHP)
	}
	if d.MaxMP != 50+50+8*12 {
		t.Fatalf("MaxMP got %d", d.MaxMP)
	}
	if d.ATK != 70 {
		t.Fatalf("ATK got %d", d.ATK)
	}
	if d.DEF != 28 {
		t.Fatalf("DEF got %d", d.DEF)
	}
	if d.CritRate < 0.05 || d.CritRate > 0.5 {
		t.Fatalf("CritRate out of range: %v", d.CritRate)
	}
	if d.CritDamage < 1.25 || d.CritDamage > 3.0 {
		t.Fatalf("CritDamage out of range: %v", d.CritDamage)
	}
	if d.PhysMitigation < 0 || d.PhysMitigation >= 1 {
		t.Fatalf("PhysMitigation out of range: %v", d.PhysMitigation)
	}
}

func TestDefaultCalculator_ApplyMaxToRuntime(t *testing.T) {
	c := DefaultCalculator{}
	d := attr.Derived{MaxHP: 100, MaxMP: 50}
	rt := attr.Runtime{CurHP: 200, CurMP: 10, Shield: -5}
	c.ApplyMaxToRuntime(&rt, d)
	if rt.CurHP != 100 || rt.CurMP != 10 || rt.Shield != 0 {
		t.Fatalf("unexpected runtime: %+v", rt)
	}
}
