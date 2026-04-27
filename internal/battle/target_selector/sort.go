package target_selector

import (
	"math"
	"sort"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// EntityCompare 比较两个候选实体在排序键上的先后。
// 返回值：<0 表示 a 应排在 b 前（在「升序度量」语义下，如血量更低、距离更近）；
// 0 表示键相等，由实体 ID 稳定打破平局；>0 表示 a 应排在 b 后。
// ref 一般为施法者，用于「距施法者距离」等排序。
type EntityCompare func(w *ecs.World, ref, a, b ecs.Entity) int

var sortStrategies = make(map[config.TargetSortType]EntityCompare)

// RegisterSortStrategy 注册或覆盖某一 [config.TargetSortType] 的比较函数；cmp 为 nil 时注销。
// 用于扩展新的排序维度（威胁值、随机种子等）而无需改 sort.go 主体。
func RegisterSortStrategy(typ config.TargetSortType, cmp EntityCompare) {
	if typ == config.SortNone {
		return
	}
	if cmp == nil {
		delete(sortStrategies, typ)
		return
	}
	sortStrategies[typ] = cmp
}

func init() {
	RegisterSortStrategy(config.SortHealth, compareHealthCurrent)
	RegisterSortStrategy(config.SortPosition, compareDistanceSquared)
}

func sortTargets(w *ecs.World, ref ecs.Entity, ents []ecs.Entity, st config.TargetSortType, ord config.SortOrder) {
	if len(ents) <= 1 || st == config.SortNone {
		return
	}
	cmp := sortStrategies[st]
	if cmp == nil {
		stableSortByEntityID(ents)
		return
	}
	sort.SliceStable(ents, func(i, j int) bool {
		c := cmp(w, ref, ents[i], ents[j])
		if c == 0 {
			return ents[i] < ents[j]
		}
		if ord == config.OrderDesc {
			return c > 0
		}
		return c < 0
	})
}

func stableSortByEntityID(ents []ecs.Entity) {
	sort.SliceStable(ents, func(i, j int) bool { return ents[i] < ents[j] })
}

func compareHealthCurrent(w *ecs.World, _ ecs.Entity, a, b ecs.Entity) int {
	ha := healthCurrent(w, a)
	hb := healthCurrent(w, b)
	switch {
	case ha < hb:
		return -1
	case ha > hb:
		return 1
	default:
		return 0
	}
}

func compareDistanceSquared(w *ecs.World, ref, a, b ecs.Entity) int {
	da := distanceSquaredFromRef(w, ref, a)
	db := distanceSquaredFromRef(w, ref, b)
	switch {
	case da < db:
		return -1
	case da > db:
		return 1
	default:
		return 0
	}
}

func distanceSquaredFromRef(w *ecs.World, ref, e ecs.Entity) float64 {
	rx, ry, rok := transformXY(w, ref)
	if !rok {
		rx, ry = 0, 0
	}
	x, y, ok := transformXY(w, e)
	if !ok {
		return math.MaxFloat64
	}
	dx := x - rx
	dy := y - ry
	return dx*dx + dy*dy
}

func transformXY(w *ecs.World, e ecs.Entity) (float64, float64, bool) {
	t, ok := w.GetComponent(e, &component.Transform2D{})
	if !ok {
		return 0, 0, false
	}
	tr := t.(*component.Transform2D)
	return tr.X, tr.Y, true
}

func healthCurrent(w *ecs.World, e ecs.Entity) int {
	h, ok := w.GetComponent(e, &component.Health{})
	if !ok {
		return 0
	}
	return h.(*component.Health).Current
}
