package skill_effect

import (
	"battle/internal/battle/buff"
	"battle/internal/battle/config"
	"errors"
)

// handleAddBuff 添加 Buff：走 [buff.AddBuff]，模板来自全局 [config.Tab.BuffConfigConfigByID]。
// IntParams[0]：Buff 模板 ID（uint32）。
func handleAddBuff(ctx *Context, desc *config.SkillEffectConfig) error {
	if len(desc.IntParams) < 1 || desc.IntParams[0] <= 0 {
		return errors.New("int param number must be greater than 0")
	}
	buffID := uint32(desc.IntParams[0])

	return buff.Add(ctx.Word, ctx.Caster, ctx.Target, buffID)
}
