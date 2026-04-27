package system

import (
	"battle/ecs"
	"battle/internal/battle/action"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// SkillCastValidationSystem 统一校验冷却、消耗、状态等
type SkillCastValidationSystem struct {
	world *ecs.World
	q     *ecs.Query3[*component.SkillSet, *component.Attributes, *component.SkillCastRequest]
}

func NewSkillCastValidationSystem() *SkillCastValidationSystem { return &SkillCastValidationSystem{} }

func (s *SkillCastValidationSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery3[*component.SkillSet, *component.Attributes, *component.SkillCastRequest](w)
}

func (s *SkillCastValidationSystem) Update(dt float64) {
	if s.world == nil || s.q == nil {
		return
	}

	s.q.ForEach(func(e ecs.Entity, set *component.SkillSet, attrs *component.Attributes, req *component.SkillCastRequest) {
		if req == nil || set == nil || attrs == nil {
			s.world.RemoveComponent(e, &component.SkillCastRequest{})
			return
		}
		//  有技能没释放完则不能释放
		if casting, ok := s.world.GetComponent(e, &component.SkillCastState{}); ok && casting.(*component.SkillCastState).IsCasting {
			s.world.RemoveComponent(e, &component.SkillCastRequest{})
			return
		}

		// 检测是否存在技能
		rs := findRuntimeSkill(set, int(req.SkillID))
		if rs == nil {
			s.world.RemoveComponent(e, &component.SkillCastRequest{})
			return
		}
		skillCfg, ok := config.Tab.SkillConfigByID[int32(rs.ConfigID)]
		if !ok || skillCfg == nil {
			s.world.RemoveComponent(e, &component.SkillCastRequest{})
			return
		}

		// 是否处于控制状态（眩晕、沉默）阻止施法
		if !action.CanAct(s.world, e) {
			s.world.RemoveComponent(e, &component.SkillCastRequest{})
			return
		}
		if cs, ok := s.world.GetComponent(e, &component.ControlState{}); ok && cs.(*component.ControlState).Flags.HasSilence() {
			s.world.RemoveComponent(e, &component.SkillCastRequest{})
			return
		}

		// 检测是在 CD 中
		if rs.CurrentCooldown > 0 {
			s.world.RemoveComponent(e, &component.SkillCastRequest{})
			return
		}

		resourceKey, cost := skillCost(skillCfg)
		if cost > 0 {
			if attrs.Get(resourceKey) < cost {
				s.world.RemoveComponent(e, &component.SkillCastRequest{})
				return
			}
			//  扣除资源
			attrs.Add(resourceKey, -cost)
		}

		s.world.RemoveComponent(e, &component.SkillCastRequest{})
		//  记录冷却
		rs.CurrentCooldown = skillCfg.CooldownFrames

		castState := &component.SkillCastState{
			IsCasting:       true,
			SkillId:         rs.ConfigID,
			Phase:           component.SkillStagePreCast,
			RemainingFrames: skillCfg.PreCastFrames,
			TargetEntity:    req.TargetEntity,
			CastPosition:    &req.CastPosition,
		}

		s.world.AddComponent(e, castState)
	})
}

func findRuntimeSkill(set *component.SkillSet, skillID int) *component.RuntimeSkill {
	for _, v := range set.Skills {
		if v != nil && v.ConfigID == skillID {
			return v
		}
	}
	return nil
}

func skillCost(desc *config.SkillBaseConfig) (resourceKey string, cost int) {
	if desc == nil {
		return "", 0
	}
	cost = desc.ConsumeValue
	if cost <= 0 {
		return "", 0
	}
	switch desc.ConsumeType {
	case config.AttrHP:
		return config.AttrHp, cost
	case config.AttrEnergy:
		return "energy", cost
	case config.AttrRage:
		return "rage", cost
	case config.AttrMP:
		fallthrough
	default:
		return config.AttrMana, cost
	}
}
