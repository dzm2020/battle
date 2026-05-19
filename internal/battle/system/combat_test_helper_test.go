package system

import (
	"path/filepath"
	"runtime"
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/system/attrs"
)

func mustLoadBattleConfig(t *testing.T) {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", "test", "battle_config"))
	if err := config.Load(dir); err != nil {
		t.Fatal(err)
	}
}

func newCombatWorld(t *testing.T) *ecs.World {
	t.Helper()
	w := ecs.NewWorld(16)
	component.Register(w)
	return w
}

func spawnCombatEntity(w *ecs.World, hp, mana int) ecs.Entity {
	e := w.CreateEntity()
	w.AddComponent(e, &component.Team{Side: component.SideTypeRed})
	a := ecs.EnsureGetComponent[*component.Attributes](w, e)
	attrs.SetRange(a, config.AttrHp, hp, hp)
	attrs.SetRange(a, config.AttrMana, mana, mana)
	return e
}
