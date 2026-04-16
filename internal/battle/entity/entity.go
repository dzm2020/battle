package entity

import (
	"battle/internal/battle/attr"
	"battle/internal/battle/buff"
	"battle/internal/battle/calc"
	"battle/internal/battle/control"
	"battle/internal/battle/cooldown"
	"battle/internal/battle/geom"
)

// Entity 战斗实体：聚合属性、Buff、技能 CD、控制位。
type Entity struct {
	ID   string
	Camp int32

	Base    attr.Base
	Derived attr.Derived
	Runtime attr.Runtime

	Pos geom.Vec2

	Control control.Flags

	SkillCD     *cooldown.Book
	KnownSkills map[string]struct{}

	Buffs *buff.Manager

	// 由 Buff 帧心跳写入，供技能伤害与移动读取。
	MoveSpeedMul      float64
	OutgoingDamageMul float64
}

// New 创建 lobby 态实体；进战需 InitBattle。
func New(id string, camp int32, b attr.Base) *Entity {
	return &Entity{
		ID:                id,
		Camp:              camp,
		Base:              b,
		KnownSkills:       make(map[string]struct{}),
		Buffs:             buff.NewManager(buff.DemoRegistry()),
		MoveSpeedMul:      1,
		OutgoingDamageMul: 1,
	}
}

func (e *Entity) AttrBase() attr.Base { return e.Base }

func (e *Entity) AttrRuntime() *attr.Runtime { return &e.Runtime }

func (e *Entity) IsDead() bool { return e.Runtime.CurHP <= 0 }

func (e *Entity) GrantSkill(id string) {
	if e.KnownSkills == nil {
		e.KnownSkills = make(map[string]struct{})
	}
	e.KnownSkills[id] = struct{}{}
}

func (e *Entity) KnowsSkill(id string) bool {
	if e.KnownSkills == nil {
		return false
	}
	_, ok := e.KnownSkills[id]
	return ok
}

// InitBattle 初始化局内资源：衍生属性、血蓝满、清空控制与 Buff。
func (e *Entity) InitBattle(cal calc.Calculator) {
	if e.SkillCD == nil {
		e.SkillCD = cooldown.NewBook()
	}
	if e.Buffs == nil {
		e.Buffs = buff.NewManager(buff.DemoRegistry())
	} else {
		e.Buffs.Reset()
	}
	e.Control = 0
	e.MoveSpeedMul = 1
	e.OutgoingDamageMul = 1
	e.refreshDerived(cal, 0)
	e.Runtime.CurHP = e.Derived.MaxHP
	e.Runtime.CurMP = e.Derived.MaxMP
	e.Runtime.Shield = 0
	cal.ApplyMaxToRuntime(&e.Runtime, e.Derived)
}

// Recalculate 外部改 Base 后重算衍生（不含 Buff 平铺攻击）。
func (e *Entity) Recalculate(cal calc.Calculator) {
	e.refreshDerived(cal, 0)
	cal.ApplyMaxToRuntime(&e.Runtime, e.Derived)
}

func (e *Entity) refreshDerived(cal calc.Calculator, bonusATK int64) {
	d := cal.DerivedFromBase(e.Base)
	d.ATK += bonusATK
	if d.ATK < 0 {
		d.ATK = 0
	}
	e.Derived = d
}

// TickBuffs 每逻辑帧调用：Buff 到期/DoT/修饰器聚合并重算攻击衍生。
func (e *Entity) TickBuffs(frame uint64, cal calc.Calculator) {
	if e.Buffs == nil || cal == nil {
		return
	}
	mods := e.Buffs.Tick(frame, e)
	e.Control = mods.Control
	e.MoveSpeedMul = mods.MoveSpeedMul
	e.OutgoingDamageMul = mods.OutgoingDamageMul
	e.refreshDerived(cal, mods.BonusATK)
	cal.ApplyMaxToRuntime(&e.Runtime, e.Derived)
	buff.ClampHPToMax(&e.Runtime, e.Derived)
}

// AddBuff 由技能或其它系统为宿主添加 Buff（服务器权威入口）。
func (e *Entity) AddBuff(frame uint64, buffID string) {
	if e.Buffs == nil {
		e.Buffs = buff.NewManager(buff.DemoRegistry())
	}
	e.Buffs.Add(frame, buffID, e)
}
