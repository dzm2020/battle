package room_factory

import (
	"battle/internal/battle/component"
	"battle/internal/battle/system"
)

func pvpBuilder(spec *Spec) error {
	//  初始化system
	system.AddCombatSystems(spec.World)

	//  创建怪物
	if err := spawnPlayersOnGridWithTeam(spec, spec.Enemy, component.SideTypeBlue); err != nil {
		return err
	}
	//  创建玩家对象
	if err := spawnPlayersOnGridWithTeam(spec, spec.Self, component.SideTypeRed); err != nil {
		return err
	}

	return nil
}
