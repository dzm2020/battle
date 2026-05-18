package resource

import "time"

type TPS struct {
	TPS   int    // 每秒更新多少帧
	Frame uint64 // 当前推进到多少帧
}

const defaultTPS = 60

// EffectiveTPS 返回有效 TPS；小于等于 0 时回退为 [defaultTPS]。
func (t *TPS) EffectiveTPS() int {
	if t == nil || t.TPS <= 0 {
		return defaultTPS
	}
	return t.TPS
}

// DeltaTime 返回当前 TPS 下单帧逻辑时长（秒），供 [ecs.World.Update] 使用。
func (t *TPS) DeltaTime() float64 {
	return 1.0 / float64(t.EffectiveTPS())
}

// FrameDuration 返回当前 TPS 下单帧墙钟时长，供房间循环 ticker 使用。
func (t *TPS) FrameDuration() time.Duration {
	return time.Second / time.Duration(t.EffectiveTPS())
}
