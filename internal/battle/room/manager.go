package room

import (
	"sync"
	"sync/atomic"
)

var manager = NewManager()

func GetManager() *Manager {
	return manager
}

// Manager 全局房间表：多房间隔离；与单个 Room 的内锁分层，避免把整张 map 锁进 Room 逻辑。
type Manager struct {
	mu    sync.RWMutex
	rooms map[uint64]IRoom
	id    atomic.Uint64
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[uint64]IRoom),
	}
}
func (m *Manager) NextID() uint64 {
	roomId := m.id.Add(1)
	return roomId
}

func (m *Manager) Add(r IRoom) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.rooms[r.ID()] = r
}

// Get 读侧查找。
func (m *Manager) Get(id uint64) (IRoom, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rooms[id]
	return r, ok
}

// Remove 仅从管理器删除引用；不调用 Room.Shutdown，适合已自行 Shutdown 的场景。
func (m *Manager) Remove(id uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.rooms, id)
}

// Destroy 撤房并移出管理表（幂等：房间不存在则忽略）。
func (m *Manager) Destroy(id uint64) {
	m.mu.RLock()
	r, ok := m.rooms[id]
	m.mu.RUnlock()
	if !ok {
		return
	}
	r.Shutdown()
	m.Remove(id)
}

// Count 当前房间数量（测试/监控）。
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.rooms)
}
