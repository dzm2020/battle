package attr

// Runtime 临时战斗状态：局内可变，后续可由 Buff、技能、Redis 临时层改写。
// 与 Derived 分离，避免把「当前血量」和「属性上限」绑死在同一结构里。
type Runtime struct {
	CurHP  int64
	CurMP  int64
	Shield int64
}
