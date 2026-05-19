package component

import "battle/internal/battle/config"

// ResourceConsumeEntry 单条资源消耗请求。
type ResourceConsumeEntry struct {
	Type   config.AttributeType
	Amount int
}

// ResourceConsumeQueue 待 [system.ResourceSystem] 消费的战斗资源；施法校验通过后入队。
type ResourceConsumeQueue struct {
	Entries []ResourceConsumeEntry
}

func (*ResourceConsumeQueue) Component() {}
