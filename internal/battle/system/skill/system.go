package skill

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

type System struct {
	world *ecs.World
	q     *ecs.Query[*component.SkillCastState]
}

func (s *System) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.SkillCastState](w)
}

func (s *System) Update(dt float64) {
}

// AddSkill 若配置表中存在 skillConfigID，则在 entity 上确保存在 [component.SkillSet]，
// 并追加一条 [RuntimeSkill]（已存在相同 ConfigID 时不重复追加，返回 true）。
func (s *System) AddSkill(w *ecs.World, entity ecs.Entity, skillConfigID int32) bool {
	if w == nil || entity == 0 || !w.EntityExists(entity) {
		return false
	}
	cfg := config.GetSkillConfigByID(skillConfigID)
	if cfg == nil {
		return false
	}
	set := ecs.EnsureGetComponent[*component.SkillSet](w, entity)
	for _, rs := range set.Skills {
		if rs != nil && rs.ConfigID == cfg.ID {
			return true
		}
	}
	set.Skills = append(set.Skills, &component.RuntimeSkill{ConfigID: cfg.ID, CurrentCooldown: 0})
	return true
}

// RemoveSkill 从 entity 的 [SkillSet] 中删除 ConfigID 等于 skillConfigID 的条目（至少删一条匹配项）。
// 若移除后无任何技能，则移除 [SkillSet] 组件。未挂载 SkillSet 或未找到匹配项时返回 false。
func (s *System) RemoveSkill(w *ecs.World, entity ecs.Entity, skillConfigID int32) bool {
	if w == nil || entity == 0 || !w.EntityExists(entity) {
		return false
	}
	c, ok := w.GetComponent(entity, &component.SkillSet{})
	if !ok {
		return false
	}
	set, ok := c.(*component.SkillSet)
	if !ok || set == nil || len(set.Skills) == 0 {
		return false
	}
	want := int(skillConfigID)
	found := false
	out := set.Skills[:0]
	for _, rs := range set.Skills {
		if rs != nil && rs.ConfigID == want {
			found = true
			continue
		}
		out = append(out, rs)
	}
	set.Skills = out
	if len(set.Skills) == 0 {
		w.RemoveComponent(entity, &component.SkillSet{})
	}
	return found
}
