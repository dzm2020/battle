package skill

import (
	"battle/internal/battle/entity"
	"battle/internal/battle/geom"
)

// ResolveEffectiveTarget 根据 TargetMode 决定「距离与阵营」校验所用的目标实体。
// TargetModeSelf 时始终返回施法者，忽略 CastInput.Target。
func ResolveEffectiveTarget(cfg SkillConfig, in CastInput) *entity.Entity {
	switch cfg.TargetMode {
	case TargetModeSelf:
		return in.Caster
	default:
		return in.Target
	}
}

// ValidateCast 施法前校验（瞬发与前摇起点共用）。
// 不负责「是否已有前摇」——该规则由 System 额外串行保证，避免配置层感知调度细节。
func ValidateCast(cfg SkillConfig, in CastInput) RejectReason {
	if !in.BattleActive {
		return RejectNotFighting
	}
	if in.Caster == nil {
		return RejectCasterMissing
	}
	if in.Caster.IsDead() {
		return RejectCasterDead
	}
	if in.Caster.SkillCD == nil {
		return RejectUninitialized
	}
	if !in.Caster.KnowsSkill(in.SkillID) {
		return RejectNotKnownSkill
	}

	if in.Caster.Control.HasStun() {
		return RejectStunned
	}
	if cfg.School == SchoolMagic && in.Caster.Control.HasSilence() {
		return RejectSilenced
	}

	if !in.Caster.SkillCD.IsReady(in.Frame, cfg.ID) {
		return RejectOnCooldown
	}
	if in.Caster.Runtime.CurMP < cfg.MPCost {
		return RejectNotEnoughMP
	}

	switch cfg.TargetMode {
	case TargetModeEnemy, TargetModeAlly:
		if in.Target == nil {
			return RejectTargetRequired
		}
		if in.Target.IsDead() {
			return RejectTargetDead
		}
	case TargetModeSelf, TargetModeNone:
		// 无强制目标
	}

	if cfg.TargetMode == TargetModeEnemy {
		if in.Target == nil {
			return RejectTargetRequired
		}
		if in.Target.Camp == in.Caster.Camp {
			return RejectTargetCamp
		}
	}
	if cfg.TargetMode == TargetModeAlly {
		if in.Target == nil {
			return RejectTargetRequired
		}
		if in.Target.Camp != in.Caster.Camp {
			return RejectTargetCamp
		}
	}

	if cfg.CastRange > 0 {
		eff := ResolveEffectiveTarget(cfg, in)
		if eff != nil {
			if geom.Dist(in.Caster.Pos, eff.Pos) > cfg.CastRange+1e-6 {
				return RejectOutOfRange
			}
		}
	}

	return RejectNone
}

// ValidateCastAfterWindup 前摇结束时的二次校验：控制、距离、蓝量、目标存活等可能已变化。
func ValidateCastAfterWindup(cfg SkillConfig, in CastInput) RejectReason {
	// 与首帧校验一致；CD 仍未进入（成功结算后才 Trigger），因此仍应 IsReady。
	return ValidateCast(cfg, in)
}
