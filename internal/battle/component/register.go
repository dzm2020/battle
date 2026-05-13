package component

import "battle/ecs"

func Register(w *ecs.World) {
	RegisterCombatTypes(w.Registry())
}

// RegisterCombatTypes 注册战斗相关组件类型 ID；创建 World 后、实例化 Query 前应调用一次。
func RegisterCombatTypes(r *ecs.ComponentRegistry) {
	r.Register(&DamageQueue{})
	r.Register(&ResolvedDamage{})
	r.Register(&Attributes{})
	r.Register(&Health{})
	r.Register(&BuffList{}) // Buff 运行时列表（内含 BuffInstance 缓冲）
	r.Register(&BuffStatModifiers{})
	r.Register(&BuffControlState{})
	r.Register(&Team{})
	r.Register(&Transform2D{})
	r.Register(&SkillSet{})
	r.Register(&SkillCastRequest{})
	r.Register(&CastIntent{})
	r.Register(&SkillCastState{})
	r.Register(&PendingHeal{})

}

// RegisterCombatTypesWorld 向战斗 [ecs.World] 注册全部战斗组件类型 ID（测试与房间初始化用）。
func RegisterCombatTypesWorld(w *ecs.World) {
	RegisterCombatTypes(w.Registry())
}
