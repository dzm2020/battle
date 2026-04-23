package system

import (
	"strings"
	"testing"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/skill"
)

func TestSkill_LoadJSONAndInstantDamage(t *testing.T) {
	const raw = `[
	  {"id":10,"resource":1,"cost":30,"cooldownFrames":2,"scope":2,"camp":0,"castFrames":0,
	   "effects":[{"kind":0,"amount":40,"damageType":2}]}
	]`
	skillConfig := skill.NewCatalogConfig()
	if err := skill.LoadCatalogConfigFromJSON([]byte(strings.TrimSpace(raw)), skillConfig); err != nil {
		t.Fatal(err)
	}
	w := ecs.NewWorld(16)
	component.RegisterCombatTypesWorld(w)
	AddCombatSystems(w, skillConfig)

	caster := w.CreateEntity()
	foe := w.CreateEntity()
	w.AddComponent(caster, &component.Team{Side: 1})
	w.AddComponent(foe, &component.Team{Side: 2})
	w.AddComponent(caster, &component.SkillUser{
		Mana:              100,
		GrantedSkillIDs:   []uint32{10},
		CooldownRemaining: nil,
	})
	w.AddComponent(foe, &component.Health{Current: 100, Max: 100})
	w.AddComponent(foe, &component.Attributes{})
	w.AddComponent(caster, &component.CastIntent{SkillID: 10, Target: foe})

	w.Update(0)
	if su, _ := w.GetComponent(caster, &component.SkillUser{}); su.(*component.SkillUser).Mana != 70 {
		t.Fatalf("mana want 70")
	}
	h := getHP(t, w, foe)
	if h.Current >= 100 {
		t.Fatalf("should take damage, hp=%d", h.Current)
	}
	if _, ok := w.GetComponent(caster, &component.CastIntent{}); ok {
		t.Fatal("intent should be consumed")
	}
	if left := suCO(t, w, caster, 10); left != 2 {
		t.Fatalf("cd want 2 frames left got %d", left)
	}
}

func TestSkill_AoETwoTargets(t *testing.T) {
	skillConfig := skill.NewCatalogConfig()
	skillConfig.Register(skill.SkillConfig{
		ID:             20,
		Resource:       skill.ResourceNone,
		Cost:           0,
		CooldownFrames: 0,
		Scope:          skill.TargetScopeMulti,
		Camp:           skill.CampEnemy,
		CastFrames:     0,
		Effects: []skill.EffectConfig{
			{Kind: skill.EffectDamage, Amount: 10, DamageType: component.DamageTrue},
		},
	})
	w := ecs.NewWorld(16)
	component.RegisterCombatTypesWorld(w)
	AddCombatSystems(w, skillConfig)
	caster := w.CreateEntity()
	a := w.CreateEntity()
	b := w.CreateEntity()
	w.AddComponent(caster, &component.Team{Side: 1})
	w.AddComponent(a, &component.Team{Side: 2})
	w.AddComponent(b, &component.Team{Side: 2})
	w.AddComponent(caster, &component.SkillUser{GrantedSkillIDs: []uint32{20}})
	w.AddComponent(a, &component.Health{Current: 50, Max: 50})
	w.AddComponent(b, &component.Health{Current: 50, Max: 50})
	w.AddComponent(a, &component.Attributes{})
	w.AddComponent(b, &component.Attributes{})
	w.AddComponent(caster, &component.CastIntent{SkillID: 20})

	w.Update(0)
	for _, e := range []ecs.Entity{a, b} {
		h := getHP(t, w, e)
		if h.Current >= 50 {
			t.Fatalf("entity %v should lose hp", e)
		}
	}
}

