package room_builder

import "battle/internal/battle/pb"

type Options struct {
	DungeonId   int32
	Self, Enemy *pb.Player
}
