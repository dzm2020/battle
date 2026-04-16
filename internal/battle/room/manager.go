package room

import "sync"

// Manager 全局房间表：多房间隔离；与单个 Room 的内锁分层，避免把整张 map 锁进 Room 逻辑。
type Manager struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
	}
}

// Create 创建空房间；id 重复返回 ErrRoomExists。
func (m *Manager) Create(id string, maxPlayers int) (*Room, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.rooms[id]; ok {
		return nil, ErrRoomExists
	}
	r := newRoom(id, maxPlayers)
	m.rooms[id] = r
	return r, nil
}

// Get 读侧查找。
func (m *Manager) Get(id string) (*Room, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.rooms[id]
	return r, ok
}

// Remove 仅从管理器删除引用；不调用 Room.Shutdown，适合已自行 Shutdown 的场景。
func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.rooms, id)
}

// Destroy 撤房并移出管理表（幂等：房间不存在则忽略）。
func (m *Manager) Destroy(id string) {
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
