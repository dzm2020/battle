package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
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

// BattleInitSystem 在 [SpawnSystem] 消费完 [runtime.BattleContext].SpawnQueue 后完成单局开战初始化：
// 快照开局存活阵营数、标记 [runtime.BattleContext].Started，并派发 [event.KindBattleStart]。
//
// 须注册在 [SpawnSystem] 之后、[BuffSystem] 之前（见 [AddSystems]），以便首帧刷怪完成后再进入战斗管线。
type BattleInitSystem struct {
	world *ecs.World
	q     *ecs.Query2[*component.Team, *component.Attributes]
}

func (s *BattleInitSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery2[*component.Team, *component.Attributes](w)
}

func (s *BattleInitSystem) Update(_ float64) {
	tps := ecs.GetResource[resource.TPS](s.world)
	if tps == nil || tps.Frame == 0 {
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

	//  构建房间
	if err = room_bootstrap.Bootstrap(s.world, spec); err != nil {
		return
	}

}
