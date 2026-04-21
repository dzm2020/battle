package buff

import (
	"strings"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
)

// ApplyBuff 根据 [config.Tables.BuffConfigConfigByID] 向实体挂载一条 [component.BuffInstance]。
// Tab 未初始化、表中无 buffId、StackBehavior 为 ignore 且已存在同 ID 实例、或 target 无效时返回 false。
func ApplyBuff(w *ecs.World, caster, target ecs.Entity, buffId uint32) bool {
	if w == nil || target == 0 || buffId == 0 {
		return false
	}
	tab := config.Tab
	if tab.BuffConfigConfigByID == nil {
		return false
	}
	desc, ok := tab.BuffConfigConfigByID[int32(buffId)]
	if !ok || desc == nil {
		return false
	}
	//  确定buff组件
	var bl *component.BuffList
	if c, ok := w.GetComponent(target, &component.BuffList{}); ok {
		bl = c.(*component.BuffList)
	} else {
		bl = &component.BuffList{}
		w.AddComponent(target, bl)
	}
	//  创建buff
	newInst := func(stacks int) component.BuffInstance {
		tc := tickIntervalOr1(findFirstInterval(desc))
		fl := desc.DurationFrame
		if fl == 0 {
			fl = -1
		}
		return component.BuffInstance{
			BuffId:        buffId,
			Stacks:        stacks,
			FramesLeft:    fl,
			TickCountdown: tc - 1,
			Caster:        caster,
		}
	}
	//  叠加buff
	return applyStackPolicy(desc, bl, buffId, newInst)
}

// applyStackPolicy 按 [BuffConfig.StackBehavior] 合并层数 / 持续时间 / 独立槽位；不含具体属性数值（数值在 [BuffSystem] 汇总）。
// ignore 且已有同 DefID 时返回 false。
func applyStackPolicy(desc *config.BuffConfig, bl *component.BuffList, defID uint32, newInst func(int) component.BuffInstance) bool {
	max := desc.MaxStack
	if max < 1 {
		max = 1
	}
	switch strings.ToLower(strings.TrimSpace(desc.StackBehavior)) {
	case "ignore":
		if findDefIndex(bl.Buffs, defID) >= 0 {
			return false
		}
		bl.Buffs = append(bl.Buffs, newInst(1))

	case "refresh":
		idx := findDefIndex(bl.Buffs, defID)
		if idx < 0 {
			bl.Buffs = append(bl.Buffs, newInst(1))
			break
		}
		b := &bl.Buffs[idx]
		if desc.DurationFrame == 0 {
			b.FramesLeft = -1
		} else {
			b.FramesLeft = desc.DurationFrame
		}

	case "add":
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
		if desc.DurationFrame == 0 {
			b.FramesLeft = -1
		} else {
			b.FramesLeft = desc.DurationFrame
		}

	default:
		// 独立叠加：同 DefID 允许多条实例并存
		bl.Buffs = append(bl.Buffs, newInst(1))
	}
	return true
}

// findDefIndex 在缓冲表中查找首个 DefID 匹配的槽位；无则返回 -1。
func findDefIndex(buf []component.BuffInstance, id uint32) int {
	for i := range buf {
		if buf[i].BuffId == id {
			return i
		}
	}
	return -1
}

// findFirstInterval：当前 [config.BuffConfig] 不含周期性伤害配置时返回 1。
func findFirstInterval(_ *config.BuffConfig) int {
	return 1
}

// tickIntervalOr1 将 DoT/HoT 间隔下限钳制为 1 帧，避免除零或无节拍。
func tickIntervalOr1(n int) int {
	if n < 1 {
		return 1
	}
	return n
}
