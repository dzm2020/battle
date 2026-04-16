package buff

import (
	"battle/internal/battle/attr"
	"battle/internal/battle/control"
)

// CombatMods 本帧由 Buff 聚合出的战斗修饰（写回 Entity 供伤害与移动读取）。
type CombatMods struct {
	Control             control.Flags
	BonusATK            int64
	MoveSpeedMul        float64
	OutgoingDamageMul   float64
	IncomingDamageTaken int64 // 本帧瞬时结算量（日志/统计）；通常由 Tick 内直接改 HP
}

// DefaultCombatMods 未挂 Buff 时的默认值。
func DefaultCombatMods() CombatMods {
	return CombatMods{
		MoveSpeedMul:      1,
		OutgoingDamageMul: 1,
	}
}

type instance struct {
	cfg        BuffConfig
	stacks     int32
	expiresAt  uint64
	nextTickAt uint64
}

// Manager 单个实体上的 Buff 列表：添加、移除、帧心跳。
type Manager struct {
	reg   *Registry
	insts []instance
}

func NewManager(reg *Registry) *Manager {
	if reg == nil {
		reg = DemoRegistry()
	}
	return &Manager{reg: reg}
}

// Registry 返回配表指针。
func (m *Manager) Registry() *Registry { return m.reg }

// Reset 清空实例（开局 / 撤房）。
func (m *Manager) Reset() {
	m.insts = m.insts[:0]
}

// Add 在 frame 时刻添加一层或刷新 Buff；若配置未知则忽略。
func (m *Manager) Add(frame uint64, buffID string, h Host) {
	if h == nil || h.IsDead() {
		return
	}
	cfg, ok := m.reg.Get(buffID)
	if !ok {
		return
	}
	switch cfg.Kind {
	case KindInstantHeal:
		rt := h.AttrRuntime()
		if rt == nil {
			return
		}
		rt.CurHP += cfg.TickDeltaHP
		return
	case KindInstantDamage:
		rt := h.AttrRuntime()
		if rt == nil {
			return
		}
		rt.CurHP += cfg.TickDeltaHP // 负数伤害
		if rt.CurHP < 0 {
			rt.CurHP = 0
		}
		return
	}

	if cfg.DurationFrames == 0 {
		return
	}
	exp := frame + cfg.DurationFrames

	switch cfg.StackPolicy {
	case StackReplace:
		m.removeByID(cfg.ID)
		m.insts = append(m.insts, m.newInst(cfg, 1, exp, frame))
	case StackRefresh:
		if i, ok := m.findByID(cfg.ID); ok {
			m.insts[i].expiresAt = maxUint64(m.insts[i].expiresAt, exp)
			return
		}
		m.insts = append(m.insts, m.newInst(cfg, 1, exp, frame))
	case StackLayer:
		if i, ok := m.findByID(cfg.ID); ok {
			maxs := cfg.MaxStacks
			if maxs <= 0 {
				maxs = 1
			}
			if m.insts[i].stacks < maxs {
				m.insts[i].stacks++
			}
			rem := m.insts[i].expiresAt - frame
			if cfg.DurationFrames > rem {
				m.insts[i].expiresAt = frame + cfg.DurationFrames
			}
			return
		}
		m.insts = append(m.insts, m.newInst(cfg, 1, exp, frame))
	default:
		m.removeByID(cfg.ID)
		m.insts = append(m.insts, m.newInst(cfg, 1, exp, frame))
	}
}

func (m *Manager) newInst(cfg BuffConfig, stacks int32, expiresAt, frame uint64) instance {
	ins := instance{cfg: cfg, stacks: stacks, expiresAt: expiresAt}
	if cfg.Kind == KindDot && cfg.TickIntervalFrames > 0 {
		ins.nextTickAt = frame + cfg.TickIntervalFrames
	}
	return ins
}

func maxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func (m *Manager) findByID(id string) (int, bool) {
	for i := range m.insts {
		if m.insts[i].cfg.ID == id {
			return i, true
		}
	}
	return 0, false
}

func (m *Manager) removeByID(id string) {
	dst := m.insts[:0]
	for _, x := range m.insts {
		if x.cfg.ID != id {
			dst = append(dst, x)
		}
	}
	m.insts = dst
}

// RemoveAll 移除宿主身上全部 Buff（驱散 / 死亡清理）。
func (m *Manager) RemoveAll() {
	m.insts = m.insts[:0]
}

// Tick 推进一帧：到期移除、DoT 跳伤、聚合修饰器。
func (m *Manager) Tick(frame uint64, h Host) CombatMods {
	mods := DefaultCombatMods()
	if h == nil || h.IsDead() {
		return mods
	}
	rt := h.AttrRuntime()
	if rt == nil {
		return mods
	}

	var alive []instance
	for _, ins := range m.insts {
		if frame >= ins.expiresAt {
			continue
		}
		if ins.cfg.Kind == KindDot && ins.cfg.TickIntervalFrames > 0 && frame >= ins.nextTickAt {
			delta := ins.cfg.TickDeltaHP * int64(ins.stacks)
			rt.CurHP += delta
			if rt.CurHP < 0 {
				rt.CurHP = 0
			}
			if delta < 0 {
				mods.IncomingDamageTaken += -delta
			}
			ins.nextTickAt = frame + ins.cfg.TickIntervalFrames
		}
		alive = append(alive, ins)
	}
	m.insts = alive

	for _, ins := range m.insts {
		st := ins.stacks
		if st < 1 {
			st = 1
		}
		switch ins.cfg.Kind {
		case KindStun:
			mods.Control |= ins.cfg.Control
		case KindSlow:
			for i := int32(0); i < st; i++ {
				mul := ins.cfg.SlowMoveMul
				if mul <= 0 {
					mul = 1
				}
				mods.MoveSpeedMul *= mul
			}
		case KindDamageAmp:
			for i := int32(0); i < st; i++ {
				om := ins.cfg.OutDamageMul
				if om <= 0 {
					om = 1
				}
				mods.OutgoingDamageMul *= om
			}
		case KindStatATK:
			mods.BonusATK += ins.cfg.StatATKFlat * int64(st)
		}
	}
	return mods
}

// ClampHPToMax 在重算 Derived 后调用，防止治疗溢出上限。
func ClampHPToMax(rt *attr.Runtime, d attr.Derived) {
	if rt == nil {
		return
	}
	if rt.CurHP > d.MaxHP {
		rt.CurHP = d.MaxHP
	}
	if rt.CurMP > d.MaxMP {
		rt.CurMP = d.MaxMP
	}
}
