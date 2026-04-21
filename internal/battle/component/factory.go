package component

import (
	"slices"

	"battle/ecs"
)

// ComponentFactory 创建零值组件实例（可由配置按字符串名挂载）；具体数值仍由模板字段或系统在后续填充。
type ComponentFactory func() ecs.Component

var componentFactories = map[string]ComponentFactory{}

// RegisterComponent 运行时注册组件工厂（同名覆盖），供配置「extraComponents」或插件扩展。
func RegisterComponent(typ string, factory ComponentFactory) {
	if typ == "" || factory == nil {
		panic("component: RegisterComponent requires non-empty typ and factory")
	}
	componentFactories[typ] = factory
}

// CreateComponent 按注册名创建一个空组件实例；未注册返回 false。
func CreateComponent(typ string) (ecs.Component, bool) {
	f, ok := componentFactories[typ]
	if !ok {
		return nil, false
	}
	return f(), true
}

// IsRegisteredComponent 是否已通过 [RegisterComponent] 注册。
func IsRegisteredComponent(typ string) bool {
	_, ok := componentFactories[typ]
	return ok
}

// RegisteredComponentNames 返回当前已注册的类型名快照（排序后），便于调试与校验。
func RegisteredComponentNames() []string {
	out := make([]string, 0, len(componentFactories))
	for k := range componentFactories {
		out = append(out, k)
	}
	slices.Sort(out)
	return out
}

func registerDefaultComponentFactories() {
	RegisterComponent("ThreatBook", func() ecs.Component { return &ThreatBook{} })
	RegisterComponent("BuffList", func() ecs.Component { return &BuffList{} })
	RegisterComponent("StatModifiers", func() ecs.Component { return &StatModifiers{} })
	RegisterComponent("ControlState", func() ecs.Component { return &ControlState{} })
	RegisterComponent("Transform2D", func() ecs.Component { return &Transform2D{} })
}

func init() {
	registerDefaultComponentFactories()
}
