package room_builder

import (
	"battle/internal/battle/component"
	"battle/internal/battle/system"
)

func pveBuilder(spec *Spec) error {
	//  初始化system
	system.AddCombatSystems(spec.World)
	//  初始化全局资源
	component.InitResource(spec.World)
	//  创建怪物
	if err := spawnMonstersForDesc(spec, component.SideTypeBlue); err != nil {
		return err
	}
	//  创建玩家对象
	if err := spawnPlayersOnGridWithTeam(spec, spec.Self, component.SideTypeRed); err != nil {
		return err
	}

	return nil
}
