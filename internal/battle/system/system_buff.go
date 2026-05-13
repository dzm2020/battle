package system

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/buff/buff_effect"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
)

type BuffSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.BuffList]
}

func (s *BuffSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.BuffList](w)
	log.Info("[buff] Buff 系统已初始化")
}

// Update 遍历含 [component.BuffList] 的实体：先清零并重算 StatModifiers/ControlState，再逐实例
// 聚合属性，最后递减 FramesLeft 并剔除到期实例。
func (s *BuffSystem) Update(dt float64) {
	//  遍历所有对象
	s.q.ForEach(func(e ecs.Entity, bl *component.BuffList) {
		//  清空buff产生的临时组件
		s.clearPreFrameEffect(e)
		//  倒序遍历：到期 RemoveBuff 会缩切片，正向 for-range 会跳过元素。
		for i := len(bl.Buffs) - 1; i >= 0; i-- {
			bi := bl.Buffs[i]
			//  配置检测
			desc, _ := config.Tab.BuffConfigConfigByID[int32(bi.BuffId)]
			if desc == nil {
				log.Error("[buff] 每帧轮询：缺少 Buff 配置 实体=%v Buff编号=%d", e, bi.BuffId)
				_ = buff.Remove(s.world, e, bi.BuffId)
				continue
			}

			//  触发效果
			bi.CoolDownFrame--
			if bi.CoolDownFrame <= 0 {

				if err := buff_effect.Apply(s.world, e, bi); err != nil {
					log.Error("buff effect apply error:%s", err)
				}
				//  重置冷却
				periodicFrame := desc.CoolingFrame
				bi.CoolDownFrame = periodicFrame - 1
				if bi.CoolDownFrame < 0 {
					bi.CoolDownFrame = 0
				}
			}

			//  检测生命周期
			if bi.DurationFrame > 0 {
				bi.DurationFrame--
				if bi.DurationFrame <= 0 {
					_ = buff.Remove(s.world, e, bi.BuffId)
				}
			}
		}

	})
}

func (s *BuffSystem) clearPreFrameEffect(e ecs.Entity) {
	w := s.world
	if w == nil || e == 0 {
		return
	}
	w.RemoveComponent(e, &component.BuffStatModifiers{})
	w.RemoveComponent(e, &component.BuffControlState{})
}
