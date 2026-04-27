package component

import "battle/ecs"

// SkillCastRequest
// @Description: 请求释放技能
type SkillCastRequest struct {
	SkillID      int32      // 技能ID
	TargetEntity ecs.Entity // 主目标实体（若无目标则存Invalid）
	CastPosition Vector2    // 释放位置（用于地面目标技能）
	Frame        int        // 请求发出的帧号（用于优先级/取消）
}

func (*SkillCastRequest) Component() {}

type SkillStage = int

const (
	SkillStageNone      SkillStage = iota
	SkillStagePreCast              // 前摇
	SkillStagePostCast             // 释放
	SkillStageAfterCast            // 后摇
)

type RuntimeSkill struct {
	ConfigID        int // 技能配置ID（对应SkillBaseConfig）
	CurrentCooldown int // 当前剩余冷却帧数（0表示可用）
}

// SkillSet
// @Description: 玩家拥有技能集合
type SkillSet struct {
	Skills []*RuntimeSkill
}

func (*SkillSet) Component() {}

// SkillCastState
// @Description: 技能释放状态
type SkillCastState struct {
	IsCasting       bool       // 是否正在释放技能
	SkillId         int        // 正在释放的技能ID
	RemainingFrames int        // 当前阶段剩余帧数
	TargetEntity    ecs.Entity // 记录释放时的主目标
	CastPosition    *Vector2
	Phase           SkillStage // 技能阶段
}

func (*SkillCastState) Component() {}

// CastIntent 由外部玩法层写入，表示“本实体希望施放某技能”；[system.SkillIntentSystem] 消费后应移除组件。
// 同一实体同一帧至多处理一次意图（后者覆盖前者由玩法层避免）。
type CastIntent struct {
	// SkillID 对应 [skill.SkillConfig].ID。
	SkillID uint32
	// Target 主目标；[skill.TargetScopeSelf] 等作用范围下可忽略；单体与链式需有效实体 ID。
	Target ecs.Entity
}

func (*CastIntent) Component() {}

// SkillUser 旧版技能会话数据（冷却系统与部分测试仍引用）；新施法管线以 [SkillSet] / [SkillCastRequest] 为准。
type SkillUser struct {
	Mana                int
	GrantedSkillIDs     []uint32
	CooldownRemaining   map[uint32]int
}

func (*SkillUser) Component() {}
