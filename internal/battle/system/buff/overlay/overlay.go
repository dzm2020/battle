package overlay

import (
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// --- 叠层策略（与 [effectHandlerDict] 相同的注册方式）---

// stackPolicyFn 按 [BuffStackBehavior] 合并实例；返回 false 表示本次施加无效（如 ignore）。
type stackPolicyFn func(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool

var stackPolicyDict = make(map[config.BuffStackBehavior]stackPolicyFn)

func registerStackPolicy(behavior config.BuffStackBehavior, fn stackPolicyFn) {
	stackPolicyDict[behavior] = fn
}

func Apply(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	fn := stackPolicyDict[desc.StackBehavior]
	if fn == nil {
		fn = stackPolicyDict[config.BuffStackUndefined]
	}
	if fn == nil {
		return stackPolicyRefresh(new, desc, bl)
	}
	return fn(new, desc, bl)
}

func init() {
	registerStackPolicy(config.BuffStackReplace, stackPolicyReplace)
	registerStackPolicy(config.BuffStackRefresh, stackPolicyRefresh)
	registerStackPolicy(config.BuffStackAdd, stackPolicyAdd)
	registerStackPolicy(config.BuffStackIgnore, stackPolicyIgnore)
	registerStackPolicy(config.BuffStackUndefined, stackPolicyRefresh)
}
