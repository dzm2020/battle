package component

import (
	"battle/ecs"
	"battle/internal/battle/pb"
)

// SpawnRequest 延迟刷单位请求；由 [system.SpawnSystem] 消费后移除并调用 unit 创建逻辑。
// 挂于任意实体（如房间刷怪点、技能召唤源）；同一实体可挂多条请求时由系统定义处理顺序。
type SpawnRequest struct {
	UnitID     int32      // 单位配置表 ID（与 unit.CreateByID 一致）
	Side       SideType   // 阵营
	CellX      int        // 网格 X；-1 表示由系统选空位
	CellY      int        // 网格 Y（对应 Transform2D.Y）；-1 表示自动选位
	TeamEntity ecs.Entity // 玩家编队时写入 Team.Entity；怪物可为 0
	Data       *pb.PlayerUnit
	Components []ecs.Component
}

type SpawnRequestQueue struct {
	Queue []*SpawnRequest
}

func (*SpawnRequestQueue) Component() {}
