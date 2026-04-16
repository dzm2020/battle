package tick

import (
	"testing"

	"battle/internal/battle/clock"
)

func TestLoop_StepOrder(t *testing.T) {
	clk := clock.New(60)
	loop := NewLoop(clk)
	var seq []int
	loop.Add(FuncSubscriber(func(c *clock.Clock) { seq = append(seq, 1) }))
	loop.Add(FuncSubscriber(func(c *clock.Clock) { seq = append(seq, 2) }))
	loop.Step()
	if clk.Frame() != 1 {
		t.Fatalf("frame %d", clk.Frame())
	}
	if len(seq) != 2 || seq[0] != 1 || seq[1] != 2 {
		t.Fatalf("seq %v", seq)
	}
}
