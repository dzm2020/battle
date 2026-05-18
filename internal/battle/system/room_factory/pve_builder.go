package room_factory

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/resource"
	"battle/internal/battle/system"
)

func pveBuilder(w *ecs.World, spec *resource.RoomSpec) error {
	system.AddCombatSystems(w)
	//  创建怪物
	if err := spawnMonstersForDesc(w, spec, component.SideTypeBlue); err != nil {
		return err
	}
	//  创建玩家对象
	if err := spawnPlayersOnGridWithTeam(w, spec.Self, component.SideTypeRed); err != nil {
		return err
	}

	return nil
}
