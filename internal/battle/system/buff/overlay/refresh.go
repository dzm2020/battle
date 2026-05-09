package overlay

import (
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
	"battle/internal/battle/utils"
)

func stackPolicyRefresh(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	idx := utils.FindDefIndex(bl.Buffs, desc.ID)
	if idx < 0 {
		bl.Buffs = append(bl.Buffs, new)
		log.Debug("[buff] 叠层·刷新策略：新增实例 模板编号=%d", desc.ID)
		return true
	}
	bl.Buffs[idx].DurationFrame = desc.DurationFrame
	log.Debug("[buff] 叠层·刷新策略：刷新持续时间 模板编号=%d 剩余帧=%d", desc.ID, desc.DurationFrame)
	return true
}
