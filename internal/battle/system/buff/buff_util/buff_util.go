package buff_util

import (
	"battle/ecs"
	"battle/internal/battle/component"
)

// FindDefIndex 在缓冲表中查找首个 DefID 匹配的槽位；无则返回 -1。
func FindDefIndex(buf []*component.BuffInstance, id uint32) int {
	for i := range buf {
		if buf[i].BuffId == id {
			return i
		}
	}
	return -1
}

func HasBuff(world *ecs.World, e ecs.Entity, buffId uint32) bool {
	c, _ := world.GetComponent(e, &component.BuffList{})
	if c == nil {
		return false
	}
	bl := c.(*component.BuffList)
	for _, one := range bl.Buffs {
		if one.BuffId == buffId {
			return true
		}
	}
	return false
}
