package skill

// 与 [internal/battle/config/skill_config] 枚举值对齐，供技能模板与 JSON 共用。

type TargetScope uint8

const (
	TargetScopeIllegal TargetScope = iota
	TargetScopeSelf
	TargetScopeSingle
	TargetScopeCone
	TargetScopeCircle
	TargetScopeLineRect
	TargetScopeMulti
	TargetScopeFullScreen
	TargetScopeChain
	TargetScopeRandom
)

type CampRelation uint8

const (
	CampEnemy CampRelation = iota
	CampAllyIncludeSelf
	CampAllyExcludeSelf
	CampEveryone
	CampSpecificSide
)

type PickRule uint8

const (
	PickNone PickRule = iota
	PickNearest
	PickFarthest
	PickHPCurrentAsc
	PickHPPercentAsc
	PickAttackHighest
)

type ResourceKind uint8

const (
	ResourceNone ResourceKind = iota
	ResourceMana
	ResourceRage
	ResourceEnergy
)

type EffectKind int

const (
	EffectDamage EffectKind = iota
	EffectHeal
	EffectApplyBuff
)
