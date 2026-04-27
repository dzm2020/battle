package room

// Phase 房间当前阶段（仅服务端权威）；由 [Room] 在各 API 内直接维护 [Room.phase] 字段。
type Phase int8

const (
	PhaseLobby     Phase = iota // 等待加入 / 准备
	PhasePreBattle              // 已开始开战流程，禁止再 Join（防止与 InitBattle 交错）
	PhaseFighting               // 战斗循环运行中
	PhaseSettled                // 已结算，等待销毁或复盘
	PhaseClosed                 // 已关闭，不可再操作
)

func (p Phase) String() string {
	switch p {
	case PhaseLobby:
		return "Lobby"
	case PhasePreBattle:
		return "PreBattle"
	case PhaseFighting:
		return "Fighting"
	case PhaseSettled:
		return "Settled"
	case PhaseClosed:
		return "Closed"
	default:
		return "Unknown"
	}
}
