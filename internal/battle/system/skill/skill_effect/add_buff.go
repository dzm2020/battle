package skill_effect

import (
	"battle/internal/battle/config"
	"battle/internal/battle/system/buff"
	"errors"
)

// handleAddBuff 添加 Buff：走 [buff.AddBuff]，模板来自全局 [config.Tab.BuffConfigConfigByID]。
// IntParams[0]：Buff 模板 ID（uint32）。
func handleAddBuff(ctx *Context, desc *config.SkillEffectConfig) error {
	if len(desc.IntParams) < 1 || desc.IntParams[0] <= 0 {
		return errors.New("int param number must be greater than 0")
	}
	buffID := uint32(desc.IntParams[0])

	return buff.Add(ctx.World, ctx.Caster, ctx.Target, buffID)
}
