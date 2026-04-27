package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// SkillCooldownSystem
// @Description: 更新技能CD
type SkillCooldownSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.SkillSet]
}

func (s *SkillCooldownSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.SkillSet](w)
}

func (s *SkillCooldownSystem) Update(dt float64) {
	if s.world == nil || s.q == nil {
		return
	}
	s.q.ForEach(func(e ecs.Entity, set *component.SkillSet) {
		if set == nil || len(set.Skills) == 0 {
			return
		}
		for _, rs := range set.Skills {
			if rs == nil || rs.CurrentCooldown <= 0 {
				continue
			}
			rs.CurrentCooldown--
			if rs.CurrentCooldown < 0 {
				rs.CurrentCooldown = 0
			}
		}
	})
}
