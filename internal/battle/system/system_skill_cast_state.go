package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/system/skill/skill_effect"
	"battle/internal/battle/system/target_selector"
)

// CastStateSystem
// @Description: CastStateSystem  释放技能
type CastStateSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.SkillCastState]
}

func (s *CastStateSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.SkillCastState](w)
}

func (s *CastStateSystem) Update(dt float64) {
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
			s.ApplyEffects(e, state.SkillId)
			//  切换到后摇阶段
			cd := config.SkillAfterCastFrames(state.SkillId)
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

func (s *CastStateSystem) ApplyEffects(caster ecs.Entity, skillID int) {
	w := s.world
	if w == nil || caster == 0 || !w.EntityExists(caster) {
		return
	}

	desc := config.GetSkillConfigByID(int32(skillID))

	if desc == nil {
		return
	}
	for _, eid := range desc.EffectIDs {
		effectDesc := config.GetSkillEffectConfigByID(int32(eid))
		if effectDesc == nil {
			continue
		}
		//  选取目标
		targets := target_selector.Select(w, caster, int32(effectDesc.TargetSelectID))
		//  执行效果
		for _, t := range targets {
			ctx := &skill_effect.Context{
				Word:     w,
				Caster:   caster,
				Target:   t,
				EffectId: int32(effectDesc.EffectID),
			}
			skill_effect.Apply(ctx)
		}
	}
}
