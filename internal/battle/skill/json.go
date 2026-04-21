package skill

import (
	"encoding/json"
	"fmt"
)

// LoadCatalogConfigFromJSON 从 JSON 数组加载技能列表（与 docs/skill-design 中示例格式一致）。
func LoadCatalogConfigFromJSON(data []byte, cat *CatalogConfig) error {
	if cat == nil {
		return fmt.Errorf("skill: nil catalog")
	}
	var raw []SkillConfig
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	for _, sk := range raw {
		cat.Register(sk)
	}
	return nil
}
