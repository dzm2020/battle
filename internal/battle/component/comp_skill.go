package component

import "battle/ecs"

// SkillCastRequest 施法请求（玩法层写入的唯一入口）。
// 由 [system.CastValidationSystem] 在本帧校验；通过后移除并写入 [SkillCastState]。
// 写入请使用 [battle/internal/battle/system/skill.SetSkillCastRequest] 或 [skill.RequestSkillCast]。
type SkillCastRequest struct {
	SkillID      int32   // 技能配置 ID（与 [RuntimeSkill].ConfigID / 技能表一致）
	TargetEntity ecs.Entity // 主目标；无目标技能可为 0
	CastPosition Vector2 // 地面/指向技能落点；非位置技能可留零值
	Frame        int     // 请求帧号（可选，用于优先级或回放）
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

// SkillSet 玩家拥有技能集合。
type SkillSet struct {
	Skills []*RuntimeSkill
}

func (*SkillSet) Component() {}

// SkillCastState 技能释放状态（校验通过后由 CastValidationSystem 写入）。
type SkillCastState struct {
	IsCasting       bool
	SkillId         int
	RemainingFrames int
	TargetEntity    ecs.Entity
	CastPosition    *Vector2
	Phase           SkillStage
}

func (*SkillCastState) Component() {}
