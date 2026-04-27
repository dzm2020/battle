package target_selector

import (
	"encoding/json"
	"math"
	"strings"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

func handleFilterCamp(ctx *Context, f config.Filter) bool {
	var p config.CampFilter
	if err := json.Unmarshal(f.Params, &p); err != nil {
		return false
	}
	if len(p.AllowedCamps) == 0 {
		return true
	}
	rel := campRelation(ctx.World, ctx.Caster, ctx.Target)
	for _, c := range p.AllowedCamps {
		if c == rel {
			return true
		}
	}
	return false
}

// campRelation 阵营关系
func campRelation(w *ecs.World, caster, target ecs.Entity) config.Camp {
	if target == caster {
		return config.CampFriend
	}
	cs, cOK := teamSide(w, caster)
	ts, tOK := teamSide(w, target)
	if !cOK || !tOK {
		return config.CampNeutral
	}
	if ts == cs {
		return config.CampFriend
	}
	return config.CampEnemy
}

// 获取阵营
func teamSide(w *ecs.World, e ecs.Entity) (side uint8, ok bool) {
	c, exists := w.GetComponent(e, &component.Team{})
	if !exists {
		return 0, false
	}
	return c.(*component.Team).Side, true
}

func handleFilterStatusMask(ctx *Context, f config.Filter) bool {
	var p config.StatusFilter
	if err := json.Unmarshal(f.Params, &p); err != nil {
		return false
	}
	if p.StatusMask == 0 {
		return true
	}
	c, ok := ctx.World.GetComponent(ctx.Target, &component.ControlState{})
	if !ok {
		return false
	}
	fl := uint8(c.(*component.ControlState).Flags)
	m := uint8(p.StatusMask)
	return (fl & m) != 0
}

func handleFilterProperty(ctx *Context, f config.Filter) bool {
	var p config.PropertyFilter
	if err := json.Unmarshal(f.Params, &p); err != nil {
		return false
	}
	cur, ok := propertyValue(ctx.World, ctx.Target, p.Property)
	if !ok {
		return false
	}
	return compareFloat(cur, p.Op, p.Value)
}

func propertyValue(w *ecs.World, e ecs.Entity, name string) (float64, bool) {
	key := strings.ToLower(strings.TrimSpace(name))
	if a, ok := w.GetComponent(e, &component.Attributes{}); ok {
		return float64(a.(*component.Attributes).Get(key)), true
	}
	return 0, false
}

func compareFloat(cur float64, op string, val float64) bool {
	switch strings.TrimSpace(op) {
	case ">":
		return cur > val
	case "<":
		return cur < val
	case "==", "=":
		return math.Abs(cur-val) < 1e-6
	case "!=":
		return math.Abs(cur-val) >= 1e-6
	case ">=":
		return cur >= val
	case "<=":
		return cur <= val
	default:
		return false
	}
}
