package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/skill"
)

// SkillChannelSystem 推进 [component.SkillCastState]；吟唱帧归零当帧结算效果并写入冷却（资源已在发起吟唱时扣除）。
// 须在 [CooldownSystem] 之后、[SkillIntentSystem] 之前注册，以便同一帧内先完成吟唱结算再处理新的施法意图；
// 且须在 [DamageSystem] 之前，使本帧技能产生的 [PendingDamage] 参与减免。
type SkillChannelSystem struct {
	world       *ecs.World
	skillConfig *skill.CatalogConfig

	qChannel *ecs.Query[*component.SkillCastState]
}

// NewSkillChannelSystem skillConfig 可为 nil（内部退化为空表）。
func NewSkillChannelSystem(skillConfig *skill.CatalogConfig) *SkillChannelSystem {
	return &SkillChannelSystem{skillConfig: skillConfig}
}

func (s *SkillChannelSystem) Initialize(w *ecs.World) {
	s.world = w
	if s.skillConfig == nil {
		s.skillConfig = skill.NewCatalogConfig()
	}
	s.qChannel = ecs.NewQuery[*component.SkillCastState](w)
}

func (s *SkillChannelSystem) Update(dt float64) {
	s.qChannel.ForEach(func(e ecs.Entity, st *component.SkillCastState) {
		if st.FramesLeft <= 0 {
			return
		}
		st.FramesLeft--
		if st.FramesLeft > 0 {
			return
		}
		sk, ok := s.skillConfig.Get(st.SkillID)
		if !ok {
			s.world.RemoveComponent(e, &component.SkillCastState{})
			return
		}
		targets := skill.ResolveTargets(s.world, e, st.PrimaryTarget, sk)
		skill.ExecuteEffects(s.world, e, targets, sk)
		if su, ok := s.world.GetComponent(e, &component.SkillUser{}); ok {
			startSkillCooldown(su.(*component.SkillUser), sk.ID, sk.CooldownFrames)
		}
		s.world.RemoveComponent(e, &component.SkillCastState{})
	})
}
