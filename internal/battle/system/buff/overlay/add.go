package overlay

import (
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
	"battle/internal/battle/utils"
)

func stackPolicyAdd(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	idx := utils.FindDefIndex(bl.Buffs, desc.ID)
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
