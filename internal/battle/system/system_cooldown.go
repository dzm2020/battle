package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// CooldownSystem 仅负责将 [component.SkillUser].CooldownRemaining 中各技能剩余帧每帧减 1，
// 并在归零时删除键。须排在 [SkillIntentSystem] 之前，以便本帧先完成“上帧进入冷却”的递减再接收新施法。
type CooldownSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.SkillUser]
}

// NewCooldownSystem 创建冷却递减系统（无外部配置）。
func NewCooldownSystem() *CooldownSystem { return &CooldownSystem{} }

func (s *CooldownSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.SkillUser](w)
}

func (s *CooldownSystem) Update(dt float64) {
	s.q.ForEach(func(e ecs.Entity, su *component.SkillUser) {
		if su.CooldownRemaining == nil {
			return
		}
		for id, left := range su.CooldownRemaining {
			left--
			if left <= 0 {
				delete(su.CooldownRemaining, id)
			} else {
				su.CooldownRemaining[id] = left
			}
		}
	})
}
