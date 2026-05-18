package room_factory

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/resource"
	"battle/internal/battle/system"
)

func pvpBuilder(w *ecs.World, spec *resource.RoomSpec) error {
	//  初始化system
	system.AddCombatSystems(w)

	//  创建怪物
	if err := spawnPlayersOnGridWithTeam(w, spec.Enemy, component.SideTypeBlue); err != nil {
		return err
	}
	//  创建玩家对象
	if err := spawnPlayersOnGridWithTeam(w, spec.Self, component.SideTypeRed); err != nil {
		return err
	}

	return nil
}
