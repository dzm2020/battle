package entity_factory

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/pb"
	"battle/internal/battle/skill"
	"fmt"
)

// SpawnUnitOptions 属性组件就绪后要挂载的技能、Buff 与额外 ECS 组件；后续可在此增加字段（如被动 id、阵营修正等）。
type SpawnUnitOptions struct {
	Attributes map[config.AttributeType]*component.Attribute
	Abilities  []int32
	BuffDefIDs []uint32
	Camp       component.SideType
	Components []ecs.Component
}

func Spawn(w *ecs.World, option SpawnUnitOptions) (ecs.Entity, error) {
	e := w.CreateEntity()
	attrsComponent := ecs.EnsureGetComponent[*component.Attributes](w, e)
	attrsComponent.Base = option.Attributes
	if err := skill.Add(w, e, option.Abilities...); err != nil {
		w.RemoveEntity(e)
		return 0, fmt.Errorf("unit: 初始技能 挂载失败: %w", err)
	}
	if err := buff.Add(w, e, e, option.BuffDefIDs...); err != nil {
		w.RemoveEntity(e)
		return 0, fmt.Errorf("unit: 初始 Buff 挂载失败: %w", err)
	}
	for _, com := range option.Components {
		w.AddComponent(e, com)
	}
	return e, nil
}

func buildAttrFromPB(stats []pb.Attribute) map[config.AttributeType]*component.Attribute {
	base := make(map[config.AttributeType]*component.Attribute)
	for _, stat := range stats {
		base[stat.Type] = &component.Attribute{
			Current: int(stat.InitValue),
			Max:     int(stat.MaxValue),
		}
	}
	return base
}

func spawnUnitFromPBUnit(w *ecs.World, unit *pb.PlayerUnit, Components ...ecs.Component) (ecs.Entity, error) {
	e := w.CreateEntity()
	if err := spawn(w, e, SpawnUnitOptions{
		Abilities:  unit.Ability,
		BuffDefIDs: unit.BuffDefIDs,
		Components: Components,
		Attributes: buildAttrFromPB(unit.Stats),
	}); err != nil {
		return 0, err
	}
	return e, nil
}
