package component

import "battle/ecs"

// SetSkillCastRequest 为实体设置施法请求（覆盖同实体上已有请求）。
func SetSkillCastRequest(w *ecs.World, e ecs.Entity, req *SkillCastRequest) {
	if w == nil || e == 0 || req == nil {
		return
	}
	w.RemoveComponent(e, &SkillCastRequest{})
	w.AddComponent(e, req)
}

// RequestSkillCast 玩法层快捷入口：挂本帧施法请求，目标与落点使用默认值。
func RequestSkillCast(w *ecs.World, caster ecs.Entity, skillID int32, target ecs.Entity) {
	SetSkillCastRequest(w, caster, &SkillCastRequest{
		SkillID:      skillID,
		TargetEntity: target,
	})
}
