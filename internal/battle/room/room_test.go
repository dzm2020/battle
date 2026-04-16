package room

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"battle/internal/battle/attr"
	"battle/internal/battle/calc"
	"battle/internal/battle/entity"
)

func TestManager_CreateDuplicate(t *testing.T) {
	m := NewManager()
	if _, err := m.Create("r1", 2); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Create("r1", 2); err != ErrRoomExists {
		t.Fatalf("got %v", err)
	}
}

func TestRoom_JoinFull(t *testing.T) {
	r := newRoom("x", 1)
	_ = r.Join("a", entity.New("e1", 1, attr.Base{Level: 1}))
	if err := r.Join("b", entity.New("e2", 1, attr.Base{Level: 1})); err != ErrRoomFull {
		t.Fatalf("got %v", err)
	}
}

func TestRoom_StartBattle_Cancel(t *testing.T) {
	r := newRoom("x", 4)
	_ = r.Join("p1", entity.New("e1", 1, attr.Base{Level: 1}))
	ctx, cancel := context.WithCancel(context.Background())
	if err := r.StartBattle(ctx, calc.DefaultCalculator{}); err != nil {
		t.Fatal(err)
	}
	cancel()
	if err := r.Settle(); err != nil {
		t.Fatal(err)
	}
	if r.Phase() != PhaseSettled {
		t.Fatalf("phase %s", r.Phase())
	}
	r.Shutdown()
	if r.Phase() != PhaseClosed {
		t.Fatalf("phase %s", r.Phase())
	}
}

func TestRoom_ConcurrentJoin(t *testing.T) {
	r := newRoom("x", 8)
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			pid := fmt.Sprintf("p%d", i)
			_ = r.Join(pid, entity.New(pid, 1, attr.Base{Level: 1}))
		}(i)
	}
	wg.Wait()
	if r.PlayerCount() != 8 {
		t.Fatalf("count %d", r.PlayerCount())
	}
}
