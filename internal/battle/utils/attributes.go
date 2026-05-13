package utils

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"strings"
)

func TransformXY(w *ecs.World, e ecs.Entity) (float64, float64, bool) {
	t, ok := w.GetComponent(e, &component.Transform2D{})
	if !ok {
		return 0, 0, false
	}
	tr := t.(*component.Transform2D)
	return tr.X, tr.Y, true
}

func HealthCurrent(w *ecs.World, e ecs.Entity) int {
	h, ok := w.GetComponent(e, &component.Health{})
	if !ok {
		return 0
	}
	return h.(*component.Health).Current
}

// CampRelation 两个实体 阵营关系
func CampRelation(w *ecs.World, caster, target ecs.Entity) config.Camp {
	if target == caster {
		return config.CampFriend
	}
	cs, cOK := GetEntityCamp(w, caster)
	ts, tOK := GetEntityCamp(w, target)
	if !cOK || !tOK {
		return config.CampNeutral
	}
	if ts == cs {
		return config.CampFriend
	}
	return config.CampEnemy
}

// GetEntityCamp 返回实体阵营
func GetEntityCamp(w *ecs.World, e ecs.Entity) (side uint8, ok bool) {
	c, exists := w.GetComponent(e, &component.Team{})
	if !exists {
		return 0, false
	}
	return c.(*component.Team).Side, true
}

func GetEntityAttributeValue(w *ecs.World, e ecs.Entity, name string) (float64, bool) {
	key := strings.ToLower(strings.TrimSpace(name))
	if a, ok := w.GetComponent(e, &component.Attributes{}); ok {
		return float64(a.(*component.Attributes).Get(key)), true
	}
	return 0, false
}

func GetAttributeFinalValue(w *ecs.World, e ecs.Entity, typ config.AttributeType) int {
	var value int
	if a, ok := w.GetComponent(e, &component.Attributes{}); ok {
		attr := a.(*component.Attributes)
		if attr.Get(typ) > 0 {
			value = attr.Get(typ)
		}
		if sm, ok := w.GetComponent(e, &component.BuffStatModifiers{}); ok {
			modifier := sm.(*component.BuffStatModifiers)
			value += int(modifier.Modifiers[typ])
		}
	}
	return value
}
