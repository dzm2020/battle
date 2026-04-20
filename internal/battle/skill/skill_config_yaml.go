package skill

import "gopkg.in/yaml.v3"

// LoadCatalogConfigFromYAML 从 YAML 数组解析多条 [SkillConfig] 并写入 [CatalogConfig]。
// 字段名、枚举数值与 JSON 一致（scope、camp、pickRule 等）。
func LoadCatalogConfigFromYAML(data []byte, catalogConfig *CatalogConfig) error {
	var defs []SkillConfig
	if err := yaml.Unmarshal(data, &defs); err != nil {
		return err
	}
	for i := range defs {
		catalogConfig.Register(defs[i])
	}
	return nil
}
