package ecs

import "testing"

type testGrid struct {
	W, H int
}

func TestResource(t *testing.T) {
	w := NewWorld(8)
	var get Resource[testGrid]
	get = get.New(w)

	if get.Has() {
		t.Fatal("expected no resource")
	}
	grid := testGrid{W: 100, H: 200}
	get.Add(&grid)

	if !get.Has() {
		t.Fatal("expected resource")
	}
	got := get.Get()
	if got == nil || got.W != 100 || got.H != 200 {
		t.Fatalf("grid: got %+v", got)
	}

	get.Remove()
	if get.Has() {
		t.Fatal("expected resource removed")
	}
}

func TestAddResource_GetResource(t *testing.T) {
	w := NewWorld(4)
	grid := testGrid{W: 10, H: 20}
	AddResource(w, &grid)

	got := GetResource[testGrid](w)
	if got == nil || got.W != 10 {
		t.Fatalf("GetResource: got %+v", got)
	}
}

func TestResources_IDAPI(t *testing.T) {
	w := NewWorld(4)
	id := ResourceID[testGrid](w)
	res := &testGrid{W: 1, H: 2}
	w.Resources().Add(id, res)

	if !w.Resources().Has(id) {
		t.Fatal("Has should be true")
	}
	if w.Resources().Get(id) != res {
		t.Fatal("Get mismatch")
	}
}

func TestAddResource_panicsOnDuplicate(t *testing.T) {
	w := NewWorld(4)
	AddResource(w, &testGrid{})
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on duplicate AddResource")
		}
	}()
	AddResource(w, &testGrid{})
}
