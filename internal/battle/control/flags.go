package control

// Flags 控制类状态位掩码（眩晕、沉默等）。
// 第 7 天 Buff 系统会写入这些位；技能校验只读，不在此包内改 Entity。
type Flags uint8

const (
	// FlagStunned 眩晕：通常禁止移动与施法（本工程第 5 天：禁止一切技能）。
	FlagStunned Flags = 1 << iota
	// FlagSilenced 沉默：禁止「魔法学派」技能；物理技能仍可按配置放行。
	FlagSilenced
	// FlagRooted 定身：第 5 天仅占位，位移校验（第 12 天）会用到。
	FlagRooted
)

func (f Flags) HasStun() bool    { return f&FlagStunned != 0 }
func (f Flags) HasSilence() bool { return f&FlagSilenced != 0 }
func (f Flags) HasRoot() bool    { return f&FlagRooted != 0 }
