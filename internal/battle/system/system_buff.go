package system

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// BuffSystem 递减持续时间、汇总属性并写入 [StatModifiers]/[ControlState]。
// 须在 [DamageSystem] 之前运行，以便本帧 DoT 写入的 [PendingDamage] 参与结算。
type BuffSystem struct {
	world *ecs.World
	q     *ecs.Query[*component.BuffList]
}

// NewBuffSystem 使用全局 [config.Tab.BuffConfigConfigByID] 解析 Buff 模板。
func NewBuffSystem() *BuffSystem {
	return &BuffSystem{}
}

func (s *BuffSystem) Initialize(w *ecs.World) {
	s.world = w
	s.q = ecs.NewQuery[*component.BuffList](w)
}

// Update 遍历含 [component.BuffList] 的实体：先清零并重算 StatModifiers/ControlState，再逐实例
// 聚合属性，最后递减 FramesLeft 并剔除到期实例。
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

	tab := config.Tab
	if tab.BuffConfigConfigByID == nil {
		bl.Buffs = nil
		s.stripBuffAux(e)
		s.world.RemoveComponent(e, &component.BuffList{})
		return
	}

	alive := make([]component.BuffInstance, 0, len(bl.Buffs))
	for i := range bl.Buffs {
		bi := bl.Buffs[i]
		bc, ok := tab.BuffConfigConfigByID[int32(bi.BuffId)]
		if !ok || bc == nil {
			continue
		}

		s.accumulateFromBuffConfig(bc, &bi, mods)

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

func (s *BuffSystem) accumulateFromBuffConfig(bc *config.BuffConfig, bi *component.BuffInstance, mods *component.StatModifiers) {
	st := bi.Stacks
	if st < 1 {
		st = 1
	}
	for _, sm := range bc.Modifiers {
		applyStatModifier(mods, sm, st)
	}
}

func applyStatModifier(mods *component.StatModifiers, sm config.StatModifier, stacks int) {
	delta := int(sm.Delta) * stacks
	switch sm.Stat {
	case config.AttrArmor:
		mods.ArmorDelta += delta
	case config.AttrMagicResist:
		mods.MRDelta += delta
	case config.AttrAttackDamage:
		mods.AttackDamageDelta += delta
	case config.AttrHitPermille:
		mods.HitDeltaPermille += delta
	case config.AttrDodgePermille:
		mods.DodgeDeltaPermille += delta
	case config.AttrCritRate:
		mods.CritRateDeltaPermille += delta
	case config.AttrCritDamage:
		mods.CritDamageDeltaPermille += delta
	default:
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
