package entity_factory

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
	"battle/internal/battle/skill"
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

type monsterBuilder func(w *ecs.World, desc *config.UnitConfig, coms ...ecs.Component) (ecs.Entity, error)

var (
	blueprints = make(map[string]monsterBuilder)
)

// ====================== 方式一：蓝图ID创建 ======================
func register(id string, builder monsterBuilder) {
	blueprints[id] = builder
}

func CreateByID(w *ecs.World, unitID int32, coms ...ecs.Component) (ecs.Entity, error) {
	unitDesc := config.GetUnitConfigByID(unitID)
	if unitDesc == nil {
		return 0, fmt.Errorf("%w: %d", ErrUnknownUnit, unitID)
	}
	if builder, ok := blueprints[unitDesc.Builder]; ok {
		return builder(w, unitDesc, coms...)
	} else {
		return defaultBuilder(w, unitDesc, coms...)
	}
}

func defaultBuilder(w *ecs.World, desc *config.UnitConfig, components ...ecs.Component) (ecs.Entity, error) {
	e := w.CreateEntity()
	//  初始化属性
	initAttrComponentFromConfig(w, e, desc.Stats)
	//  初始化技能
	if err := skill.Add(w, e, desc.Ability...); err != nil {
		w.RemoveEntity(e)
		return 0, fmt.Errorf("unit: 初始技能  挂载失败（单位模板 %d）: %w", desc.ID, err)
	}
	//  初始化buff
	if err := buff.Add(w, e, e, desc.BuffDefIDs...); err != nil {
		w.RemoveEntity(e)
		return 0, fmt.Errorf("unit: 初始 Buff 挂载失败（单位模板 %d）: %w", desc.ID, err)
	}
	//  初始化组件
	for _, com := range components {
		w.AddComponent(e, com)
	}

	return e, nil
}

func initAttrComponentFromConfig(w *ecs.World, e ecs.Entity, attrConfigIds []int32) {
	attrs := ecs.EnsureGetComponent[*component.Attributes](w, e)
	for _, attrRowID := range attrConfigIds {
		attrDesc := config.GetAttributeConfigByID(attrRowID)
		if attrDesc == nil {
			w.RemoveEntity(e)
			log.Error("属性表缺少 id=%d", attrRowID)
			return
		}
		key := string(attrDesc.Type)
		cur := int(attrDesc.InitValue)
		maxV := int(attrDesc.MaxValue)
		if maxV < cur {
			maxV = cur
		}
		attrs.SetRange(key, cur, maxV)

	}
	//  加入组件
	w.AddComponent(e, attrs)
}
