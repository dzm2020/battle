package config

import (
	"encoding/json"
	"fmt"
	"os"
	path2 "path"
)

var Tab *Tables

// Load 从目录加载全部配表；失败返回错误且不修改 [Tab]。
func Load(path string) error {
	tab := new(Tables)
	if err := tab.load(path); err != nil {
		return err
	}
	Tab = tab
	return nil
}

// MustLoad 等同 [Load]，失败时 panic（仅启动/单测快捷路径）。
func MustLoad(path string) {
	if err := Load(path); err != nil {
		panic(err)
	}
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
	SkillEffectConfigByID  map[int32]*SkillEffectConfig
	TargetSelectConfigByID map[int32]*TargetSelectConfig
	UnitConfigByID         map[int32]*UnitConfig
	DungeonConfigByID      map[int32]*DungeonConfig
	MapConfigByID          map[int32]*MapConfig
}

func (t *Tables) load(path string) error {
	loaders := []struct {
		name string
		fn   func(string) error
	}{
		{"Attribute.json", t.loadAttributeConfig},
		{"Buff.json", t.loadBuffConfig},
		{"Skill.json", t.loadSkillConfig},
		{"SkillEffect.json", t.loadSkillEffectConfig},
		{"TargetSelect.json", t.loadTargetSelectConfig},
		{"Unit.json", t.loadUnitConfig},
		{"Dungeon.json", t.loadDungeonConfig},
		{"Map.json", t.loadMapConfig},
	}
	for _, l := range loaders {
		if err := l.fn(path); err != nil {
			return fmt.Errorf("config: load %s: %w", l.name, err)
		}
	}
	return nil
}

func (t *Tables) loadAttributeConfig(path string) error {
	if t.AttributeConfigByID == nil {
		t.AttributeConfigByID = make(map[int32]*AttributeConfig)
	}
	return ReadJSONFromFile(path2.Join(path, "Attribute.json"), &t.AttributeConfigByID)
}

func (t *Tables) loadBuffConfig(path string) error {
	if t.BuffConfigConfigByID == nil {
		t.BuffConfigConfigByID = make(map[int32]*BuffConfig)
	}
	return ReadJSONFromFile(path2.Join(path, "Buff.json"), &t.BuffConfigConfigByID)
}

func (t *Tables) loadSkillConfig(path string) error {
	if t.SkillConfigByID == nil {
		t.SkillConfigByID = make(map[int32]*SkillBaseConfig)
	}
	return ReadJSONFromFile(path2.Join(path, "Skill.json"), &t.SkillConfigByID)
}

func (t *Tables) loadSkillEffectConfig(path string) error {
	if t.SkillEffectConfigByID == nil {
		t.SkillEffectConfigByID = make(map[int32]*SkillEffectConfig)
	}
	return ReadJSONFromFile(path2.Join(path, "SkillEffect.json"), &t.SkillEffectConfigByID)
}

func (t *Tables) loadTargetSelectConfig(path string) error {
	if t.TargetSelectConfigByID == nil {
		t.TargetSelectConfigByID = make(map[int32]*TargetSelectConfig)
	}
	return ReadJSONFromFile(path2.Join(path, "TargetSelect.json"), &t.TargetSelectConfigByID)
}

func (t *Tables) loadUnitConfig(path string) error {
	if t.UnitConfigByID == nil {
		t.UnitConfigByID = make(map[int32]*UnitConfig)
	}
	return ReadJSONFromFile(path2.Join(path, "Unit.json"), &t.UnitConfigByID)
}

func (t *Tables) loadDungeonConfig(path string) error {
	if t.DungeonConfigByID == nil {
		t.DungeonConfigByID = make(map[int32]*DungeonConfig)
	}
	return ReadJSONFromFile(path2.Join(path, "Dungeon.json"), &t.DungeonConfigByID)
}

func (t *Tables) loadMapConfig(path string) error {
	if t.MapConfigByID == nil {
		t.MapConfigByID = make(map[int32]*MapConfig)
	}
	return ReadJSONFromFile(path2.Join(path, "Map.json"), &t.MapConfigByID)
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

func GetUnitConfigByID(id int32) *UnitConfig {
	return Tab.UnitConfigByID[id]
}

func GetAttributeConfigByID(id int32) *AttributeConfig {
	return Tab.AttributeConfigByID[id]
}

func GetDungeonConfigByID(id int32) *DungeonConfig {
	if Tab == nil || Tab.DungeonConfigByID == nil {
		return nil
	}
	return Tab.DungeonConfigByID[id]
}

func GetMapConfigByID(id int32) *MapConfig {
	if Tab == nil || Tab.MapConfigByID == nil {
		return nil
	}
	return Tab.MapConfigByID[id]
}
