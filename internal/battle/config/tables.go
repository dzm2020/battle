package config

var Tab = new(Tables)

func Load() {
	Tab.load()
}

type Tables struct {
	AttributeConfigByID  map[int32]AttributeConfig
	BuffConfigConfigByID map[int32]*BuffConfig
	SkillConfigByID      map[int32]*SkillBaseConfig
	// SkillEffectConfigByID 技能效果行（[SkillEffectConfig.EffectID] → 配置）。
	SkillEffectConfigByID map[int32]*SkillEffectConfig
	// TargetSelectConfigByID 选目标规则（[TargetSelectConfig.ID] → 配置）。
	TargetSelectConfigByID map[int32]*TargetSelectConfig
	UnitConfig             *UnitConfig
}

// todo  完成配置加载
func (t *Tables) load() {
	if t.SkillConfigByID == nil {
		t.SkillConfigByID = make(map[int32]*SkillBaseConfig)
	}
	if t.SkillEffectConfigByID == nil {
		t.SkillEffectConfigByID = make(map[int32]*SkillEffectConfig)
	}
	if t.TargetSelectConfigByID == nil {
		t.TargetSelectConfigByID = make(map[int32]*TargetSelectConfig)
	}
}

// CombatBundle 预留：目录下若干 JSON/YAML 汇总后的结果（供校验器交叉引用）。
type CombatBundle struct{}

// LoadCombatBundleFromDir 读取目录战前配置；当前为占位实现。
func LoadCombatBundleFromDir(dir string) (*CombatBundle, error) {
	_ = dir
	return &CombatBundle{}, nil
}

// ValidateCombatBundle 交叉引用校验；当前无规则时返回 nil。
func ValidateCombatBundle(_ *CombatBundle) []error {
	return nil
}

func GetSkillConfigByID(id int32) *SkillBaseConfig {
	return Tab.SkillConfigByID[id]
}

func GetTargetSelectConfigByID(id int32) *TargetSelectConfig {
	return Tab.TargetSelectConfigByID[id]
}
