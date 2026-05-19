package resource

import "battle/ecs"

// Phase 房间当前阶段（服务端权威）；存于 World Resource [RoomPhase]。
type Phase int8

const (
	PhaseLobby     Phase = iota // 等待加入 / 准备
	PhasePreBattle              // 已开始开战流程，禁止再 Join（防止与 InitBattle 交错）
	PhaseFighting               // 战斗循环运行中
	PhaseSettled                // 已结算，等待销毁或复盘
	PhaseClosed                 // 已关闭，不可再操作
)

type RoomPhase struct {
	Phase Phase
}

// SetPhase 更新本局 [RoomPhase]；无 Resource 时无操作。
func SetPhase(w *ecs.World, phase Phase) {
	if w == nil {
		return
	}
	p := ecs.GetResource[RoomPhase](w)
	if p != nil {
		p.Phase = phase
	}
}
