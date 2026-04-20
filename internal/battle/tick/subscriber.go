package tick

import "battle/internal/battle/clock"

// Subscriber 每逻辑帧回调一次；技能、Buff 等以组合方式接入，无需继承巨型 Loop。
type Subscriber interface {
	OnTick(c *clock.Clock)
}

// FuncSubscriber 函数式订阅，减少样板类型定义。
type FuncSubscriber func(c *clock.Clock)

func (f FuncSubscriber) OnTick(c *clock.Clock) { f(c) }
