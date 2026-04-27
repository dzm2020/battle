package system

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
	"battle/internal/battle/log"
)

// BufferSystem 递减持续时间、汇总属性并写入 [StatModifiers]/[ControlState]。
// 须在 [DamageSystem] 之前运行，以便本帧 DoT 写入的 [PendingDamage] 参与结算。
type BufferSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.BuffList]
}

func (s *BufferSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.BuffList](w)
	log.Info("[buff] Buff 系统已初始化")
}

// Update 遍历含 [component.BuffList] 的实体：先清零并重算 StatModifiers/ControlState，再逐实例
// 聚合属性，最后递减 FramesLeft 并剔除到期实例。
func (s *BufferSystem) Update(dt float64) {
	s.q.ForEach(func(e ecs.Entity, bl *component.BuffList) {
		buff.Tick(s.world, e, bl)
	})
}
