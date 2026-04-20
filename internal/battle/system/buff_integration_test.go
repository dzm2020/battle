package system

import (
	"testing"

	"battle/ecs"
	"battle/internal/battle/action"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
	"battle/internal/battle/control"
)

func setupCombatWorld(buffConfig *buff.DefinitionConfig) *ecs.World {
	w := ecs.NewWorld(64)
	component.RegisterCombatTypesWorld(w)
	AddCombatSystems(w, buffConfig, nil)
	return w
}

func TestBuff_ManyIndependentBuffs(t *testing.T) {
	buffConfig := buff.NewDefinitionConfig()
	for id := uint32(1); id <= 20; id++ {
		id := id
		buffConfig.Register(buff.DescriptorConfig{
			ID:             id,
			MaxStacks:      1,
			Policy:         buff.StackIndependent,
			DurationFrames: 5,
			Effects: []buff.EffectConfig{
				{Kind: buff.EffectStatMod, ArmorDeltaPerStack: 1},
			},
		})
	}
	w := setupCombatWorld(buffConfig)
	e := w.CreateEntity()
	for id := uint32(1); id <= 20; id++ {
		if !buff.ApplyBuff(w, buffConfig, e, id) {
			t.Fatalf("ApplyBuff %d", id)
		}
	}
	bl, ok := w.GetComponent(e, &component.BuffList{})
	if !ok {
		t.Fatal("no BuffList")
	}
	if n := len(bl.(*component.BuffList).Buffs); n != 20 {
		t.Fatalf("buff count want 20 got %d", n)
	}
	w.Update(0)
	sm, ok := w.GetComponent(e, &component.StatModifiers{})
	if !ok {
		t.Fatal("no StatModifiers")
	}
	if sm.(*component.StatModifiers).ArmorDelta != 20 {
		t.Fatalf("armor delta want 20 got %d", sm.(*component.StatModifiers).ArmorDelta)
	}
}

func TestBuff_DoTPerFrame(t *testing.T) {
	buffConfig := buff.NewDefinitionConfig()
	buffConfig.Register(buff.DescriptorConfig{
		ID:             101,
		MaxStacks:      1,
		Policy:         buff.StackMerge,
		DurationFrames: 10,
		Effects: []buff.EffectConfig{
			{
				Kind:               buff.EffectDoT,
				DamagePerTick:      20,
				DamageType:         component.DamageMagical,
				TickIntervalFrames: 1,
			},
		},
	})
	w := setupCombatWorld(buffConfig)
	e := w.CreateEntity()
	w.AddComponent(e, &component.Health{Current: 1000, Max: 1000})
	w.AddComponent(e, &component.Attributes{MagicResist: 100})
	buff.ApplyBuff(w, buffConfig, e, 101)

	w.Update(0)
	h := getHealth(t, w, e)
	if h.Current != 990 {
		t.Fatalf("after 1 frame DoT want 990 got %d", h.Current)
	}
}

func TestBuff_StunBlocksAct(t *testing.T) {
	buffConfig := buff.NewDefinitionConfig()
	buffConfig.Register(buff.DescriptorConfig{
		ID:             202,
		MaxStacks:      1,
		Policy:         buff.StackMerge,
		DurationFrames: 30,
		Effects: []buff.EffectConfig{
			{Kind: buff.EffectControl, Control: control.FlagStunned},
		},
	})
	w := setupCombatWorld(buffConfig)
	e := w.CreateEntity()
	if !action.CanAct(w, e) {
		t.Fatal("should act without buff")
	}
	buff.ApplyBuff(w, buffConfig, e, 202)
	w.Update(0)
	if action.CanAct(w, e) {
		t.Fatal("stun should block action")
	}
}

func TestBuff_StackRefresh(t *testing.T) {
	buffConfig := buff.NewDefinitionConfig()
	buffConfig.Register(buff.DescriptorConfig{
		ID:             303,
		MaxStacks:      3,
		Policy:         buff.StackRefresh,
		DurationFrames: 3,
		Effects: []buff.EffectConfig{
			{Kind: buff.EffectStatMod, ArmorDeltaPerStack: 5},
		},
	})
	w := setupCombatWorld(buffConfig)
	e := w.CreateEntity()
	buff.ApplyBuff(w, buffConfig, e, 303)
	bl1 := getBuffList(t, w, e)
	if bl1.Buffs[0].Stacks != 1 {
		t.Fatalf("stacks want 1 got %d", bl1.Buffs[0].Stacks)
	}
	buff.ApplyBuff(w, buffConfig, e, 303)
	bl2 := getBuffList(t, w, e)
	if bl2.Buffs[0].Stacks != 1 {
		t.Fatalf("refresh should not add stack, want 1 got %d", bl2.Buffs[0].Stacks)
	}
	if bl2.Buffs[0].FramesLeft != 3 {
		t.Fatalf("duration refreshed, want 3 got %d", bl2.Buffs[0].FramesLeft)
	}
}

func TestBuff_StackMerge(t *testing.T) {
	buffConfig := buff.NewDefinitionConfig()
	buffConfig.Register(buff.DescriptorConfig{
		ID:             404,
		MaxStacks:      3,
		Policy:         buff.StackMerge,
		DurationFrames: 10,
		Effects: []buff.EffectConfig{
			{Kind: buff.EffectStatMod, ArmorDeltaPerStack: 2},
		},
	})
	w := setupCombatWorld(buffConfig)
	e := w.CreateEntity()
	buff.ApplyBuff(w, buffConfig, e, 404)
	buff.ApplyBuff(w, buffConfig, e, 404)
	bl := getBuffList(t, w, e)
	if bl.Buffs[0].Stacks != 2 {
		t.Fatalf("merge stacks want 2 got %d", bl.Buffs[0].Stacks)
	}
	w.Update(0)
	sm, _ := w.GetComponent(e, &component.StatModifiers{})
	if sm.(*component.StatModifiers).ArmorDelta != 4 {
		t.Fatalf("armor want 4 got %d", sm.(*component.StatModifiers).ArmorDelta)
	}
}

func getHealth(t *testing.T, w *ecs.World, e ecs.Entity) *component.Health {
	t.Helper()
	c, ok := w.GetComponent(e, &component.Health{})
	if !ok {
		t.Fatal("no Health")
	}
	return c.(*component.Health)
}

func getBuffList(t *testing.T, w *ecs.World, e ecs.Entity) *component.BuffList {
	t.Helper()
	c, ok := w.GetComponent(e, &component.BuffList{})
	if !ok {
		t.Fatal("no BuffList")
	}
	return c.(*component.BuffList)
}
