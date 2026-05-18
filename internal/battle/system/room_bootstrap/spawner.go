package room_bootstrap

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"battle/internal/battle/resource"
	"fmt"

	"github.com/duke-git/lancet/v2/maputil"
)

// Spawner 按副本类型将单位入队到 [resource.SpawnRequestQueue]。
type Spawner func(w *ecs.World, spec *resource.RoomSpec) error

var (
	spawners        = maputil.NewConcurrentMap[int32, Spawner](2)
	defaultSpawner Spawner = spawnerPVE
)

func init() {
	RegisterSpawner(config.DungeonTypePVE, spawnerPVE)
	RegisterSpawner(config.DungeonTypePVP, spawnerPVP)
}

// RegisterSpawner 注册副本类型对应的单位入队逻辑。
func RegisterSpawner(dungeonType config.DungeonType, spawner Spawner) {
	if spawner == nil {
		panic("room_bootstrap: RegisterSpawner with nil spawner")
	}
	spawners.Set(dungeonType, spawner)
}

// SetDefaultSpawner 设置未匹配副本类型时使用的入队逻辑。
func SetDefaultSpawner(spawner Spawner) {
	defaultSpawner = spawner
}

func runSpawner(w *ecs.World, spec *resource.RoomSpec, dungeonType int32) error {
	spawner, ok := spawners.Get(dungeonType)
	if !ok || spawner == nil {
		spawner = defaultSpawner
	}
	if spawner == nil {
		return fmt.Errorf("room_bootstrap: no spawner for dungeon type %d", dungeonType)
	}
	return spawner(w, spec)
}
