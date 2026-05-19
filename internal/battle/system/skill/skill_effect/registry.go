package skill_effect

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"fmt"

	"github.com/duke-git/lancet/v2/maputil"
)

type Context struct {
	World          *ecs.World
	Caster, Target ecs.Entity
	EffectId       int32
}

type handler func(ctx *Context, desc *config.SkillEffectConfig) error

var (
	handlers = maputil.NewConcurrentMap[config.EffectType, handler](4)
)

func Apply(ctx *Context) error {
	w := ctx.World
	target := ctx.Target
	if w == nil || target == 0 || !w.EntityExists(target) {
		return fmt.Errorf("entity not exists")
	}
	effectDesc := config.GetSkillEffectConfigByID(ctx.EffectId)
	if effectDesc == nil {
		return fmt.Errorf("effect not exists")
	}
	if fn, ok := handlers.Get(effectDesc.EffectType); ok && fn != nil {
		return fn(ctx, effectDesc)
	} else {
		return fmt.Errorf("effect not exists")
	}
}

func init() {
	handlers.Set(config.EffectDamage, handleDamage)
	handlers.Set(config.EffectAddBuff, handleAddBuff)
}
