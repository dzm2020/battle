package room

import (
	"context"
	"sync"

	"battle/ecs"
	"battle/internal/battle/clock"
	"battle/internal/battle/component"
	"battle/internal/battle/skill"
	"battle/internal/battle/system"
	"battle/internal/battle/tick"
)

// Room 单局战斗隔离单元：独立 [ecs.World]、阶段、Clock/Loop。
// 生命周期阶段迁移集中在 phase_fsm.go（transitionPhase / advancePhaseLocked）。
// 不依赖网络层；Gateway 只应持有 roomID 并转调 Manager/Room API。
// 流程：大厅用 [Room.World] 创建实体并 [Join]；[StartBattle] 注册战斗系统并启动 tick；[Settle] 停循环；[Shutdown] 清场。
type Room struct {
	id         string
	maxPlayers int
	tps        int

	mu sync.RWMutex
	// 房间阶段
	phase Phase
	// 房间对象
	players map[string]ecs.Entity
	// ecs系统
	world *ecs.World
	// 逻辑帧驱动
	clk  *clock.Clock
	loop *tick.Loop

	cancel context.CancelFunc
	runWG  sync.WaitGroup
}

func newRoom(id string, maxPlayers int) *Room {
	if maxPlayers <= 0 {
		maxPlayers = 4
	}
	w := ecs.NewWorld(int32(maxPlayers + 8))
	component.RegisterCombatTypesWorld(w)
	return &Room{
		id:         id,
		maxPlayers: maxPlayers,
		tps:        60,
		phase:      PhaseLobby,
		players:    make(map[string]ecs.Entity),
		world:      w,
	}
}

func (r *Room) ID() string { return r.id }

func (r *Room) Phase() Phase {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.phase
}

func (r *Room) MaxPlayers() int { return r.maxPlayers }

func (r *Room) PlayerCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.players)
}

// World 该房间独占的 ECS 世界；大厅阶段即可 CreateEntity / 挂组件。
func (r *Room) World() *ecs.World {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.world
}

// Join 仅在 Lobby 阶段允许；sessionPlayerID 为连接侧玩家会话标识（非实体 ID）。
// e 必须为本房间 World 内仍存活的实体。
func (r *Room) Join(sessionPlayerID string, e ecs.Entity) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.phase == PhaseClosed {
		return ErrRoomClosed
	}
	if r.phase != PhaseLobby {
		return ErrWrongPhase
	}
	if _, ok := r.players[sessionPlayerID]; ok {
		return ErrDuplicatePlayer
	}
	if len(r.players) >= r.maxPlayers {
		return ErrRoomFull
	}
	if r.world == nil || !r.world.EntityExists(e) {
		return ErrInvalidEntity
	}
	for _, existing := range r.players {
		if existing == e {
			return ErrDuplicateEntity
		}
	}
	r.players[sessionPlayerID] = e
	return nil
}

// Leave 移除玩家并销毁对应实体；仅在 Lobby 开放。
func (r *Room) Leave(sessionPlayerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.phase == PhaseClosed {
		return ErrRoomClosed
	}
	if r.phase != PhaseLobby {
		return ErrWrongPhase
	}
	e, ok := r.players[sessionPlayerID]
	if !ok {
		return ErrPlayerNotFound
	}
	delete(r.players, sessionPlayerID)
	if r.world != nil {
		r.world.RemoveEntity(e)
	}
	return nil
}

// StartBattle 注册战斗管线、挂载 tick→World.Update 并启动独立循环协程；ctx 用于上层整体关服/撤房时取消。
// skillConfig 可为 nil，将使用空技能表；Buff 定义来自全局 [config.Tab]。
func (r *Room) StartBattle(ctx context.Context, skillConfig *skill.CatalogConfig) error {
	if ctx == nil {
		ctx = context.Background()
	}

	r.mu.Lock()
	if r.phase == PhaseClosed {
		r.mu.Unlock()
		return ErrRoomClosed
	}
	if r.phase != PhaseLobby {
		r.mu.Unlock()
		return ErrWrongPhase
	}
	if len(r.players) == 0 {
		r.mu.Unlock()
		return ErrNoPlayers
	}
	if r.cancel != nil {
		r.mu.Unlock()
		return ErrWrongPhase
	}

	if err := advancePhaseLocked(r, phaseEvStartBattle); err != nil {
		r.mu.Unlock()
		return err
	}
	r.clk = clock.New(r.tps)
	r.loop = tick.NewLoop(r.clk)
	loopCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	loop := r.loop
	w := r.world
	r.mu.Unlock()

	system.AddCombatSystems(w, skillConfig)

	dt := 1.0 / float64(r.clk.TPS())
	loop.Add(tick.FuncSubscriber(func(_ *clock.Clock) {
		w.Update(dt)
	}))

	r.mu.Lock()
	if err := advancePhaseLocked(r, phaseEvBattleLive); err != nil {
		r.mu.Unlock()
		return err
	}

	r.runWG.Add(1)
	go func() {
		defer r.runWG.Done()
		_ = loop.Run(loopCtx)
	}()
	r.mu.Unlock()
	return nil
}

// Settle 结束战斗循环并进入结算阶段；后续可读快照再 Destroy。
func (r *Room) Settle() error {
	var cancel context.CancelFunc
	r.mu.Lock()
	if r.phase != PhaseFighting {
		r.mu.Unlock()
		return ErrWrongPhase
	}
	cancel = r.cancel
	if err := advancePhaseLocked(r, phaseEvSettle); err != nil {
		r.mu.Unlock()
		return err
	}
	r.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	r.runWG.Wait()

	r.mu.Lock()
	r.cancel = nil
	r.loop = nil
	r.clk = nil
	r.mu.Unlock()
	return nil
}

// Shutdown 撤房：取消循环（若仍在跑）、清空世界上全部实体与玩家表；可从任意阶段调用（幂等）。
func (r *Room) Shutdown() {
	var cancel context.CancelFunc
	r.mu.Lock()
	if r.phase == PhaseClosed {
		r.mu.Unlock()
		return
	}
	if r.cancel != nil {
		cancel = r.cancel
	}
	w := r.world
	r.mu.Unlock()

	if cancel != nil {
		cancel()
		r.runWG.Wait()
	}

	if w != nil {
		w.RemoveAllEntities()
	}

	r.mu.Lock()
	if err := advancePhaseLocked(r, phaseEvShutdown); err != nil {
		r.mu.Unlock()
		return
	}
	r.players = make(map[string]ecs.Entity)
	r.cancel = nil
	r.loop = nil
	r.clk = nil
	r.mu.Unlock()
}

// Loop 返回当前战斗循环（仅 Fighting 阶段有效；单测可用 [tick.Loop.Step]）。
func (r *Room) Loop() *tick.Loop {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.loop
}

// SnapshotPlayers 返回 session → 实体 ID 的拷贝（用于广播/调试）。
func (r *Room) SnapshotPlayers() map[string]ecs.Entity {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]ecs.Entity, len(r.players))
	for k, v := range r.players {
		out[k] = v
	}
	return out
}
