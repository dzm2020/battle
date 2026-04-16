package skill

import (
	"sync"

	"battle/internal/battle/clock"
	"battle/internal/battle/entity"
	"battle/internal/battle/tick"
	"battle/internal/battle/timer"
)

// tagWindup 定时器标签：仅用于区分技能前摇，与其他系统 Tag 错开即可。
const tagWindup timer.Tag = 11

// System 技能子系统：同时实现 tick.Subscriber，由房间在 StartBattle 时 loop.Add。
//
// 职责边界：
//   - 维护「配置表 + 校验 + 前摇调度 + 扣费顺序」；
//   - 具体伤害/Buff 在 EffectApplier 中扩展（第 6～7 天）；
//   - 不持有 Room 指针，避免 skill ↔ room 双向依赖；战斗是否仍进行由 CastInput.BattleActive 表达。
//
// 线程模型：TryCast 可与 OnTick 并发；内部使用同一把互斥锁串行化 timer 与 pending。
// 工业部署推荐：将 TryCast 仅暴露在房间邮箱消费者上，网络层只入队。
type System struct {
	mu sync.Mutex

	reg     *Registry
	applier EffectApplier
	tm      *timer.Manager

	pending    map[timer.Handle]*windupJob
	windupBusy map[string]struct{}
}

type windupJob struct {
	Caster  *entity.Entity
	Target  *entity.Entity
	SkillID string
	Config  SkillConfig
}

// NewSystem 创建技能系统；applier 为 nil 时使用 DefaultApplier。
func NewSystem(reg *Registry, applier EffectApplier) *System {
	if reg == nil {
		reg = NewRegistry()
	}
	if applier == nil {
		applier = DefaultApplier{}
	}
	return &System{
		reg:        reg,
		applier:    applier,
		tm:         timer.NewManager(),
		pending:    make(map[timer.Handle]*windupJob),
		windupBusy: make(map[string]struct{}),
	}
}

// Registry 返回配置表指针（只读使用）。
func (s *System) Registry() *Registry { return s.reg }

// ResetForBattle 清空前摇与定时器；应在每局开战前、或结算/撤房后调用，防止状态泄漏到下一局。
func (s *System) ResetForBattle() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tm = timer.NewManager()
	s.pending = make(map[timer.Handle]*windupJob)
	s.windupBusy = make(map[string]struct{})
}

// TryCast 处理一次施法请求（瞬发立即 commit；有前摇则登记 timer）。
func (s *System) TryCast(in CastInput) Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	//  获取技能配置
	cfg, ok := s.reg.Get(in.SkillID)
	if !ok {
		return Result{Reason: RejectUnknownSkill, Stage: StageRejected}
	}
	//  检测是否能释放技能
	if reason := ValidateCast(cfg, in); reason != RejectNone {
		return Result{Reason: reason, Stage: StageRejected}
	}

	if cfg.WindupFrames > 0 {
		//  角色正在前摇不能释放技能
		if _, busy := s.windupBusy[in.Caster.ID]; busy {
			return Result{Reason: RejectWindupBusy, Stage: StageRejected}
		}
		//  加入到前摇定时器
		expire := in.Frame + uint64(cfg.WindupFrames)
		h := s.tm.AddOneShot(expire, tagWindup)
		s.pending[h] = &windupJob{
			Caster:  in.Caster,
			Target:  in.Target,
			SkillID: in.SkillID,
			Config:  cfg,
		}
		s.windupBusy[in.Caster.ID] = struct{}{}
		return Result{OK: true, Stage: StageWindupScheduled, WindupEndsAtFrame: expire}
	}
	//  选择目标
	eff := ResolveEffectiveTarget(cfg, in)
	//  触发战斗效果
	s.commit(cfg, in, eff)
	return Result{OK: true, Stage: StageApplied}
}

func (s *System) commit(cfg SkillConfig, in CastInput, eff *entity.Entity) {
	// 先效果后扣费：便于第 6 天按「当前资源」结算；若策划要求先扣再算，可在此调换顺序。
	s.applier.OnSkillApplied(ApplyContext{
		Frame:  in.Frame,
		Caster: in.Caster,
		Target: eff,
		Config: cfg,
	})
	//  扣除消耗
	s.applyCosts(cfg, in)
}

func (s *System) applyCosts(cfg SkillConfig, in CastInput) {
	in.Caster.Runtime.CurMP -= cfg.MPCost
	in.Caster.SkillCD.Trigger(in.Frame, cfg.ID, cfg.CooldownFrames)
}

// OnTick 驱动前摇到期结算；必须与同一房间的 Loop 绑定，保证帧号与 timer 一致。
func (s *System) OnTick(c *clock.Clock) {
	frame := c.Frame()

	s.mu.Lock()
	evs := s.tm.ProcessFrame(frame)
	var jobs []*windupJob
	for _, ev := range evs {
		if j, ok := s.pending[ev.ID]; ok {
			delete(s.pending, ev.ID)
			delete(s.windupBusy, j.Caster.ID)
			jobs = append(jobs, j)
		}
	}
	s.mu.Unlock()
	//  前摇结束自动释放技能
	for _, j := range jobs {
		in := CastInput{
			Frame:        frame,
			BattleActive: true,
			Caster:       j.Caster,
			Target:       j.Target,
			SkillID:      j.SkillID,
		}
		if reason := ValidateCastAfterWindup(j.Config, in); reason != RejectNone {
			continue
		}
		eff := ResolveEffectiveTarget(j.Config, in)
		s.mu.Lock()
		s.commit(j.Config, in, eff)
		s.mu.Unlock()
	}
}

var _ tick.Subscriber = (*System)(nil)
