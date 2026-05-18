package buff_effect

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
	"fmt"
)

// handleBufferEffectStatChange 属性修改效果
func handleStatChange(ctx *Context) error {
	world := ctx.world
	e := ctx.e
	desc := ctx.desc
	buff := ctx.buff

	mods := ecs.EnsureGetComponent[*component.BuffStatModifiers](world, e)
	if len(desc.ParamsString) < 1 || len(desc.Params) < 1 {
		return fmt.Errorf("[buff] 属性变更效果：参数不足 实体=%v Buff编号=%d", e, buff.BuffId)
	}
	key := desc.ParamsString[0]
	delta := int(desc.Params[0]) * buff.Stacks

	if mods.Modifiers == nil {
		mods.Modifiers = make(map[config.AttributeType]int32)
	}
	mods.Modifiers[key] += int32(delta)

	log.Info("[buff] 属性变更效果 实体=%v Buff编号=%d 键=%s mods=%+v", e, buff.BuffId, key, mods)
	return nil
}
