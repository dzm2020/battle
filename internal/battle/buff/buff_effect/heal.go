package buff_effect

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/log"
	"fmt"
)

// handleHeal 治疗效果
func handleHeal(ctx *Context) error {
	world := ctx.world
	e := ctx.e
	desc := ctx.desc
	buff := ctx.buff

	if desc == nil || len(desc.Params) < 1 {
		return fmt.Errorf("[buff] 治疗效果：参数不足 实体=%v Buff编号=%d", e, buff.BuffId)
	}
	st := max(buff.Stacks, 1)
	heal := int(desc.Params[0]) * st
	if heal <= 0 {
		return fmt.Errorf("[buff] 治疗效果：跳过非正数额 实体=%v Buff编号=%d 数额=%d", e, buff.BuffId, heal)
	}
	com := ecs.EnsureGetComponent[*component.PendingHeal](world, e)
	com.Amount += heal
	log.Debug("[buff] 持续治疗：累计写入 实体=%v Buff编号=%d 本帧增量=%d 当前合计=%d", e, buff.BuffId, heal, com.Amount)
	return nil
}
