package component

import "battle/internal/battle/config"

type BuffStatModifiers struct {
	Modifiers map[config.AttributeType]int32
}

func (*BuffStatModifiers) Component() {}
