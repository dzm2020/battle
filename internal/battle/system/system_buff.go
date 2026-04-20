package system

import (
	"battle/ecs"
	"battle/internal/battle/buff"
	"battle/internal/battle/component"
)

// BuffSystem 递减持续时间、结算 DoT/HoT、汇总属性与控制位并写入 [StatModifiers]/[ControlState]。
// 须在 [DamageSystem] 之前运行，以便本帧 DoT 写入的 [PendingDamage] 参与结算。
type BuffSystem struct {
	world     *ecs.World
	buffConfig *buff.DefinitionConfig
	q         *ecs.Query[*component.BuffList]
}

// NewBuffSystem buffConfig 可为 nil（会使用空 [buff.DefinitionConfig]），正常运行时应传入已 Register 的配置表。
func NewBuffSystem(buffConfig *buff.DefinitionConfig) *BuffSystem {
	return &BuffSystem{buffConfig: buffConfig}
}

func (s *BuffSystem) Initialize(w *ecs.World) {
	s.world = w
	if s.buffConfig == nil {
		s.buffConfig = buff.NewDefinitionConfig()
	}
	s.q = ecs.NewQuery[*component.BuffList](w)
}

// Update 遍历含 [component.BuffList] 的实体：先清零并重算 StatModifiers/ControlState，再逐实例
// 聚合属性与控制、触发 DoT/HoT，最后递减 FramesLeft 并剔除到期实例。
func (s *BuffSystem) Update(dt float64) {
	s.q.ForEach(func(e ecs.Entity, bl *component.BuffList) {
		s.tickEntity(e, bl)
	})
}

// tickEntity 处理单个实体的 Buff 缓冲表：未知 DefID 的实例被丢弃；列表空时移除 BuffList 及派生组件。
func (s *BuffSystem) tickEntity(e ecs.Entity, bl *component.BuffList) {
	mods := s.ensureStatMods(e)
	ctrl := s.ensureControl(e)
	*mods = component.StatModifiers{}
	ctrl.Flags = 0

	if len(bl.Buffs) == 0 {
		s.stripBuffAux(e)
		s.world.RemoveComponent(e, &component.BuffList{})
		return
	}

	alive := make([]component.BuffInstance, 0, len(bl.Buffs))
	for i := range bl.Buffs {
		bi := bl.Buffs[i]
		descConfig, ok := s.buffConfig.Get(bi.DefID)
		if !ok {
			continue
		}

		s.accumulateStatic(&descConfig, &bi, mods, ctrl)
		s.tickDOTHOT(e, &bi, &descConfig)

		keep := true
		if bi.FramesLeft >= 0 {
			bi.FramesLeft--
			if bi.FramesLeft <= 0 {
				keep = false
			}
		}
		if keep {
			alive = append(alive, bi)
		}
	}

	bl.Buffs = alive
	if len(bl.Buffs) == 0 {
		s.stripBuffAux(e)
		s.world.RemoveComponent(e, &component.BuffList{})
	}
}

// accumulateStatic 将本帧仍存活实例上的 StatMod/Control 效果累加到 mods、ctrl。
func (s *BuffSystem) accumulateStatic(desc *buff.DescriptorConfig, bi *component.BuffInstance, mods *component.StatModifiers, ctrl *component.ControlState) {
	st := bi.Stacks
	if st < 1 {
		st = 1
	}
	for _, ef := range desc.Effects {
		switch ef.Kind {
		case buff.EffectStatMod:
			mods.ArmorDelta += ef.ArmorDeltaPerStack * st
			mods.MRDelta += ef.MRDeltaPerStack * st
			mods.PhysicalPowerDelta += ef.PowerDeltaPerStack * st
		case buff.EffectControl:
			ctrl.Flags |= ef.Control
		default:
		}
	}
}

// tickDOTHOT 按 TickCountdown/间隔推进 DoT（[component.MergePendingDamage]）与 HoT（直接改生命）。
// 多条 DoT/HoT 共用 DescriptorConfig 内首个 Tick 间隔（与 BuffInstance.TickCountdown 一致）。
func (s *BuffSystem) tickDOTHOT(e ecs.Entity, bi *component.BuffInstance, desc *buff.DescriptorConfig) {
	interval := 1
	hasTick := false
	for _, ef := range desc.Effects {
		if ef.Kind == buff.EffectDoT || ef.Kind == buff.EffectHoT {
			hasTick = true
			interval = ef.TickIntervalFrames
			if interval < 1 {
				interval = 1
			}
			break
		}
	}
	if !hasTick {
		return
	}

	bi.TickCountdown--
	if bi.TickCountdown >= 0 {
		return
	}

	for _, ef := range desc.Effects {
		st := bi.Stacks
		if st < 1 {
			st = 1
		}
		switch ef.Kind {
		case buff.EffectDoT:
			component.MergePendingDamage(s.world, e, ef.DamagePerTick*st, ef.DamageType)
		case buff.EffectHoT:
			heal := ef.HealPerTick * st
			if heal <= 0 {
				continue
			}
			if h, ok := s.world.GetComponent(e, &component.Health{}); ok {
				hp := h.(*component.Health)
				hp.Current += heal
				if hp.Current > hp.Max {
					hp.Current = hp.Max
				}
			}
		}
	}

	bi.TickCountdown = interval - 1
	if bi.TickCountdown < 0 {
		bi.TickCountdown = 0
	}
}

// ensureStatMods 保证实体上存在 StatModifiers 指针以便原地清零与累加。
func (s *BuffSystem) ensureStatMods(e ecs.Entity) *component.StatModifiers {
	if c, ok := s.world.GetComponent(e, &component.StatModifiers{}); ok {
		return c.(*component.StatModifiers)
	}
	sm := &component.StatModifiers{}
	s.world.AddComponent(e, sm)
	return sm
}

// ensureControl 同上，对应控制位聚合。
func (s *BuffSystem) ensureControl(e ecs.Entity) *component.ControlState {
	if c, ok := s.world.GetComponent(e, &component.ControlState{}); ok {
		return c.(*component.ControlState)
	}
	cs := &component.ControlState{}
	s.world.AddComponent(e, cs)
	return cs
}

// stripBuffAux 在无 Buff 时移除本帧派生的 StatModifiers、ControlState，避免残留上一次结果。
func (s *BuffSystem) stripBuffAux(e ecs.Entity) {
	s.world.RemoveComponent(e, &component.StatModifiers{})
	s.world.RemoveComponent(e, &component.ControlState{})
}
