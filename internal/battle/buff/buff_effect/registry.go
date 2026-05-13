package buff_effect

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
	"fmt"

	"github.com/duke-git/lancet/v2/maputil"
)

// effectHandlerFn 按模板中某一种 [BufferEffectType] 聚合处理：DoT/HoT 仍在 [TickPeriodicBuffEffects]。
type (
	handler func(ctx *Context) error
	Context struct {
		world *ecs.World
		e     ecs.Entity
		buff  *component.BuffInstance
		desc  *config.BuffConfig
	}
)

var handlers = maputil.NewConcurrentMap[config.BufferEffectType, handler](8)

func Apply(world *ecs.World, e ecs.Entity, buff *component.BuffInstance) error {
	desc, _ := config.Tab.BuffConfigConfigByID[int32(buff.BuffId)]
	if desc == nil {
		return fmt.Errorf("effect not exists")
	}
	log.Debug("[buff] 触发周期效果 实体=%v Buff编号=%d 效果类型=%d 层数=%d", e, buff.BuffId, desc.EffectType, buff.Stacks)
	if fn, ok := handlers.Get(desc.EffectType); ok && fn != nil {
		ctx := &Context{
			world: world,
			e:     e,
			buff:  buff,
			desc:  desc,
		}
		return fn(ctx)
	} else {
		return fmt.Errorf("effect not exists")
	}
}

func init() {
	handlers.Set(config.BufferEffectStatChange, handleStatChange)
	handlers.Set(config.BufferEffectControl, handleControl)
	handlers.Set(config.BufferEffectDamage, handleDamage)
	handlers.Set(config.BufferEffectHeal, handleHeal)
}
