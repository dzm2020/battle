package test

import (
	"path/filepath"
	"runtime"
	"slices"
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/control"
	"battle/internal/battle/target_selector"
)

func battleConfigDirForTarget(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "battle_config"))
}

func newTargetTestWorld(t *testing.T) *ecs.World {
	t.Helper()
	w := ecs.NewWorld(16)
	component.RegisterCombatTypesWorld(w)
	return w
}

// 带 [Health] + [Team] + 属性「hp」+ 可选 [Transform2D]；attrHP 供 property 筛选使用。
func spawnUnit(
	w *ecs.World,
	side uint8,
	hpCur, hpMax, attrHP int,
	x, y float64,
) ecs.Entity {
	e := w.CreateEntity()
	w.AddComponent(e, &component.Team{Side: side})
	w.AddComponent(e, &component.Health{Current: hpCur, Max: hpMax})
	a := ecs.EnsureGetComponent[*component.Attributes](w, e)
	a.Set("hp", attrHP)
	if x != 0 || y != 0 {
		w.AddComponent(e, &component.Transform2D{X: x, Y: y})
	} else {
		w.AddComponent(e, &component.Transform2D{X: 0, Y: 0})
	}
	return e
}

func sortEntityIDs(ents []ecs.Entity) {
	slices.Sort(ents)
}

func sameEntitySet(t *testing.T, got, want []ecs.Entity) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("数量不同 got=%d want=%d", len(got), len(want))
	}
	g := append([]ecs.Entity(nil), got...)
	w := append([]ecs.Entity(nil), want...)
	sortEntityIDs(g)
	sortEntityIDs(w)
	if !slices.Equal(g, w) {
		t.Fatalf("实体集合不同\ngot  %v\nwant %v", g, w)
	}
}

func TestTargetSelect(t *testing.T) {
	dir := battleConfigDirForTarget(t)
	config.Load(dir)

	t.Run("空或非法入参", func(t *testing.T) {
		var w *ecs.World
		if target_selector.Select(w, 1, 1) != nil {
			t.Fatal("nil world 应返回 nil")
		}
		ww := newTargetTestWorld(t)
		e := spawnUnit(ww, 0, 100, 100, 100, 0, 0)
		if target_selector.Select(ww, 0, 1) != nil {
			t.Fatal("caster=0 应返回 nil")
		}
		ww.RemoveEntity(e)
		if target_selector.Select(ww, e, 1) != nil {
			t.Fatal("不存在的实体作施法者应返回 nil")
		}
		if target_selector.Select(ww, 1, 9999) != nil {
			t.Fatal("未知选择表 id 应返回 nil")
		}
	})

	t.Run("max_count=0 不选", func(t *testing.T) {
		w := newTargetTestWorld(t)
		c := spawnUnit(w, 0, 100, 100, 100, 0, 0)
		if target_selector.Select(w, c, 99) != nil {
			t.Fatal("期望 nil")
		}
	})

	t.Run("IncludeSelf false 仅有自己", func(t *testing.T) {
		w := newTargetTestWorld(t)
		c := spawnUnit(w, 0, 100, 100, 100, 0, 0)
		if len(target_selector.Select(w, c, 1)) != 0 {
			t.Fatal("不含自己时应无目标")
		}
	})

	t.Run("IncludeSelf true 可选自己", func(t *testing.T) {
		w := newTargetTestWorld(t)
		c := spawnUnit(w, 0, 100, 100, 100, 0, 0)
		got := target_selector.Select(w, c, 2)
		if len(got) != 1 || got[0] != c {
			t.Fatalf("期望仅选中施法者自身 got=%v", got)
		}
	})

	t.Run("筛选·阵营敌方", func(t *testing.T) {
		w := newTargetTestWorld(t)
		caster := spawnUnit(w, 0, 100, 100, 100, 0, 0)
		e1 := spawnUnit(w, 1, 50, 50, 50, 5, 0)
		e2 := spawnUnit(w, 1, 50, 50, 50, 8, 0)
		got := target_selector.Select(w, caster, 10)
		sameEntitySet(t, got, []ecs.Entity{e1, e2})
	})

	t.Run("筛选·阵营友方不含自己", func(t *testing.T) {
		w := newTargetTestWorld(t)
		caster := spawnUnit(w, 0, 100, 100, 100, 0, 0)
		ally := spawnUnit(w, 0, 80, 80, 80, 1, 0)
		_ = spawnUnit(w, 1, 40, 40, 40, 2, 0)
		got := target_selector.Select(w, caster, 11)
		sameEntitySet(t, got, []ecs.Entity{ally})
	})

	t.Run("筛选·状态眩晕位", func(t *testing.T) {
		w := newTargetTestWorld(t)
		caster := spawnUnit(w, 0, 100, 100, 100, 0, 0)
		stunned := spawnUnit(w, 1, 40, 40, 40, 0, 0)
		_ = spawnUnit(w, 1, 40, 40, 40, 1, 0)
		w.AddComponent(stunned, &component.ControlState{Flags: control.FlagStunned})
		got := target_selector.Select(w, caster, 12)
		sameEntitySet(t, got, []ecs.Entity{stunned})
	})

	t.Run("筛选·属性 hp 小于", func(t *testing.T) {
		w := newTargetTestWorld(t)
		caster := spawnUnit(w, 0, 100, 100, 100, 0, 0)
		low := spawnUnit(w, 1, 40, 40, 30, 0, 0)
		_ = spawnUnit(w, 1, 90, 90, 90, 0, 0)
		got := target_selector.Select(w, caster, 13)
		sameEntitySet(t, got, []ecs.Entity{low})
	})

	t.Run("排序·当前生命升序取一", func(t *testing.T) {
		w := newTargetTestWorld(t)
		caster := spawnUnit(w, 0, 100, 100, 100, 0, 0)
		weaker := spawnUnit(w, 1, 30, 100, 50, 0, 0)
		_ = spawnUnit(w, 1, 80, 100, 90, 0, 0)
		got := target_selector.Select(w, caster, 14)
		if len(got) != 1 || got[0] != weaker {
			t.Fatalf("期望当前生命最低者 got=%v want=%v", got, weaker)
		}
	})

	t.Run("排序·当前生命降序取一", func(t *testing.T) {
		w := newTargetTestWorld(t)
		caster := spawnUnit(w, 0, 100, 100, 100, 0, 0)
		stronger := spawnUnit(w, 1, 90, 100, 90, 0, 0)
		_ = spawnUnit(w, 1, 30, 100, 50, 0, 0)
		got := target_selector.Select(w, caster, 15)
		if len(got) != 1 || got[0] != stronger {
			t.Fatalf("期望当前生命最高者 got=%v want=%v", got, stronger)
		}
	})

	t.Run("排序·距施法者最近", func(t *testing.T) {
		w := newTargetTestWorld(t)
		caster := spawnUnit(w, 0, 100, 100, 100, 0, 0)
		near := spawnUnit(w, 1, 50, 50, 50, 2, 0)
		_ = spawnUnit(w, 1, 50, 50, 50, 100, 0)
		got := target_selector.Select(w, caster, 16)
		if len(got) != 1 || got[0] != near {
			t.Fatalf("期望距离最近者 got=%v want=%v", got, near)
		}
	})
}
