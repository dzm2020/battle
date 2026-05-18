package room_bootstrap

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"fmt"

	"github.com/duke-git/lancet/v2/maputil"
)

// Installer 按副本类型向 World 挂载战斗 System 管线。
type Installer func(w *ecs.World)

var (
	installers        = maputil.NewConcurrentMap[int32, Installer](2)
	defaultInstaller Installer
)

// RegisterInstaller 注册副本类型对应的 System 安装器（通常在 system 包 init 中调用）。
func RegisterInstaller(dungeonType config.DungeonType, installer Installer) {
	if installer == nil {
		panic("room_bootstrap: RegisterInstaller with nil installer")
	}
	installers.Set(dungeonType, installer)
}

// SetDefaultInstaller 设置未匹配副本类型时使用的 System 安装器。
func SetDefaultInstaller(installer Installer) {
	defaultInstaller = installer
}

func runInstaller(w *ecs.World, dungeonType int32) error {
	installer, ok := installers.Get(dungeonType)
	if !ok || installer == nil {
		installer = defaultInstaller
	}
	if installer == nil {
		return fmt.Errorf("room_bootstrap: no installer for dungeon type %d", dungeonType)
	}
	installer(w)
	return nil
}
