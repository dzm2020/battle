package system

import "battle/internal/battle/component"

func startSkillCooldown(su *component.SkillUser, skillID uint32, frames int) {
	if frames <= 0 {
		return
	}
	if su.CooldownRemaining == nil {
		su.CooldownRemaining = make(map[uint32]int)
	}
	su.CooldownRemaining[skillID] = frames
}
