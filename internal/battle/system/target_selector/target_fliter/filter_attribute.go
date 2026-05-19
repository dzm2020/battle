package target_fliter

import (
	"battle/internal/battle/config"
	"battle/internal/battle/system/attrs"
	"battle/internal/battle/system/utils"
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
	return utils.CompareFloat64(cur, p.Op, p.Value)
}
