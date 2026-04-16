package skill

// RejectReason 技能被拒绝的原因（服务器日志、网关错误码映射、防作弊审计）。
// 值为 0 表示「未拒绝」。
type RejectReason int

const (
	RejectNone RejectReason = iota
	RejectNotFighting
	RejectSkillSystemOff
	RejectCasterMissing
	RejectUninitialized
	RejectCasterDead
	RejectUnknownSkill
	RejectNotKnownSkill
	RejectStunned
	RejectSilenced
	RejectOnCooldown
	RejectNotEnoughMP
	RejectTargetRequired
	RejectTargetDead
	RejectTargetCamp
	RejectOutOfRange
	RejectWindupBusy // 同一实体已有一条前摇在进行中
)

// String 便于日志输出。
func (r RejectReason) String() string {
	switch r {
	case RejectNone:
		return "ok"
	case RejectNotFighting:
		return "not_fighting"
	case RejectSkillSystemOff:
		return "skill_system_off"
	case RejectCasterMissing:
		return "caster_missing"
	case RejectUninitialized:
		return "entity_not_initialized_for_battle"
	case RejectCasterDead:
		return "caster_dead"
	case RejectUnknownSkill:
		return "unknown_skill_config"
	case RejectNotKnownSkill:
		return "skill_not_learned"
	case RejectStunned:
		return "stunned"
	case RejectSilenced:
		return "silenced"
	case RejectOnCooldown:
		return "on_cooldown"
	case RejectNotEnoughMP:
		return "not_enough_mp"
	case RejectTargetRequired:
		return "target_required"
	case RejectTargetDead:
		return "target_dead"
	case RejectTargetCamp:
		return "target_camp_invalid"
	case RejectOutOfRange:
		return "out_of_range"
	case RejectWindupBusy:
		return "windup_busy"
	default:
		return "unknown_reject"
	}
}
