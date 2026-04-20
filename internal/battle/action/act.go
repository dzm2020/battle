// Package action 提供与 Buff 派生控制状态相关的行动判定（与 [component.ControlState]、
// [system.BuffSystem] 配合；沉默等位标志可在此包扩展专用 API）。
package action

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// CanAct 若无控制组件或未被眩晕（[control.FlagStunned]），则允许一般行动；技能施放可再校验沉默等。
func CanAct(w *ecs.World, e ecs.Entity) bool {
	c, ok := w.GetComponent(e, &component.ControlState{})
	if !ok {
		return true
	}
	return !c.(*component.ControlState).Flags.HasStun()
}
