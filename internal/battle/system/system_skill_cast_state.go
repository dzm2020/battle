package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/skill"
)

// SkillCastStateSystem
// @Description: 更新技能状态
type SkillCastStateSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.SkillCastState]
}

func NewSkillCastStateSystem() *SkillCastStateSystem { return &SkillCastStateSystem{} }

func (s *SkillCastStateSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.SkillCastState](w)
}

func (s *SkillCastStateSystem) Update(dt float64) {
	if s.world == nil || s.q == nil {
		return
	}
	s.q.ForEach(func(e ecs.Entity, state *component.SkillCastState) {
		if state == nil {
			return
		}
		if state.Phase == component.SkillStageNone {
			return
		}
		//  阶段冷却中
		if state.RemainingFrames > 0 {
			state.RemainingFrames--
		}
		if state.RemainingFrames > 0 {
			return
		}

		switch state.Phase {
		case component.SkillStagePreCast:
			state.Phase = component.SkillStagePostCast
			fallthrough
		case component.SkillStagePostCast:
			skill.ApplySkillEffects(s.world, e, state.TargetEntity, state.SkillId)
			//  切换到后摇阶段
			cd := afterCastFrames(state.SkillId)
			state.Phase = component.SkillStageAfterCast
			state.IsCasting = false
			state.RemainingFrames = cd
			if cd > 0 {
				return
			}
			fallthrough
		case component.SkillStageAfterCast:
			//  后摇结束
			s.world.RemoveComponent(e, &component.SkillCastState{})
		default:
			s.world.RemoveComponent(e, &component.SkillCastState{})
		}
	})
}

func afterCastFrames(skillID int) int {
	cfg, ok := config.Tab.SkillConfigByID[int32(skillID)]
	if !ok || cfg == nil || cfg.AfterCastFrames <= 0 {
		return 0
	}
	return cfg.AfterCastFrames
}
