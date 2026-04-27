package skill

import (
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

func TestApplySkillEffects_Damage(t *testing.T) {
	prevSkill := config.Tab.SkillConfigByID
	prevEff := config.Tab.SkillEffectConfigByID
	t.Cleanup(func() {
		config.Tab.SkillConfigByID = prevSkill
		config.Tab.SkillEffectConfigByID = prevEff
	})

	config.Tab.SkillConfigByID = map[int32]*config.SkillBaseConfig{
		1: {ID: 1, EffectIDs: []int{10}},
	}
	config.Tab.SkillEffectConfigByID = map[int32]*config.SkillEffectConfig{
		10: {
			EffectID:       10,
			EffectType:     config.EffectDamage,
			IntParams:      []int{42, int(component.DamageTrue)},
			TargetSelectID: 0,
		},
	}

	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)
	caster := w.CreateEntity()
	victim := w.CreateEntity()
	w.AddComponent(victim, &component.Health{Current: 100, Max: 100})
	w.AddComponent(victim, &component.Attributes{})

	ApplySkillEffects(w, caster, victim, 1)

	pd, ok := w.GetComponent(victim, &component.PendingDamage{})
	if !ok {
		t.Fatal("expected PendingDamage on victim")
	}
	got := pd.(*component.PendingDamage)
	if got.Amount != 42 || got.Type != component.DamageTrue || got.Source != caster {
		t.Fatalf("PendingDamage = %+v", got)
	}
}

func TestApplySkillEffects_AddBuff(t *testing.T) {
	prevSkill := config.Tab.SkillConfigByID
	prevEff := config.Tab.SkillEffectConfigByID
	prevBuff := config.Tab.BuffConfigConfigByID
	t.Cleanup(func() {
		config.Tab.SkillConfigByID = prevSkill
		config.Tab.SkillEffectConfigByID = prevEff
		config.Tab.BuffConfigConfigByID = prevBuff
	})

	const buffDefID uint32 = 900
	config.Tab.BuffConfigConfigByID = map[int32]*config.BuffConfig{
		int32(buffDefID): {
			ID:            buffDefID,
			MaxStack:      1,
			StackBehavior: config.BuffStackAdd,
			DurationFrame: 3,
			EffectType:    config.BufferEffectStatChange,
			ParamsString:  []string{config.AttrArmor},
			Params:        []float64{1},
		},
	}
	config.Tab.SkillConfigByID = map[int32]*config.SkillBaseConfig{
		2: {ID: 2, EffectIDs: []int{20}},
	}
	config.Tab.SkillEffectConfigByID = map[int32]*config.SkillEffectConfig{
		20: {
			EffectID:       20,
			EffectType:     config.EffectAddBuff,
			IntParams:      []int{int(buffDefID)},
			TargetSelectID: 0,
		},
	}

	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)
	caster := w.CreateEntity()
	target := w.CreateEntity()
	w.AddComponent(target, &component.Health{Current: 10, Max: 10})

	ApplySkillEffects(w, caster, target, 2)

	bl, ok := w.GetComponent(target, &component.BuffList{})
	if !ok {
		t.Fatal("expected BuffList on target")
	}
	if n := len(bl.(*component.BuffList).Buffs); n != 1 {
		t.Fatalf("buff count want 1 got %d", n)
	}
}
