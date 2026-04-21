package config

var Tab = new(Tables)

func Load() {
	Tab.load()
}

type Tables struct {
	AttributeConfigByID  map[int32]AttributeConfig
	BuffConfigConfigByID map[int32]*BuffConfig
	SkillConfig          *SkillConfig
	UnitConfig           *UnitConfig
}

// todo  完成配置加载
func (t *Tables) load() {

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
