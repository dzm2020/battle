package clock

import "time"

// Clock 逻辑战斗时钟：以固定 TPS 推进的帧计数器，与墙上时钟解耦。
// 战斗逻辑只认 Frame / LogicalMs，避免混用 time.Now 导致回放与单测不稳定。
type Clock struct {
	frame uint64
	tps   int
}

// New 创建时钟；常见战斗 TPS 为 60。
func New(tps int) *Clock {
	if tps <= 0 {
		tps = 60
	}
	return &Clock{tps: tps}
}

func (c *Clock) TPS() int { return c.tps }

// Frame 当前逻辑帧，从 1 开始计数（第 1 次 Advance 后为 1）。
func (c *Clock) Frame() uint64 { return c.frame }

// Advance 推进一帧；仅应由 tick.Loop 或单测 Step 调用，避免多处自增。
func (c *Clock) Advance() { c.frame++ }

// LogicalMs 已经运行了多久
func (c *Clock) LogicalMs() int64 {
	if c.tps <= 0 {
		return 0
	}
	return int64(c.frame) * 1000 / int64(c.tps)
}

// FrameDuration 单帧墙钟时长，供 Loop 与真实时间对齐。
func (c *Clock) FrameDuration() time.Duration {
	if c.tps <= 0 {
		return time.Second / 60
	}
	return time.Second / time.Duration(c.tps)
}
