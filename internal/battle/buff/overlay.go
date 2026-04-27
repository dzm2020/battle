package buff

import (
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
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
		log.Debug("[buff] 叠层·替换策略：新增实例 模板编号=%d", desc.ID)
	} else {
		bl.Buffs[idx] = new
		log.Debug("[buff] 叠层·替换策略：替换已有槽位 模板编号=%d", desc.ID)
	}
	return true
}

func stackPolicyRefresh(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	idx := findDefIndex(bl.Buffs, desc.ID)
	if idx < 0 {
		bl.Buffs = append(bl.Buffs, new)
		log.Debug("[buff] 叠层·刷新策略：新增实例 模板编号=%d", desc.ID)
		return true
	}
	bl.Buffs[idx].DurationFrame = desc.DurationFrame
	log.Debug("[buff] 叠层·刷新策略：刷新持续时间 模板编号=%d 剩余帧=%d", desc.ID, desc.DurationFrame)
	return true
}

func stackPolicyAdd(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	idx := findDefIndex(bl.Buffs, desc.ID)
	if idx < 0 {
		bl.Buffs = append(bl.Buffs, new)
		log.Debug("[buff] 叠层·叠加策略：新增实例 模板编号=%d", desc.ID)
		return true
	}
	b := bl.Buffs[idx]
	b.Stacks++
	b.DurationFrame = desc.DurationFrame
	log.Debug("[buff] 叠层·叠加策略：当前层数=%d 模板编号=%d", b.Stacks, desc.ID)
	return true
}

func stackPolicyIgnore(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	if findDefIndex(bl.Buffs, desc.ID) >= 0 {
		log.Debug("[buff] 叠层·忽略策略：已存在相同模板 模板编号=%d，忽略本次施加", desc.ID)
		return false
	}
	bl.Buffs = append(bl.Buffs, new)
	log.Debug("[buff] 叠层·忽略策略：新增实例 模板编号=%d", desc.ID)
	return true
}
