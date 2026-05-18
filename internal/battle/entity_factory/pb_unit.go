package entity_factory

import (
	"errors"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/pb"
)

var ErrNilPBUnit = errors.New("entity_factory: player unit is nil")

// CreateFromPBUnit 按协议单位数据创建实体。
func CreateFromPBUnit(w *ecs.World, u *pb.PlayerUnit, extra ...ecs.Component) (ecs.Entity, error) {
	if w == nil {
		return 0, ErrNilWorld
	}
	if u == nil {
		return 0, ErrNilPBUnit
	}
	return Create(w, Spec{
		Abilities:  u.Ability,
		BuffDefIDs: u.BuffDefIDs,
		Components: extra,
		Attributes: buildAttributesFromPB(u.Stats),
	})
}

func buildAttributesFromPB(stats []pb.Attribute) map[config.AttributeType]*component.Attribute {
	base := make(map[config.AttributeType]*component.Attribute)
	for _, stat := range stats {
		base[stat.Type] = &component.Attribute{
			Current: int(stat.InitValue),
			Max:     int(stat.MaxValue),
		}
	}
	return base
}
