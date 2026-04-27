package test

import (
	"battle/internal/battle/config"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	config.Load("./battle_config")

	tab := config.Tab
	if tab == nil {
		t.Fatal("Tab is nil")
	}
	if len(tab.BuffConfigConfigByID) == 0 {
		t.Fatal("Buff.config empty")
	}
	if _, ok := tab.BuffConfigConfigByID[900]; !ok {
		t.Fatal("missing buff 900")
	}
	if len(tab.SkillConfigByID) == 0 {
		t.Fatal("Skill.config empty")
	}
	if _, ok := tab.SkillConfigByID[1]; !ok {
		t.Fatal("missing skill 1")
	}
	if _, ok := tab.SkillEffectConfigByID[10]; !ok {
		t.Fatal("missing skill effect 10")
	}
}
