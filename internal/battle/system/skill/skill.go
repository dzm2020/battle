package skill

import (
	"fmt"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

func Add(w *ecs.World, entity ecs.Entity, skillID ...int32) error {
	for _, id := range skillID {
		if err := AddOne(w, entity, id); err != nil {
			return err
		}
	}
	return nil
}

func AddOne(w *ecs.World, entity ecs.Entity, skillID int32) error {
	if entity == 0 || !w.EntityExists(entity) {
		return fmt.Errorf("skill: invalid entity %v", entity)
	}
	cfg := config.GetSkillConfigByID(skillID)
	if cfg == nil {
		return fmt.Errorf("skill: unknown skill id %d", skillID)
	}
	set := ecs.EnsureGetComponent[*component.SkillSet](w, entity)
	for _, rs := range set.Skills {
		if rs != nil && rs.ConfigID == cfg.ID {
			return nil
		}
	}
	set.Skills = append(set.Skills, &component.RuntimeSkill{ConfigID: cfg.ID, CurrentCooldown: 0})
	return nil
}

func Remove(w *ecs.World, entity ecs.Entity, skillID int32) error {
	if w == nil {
		return fmt.Errorf("skill: nil world")
	}
	if entity == 0 || !w.EntityExists(entity) {
		return fmt.Errorf("skill: invalid entity %v", entity)
	}
	c, ok := w.GetComponent(entity, &component.SkillSet{})
	if !ok {
		return fmt.Errorf("skill: entity %v has no SkillSet", entity)
	}
	set, ok := c.(*component.SkillSet)
	if !ok || set == nil || len(set.Skills) == 0 {
		return fmt.Errorf("skill: empty or invalid SkillSet on entity %v", entity)
	}
	want := int(skillID)
	found := false
	out := set.Skills[:0]
	for _, rs := range set.Skills {
		if rs != nil && rs.ConfigID == want {
			found = true
			continue
		}
		out = append(out, rs)
	}
	if !found {
		return fmt.Errorf("skill: skill %d not on entity %v", skillID, entity)
	}
	set.Skills = out
	if len(set.Skills) == 0 {
		w.RemoveComponent(entity, &component.SkillSet{})
	}
	return nil
}

// SetSkillCastRequest 为实体设置施法请求（覆盖同实体上已有请求）。
func SetSkillCastRequest(w *ecs.World, e ecs.Entity, req *component.SkillCastRequest) {
	if w == nil || e == 0 || req == nil {
		return
	}
	w.RemoveComponent(e, &component.SkillCastRequest{})
	w.AddComponent(e, req)
}

// RequestSkillCast 玩法层快捷入口：挂本帧施法请求。
func RequestSkillCast(w *ecs.World, caster ecs.Entity, skillID int32, target ecs.Entity) {
	SetSkillCastRequest(w, caster, &component.SkillCastRequest{
		SkillID:      skillID,
		TargetEntity: target,
	})
}
