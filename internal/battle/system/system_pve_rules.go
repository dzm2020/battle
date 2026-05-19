package system

import "battle/ecs"

// PVERulesSystem PVE 专用扩展点：副本波次、星级、掉落等规则在此帧更新。
// 当前为占位，与 [AddPVPSystems] 的 [PVPRulesSystem] 区分管线。
type PVERulesSystem struct {
	world *ecs.World
}

func (s *PVERulesSystem) Initialize(w *ecs.World) { s.world = w }

func (s *PVERulesSystem) Update(_ float64) {}
