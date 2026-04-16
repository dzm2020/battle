# 第 6 天：伤害公式与战斗结算

## 1. 本节目标与边界

在 **服务器权威** 前提下，把「打多少、扣在哪、算不算死」从技能流水线里拆到 **可测试的纯函数 + 单一结算入口**，避免在 `skill.System` 里堆公式。

**本日交付（与本仓库实现对齐）：**

- `calc.PhysicalHit`：基于攻击方 `Derived.ATK`、目标 `PhysMitigation`、以及施法者 **出伤乘子**（为第 7 天 `OutgoingDamageMul` 预留）的 **物理向简化伤害**；
- `skill.BattleApplier`：在 `EffectApplier.OnSkillApplied` 中完成 **扣血**、**治疗（含上限裁剪）**、以及第 7 天已接上的 **挂 Buff**；
- `Entity.IsDead()`：`Runtime.CurHP <= 0` 的死亡判定，供技能校验、Buff 心跳等只读使用。

**刻意简化或留给后续天的内容：**

- **暴击**：`Derived` 中已有 `CritRate` / `CritDamage`，本仓库 `PhysicalHit` 尚未接入随机暴击（可第 13 天整合单测 RNG 注入后再做）；
- **真实伤害 / 元素 / 护盾吸收顺序**：当前结算直接改 `Runtime.CurHP`；`Runtime.Shield` 字段已预留，吸收链未实现；
- **复活**：未实现 `Revive` API；死亡后状态机（倒地无敌、复活读条）可归入后续「结算 / PVE」天；
- **阵营**：第 5 天已在 `ValidateCast` 中校验敌我目标；伤害结算本身不重复判断阵营（由「能否放出技能」保证）。

---

## 2. 与第 5 天的衔接

| 第 5 天 | 第 6 天 |
|--------|--------|
| `EffectApplier` 默认为空 | 使用 `BattleApplier` 写入伤害与治疗 |
| `SkillConfig` 无伤害字段 | 增加 `Damage`、`Heal`（配表驱动结算强度） |
| 校验层禁止死目标 / 错阵营 | 结算层仍应防御性判断 `IsDead()`，避免前摇结束后双杀边界 bug |

**原则不变：** `skill.System` 只做校验、调度、扣费顺序；**数值结果**只出现在 `OnSkillApplied`（或未来独立的 `DamagePipeline`）。

---

## 3. 伤害模型（当前实现）

### 3.1 物理向一击：`calc.PhysicalHit`

输入：

- `attacker`、`target` 的 `attr.Derived`（只读快照，结算瞬间的属性）；
- `outgoingMul`：出伤乘子；`<= 0` 时按 `1` 处理。

计算要点：

1. `raw = ATK(attacker) * outgoingMul`；
2. `mit = clamp(target.PhysMitigation, 0, 0.95)`，避免除零或「完全免伤」导致无意义分支；
3. `dmg = raw * (1 - mit)`，且 **至少为 1**（教学用下限，线上可按策划表改为 0 或最小破防值）。

代码位置：`internal/battle/calc/damage.go`。

### 3.2 技能倍率 / 固定附加值

本仓库采用 **「公式一击 + 配表附加值」** 的折中写法，便于配表只填整数：

- `finalDamage = PhysicalHit(caster.Derived, target.Derived, caster.OutgoingDamageMul) + SkillConfig.Damage`

其中 `SkillConfig.Damage` 表示技能附带的 **额外基数**（可理解为倍率已折进表或策划填的 flat bonus）。若你希望严格「倍率 × ATK」，可把 `Damage` 改为 float 倍率字段，或在 `PhysicalHit` 外包一层 `SkillDamageSpec`。

### 3.3 暴击、真实伤害（规划）

| 类型 | 建议落点 |
|------|----------|
| 暴击 | `calc` 包增加 `PhysicalHitWithCrit(att, def, mul, roll01 float64)` 或注入 `rand.RNG` |
| 真实伤害 | 单独函数 `TrueDamage(amount int64)` 跳过 `PhysMitigation`，或 `DamageType` 枚举分支 |
| 治疗暴击 | 与伤害对称的 `HealCrit` 或统一 `CombatRoll` 结构 |

---

## 4. 战斗结算：扣血、治疗、死亡

### 4.1 扣血

`BattleApplier` 在 `Damage > 0` 且目标存活时：

- 按上一节得到 `dmg`，执行 `target.Runtime.CurHP -= dmg`；
- `CurHP` 裁剪到 `>= 0`，避免出现负血量与客户端展示不一致。

### 4.2 治疗

- `Heal > 0` 时 `CurHP += Heal`，再调用 `buff.ClampHPToMax`（与 Buff 瞬时治疗共用上限逻辑），不超过 `Derived.MaxHP`。

### 4.3 死亡判定

- `Entity.IsDead()`：`Runtime.CurHP <= 0`；
- 技能拒绝、Buff 是否继续 tick 等，应统一以 **服务端实体状态** 为准，不信任客户端上报血量。

### 4.4 复活（未实现时的约定）

建议在 `entity` 或 `room` 层增加显式 API，例如：

- `Revive(e *Entity, hpRatio float64, cal Calculator)`：设置 `CurHP`、清除部分 Debuff、切换阶段；
- 与 `InitBattle` 区分：`InitBattle` 是 **开局拉满**，`Revive` 是 **局中复活**。

当前仓库未提供该方法，避免与 `InitBattle` 语义混淆。

---

## 5. 配置与调用示例

### 5.1 配表字段（`SkillConfig`）

| 字段 | 含义 |
|------|------|
| `Damage` | 在 `PhysicalHit` 结果之上的 **额外伤害基数**；`0` 表示本技能不造成该路径的物理结算 |
| `Heal` | 对 **有效目标** 的治疗量；自疗技能使用 `TargetModeSelf` |

### 5.2 装配 `BattleApplier`

```go
reg := skill.MustDemoRegistry()
sys := skill.NewSystem(reg, skill.BattleApplier{})
```

仅跑技能校验单测、不关心数值时，可继续使用 `skill.DefaultApplier{}`。

---

## 6. 单测与扩展建议

- **单测**：在 `calc` 包对 `PhysicalHit` 做表驱动测试（不同 DEF → 不同 `PhysMitigation` → 不同承伤）；对 `BattleApplier` 可做集成测试（施法后目标 `CurHP` 变化）。
- **日志**：生产环境建议在结算点打结构化日志：`frame`、`skill_id`、`attacker_id`、`target_id`、`damage`、`hp_after`。
- **第 7 天**：`OutgoingDamageMul` 由 Buff 帧心跳写入，`BattleApplier` 已读取；护盾、无敌、伤害类型枚举可继续叠在同一 applier 或拆子模块。

---

## 7. 与 `record.md` 编码任务对照

| `record.md` 任务 | 本仓库对应 |
|------------------|------------|
| 全套伤害计算函数 | `calc.PhysicalHit`（可继续扩展暴击/真实伤） |
| 扣血、死亡 | `BattleApplier` 扣 `CurHP`；`Entity.IsDead()` |
| 复活 | 预留设计，待实现 API |

---

## 8. 运行

```bash
go test ./internal/battle/calc/... ./internal/battle/skill/...
go run .
```

`main` 中若使用 `BattleApplier`，即可在演示流程里观察到技能伤害与后续 Buff 的联动（第 7 天内容）。
