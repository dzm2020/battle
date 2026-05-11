package room

import (
	"battle/internal/battle/config"
	"battle/internal/battle/entity_factory"
	"battle/internal/battle/pb"
)

func CreateRoom(dungeonId int32, players []*pb.Player) (*Room, error) {
	desc := config.GetDungeonConfigByID(dungeonId)
	if desc == nil {
		return nil, ErrNoDungeonConfig
	}
	roomId := GetManager().NextID()
	//  创建房间
	r := newRoom(roomId)
	//  添加怪物
	for _, monsterId := range desc.Monster {
		_, _ = entity_factory.CreateByID(r.World(), monsterId)
	}
	//  添加玩家
	for _, player := range players {
		_, _ = entity_factory.CreateByPlayer(r.World(), player)
	}
	return r, nil
}
