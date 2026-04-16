package cooldown

import "testing"

func TestBook_Trigger(t *testing.T) {
	b := NewBook()
	if !b.IsReady(1, "s1") {
		t.Fatal()
	}
	b.Trigger(10, "s1", 30)
	if b.IsReady(10, "s1") {
		t.Fatal()
	}
	if !b.IsReady(40, "s1") {
		t.Fatal()
	}
}
