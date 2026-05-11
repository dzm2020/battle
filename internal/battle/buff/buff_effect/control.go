package buff_effect

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
	"battle/internal/battle/utils"
	"strings"
)

// handleBufferEffectControl 控制效果
func handleBufferEffectControl(world *ecs.World, e ecs.Entity, buff *component.BuffInstance, desc *config.BuffConfig) {
	if desc == nil || len(desc.ParamsString) < 1 {
		log.Error("[buff] 控制效果：参数不足 实体=%v Buff编号=%d", e, buff.BuffId)
		return
	}
	ctrl := ecs.EnsureGetComponent[*component.ControlState](world, e)
	tag := strings.ToLower(strings.TrimSpace(desc.ParamsString[0]))
	switch tag {
	case "stun", "stunned":
		ctrl.Flags |= utils.FlagStunned
		log.Debug("[buff] 控制效果：眩晕 实体=%v Buff编号=%d", e, buff.BuffId)
	case "silence", "silenced":
		ctrl.Flags |= utils.FlagSilenced
		log.Debug("[buff] 控制效果：沉默 实体=%v Buff编号=%d", e, buff.BuffId)
	case "root", "rooted":
		ctrl.Flags |= utils.FlagRooted
		log.Debug("[buff] 控制效果：禁锢 实体=%v Buff编号=%d", e, buff.BuffId)
	default:
		log.Debug("[buff] 控制效果：未识别的标签 实体=%v Buff编号=%d 标签=%s", e, buff.BuffId, tag)
	}
}
