package buff

import (
	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
	"slices"
)

// AddBuff 根据 [config.Tab.BuffConfigConfigByID] 向实体挂载一条 [component.BuffInstance]。
// Tab 未初始化、表中无 buffId、StackBehavior 为 ignore 且已存在同 ID 实例、或 target 无效时返回 false。
func AddBuff(w *ecs.World, caster, target ecs.Entity, buffId uint32) bool {
	if w == nil || target == 0 || buffId == 0 {
		log.Debug("[buff] 添加 Buff 跳过：目标或 Buff 编号无效 目标=%v Buff编号=%d", target, buffId)
		return false
	}
	tab := config.Tab
	if tab == nil || tab.BuffConfigConfigByID == nil {
		log.Debug("[buff] 添加 Buff 跳过：配置表未就绪")
		return false
	}
	desc, ok := tab.BuffConfigConfigByID[int32(buffId)]
	if !ok || desc == nil {
		log.Debug("[buff] 添加 Buff 跳过：表中无 Buff 定义 Buff编号=%d", buffId)
		return false
	}
	com := ecs.EnsureGetComponent[*component.BuffList](w, target)
	newBuf := newBuffInstance(caster, buffId, 1)
	if newBuf == nil {
		log.Debug("[buff] 添加 Buff 跳过：创建 Buff 实例失败 Buff编号=%d", buffId)
		return false
	}
	if !applyStackPolicy(newBuf, desc, com) {
		log.Debug("[buff] 添加 Buff 跳过：叠层策略拒绝 叠层行为=%d Buff编号=%d 目标=%v", desc.StackBehavior, buffId, target)
		return false
	}
	stacks := newBuf.Stacks
	if idx := findDefIndex(com.Buffs, buffId); idx >= 0 {
		stacks = com.Buffs[idx].Stacks
	}
	log.Info("[buff] 添加 Buff 成功 施法者=%v 目标=%v Buff编号=%d 层数=%d", caster, target, buffId, stacks)
	return true
}

// RemoveBuff 从列表中移除指定 Buff 定义；若列表为空则重置 [component.BuffList] 组件。
func RemoveBuff(w *ecs.World, e ecs.Entity, bl *component.BuffList, buffId uint32) {
	if w == nil || bl == nil {
		return
	}
	idx := findDefIndex(bl.Buffs, buffId)
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

// Tick 单实体 Buff 列表的一帧更新：清零派生组件、执行周期效果、递减持续时间并移除到期实例。
func Tick(w *ecs.World, e ecs.Entity, bl *component.BuffList) {
	if w == nil || bl == nil || e == 0 {
		return
	}
	stripBuffAux(w, e)

	if len(bl.Buffs) == 0 {
		return
	}

	for i := range bl.Buffs {
		bi := bl.Buffs[i]
		desc, ok := config.Tab.BuffConfigConfigByID[int32(bi.BuffId)]
		if !ok || desc == nil {
			log.Debug("[buff] 每帧轮询：缺少 Buff 配置 实体=%v Buff编号=%d", e, bi.BuffId)
			continue
		}

		tickApplyBuffEffect(w, e, desc, bi)

		if bi.DurationFrame > 0 {
			bi.DurationFrame--
			if bi.DurationFrame <= 0 {
				RemoveBuff(w, e, bl, bi.BuffId)
			}
		}
	}
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

func tickApplyBuffEffect(w *ecs.World, e ecs.Entity, desc *config.BuffConfig, buff *component.BuffInstance) {
	if buff == nil || e == 0 {
		return
	}
	buff.CoolDownFrame--
	if buff.CoolDownFrame >= 0 {
		return
	}
	log.Debug("[buff] 触发周期效果 实体=%v Buff编号=%d 效果类型=%d 层数=%d", e, buff.BuffId, desc.EffectType, buff.Stacks)
	applyEffect(w, e, buff, desc)
	periodicFrame := desc.CoolingFrame
	buff.CoolDownFrame = periodicFrame - 1
	if buff.CoolDownFrame < 0 {
		buff.CoolDownFrame = 0
	}
}

func stripBuffAux(w *ecs.World, e ecs.Entity) {
	if w == nil || e == 0 {
		return
	}
	w.RemoveComponent(e, &component.StatModifiers{})
	w.RemoveComponent(e, &component.ControlState{})
	w.RemoveComponent(e, &component.PendingDamageBuff{})
	w.RemoveComponent(e, &component.PendingHealBuff{})
}
