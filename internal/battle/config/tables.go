package config

var Tab = new(Tables)

func Load() {
	Tab.load()
}

type Tables struct {
	// AttributeConfigByID 属性配置表，键为 [AttributeConfig].ID，加载后注入；供 [NewAttributesFromAttributeConfigIDs] 等解析。
	AttributeConfigByID map[int32]AttributeConfig
	BuffConfig          *BuffConfig
	SkillConfig         *SkillConfig
	UnitConfig          *UnitConfig
}

// todo  完成配置加载
func (t *Tables) load() {

}
