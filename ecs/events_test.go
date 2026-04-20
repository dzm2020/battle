package ecs

import "testing"

func TestEvents_EntityAndComponent(t *testing.T) {
	w := NewWorld(4)

	var created, destroyed, added, removed int
	cancelC := w.Subscribe(EventEntityCreated, func(e Event) {
		if e.Kind != EventEntityCreated || e.Entity == 0 {
			t.Errorf("bad created: %+v", e)
		}
		created++
	})
	cancelD := w.Subscribe(EventEntityDestroyed, func(e Event) { destroyed++ })
	cancelA := w.Subscribe(EventComponentAdded, func(e Event) { added++ })
	cancelR := w.Subscribe(EventComponentRemoved, func(e Event) { removed++ })
	defer cancelC()
	defer cancelD()
	defer cancelA()
	defer cancelR()

	e := w.CreateEntity()
	if created != 1 {
		t.Fatalf("created events: want 1 got %d", created)
	}

	w.AddComponent(e, &Position{X: 1, Y: 2})
	if added != 1 {
		t.Fatalf("added events: want 1 got %d", added)
	}

	w.AddComponent(e, &Position{X: 9, Y: 9})
	if added != 1 {
		t.Fatalf("duplicate add should not emit: added=%d", added)
	}

	w.RemoveComponent(e, &Position{})
	if removed != 1 {
		t.Fatalf("removed events: want 1 got %d", removed)
	}

	w.RemoveEntity(e)
	if destroyed != 1 {
		t.Fatalf("destroyed events: want 1 got %d", destroyed)
	}
}

func TestEvents_SubscribeCancel(t *testing.T) {
	w := NewWorld(2)
	var n int
	cancel := w.Subscribe(EventEntityCreated, func(e Event) { n++ })
	e := w.CreateEntity()
	if n != 1 {
		t.Fatalf("want 1 got %d", n)
	}
	cancel()
	w.RemoveEntity(e)
	_ = w.CreateEntity()
	if n != 1 {
		t.Fatalf("after cancel should not fire: n=%d", n)
	}
}
