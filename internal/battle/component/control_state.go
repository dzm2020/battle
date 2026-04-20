package component

import "battle/internal/battle/control"

// ControlState 本帧由 [system.BuffSystem] 从 [BuffList] 中所有 EffectControl 项按位或汇聚而成；
// 无相关 Buff 时仍可能保留零值组件一帧，在 Buff 列表清空时与 [StatModifiers] 一并移除。
// 行为门控见 [action.CanAct] 与 [control.Flags]（如眩晕、沉默、定身占位）。
type ControlState struct {
	Flags control.Flags
}

func (*ControlState) Component() {}
