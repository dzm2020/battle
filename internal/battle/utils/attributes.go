package utils

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/system/attrs"
	"strings"
)

func TransformXY(w *ecs.World, e ecs.Entity) (int, int, bool) {
	t, ok := w.GetComponent(e, &component.Transform2D{})
	if !ok {
		return 0, 0, false
	}
	tr := t.(*component.Transform2D)
	return tr.X, tr.Y, true
}

// HealthCurrent 读取实体当前生命（来自 [component.Attributes] 的 hp，唯一数据源）。
func HealthCurrent(w *ecs.World, e ecs.Entity) int {
	a, ok := w.GetComponent(e, &component.Attributes{})
	if !ok {
		return 0
	}
	return attrs.HP(a.(*component.Attributes))
}

// HealthMax 读取实体生命上限。
func HealthMax(w *ecs.World, e ecs.Entity) int {
	a, ok := w.GetComponent(e, &component.Attributes{})
	if !ok {
		return 0
	}
	return attrs.HPMax(a.(*component.Attributes))
}

// CampRelation 两个实体阵营关系。
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

// GetEntityCamp 返回实体阵营。
func GetEntityCamp(w *ecs.World, e ecs.Entity) (component.SideType, bool) {
	c, exists := w.GetComponent(e, &component.Team{})
	if !exists {
		return "", false
	}
	return c.(*component.Team).Side, true
}

func GetEntityAttributeValue(w *ecs.World, e ecs.Entity, name string) (float64, bool) {
	key := config.AttributeType(strings.ToLower(strings.TrimSpace(name)))
	if a, ok := w.GetComponent(e, &component.Attributes{}); ok {
		return float64(attrs.Current(a.(*component.Attributes), key)), true
	}
	return 0, false
}

func GetAttributeFinalValue(w *ecs.World, e ecs.Entity, typ config.AttributeType) int {
	var value int
	if a, ok := w.GetComponent(e, &component.Attributes{}); ok {
		attr := a.(*component.Attributes)
		if attrs.Current(attr, typ) > 0 {
			value = attrs.Current(attr, typ)
		}
		if sm, ok := w.GetComponent(e, &component.BuffStatModifiers{}); ok {
			modifier := sm.(*component.BuffStatModifiers)
			value += int(modifier.Modifiers[typ])
		}
	}
	return value
}
