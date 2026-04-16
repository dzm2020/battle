package skill

// School 技能学派：用于沉默等「按类别封禁」的校验，与客户端表现解耦。
type School uint8

const (
	SchoolPhysical School = iota // 物理：默认不受沉默影响（可按项目改为部分技能受影响）
	SchoolMagic                  // 魔法：受沉默影响
)

// TargetMode 目标选取规则（服务器权威）。客户端只上报「想选谁」，最终是否合法由本配置决定。
type TargetMode uint8

const (
	// TargetModeNone 无需目标（例如全屏 BUFF 占位、仅自身逻辑由 Applier 处理）。
	TargetModeNone TargetMode = iota
	// TargetModeSelf 仅自身：CastInput.Target 可为 nil，生效目标在 Applier 内回退为 Caster。
	TargetModeSelf
	// TargetModeEnemy 敌方单体：需目标、敌对阵营、距离与存活校验。
	TargetModeEnemy
	// TargetModeAlly 友方单体：用于治疗/增益；需存活，阵营与自身关系合法。
	TargetModeAlly
)

// SkillConfig 技能静态配置（配表入口）。
// 第 5 天用内存表注册；后续可替换为 Protobuf / JSON + 版本号热更。
type SkillConfig struct {
	// ID 全局唯一技能 ID，与客户端协议、CD key 一致。
	ID string
	// Name 调试与日志用。
	Name string
	// School 学派：沉默等控制效果按学派过滤。
	School School
	// MPCost 每次成功结算（前摇结束或瞬发）消耗的魔法；0 表示无消耗。
	MPCost int64
	// CooldownFrames 冷却长度，以逻辑帧计（与第 3 天帧循环对齐）。
	CooldownFrames uint64
	// CastRange 施法距离；0 表示不校验距离（常用于纯自身技能）。
	CastRange float64
	// WindupFrames 前摇：>0 时先进入 windup，由 timer 在到期帧尝试结算；0 为瞬发。
	WindupFrames uint32
	// TargetMode 目标模式，决定 Target 是否必填及阵营关系。
	TargetMode TargetMode
}

// Registry 技能配置注册表（只读快照）。战斗中途替换指针需自行保证房间隔离。
type Registry struct {
	byID map[string]SkillConfig
}

// NewRegistry 创建空注册表。
func NewRegistry() *Registry {
	return &Registry{byID: make(map[string]SkillConfig)}
}

// Register 注册或覆盖技能配置；启动阶段调用，战斗线程勿写。
func (r *Registry) Register(cfg SkillConfig) {
	r.byID[cfg.ID] = cfg
}

// Get 查找配置。
func (r *Registry) Get(id string) (SkillConfig, bool) {
	c, ok := r.byID[id]
	return c, ok
}

// MustDemoRegistry 返回教学用默认技能表（普攻 + 火球带前摇）。
func MustDemoRegistry() *Registry {
	reg := NewRegistry()
	reg.Register(SkillConfig{
		ID:             "strike",
		Name:           "普攻",
		School:         SchoolPhysical,
		MPCost:         0,
		CooldownFrames: 20,
		CastRange:      2.5,
		WindupFrames:   0,
		TargetMode:     TargetModeEnemy,
	})
	reg.Register(SkillConfig{
		ID:             "fireball",
		Name:           "火球术",
		School:         SchoolMagic,
		MPCost:         20,
		CooldownFrames: 90,
		CastRange:      12,
		WindupFrames:   12,
		TargetMode:     TargetModeEnemy,
	})
	reg.Register(SkillConfig{
		ID:             "focus",
		Name:           "专注",
		School:         SchoolMagic,
		MPCost:         10,
		CooldownFrames: 60,
		CastRange:      0,
		WindupFrames:   0,
		TargetMode:     TargetModeSelf,
	})
	return reg
}
