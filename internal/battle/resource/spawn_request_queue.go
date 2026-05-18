package resource

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/pb"
	"fmt"
)

// SpawnRequest 延迟刷单位请求；写入 [runtime.BattleContext].SpawnQueue，由 [system.SpawnSystem] 消费。
type SpawnRequest struct {
	UnitID     int32              // 单位配置表 ID（与 entity_factory.CreateByConfigID 一致）
	Side       component.SideType // 阵营
	CellX      int                // 网格 X；-1 表示由系统选空位
	CellY      int                // 网格 Y（对应 Transform2D.Y）；-1 表示自动选位
	TeamEntity ecs.Entity         // 玩家编队时写入 Team.Entity；怪物可为 0
	Data       *pb.PlayerUnit
	Components []ecs.Component
}

// SpawnRequestQueue 刷怪请求队列；仅存于 [runtime.BattleContext]，不作为 ECS 组件。
type SpawnRequestQueue struct {
	Queue []*SpawnRequest
}

// SpawnQueue 返回刷怪请求队列。
func SpawnQueue(w *ecs.World) (*SpawnRequestQueue, bool) {
	q := ecs.GetResource[SpawnRequestQueue](w)
	return q, q != nil
}

// EnqueueSpawn 向本局刷怪队列追加请求。
func EnqueueSpawn(w *ecs.World, req *SpawnRequest) error {
	queue, ok := SpawnQueue(w)
	if queue == nil || !ok {
		return fmt.Errorf("enqueue spawn request queue is nil")
	}
	queue.Queue = append(queue.Queue, req)
	return nil
}
