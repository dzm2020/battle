package buff

import "encoding/json"

// LoadDefinitionConfigFromJSON 从 JSON 数组解析多条 [DescriptorConfig] 并写入 [DefinitionConfig]；
// 同 DefID 后写入覆盖先前项；字段须与 struct 的 json 标签一致。
func LoadDefinitionConfigFromJSON(data []byte, config *DefinitionConfig) error {
	var defs []DescriptorConfig
	if err := json.Unmarshal(data, &defs); err != nil {
		return err
	}
	for i := range defs {
		config.Register(defs[i])
	}
	return nil
}
