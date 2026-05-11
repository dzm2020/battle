package skill_effect

import (
	"battle/ecs"
	"battle/internal/battle/config"
	"battle/internal/battle/event"
)

// handleSkillEffectAddBuff 添加 Buff：走 [buff.AddBuff]，模板来自全局 [config.Tab.BuffConfigConfigByID]。
// IntParams[0]：Buff 模板 ID（uint32）。
func handleSkillEffectAddBuff(w *ecs.World, caster, target ecs.Entity, eff *config.SkillEffectConfig) {
	if len(eff.IntParams) < 1 || eff.IntParams[0] <= 0 {
		return
	}
	buffID := uint32(eff.IntParams[0])
	
	w.EmitEvent(ecs.Event{
		Kind: event.KindAddBuffRequest,
		Payload: event.AddBuffRequestPayLoad{
			Caster: caster,
			Target: target,
			BuffId: buffID,
		},
	})
}
