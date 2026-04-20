package system

import (
	"battle/ecs"
	"battle/internal/battle/action"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
	"battle/internal/battle/skill"
	"slices"
)

// SkillIntentSystem 消费 [component.CastIntent]：校验、扣费、瞬发结算或创建吟唱状态。
// 须在 [SkillChannelSystem] 之后注册，以便本帧吟唱先结算后再接受新意图；须在 [DamageSystem] 之前，
// 以便瞬发技能产生的 [PendingDamage] 参与减免。
type SkillIntentSystem struct {
	world       *ecs.World
	skillConfig *skill.CatalogConfig
	buffConfig  *buff.DefinitionConfig

	qIntent *ecs.Query[*component.CastIntent]
}

// NewSkillIntentSystem skillConfig / buffConfig 可为 nil（内部退化为空表）。
func NewSkillIntentSystem(skillConfig *skill.CatalogConfig, buffConfig *buff.DefinitionConfig) *SkillIntentSystem {
	return &SkillIntentSystem{skillConfig: skillConfig, buffConfig: buffConfig}
}

func (s *SkillIntentSystem) Initialize(w *ecs.World) {
	s.world = w
	if s.skillConfig == nil {
		s.skillConfig = skill.NewCatalogConfig()
	}
	if s.buffConfig == nil {
		s.buffConfig = buff.NewDefinitionConfig()
	}
	s.qIntent = ecs.NewQuery[*component.CastIntent](w)
}

func (s *SkillIntentSystem) Update(dt float64) {
	s.qIntent.ForEach(func(e ecs.Entity, intent *component.CastIntent) {
		s.tryCast(e, intent)
	})
}

func (s *SkillIntentSystem) tryCast(caster ecs.Entity, intent *component.CastIntent) {
	removeIntent := func() { s.world.RemoveComponent(caster, &component.CastIntent{}) }

	if !action.CanAct(s.world, caster) {
		removeIntent()
		return
	}
	if _, busy := s.world.GetComponent(caster, &component.SkillCastState{}); busy {
		removeIntent()
		return
	}

	suComp, ok := s.world.GetComponent(caster, &component.SkillUser{})
	if !ok {
		removeIntent()
		return
	}
	su := suComp.(*component.SkillUser)

	sk, ok := s.skillConfig.Get(intent.SkillID)
	if !ok {
		removeIntent()
		return
	}
	if !slices.Contains(su.GrantedSkillIDs, intent.SkillID) {
		removeIntent()
		return
	}
	if left, ok := su.CooldownRemaining[intent.SkillID]; ok && left > 0 {
		removeIntent()
		return
	}
	if !skill.ValidCastTargets(s.world, caster, intent.Target, sk) {
		removeIntent()
		return
	}
	if !paySkillCost(su, sk) {
		removeIntent()
		return
	}

	if sk.CastFrames > 0 {
		s.world.AddComponent(caster, &component.SkillCastState{
			SkillID:       intent.SkillID,
			PrimaryTarget: intent.Target,
			FramesLeft:    sk.CastFrames,
		})
		removeIntent()
		return
	}

	targets := skill.ResolveTargets(s.world, caster, intent.Target, sk)
	skill.ExecuteEffects(s.world, targets, sk, s.buffConfig)
	startSkillCooldown(su, sk.ID, sk.CooldownFrames)
	removeIntent()
}

func paySkillCost(su *component.SkillUser, sk skill.SkillConfig) bool {
	if sk.Resource == skill.ResourceNone {
		return sk.Cost == 0
	}
	if sk.Cost < 0 {
		return false
	}
	switch sk.Resource {
	case skill.ResourceMana:
		if su.Mana < sk.Cost {
			return false
		}
		su.Mana -= sk.Cost
	case skill.ResourceRage:
		if su.Rage < sk.Cost {
			return false
		}
		su.Rage -= sk.Cost
	case skill.ResourceEnergy:
		if su.Energy < sk.Cost {
			return false
		}
		su.Energy -= sk.Cost
	default:
		return false
	}
	return true
}
