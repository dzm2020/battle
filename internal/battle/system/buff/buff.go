package buff

import (
	"battle/internal/battle/system/buff/buff_util"
	"battle/internal/battle/system/buff/overlay"
	"fmt"

	"battle/ecs"
	"battle/internal/battle/component"
	"battle/internal/battle/config"
	"battle/internal/battle/log"
	"slices"
)

func Add(w *ecs.World, caster, target ecs.Entity, buffId ...uint32) error {
	for _, id := range buffId {
		if err := AddOne(w, caster, target, id); err != nil {
			return err
		}
	}
	return nil
}

func AddOne(w *ecs.World, caster, target ecs.Entity, buffId uint32) error {
	if target == 0 {
		log.Error("[buff] 添加 Buff 跳过：目标或 Buff 编号无效 目标=%v Buff编号=%d", target, buffId)
		return fmt.Errorf("buff: invalid target %v for buff %d", target, buffId)
	}
	tab := config.Tab
	desc, ok := tab.BuffConfigConfigByID[int32(buffId)]
	if !ok || desc == nil {
		log.Error("[buff] 添加 Buff 跳过：表中无 Buff 定义 Buff编号=%d", buffId)
		return fmt.Errorf("buff: unknown buff id %d", buffId)
	}
	//  创建buff示例
	com := ecs.EnsureGetComponent[*component.BuffList](w, target)
	newBuf := newBuffInstance(caster, buffId, 1)
	if newBuf == nil {
		log.Error("[buff] 添加 Buff 跳过：创建 Buff 实例失败 Buff编号=%d", buffId)
		return fmt.Errorf("buff: create instance failed for buff %d", buffId)
	}
	//  叠加
	if !overlay.Apply(newBuf, desc, com) {
		log.Debug("[buff] 添加 Buff 跳过：叠层策略拒绝 叠层行为=%d Buff编号=%d 目标=%v", desc.StackBehavior, buffId, target)
		return fmt.Errorf("buff: overlay rejected stackBehavior=%d buffId=%d target=%v", desc.StackBehavior, buffId, target)
	}
	stacks := newBuf.Stacks
	if idx := buff_util.FindDefIndex(com.Buffs, buffId); idx >= 0 {
		stacks = com.Buffs[idx].Stacks
	}
	log.Info("[buff] 添加 Buff 成功 施法者=%v 目标=%v Buff编号=%d 层数=%d", caster, target, buffId, stacks)
	return nil
}

// Remove 从列表中移除指定 Buff 定义；若列表为空则重置 [component.BuffList] 组件。
func Remove(w *ecs.World, e ecs.Entity, buffId uint32) error {
	if w == nil {
		return fmt.Errorf("buff: nil world")
	}
	c, ok := w.GetComponent(e, &component.BuffList{})
	if !ok || c == nil {
		return fmt.Errorf("buff: entity %v has no BuffList", e)
	}
	bl := c.(*component.BuffList)

	idx := buff_util.FindDefIndex(bl.Buffs, buffId)
	if idx < 0 {
		log.Debug("[buff] 移除 Buff：槽位不存在 实体=%v Buff编号=%d", e, buffId)
		return fmt.Errorf("buff: buff %d not on entity %v", buffId, e)
	}
	log.Info("[buff] 移除 Buff 实体=%v Buff编号=%d 移除后剩余实例数=%d", e, buffId, len(bl.Buffs)-1)
	bl.Buffs = slices.Delete(bl.Buffs, idx, idx+1)

	if len(bl.Buffs) == 0 {
		w.RemoveComponent(e, &component.BuffList{})
		w.AddComponent(e, &component.BuffList{})
	}
	return nil
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
