package buff

import (
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

func init() {
	registerStackPolicy(config.BuffStackReplace, stackPolicyReplace)
	registerStackPolicy(config.BuffStackRefresh, stackPolicyRefresh)
	registerStackPolicy(config.BuffStackAdd, stackPolicyAdd)
	registerStackPolicy(config.BuffStackIgnore, stackPolicyIgnore)
	// 未配置或未知枚举值（[BuffStackUndefined]）：与 refresh 相同
	registerStackPolicy(config.BuffStackUndefined, stackPolicyRefresh)
}

// --- 叠层策略（与 [effectHandlerDict] 相同的注册方式）---

// stackPolicyFn 按 [BuffStackBehavior] 合并实例；返回 false 表示本次施加无效（如 ignore）。
type stackPolicyFn func(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool

var stackPolicyDict = make(map[config.BuffStackBehavior]stackPolicyFn)

func registerStackPolicy(behavior config.BuffStackBehavior, fn stackPolicyFn) {
	stackPolicyDict[behavior] = fn
}

// applyStackPolicy 按 [BuffConfig.StackBehavior] 合并层数 / 持续时间 / 独立槽位；不含具体属性数值（数值在 [BuffSystem] 汇总）。
// 具体逻辑见 [stackPolicyDict] / [registerStackPolicy]。[BuffStackIgnore] 且已有同 DefID 时返回 false。
func applyStackPolicy(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	fn := stackPolicyDict[desc.StackBehavior]
	if fn == nil {
		fn = stackPolicyDict[config.BuffStackUndefined]
	}
	if fn == nil {
		return stackPolicyRefresh(new, desc, bl)
	}
	return fn(new, desc, bl)
}

func stackPolicyReplace(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	idx := findDefIndex(bl.Buffs, desc.ID)
	if idx < 0 {
		bl.Buffs = append(bl.Buffs, new)
	} else {
		bl.Buffs[idx] = new
	}
	return true
}

func stackPolicyRefresh(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	idx := findDefIndex(bl.Buffs, desc.ID)
	if idx < 0 {
		bl.Buffs = append(bl.Buffs, new)
		return true
	}
	bl.Buffs[idx].DurationFrame = desc.DurationFrame
	return true
}

func stackPolicyAdd(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	idx := findDefIndex(bl.Buffs, desc.ID)
	if idx < 0 {
		bl.Buffs = append(bl.Buffs, new)
		return true
	}
	b := bl.Buffs[idx]
	b.Stacks++
	b.DurationFrame = desc.DurationFrame
	return true
}

func stackPolicyIgnore(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	if findDefIndex(bl.Buffs, desc.ID) >= 0 {
		return false
	}
	bl.Buffs = append(bl.Buffs, new)
	return true
}
