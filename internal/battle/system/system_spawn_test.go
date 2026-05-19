package system

import (
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/land"
	"battle/internal/battle/pb"
	"battle/internal/battle/resource"
)

func TestFlushSpawnQueue_PlayerUnitsMapping(t *testing.T) {
	mustLoadBattleConfig(t)
	w := newCombatWorld(t)
	grid, err := land.CreateGridByID(1)
	if err != nil {
		t.Fatal(err)
	}
	ecs.AddResource(w, grid)
	ecs.AddResource(w, &resource.SpawnRequestQueue{})

	teamEntity := w.CreateEntity()
	w.AddComponent(teamEntity, &component.Player{
		ID:    1,
		Units: make(map[uint32]ecs.Entity),
	})

	unit := &pb.PlayerUnit{ID: 1}
	if err := resource.EnqueueSpawn(w, &resource.SpawnRequest{
		UnitID:     1,
		Side:       component.SideTypeRed,
		CellX:      0,
		CellY:      0,
		TeamEntity: teamEntity,
		Data:       unit,
	}); err != nil {
		t.Fatal(err)
	}

	FlushSpawnQueue(w)

	pc, ok := w.GetComponent(teamEntity, &component.Player{})
	if !ok {
		t.Fatal("Player component missing")
	}
	player := pc.(*component.Player)
	if len(player.Units) != 1 {
		t.Fatalf("Units map len: got %d want 1", len(player.Units))
	}
	spawned, ok := player.Units[1]
	if !ok || spawned == 0 {
		t.Fatal("expected unit id 1 mapped to spawned entity")
	}
	if !w.EntityExists(spawned) {
		t.Fatal("mapped entity should exist in world")
	}
}
