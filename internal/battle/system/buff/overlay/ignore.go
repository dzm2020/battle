package overlay

import (
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
	"battle/internal/battle/system/buff/buff_util"
)

func stackPolicyIgnore(new *component.BuffInstance, desc *config.BuffConfig, bl *component.BuffList) bool {
	if buff_util.FindDefIndex(bl.Buffs, desc.ID) >= 0 {
		log.Debug("[buff] 叠层·忽略策略：已存在相同模板 模板编号=%d，忽略本次施加", desc.ID)
		return false
	}
	bl.Buffs = append(bl.Buffs, new)
	log.Debug("[buff] 叠层·忽略策略：新增实例 模板编号=%d", desc.ID)
	return true
}
