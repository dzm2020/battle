package buff

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// tickIntervalOr1 将 DoT/HoT 间隔下限钳制为 1 帧，避免除零或无节拍。
func tickIntervalOr1(n int) int {
	if n < 1 {
		return 1
	}
	return n
}

// ApplyBuff 在实体上施加（或叠层）指定 DefID：先查 [DefinitionConfig]，再按 [DescriptorConfig.Policy]
// 修改或追加 [component.BuffList].Buffs；新建实例时会初始化 FramesLeft、TickCountdown（取首个 DoT/HoT 间隔）。
// 返回 false 表示 defID 未注册。
func ApplyBuff(w *ecs.World, buffConfig *DefinitionConfig, e ecs.Entity, defID uint32) bool {
	desc, ok := buffConfig.Get(defID)
	if !ok {
		return false
	}

	var bl *component.BuffList
	if c, ok := w.GetComponent(e, &component.BuffList{}); ok {
		bl = c.(*component.BuffList)
	} else {
		bl = &component.BuffList{}
		w.AddComponent(e, bl)
	}

	max := desc.MaxStacks
	if max < 1 {
		max = 1
	}

	newInst := func(stacks int) component.BuffInstance {
		tc := tickIntervalOr1(findFirstInterval(&desc))
		return component.BuffInstance{
			DefID:         defID,
			Stacks:        stacks,
			FramesLeft:    desc.DurationFrames,
			TickCountdown: tc - 1,
		}
	}

	switch desc.Policy {
	case StackIndependent:
		bl.Buffs = append(bl.Buffs, newInst(1))

	case StackRefresh:
		idx := findDefIndex(bl.Buffs, defID)
		if idx < 0 {
			bl.Buffs = append(bl.Buffs, newInst(1))
			break
		}
		b := &bl.Buffs[idx]
		b.FramesLeft = desc.DurationFrames

	case StackMerge:
		idx := findDefIndex(bl.Buffs, defID)
		if idx < 0 {
			bl.Buffs = append(bl.Buffs, newInst(1))
			break
		}
		b := &bl.Buffs[idx]
		b.Stacks++
		if b.Stacks > max {
			b.Stacks = max
		}
		b.FramesLeft = desc.DurationFrames
	default:
		bl.Buffs = append(bl.Buffs, newInst(1))
	}
	return true
}

// findDefIndex 在缓冲表中查找首个 DefID 匹配的槽位；无则返回 -1。
func findDefIndex(buf []component.BuffInstance, id uint32) int {
	for i := range buf {
		if buf[i].DefID == id {
			return i
		}
	}
	return -1
}

// findFirstInterval 取 Effects 里首个 DoT/HoT 的 TickIntervalFrames，用于初始化 BuffInstance.TickCountdown。
func findFirstInterval(desc *DescriptorConfig) int {
	for _, ef := range desc.Effects {
		switch ef.Kind {
		case EffectDoT, EffectHoT:
			return tickIntervalOr1(ef.TickIntervalFrames)
		}
	}
	return 1
}
