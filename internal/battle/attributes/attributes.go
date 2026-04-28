package attributes

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
)

type Attribute struct {
	ID        int32            // 配置项唯一ID
	Type      config.Attribute // 属性类型
	InitValue int32            // 属性初始值
	MaxValue  int32            // 属性最大值
}

func InitFromConfig(w *ecs.World, e ecs.Entity, attrConfigIds []int32) {
	attrs := ecs.EnsureGetComponent[*component.Attributes](w, e)
	var health *component.Health
	//  初始化属性
	for _, attrRowID := range attrConfigIds {
		attrDesc := config.GetAttributeConfigByID(attrRowID)
		if attrDesc == nil {
			w.RemoveEntity(e)
			log.Error("属性表缺少 id=%d", attrRowID)
			return
		}
		key := string(attrDesc.Type)
		v := int(attrDesc.InitValue)
		attrs.Set(key, v)
		if attrDesc.Type == config.AttrHp {
			cur := int(attrDesc.InitValue)
			maxHP := int(attrDesc.MaxValue)
			if maxHP < cur {
				maxHP = cur
			}
			health = &component.Health{Current: cur, Max: maxHP}
		}
	}
	//  加入组件
	w.AddComponent(e, attrs)
	if health != nil {
		w.AddComponent(e, health)
	}
}

func InitFromStats(w *ecs.World, e ecs.Entity, stats []Attribute) {
	if len(stats) == 0 {
		return
	}
	attrs := ecs.EnsureGetComponent[*component.Attributes](w, e)
	var health *component.Health
	for _, attr := range stats {
		key := string(attr.Type)
		attrs.Set(key, int(attr.InitValue))
		if attr.Type == config.AttrHp {
			cur := int(attr.InitValue)
			maxHP := int(attr.MaxValue)
			if maxHP < cur {
				maxHP = cur
			}
			health = &component.Health{Current: cur, Max: maxHP}
		}
	}
	w.AddComponent(e, attrs)
	if health != nil {
		w.AddComponent(e, health)
	}
}
