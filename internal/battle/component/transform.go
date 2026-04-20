package component

// Transform2D 可选空间坐标，供技能按距离排序等逻辑使用；无坐标时距离类排序退化为稳定顺序。
type Transform2D struct {
	X float64
	Y float64
}

func (*Transform2D) Component() {}
