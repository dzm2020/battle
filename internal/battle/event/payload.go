package event

import "battle/ecs"

type AddBuffRequestPayLoad struct {
	Caster, Target ecs.Entity
	BuffId         uint32
}

type RemoveBuffRequestPayLoad struct {
	Caster, Target ecs.Entity
	BuffId         uint32
}

type AddSkillRequestPayLoad struct {
	Caster, Target ecs.Entity
	SkillID        int32
}

type RemoveSkillRequestPayLoad struct {
	Caster, Target ecs.Entity
	SkillID        int32
}
