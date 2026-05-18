package entity_factory

import (
	"errors"
	"fmt"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
)

var (
	ErrNilWorld    = errors.New("entity_factory: world is nil")
	ErrUnknownUnit = errors.New("entity_factory: 未知的单位模板 id")
)

type configBuilder func(w *ecs.World, desc *config.UnitConfig, extra ...ecs.Component) (ecs.Entity, error)

var configBuilders = make(map[string]configBuilder)

// RegisterConfigBuilder 为指定 UnitConfig.Builder 注册自定义装配逻辑。
func RegisterConfigBuilder(name string, b configBuilder) {
	if b == nil {
		panic("entity_factory: RegisterConfigBuilder with nil builder")
	}
	configBuilders[name] = b
}

// CreateByConfigID 按单位表 ID 创建实体（读 [config.UnitConfig]）。
func CreateByConfigID(w *ecs.World, unitID int32, extra ...ecs.Component) (ecs.Entity, error) {
	if w == nil {
		return 0, ErrNilWorld
	}
	desc := config.GetUnitConfigByID(unitID)
	if desc == nil {
		return 0, fmt.Errorf("%w: %d", ErrUnknownUnit, unitID)
	}
	if b, ok := configBuilders[desc.Builder]; ok {
		return b(w, desc, extra...)
	}
	return createDefaultFromConfig(w, desc, extra...)
}

func createDefaultFromConfig(w *ecs.World, desc *config.UnitConfig, extra ...ecs.Component) (ecs.Entity, error) {
	e, err := Create(w, Spec{
		Abilities:  desc.Ability,
		BuffDefIDs: desc.BuffDefIDs,
		Components: extra,
		Attributes: buildAttributesFromConfigRows(desc.Stats),
	})
	if err != nil {
		return 0, fmt.Errorf("单位模板 %s: %w", desc.ID, err)
	}
	return e, nil
}

func buildAttributesFromConfigRows(attrConfigIds []int32) map[config.AttributeType]*component.Attribute {
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
