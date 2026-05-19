package skill_effect

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"errors"
)

// handleDamage 造成伤害：写入目标的 [component.DamageQueue]，由 [system.DamageSystem] 结算。
// IntParams[0]：伤害量；IntParams[1] 可选：[component.DamageType]（0 物 / 1 法 / 2 真），缺省为物理。
func handleDamage(ctx *Context, desc *config.SkillEffectConfig) error {
	if len(desc.IntParams) < 1 {
		return errors.New("int param number must be greater than 0")
	}
	amt := desc.IntParams[0]
	if amt <= 0 {
		return errors.New("int param number must be greater than 0")
	}

	dmgType := component.DamagePhysical
	if len(desc.IntParams) >= 2 {
		switch desc.IntParams[1] {
		case int(component.DamagePhysical):
			dmgType = component.DamagePhysical
		case int(component.DamageMagic):
			dmgType = component.DamageMagic
		case int(component.DamageTrue):
			dmgType = component.DamageTrue
		}
	}

	q := ecs.EnsureGetComponent[*component.DamageQueue](ctx.World, ctx.Target)
	component.DamageQueueAppend(q, &component.PendingDamage{
		Source:    ctx.Caster,
		RawDamage: float64(amt),
		Type:      dmgType,
	})

	return nil
}
