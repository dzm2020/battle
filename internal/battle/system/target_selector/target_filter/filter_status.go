package target_filter

import (
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"encoding/json"
)

func statusFilter(ctx *Context, f config.Filter) bool {
	var p config.StatusFilter
	if err := json.Unmarshal(f.Params, &p); err != nil {
		return false
	}
	if p.StatusMask == 0 {
		return true
	}
	c, ok := ctx.World.GetComponent(ctx.Target, &component.BuffControlState{})
	if !ok {
		return false
	}
	fl := uint8(c.(*component.BuffControlState).Flags)
	m := uint8(p.StatusMask)
	return (fl & m) != 0
}
