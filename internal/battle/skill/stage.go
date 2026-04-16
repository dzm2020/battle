package skill

// Stage 描述一次施法请求在服务端所处的流水线阶段（便于日志与联调）。
type Stage uint8

const (
	// StageRejected 在校验阶段失败，未产生任何战斗效果。
	StageRejected Stage = iota
	// StageWindupScheduled 已通过校验并登记前摇，等待第 3 天 timer 在到期帧结算。
	StageWindupScheduled
	// StageApplied 已成功结算（瞬发或前摇完成后的那一次 Apply）。
	StageApplied
)
