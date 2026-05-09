package utils

import "battle/internal/battle/component"

// FindDefIndex 在缓冲表中查找首个 DefID 匹配的槽位；无则返回 -1。
func FindDefIndex(buf []*component.BuffInstance, id uint32) int {
	for i := range buf {
		if buf[i].BuffId == id {
			return i
		}
	}
	return -1
}
