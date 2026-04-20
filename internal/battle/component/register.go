package component

import "battle/ecs"

// RegisterCombatTypes 注册战斗相关组件类型 ID；创建 World 后、实例化 Query 前应调用一次。
func RegisterCombatTypes(r *ecs.ComponentRegistry) {
	r.Register(&PendingDamage{})
	r.Register(&ResolvedDamage{})
	r.Register(&Health{})
	r.Register(&Attributes{})
	r.Register(&BuffList{})       // Buff 运行时列表（内含 BuffInstance 缓冲）
	r.Register(&StatModifiers{}) // Buff 汇总后的属性增量
	r.Register(&ControlState{})   // Buff 汇总后的控制位
}

// RegisterCombatTypesWorld 便捷封装。
func RegisterCombatTypesWorld(w *ecs.World) {
	RegisterCombatTypes(w.Registry())
}
