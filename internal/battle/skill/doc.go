// Package skill 实现战斗服「技能」子系统：配表、分层校验、前摇调度与效果钩子。
//
// 设计要点（第 5 天）：
//   - 与 room 解耦：不反向依赖房间类型，由调用方提供 CastInput.BattleActive 与逻辑帧；
//   - 与 entity 解耦：实体只暴露 SkillCD / KnownSkills / Pos / Control 等挂点；
//   - 与伤害解耦：具体结算走 EffectApplier，便于第 6～7 天替换实现。
//
// 使用顺序建议：
//   1. Registry 注册 SkillConfig（或从配表加载）；
//   2. NewSystem + 将 System 作为 tick.Subscriber 挂到房间 Loop；
//   3. 开战前 System.ResetForBattle，开战时实体 InitBattle + GrantSkill；
//   4. 每帧 Loop 驱动 OnTick；施法请求 TryCast（推荐最终改为邮箱串行）。
package skill
