package config

import (
	"encoding/json"
	"os"
	path2 "path"
)

var Tab *Tables

func Load(path string) {
	tab := new(Tables)
	tab.load(path)
	Tab = tab
}
func ReadJSONFromFile(filePath string, v interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

type Tables struct {
	AttributeConfigByID    map[int32]*AttributeConfig
	BuffConfigConfigByID   map[int32]*BuffConfig
	SkillConfigByID        map[int32]*SkillBaseConfig
	SkillEffectConfigByID  map[int32]*SkillEffectConfig  // 技能效果（[EffectID] → 行）
	TargetSelectConfigByID map[int32]*TargetSelectConfig // 选目标规则
	UnitConfigByID         map[string]*UnitConfig        // 单位模板（[UnitConfig.ID] → 行）
}

func (t *Tables) load(path string) {
	t.loadAttributeConfig(path)
	t.loadBuffConfig(path)
	t.loadSkillConfig(path)
	t.loadSkillEffectConfig(path)
	t.loadTargetSelectConfig(path)
	t.loadUnitConfig(path)
}

func (t *Tables) loadAttributeConfig(path string) {
	if t.AttributeConfigByID == nil {
		t.AttributeConfigByID = make(map[int32]*AttributeConfig)
	}
	err := ReadJSONFromFile(path2.Join(path, "Attribute.json"), &t.AttributeConfigByID)
	if err != nil {
		panic(err)
	}
}

func (t *Tables) loadBuffConfig(path string) {
	if t.BuffConfigConfigByID == nil {
		t.BuffConfigConfigByID = make(map[int32]*BuffConfig)
	}
	err := ReadJSONFromFile(path2.Join(path, "Buff.json"), &t.BuffConfigConfigByID)
	if err != nil {
		panic(err)
	}
}

func (t *Tables) loadSkillConfig(path string) {
	if t.SkillConfigByID == nil {
		t.SkillConfigByID = make(map[int32]*SkillBaseConfig)
	}
	err := ReadJSONFromFile(path2.Join(path, "Skill.json"), &t.SkillConfigByID)
	if err != nil {
		panic(err)
	}
}

func (t *Tables) loadSkillEffectConfig(path string) {
	if t.SkillEffectConfigByID == nil {
		t.SkillEffectConfigByID = make(map[int32]*SkillEffectConfig)
	}
	err := ReadJSONFromFile(path2.Join(path, "SkillEffect.json"), &t.SkillEffectConfigByID)
	if err != nil {
		panic(err)
	}
}

func (t *Tables) loadTargetSelectConfig(path string) {
	if t.TargetSelectConfigByID == nil {
		t.TargetSelectConfigByID = make(map[int32]*TargetSelectConfig)
	}
	err := ReadJSONFromFile(path2.Join(path, "TargetSelect.json"), &t.TargetSelectConfigByID)
	if err != nil {
		panic(err)
	}
}

func (t *Tables) loadUnitConfig(path string) {
	if t.UnitConfigByID == nil {
		t.UnitConfigByID = make(map[string]*UnitConfig)
	}
	err := ReadJSONFromFile(path2.Join(path, "Unit.json"), &t.UnitConfigByID)
	if err != nil {
		panic(err)
	}
}

func GetSkillConfigByID(id int32) *SkillBaseConfig {
	return Tab.SkillConfigByID[id]
}
func GetTargetSelectConfigByID(id int32) *TargetSelectConfig {
	return Tab.TargetSelectConfigByID[id]
}
func GetSkillEffectConfigByID(id int32) *SkillEffectConfig {
	return Tab.SkillEffectConfigByID[id]
}
