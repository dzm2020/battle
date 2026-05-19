package test

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/system/attrs"
	"battle/internal/battle/pb"
	"battle/internal/battle/room"
	"battle/internal/battle/system/room_bootstrap"
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

	player := &pb.Player{
		ID: 1,
		Units: map[uint32]*pb.PlayerUnit{
			1: {
				ID: 1,
				Stats: []pb.Attribute{
					{Type: config.AttrHp, InitValue: 50, MaxValue: 50},
				},
				Ability: []int32{1},
			},
		},
	}

	r, err := room.CreateRoom(1, &room_bootstrap.Spec{Self: player})
	if err != nil {
		t.Fatal(err)
	}
	if r.Phase() != room.PhaseFighting {
		t.Fatalf("CreateRoom 已自动开战，期望 Fighting，实际 %v", r.Phase())
	}
	if r.Loop() == nil {
		t.Fatal("Loop 不应为 nil")
	}

	w := r.World()
	// 同步推进一帧，确保 SpawnSystem 消费入队请求
	w.Update(1.0 / 60.0)

	q := ecs.NewQuery[*component.Attributes](w)
	n := 0
	q.ForEach(func(_ ecs.Entity, a *component.Attributes) {
		if attrs.HPMax(a) > 0 {
			n++
		}
	})
	if n != 2 {
		t.Fatalf("副本内应有 2 个带生命单位（1 玩家 + 1 怪），实际 %d", n)
	}

	if err := r.StartBattle(context.Background()); !errors.Is(err, room.ErrWrongPhase) {
		t.Fatalf("重复 StartBattle 应返回 ErrWrongPhase，实际 %v", err)
	}
}

func TestCreateRoom_RejectsPVPDungeonWithoutEnemy(t *testing.T) {
	dir := battleConfigDirForRoom(t)
	config.Load(dir)

	_, err := room.CreateRoom(2, &room_bootstrap.Spec{Self: &pb.Player{ID: 1}})
	if !errors.Is(err, room.ErrUseCreatePVPRoom) {
		t.Fatalf("PVP 副本缺少 Enemy 应返回 ErrUseCreatePVPRoom，实际 %v", err)
	}
}

func TestCreatePVPRoom(t *testing.T) {
	dir := battleConfigDirForRoom(t)
	config.Load(dir)

	unitStats := []pb.Attribute{{Type: config.AttrHp, InitValue: 50, MaxValue: 50}}
	red := &pb.Player{
		ID: 1,
		Units: map[uint32]*pb.PlayerUnit{
			1: {ID: 1, Stats: unitStats, Ability: []int32{1}},
		},
	}
	blue := &pb.Player{
		ID: 2,
		Units: map[uint32]*pb.PlayerUnit{
			1: {ID: 1, Stats: unitStats, Ability: []int32{1}},
		},
	}

	r, err := room.CreatePVPRoom(2, red, blue)
	if err != nil {
		t.Fatal(err)
	}
	w := r.World()
	w.Update(1.0 / 60.0)

	var nRed, nBlue int
	ecs.NewQuery[*component.Team](w).ForEach(func(_ ecs.Entity, team *component.Team) {
		switch team.Side {
		case component.SideTypeRed:
			nRed++
		case component.SideTypeBlue:
			nBlue++
		}
	})
	if nRed != 1 || nBlue != 1 {
		t.Fatalf("PVP 双方单位各应带 Team 各 1，实际 red=%d blue=%d", nRed, nBlue)
	}
}
