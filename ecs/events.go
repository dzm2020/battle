package ecs

import "sync"

// EventKind 世界内广播的事件类别。
type EventKind uint8

const (
	EventEntityCreated EventKind = iota + 1
	EventEntityDestroyed
	EventComponentAdded
	EventComponentRemoved
	// EventDamageApplied HealthSystem 扣减生命后派发；Entity 为受击者，Attacker 为来源（可为 0）；IntPayload 为结算伤害量。
	EventDamageApplied
	// EventDeath 实体因生命耗尽等处决逻辑即将移除或已标记死亡时派发。
	EventDeath
	// EventDamageMissed 伤害经命中/闪避判定未命中；Entity 为受击者，Attacker 为攻击来源。
	EventDamageMissed
	// EventHealApplied HealSystem 增加生命后派发；Entity 为受疗者，Attacker 作为治疗来源字段复用（Source）；IntPayload 为治疗量。
	EventHealApplied
	// EventBattleEnd BattleEndSystem 判定对局结束时派发；IntPayload 为获胜方 [component.Team].Side（0–255）；若全员阵亡（含同归于尽）则为 -1 表示平局。
	EventBattleEnd
)

// Event 单次广播的载荷；ComponentID / Component 仅在组件相关事件中有效。
type Event struct {
	Kind        EventKind
	Entity      Entity // 主实体（如受击者、死亡者）
	Attacker    Entity // 可选：伤害/治疗来源；无则 0
	ComponentID uint8
	Component   Component
	IntPayload  int // 数值载荷（伤害量、治疗量等）
}

// Subscribe 订阅指定类别事件；返回 cancel，调用后停止接收（可在任意 goroutine 调用 cancel，但 emit 与业务更新通常应在同一线程）。
func (w *World) Subscribe(kind EventKind, fn func(Event)) (cancel func()) {
	return w.events.add(kind, fn)
}

func newEventBus() *eventBus {
	return &eventBus{
		subs: make(map[EventKind]map[uint64]func(Event)),
	}
}

type eventBus struct {
	mu     sync.Mutex
	nextID uint64
	subs   map[EventKind]map[uint64]func(Event)
}

func (b *eventBus) add(kind EventKind, fn func(Event)) (cancel func()) {
	if fn == nil {
		return func() {}
	}
	b.mu.Lock()
	id := b.nextID
	b.nextID++
	if b.subs[kind] == nil {
		b.subs[kind] = make(map[uint64]func(Event))
	}
	b.subs[kind][id] = fn
	b.mu.Unlock()
	return func() {
		b.mu.Lock()
		if m, ok := b.subs[kind]; ok {
			delete(m, id)
			if len(m) == 0 {
				delete(b.subs, kind)
			}
		}
		b.mu.Unlock()
	}
}

func (b *eventBus) emit(e Event) {
	b.mu.Lock()
	list := make([]func(Event), 0, 8)
	if m := b.subs[e.Kind]; m != nil {
		for _, fn := range m {
			list = append(list, fn)
		}
	}
	b.mu.Unlock()
	for _, fn := range list {
		fn(e)
	}
}
