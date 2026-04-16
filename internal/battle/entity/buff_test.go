package entity

import (
	"testing"

	"battle/internal/battle/attr"
	"battle/internal/battle/calc"
)

func TestEntityTickBuffs_StatATKLayers(t *testing.T) {
	cal := calc.DefaultCalculator{}
	e := New("x", 1, attr.Base{Level: 1, STR: 0})
	e.InitBattle(cal)
	baseAtk := e.Derived.ATK
	e.AddBuff(1, "demo_strong")
	e.AddBuff(1, "demo_strong")
	e.TickBuffs(1, cal)
	if e.Derived.ATK != baseAtk+30 {
		t.Fatalf("atk want %d got %d", baseAtk+30, e.Derived.ATK)
	}
}
