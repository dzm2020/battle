package buff_effect

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
)

// handleBufferEffectDamage 伤害效果（周期触发时的结算也可走同一 dict；列表合并见 [tickApplyBuffEffect]）。
func handleBufferEffectDamage(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	if desc == nil || len(desc.Params) < 1 {
		log.Error("[buff] 伤害效果：参数不足 实体=%v Buff编号=%d", e, buff.BuffId)
		return
	}
	st := max(buff.Stacks, 1)
	amt := int(desc.Params[0]) * st
	if amt <= 0 {
		log.Error("[buff] 伤害效果：跳过非正数额 实体=%v Buff编号=%d 数额=%d", e, buff.BuffId, amt)
		return
	}
	com := ecs.EnsureGetComponent[*component.PendingDamageBuff](world, e)
	com.Amount += amt
	log.Debug("[buff] 持续伤害：累计写入 实体=%v Buff编号=%d 本帧增量=%d 当前合计=%d", e, buff.BuffId, amt, com.Amount)
}
