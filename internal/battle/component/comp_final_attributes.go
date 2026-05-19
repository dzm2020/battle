package component

import "battle/internal/battle/config"

// FinalAttributes 战斗用最终属性（基础 [Attributes] + [BuffStatModifiers]），由 [system.AttributeSystem] 每帧写入。
// 命中、护甲、暴击等战斗结算请读本组件；生命/法力的当前消耗仍写在 [Attributes]（[HealthSystem] / 施法校验）。
type FinalAttributes struct {
	Values map[config.AttributeType]int
}

func (*FinalAttributes) Component() {}
