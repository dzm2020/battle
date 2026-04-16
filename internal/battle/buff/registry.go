package buff

import (
	"sync"

	"battle/internal/battle/control"
)

// Registry Buff 配表：技能 applier、GM、关卡脚本共用。
type Registry struct {
	mu sync.RWMutex
	m  map[string]BuffConfig
}

func NewRegistry() *Registry {
	return &Registry{m: make(map[string]BuffConfig)}
}

func (r *Registry) Register(c BuffConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[c.ID] = c
}

func (r *Registry) Get(id string) (BuffConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.m[id]
	return c, ok
}

// DemoRegistry 眩晕 / 减速 / 增伤 / 毒伤 演示配表。
func DemoRegistry() *Registry {
	r := NewRegistry()
	for _, c := range []BuffConfig{
		{
			ID: "demo_stun", Name: "眩晕", Kind: KindStun,
			DurationFrames: 120, Control: control.FlagStunned,
			StackPolicy: StackReplace,
		},
		{
			ID: "demo_slow", Name: "减速", Kind: KindSlow,
			DurationFrames: 180, SlowMoveMul: 0.55,
			StackPolicy: StackRefresh,
		},
		{
			ID: "demo_amp", Name: "增伤", Kind: KindDamageAmp,
			DurationFrames: 90, OutDamageMul: 1.25, MaxStacks: 3,
			StackPolicy: StackLayer,
		},
		{
			ID: "demo_poison", Name: "毒", Kind: KindDot,
			DurationFrames: 300, TickIntervalFrames: 60, TickDeltaHP: -8,
			StackPolicy: StackRefresh,
		},
		{
			ID: "demo_strong", Name: "强攻", Kind: KindStatATK,
			DurationFrames: 120, StatATKFlat: 15, MaxStacks: 5,
			StackPolicy: StackLayer,
		},
		{
			ID: "demo_instant_heal", Name: "瞬疗", Kind: KindInstantHeal,
			TickDeltaHP: 30, // 正数治疗量
			StackPolicy: StackReplace,
		},
	} {
		r.Register(c)
	}
	return r
}
