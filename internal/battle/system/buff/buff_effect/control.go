package buff_effect

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/log"
	"fmt"
	"strings"
)

// handleControl 控制效果
func handleControl(ctx *Context) error {
	world := ctx.world
	e := ctx.e
	desc := ctx.desc
	buff := ctx.buff

	if desc == nil || len(desc.ParamsString) < 1 {
		return fmt.Errorf("[buff] 控制效果：参数不足 实体=%v Buff编号=%d", e, buff.BuffId)
	}
	ctrl := ecs.EnsureGetComponent[*component.BuffControlState](world, e)
	tag := strings.ToLower(strings.TrimSpace(desc.ParamsString[0]))
	switch tag {
	case "stun", "stunned":
		ctrl.Flags |= component.FlagStunned
		log.Debug("[buff] 控制效果：眩晕 实体=%v Buff编号=%d", e, buff.BuffId)
	case "silence", "silenced":
		ctrl.Flags |= component.FlagSilenced
		log.Debug("[buff] 控制效果：沉默 实体=%v Buff编号=%d", e, buff.BuffId)
	case "root", "rooted":
		ctrl.Flags |= component.FlagRooted
		log.Debug("[buff] 控制效果：禁锢 实体=%v Buff编号=%d", e, buff.BuffId)
	default:
		log.Debug("[buff] 控制效果：未识别的标签 实体=%v Buff编号=%d 标签=%s", e, buff.BuffId, tag)
	}
	return nil
}
