package skill

// CatalogConfig 技能模板表：按技能 ID 索引。
type CatalogConfig struct {
	byID map[uint32]SkillConfig
}

func NewCatalogConfig() *CatalogConfig {
	return &CatalogConfig{byID: make(map[uint32]SkillConfig)}
}

func (c *CatalogConfig) Register(sk SkillConfig) {
	c.byID[sk.ID] = sk
}

func (c *CatalogConfig) Get(id uint32) (SkillConfig, bool) {
	sk, ok := c.byID[id]
	return sk, ok
}
