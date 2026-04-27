package component

// Health 生命；单位一般与 [Attributes] 同加于战斗实体。
type Health struct {
	Current int
	Max     int
}

func (*Health) Component() {}
