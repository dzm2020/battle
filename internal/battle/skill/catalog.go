package skill

// CatalogConfig 旧版技能表占位；若项目仍通过 [system.AddCombatSystems] 传入，保留空实现即可。
type CatalogConfig struct{}

func NewCatalogConfig() *CatalogConfig {
	return &CatalogConfig{}
}
