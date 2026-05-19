package target_filter

import (
	"strings"

	"battle/ecs"
	"battle/internal/battle/config"

	"github.com/duke-git/lancet/v2/maputil"
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

var handlers = maputil.NewConcurrentMap[string, handler](8)

func typeKey(t config.FilterType) string {
	return strings.ToLower(strings.TrimSpace(string(t)))
}

func init() {
	handlers.Set(typeKey(config.FilterCamp), campFilter)
	handlers.Set(typeKey(config.FilterStatusMask), statusFilter)
	handlers.Set(typeKey(config.FilterProperty), attributeFilter)
}

func Apply(ctx *Context, filters ...config.Filter) bool {
	for _, f := range filters {
		fn, ok := handlers.Get(typeKey(f.Type))
		if !ok || fn == nil {
			return false
		}
		if !fn(ctx, f) {
			return false
		}
	}
	return true
}
