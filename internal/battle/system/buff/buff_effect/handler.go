package buff_effect

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
)

// effectHandlerFn 按模板中某一种 [BufferEffectType] 聚合处理：DoT/HoT 仍在 [TickPeriodicBuffEffects]。
type effectHandlerFn func(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig)

var effectHandlerDict = make(map[config.BufferEffectType]effectHandlerFn)

func registerEffectHandler(typ config.BufferEffectType, fn effectHandlerFn) {
	effectHandlerDict[typ] = fn
}

func Apply(world *ecs.World, e ecs.Entity, buff *component.BuffInstance) {
	desc, _ := config.Tab.BuffConfigConfigByID[int32(buff.BuffId)]
	if desc == nil {
		log.Error("[buff] 每帧轮询：缺少 Buff 配置 实体=%v Buff编号=%d", e, buff.BuffId)
		return
	}

	log.Debug("[buff] 触发周期效果 实体=%v Buff编号=%d 效果类型=%d 层数=%d", e, buff.BuffId, desc.EffectType, buff.Stacks)

	if fn := effectHandlerDict[desc.EffectType]; fn != nil {
		fn(world, e, buff, desc)
	} else {
		log.Error("[buff] 未注册的效果类型 实体=%v Buff编号=%d 效果类型=%d", e, buff.BuffId, desc.EffectType)
	}
}

func init() {
	registerEffectHandler(config.BufferEffectStatChange, handleBufferEffectStatChange)
	registerEffectHandler(config.BufferEffectControl, handleBufferEffectControl)
	registerEffectHandler(config.BufferEffectDamage, handleBufferEffectDamage)
	registerEffectHandler(config.BufferEffectHeal, handleBufferEffectHeal)
}
