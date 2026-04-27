package system

import (
	"battle/ecs"
	"battle/internal/battle/skill"
)

type SkillChannelSystem struct {
	cfg *skill.CatalogConfig
}

func NewSkillChannelSystem(cfg *skill.CatalogConfig) *SkillChannelSystem {
	return &SkillChannelSystem{cfg: cfg}
}

func (s *SkillChannelSystem) Initialize(w *ecs.World) { _ = w }

func (s *SkillChannelSystem) Update(dt float64) { _ = s.cfg }
