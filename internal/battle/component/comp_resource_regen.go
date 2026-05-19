package component

import "battle/internal/battle/config"

// ResourceRegen 战斗资源每帧自然恢复量（写入 [Attributes] 的 Current，不超过 Max）。
type ResourceRegen struct {
	PerFrame map[config.AttributeType]int
}

func (*ResourceRegen) Component() {}
