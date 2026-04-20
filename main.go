package main

import (
	"github.com/edwinsyarief/teishoku"
)

type Position struct {
	X, Y float64
}

type Velocity struct {
	DX, DY float64
}

func main() {

	// 创建 World
	world := teishoku.NewWorld(10000)

	// 使用 Builder 创建实体
	builder := teishoku.NewBuilder2[Position, Velocity](world)
	builder.NewEntities(100)

	// 查询和迭代
	query := teishoku.NewFilter2[Position, Velocity](world)
	for query.Next() {
		pos, vel := query.Get()
		pos.X += vel.DX
		pos.Y += vel.DY
	}
}
