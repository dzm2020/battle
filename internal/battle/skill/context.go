package skill

import "battle/internal/battle/entity"

// CastInput 一次施法请求所需的全部「只读快照」。
//
// 重要：Frame / BattleActive 应由房间线程或「战斗单线程」在入队前快照，
// 网络协程若直接调用 TryCast，存在 1～2 帧的竞态窗口；工业实践通常将请求投递到
// 房间 mailbox，在 Loop 内统一 TryCast（第 12 天可加强制串行与频率限制）。
type CastInput struct {
	Frame uint64
	// BattleActive 为 false 时一律拒绝（例如房间已结算）。
	BattleActive bool
	Caster       *entity.Entity
	Target       *entity.Entity
	SkillID      string
}
