package test

import (
	"path/filepath"
	"runtime"
	"testing"

	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// battleConfigDir 返回与本文件同目录下的 battle_config 绝对路径（不依赖进程工作目录）。
func battleConfigDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "battle_config"))
}

func restoreConfigTab(t *testing.T, prev *config.Tables) {
	t.Helper()
	t.Cleanup(func() {
		config.Tab = prev
	})
}

func loadFixtureConfig(t *testing.T) {
	t.Helper()
	prev := config.Tab
	config.Load(battleConfigDir(t))
	restoreConfigTab(t, prev)
}

func swapConfigTab(t *testing.T, tab *config.Tables) {
	t.Helper()
	prev := config.Tab
	config.Tab = tab
	restoreConfigTab(t, prev)
}

func TestBuff_Fixture_AddBuff900(t *testing.T) {
	loadFixtureConfig(t)

	w := ecs.NewWorld(64)
	component.RegisterCombatTypesWorld(w)
	caster := w.CreateEntity()
	target := w.CreateEntity()

	if ok := buff.AddBuff(w, caster, target, 900); !ok {
		t.Fatal("AddBuff(900) 期望成功")
	}
	bl, ok := w.GetComponent(target, &component.BuffList{})
	if !ok {
		t.Fatal("目标缺少 BuffList")
	}
	list := bl.(*component.BuffList)
	if len(list.Buffs) != 1 {
		t.Fatalf("实例数期望 1，实际 %d", len(list.Buffs))
	}
	if list.Buffs[0].BuffId != 900 || list.Buffs[0].Stacks != 1 {
		t.Fatalf("实例异常: %+v", list.Buffs[0])
	}
	if buff.AddBuff(w, caster, target, 0) {
		t.Fatal("Buff 编号为 0 时应失败")
	}
}

func TestBuff_StackAdd_SecondApply_IncrementsStacks(t *testing.T) {
	loadFixtureConfig(t)

	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)
	caster := w.CreateEntity()
	target := w.CreateEntity()

	if !buff.AddBuff(w, caster, target, 900) || !buff.AddBuff(w, caster, target, 900) {
		t.Fatal("fixture 中 900 为叠加策略，两次施加应成功")
	}
	bl, _ := w.GetComponent(target, &component.BuffList{})
	list := bl.(*component.BuffList)
	if len(list.Buffs) != 1 {
		t.Fatalf("叠加后仍为单槽，实例数=%d", len(list.Buffs))
	}
	if list.Buffs[0].Stacks != 2 {
		t.Fatalf("层数期望 2，实际 %d", list.Buffs[0].Stacks)
	}
}

func TestBuff_StackIgnore_SecondApplyFails(t *testing.T) {
	tab := &config.Tables{
		BuffConfigConfigByID: map[int32]*config.BuffConfig{
			701: {
				ID:            701,
				DurationFrame: 60,
				MaxStack:      5,
				StackBehavior: config.BuffStackIgnore,
				EffectType:    config.BufferEffectStatChange,
				ParamsString:  []string{config.AttrArmor},
				Params:        []float64{1},
				CoolingFrame:  0,
			},
		},
	}
	swapConfigTab(t, tab)

	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)
	caster := w.CreateEntity()
	target := w.CreateEntity()

	if !buff.AddBuff(w, caster, target, 701) {
		t.Fatal("首次施加应成功")
	}
	if buff.AddBuff(w, caster, target, 701) {
		t.Fatal("忽略策略下第二次施加应失败")
	}
	bl, _ := w.GetComponent(target, &component.BuffList{})
	if n := len(bl.(*component.BuffList).Buffs); n != 1 {
		t.Fatalf("实例数期望 1，实际 %d", n)
	}
}

func TestBuff_Tick_AppliesArmorStatModifier(t *testing.T) {
	loadFixtureConfig(t)

	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)
	caster := w.CreateEntity()
	target := w.CreateEntity()

	if !buff.AddBuff(w, caster, target, 900) {
		t.Fatal("施加失败")
	}
	bl, _ := w.GetComponent(target, &component.BuffList{})
	list := bl.(*component.BuffList)

	buff.Tick(w, target, list)

	sm, ok := w.GetComponent(target, &component.StatModifiers{})
	if !ok {
		t.Fatal("Tick 后应写入 StatModifiers")
	}
	mod := sm.(*component.StatModifiers)
	want := int(3 * list.Buffs[0].Stacks)
	if mod.ArmorDelta != want {
		t.Fatalf("护甲增量期望 %d（fixture params 3×层数），实际 %d", want, mod.ArmorDelta)
	}
}

func TestBuff_RemoveBuff_Manual(t *testing.T) {
	loadFixtureConfig(t)

	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)
	target := w.CreateEntity()
	if !buff.AddBuff(w, target, target, 900) {
		t.Fatal("施加失败")
	}
	bl, _ := w.GetComponent(target, &component.BuffList{})
	list := bl.(*component.BuffList)

	buff.RemoveBuff(w, target, list, 900)

	if len(list.Buffs) != 0 {
		t.Fatalf("移除后实例数应为 0，实际 %d", len(list.Buffs))
	}
}
