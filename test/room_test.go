package test

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"battle/ecs"
	"battle/internal/battle/attributes"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/room"
	"battle/internal/battle/unit"
)

func battleConfigDirForRoom(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "battle_config"))
}

func TestRoom(t *testing.T) {
	dir := battleConfigDirForRoom(t)
	config.Load(dir)

	player := &unit.Player{
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

	r, err := room.CreateRoom(1, []*unit.Player{player})
	if err != nil {
		t.Fatal(err)
	}
	if r.Phase() != room.PhaseLobby {
		t.Fatalf("创建后阶段应为 Lobby，实际 %s", r.Phase())
	}

	w := r.World()
	q := ecs.NewQuery[*component.Health](w)
	n := 0
	q.ForEach(func(_ ecs.Entity, _ *component.Health) { n++ })
	// 副本配置含 1 只怪物（Unit 模板 id=1）+ 玩家单位各 1 个实体
	if n != 2 {
		t.Fatalf("副本内应有 2 个带生命单位（1 玩家 + 1 怪），实际 %d", n)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := r.StartBattle(ctx); err != nil {
		t.Fatal(err)
	}
	if r.Phase() != room.PhaseFighting {
		t.Fatalf("开战后期望 Fighting，实际 %s", r.Phase())
	}
	if r.Loop() == nil {
		t.Fatal("Loop 不应为 nil")
	}

	if err := r.StartBattle(ctx); !errors.Is(err, room.ErrWrongPhase) {
		t.Fatalf("重复 StartBattle 应返回 ErrWrongPhase，实际 %v", err)
	}
}
