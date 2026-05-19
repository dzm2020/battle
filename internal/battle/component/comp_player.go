package component

import (
	"battle/ecs"
	"battle/internal/battle/pb"
)

type Player struct {
	ID    uint32         // 玩家ID
	Base  *pb.PlayerBase // 玩家基础数据
	Units map[uint32]ecs.Entity
}

func (*Player) Component() {}
