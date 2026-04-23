package room

import (
	"context"
	"testing"
	"time"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/skill"
)

func TestRoom_JoinLeaveAndDuplicateRules(t *testing.T) {
	m := NewManager()
	r, err := m.Create("r1", 8)
	if err != nil {
		t.Fatal(err)
	}
	w := r.World()
	if err := r.Join("ghost", ecs.Entity(1<<60)); err != ErrInvalidEntity {
		t.Fatalf("want ErrInvalidEntity for unknown entity, got %v", err)
	}
	a := w.CreateEntity()
	b := w.CreateEntity()
	if err := r.Join("alice", a); err != nil {
		t.Fatal(err)
	}
	if err := r.Join("bob", a); err != ErrDuplicateEntity {
		t.Fatalf("want ErrDuplicateEntity, got %v", err)
	}
	if err := r.Join("bob", b); err != nil {
		t.Fatal(err)
	}
	if err := r.Leave("alice"); err != nil {
		t.Fatal(err)
	}
	a2 := w.CreateEntity()
	if err := r.Join("alice", a2); err != nil {
		t.Fatal(err)
	}
	r.Shutdown()
	if r.Phase() != PhaseClosed {
		t.Fatalf("phase want Closed got %v", r.Phase())
	}
}

func TestRoom_StartBattleSettleShutdown(t *testing.T) {
	m := NewManager()
	r, err := m.Create("r2", 4)
	if err != nil {
		t.Fatal(err)
	}
	e := r.World().CreateEntity()
	if err := r.Join("p1", e); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	if err := r.StartBattle(ctx, skill.NewCatalogConfig()); err != nil {
		t.Fatal(err)
	}
	if r.Phase() != PhaseFighting {
		t.Fatalf("phase want Fighting got %v", r.Phase())
	}
	cancel()
	time.Sleep(80 * time.Millisecond)

	if err := r.Settle(); err != nil {
		t.Fatal(err)
	}
	if r.Phase() != PhaseSettled {
		t.Fatalf("phase want Settled got %v", r.Phase())
	}
	r.Shutdown()
	if r.Phase() != PhaseClosed {
		t.Fatalf("phase want Closed got %v", r.Phase())
	}
}

func TestRoom_StartBattleNeedsPlayers(t *testing.T) {
	m := NewManager()
	r, _ := m.Create("empty", 2)
	err := r.StartBattle(context.Background(), nil)
	if err != ErrNoPlayers {
		t.Fatalf("want ErrNoPlayers, got %v", err)
	}
	r.Shutdown()
}

func TestRoom_TickUpdatesCombatWorld(t *testing.T) {
	m := NewManager()
	r, _ := m.Create("r3", 2)
	w := r.World()
	caster := w.CreateEntity()
	foe := w.CreateEntity()
	w.AddComponent(caster, &component.Team{Side: 1})
	w.AddComponent(foe, &component.Team{Side: 2})
	w.AddComponent(foe, &component.Health{Current: 100, Max: 100})
	w.AddComponent(foe, &component.Attributes{Values: map[string]int{config.AttrHp: 100}})
	skillCfg := skill.NewCatalogConfig()
	skillCfg.Register(skill.SkillConfig{
		ID:         99,
		Resource:   skill.ResourceNone,
		Cost:       0,
		Scope:      skill.TargetScopeSingle,
		Camp:       skill.CampEnemy,
		CastFrames: 0,
		Effects: []skill.EffectConfig{
			{Kind: skill.EffectDamage, Amount: 20, DamageType: component.DamagePhysical},
		},
	})
	// Same wiring as [Room.StartBattle]: systems + dt + tick subscriber.
	if err := joinAndStartBattleForTest(t, r, caster, skillCfg); err != nil {
		t.Fatal(err)
	}
	defer r.Shutdown()

	loop := r.Loop()
	if loop == nil {
		t.Fatal("loop is nil")
	}
	w.AddComponent(caster, &component.SkillUser{GrantedSkillIDs: []uint32{99}})
	w.AddComponent(caster, &component.CastIntent{SkillID: 99, Target: foe})

	loop.Step()
	h, ok := w.GetComponent(foe, &component.Health{})
	if !ok {
		t.Fatal("foe health missing")
	}
	if hp := h.(*component.Health); hp.Current >= 100 {
		t.Fatalf("expected damage, hp=%d", hp.Current)
	}
}

// joinAndStartBattleForTest mirrors [Room.StartBattle] but uses ctx cancel + wait instead of sleeping in tests.
func joinAndStartBattleForTest(t *testing.T, r *Room, ent ecs.Entity, skillCfg *skill.CatalogConfig) error {
	t.Helper()
	if err := r.Join("solo", ent); err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	err := r.StartBattle(ctx, skillCfg)
	if err != nil {
		return err
	}
	cancel()
	time.Sleep(80 * time.Millisecond)
	// battle loop exited; phase still Fighting — caller uses [Loop] inline [Step].
	return nil
}

func TestManager_Destroy(t *testing.T) {
	m := NewManager()
	r, _ := m.Create("gone", 2)
	_ = r.World().CreateEntity()
	if m.Count() != 1 {
		t.Fatalf("count want 1")
	}
	m.Destroy("gone")
	if m.Count() != 0 {
		t.Fatalf("count want 0 after Destroy")
	}
}
