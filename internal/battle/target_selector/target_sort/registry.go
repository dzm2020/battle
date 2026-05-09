package target_sort

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"sort"
)

type compare func(w *ecs.World, ref, a, b ecs.Entity) int

var dict = make(map[config.TargetSortType]compare)

func registry(typ config.TargetSortType, cmp compare) {
	if typ == config.SortNone {
		return
	}
	if cmp == nil {
		delete(dict, typ)
		return
	}
	dict[typ] = cmp
}

func init() {
	registry(config.SortHealth, compareHealthCurrent)
	registry(config.SortPosition, compareDistanceSquared)
}

func Do(w *ecs.World, ref ecs.Entity, entityList []ecs.Entity, st config.TargetSortType, ord config.SortOrder) {
	if len(entityList) <= 1 || st == config.SortNone {
		return
	}
	cmp := dict[st]
	if cmp == nil {
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
