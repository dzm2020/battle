package skill

// Result 施法结果：网关可将 OK/Reason 映射给客户端；Stage 用于区分「已进前摇」与「已命中」。
type Result struct {
	OK bool
	// Reason 拒绝原因；OK 为 true 时为 RejectNone。
	Reason RejectReason
	Stage Stage
	// WindupEndsAtFrame 当前实现下前摇预计结束的逻辑帧（调度成功时有效）。
	WindupEndsAtFrame uint64
}
