package room

import (
	"battle/internal/battle/config"
	"battle/internal/battle/factory/room_factory"
	"battle/internal/battle/pb"
	"context"
	"errors"
	"sync"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/land"
	"battle/internal/battle/system/runtime"
	"battle/internal/battle/tick"
)

// Phase 房间当前阶段（仅服务端权威）；由 [Room] 在各 API 内直接维护 [Room.phase] 字段。
type Phase int8

const (
	PhaseLobby     Phase = iota // 等待加入 / 准备
	PhasePreBattle              // 已开始开战流程，禁止再 Join（防止与 InitBattle 交错）
	PhaseFighting               // 战斗循环运行中
	PhaseSettled                // 已结算，等待销毁或复盘
	PhaseClosed                 // 已关闭，不可再操作
)

var (
	ErrNoDungeonConfig = errors.New("no dungeon config")
	ErrNoMapConfig     = errors.New("no map config")
)

// CreateRoom 根据 dungeonId 加载副本配置，并按 [config.DungeonConfig.Type] 选择已注册的装配逻辑创建房间。
func CreateRoom(dungeonId int32, spec *room_factory.Spec) (*Room, error) {
	r := &Room{
		id:    GetManager().NextID(),
		tps:   60,
		phase: PhaseLobby,
		world: ecs.NewWorld(100),
	}

	GetManager().Add(r)

	desc := config.GetDungeonConfigByID(dungeonId)
	if desc == nil {
		return nil, ErrNoDungeonConfig
	}
	if desc.Type == config.DungeonTypePVP && (spec == nil || spec.Self == nil || spec.Enemy == nil) {
		return nil, ErrUseCreatePVPRoom
	}

	grid, err := land.CreateGridByID(desc.MapID)
	if err != nil {
		return nil, err
	}

	component.Init(r.world)
	r.SetGrid(grid)
	runtime.Install(r.world, runtime.New(grid))

	spec.World = r.world
	spec.Desc = desc
	//  构建房间
	if err = room_factory.Create(spec); err != nil {
		return nil, err
	}

	if err = r.StartBattle(context.Background()); err != nil {
		return nil, err
	}

	return r, nil
}

// CreatePVPRoom 创建 PVP 房间（须同时提供 Self 与 Enemy）。
func CreatePVPRoom(dungeonId int32, self, enemy *pb.Player) (*Room, error) {
	return CreateRoom(dungeonId, &room_factory.Spec{Self: self, Enemy: enemy})
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
	loop   *tick.Loop
	cancel context.CancelFunc
	runWG  sync.WaitGroup
}

// phaseIs 判断当前是否处于指定阶段；调用方必须已持有 [Room.mu] 读锁或写锁。
func (r *Room) phaseIs(p Phase) bool {
	return r.phase == p
}

func (r *Room) ID() uint64 { return r.id }

// SetGrid 设置空间网格，并同步到 [runtime.BattleContext].Grid（须已 [runtime.Install] 或随后 Install）。
func (r *Room) SetGrid(base *land.Grid) {
	r.grid = base
	if r.world == nil {
		return
	}
	if ctx, ok := runtime.Get(r.world); ok {
		ctx.Grid = base
		return
	}
	runtime.Install(r.world, runtime.New(base))
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

	r.loop = tick.NewLoop(tick.New(r.tps))

	loopCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	loop := r.loop
	w := r.world

	dt := 1.0 / float64(loop.Clock().TPS())
	loop.Add(tick.FuncSubscriber(func(_ *tick.Clock) {
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
