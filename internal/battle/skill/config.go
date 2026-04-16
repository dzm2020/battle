package skill

// School 技能学派：沉默仅拦截魔法。
type School uint8

const (
	SchoolPhysical School = iota
	SchoolMagic
)

// TargetMode 目标选取模式。
type TargetMode uint8

const (
	TargetModeNone TargetMode = iota
	TargetModeSelf
	TargetModeEnemy
	TargetModeAlly
)

// SkillConfig 技能配表。
type SkillConfig struct {
	ID             string
	School         School
	MPCost         int64
	CooldownFrames uint64
	CastRange      float64
	WindupFrames   uint64
	TargetMode     TargetMode
	Damage         int64    // 结算时对有效目标造成的物理向伤害基数；0 表示无直接伤害
	Heal           int64    // 对有效目标治疗
	TargetBuffIDs  []string // 命中后对目标施加的 Buff ID 列表
}

// Registry 技能表。
type Registry struct {
	byID map[string]SkillConfig
}

func NewRegistry() *Registry {
	return &Registry{byID: make(map[string]SkillConfig)}
}

func (r *Registry) Register(c SkillConfig) {
	r.byID[c.ID] = c
}

func (r *Registry) Get(id string) (SkillConfig, bool) {
	c, ok := r.byID[id]
	return c, ok
}

// MustDemoRegistry 第 5 天演示三条技能。
func MustDemoRegistry() *Registry {
	r := NewRegistry()
	r.Register(SkillConfig{
		ID: "strike", School: SchoolPhysical, MPCost: 0, CooldownFrames: 30,
		CastRange: 2, WindupFrames: 0, TargetMode: TargetModeEnemy,
		Damage: 10,
	})
	r.Register(SkillConfig{
		ID: "fireball", School: SchoolMagic, MPCost: 15, CooldownFrames: 90,
		CastRange: 10, WindupFrames: 45, TargetMode: TargetModeEnemy,
		Damage: 25, TargetBuffIDs: []string{"demo_poison"},
	})
	r.Register(SkillConfig{
		ID: "focus", School: SchoolPhysical, MPCost: 5, CooldownFrames: 120,
		CastRange: 0, WindupFrames: 0, TargetMode: TargetModeSelf,
		Heal: 0, TargetBuffIDs: []string{"demo_amp", "demo_strong"},
	})
	// 第 7 天：演示控制类技能（对敌上眩晕 / 减速）
	r.Register(SkillConfig{
		ID: "hammer_stun", School: SchoolPhysical, MPCost: 10, CooldownFrames: 180,
		CastRange: 2, WindupFrames: 0, TargetMode: TargetModeEnemy,
		Damage: 5, TargetBuffIDs: []string{"demo_stun"},
	})
	r.Register(SkillConfig{
		ID: "ice_slow", School: SchoolMagic, MPCost: 8, CooldownFrames: 60,
		CastRange: 8, WindupFrames: 0, TargetMode: TargetModeEnemy,
		Damage: 3, TargetBuffIDs: []string{"demo_slow"},
	})
	return r
}
