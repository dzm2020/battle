package calc

import "battle/internal/battle/attr"

// PhysicalHit 普攻向物理伤害（简化版）：攻击 × outgoingMul 再套目标物免。
func PhysicalHit(attacker attr.Derived, target attr.Derived, outgoingMul float64) int64 {
	if outgoingMul <= 0 {
		outgoingMul = 1
	}
	raw := float64(attacker.ATK) * outgoingMul
	mit := target.PhysMitigation
	if mit < 0 {
		mit = 0
	}
	if mit > 0.95 {
		mit = 0.95
	}
	dmg := raw * (1 - mit)
	if dmg < 1 {
		return 1
	}
	return int64(dmg)
}
