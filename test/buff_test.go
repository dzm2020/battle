package test

import (
	"path/filepath"
	"runtime"
	"testing"

	"battle/ecs"
	"battle/internal/battle/attributes"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/control"
	"battle/internal/battle/room"
	"battle/internal/battle/system"
	"battle/internal/battle/unit"
)

func battleConfigDirForBuff(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "battle_config"))
}

// testPlayer 与 room_test 一致：50 HP，与怪物（模板 hp=100）区分。
func testPlayerForBuff() *unit.Player {
	return &unit.Player{
		ID: 1,
		Units: map[uint32]*unit.PlayerUnit{
			1: {
				ID: 1,
				Stats: []attributes.Attribute{
					{Type: config.AttrHp, InitValue: 50, MaxValue: 50},
				},
				Ability: []int32{1},
			},
		},
	}
}

func findPlayerEntity(t *testing.T, w *ecs.World) ecs.Entity {
	t.Helper()
	var ent ecs.Entity
	ecs.NewQuery[*component.Health](w).ForEach(func(e ecs.Entity, h *component.Health) {
		if h.Current == 50 {
			ent = e
		}
	})
	if ent == 0 {
		t.Fatal("未找到玩家实体（期望 Current HP=50）")
	}
	return ent
}

func buffStacks(t *testing.T, w *ecs.World, e ecs.Entity, buffID uint32) int {
	t.Helper()
	c, ok := w.GetComponent(e, &component.BuffList{})
	if !ok {
		t.Fatal("缺少 BuffList")
	}
	bl, _ := c.(*component.BuffList)
	if bl == nil {
		t.Fatal("BuffList 类型断言失败")
	}
	for _, bi := range bl.Buffs {
		if bi.BuffId == buffID {
			return bi.Stacks
		}
	}
	t.Fatalf("未找到 BuffId=%d", buffID)
	return 0
}

// TestBuff 沿房间创建后的 ECS 世界注册战斗管线，每帧 [World.Update] 驱动 Buff 汇总，
// 覆盖叠层策略与若干效果类型（属性 / 控制）。
func TestBuff(t *testing.T) {
	dir := battleConfigDirForBuff(t)
	dt := 1.0 / 60.0

	t.Run("叠加·层数相加900", func(t *testing.T) {
		config.Load(dir)
		r, err := room.CreateRoom(1, []*unit.Player{testPlayerForBuff()})
		if err != nil {
			t.Fatal(err)
		}
		w := r.World()
		system.AddCombatSystems(w)
		e := findPlayerEntity(t, w)

		for range 3 {
			if !buff.AddBuff(w, e, e, 900) {
				t.Fatal("AddBuff 900 失败")
			}
		}
		w.Update(dt)

		if buffStacks(t, w, e, 900) != 3 {
			t.Fatalf("900 期望层数 3")
		}
		c, ok := w.GetComponent(e, &component.StatModifiers{})
		if !ok {
			t.Fatal("缺少 StatModifiers")
		}
		sm, _ := c.(*component.StatModifiers)
		// Params[0]=3，每层叠一层：delta = 3 * Stacks = 9
		if sm.ArmorDelta != 9 {
			t.Fatalf("护甲增量期望 9，实际 %d", sm.ArmorDelta)
		}
	})

	t.Run("叠加·刷新持续时间901", func(t *testing.T) {
		config.Load(dir)
		r, err := room.CreateRoom(1, []*unit.Player{testPlayerForBuff()})
		if err != nil {
			t.Fatal(err)
		}
		w := r.World()
		system.AddCombatSystems(w)
		e := findPlayerEntity(t, w)

		if !buff.AddBuff(w, e, e, 901) || !buff.AddBuff(w, e, e, 901) {
			t.Fatal("AddBuff 901 失败")
		}
		w.Update(dt)

		if buffStacks(t, w, e, 901) != 1 {
			t.Fatalf("901 refresh 期望层数仍为 1")
		}
		c, ok := w.GetComponent(e, &component.StatModifiers{})
		if !ok {
			t.Fatal("缺少 StatModifiers")
		}
		sm, _ := c.(*component.StatModifiers)
		if sm.AttackDamageDelta != -10 {
			t.Fatalf("虚弱期望攻击力增量 -10，实际 %d", sm.AttackDamageDelta)
		}
	})

	t.Run("叠加·已有则忽略902", func(t *testing.T) {
		config.Load(dir)
		r, err := room.CreateRoom(1, []*unit.Player{testPlayerForBuff()})
		if err != nil {
			t.Fatal(err)
		}
		w := r.World()
		system.AddCombatSystems(w)
		e := findPlayerEntity(t, w)

		if !buff.AddBuff(w, e, e, 902) {
			t.Fatal("首次 AddBuff 902 失败")
		}
		if buff.AddBuff(w, e, e, 902) {
			t.Fatal("第二次施加应被忽略")
		}
		w.Update(dt)

		if buffStacks(t, w, e, 902) != 1 {
			t.Fatalf("902 ignore 期望层数 1")
		}
		c, ok := w.GetComponent(e, &component.StatModifiers{})
		if !ok {
			t.Fatal("缺少 StatModifiers")
		}
		sm, _ := c.(*component.StatModifiers)
		if sm.ArmorDelta != 1 {
			t.Fatalf("护甲增量期望 1，实际 %d", sm.ArmorDelta)
		}
	})

	t.Run("叠加·替换实例903", func(t *testing.T) {
		config.Load(dir)
		r, err := room.CreateRoom(1, []*unit.Player{testPlayerForBuff()})
		if err != nil {
			t.Fatal(err)
		}
		w := r.World()
		system.AddCombatSystems(w)
		e := findPlayerEntity(t, w)

		if !buff.AddBuff(w, e, e, 903) || !buff.AddBuff(w, e, e, 903) {
			t.Fatal("AddBuff 903 失败")
		}
		w.Update(dt)

		if buffStacks(t, w, e, 903) != 1 {
			t.Fatalf("903 replace 期望层数重置为 1")
		}
		c, ok := w.GetComponent(e, &component.StatModifiers{})
		if !ok {
			t.Fatal("缺少 StatModifiers")
		}
		sm, _ := c.(*component.StatModifiers)
		if sm.ArmorDelta != 2 {
			t.Fatalf("护甲增量期望 2，实际 %d", sm.ArmorDelta)
		}
	})

	t.Run("效果·控制眩晕904", func(t *testing.T) {
		config.Load(dir)
		r, err := room.CreateRoom(1, []*unit.Player{testPlayerForBuff()})
		if err != nil {
			t.Fatal(err)
		}
		w := r.World()
		system.AddCombatSystems(w)
		e := findPlayerEntity(t, w)

		if !buff.AddBuff(w, e, e, 904) {
			t.Fatal("AddBuff 904 失败")
		}
		w.Update(dt)

		c, ok := w.GetComponent(e, &component.ControlState{})
		if !ok {
			t.Fatal("缺少 ControlState")
		}
		cs, _ := c.(*component.ControlState)
		if cs.Flags&control.FlagStunned == 0 {
			t.Fatal("期望眩晕位已置位")
		}
	})
}
