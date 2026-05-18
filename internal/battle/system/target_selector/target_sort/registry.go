package target_sort

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"sort"

	"github.com/duke-git/lancet/v2/maputil"
)

type compare func(w *ecs.World, ref, a, b ecs.Entity) int

var comparators = maputil.NewConcurrentMap[config.TargetSortType, compare](4)

func Apply(w *ecs.World, ref ecs.Entity, entityList []ecs.Entity, st config.TargetSortType, ord config.SortOrder) {
	if len(entityList) <= 1 || st == config.SortNone {
		return
	}
	cmp, ok := comparators.Get(st)
	if !ok || cmp == nil {
		sort.SliceStable(entityList, func(i, j int) bool { return entityList[i] < entityList[j] })
		return
	}
	sort.SliceStable(entityList, func(i, j int) bool {
		c := cmp(w, ref, entityList[i], entityList[j])
		if c == 0 {
			return entityList[i] < entityList[j]
		}
		if ord == config.OrderDesc {
			return c > 0
		}
		return c < 0
	})
}

func init() {
	comparators.Set(config.SortHealth, compareHealthCurrent)
	comparators.Set(config.SortPosition, compareDistanceSquared)
}
