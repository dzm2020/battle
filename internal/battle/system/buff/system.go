package buff

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/event"
	"battle/internal/battle/log"
	"battle/internal/battle/system/buff/buff_effect"
	"battle/internal/battle/system/buff/overlay"
	"battle/internal/battle/utils"
	"slices"
)

type System struct {
	world *ecs.World
	q     *ecs.Query[*component.BuffList]
}

func (s *System) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.BuffList](w)
	log.Info("[buff] Buff 系统已初始化")

	s.world.Subscribe(event.KindAddBuffRequest, func(e ecs.Event) {
		payload := e.Payload.(*event.AddBuffRequestPayLoad)
		s.AddBuff(payload.Caster, payload.Target, payload.BuffId)
	})

	s.world.Subscribe(event.KindRemoveBuffRequest, func(e ecs.Event) {
		payload := e.Payload.(*event.RemoveBuffRequestPayLoad)
		s.RemoveBuff(payload.Target, payload.BuffId)
	})
}

// Update 遍历含 [component.BuffList] 的实体：先清零并重算 StatModifiers/ControlState，再逐实例
// 聚合属性，最后递减 FramesLeft 并剔除到期实例。
func (s *System) Update(dt float64) {
	//  遍历所有对象
	s.q.ForEach(func(e ecs.Entity, bl *component.BuffList) {
		//  清空buff产生的组件
		s.stripBuffAux(e)
		//  倒序遍历：到期 RemoveBuff 会缩切片，正向 for-range 会跳过元素。
		for i := len(bl.Buffs) - 1; i >= 0; i-- {
			bi := bl.Buffs[i]
			//  配置检测
			desc, _ := config.Tab.BuffConfigConfigByID[int32(bi.BuffId)]
			if desc == nil {
				log.Error("[buff] 每帧轮询：缺少 Buff 配置 实体=%v Buff编号=%d", e, bi.BuffId)
				s.RemoveBuff(e, bi.BuffId)
				continue
			}

			//  触发效果
			bi.CoolDownFrame--
			if bi.CoolDownFrame <= 0 {

				buff_effect.Apply(s.world, e, bi)
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
					s.RemoveBuff(e, bi.BuffId)
				}
			}
		}

	})
}

func (s *System) AddBuff(caster, target ecs.Entity, buffId uint32) bool {
	w := s.world
	if target == 0 {
		log.Error("[buff] 添加 Buff 跳过：目标或 Buff 编号无效 目标=%v Buff编号=%d", target, buffId)
		return false
	}
	tab := config.Tab
	desc, ok := tab.BuffConfigConfigByID[int32(buffId)]
	if !ok || desc == nil {
		log.Error("[buff] 添加 Buff 跳过：表中无 Buff 定义 Buff编号=%d", buffId)
		return false
	}
	//  创建buff示例
	com := ecs.EnsureGetComponent[*component.BuffList](w, target)
	newBuf := newBuffInstance(caster, buffId, 1)
	if newBuf == nil {
		log.Error("[buff] 添加 Buff 跳过：创建 Buff 实例失败 Buff编号=%d", buffId)
		return false
	}
	//  叠加
	if !overlay.Apply(newBuf, desc, com) {
		log.Debug("[buff] 添加 Buff 跳过：叠层策略拒绝 叠层行为=%d Buff编号=%d 目标=%v", desc.StackBehavior, buffId, target)
		return false
	}
	stacks := newBuf.Stacks
	if idx := utils.FindDefIndex(com.Buffs, buffId); idx >= 0 {
		stacks = com.Buffs[idx].Stacks
	}
	log.Info("[buff] 添加 Buff 成功 施法者=%v 目标=%v Buff编号=%d 层数=%d", caster, target, buffId, stacks)
	return true
}

// RemoveBuff 从列表中移除指定 Buff 定义；若列表为空则重置 [component.BuffList] 组件。
func (s *System) RemoveBuff(e ecs.Entity, buffId uint32) {
	w := s.world
	c, _ := w.GetComponent(e, &component.BuffList{})
	if w == nil || c == nil {
		return
	}
	bl := c.(*component.BuffList)

	idx := utils.FindDefIndex(bl.Buffs, buffId)
	if idx < 0 {
		log.Debug("[buff] 移除 Buff：槽位不存在 实体=%v Buff编号=%d", e, buffId)
		return
	}
	log.Info("[buff] 移除 Buff 实体=%v Buff编号=%d 移除后剩余实例数=%d", e, buffId, len(bl.Buffs)-1)
	bl.Buffs = slices.Delete(bl.Buffs, idx, idx+1)

	if len(bl.Buffs) == 0 {
		w.RemoveComponent(e, &component.BuffList{})
		w.AddComponent(e, &component.BuffList{})
	}
}

func (s *System) stripBuffAux(e ecs.Entity) {
	w := s.world
	if w == nil || e == 0 {
		return
	}
	w.RemoveComponent(e, &component.StatModifiers{})
	w.RemoveComponent(e, &component.ControlState{})
	w.RemoveComponent(e, &component.PendingDamageBuff{})
	w.RemoveComponent(e, &component.PendingHealBuff{})
}

func newBuffInstance(caster ecs.Entity, buffId uint32, stacks int) *component.BuffInstance {
	tab := config.Tab
	desc, ok := tab.BuffConfigConfigByID[int32(buffId)]
	if !ok || desc == nil {
		return nil
	}
	return &component.BuffInstance{
		BuffId:        buffId,
		Stacks:        stacks,
		DurationFrame: desc.DurationFrame,
		Caster:        caster,
	}
}
