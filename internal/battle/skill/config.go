package skill

import "battle/internal/battle/component"

// SkillConfig 技能静态模板（内存注册或 JSON 加载）。
type SkillConfig struct {
	ID             uint32          `json:"id"`
	Resource       ResourceKind    `json:"resource"`
	Cost           int             `json:"cost"`
	CooldownFrames int             `json:"cooldownFrames"`
	Scope          TargetScope     `json:"scope"`
	Camp           CampRelation    `json:"camp"`
	PickRule       PickRule        `json:"pickRule,omitempty"`
	CampSide       uint8           `json:"campSide,omitempty"`
	AOERadius      float64         `json:"aoeRadius,omitempty"`
	MaxTargets     int             `json:"maxTargets,omitempty"`
	CastFrames     int             `json:"castFrames"`
	Effects        []EffectConfig  `json:"effects"`
}

// EffectConfig 单条技能效果。
type EffectConfig struct {
	Kind       EffectKind          `json:"kind"`
	Amount     int                 `json:"amount,omitempty"`
	DamageType component.DamageType `json:"damageType,omitempty"`
	BuffDefID  uint32              `json:"buffDefId,omitempty"`
}
