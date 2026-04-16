package calc

import (
	"math"

	"battle/internal/battle/attr"
)

// Calculator 属性计算抽象：Entity 只依赖接口，便于第 7 天 Buff、配表替换实现。
type Calculator interface {
	DerivedFromBase(b attr.Base) attr.Derived
	// ApplyMaxToRuntime 根据衍生上限裁剪 CurHP/CurMP，防止换装备/降级后出现超上限。
	ApplyMaxToRuntime(rt *attr.Runtime, d attr.Derived)
}

// DefaultCalculator 第 2 天固定公式实现；系数集中在本文件，方便日后迁到配置表。
type DefaultCalculator struct{}

func (DefaultCalculator) DerivedFromBase(b attr.Base) attr.Derived {
	l := float64(b.Level)
	str := float64(b.STR)
	agi := float64(b.AGI)
	intel := float64(b.INT)
	vit := float64(b.VIT)

	maxHP := int64(100 + 20*l + 15*vit)
	maxMP := int64(50 + 5*l + 8*intel)
	atk := int64(3*b.Level + 2*b.STR)
	def := int64(b.Level + b.VIT)

	critRate := 0.05 + 0.001*agi
	critRate = clamp(critRate, 0, 0.5)

	critDmg := 1.5 + 0.01*str
	critDmg = clamp(critDmg, 1.25, 3.0)

	defF := float64(def)
	mitigation := defF / (defF + 500)
	mitigation = clamp(mitigation, 0, 0.95)

	return attr.Derived{
		MaxHP:          maxHP,
		MaxMP:          maxMP,
		ATK:            atk,
		DEF:            def,
		CritRate:       critRate,
		CritDamage:     critDmg,
		PhysMitigation: mitigation,
	}
}

func (DefaultCalculator) ApplyMaxToRuntime(rt *attr.Runtime, d attr.Derived) {
	if rt.CurHP > d.MaxHP {
		rt.CurHP = d.MaxHP
	}
	if rt.CurMP > d.MaxMP {
		rt.CurMP = d.MaxMP
	}
	if rt.CurHP < 0 {
		rt.CurHP = 0
	}
	if rt.CurMP < 0 {
		rt.CurMP = 0
	}
	if rt.Shield < 0 {
		rt.Shield = 0
	}
}

func clamp(x, lo, hi float64) float64 {
	return math.Min(hi, math.Max(lo, x))
}
