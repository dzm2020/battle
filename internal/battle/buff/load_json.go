package buff

import "encoding/json"

// LoadDescriptorsFromJSON 从 JSON 数组解析多条 [Descriptor] 并注册到 reg，
// 同 DefID 后写入覆盖先注册项。字段须与 struct 的 json 标签一致（枚举为数值）。
func LoadDescriptorsFromJSON(data []byte, reg *DefinitionRegistry) error {
	var defs []Descriptor
	if err := json.Unmarshal(data, &defs); err != nil {
		return err
	}
	for i := range defs {
		reg.Register(defs[i])
	}
	return nil
}
