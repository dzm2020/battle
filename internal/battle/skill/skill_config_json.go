package skill

import "encoding/json"

// LoadCatalogConfigFromJSON 从 JSON 数组解析多条 [SkillConfig] 并写入 [CatalogConfig]；同 ID 后写入覆盖先前项。
// 每条对象须含合法 scope、camp 等字段，见 skill_target_spec.go。
func LoadCatalogConfigFromJSON(data []byte, catalogConfig *CatalogConfig) error {
	var defs []SkillConfig
	if err := json.Unmarshal(data, &defs); err != nil {
		return err
	}
	for i := range defs {
		catalogConfig.Register(defs[i])
	}
	return nil
}
