package skill

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/target_selector"
)

// AddSkill 若配置表中存在 skillConfigID，则在 entity 上确保存在 [component.SkillSet]，
// 并追加一条 [RuntimeSkill]（已存在相同 ConfigID 时不重复追加，返回 true）。
func AddSkill(w *ecs.World, entity ecs.Entity, skillConfigID int32) bool {
	if w == nil || entity == 0 || !w.EntityExists(entity) {
		return false
	}
	cfg := config.GetSkillConfigByID(skillConfigID)
	if cfg == nil {
		return false
	}
	set := ecs.EnsureGetComponent[*component.SkillSet](w, entity)
	for _, rs := range set.Skills {
		if rs != nil && rs.ConfigID == cfg.ID {
			return true
		}
	}
	set.Skills = append(set.Skills, &component.RuntimeSkill{ConfigID: cfg.ID, CurrentCooldown: 0})
	return true
}

// RemoveSkill 从 entity 的 [SkillSet] 中删除 ConfigID 等于 skillConfigID 的条目（至少删一条匹配项）。
// 若移除后无任何技能，则移除 [SkillSet] 组件。未挂载 SkillSet 或未找到匹配项时返回 false。
func RemoveSkill(w *ecs.World, entity ecs.Entity, skillConfigID int32) bool {
	if w == nil || entity == 0 || !w.EntityExists(entity) {
		return false
	}
	c, ok := w.GetComponent(entity, &component.SkillSet{})
	if !ok {
		return false
	}
	set, ok := c.(*component.SkillSet)
	if !ok || set == nil || len(set.Skills) == 0 {
		return false
	}
	want := int(skillConfigID)
	found := false
	out := set.Skills[:0]
	for _, rs := range set.Skills {
		if rs != nil && rs.ConfigID == want {
			found = true
			continue
		}
		out = append(out, rs)
	}
	set.Skills = out
	if len(set.Skills) == 0 {
		w.RemoveComponent(entity, &component.SkillSet{})
	}
	return found
}

// ApplySkillEffects 按 [config.SkillBaseConfig.EffectIDs] 顺序执行技能效果（伤害、加 Buff 等）。
// caster：施法者；skillID：技能配置 ID。
func ApplySkillEffects(w *ecs.World, caster ecs.Entity, skillID int) {
	if w == nil || caster == 0 || !w.EntityExists(caster) {
		return
	}

	desc := config.GetSkillConfigByID(int32(skillID))

	if desc == nil {
		return
	}
	for _, eid := range desc.EffectIDs {
		effectDesc := config.GetSkillEffectConfigByID(int32(eid))
		if effectDesc == nil {
			continue
		}
		//  选取目标
		targets := target_selector.Select(w, caster, int32(effectDesc.TargetSelectID))
		//  执行效果
		for _, t := range targets {
			applySkillEffect(w, caster, t, effectDesc)
		}
	}
}
