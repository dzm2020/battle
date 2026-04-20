package ecs

import "sync/atomic"

// ========== 实体管理 ==========

// Entity 实体标识符
type Entity uint64

// entityCounter 全局实体计数器
var entityCounter uint64 = 0

// NewEntity 创建新实体
func NewEntity() Entity {
	return Entity(atomic.AddUint64(&entityCounter, 1))
}
