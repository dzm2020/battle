package skill

import "fmt"

// CatalogConfig 技能静态模板配置表（典型用法：战斗房间初始化时加载 JSON/YAML，整场战斗只读）。
type CatalogConfig struct {
	byID map[uint32]SkillConfig
}

// NewCatalogConfig 创建空的技能配置表。
func NewCatalogConfig() *CatalogConfig {
	return &CatalogConfig{byID: make(map[uint32]SkillConfig)}
}

// Register 注册或覆盖某一技能 ID；ID 为 0 将 panic。
func (c *CatalogConfig) Register(s SkillConfig) {
	if s.ID == 0 {
		panic("skill: SkillConfig.ID must be non-zero")
	}
	c.byID[s.ID] = s
}

// Get 查询技能配置；未找到时 ok 为 false。
func (c *CatalogConfig) Get(id uint32) (SkillConfig, bool) {
	v, ok := c.byID[id]
	return v, ok
}

// MustGet 未找到时 panic。
func (c *CatalogConfig) MustGet(id uint32) SkillConfig {
	v, ok := c.Get(id)
	if !ok {
		panic(fmt.Sprintf("skill: unknown id %d", id))
	}
	return v
}
