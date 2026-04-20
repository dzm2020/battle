package room

// Phase 房间生命周期阶段（仅服务端权威）；合法迁移见 [phase_fsm.go]。
type Phase int8

const (
	PhaseLobby Phase = iota // 等待加入 / 准备
	PhasePreBattle          // 已开始开战流程，禁止再 Join（防止与 InitBattle 交错）
	PhaseFighting           // 战斗循环运行中
	PhaseSettled            // 已结算，等待销毁或复盘
	PhaseClosed             // 已关闭，不可再操作
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
