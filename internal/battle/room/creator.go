package room

import (
	"battle/internal/battle/config"
	"battle/internal/battle/unit"
)

func CreateRoom(dungeonId int32, players []*unit.Player) (*Room, error) {
	desc := config.GetDungeonConfigByID(dungeonId)
	if desc == nil {
		return nil, ErrNoDungeonConfig
	}
	roomId := GetManager().NextID()
	//  创建房间
	r := newRoom(roomId)
	//  添加怪物
	for _, monsterId := range desc.Monster {
		_, _ = unit.SpawnUnitByID(r.World(), monsterId)
	}
	//  添加玩家
	for _, player := range players {
		_, _ = unit.SpawnUnitByPlayer(r.World(), player)
	}
	return r, nil
}
