package tick

import (
	"context"
	"sync"
	"time"

	"battle/internal/battle/clock"
)

// Loop 固定帧战斗主循环：驱动 Clock Advance，并按注册顺序通知订阅者。
// 不负责具体战斗逻辑，与 timer / entity / room 解耦。
type Loop struct {
	clk  *clock.Clock
	mu   sync.RWMutex
	subs []Subscriber
}

func NewLoop(clk *clock.Clock) *Loop {
	return &Loop{clk: clk}
}

func (l *Loop) Clock() *clock.Clock { return l.clk }

// Add 追加订阅者；应在 Run/Step 前完成注册（战斗服通常在房间创建时组装）。
func (l *Loop) Add(s Subscriber) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.subs = append(l.subs, s)
}

// Step 推进一帧并通知订阅者；供单测与确定性演示使用，不睡眠。
func (l *Loop) Step() {
	l.clk.Advance()
	l.mu.RLock()
	subs := append([]Subscriber(nil), l.subs...)
	l.mu.RUnlock()
	for _, s := range subs {
		s.OnTick(l.clk)
	}
}

// Run 按 FrameDuration 对齐真实时间阻塞运行，直到 ctx 取消。
func (l *Loop) Run(ctx context.Context) error {
	ticker := time.NewTicker(l.clk.FrameDuration())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			l.Step()
		}
	}
}
