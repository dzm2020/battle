package timer

import (
	"sync"
)

// Tag 业务自定义标签，用于区分技能 CD 定时、Buff 到期、延时任务等。
type Tag int32

// Handle 定时器句柄，用于取消。
type Handle uint64

// Event 单次到期事件。
type Event struct {
	ID  Handle
	Tag Tag
}

type entry struct {
	id          Handle
	expireFrame uint64
	repeatEvery uint64 // 0 表示一次性；>0 表示每次触发后顺延 repeatEvery 帧
	tag         Tag
	removed     bool
}

// Manager 帧驱动定时器：仅依赖「当前逻辑帧」，不依赖 wall clock，与战斗循环天然对齐。
type Manager struct {
	mu      sync.Mutex
	next    Handle
	entries []*entry
}

func NewManager() *Manager {
	return &Manager{}
}

// AddOneShot 在 expireFrame 这一帧及之后首次 ProcessFrame 时触发（currentFrame >= expireFrame）。
// 若 expireFrame 为 0，则下一帧即可触发（通常应传入「当前帧 + 延迟」）。
func (m *Manager) AddOneShot(expireFrame uint64, tag Tag) Handle {
	return m.add(expireFrame, 0, tag)
}

// AddRepeat 在 firstExpire 首次触发，之后每 repeatEvery 帧再触发；repeatEvery 必须 > 0。
func (m *Manager) AddRepeat(firstExpire, repeatEvery uint64, tag Tag) Handle {
	if repeatEvery == 0 {
		repeatEvery = 1
	}
	return m.add(firstExpire, repeatEvery, tag)
}

func (m *Manager) add(expireFrame, repeatEvery uint64, tag Tag) Handle {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.next++
	h := m.next
	m.entries = append(m.entries, &entry{
		id:          h,
		expireFrame: expireFrame,
		repeatEvery: repeatEvery,
		tag:         tag,
	})
	return h
}

// Cancel 标记移除；本帧若已入队触发则仍可能派发（与常见游戏「本帧已结算」一致，可按需改为提前扫描）。
func (m *Manager) Cancel(id Handle) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, e := range m.entries {
		if e.id == id {
			e.removed = true
			break
		}
	}
}

// ProcessFrame 在当前逻辑帧下结算到期定时器，返回本次触发的事件列表。
func (m *Manager) ProcessFrame(currentFrame uint64) []Event {
	m.mu.Lock()
	defer m.mu.Unlock()

	var out []Event
	alive := m.entries[:0]
	for _, e := range m.entries {
		if e.removed {
			continue
		}
		if e.expireFrame > currentFrame {
			alive = append(alive, e)
			continue
		}
		out = append(out, Event{ID: e.id, Tag: e.tag})
		if e.repeatEvery > 0 {
			e.expireFrame = currentFrame + e.repeatEvery
			alive = append(alive, e)
		}
	}
	m.entries = alive
	return out
}
