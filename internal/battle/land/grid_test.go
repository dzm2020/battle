package land

import (
	"math"
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
	u1.Pos.X, u1.Pos.Z = 1, 1
	g.AddUnit(u1)

	u2 := &Unit{ID: 2}
	u2.Pos.X, u2.Pos.Z = 50, 50
	g.AddUnit(u2)

	near := g.GetNearbyUnits(1, 1, 10)
	if len(near) != 1 || near[0].ID != 1 {
		t.Fatalf("near (1,1) want unit 1, got %v", len(near))
	}

	g.UpdateUnit(u1, 48, 48)
	near = g.GetNearbyUnits(50, 50, 10)
	if len(near) != 2 {
		t.Fatalf("after move want 2 near (50,50), got %d", len(near))
	}

	g.RemoveUnit(u1)
	near = g.GetNearbyUnits(50, 50, 10)
	if len(near) != 1 || near[0].ID != 2 {
		t.Fatalf("after remove want only u2")
	}
}

func TestSpatialGrid_CellIndexClamp(t *testing.T) {
	g, err := NewSpatialGrid(0, 0, 10, 10, 2)
	if err != nil {
		t.Fatal(err)
	}
	u := &Unit{ID: 99}
	u.Pos.X, u.Pos.Z = -100, -100
	g.AddUnit(u)
	// 仍落在合法格内（被裁剪到角格）
	cx, cz := g.cellIndex(u.Pos.X, u.Pos.Z)
	if cx != 0 || cz != 0 {
		t.Fatalf("clamped corner want (0,0) got (%d,%d)", cx, cz)
	}
}

func TestGetNearbyUnits_Radius(t *testing.T) {
	g, _ := NewSpatialGrid(0, 0, 20, 20, 1)
	u := &Unit{ID: 1}
	u.Pos.X, u.Pos.Z = 10, 0
	g.AddUnit(u)
	list := g.GetNearbyUnits(0, 0, 9.9)
	if len(list) != 0 {
		t.Fatal("outside radius should be empty")
	}
	list = g.GetNearbyUnits(0, 0, 10.1)
	if len(list) != 1 {
		t.Fatalf("inside radius want 1, got %d", len(list))
	}
	if math.Abs(list[0].Pos.X-10) > 1e-9 {
		t.Fatal("wrong unit")
	}
}
