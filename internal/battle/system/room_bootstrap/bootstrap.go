package room_bootstrap

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"battle/internal/battle/resource"
	"fmt"
)

// Bootstrap 按副本类型装配房间：先挂载已注册的 [Installer]，再执行 [Spawner] 入队单位。
func Bootstrap(w *ecs.World, spec *resource.RoomSpec) error {
	desc := config.GetDungeonConfigByID(spec.DungeonId)
	if desc == nil {
		return fmt.Errorf("room_bootstrap: dungeon config not found: %d", spec.DungeonId)
	}
	if err := runInstaller(w, desc.Type); err != nil {
		return err
	}
	return runSpawner(w, spec, desc.Type)
}
