package buff

import "fmt"

// DefinitionConfig 维护 DefID→[DescriptorConfig] 的查找表；[ApplyBuff] 与 [BuffSystem] 依赖同一张表解析效果。
// 通常在房间/战斗初始化时注入，同一场战斗内视为只追加、不中途改语义（覆盖同 ID 需自行约定）。
type DefinitionConfig struct {
	byID map[uint32]DescriptorConfig
}

// NewDefinitionConfig 创建空的 Buff 定义配置表。
func NewDefinitionConfig() *DefinitionConfig {
	return &DefinitionConfig{byID: make(map[uint32]DescriptorConfig)}
}

// Register 注册或覆盖某一 DefID；ID 为 0 会 panic。
func (r *DefinitionConfig) Register(d DescriptorConfig) {
	if d.ID == 0 {
		panic("buff: DescriptorConfig.ID must be non-zero")
	}
	r.byID[d.ID] = d
}

// Get 查询定义；未找到时 ok 为 false，[BuffSystem] 会丢弃持有该 DefID 的实例。
func (r *DefinitionConfig) Get(id uint32) (DescriptorConfig, bool) {
	d, ok := r.byID[id]
	return d, ok
}

// MustGet 与 Get 相同，未找到时 panic（仅建议在确信 ID 合法时使用）。
func (r *DefinitionConfig) MustGet(id uint32) DescriptorConfig {
	d, ok := r.Get(id)
	if !ok {
		panic(fmt.Sprintf("buff: unknown def id %d", id))
	}
	return d
}
