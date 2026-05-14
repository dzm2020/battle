package room

import (
	"context"
	"sync"

	"battle/ecs"
	"battle/internal/battle/clock"
	"battle/internal/battle/component"
	"battle/internal/battle/land"
	"battle/internal/battle/tick"
)

// New 分配房间 id、创建 ECS 世界并完成 [component.Register]；尚未装配地图与单位。
// 由 [battle/internal/battle/room_builder] 在创建流程中调用后再 [SetGrid] 与刷实体。
func New() *Room {
	r := &Room{
		id:    GetManager().NextID(),
		tps:   60,
		phase: PhaseLobby,
		world: ecs.NewWorld(100),
	}
	r.world.SetContext(r)
	component.Register(r.world)
	return r
}

// Room 单局战斗隔离单元：独立 [ecs.World]、阶段字段 [Room.phase]、Clock/Loop。
// 不依赖网络层；Gateway 只应持有 roomID 并转调 Manager/Room API。
// 流程：大厅用 [Room.World] 创建实体并 [Join]；[StartBattle] 注册战斗系统并启动 tick；[Settle] 停循环；[Shutdown] 清场。
type Room struct {
	id    uint64
	tps   int
	phase Phase
	// 地图
	grid *land.Grid
	// ecs系统
	world *ecs.World
	// 逻辑帧驱动
	clk    *clock.Clock
	loop   *tick.Loop
	cancel context.CancelFunc
	runWG  sync.WaitGroup
}

// phaseIs 判断当前是否处于指定阶段；调用方必须已持有 [Room.mu] 读锁或写锁。
func (r *Room) phaseIs(p Phase) bool {
	return r.phase == p
}

func (r *Room) ID() uint64 { return r.id }

// SetGrid 设置空间网格；传入 [land.Grid]，内部会绑定本房间的 [Room.World]。
func (r *Room) SetGrid(base *land.Grid) {
	r.grid = base
}

func (r *Room) Grid() *land.Grid { return r.grid }

func (r *Room) Phase() Phase {
	return r.phase
}

// World 该房间独占的 ECS 世界；大厅阶段即可 CreateEntity / 挂组件。
func (r *Room) World() *ecs.World {
	return r.world
}

// StartBattle 注册战斗管线、挂载 tick→World.Update 并启动独立循环协程；ctx 用于上层整体关服/撤房时取消。
func (r *Room) StartBattle(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if r.phaseIs(PhaseClosed) {
		return ErrRoomClosed
	}
	if !r.phaseIs(PhaseLobby) {

		return ErrWrongPhase
	}

	if r.cancel != nil {
		return ErrWrongPhase
	}

	r.setPhase(PhasePreBattle)

	r.clk = clock.New(r.tps)
	r.loop = tick.NewLoop(r.clk)
	loopCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	loop := r.loop
	w := r.world

	dt := 1.0 / float64(r.clk.TPS())
	loop.Add(tick.FuncSubscriber(func(_ *clock.Clock) {
		w.Update(dt)
	}))

	r.setPhase(PhaseFighting)

	r.runWG.Add(1)
	go func() {
		defer r.runWG.Done()
		_ = loop.Run(loopCtx)
		r.destroy()
	}()
	return nil
}

// Settle 结束战斗循环并进入结算阶段；后续可读快照再 Destroy。
func (r *Room) Settle() error {
	if !r.phaseIs(PhaseFighting) {
		return ErrWrongPhase
	}
	r.setPhase(PhaseSettled)
	//  todo 做一些其他处理
	if r.cancel != nil {
		r.cancel()
	}
	return nil
}

func (r *Room) destroy() {
	w := r.world
	if w != nil {
		w.RemoveAllEntities()
	}
	r.setPhase(PhaseClosed)
	r.cancel = nil
	r.loop = nil
	r.clk = nil
}

// Loop 返回当前战斗循环（仅 Fighting 阶段有效；单测可用 [tick.Loop.Step]）。
func (r *Room) Loop() *tick.Loop {
	return r.loop
}

func (r *Room) setPhase(next Phase) {
	r.phase = next
}

// Shutdown 强制销毁房间
func (r *Room) Shutdown() {
	if r.phaseIs(PhaseClosed) {
		return
	}
	if r.cancel != nil {
		r.cancel()
	}
}
