// Package attrs 提供对 [component.Attributes] 的读写辅助，以及基于 World 的属性/阵营查询。
//
// 约定：运行时属性变更应在 System（或出生装配 entity_factory）中通过本包完成；
// [component] 包仅定义数据结构，不包含修改逻辑。
package attrs

import (
	"strings"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// Current 读取 Current；nil、空 map 或未登记键返回 0。
func Current(a *component.Attributes, key config.AttributeType) int {
	if a == nil || a.Base == nil {
		return 0
	}
	if x := a.Base[key]; x != nil {
		return x.Current
	}
	return 0
}

// SetCurrent 写入 Current。键已存在时保留 Max，并在 Max>0 时将 Current 夹在 [0, Max]；
// 新键则 Current、Max 均为 v。
func SetCurrent(a *component.Attributes, key config.AttributeType, v int) {
	if a == nil {
		return
	}
	if a.Base == nil {
		a.Base = make(map[config.AttributeType]*component.Attribute)
	}
	if existing := a.Base[key]; existing != nil {
		existing.Current = v
		if existing.Max > 0 && existing.Current > existing.Max {
			existing.Current = existing.Max
		}
		if existing.Current < 0 {
			existing.Current = 0
		}
		return
	}
	a.Base[key] = &component.Attribute{Current: v, Max: v}
}

// Add 在给定键上累加 Current。
func Add(a *component.Attributes, key config.AttributeType, delta int) {
	SetCurrent(a, key, Current(a, key)+delta)
}

// Sub 在给定键上减少 Current。
func Sub(a *component.Attributes, key config.AttributeType, delta int) {
	Add(a, key, -delta)
}

// SetRange 同时设置 Current 与 Max（用于单位初始化与读表生成属性）。
func SetRange(a *component.Attributes, key config.AttributeType, current, max int) {
	if max < current {
		max = current
	}
	if a == nil {
		return
	}
	if a.Base == nil {
		a.Base = make(map[config.AttributeType]*component.Attribute)
	}
	a.Base[key] = &component.Attribute{Current: current, Max: max}
}

// Max 读取 Max；未登记键返回 0。
func Max(a *component.Attributes, key config.AttributeType) int {
	if a == nil || a.Base == nil {
		return 0
	}
	if x := a.Base[key]; x != nil {
		return x.Max
	}
	return 0
}

// HP 读取生命值 Current。
func HP(a *component.Attributes) int {
	return Current(a, config.AttrHp)
}

// HPMax 读取生命值上限。
func HPMax(a *component.Attributes) int {
	return Max(a, config.AttrHp)
}

// TransformXY 读取实体平面坐标。
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
	return HP(a.(*component.Attributes))
}

// HealthMax 读取实体生命上限。
func HealthMax(w *ecs.World, e ecs.Entity) int {
	a, ok := w.GetComponent(e, &component.Attributes{})
	if !ok {
		return 0
	}
	return HPMax(a.(*component.Attributes))
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

// GetEntityAttributeValue 按属性名（不区分大小写）读取 Current。
func GetEntityAttributeValue(w *ecs.World, e ecs.Entity, name string) (float64, bool) {
	key := config.AttributeType(strings.ToLower(strings.TrimSpace(name)))
	if a, ok := w.GetComponent(e, &component.Attributes{}); ok {
		return float64(Current(a.(*component.Attributes), key)), true
	}
	return 0, false
}

// GetAttributeFinalValue 读取属性 Current 与 Buff 修正后的最终值。
func GetAttributeFinalValue(w *ecs.World, e ecs.Entity, typ config.AttributeType) int {
	var value int
	if a, ok := w.GetComponent(e, &component.Attributes{}); ok {
		attr := a.(*component.Attributes)
		if Current(attr, typ) > 0 {
			value = Current(attr, typ)
		}
		if sm, ok := w.GetComponent(e, &component.BuffStatModifiers{}); ok {
			modifier := sm.(*component.BuffStatModifiers)
			value += int(modifier.Modifiers[typ])
		}
	}
	return value
}
