package land

import (
	"testing"
)

func TestNewSpatialGrid_Invalid(t *testing.T) {
	if _, err := NewSpatialGrid(0, 0, 10, 10, 0); err == nil {
		t.Fatal("want error for zero cellSize")
	}
	if _, err := NewSpatialGrid(10, 0, 5, 10, 5); err == nil {
		t.Fatal("want error for maxX<=minX")
	}
	if _, err := NewSpatialGrid(0, 10, 10, 5, 5); err == nil {
		t.Fatal("want error for maxZ<=minZ")
	}
}

func TestSpatialGrid_AddUpdateNearby(t *testing.T) {
	g, err := NewSpatialGrid(0, 0, 100, 100, 5)
	if err != nil {
		t.Fatal(err)
	}

	u1 := &Unit{ID: 1}
	if err := g.AddUnit(u1, 0, 0); err != nil {
		t.Fatal(err)
	}

	u2 := &Unit{ID: 2}
	if err := g.AddUnit(u2, 10, 10); err != nil {
		t.Fatal(err)
	}

	// (0,0) 与 (10,10) 格距平方 200；半径 14 → 196 < 200，仅 u1
	near := g.GetNearbyUnits(0, 0, 14)
	if len(near) != 1 || near[0].ID != 1 {
		t.Fatalf("near (0,0) r=14 want unit 1 only, got %d units", len(near))
	}

	near = g.GetNearbyUnits(0, 0, 15)
	if len(near) != 2 {
		t.Fatalf("near (0,0) r=15 want 2 units, got %d", len(near))
	}

	g.UpdateUnit(u1, 10, 10)
	near = g.GetNearbyUnits(10, 10, 1)
	if len(near) != 2 {
		t.Fatalf("after move want 2 at (10,10), got %d", len(near))
	}

	g.RemoveUnit(u1)
	near = g.GetNearbyUnits(10, 10, 0)
	if len(near) != 1 || near[0].ID != 2 {
		t.Fatalf("after remove want only u2")
	}
}

func TestGetNearbyUnits_Radius(t *testing.T) {
	g, _ := NewSpatialGrid(0, 0, 20, 20, 1)
	u := &Unit{ID: 1}
	_ = g.AddUnit(u, 10, 0)
	list := g.GetNearbyUnits(0, 0, 9)
	if len(list) != 0 {
		t.Fatal("outside radius should be empty")
	}
	list = g.GetNearbyUnits(0, 0, 10)
	if len(list) != 1 {
		t.Fatalf("inside radius want 1, got %d", len(list))
	}
	if list[0].cellX != 10 || list[0].cellZ != 0 {
		t.Fatal("wrong unit cell")
	}
}

func TestGetNearbyUnits_NilGrid(t *testing.T) {
	var g *Grid
	if out := g.GetNearbyUnits(0, 0, 5); out != nil {
		t.Fatal("nil grid want nil slice")
	}
}

func TestForEachCell_AscDesc(t *testing.T) {
	g, err := NewSpatialGrid(0, 0, 5, 5, 2)
	if err != nil {
		t.Fatal(err)
	}
	var asc [][2]int
	g.ForEachCellAsc(func(cx, cz int, cell *GridCell) {
		if cell == nil {
			t.Fatal("nil cell")
		}
		asc = append(asc, [2]int{cx, cz})
	})
	var desc [][2]int
	g.ForEachCellDesc(func(cx, cz int, cell *GridCell) {
		desc = append(desc, [2]int{cx, cz})
	})
	if len(asc) != len(desc) || len(asc) != g.Width()*g.Height() {
		t.Fatalf("visit count: asc=%d desc=%d want %d", len(asc), len(desc), g.Width()*g.Height())
	}
	for i := range asc {
		j := len(desc) - 1 - i
		if asc[i][0] != desc[j][0] || asc[i][1] != desc[j][1] {
			t.Fatalf("reverse mismatch i=%d asc=%v desc[j]=%v", i, asc[i], desc[j])
		}
	}
}

func TestFirstFreeCell_AscDesc(t *testing.T) {
	g, err := NewSpatialGrid(0, 0, 5.1, 5.1, 5)
	if err != nil {
		t.Fatal(err)
	}
	if g.Width() != 2 || g.Height() != 2 {
		t.Fatalf("want 2x2 grid, got %dx%d", g.Width(), g.Height())
	}
	x0, z0, ok := g.FirstFreeCellAsc()
	if !ok || x0 != 0 || z0 != 0 {
		t.Fatalf("empty grid asc want (0,0) ok got (%d,%d) ok=%v", x0, z0, ok)
	}
	x1, z1, ok := g.FirstFreeCellDesc()
	if !ok || x1 != 1 || z1 != 1 {
		t.Fatalf("empty grid desc want (1,1) ok got (%d,%d) ok=%v", x1, z1, ok)
	}

	u := &Unit{ID: 1}
	_ = g.AddUnit(u, 0, 0)
	x0, z0, ok = g.FirstFreeCellAsc()
	if !ok || x0 != 0 || z0 != 1 {
		t.Fatalf("one unit at (0,0) asc want (0,1) got (%d,%d) ok=%v", x0, z0, ok)
	}
	x1, z1, ok = g.FirstFreeCellDesc()
	if !ok || x1 != 1 || z1 != 1 {
		t.Fatalf("desc still want (1,1) got (%d,%d) ok=%v", x1, z1, ok)
	}

	for id, pos := range []struct{ cx, cz int }{
		{0, 1},
		{1, 0},
		{1, 1},
	} {
		u2 := &Unit{ID: uint64(10 + id)}
		_ = g.AddUnit(u2, pos.cx, pos.cz)
	}
	_, _, ok = g.FirstFreeCellAsc()
	if ok {
		t.Fatal("full grid asc want ok=false")
	}
	_, _, ok = g.FirstFreeCellDesc()
	if ok {
		t.Fatal("full grid desc want ok=false")
	}
}
