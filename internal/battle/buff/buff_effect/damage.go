package buff_effect

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/log"
	"fmt"
)

// handleDamage 伤害效果（周期触发时的结算也可走同一 dict；列表合并见 [tickApplyBuffEffect]）。
func handleDamage(ctx *Context) error {
	world := ctx.world
	e := ctx.e
	desc := ctx.desc
	buff := ctx.buff

	if desc == nil || len(desc.Params) < 1 {
		return fmt.Errorf("[buff] 伤害效果：参数不足 实体=%v Buff编号=%d", e, buff.BuffId)
	}
	st := max(buff.Stacks, 1)
	amt := int(desc.Params[0]) * st
	if amt <= 0 {
		return fmt.Errorf("[buff] 伤害效果：跳过非正数额 实体=%v Buff编号=%d 数额=%d", e, buff.BuffId, amt)
	}

	damageQueue := ecs.EnsureGetComponent[*component.DamageQueue](world, e)
	damageQueue.Add(&component.PendingDamage{
		Source:    buff.Caster,
		RawDamage: float64(amt),
		Type:      component.DamageMagic,
	})

	log.Debug("[buff] 持续伤害：累计写入 实体=%v Buff编号=%d 本帧增量=%d ", e, buff.BuffId, amt)
	return nil
}
