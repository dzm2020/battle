package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/land"
	"battle/internal/battle/log"
	"battle/internal/battle/resource"
	"battle/internal/battle/system/entity_factory"
)

// SpawnSystem 消费 [runtime.BattleContext].SpawnQueue，按请求创建单位并登记到 Grid。
type SpawnSystem struct {
	world *ecs.World
}

func (s *SpawnSystem) Initialize(w *ecs.World) {
	s.world = w
}

func (s *SpawnSystem) Update(_ float64) {
	grid, _ := resource.Grid(s.world)

	queue, _ := resource.SpawnQueue(s.world)

	var pending []*resource.SpawnRequest
	for _, request := range queue.Queue {
		if s.fulfill(grid, request) {
			continue
		}
		pending = append(pending, request)
	}
	queue.Queue = pending
}

func (s *SpawnSystem) fulfill(grid *land.Grid, req *resource.SpawnRequest) bool {
	if req.UnitID == 0 && (req.Data == nil || req.Data.ID == 0) {
		return true
	}

	cellX, cellY, ok := resolveSpawnCell(grid, req)
	if !ok {
		return false
	}

	components := req.Components
	components = append(components, &component.Team{Side: req.Side, Entity: req.TeamEntity})
	components = append(components, &component.Transform2D{X: cellX, Y: cellY})

	var (
		e   ecs.Entity
		err error
	)
	if req.Data != nil {
		data := req.Data
		e, err = entity_factory.CreateFromPBUnit(s.world, data, components...)
	} else {
		e, err = entity_factory.CreateByConfigID(s.world, req.UnitID, components...)
	}
	if err != nil {
		log.Error("[spawn] 创建单位失败 unitID=%d dataID=%d: %v", req.UnitID, req.Data.ID, err)
		return true
	}
	if err = grid.AddUnit(uint64(e), cellX, cellY); err != nil {
		log.Error("[spawn] 登记网格失败 entity=%v: %v", e, err)
		return true
	}

	if req.TeamEntity != 0 && req.Data != nil && req.Data.ID != 0 {
		if pc, ok := s.world.GetComponent(req.TeamEntity, &component.Player{}); ok {
			pc.(*component.Player).Units[req.Data.ID] = e
		}
	}
	return true
}

func resolveSpawnCell(grid *land.Grid, req *resource.SpawnRequest) (cellX, cellY int, ok bool) {
	if req.CellX >= 0 && req.CellY >= 0 {
		return req.CellX, req.CellY, true
	}
	return grid.PickFreeCell(req.Side == component.SideTypeRed)
}
