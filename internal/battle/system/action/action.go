// Package action 提供实体行动资格判定（控制状态等）。
package action

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// CanAct 若无控制组件或未被眩晕（[component.BuffControlState]），则允许一般行动；技能施放可再校验沉默等。
func CanAct(w *ecs.World, e ecs.Entity) bool {
	c, ok := w.GetComponent(e, &component.BuffControlState{})
	if !ok {
		return true
	}
	return !c.(*component.BuffControlState).Flags.HasStun()
}
