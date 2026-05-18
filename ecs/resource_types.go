package ecs

// ResID 资源类型标识（ID 式 API 使用）。
type ResID struct {
	id uint8
}

// Index 返回资源在内部表中的索引。
func (id ResID) Index() uint8 {
	return id.id
}
