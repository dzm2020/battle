package skill

import (
	"testing"

	"battle/internal/battle/attr"
	"battle/internal/battle/calc"
	"battle/internal/battle/control"
	"battle/internal/battle/entity"
	"battle/internal/battle/geom"
)

func TestValidateCast_stunBlocks(t *testing.T) {
	cfg := SkillConfig{ID: "s", School: SchoolPhysical, TargetMode: TargetModeSelf}
	c := entity.New("c", 1, attr.Base{Level: 1})
	c.InitBattle(calc.DefaultCalculator{})
	c.GrantSkill("s")
	c.Control = control.FlagStunned
	in := CastInput{Frame: 1, BattleActive: true, Caster: c, SkillID: "s"}
	if r := ValidateCast(cfg, in); r != RejectStunned {
		t.Fatalf("got %v", r)
	}
}

func TestValidateCast_silenceBlocksMagic(t *testing.T) {
	cfg := SkillConfig{ID: "m", School: SchoolMagic, TargetMode: TargetModeSelf}
	c := entity.New("c", 1, attr.Base{Level: 1})
	c.InitBattle(calc.DefaultCalculator{})
	c.GrantSkill("m")
	c.Control = control.FlagSilenced
	in := CastInput{Frame: 1, BattleActive: true, Caster: c, SkillID: "m"}
	if r := ValidateCast(cfg, in); r != RejectSilenced {
		t.Fatalf("got %v", r)
	}
}

func TestValidateCast_distance(t *testing.T) {
	cfg := SkillConfig{
		ID: "x", School: SchoolPhysical, CastRange: 1, TargetMode: TargetModeEnemy,
		CooldownFrames: 1,
	}
	c := entity.New("c", 1, attr.Base{Level: 1})
	c.InitBattle(calc.DefaultCalculator{})
	c.GrantSkill("x")
	c.Pos = geom.Vec2{X: 0, Y: 0}
	c.Runtime.CurMP = 100
	tgt := entity.New("t", 2, attr.Base{Level: 1})
	tgt.InitBattle(calc.DefaultCalculator{})
	tgt.Pos = geom.Vec2{X: 5, Y: 0}
	in := CastInput{Frame: 1, BattleActive: true, Caster: c, Target: tgt, SkillID: "x"}
	if r := ValidateCast(cfg, in); r != RejectOutOfRange {
		t.Fatalf("got %v", r)
	}
}
