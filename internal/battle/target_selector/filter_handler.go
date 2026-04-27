package target_selector

import (
	"strings"

	"battle/ecs"
	"battle/internal/battle/config"
)

func init() {
	registerFilterHandler(string(config.FilterCamp), handleFilterCamp)
	registerFilterHandler(string(config.FilterStatusMask), handleFilterStatusMask)
	registerFilterHandler(string(config.FilterProperty), handleFilterProperty)
}

func evalRootFilters(ctx *Context, filters []config.Filter) bool {
	if len(filters) == 0 {
		return true
	}
	for _, f := range filters {
		if !EvalFilter(ctx, f) {
			return false
		}
	}
	return true
}

// Context 单次筛选时的上下文（施法者、主目标、当前候选实体及其生命组件）。
type Context struct {
	World  *ecs.World
	Caster ecs.Entity
	Target ecs.Entity
}

// FilterHandler 处理某一类 [config.Filter.Type]（叶子或组合均可）。
type FilterHandler func(ctx *Context, f config.Filter) bool

var filterHandlers = make(map[string]FilterHandler)

func filterTypeKey(t config.FilterType) string {
	return strings.ToLower(strings.TrimSpace(string(t)))
}

// RegisterFilterHandler 注册或覆盖某类 Filter 的处理逻辑（扩展点，与 buff handler 一致）。
func RegisterFilterHandler(typ string, fn FilterHandler) {
	if typ == "" || fn == nil {
		return
	}
	filterHandlers[filterTypeKey(config.FilterType(typ))] = fn
}

func registerFilterHandler(typ string, fn FilterHandler) {
	RegisterFilterHandler(typ, fn)
}

// EvalFilter 执行单条 Filter；未知 type 时默认放行
func EvalFilter(ctx *Context, f config.Filter) bool {
	if ctx == nil || ctx.World == nil {
		return false
	}
	if fn := filterHandlers[filterTypeKey(f.Type)]; fn != nil {
		return fn(ctx, f)
	}
	return true
}

//func handleFilterAnd(ctx *Context, f config.Filter) bool {
//	var subs []config.Filter
//	if err := json.Unmarshal(f.Params, &subs); err != nil || len(subs) == 0 {
//		return false
//	}
//	for _, s := range subs {
//		if !EvalFilter(ctx, s) {
//			return false
//		}
//	}
//	return true
//}
//
//func handleFilterOr(ctx *Context, f config.Filter) bool {
//	var subs []config.Filter
//	if err := json.Unmarshal(f.Params, &subs); err != nil || len(subs) == 0 {
//		return false
//	}
//	for _, s := range subs {
//		if EvalFilter(ctx, s) {
//			return true
//		}
//	}
//	return false
//}
//
//func handleFilterNot(ctx *Context, f config.Filter) bool {
//	var sub config.Filter
//	if err := json.Unmarshal(f.Params, &sub); err != nil {
//		return false
//	}
//	return !EvalFilter(ctx, sub)
//}
