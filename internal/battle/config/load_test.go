package config

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testBattleConfigDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", "test", "battle_config"))
}

func TestLoad_ok(t *testing.T) {
	dir := testBattleConfigDir(t)
	if err := Load(dir); err != nil {
		t.Fatal(err)
	}
	if Tab == nil || Tab.SkillConfigByID == nil {
		t.Fatal("Tab not populated")
	}
}

func TestLoad_missingDir(t *testing.T) {
	err := Load(filepath.Join(testBattleConfigDir(t), "nonexistent_subdir_xyz"))
	if err == nil {
		t.Fatal("expected error for missing config dir")
	}
}
