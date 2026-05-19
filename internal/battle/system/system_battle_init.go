package system

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"battle/internal/battle/land"
	"battle/internal/battle/resource"
	"battle/internal/battle/system/room_bootstrap"
)

func init() {
	room_bootstrap.SetDefaultInstaller(AddPVESystems)

	room_bootstrap.RegisterInstaller(config.DungeonTypePVE, AddPVESystems)
	room_bootstrap.RegisterInstaller(config.DungeonTypePVP, AddPVPSystems)
}

// BattleInitSystem 在首帧（[resource.TPS].Frame == 0）创建网格、Bootstrap 并按副本挂载战斗 System。
// Bootstrap 入队后立即 [FlushSpawnQueue]，使单位在同帧生成（新 System 尚不参与当次 World.Update 遍历）。
//
// 由 [room.Create] 单独注册，不包含在 [AddCoreCombatSystems] 内。
type BattleInitSystem struct {
	world       *ecs.World
	initialized bool
}

func (s *BattleInitSystem) Initialize(w *ecs.World) {
	s.world = w
}

func (s *BattleInitSystem) Update(_ float64) {
	if s.initialized {
		return
	}
	tps := ecs.GetResource[resource.TPS](s.world)
	if tps == nil || tps.Frame != 0 {
		return
	}

	spec := ecs.GetResource[resource.RoomSpec](s.world)
	if spec == nil {
		return
	}

	desc := config.GetDungeonConfigByID(spec.DungeonId)
	if desc == nil {
		return
	}

	grid, err := land.CreateGridByID(desc.MapID)
	if err != nil {
		return
	}

	ecs.AddResource(s.world, grid)

	if err = room_bootstrap.Bootstrap(s.world, spec); err != nil {
		return
	}
	FlushSpawnQueue(s.world)
	s.initialized = true
}
