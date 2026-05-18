package component

import "battle/ecs"

// Init 注册战斗组件类型；局内资源请用 [system/runtime.Install] 注入 [runtime.BattleContext]。
func Init(w *ecs.World) {
	InitCombatTypes(w.Registry())
}

// InitCombatTypes 注册战斗相关组件类型 ID；创建 World 后、实例化 Query 前应调用一次。
func InitCombatTypes(r *ecs.ComponentRegistry) {
	r.Register(&DamageQueue{})
	r.Register(&ResolvedDamage{})
	r.Register(&Attributes{})
	r.Register(&BuffList{}) // Buff 运行时列表（内含 BuffInstance 缓冲）
	r.Register(&BuffStatModifiers{})
	r.Register(&BuffControlState{})
	r.Register(&Team{})
	r.Register(&Transform2D{})
	r.Register(&SkillSet{})
	r.Register(&SkillCastRequest{})
	r.Register(&SkillCastState{})
	r.Register(&PendingHeal{})
}

// RegisterCombatTypesWorld 向战斗 World 注册全部战斗组件类型（测试用，不含 BattleContext）。
func RegisterCombatTypesWorld(w *ecs.World) {
	InitCombatTypes(w.Registry())
}
