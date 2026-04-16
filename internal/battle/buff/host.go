package buff

import "battle/internal/battle/attr"

// Host Buff 心跳所需的宿主只读视图（由 *entity.Entity 实现，避免 buff 依赖 entity 包产生循环引用）。
type Host interface {
	AttrBase() attr.Base
	AttrRuntime() *attr.Runtime
	IsDead() bool
}
