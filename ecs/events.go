package ecs

import "sync"

// EventKind 世界内广播的事件类别。
type EventKind uint8

const (
	EventKindEntityCreated EventKind = iota + 1
	EventKindEntityDestroyed
	EventKindComponentAdded
	EventKindComponentRemoved
)

type EventEntityCreated struct {
	E Entity
}

type EventEntityDestroyed struct {
	E Entity
}

type EventEntityComponentAdded struct {
	E      Entity
	CompID uint8
	Comp   Component
}
type EventEntityComponentRemoved struct {
	E      Entity
	CompID uint8
	Comp   Component
}

// FirstUserEventKind 业务自定义事件 Kind 建议从此值起递增分配，避免与框架内置 Kind 冲突。
const FirstUserEventKind EventKind = 32

// Event 单次广播：Kind 为订阅键；Payload 为任意类型，由业务定义并在订阅端断言。
type Event struct {
	Kind    EventKind
	Payload any
}

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
