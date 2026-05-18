package component

import (
	"battle/ecs"
	"battle/internal/battle/pb"
)

// SpawnRequest 延迟刷单位请求；写入 [runtime.BattleContext].SpawnQueue，由 [system.SpawnSystem] 消费。
type SpawnRequest struct {
	UnitID     int32      // 单位配置表 ID（与 entity_factory.CreateByConfigID 一致）
	Side       SideType   // 阵营
	CellX      int        // 网格 X；-1 表示由系统选空位
	CellY      int        // 网格 Y（对应 Transform2D.Y）；-1 表示自动选位
	TeamEntity ecs.Entity // 玩家编队时写入 Team.Entity；怪物可为 0
	Data       *pb.PlayerUnit
	Components []ecs.Component
}

// SpawnRequestQueue 刷怪请求队列；仅存于 [runtime.BattleContext]，不作为 ECS 组件。
type SpawnRequestQueue struct {
	Queue []*SpawnRequest
}
