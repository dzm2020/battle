package overlay

import (
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
	"battle/internal/battle/system/buff/utils"
)

func stackPolicyReplace(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	idx := utils.FindDefIndex(bl.Buffs, desc.ID)
	if idx < 0 {
		bl.Buffs = append(bl.Buffs, new)
		log.Debug("[buff] 叠层·替换策略：新增实例 模板编号=%d", desc.ID)
	} else {
		bl.Buffs[idx] = new
		log.Debug("[buff] 叠层·替换策略：替换已有槽位 模板编号=%d", desc.ID)
	}
	return true
}
