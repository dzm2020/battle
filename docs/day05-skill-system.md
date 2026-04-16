# 第 5 天：技能系统（服务端权威全流程）

## 1. 本节目标与边界

本日实现的是 **纯服务端** 技能流水线，与 UE GAS 的关系仅限于「思想借鉴」：**激活 → 校验 →（可选前摇）→ 生效 → 扣费**，不在此引入客户端预测与回滚。

**本日交付：**

- 配表结构 `SkillConfig` 与 `Registry`；
- 可扩展的拒绝原因 `RejectReason`（便于网关映射错误码与风控日志）；
- `ValidateCast` / `ValidateCastAfterWindup` 分层校验；
- `System`：瞬发与前摇（`timer.Manager` + `tick.Subscriber`）；
- `EffectApplier` 钩子：第 6 天伤害、第 7 天 Buff 在此接入；
- `Room` 侧 `SetSkillSystem` / `TryCastSkill` 胶水方法。

**刻意不做的（留给后续天）：**

- 弹道、碰撞、服务器物理步进；
- 完整 Buff 叠加规则（仅预留 `control.Flags`）；
- 战斗邮箱串行化与频率限制（第 12 天）——文档中说明推荐做法。

---

## 2. 模块划分与依赖方向

```
geom        （无战斗依赖）
control     （无战斗依赖）
entity      → attr, calc, cooldown, control, geom
skill       → entity, geom, clock, tick, timer
room        → skill（胶水：TryCastSkill）
```

**原则：**

- `skill` **不引用** `room`，避免双向依赖；房间是否仍在战斗用 `CastInput.BattleActive` 表达。
- `Entity` **不引用** `skill`，只用通用字段（`SkillCD`、`KnownSkills`、`Pos`、`Control`）挂数据，防止 `entity ↔ skill` 循环引用。
- 伤害/Buff **不进** `System` 核心判断，统一走 `EffectApplier`。

---

## 3. 服务端技能流水线（推荐心智模型）

### 3.1 客户端请求

客户端通常只上报：**技能 ID、目标 ID（或坐标）**。服务器 **不信任** 客户端算好的伤害、CD、命中结果。

### 3.2 服务器阶段

| 阶段 | 说明 |
|------|------|
| 路由 | Gateway 将请求路由到对应 `Room`，解析出 `sessionPlayerID` → `Entity`。 |
| 快照 | 读取当前逻辑帧 `Frame`、阶段是否仍 `Fighting`（本仓库为 `TryCastSkill` 内快照）。 |
| 校验 | `ValidateCast`：战斗状态、死亡、学会技能、眩晕、沉默、CD、蓝量、目标、阵营、距离。 |
| 前摇 | `WindupFrames > 0`：登记 `timer.AddOneShot`，并禁止同实体并发前摇（`RejectWindupBusy`）。 |
| 二次校验 | 前摇结束帧 `ValidateCastAfterWindup`：防止「读条期间被控 / 目标已死 / 走出范围」。 |
| 生效 | `EffectApplier.OnSkillApplied`：伤害、治疗、Buff、召唤等（本日默认空实现）。 |
| 扣费 | `applyCosts`：扣蓝、`SkillCD.Trigger`；**成功结算后才进 CD**，前摇起点不扣蓝。 |

### 3.3 与 GAS 的类比（无 UE 代码）

- **Ability 配置** ≈ `SkillConfig`；
- **Cost/Cooldown** ≈ `MPCost` + `CooldownFrames`；
- **Target** ≈ `TargetMode` + 服务器目标解析；
- **Task / 异步段** ≈ 本日的 `WindupFrames` + `timer`；
- **Execution Calculation** ≈ `EffectApplier`（第 6 天填充）。

---

## 4. 配置字段说明（`SkillConfig`）

