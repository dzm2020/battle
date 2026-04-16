package timer

import "testing"

func TestManager_OneShot(t *testing.T) {
	m := NewManager()
	m.AddOneShot(3, 7)

	if ev := m.ProcessFrame(2); len(ev) != 0 {
		t.Fatalf("frame2: %+v", ev)
	}
	ev := m.ProcessFrame(3)
	if len(ev) != 1 || ev[0].Tag != 7 {
		t.Fatalf("frame3: %+v", ev)
	}
	if ev := m.ProcessFrame(4); len(ev) != 0 {
		t.Fatalf("frame4 duplicate: %+v", ev)
	}
}

func TestManager_Repeat(t *testing.T) {
	m := NewManager()
	m.AddRepeat(2, 2, 1)

	var tags []Tag
	for f := uint64(1); f <= 6; f++ {
		for _, e := range m.ProcessFrame(f) {
			tags = append(tags, e.Tag)
		}
	}
	if len(tags) != 3 {
		t.Fatalf("got %v", tags)
	}
}

func TestManager_Cancel(t *testing.T) {
	m := NewManager()
	h := m.AddOneShot(5, 9)
	m.Cancel(h)
	if ev := m.ProcessFrame(10); len(ev) != 0 {
		t.Fatalf("cancelled fired: %+v", ev)
	}
}
