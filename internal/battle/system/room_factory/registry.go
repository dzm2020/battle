package room_factory

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"battle/internal/battle/resource"
	"fmt"

	"github.com/duke-git/lancet/v2/maputil"
)

// RoomBootstrap 按副本类型装配 [room.Room]：须在 [component.Init] 与 [runtime.Install] 之后调用；通过 SpawnQueue 入队单位。
type builder func(w *ecs.World, spec *resource.RoomSpec) error

var (
	builders               = maputil.NewConcurrentMap[int32, builder](1)
	defaultBuilder builder = pveBuilder
)

func init() {
	registerBuilder(config.DungeonTypePVE, pveBuilder)
	registerBuilder(config.DungeonTypePVP, pvpBuilder)
}

func registerBuilder(dungeonType config.DungeonType, b builder) {
	if b == nil {
		panic("room_builder: RegisterRoomBootstrap with nil bootstrap")
	}
	builders.Set(dungeonType, b)
}

func getBuilder(t int32) builder {
	if b, ok := builders.Get(t); ok {
		return b
	}
	return defaultBuilder
}

func Create(w *ecs.World, spec *resource.RoomSpec) error {
	desc := config.GetDungeonConfigByID(spec.DungeonId)
	if desc == nil {
		return fmt.Errorf("room_builder: CreateRoom: no desc found")
	}
	builder := getBuilder(desc.Type)
	return builder(w, spec)
}
