package room

import (
	"context"
	"sync"

	"battle/internal/battle/calc"
	"battle/internal/battle/clock"
	"battle/internal/battle/entity"
	"battle/internal/battle/skill"
	"battle/internal/battle/tick"
)

// Room 单局战斗隔离单元：玩家实体、阶段、独立 Clock/Loop。
// 不依赖网络层；Gateway 只应持有 roomID 并转调 Manager/Room API。
type Room struct {
	id        string
	maxPlayers int

	mu      sync.RWMutex
	phase   Phase
	players map[string]*entity.Entity

	clk    *clock.Clock
	loop   *tick.Loop
	cancel context.CancelFunc
	runWG  sync.WaitGroup

	// skillSys 可选；SetSkillSystem 仅在 Lobby 阶段注入，StartBattle 时注册为 tick 订阅者。
	skillSys *skill.System
}

func newRoom(id string, maxPlayers int) *Room {
	if maxPlayers <= 0 {
		maxPlayers = 4
	}
	return &Room{
		id:         id,
		maxPlayers: maxPlayers,
		phase:      PhaseLobby,
		players:    make(map[string]*entity.Entity),
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

// Join 仅在 Lobby 阶段允许；sessionPlayerID 为连接侧玩家会话标识（非实体 ID）。
func (r *Room) Join(sessionPlayerID string, e *entity.Entity) error {
	if e == nil {
		return ErrInvalidEntity
	}
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
	r.players[sessionPlayerID] = e
	return nil
}

// SetSkillSystem 绑定技能子系统；必须在 Lobby 阶段调用，且每房间独立实例以防状态串线。
func (r *Room) SetSkillSystem(s *skill.System) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.phase != PhaseLobby {
		return ErrWrongPhase
	}
	r.skillSys = s
	return nil
}

// TryCastSkill 由网关/逻辑服转调；内部快照当前逻辑帧与战斗阶段。
// 注意：高并发下建议改为「入队 + Loop 内 TryCast」，此处为第 5 天可读性优先的实现。
func (r *Room) TryCastSkill(sessionPlayerID, skillID string, target *entity.Entity) skill.Result {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.skillSys == nil {
		return skill.Result{Reason: skill.RejectSkillSystemOff, Stage: skill.StageRejected}
	}
	if r.phase != PhaseFighting {
		return skill.Result{Reason: skill.RejectNotFighting, Stage: skill.StageRejected}
	}
	e := r.players[sessionPlayerID]
	var fr uint64
	if r.clk != nil {
		fr = r.clk.Frame()
	}
	if e == nil {
		return skill.Result{Reason: skill.RejectCasterMissing, Stage: skill.StageRejected}
	}
	return r.skillSys.TryCast(skill.CastInput{
		Frame:        fr,
		BattleActive: true,
		Caster:       e,
		Target:       target,
		SkillID:      skillID,
	})
}

// Leave 移除玩家；仅在 Lobby 开放（战斗中离开涉及掉线逻辑，放到后续天）。
func (r *Room) Leave(sessionPlayerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.phase == PhaseClosed {
		return ErrRoomClosed
	}
	if r.phase != PhaseLobby {
		return ErrWrongPhase
	}
	delete(r.players, sessionPlayerID)
	return nil
}

// StartBattle 初始化实体战斗状态并启动独立 tick 循环协程；ctx 用于上层整体关服/撤房时取消。
func (r *Room) StartBattle(ctx context.Context, cal calc.Calculator) error {
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

	r.phase = PhasePreBattle
	r.clk = clock.New(60)
	r.loop = tick.NewLoop(r.clk)
	loopCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	loop := r.loop
	ents := make([]*entity.Entity, 0, len(r.players))
	for _, e := range r.players {
		ents = append(ents, e)
	}
	r.mu.Unlock()

	for _, e := range ents {
		e.InitBattle(cal)
	}

	r.mu.Lock()
	r.phase = PhaseFighting
	if r.skillSys != nil {
		r.skillSys.ResetForBattle()
		loop.Add(r.skillSys)
	}
	r.runWG.Add(1)
	go func() {
		defer r.runWG.Done()
		_ = loop.Run(loopCtx)
	}()
	r.mu.Unlock()
	return nil
}

// Settle 结束战斗循环并进入结算阶段；后续可由结算服务读快照再 Destroy。
func (r *Room) Settle() error {
	var cancel context.CancelFunc
	r.mu.Lock()
	if r.phase != PhaseFighting {
		r.mu.Unlock()
		return ErrWrongPhase
	}
	cancel = r.cancel
	r.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	r.runWG.Wait()

	r.mu.Lock()
	if r.skillSys != nil {
		r.skillSys.ResetForBattle()
	}
	r.phase = PhaseSettled
	r.cancel = nil
	r.loop = nil
	r.clk = nil
	r.mu.Unlock()
	return nil
}

// Shutdown 撤房：取消循环（若仍在跑）、等待退出、清空引用；可从任意阶段调用（幂等）。
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
	r.mu.Unlock()

	if cancel != nil {
		cancel()
		r.runWG.Wait()
	}

	r.mu.Lock()
	if r.skillSys != nil {
		r.skillSys.ResetForBattle()
	}
	r.phase = PhaseClosed
	r.players = make(map[string]*entity.Entity)
	r.cancel = nil
	r.loop = nil
	r.clk = nil
	r.mu.Unlock()
}

// Loop 返回当前战斗循环（仅 Fighting 阶段有效，供调试或 Step 扩展；多数逻辑应通过 Subscriber 接入）。
func (r *Room) Loop() *tick.Loop {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.loop
}
