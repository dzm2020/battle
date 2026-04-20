package room

import "testing"

func TestPhaseFSM_Table(t *testing.T) {
	cases := []struct {
		from Phase
		ev   phaseEvent
		want Phase
		ok   bool
	}{
		{PhaseLobby, phaseEvStartBattle, PhasePreBattle, true},
		{PhaseLobby, phaseEvShutdown, PhaseClosed, true},
		{PhasePreBattle, phaseEvBattleLive, PhaseFighting, true},
		{PhasePreBattle, phaseEvShutdown, PhaseClosed, true},
		{PhaseFighting, phaseEvSettle, PhaseSettled, true},
		{PhaseFighting, phaseEvShutdown, PhaseClosed, true},
		{PhaseSettled, phaseEvShutdown, PhaseClosed, true},

		{PhaseLobby, phaseEvSettle, PhaseLobby, false},
		{PhaseLobby, phaseEvBattleLive, PhaseLobby, false},
		{PhaseSettled, phaseEvSettle, PhaseSettled, false},
		{PhaseFighting, phaseEvStartBattle, PhaseFighting, false},
	}
	for _, tc := range cases {
		got, err := transitionPhase(tc.from, tc.ev)
		if tc.ok {
			if err != nil || got != tc.want {
				t.Fatalf("from=%v ev=%v want (%v,nil) got (%v,%v)", tc.from, tc.ev, tc.want, got, err)
			}
		} else {
			if err == nil {
				t.Fatalf("from=%v ev=%v want error got phase %v", tc.from, tc.ev, got)
			}
		}
	}
}

func TestPhaseFSM_Closed(t *testing.T) {
	_, err := transitionPhase(PhaseClosed, phaseEvShutdown)
	if err != ErrRoomClosed {
		t.Fatalf("want ErrRoomClosed got %v", err)
	}
}
