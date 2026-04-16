package clock

import "testing"

func TestClock_LogicalMs(t *testing.T) {
	c := New(60)
	c.Advance()
	c.Advance()
	if c.Frame() != 2 {
		t.Fatalf("frame %d", c.Frame())
	}
	if c.LogicalMs() != 33 { // 2000/60 = 33 整除
		t.Fatalf("ms %d", c.LogicalMs())
	}
}
