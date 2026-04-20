package room

// phaseEvent 驱动房间生命周期状态迁移的事件（仅限本包使用；外部通过 Room 的方法间接触发）。
type phaseEvent uint8

const (
	phaseEvStartBattle phaseEvent = iota // Lobby → PreBattle（开战准备，禁止 Join）
	phaseEvBattleLive                    // PreBattle → Fighting（tick 循环已挂载并启动）
	phaseEvSettle                        // Fighting → Settled（停循环）
	phaseEvShutdown                      // 任意运行中 → Closed（撤房）
)

// transitionPhase 查询状态机迁移结果；不涉及副作用。Closed 上对非关闭类事件一律非法。
func transitionPhase(from Phase, ev phaseEvent) (Phase, error) {
	switch from {
	case PhaseLobby:
		switch ev {
		case phaseEvStartBattle:
			return PhasePreBattle, nil
		case phaseEvShutdown:
			return PhaseClosed, nil
		}
	case PhasePreBattle:
		switch ev {
		case phaseEvBattleLive:
			return PhaseFighting, nil
		case phaseEvShutdown:
			return PhaseClosed, nil
		}
	case PhaseFighting:
		switch ev {
		case phaseEvSettle:
			return PhaseSettled, nil
		case phaseEvShutdown:
			return PhaseClosed, nil
		}
	case PhaseSettled:
		switch ev {
		case phaseEvShutdown:
			return PhaseClosed, nil
		}
	case PhaseClosed:
		return PhaseClosed, ErrRoomClosed
	}
	return from, ErrWrongPhase
}

// advancePhaseLocked 在当前已持 [Room.mu] 的前提下应用迁移并更新 [Room.phase]。
func advancePhaseLocked(r *Room, ev phaseEvent) error {
	next, err := transitionPhase(r.phase, ev)
	if err != nil {
		return err
	}
	r.phase = next
	return nil
}
