// Package runtime 提供单局战斗在 [ecs.World] 上的全局资源（BattleContext）。
//
// 注入时机：开房时由 [room.Room] 调用 [Install] 一次；[Room.SetGrid] 仅更新 Context.Grid。
// 消费方：system（如 SpawnSystem）、room_builder（入队 SpawnRequest）通过 [Get] / [MustGet] 访问。
// 勿再对 *land.Grid、*component.SpawnRequestQueue 单独 InsertResource。
package runtime

import (
	"errors"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/land"
)

var (
	ErrNoContext   = errors.New("runtime: BattleContext 未注入")
	ErrNilGrid     = errors.New("runtime: Grid 未设置")
	ErrNilSpawnQueue = errors.New("runtime: SpawnQueue 未初始化")
)

// BattleContext 单局战斗 World 级单例；通过 [Install] 写入 ecs.World.resources。
type BattleContext struct {
	Grid       *land.Grid
	SpawnQueue *component.SpawnRequestQueue
}

// New 创建战斗上下文；grid 可为 nil（仅测试或延迟 SetGrid 时）。
func New(grid *land.Grid) *BattleContext {
	return &BattleContext{
		Grid:       grid,
		SpawnQueue: &component.SpawnRequestQueue{},
	}
}

// Install 将 BattleContext 注入 World（每局一次；重复调用会覆盖）。
func Install(w *ecs.World, ctx *BattleContext) {
	if w == nil || ctx == nil {
		return
	}
	if ctx.SpawnQueue == nil {
		ctx.SpawnQueue = &component.SpawnRequestQueue{}
	}
	ecs.InsertResource(w, ctx)
}

// Get 读取已注入的 BattleContext。
func Get(w *ecs.World) (*BattleContext, bool) {
	return ecs.GetResource[*BattleContext](w)
}

// MustGet 读取 BattleContext；未注入时返回 [ErrNoContext]。
func MustGet(w *ecs.World) (*BattleContext, error) {
	ctx, ok := Get(w)
	if !ok || ctx == nil {
		return nil, ErrNoContext
	}
	return ctx, nil
}

// Grid 返回空间网格。
func Grid(w *ecs.World) (*land.Grid, bool) {
	ctx, ok := Get(w)
	if !ok || ctx == nil || ctx.Grid == nil {
		return nil, false
	}
	return ctx.Grid, true
}

// SpawnQueue 返回刷怪请求队列。
func SpawnQueue(w *ecs.World) (*component.SpawnRequestQueue, bool) {
	ctx, ok := Get(w)
	if !ok || ctx == nil || ctx.SpawnQueue == nil {
		return nil, false
	}
	return ctx.SpawnQueue, true
}

// EnqueueSpawn 向本局刷怪队列追加请求。
func EnqueueSpawn(w *ecs.World, req *component.SpawnRequest) error {
	ctx, err := MustGet(w)
	if err != nil {
		return err
	}
	if req == nil {
		return nil
	}
	ctx.SpawnQueue.Queue = append(ctx.SpawnQueue.Queue, req)
	return nil
}
