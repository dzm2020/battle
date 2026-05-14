package room_builder

import (
	"battle/internal/battle/config"
	"battle/internal/battle/land"
	"battle/internal/battle/room"
	"errors"

	"github.com/duke-git/lancet/v2/maputil"
)

var (
	ErrNoDungeonConfig = errors.New("no dungeon config")
	ErrNoMapConfig     = errors.New("no map config")
)

// RoomBootstrap 按副本类型装配 [room.Room]：应调用 [room.Room.SetGrid]、创建怪物/玩家实体等；在 [component.Register] 之后调用。
type builder func(r *room.Room, desc *config.DungeonConfig, options *Options) error

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

func doBuilder(t int32) builder {
	if b, ok := builders.Get(t); ok {
		return b
	}
	return defaultBuilder
}

// CreateRoom 根据 dungeonId 加载副本配置，并按 [config.DungeonConfig.Type] 选择已注册的装配逻辑创建房间。
func CreateRoom(dungeonId int32, options *Options) (*room.Room, error) {
	desc := config.GetDungeonConfigByID(dungeonId)
	if desc == nil {
		return nil, ErrNoDungeonConfig
	}

	grid, err := land.CreateGridByID(desc.MapID)
	if err != nil {
		return nil, err
	}
	r := room.New()
	r.SetGrid(grid)
	if err = doBuilder(desc.Type)(r, desc, options); err != nil {
		return nil, err
	}
	room.GetManager().Add(r)
	return r, nil
}
