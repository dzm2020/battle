// Package entity_factory 从配置表或 PB 组装战斗实体（挂组件、初始技能与 Buff）。
//
// 边界约定：
//   - 本包是战斗内「出生装配」的唯一入口；允许在 System 外执行 CreateEntity 与初始 skill/buff 挂载。
//   - 运行时施法、战斗中获得 Buff 等由 system / skill / buff 处理，不应再调用本包。
//   - 空间落点（格子）由 room_bootstrap / SpawnSystem + land.Grid 负责，本包不处理坐标。
package entity_factory

import (
	"battle/internal/battle/system/buff"
	"battle/internal/battle/system/skill"
	"fmt"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// Spec 出生装配参数：属性、初始技能/Buff、额外组件。
type Spec struct {
	Attributes map[config.AttributeType]*component.Attribute
	Abilities  []int32
	BuffDefIDs []uint32
	Components []ecs.Component
}

// Create 按 Spec 创建实体并完成初始装配；失败时回滚删除实体。
func Create(w *ecs.World, spec Spec) (ecs.Entity, error) {
	if w == nil {
		return 0, ErrNilWorld
	}
	e := w.CreateEntity()
	attrs := ecs.EnsureGetComponent[*component.Attributes](w, e)
	attrs.Base = spec.Attributes
	if err := attachInitialLoadout(w, e, spec.Abilities, spec.BuffDefIDs); err != nil {
		w.RemoveEntity(e)
		return 0, err
	}
	for _, com := range spec.Components {
		w.AddComponent(e, com)
	}
	return e, nil
}

func attachInitialLoadout(w *ecs.World, e ecs.Entity, abilities []int32, buffDefIDs []uint32) error {
	if err := skill.Add(w, e, abilities...); err != nil {
		return fmt.Errorf("entity_factory: 初始技能: %w", err)
	}
	if err := buff.Add(w, e, e, buffDefIDs...); err != nil {
		return fmt.Errorf("entity_factory: 初始 Buff: %w", err)
	}
	return nil
}
