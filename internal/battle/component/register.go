// Package component 定义战斗 ECS 组件（纯数据）。
//
// 组件类型文件统一使用 comp_ 前缀，例如 comp_attributes.go、comp_skill.go。
// register.go 为组件注册入口，不使用 comp_ 前缀；施法写入见 system/skill 包。
package component

import "battle/ecs"

// Register 注册战斗组件类型；局内资源（网格、刷怪队列、TPS 等）由 [room] 或单测通过 ecs.AddResource 注入。
func Register(w *ecs.World) {
	RegisterCombatTypes(w.Registry())
}

// RegisterCombatTypes 注册战斗相关组件类型 ID；创建 World 后、实例化 Query 前应调用一次。
func RegisterCombatTypes(r *ecs.ComponentRegistry) {
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
