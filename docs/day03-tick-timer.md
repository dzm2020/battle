# 第 3 天：固定帧循环、全局时间与通用计时器

## 目标

- **固定逻辑帧（如 60 TPS）** 驱动战斗：技能 CD、Buff、AI 都落在同一套时序上。
- **逻辑时间与墙钟分离**：战斗结算只使用 `Frame` / `LogicalMs`，便于回放与单测。
- **低耦合**：`tick.Loop` 只负责推进时钟与广播；`timer.Manager`、`cooldown.Book` 不依赖 `entity` / 网络。

## 包结构与依赖

```
internal/battle/clock   → 无战斗依赖
internal/battle/tick    → clock
internal/battle/timer   → 无 clock 依赖（仅 uint64 帧号）
internal/battle/cooldown→ 无其他战斗包依赖
```

依赖方向：**Loop → Clock**；**Timer/Cooldown 被订阅者在 OnTick 里调用**，而不是反过来依赖 Loop，避免环。

## 核心类型

### `clock.Clock`

- `Advance()`：帧号 +1，**仅由** `tick.Loop.Step`（或你自管的单测驱动）调用。
- `Frame()`：当前逻辑帧（第一次 `Step` 后为 1）。
- `LogicalMs()`：`frame * 1000 / TPS`（整除，确定性）。
- `FrameDuration()`：`time.Second / TPS`，供 `Loop.Run` 与真实时间对齐。

### `tick.Loop`

- `Add(Subscriber)`：注册系统（技能、Buff、房间、定时器等）。
- `Step()`：一帧：先 `Advance`，再按注册顺序 `OnTick`（**同步**执行，第 4 天房间内需注意耗时）。
- `Run(ctx)`：`time.Ticker` 按 `FrameDuration` 调 `Step`，直到 `ctx` 取消。

### `timer.Manager`

- **帧驱动**：`AddOneShot(expireFrame, tag)`、`AddRepeat(firstExpire, repeatEvery, tag)`。
- 在 `currentFrame >= expireFrame` 时于 `ProcessFrame(currentFrame)` 中触发。
- 适合：**延时结算**、**Buff 周期 tick**（重复）、**一次性到期**。
- `Tag`：业务自定义枚举，便于日志与调试。

### `cooldown.Book`

- 按 **技能 key** 记录「下次可释放帧」：`Trigger(frame, key, cdFrames)`。
- `IsReady(frame, key)`：与 `timer` 分工——CD 查询密集、按 key；延时任务用 `Handle` + `Tag` 更清晰。

## 与第 2 天 Entity 的关系

- Entity **不**内嵌 Loop/Timer；由 **BattleRoom（第 4 天）** 或 **技能系统（第 5 天）** 持有 `Clock` / `Loop`，对实体做只读或受控写入。
- 扣血、蓝耗仍在后续天在「帧末」或「技能阶段」统一处理，避免在多个 Subscriber 里重复改同一字段。

## 使用约定（推荐）

1. **每房间一个 `Clock` + `Loop`**（或全局战斗线程 + 多房间队列——进阶，本阶段先房间级）。
2. 所有「还有几秒」的展示：用 `nextReadyFrame - currentFrame` 换算，或 `LogicalMs` 差值，**不要**在战斗逻辑里用 `time.Now` 做 CD。
3. `Subscriber.OnTick` 内避免长时间阻塞；重逻辑可投递到 worker（后续优化）。

## 运行与测试

```bash
go test ./...
go run .
```

`main` 中 `day03Demo` 使用 `Step` 做**确定性**打印（不依赖真实睡眠）；`Loop.Run` 可在接入真实战斗线程时使用。

## 与后续天数衔接

| 天数 | 衔接 |
|------|------|
| 第 4 天 BattleRoom | 创建房间时 `NewLoop` + 注册房间 Subscriber，销毁时 `cancel ctx`。 |
| 第 5 天技能 | `cooldown.Book` 做释放校验；读条可用 `timer` 一次性到期。 |
| 第 7 天 Buff | `AddRepeat` 做周期伤害；到期移除用 `AddOneShot`。 |
| 第 10 天同步 | 快照可带 `Frame` 或 `LogicalMs` 作为版本参考。 |
