package test

import (
	"path/filepath"
	"runtime"
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/control"
	"battle/internal/battle/skill"
	"battle/internal/battle/system"
)

func battleConfigDirForSkill(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "battle_config"))
}

func newSkillCombatWorld(t *testing.T) *ecs.World {
	t.Helper()
	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)
	return w
}

// 施法者 / 目标：含 Team、Health、Attributes（含 hp、mana），用于校验与伤害结算。
func spawnCombatUnit(w *ecs.World, side uint8, hp, mana int) ecs.Entity {
	e := w.CreateEntity()
	w.AddComponent(e, &component.Team{Side: side})
	w.AddComponent(e, &component.Health{Current: hp, Max: hp})
	a := ecs.EnsureGetComponent[*component.Attributes](w, e)
	a.Set(config.AttrHp, hp)
	a.Set(config.AttrMana, mana)
	return e
}

func healthCurrent(t *testing.T, w *ecs.World, e ecs.Entity) int {
	t.Helper()
	c, ok := w.GetComponent(e, &component.Health{})
	if !ok {
		t.Fatal("缺少 Health")
	}
	return c.(*component.Health).Current
}

func manaCurrent(t *testing.T, w *ecs.World, e ecs.Entity) int {
	t.Helper()
	a, ok := w.GetComponent(e, &component.Attributes{})
	if !ok {
		t.Fatal("缺少 Attributes")
	}
	return a.(*component.Attributes).Get(config.AttrMana)
}

func hasBuff(w *ecs.World, e ecs.Entity, buffID uint32) bool {
	c, ok := w.GetComponent(e, &component.BuffList{})
	if !ok {
		return false
	}
	bl := c.(*component.BuffList)
	for _, b := range bl.Buffs {
		if b != nil && b.BuffId == buffID {
			return true
		}
	}
	return false
}

func castSkillRequest(w *ecs.World, caster ecs.Entity, skillID int32, target ecs.Entity) {
	w.AddComponent(caster, &component.SkillCastRequest{
		SkillID:      skillID,
		TargetEntity: target,
	})
}

// TestSkill 覆盖：释放前校验、目标选取（见 SkillEffect 与 TargetSelect 表）、效果生效（伤害 / 加 Buff）。
func TestSkill(t *testing.T) {
	dir := battleConfigDirForSkill(t)
	config.Load(dir)
	dt := 1.0 / 60.0

	t.Run("管线·真实伤害命中敌方", func(t *testing.T) {
		w := newSkillCombatWorld(t)
		system.AddCombatSystems(w)
		caster := spawnCombatUnit(w, 0, 100, 100)
		target := spawnCombatUnit(w, 1, 100, 0)
		if !skill.AddSkill(w, caster, 1) {
			t.Fatal("AddSkill 失败")
		}
		castSkillRequest(w, caster, 1, target)
		w.Update(dt)

		if healthCurrent(t, w, target) != 58 {
			t.Fatalf("目标当前生命期望 58（100-42 真伤），实际 %d", healthCurrent(t, w, target))
		}
	})

	t.Run("校验·冷却中拒绝再次释放", func(t *testing.T) {
		w := newSkillCombatWorld(t)
		system.AddCombatSystems(w)
		caster := spawnCombatUnit(w, 0, 100, 100)
		target := spawnCombatUnit(w, 1, 100, 0)
		if !skill.AddSkill(w, caster, 1) {
			t.Fatal("AddSkill 失败")
		}
		castSkillRequest(w, caster, 1, target)
		w.Update(dt)
		h1 := healthCurrent(t, w, target)

		castSkillRequest(w, caster, 1, target)
		w.Update(dt)
		if healthCurrent(t, w, target) != h1 {
			t.Fatal("冷却中不应再次结算伤害")
		}
	})

	t.Run("校验·眩晕无法施法", func(t *testing.T) {
		w := newSkillCombatWorld(t)
		system.AddCombatSystems(w)
		caster := spawnCombatUnit(w, 0, 100, 100)
		target := spawnCombatUnit(w, 1, 100, 0)
		w.AddComponent(caster, &component.ControlState{Flags: control.FlagStunned})
		if !skill.AddSkill(w, caster, 1) {
			t.Fatal("AddSkill 失败")
		}
		castSkillRequest(w, caster, 1, target)
		w.Update(dt)

		if healthCurrent(t, w, target) != 100 {
			t.Fatalf("眩晕时应未造成伤害")
		}
	})

	t.Run("校验·沉默无法施法", func(t *testing.T) {
		w := newSkillCombatWorld(t)
		system.AddCombatSystems(w)
		caster := spawnCombatUnit(w, 0, 100, 100)
		target := spawnCombatUnit(w, 1, 100, 0)
		w.AddComponent(caster, &component.ControlState{Flags: control.FlagSilenced})
		if !skill.AddSkill(w, caster, 1) {
			t.Fatal("AddSkill 失败")
		}
		castSkillRequest(w, caster, 1, target)
		w.Update(dt)

		if healthCurrent(t, w, target) != 100 {
			t.Fatalf("沉默时应未造成伤害")
		}
	})

	t.Run("校验·未学会的技能请求被丢弃", func(t *testing.T) {
		w := newSkillCombatWorld(t)
		system.AddCombatSystems(w)
		caster := spawnCombatUnit(w, 0, 100, 100)
		target := spawnCombatUnit(w, 1, 100, 0)
		castSkillRequest(w, caster, 1, target)
		w.Update(dt)

		if healthCurrent(t, w, target) != 100 {
			t.Fatal("无 SkillSet 时不应造成伤害")
		}
	})

	t.Run("校验·法力不足", func(t *testing.T) {
		w := newSkillCombatWorld(t)
		system.AddCombatSystems(w)
		caster := spawnCombatUnit(w, 0, 100, 30)
		target := spawnCombatUnit(w, 1, 100, 0)
		if !skill.AddSkill(w, caster, 3) {
			t.Fatal("AddSkill 失败")
		}
		castSkillRequest(w, caster, 3, target)
		w.Update(dt)

		if healthCurrent(t, w, target) != 100 {
			t.Fatal("法力不足时不应造成伤害")
		}
		if manaCurrent(t, w, caster) != 30 {
			t.Fatal("法力不应被扣除")
		}
	})

	t.Run("校验·法力消耗成功", func(t *testing.T) {
		w := newSkillCombatWorld(t)
		system.AddCombatSystems(w)
		caster := spawnCombatUnit(w, 0, 100, 100)
		target := spawnCombatUnit(w, 1, 100, 0)
		if !skill.AddSkill(w, caster, 3) {
			t.Fatal("AddSkill 失败")
		}
		castSkillRequest(w, caster, 3, target)
		w.Update(dt)

		if manaCurrent(t, w, caster) != 50 {
			t.Fatalf("法力期望剩余 50，实际 %d", manaCurrent(t, w, caster))
		}
		if healthCurrent(t, w, target) != 58 {
			t.Fatalf("伤害仍应生效")
		}
	})

	t.Run("效果·对选取目标施加 Buff", func(t *testing.T) {
		w := newSkillCombatWorld(t)
		system.AddCombatSystems(w)
		caster := spawnCombatUnit(w, 0, 100, 100)
		target := spawnCombatUnit(w, 1, 100, 0)
		if !skill.AddSkill(w, caster, 2) {
			t.Fatal("AddSkill 失败")
		}
		castSkillRequest(w, caster, 2, target)
		w.Update(dt)

		if !hasBuff(w, target, 900) {
			t.Fatal("SkillEffect 20 应向目标添加 Buff 900")
		}
	})
}
