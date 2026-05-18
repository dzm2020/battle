package room

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/resource"
	"battle/internal/battle/system"
	"battle/internal/battle/tick"
	"context"
	"sync"
)

// Create 根据 dungeonId 加载副本配置，并按 [config.DungeonConfig.Type] 选择已注册的装配逻辑创建房间。
func Create(spec *resource.RoomSpec) (*Room, error) {
	w := ecs.NewWorld(100)
	r := &Room{
		id:    GetManager().NextID(),
		world: w,
	}

	component.Register(r.world)

	w.AddSystem(&system.BattleInitSystem{})

	ecs.AddResource(w, &resource.RoomID{ID: r.id})
	ecs.AddResource(w, &resource.RoomPhase{Phase: resource.PhaseLobby})
	ecs.AddResource(w, &resource.TPS{TPS: 60, Frame: 0})
	ecs.AddResource(w, &resource.SpawnRequestQueue{})
	ecs.AddResource(w, spec)

	if err := r.StartBattle(context.Background()); err != nil {
		return nil, err
	}

	GetManager().Add(r)
	return r, nil
}

type Room struct {
	id    uint64
	world *ecs.World
	// 逻辑帧驱动
	loop   *tick.Loop
	cancel context.CancelFunc
	runWG  sync.WaitGroup
}

// World 该房间独占的 ECS 世界；大厅阶段即可 CreateEntity / 挂组件。
func (r *Room) World() *ecs.World {
	return r.world
}

// StartBattle 注册战斗管线、挂载 tick→World.Update 并启动独立循环协程；ctx 用于上层整体关服/撤房时取消。
func (r *Room) StartBattle(ctx context.Context) error {
	if r.cancel != nil {
		return ErrWrongPhase
	}

	r.loop = tick.NewLoop(tick.New(r.tps))

	loopCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	loop := r.loop
	w := r.world

	dt := 1.0 / float64(loop.Clock().TPS())
	loop.Add(tick.FuncSubscriber(func(_ *tick.Clock) {
		w.Update(dt)
	}))

	r.runWG.Add(1)
	go func() {
		defer r.runWG.Done()
		_ = loop.Run(loopCtx)
		r.destroy()
	}()
	return nil
}

func (r *Room) destroy() {
	w := r.world
	if w != nil {
		w.RemoveAllEntities()
	}
	r.cancel = nil
	r.loop = nil
}

// Loop 返回当前战斗循环（仅 Fighting 阶段有效；单测可用 [tick.Loop.Step]）。
func (r *Room) Loop() *tick.Loop {
	return r.loop
}

// Shutdown 强制销毁房间
func (r *Room) Shutdown() {
	if r.cancel != nil {
		r.cancel()
	}
}
