package room

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/resource"
	"battle/internal/battle/system"
	"context"
	"sync"
	"time"
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
	w.AddSystem(&system.SpawnSystem{})

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

	cancel context.CancelFunc
	runWG  sync.WaitGroup
}

func (r *Room) ID() uint64 {
	return r.id
}

// World 该房间独占的 ECS 世界；大厅阶段即可 CreateEntity / 挂组件。
func (r *Room) World() *ecs.World {
	return r.world
}

// StartBattle 启动战斗循环协程；ctx 用于上层整体关服/撤房时取消。
func (r *Room) StartBattle(ctx context.Context) error {
	if r.cancel != nil {
		return ErrWrongPhase
	}

	loopCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	r.runWG.Add(1)
	go func() {
		defer r.runWG.Done()
		r.runLoop(loopCtx)
		r.destroy()
	}()
	return nil
}

func (r *Room) runLoop(ctx context.Context) {
	tpsRes := ecs.GetResource[resource.TPS](r.world)
	if tpsRes == nil {
		return
	}

	var (
		ticker     *time.Ticker
		currentTPS int
	)

	resetTicker := func(tps *resource.TPS) {
		effective := tps.EffectiveTPS()
		if ticker != nil && effective == currentTPS {
			return
		}
		if ticker != nil {
			ticker.Stop()
			select {
			case <-ticker.C:
			default:
			}
		}
		currentTPS = effective
		ticker = time.NewTicker(tps.FrameDuration())
	}

	resetTicker(tpsRes)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tpsRes = ecs.GetResource[resource.TPS](r.world)
			if tpsRes == nil {
				return
			}
			resetTicker(tpsRes)

			r.world.Update(tpsRes.DeltaTime())
			tpsRes.Frame++
		}
	}
}

func (r *Room) destroy() {
	w := r.world
	if w != nil {
		w.RemoveAllEntities()
	}
	r.cancel = nil
}

// Shutdown 强制销毁房间
func (r *Room) Shutdown() {
	if r.cancel != nil {
		r.cancel()
	}
}
