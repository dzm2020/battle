package entity_factory

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
	"errors"
	"fmt"
)

var (
	ErrNilWorld      = errors.New("unit: world is nil")
	ErrNilPlayer     = errors.New("unit: player is nil")
	ErrNoPlayerUnits = errors.New("unit: player.Units 为空或全部为 nil")
	ErrUnknownUnit   = errors.New("unit: 未知的单位模板 id")
)

func init() {

}

type builder func(w *ecs.World, desc *config.UnitConfig, coms ...ecs.Component) (ecs.Entity, error)

var (
	blueprints = make(map[string]builder)
)

// ====================== 方式一：蓝图ID创建 ======================
func register(id string, builder builder) {
	blueprints[id] = builder
}

func CreateByID(w *ecs.World, unitID int32, components ...ecs.Component) (ecs.Entity, error) {
	unitDesc := config.GetUnitConfigByID(unitID)
	if unitDesc == nil {
		return 0, fmt.Errorf("%w: %d", ErrUnknownUnit, unitID)
	}
	if builder, ok := blueprints[unitDesc.Builder]; ok {
		return builder(w, unitDesc, components...)
	} else {
		return defaultBuilder(w, unitDesc, components...)
	}
}

func defaultBuilder(w *ecs.World, desc *config.UnitConfig, components ...ecs.Component) (ecs.Entity, error) {
	e, err := spawnUnitFromConfigDesc(w, desc, components...)
	if err != nil {
		return 0, fmt.Errorf("单位模板 %s: %w", desc.ID, err)
	}
	return e, nil
}

func spawnUnitFromConfigDesc(w *ecs.World, desc *config.UnitConfig, Components ...ecs.Component) (ecs.Entity, error) {
	if err := Spawn(w, SpawnUnitOptions{
		Abilities:  desc.Ability,
		BuffDefIDs: desc.BuffDefIDs,
		Components: Components,
		Attributes: buildAttrFromConfig(desc.Stats),
	}); err != nil {
		return 0, err
	}
	return e, nil
}

func buildAttrFromConfig(attrConfigIds []int32) map[config.AttributeType]*component.Attribute {
	base := make(map[config.AttributeType]*component.Attribute)
	for _, attrRowID := range attrConfigIds {
		attrDesc := config.GetAttributeConfigByID(attrRowID)
		if attrDesc == nil {
			log.Error("属性表缺少 id=%d", attrRowID)
			continue
		}
		base[attrDesc.Type] = &component.Attribute{
			Current: int(attrDesc.InitValue),
			Max:     int(attrDesc.MaxValue),
		}
	}
	return base
}
