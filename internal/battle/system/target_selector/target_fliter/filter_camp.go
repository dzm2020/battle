package target_fliter

import (
	"battle/internal/battle/config"
	"battle/internal/battle/utils"
	"encoding/json"
)

func campFilter(ctx *Context, f config.Filter) bool {
	var p config.CampFilter
	if err := json.Unmarshal(f.Params, &p); err != nil {
		return false
	}
	if len(p.AllowedCamps) == 0 {
		return true
	}
	rel := utils.CampRelation(ctx.World, ctx.Caster, ctx.Target)
	for _, c := range p.AllowedCamps {
		if c == rel {
			return true
		}
	}
	return false
}
