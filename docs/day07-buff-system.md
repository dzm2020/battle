# 第 7 天：Buff / Debuff 系统

## 1. 本节目标与边界

在 **固定逻辑帧** 下，把「持续效果、周期结算、控制与属性修饰」从技能 `System` 中拆出，形成 **可配表、可单测、与实体解耦** 的 Buff 子系统；技能命中只负责 **按 ID 挂 Buff**，具体叠加与心跳由 **Buff.Manager** 与 **tick 循环**驱动。

**本日交付（与本仓库实现对齐）：**

- `buff` 包：`BuffConfig`、`Kind`、`StackPolicy`、`Registry`、`Manager`（`Add` / `Tick` / `Reset`）、`Host` 接口、`CombatMods`；
- 类型覆盖：`KindInstantHeal` / `KindInstantDamage`、持续 **眩晕（控制位）**、**减速（移速乘子）**、**增伤（出伤乘子叠层）**、**攻强（平铺 ATK 叠层）**、**毒（DoT）**；
- 叠加：`StackRefresh`、`StackLayer`、`StackReplace`；`BuffConfig.Priority` 预留（同组互斥覆盖未接业务）；
- `Entity.TickBuffs`：每帧聚合修饰并重算 `Derived`、写 `Control` / `MoveSpeedMul` / `OutgoingDamageMul`；
- `Room.StartBattle`：向 `tick.Loop` 注册 **Buff 心跳订阅者**（在技能 `System` 之前注册，先刷新 Buff 再处理前摇）；
- `SkillConfig.TargetBuffIDs` + `BattleApplier`：命中后对目标 `AddBuff`；
- 单测：`buff` 包内 `mockHost` 测 Manager；`entity` 包内测 `TickBuffs` 与叠层加攻。

**刻意简化或留给后续天的内容：**

- **Buff 分组与互斥、纯 Priority 覆盖链**：仅配表字段占位；
- **驱散 / 免疫 / 不可选中**：未实现标签与过滤；
- **HoT 独立 Kind**：可用 `KindDot` + 正 `TickDeltaHP` 扩展，当前演示以瞬时治疗为主；
- **移速乘子消费端**：已写入 `Entity.MoveSpeedMul`，位移校验（第 12 天）或客户端同步再读取；
- **控制位多来源合并**：当前每帧由 Buff **完全覆盖**写入 `Entity.Control`（非 Buff 来源的控制需后续改为掩码分层）。

---

## 2. 与第 5～6 天的衔接

| 模块 | 关系 |
|------|------|
| 第 5 天 `ValidateCast` | 读取 `Caster.Control`（如 `HasStun()` / `HasSilence()`）；Buff 的 `TickBuffs` 在每帧更新 `Control`，从而影响下一帧能否施法。 |
| 第 6 天 `BattleApplier` | `calc.PhysicalHit` 使用 `caster.OutgoingDamageMul`；增伤 Buff 在每帧 `TickBuffs` 中写入该字段。 |
| `EffectApplier` | 仍承担「命中瞬间」的挂 Buff；**持续过程**不在 `OnSkillApplied` 内手写，统一进 `Manager.Tick`。 |

---

## 3. 模块划分与依赖方向

```
control     （控制位枚举，无战斗业务）
attr        （纯数据结构）
buff        → attr, control；声明 Host 接口，不依赖 entity（避免与 entity→buff 循环引用）
entity      → buff, calc, …
skill       → buff（BattleApplier 里 AddBuff）
room        → tick（Buff 订阅者内调 Entity.TickBuffs）
```

**`buff.Host` 接口**（`internal/battle/buff/host.go`）供 `Manager.Add` / `Tick` 只读访问宿主血量与基础属性；由 `*entity.Entity` 实现，保证 **`buff` 包不 import `entity`**，`buff` 单测可用 `mockHost`。

---

## 4. 配表：`BuffConfig` 与 `Kind`

| `Kind` | 行为概要 |
|--------|----------|
| `KindInstantHeal` / `KindInstantDamage` | `Add` 时直接改 `Runtime.CurHP`，**不进入**持续列表。 |
| `KindStun` | 持续期间通过 `Control` 字段贡献眩晕位（与 `control.FlagStunned` 对齐）。 |
| `KindSlow` | 每层对 `MoveSpeedMul` **乘算** `SlowMoveMul`（多层的乘积）。 |
| `KindDamageAmp` | 每层对 `OutgoingDamageMul` **乘算** `OutDamageMul`（多层的乘积）。 |
| `KindStatATK` | 按层数放大 `StatATKFlat`，在 `CombatMods.BonusATK` 中累加，再由 `Entity` 重算 `Derived.ATK`。 |
| `KindDot` | 按 `TickIntervalFrames` 间隔执行 `TickDeltaHP`（通常为负）；层数参与每跳数值。 |

**时间与周期：**

