package system

import (
	"battle/ecs"
	"battle/internal/battle/skill"
)

type SkillIntentSystem struct {
	cfg *skill.CatalogConfig
}

func NewSkillIntentSystem(cfg *skill.CatalogConfig) *SkillIntentSystem {
	return &SkillIntentSystem{cfg: cfg}
}

func (s *SkillIntentSystem) Initialize(w *ecs.World) { _ = w }

func (s *SkillIntentSystem) Update(dt float64) { _ = s.cfg }
