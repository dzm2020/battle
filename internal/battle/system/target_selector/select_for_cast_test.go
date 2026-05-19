package target_selector

import (
	"path/filepath"
	"runtime"
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/system/attrs"
)

func loadTargetConfig(t *testing.T) {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", "..", "test", "battle_config"))
	config.Load(dir)
}

func spawnUnit(w *ecs.World, side component.SideType, hpCur, hpMax int) ecs.Entity {
	e := w.CreateEntity()
	w.AddComponent(e, &component.Team{Side: side})
	a := ecs.EnsureGetComponent[*component.Attributes](w, e)
	attrs.SetRange(a, config.AttrHp, hpCur, hpMax)
	w.AddComponent(e, &component.Transform2D{X: 0, Y: 0})
	return e
}

func TestSelectForCast_PrimaryOverridesTableSort(t *testing.T) {
	loadTargetConfig(t)
	w := ecs.NewWorld(8)
	component.RegisterCombatTypes(w.Registry())

	caster := spawnUnit(w, component.SideTypeRed, 100, 100)
	weaker := spawnUnit(w, component.SideTypeBlue, 30, 100)
	stronger := spawnUnit(w, component.SideTypeBlue, 90, 100)

	got := Select(w, caster, 15)
	if len(got) != 1 || got[0] != stronger {
		t.Fatalf("Select(15) got=%v want=%v", got, stronger)
	}

	castGot := SelectForCast(w, caster, weaker, 15)
	if len(castGot) != 1 || castGot[0] != weaker {
		t.Fatalf("SelectForCast got=%v want=%v", castGot, weaker)
	}
}
