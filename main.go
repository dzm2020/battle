package main

import (
	"context"
	"fmt"
	"log"

	"battle/internal/battle/attr"
	"battle/internal/battle/calc"
	"battle/internal/battle/entity"
	"battle/internal/battle/room"
	"battle/internal/battle/skill"
)

func main() {
	m := room.NewManager()
	r, err := m.Create("demo", 2)
	if err != nil {
		log.Fatal(err)
	}

	caster := entity.New("hero", 1, attr.Base{Level: 3, STR: 8, AGI: 5, INT: 4, VIT: 6})
	target := entity.New("mob", 2, attr.Base{Level: 2, STR: 4, AGI: 3, INT: 2, VIT: 5})
	target.Pos.X = 1

	for _, id := range []string{"strike", "hammer_stun", "ice_slow", "focus"} {
		caster.GrantSkill(id)
	}

	_ = r.Join("p1", caster)
	_ = r.Join("p2", target)

	reg := skill.MustDemoRegistry()
	sys := skill.NewSystem(reg, skill.BattleApplier{})
	_ = r.SetSkillSystem(sys)

	if err := r.StartBattle(context.Background(), calc.DefaultCalculator{}); err != nil {
		log.Fatal(err)
	}

	res := r.TryCastSkill("p1", "hammer_stun", target)
	fmt.Println("hammer_stun:", res.OK, res.Reason.String(), res.Stage)

	l := r.Loop()
	for i := 0; i < 3; i++ {
		l.Step()
	}
	if target.Control.HasStun() {
		fmt.Println("target stunned (tick ok)")
	}

	_ = r.Settle()
	r.Shutdown()
	fmt.Println("day07 buff demo done")
}
