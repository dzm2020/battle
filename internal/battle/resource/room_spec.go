package resource

import (
	"battle/internal/battle/pb"
)

type RoomSpec struct {
	DungeonId   int32
	Self, Enemy *pb.Player
}