| 字段 | 含义 |
|------|------|
| `ID` | 与协议、CD Key、客户端展示 ID 对齐。 |
| `School` | `Physical` / `Magic`：**沉默**只挡 `Magic`（可按项目调整）。 |
| `MPCost` | 每次 **成功结算** 消耗的当前魔法。 |
| `CooldownFrames` | 与第 3 天逻辑帧对齐的冷却长度。 |
| `CastRange` | `0` 表示不做距离校验（常用在纯自身技能）。 |
| `WindupFrames` | `0` 瞬发；`>0` 前摇，到期帧尝试结算。 |
| `TargetMode` | `None` / `Self` / `Enemy` / `Ally`：决定目标是否必填及阵营关系。 |

`MustDemoRegistry` 提供三条示例技能：`strike`（瞬发近战）、`fireball`（魔法+前摇+耗蓝）、`focus`（自身、耗蓝）。

---

## 5. 校验顺序与拒绝原因

实现顺序见 `ValidateCast` 源码，建议 **保持「便宜的条件在前」**（如阶段、死亡），昂贵检查（距离）在后。

**`RejectReason` 列表**用于：

- 日志聚合与告警；
- 网关 → 客户端错误码映射；
- 反作弊审计（例如异常高频 `RejectNotKnownSkill`）。

---

## 6. 并发与线程模型（重要）

当前 `System` 使用 **一把互斥锁** 串行化 `TryCast` 与 `OnTick` 内对 `timer` / `pending` 的访问，适合学习与小规模战斗服。

**工业级推荐：**

1. **所有施法请求进入房间邮箱**，由 `tick.Loop` 同线程顺序 `TryCast`，消除网络线程与战斗线程竞态；
2. `Frame` 不再由调用方「快照」，而是 **Loop 内唯一真源**；
3. 对 `EffectApplier` 若接入外部系统（DB、RPC），应 **异步化或短事务**，避免卡住整帧；
4. 第 12 天增加 **施法频率限制**、**位移合法性**，与技能校验层组合。

---

## 7. 与 `Room` 的集成

- **Lobby**：`Room.SetSkillSystem` 绑定本局 `skill.System`（每局独立实例，避免上局 pending 泄漏）。
- **Fighting**：`StartBattle` 内 `ResetForBattle` 后 `loop.Add(skillSys)`，保证帧驱动与 `timer` 一致。
- **Settle / Shutdown**：再次 `ResetForBattle`，清空前摇与定时器，防止下一局误触发。

`TryCastSkill` 为 **胶水 API**：生产环境可替换为「只入队，不在网络线程直接调用 `TryCast`」。

---

## 8. 与后续天的衔接

| 天数 | 接入点 |
|------|--------|
| 第 6 天 | 在 `EffectApplier` 内实现伤害公式，读 `Derived` / `Runtime`。 |
| 第 7 天 | Buff 修改 `Control`、`CastRange` 修正、或动态 `SkillConfig` 修饰（建议 Modifier 管道）。 |
| 第 8 天 | `System` 可改为 ECS System，校验读组件快照。 |
| 第 9 天 | `Pos` 与 AOI 同步；距离校验可换「网格近似距离」。 |
| 第 12 天 | 施法邮箱、频率限制、非法包审计。 |

---

## 9. 运行与测试

```bash
go test ./internal/battle/skill/...
go run .
```

`main` 中 `day05Demo` 使用 **纯 `Loop.Step`** 演示：普攻、沉默挡魔法、火球前摇后扣蓝、自身技能，输出确定、不依赖真实睡眠。

---

## 10. 文件索引

| 路径 | 说明 |
|------|------|
| `internal/battle/skill/config.go` | `SkillConfig`、`Registry`、演示配表 |
| `internal/battle/skill/reject.go` | `RejectReason` |
| `internal/battle/skill/stage.go` | `Stage` |
| `internal/battle/skill/result.go` | `Result` |
| `internal/battle/skill/context.go` | `CastInput` |
| `internal/battle/skill/applier.go` | `EffectApplier` |
| `internal/battle/skill/validate.go` | 校验与目标解析 |
| `internal/battle/skill/system.go` | `System` + `tick.Subscriber` |
| `internal/battle/control/flags.go` | 眩晕/沉默位 |
| `internal/battle/geom/vec2.go` | 距离 |
| `internal/battle/room/room.go` | `SetSkillSystem`、`TryCastSkill` |
