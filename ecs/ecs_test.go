package ecs

import "testing"

type testPosition struct {
	X, Y float64
}

func (*testPosition) Component() {}

type testVelocity struct {
	DX, DY float64
}

func (*testVelocity) Component() {}

type testHealth struct {
	Current, Max int
}

func (*testHealth) Component() {}

type movementSystem struct {
	query *Query2[*testPosition, *testVelocity]
}

func (s *movementSystem) Initialize(w *World) {
	s.query = NewQuery2[*testPosition, *testVelocity](w)
}

func (s *movementSystem) Update(dt float64) {
	s.query.ForEach(func(_ Entity, pos *testPosition, vel *testVelocity) {
		pos.X += vel.DX * dt
		pos.Y += vel.DY * dt
	})
}

type healSystem struct {
	query *Query[*testHealth]
}

func (s *healSystem) Initialize(w *World) {
	s.query = NewQuery[*testHealth](w)
}

func (s *healSystem) Update(dt float64) {
	heal := int(10 * dt)
	s.query.ForEach(func(_ Entity, h *testHealth) {
		if h.Current >= h.Max {
			return
		}
		h.Current += heal
		if h.Current > h.Max {
			h.Current = h.Max
		}
	})
}

func TestWorld_SystemsAndQuery(t *testing.T) {
	world := NewWorld(8)
	world.Registry().Register(&testPosition{})
	world.Registry().Register(&testVelocity{})
	world.Registry().Register(&testHealth{})

	player := world.CreateEntity()
	world.AddComponent(player, &testPosition{X: 100, Y: 200})
	world.AddComponent(player, &testVelocity{DX: 10, DY: 5})
	world.AddComponent(player, &testHealth{Current: 80, Max: 100})

	world.AddSystem(&movementSystem{})
	world.AddSystem(&healSystem{})

	const dt = 0.1
	for i := 0; i < 3; i++ {
		world.Update(dt)
	}

	pos, ok := world.GetComponent(player, &testPosition{})
	if !ok {
		t.Fatal("missing position")
	}
	p := pos.(*testPosition)
	wantX := 100 + 10*dt*3
	wantY := 200 + 5*dt*3
	if p.X != wantX || p.Y != wantY {
		t.Fatalf("position: got (%.2f, %.2f) want (%.2f, %.2f)", p.X, p.Y, wantX, wantY)
	}

	hp, ok := world.GetComponent(player, &testHealth{})
	if !ok {
		t.Fatal("missing health")
	}
	h := hp.(*testHealth)
	if h.Current != 83 {
		t.Fatalf("health current: got %d want 83", h.Current)
	}

	entities := NewQuery[*testHealth](world).Collect()
	if len(entities) != 1 {
		t.Fatalf("query collect: got %d entities", len(entities))
	}
}

// 保留 demo 用组件名，供 events_test 使用。
type Position struct {
	X, Y float64
}

func (*Position) Component() {}
