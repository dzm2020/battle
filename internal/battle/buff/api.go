package buff

import "battle/ecs"

// AddBuff 向 target 挂载 buff 模板（便捷封装，等价于 [Manager.AddBuff]）。
func AddBuff(w *ecs.World, caster, target ecs.Entity, buffId uint32) bool {
	if w == nil {
		return false
	}
	return NewManager(w).AddBuff(caster, target, buffId)
}
