package resource

// BattleState 单局战斗进度；由开战/结束相关 System 读写。
type BattleState struct {
	Started      bool
	OpeningSides int
}
