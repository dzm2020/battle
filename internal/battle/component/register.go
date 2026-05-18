package component

import (
	"battle/ecs"
)

func Init(w *ecs.World) {
	InitCombatTypes(w.Registry())
	InitResource(w)
}

// InitCombatTypes 注册战斗相关组件类型 ID；创建 World 后、实例化 Query 前应调用一次。
func InitCombatTypes(r *ecs.ComponentRegistry) {
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

// InitResource 注入局内单例资源（如空间网格）；开房或 [room.Room.SetGrid] 时调用。
func InitResource(w *ecs.World) {
	if w == nil {
		return
	}
	ecs.InsertResource(w, &SpawnRequestQueue{})
}

// RegisterCombatTypesWorld 向战斗 World 注册全部战斗组件类型（测试用）。
func RegisterCombatTypesWorld(w *ecs.World) {
	InitCombatTypes(w.Registry())
}