func TestSkill_ChannelThenResolve(t *testing.T) {
	skillConfig := skill.NewCatalogConfig()
	skillConfig.Register(skill.SkillConfig{
		ID:             30,
		Resource:       skill.ResourceMana,
		Cost:           5,
		CooldownFrames: 0,
		Scope:          skill.TargetScopeSelf,
		Camp:           skill.CampEnemy,
		CastFrames:     2,
		Effects: []skill.EffectConfig{
			{Kind: skill.EffectHeal, Amount: 7},
		},
	})
	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)
	AddCombatSystems(w, skillConfig)
	caster := w.CreateEntity()
	w.AddComponent(caster, &component.Team{Side: 1})
	w.AddComponent(caster, &component.SkillUser{Mana: 50, GrantedSkillIDs: []uint32{30}})
	w.AddComponent(caster, &component.Health{Current: 10, Max: 100})
	w.AddComponent(caster, &component.CastIntent{SkillID: 30})

	w.Update(0)
	su := getSU(t, w, caster)
	if su.Mana != 45 {
		t.Fatalf("mana deducted at channel start want 45 got %d", su.Mana)
	}
	if _, ok := w.GetComponent(caster, &component.SkillCastState{}); !ok {
		t.Fatal("expect channel state")
	}
	w.Update(0)
	if _, ok := w.GetComponent(caster, &component.SkillCastState{}); !ok {
		t.Fatal("still channeling")
	}
	w.Update(0)
	if _, ok := w.GetComponent(caster, &component.SkillCastState{}); ok {
		t.Fatal("channel done")
	}
	h := getHP(t, w, caster)
	if h.Current != 17 {
		t.Fatalf("heal after channel want 17 got %d", h.Current)
	}
}

func TestSkill_ApplyBuffEffect(t *testing.T) {
	prev := config.Tab.BuffConfigConfigByID
	config.Tab.BuffConfigConfigByID = map[int32]*config.BuffConfig{
		900: {
			ID:             900,
			MaxStack:       1,
			StackBehavior:  config.BuffStackAdd,
			DurationFrame:  5,
			EffectType:     config.BufferEffectStatChange,
			ParamsString:   []string{config.AttrArmor},
			Params:         []float64{3},
		},
	}
	t.Cleanup(func() { config.Tab.BuffConfigConfigByID = prev })

	skillConfig := skill.NewCatalogConfig()
	skillConfig.Register(skill.SkillConfig{
		ID:         40,
		Resource:   skill.ResourceNone,
		Scope:      skill.TargetScopeSelf,
		Camp:       skill.CampEnemy,
		CastFrames: 0,
		Effects:    []skill.EffectConfig{{Kind: skill.EffectApplyBuff, BuffDefID: 900}},
	})
	w := ecs.NewWorld(8)
	component.RegisterCombatTypesWorld(w)
	AddCombatSystems(w, skillConfig)
	e := w.CreateEntity()
	w.AddComponent(e, &component.Team{Side: 1})
	w.AddComponent(e, &component.SkillUser{GrantedSkillIDs: []uint32{40}})
	w.AddComponent(e, &component.CastIntent{SkillID: 40})
	w.Update(0)
	bl, ok := w.GetComponent(e, &component.BuffList{})
	if !ok {
		t.Fatal("buff list")
	}
	if len(bl.(*component.BuffList).Buffs) != 1 {
		t.Fatal("expected one buff instance")
	}
}

func getHP(t *testing.T, w *ecs.World, e ecs.Entity) *component.Health {
	t.Helper()
	c, ok := w.GetComponent(e, &component.Health{})
	if !ok {
		t.Fatal("no health")
	}
	return c.(*component.Health)
}

func getSU(t *testing.T, w *ecs.World, e ecs.Entity) *component.SkillUser {
	t.Helper()
	c, ok := w.GetComponent(e, &component.SkillUser{})
	if !ok {
		t.Fatal("no SkillUser")
	}
	return c.(*component.SkillUser)
}

func suCO(t *testing.T, w *ecs.World, e ecs.Entity, id uint32) int {
	t.Helper()
	su := getSU(t, w, e)
	return su.CooldownRemaining[id]
}
