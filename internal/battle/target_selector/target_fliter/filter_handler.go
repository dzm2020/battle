package target_fliter

import (
	"strings"

	"battle/ecs"
	"battle/internal/battle/config"
)

type (
	// Context 单次筛选时的上下文（施法者、主目标、当前候选实体及其生命组件）。
	Context struct {
		World  *ecs.World
		Caster ecs.Entity
		Target ecs.Entity
	}
	handler func(ctx *Context, f config.Filter) bool
)

var (
	handlers = make(map[string]handler)
)

func registry(typ string, fn handler) {
	if typ == "" || fn == nil {
		return
	}
	handlers[typeKey(config.FilterType(typ))] = fn
}

func typeKey(t config.FilterType) string {
	return strings.ToLower(strings.TrimSpace(string(t)))
}

func init() {
	registry(string(config.FilterCamp), campFilter)
	registry(string(config.FilterStatusMask), statusFilter)
	registry(string(config.FilterProperty), attributeFilter)
}

// Do true:满足条件  false：不满足条件
func Do(ctx *Context, filters ...config.Filter) bool {
	if len(filters) == 0 {
		return true
	}
	for _, f := range filters {
		fn := handlers[typeKey(f.Type)]
		if fn == nil {
			return false
		}
		if !fn(ctx, f) {
			return false
		}
	}
	return true
}
