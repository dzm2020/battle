package room_builder

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"battle/internal/battle/land"
	"battle/internal/battle/pb"
	"battle/internal/battle/room"
	"errors"

	"github.com/duke-git/lancet/v2/maputil"
)

var (
	ErrNoDungeonConfig = errors.New("no dungeon config")
	ErrNoMapConfig     = errors.New("no map config")
)

type Spec struct {
	World       *ecs.World
	Desc        *config.DungeonConfig
	Self, Enemy *pb.Player
}

// RoomBootstrap 按副本类型装配 [room.Room]：应调用 [room.Room.SetGrid]、创建怪物/玩家实体等；在 [component.Register] 之后调用。
type builder func(ctx *Spec) error

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

func Build(t int32, spec *Spec) error {
	builder := getBuilder(t)
	return builder(spec)
}
