package target_filter

import (
	"battle/internal/battle/config"
	"battle/internal/battle/system/attrs"
	"battle/internal/battle/system/combatmath"
	"encoding/json"
)

func attributeFilter(ctx *Context, f config.Filter) bool {
	var p config.PropertyFilter
	if err := json.Unmarshal(f.Params, &p); err != nil {
		return false
	}
	cur, ok := attrs.GetEntityAttributeValue(ctx.World, ctx.Target, p.Property)
	if !ok {
		return false
	}
	return combatmath.CompareFloat64(cur, p.Op, p.Value)
}
