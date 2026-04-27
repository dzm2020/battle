package buff

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/log"
)

// System 递减持续时间、汇总属性并写入 [StatModifiers]/[ControlState]。
// 须在 [DamageSystem] 之前运行，以便本帧 DoT 写入的 [PendingDamage] 参与结算。
type System struct {
	world   *ecs.World
	q       *ecs.Query[*component.BuffList]
	manager *Manager
}

// NewBuffSystem 使用全局 [config.Tab.BuffConfigConfigByID] 解析 Buff 模板。
func NewBuffSystem() *System {
	return &System{}
}

func (s *System) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.BuffList](w)
	s.manager = NewManager(w)
	log.Info("[buff] Buff 系统已初始化")
}

// Update 遍历含 [component.BuffList] 的实体：先清零并重算 StatModifiers/ControlState，再逐实例
// 聚合属性，最后递减 FramesLeft 并剔除到期实例。
func (s *System) Update(dt float64) {
	s.q.ForEach(func(e ecs.Entity, bl *component.BuffList) {
		s.manager.Tick(e, bl)
	})
}
