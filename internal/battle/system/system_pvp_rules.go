package system

import "battle/ecs"

// PVPRulesSystem PVP 专用扩展点：同步校验、投降、断线判负等。
// 当前为占位，与 [AddPVESystems] 的 [PVERulesSystem] 区分管线。
type PVPRulesSystem struct {
	world *ecs.World
}

func (s *PVPRulesSystem) Initialize(w *ecs.World) { s.world = w }

func (s *PVPRulesSystem) Update(_ float64) {}
