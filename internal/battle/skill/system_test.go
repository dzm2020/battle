package skill

import (
	"testing"

	"battle/internal/battle/attr"
	"battle/internal/battle/calc"
	"battle/internal/battle/clock"
	"battle/internal/battle/entity"
	"battle/internal/battle/geom"
	"battle/internal/battle/tick"
)

func TestSystem_WindupCommits(t *testing.T) {
	clk := clock.New(60)
	loop := tick.NewLoop(clk)
	reg := MustDemoRegistry()
	sys := NewSystem(reg, DefaultApplier{})
	loop.Add(sys)

	cal := calc.DefaultCalculator{}
	c := entity.New("c", 1, attr.Base{Level: 10, INT: 30})
	e := entity.New("e", 2, attr.Base{Level: 1})
	c.GrantSkill("fireball")
	c.InitBattle(cal)
	e.InitBattle(cal)
	c.Pos = geom.Vec2{X: 0, Y: 0}
	e.Pos = geom.Vec2{X: 1, Y: 0}
	mpBefore := c.Runtime.CurMP

	loop.Step()
	res := sys.TryCast(CastInput{Frame: clk.Frame(), BattleActive: true, Caster: c, Target: e, SkillID: "fireball"})
	if !res.OK || res.Stage != StageWindupScheduled {
		t.Fatalf("res %+v", res)
	}
	if c.Runtime.CurMP != mpBefore {
		t.Fatalf("mp should not change before windup: %d", c.Runtime.CurMP)
	}

	for clk.Frame() < res.WindupEndsAtFrame {
		loop.Step()
	}
	if c.Runtime.CurMP != mpBefore-20 {
		t.Fatalf("mp after windup want %d got %d", mpBefore-20, c.Runtime.CurMP)
	}
}
