package ecs

import (
	"fmt"
	"testing"
)

// ========== 定义组件 ==========

type Position struct {
	X, Y float64
}

func (p *Position) Component() {} // 实现 Component 接口

type Velocity struct {
	DX, DY float64
}

func (v *Velocity) Component() {}

type Health struct {
	Current, Max int
}

func (h *Health) Component() {}

// ========== 定义系统 ==========

// MovementSystem 移动系统
type MovementSystem struct {
	world *World
	query *Query2[*Position, *Velocity]
}

func (s *MovementSystem) Initialize(w *World) {
	s.world = w
	s.query = NewQuery2[*Position, *Velocity](w)
}

func (s *MovementSystem) Update(dt float64) {
	s.query.ForEach(func(e Entity, pos *Position, vel *Velocity) {
		pos.X += vel.DX * dt
		pos.Y += vel.DY * dt
		// 更新回世界（因为组件是值类型，需要重新设置）
		fmt.Printf("Entity %d moved to (%.2f, %.2f)\n", e, pos.X, pos.Y)
	})
}

// HealthSystem 生命值系统
type HealthSystem struct {
	world *World
	query *Query[*Health]
}

func (s *HealthSystem) Initialize(w *World) {
	s.world = w
	s.query = NewQuery[*Health](w)
}

func (s *HealthSystem) Update(dt float64) {
	// 每秒恢复 10% 生命值
	healAmount := int(10 * dt)

	s.query.ForEach(func(e Entity, health *Health) {
		if health.Current < health.Max {
			health.Current += healAmount
			if health.Current > health.Max {
				health.Current = health.Max
			}
			fmt.Printf("Entity %d healed to %d/%d\n", e, health.Current, health.Max)
		}
	})
}

// ========== 主程序 ==========

func TestEcs(t *testing.T) {

	// 创建世界
	world := NewWorld(10)

	// 注册组件类型（可选，首次添加时会自动注册）
	world.Registry().Register(&Position{})
	world.Registry().Register(&Velocity{})
	world.Registry().Register(&Health{})

	// 创建实体
	player := world.CreateEntity()
	world.AddComponent(player, &Position{X: 100, Y: 200})
	world.AddComponent(player, &Velocity{DX: 10, DY: 5})
	world.AddComponent(player, &Health{Current: 80, Max: 100})

	enemy := world.CreateEntity()
	world.AddComponent(enemy, &Position{X: 300, Y: 400})
	world.AddComponent(enemy, &Velocity{DX: -5, DY: -3})
	world.AddComponent(enemy, &Health{Current: 50, Max: 100})

	// 添加系统
	movementSys := &MovementSystem{}
	healthSys := &HealthSystem{}
	world.AddSystem(movementSys)
	world.AddSystem(healthSys)

	// 模拟 3 帧
	for frame := 0; frame < 3; frame++ {
		fmt.Printf("\n=== Frame %d ===\n", frame+1)
		world.Update(0.1) // dt = 0.1 秒
	}

	// 使用 Query 收集实体
	fmt.Println("\n=== Query Results ===")
	query := NewQuery[*Health](world)
	entities := query.Collect()
	fmt.Printf("Entities with Health: %v\n", entities)
}
