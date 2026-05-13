package component

// Health 生命镜像组件（与 [Attributes] 中 hp 同步；选目标、战斗结束等逻辑可读此组件）。
type Health struct {
	Current int
	Max     int
}

func (*Health) Component() {}
