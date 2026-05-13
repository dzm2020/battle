package room

import (
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/entity_factory"
	"battle/internal/battle/land"
	"battle/internal/battle/pb"
	"context"
	"sync"

	"battle/ecs"
	"battle/internal/battle/clock"
	"battle/internal/battle/system"
	"battle/internal/battle/tick"
)

func newRoom(dungeonId int32, players []*pb.Player) (*Room, error) {
	desc := config.GetDungeonConfigByID(dungeonId)
	if desc == nil {
		return nil, ErrNoDungeonConfig
	}

	grid, err := land.NewSpatialGrid(1, 10, 1, 10, 1)
	if err != nil {
		return nil, err
	}
	r := &Room{
		id:    GetManager().NextID(),
		tps:   60,
		phase: PhaseLobby,
		world: ecs.NewWorld(100),
		grid:  grid,
	}
	r.world.SetContext(r.world)
	component.Register(r.world)
	//  添加怪物
	for _, monsterId := range desc.Monster {
		e, _ := entity_factory.CreateByID(r.World(), monsterId)
		cellX, cellZ, _ := grid.FirstFreeCellDesc()
		_ = grid.AddUnit(uint64(e), cellX, cellZ)
		r.world.AddComponent(e, &component.Transform2D{
			X: cellX,
			Y: cellZ,
		})
	}

	//  添加玩家
	for _, player := range players {
		_, _ = entity_factory.CreateByPlayer(r.World(), player)
		cellX, cellZ, _ := grid.FirstFreeCellDesc()
		_ = grid.AddUnit(uint64(e), cellX, cellZ)
		r.world.AddComponent(e, &component.Transform2D{
			X: cellX,
			Y: cellZ,
		})
	}

	GetManager().Add(r)
	return r, nil
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

	system.AddSystems(w)

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