- `DurationFrames`：从 `Add` 时的 `frame` 起算，`expiresAt = frame + DurationFrames`；`Tick` 中 `frame >= expiresAt` 的实例被移除。
- DoT **首次结算帧**：`nextTickAt = frame + TickIntervalFrames`（添加当帧不跳伤，便于与技能命中帧对齐）。

---

## 5. 叠加策略：`StackPolicy`

| 策略 | 语义（同 `BuffConfig.ID` 再次 `Add`） |
|------|----------------------------------------|
| `StackRefresh` | 若已存在：仅 **延长到期时间** `expiresAt = max(旧, frame + Duration)`；层数不变。 |
| `StackLayer` | 若已存在：**层数 +1**（不超过 `MaxStacks`），并对剩余时间与 **新持续时间** 取较长者刷新到期。 |
| `StackReplace` | **移除**同 ID 旧实例，再 **新建** 一层（常用于强控唯一实例）。 |

瞬时类不走上述持续分支；未知 `BuffID` 时 `Add` 直接忽略。

---

## 6. 帧心跳：`Manager.Tick` 与 `CombatMods`

每逻辑帧对存活实例：

1. 移除已到期实例；
2. 对 `KindDot` 判断是否到达 `nextTickAt`，若是则改 `CurHP`、推进 `nextTickAt`；
3. 第二遍扫描存活实例，聚合 `CombatMods`：**控制位按位或**、减速/增伤 **按层乘算**、攻强 **按层相加**。

`Entity.TickBuffs` 将 `CombatMods` 写回实体字段，并调用 `refreshDerived(cal, mods.BonusATK)`，最后用 `ClampHPToMax` 防止治疗溢出上限。

---

## 7. 与战斗循环、技能的集成

### 7.1 `Room` 与 `tick.Loop`

`StartBattle` 在 `PhaseFighting` 时：

1. 注册 `tick.FuncSubscriber`：在 `c.Frame()` 下遍历本房间玩家实体，调用 `e.TickBuffs(fr, cal)`；
2. 再注册 `skill.System`（若已 `SetSkillSystem`）。

这样同一帧内 **Buff 状态先于技能前摇结算** 更新，避免「本帧已眩晕仍被前摇插队」的时序歧义（工业上仍可再细化为「子阶段」枚举）。

### 7.2 技能挂 Buff

`SkillConfig.TargetBuffIDs`：`BattleApplier` 在伤害/治疗之后，对 **有效目标** 逐个 `AddBuff(ctx.Frame, buffID)`。

演示技能与 Buff 示例见 `skill.MustDemoRegistry()` 与 `buff.DemoRegistry()`（如 `hammer_stun` → `demo_stun`，`focus` → `demo_amp` + `demo_strong`，`fireball` → `demo_poison`）。

---

## 8. 演示配表 ID（`DemoRegistry`）

| Buff ID | 演示点 |
|---------|--------|
| `demo_stun` | 眩晕 + `StackReplace` |
| `demo_slow` | 减速 + `StackRefresh` |
| `demo_amp` | 增伤叠层 + `StackLayer` + `MaxStacks` |
| `demo_poison` | DoT + `StackRefresh` |
| `demo_strong` | 平铺 ATK 叠层 |
| `demo_instant_heal` | 瞬时治疗 |

---

## 9. 单测与运行

```bash
go test ./internal/battle/buff/... ./internal/battle/entity/...
go run .
```

- `internal/battle/buff/manager_test.go`：`mockHost` 测修饰聚合、DoT 节拍、瞬时治疗、到期清除。  
- `internal/battle/entity/buff_test.go`：`TickBuffs` 与 `demo_strong` 叠层加攻。

---

## 10. 与 `record.md` 编码任务对照

| `record.md` 任务 | 本仓库对应 |
|------------------|------------|
| Buff 组件：添加、移除、心跳、周期跳伤、控制 | `buff.Manager`：`Add`、`Tick`、`Reset`；DoT；`KindStun` + `Control` |
| 眩晕、减速、增伤 demo | `demo_stun` / `demo_slow` / `demo_amp` 及关联技能 |
| 叠加、覆盖、刷新、优先级、定时消除 | `StackRefresh` / `StackLayer` / `StackReplace`；到期移除；`Priority` 预留 |

---

## 11. 扩展建议（工业向）

1. **Buff 实例 ID**：除配表 `ID` 外为每次施加生成 `instanceUUID`，便于日志与精确驱散。  
2. **施法者与来源**：实例上记录 `sourceEntityID`、`skillID`，用于统计与击杀归属。  
3. **与邮箱模型结合**：`AddBuff` 仅由战斗线程调用；网络线程只投递事件（第 12 天）。  
4. **第 8 天 ECS**：可将 `Manager` 迁为组件，`Tick` 由 `BuffSystem` 批量调度。
