package ecs

import "sync"

// EventKind 世界内广播的事件类别。
type EventKind uint8

const (
	EventEntityCreated EventKind = iota + 1
	EventEntityDestroyed
	EventComponentAdded
	EventComponentRemoved
	// EventDamageApplied 已由结算系统写入生命值后派发；IntPayload 为本次结算伤害量。
	EventDamageApplied
	// EventDeath 实体因生命耗尽等处决逻辑即将移除或已标记死亡时派发。
	EventDeath
)

// Event 单次广播的载荷；ComponentID / Component 仅在组件相关事件中有效。
type Event struct {
	Kind        EventKind
	Entity      Entity
	ComponentID uint8
	Component   Component
	IntPayload  int // 用于 EventDamageApplied 等；其它事件可为 0。
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
