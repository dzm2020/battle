package cooldown

// Book 技能 CD 登记簿：以「第几帧起可再次释放」为权威，不存 wall time。
// 与 timer.Manager 分工：CD 适合按 key 查询可否施法；延时任务适合 Handle/Tag。
type Book struct {
	nextReady map[string]uint64
}

func NewBook() *Book {
	return &Book{nextReady: make(map[string]uint64)}
}

// IsReady 当前帧是否可释放（未登记或当前帧 >= nextReady）。
func (b *Book) IsReady(frame uint64, skillKey string) bool {
	t, ok := b.nextReady[skillKey]
	return !ok || frame >= t
}

// NextReadyFrame 返回技能下次可用帧；从未使用过则 ok=false。
func (b *Book) NextReadyFrame(skillKey string) (frame uint64, ok bool) {
	t, ok := b.nextReady[skillKey]
	return t, ok
}

// Trigger 从当前帧开始冷却 cdFrames 帧（简单模型：下次可用 = frame + cdFrames）。
func (b *Book) Trigger(frame uint64, skillKey string, cdFrames uint64) {
	b.nextReady[skillKey] = frame + cdFrames
}

// Reset 清除某技能 CD（例如死亡清空、GM）。
func (b *Book) Reset(skillKey string) {
	delete(b.nextReady, skillKey)
}
