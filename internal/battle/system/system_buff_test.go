package system

import (
	"testing"

	"battle/internal/battle/component"
	"battle/internal/battle/system/buff"
)

func TestBuffSystem_DecrementsDuration(t *testing.T) {
	mustLoadBattleConfig(t)
	w := newCombatWorld(t)
	e := spawnCombatEntity(w, 100, 0)
	if err := buff.Add(w, 0, e, 900); err != nil {
		t.Fatal(err)
	}
	bl, ok := w.GetComponent(e, &component.BuffList{})
	if !ok || len(bl.(*component.BuffList).Buffs) == 0 {
		t.Fatal("expected buff instance")
	}
	before := bl.(*component.BuffList).Buffs[0].DurationFrame

	sys := &BuffSystem{}
	sys.Initialize(w)
	sys.Update(0)

	bl2 := bl.(*component.BuffList)
	if len(bl2.Buffs) == 0 {
		t.Fatal("buff should remain after one tick")
	}
	if bl2.Buffs[0].DurationFrame != before-1 {
		t.Fatalf("duration: got %d want %d", bl2.Buffs[0].DurationFrame, before-1)
	}
}
