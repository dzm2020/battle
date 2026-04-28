package room

import (
	"battle/ecs"
	"context"
)

type IRoom interface {
	ID() uint64
	DungeonId() int32
	Phase() Phase
	World() *ecs.World
	StartBattle(ctx context.Context) error
	Shutdown()
}
