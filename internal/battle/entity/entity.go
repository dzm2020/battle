package entity

import (
	"battle/internal/battle/attr"
	"battle/internal/battle/calc"
	"battle/internal/battle/control"
	"battle/internal/battle/cooldown"
	"battle/internal/battle/geom"
)

// Entity 战斗实体：聚合身份、阵营、基础属性、缓存衍生属性与局内状态。
// 不内嵌计算逻辑，Recalculate 显式注入 Calculator，降低与公式实现的耦合。
// 技能、Buff 相关字段由外部系统写入；Entity 仅提供数据挂点，不引用 skill 包以避免循环依赖。
type Entity struct {
	ID      string
	Camp    int8
	Base    attr.Base
	Derived attr.Derived
	Runtime attr.Runtime

	// Pos 逻辑位置（服务器权威），用于距离校验；第 9 天可与 AOI 对齐。
	Pos geom.Vec2
	// Control 眩晕/沉默等控制位；Buff 系统（第 7 天）在 Tick 中更新。
	Control control.Flags
	// SkillCD 每实体技能冷却表，key 为 SkillConfig.ID；由 InitBattle 重置。
	SkillCD *cooldown.Book
	// KnownSkills 已学会技能集合；网关/养成服在进房前写入。
	KnownSkills map[string]struct{}
}

// New 仅构造数据；满血满蓝需调用 InitBattle 或自行赋值 Runtime。
func New(id string, camp int8, base attr.Base) *Entity {
	return &Entity{
		ID:   id,
		Camp: camp,
		Base: base,
	}
}

// Recalculate 根据当前 Base 刷新 Derived，并按上限裁剪 Runtime。
func (e *Entity) Recalculate(c calc.Calculator) {
	e.Derived = c.DerivedFromBase(e.Base)
	c.ApplyMaxToRuntime(&e.Runtime, e.Derived)
}

// InitBattle 进入战斗时重置临时状态（第 4 天进房时可调用）。
func (e *Entity) InitBattle(c calc.Calculator) {
	e.Recalculate(c)
	e.Runtime.CurHP = e.Derived.MaxHP
	e.Runtime.CurMP = e.Derived.MaxMP
	e.Runtime.Shield = 0
	e.Control = 0
	e.SkillCD = cooldown.NewBook()
}

// GrantSkill 将技能 ID 加入已知集合（非并发安全：应在单线程或持房间锁时调用）。
func (e *Entity) GrantSkill(skillID string) {
	if e.KnownSkills == nil {
		e.KnownSkills = make(map[string]struct{})
	}
	e.KnownSkills[skillID] = struct{}{}
}

// KnowsSkill 是否已学会该技能。
func (e *Entity) KnowsSkill(skillID string) bool {
	if e.KnownSkills == nil {
		return false
	}
	_, ok := e.KnownSkills[skillID]
	return ok
}

// IsDead 仅根据当前血量判断；后续可与「复活无敌」等状态位扩展。
func (e *Entity) IsDead() bool {
	return e.Runtime.CurHP <= 0
}
