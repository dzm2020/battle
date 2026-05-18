package resource

type TPS struct {
	TPS   int    // 每秒更新多少帧
	Frame uint64 // 当前推进到多少帧
}
