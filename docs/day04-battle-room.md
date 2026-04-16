# 第 4 天：战斗房间（BattleRoom）与管理器

## 目标

- **房间生命周期**：创建 → 大厅加入 → 开战（初始化实体 + 启动 Tick）→ 结算停循环 → 撤房销毁。
- **多房间隔离**：每个 `Room` 独立玩家表、阶段、`Clock`/`Loop`；`Manager` 只做全局索引。
- **并发安全**：`Manager` 用读写锁保护 `map`；`Room` 用互斥锁保护阶段与玩家表；**不在持锁时调用可能阻塞的外部逻辑**（`InitBattle` 在释放锁后执行）。

## 包结构

| 类型 | 路径 | 职责 |
|------|------|------|
| `Phase` | `internal/battle/room/phase.go` | 阶段枚举；含短暂的 `PreBattle` 防止与 `InitBattle` 交错 Join。 |
| `Room` | `internal/battle/room/room.go` | Join/Leave、StartBattle、Settle、Shutdown；挂载 `tick.Loop`。 |
| `Manager` | `internal/battle/room/manager.go` | Create/Get/Remove/Destroy/Count。 |
| 错误 | `internal/battle/room/errors.go` | 可判定错误，供 Gateway 映射协议错误码。 |

依赖方向：`room` → `entity` / `calc` / `clock` / `tick`。**不依赖** `main`、网络、`timer`（后续由 Subscriber 接入）。

## 阶段说明

1. **Lobby**：允许 `Join` / `Leave`（战斗中掉线逻辑后续天扩展）。
2. **PreBattle**：`StartBattle` 已接管流程，**禁止再 Join**，避免与 `InitBattle` 并发写实体。
3. **Fighting**：`Loop.Run` 在独立协程中运行；`ctx` 可由上层或 `Settle`/`Shutdown` 取消。
4. **Settled**：循环已停，可做结算快照、写日志、发奖（本仓库仅占位）。
5. **Closed**：资源已释放；`Shutdown` 幂等。

## Manager 与 Room 的分工

- **Manager**：进程级「房间字典」，适合按 `roomID` 路由网关消息。
- **Room**：单局权威；后续第 5～7 天的技能/Buff **以 Subscriber 形式挂到 `Loop`**，而不是写进 `Manager`。

推荐加锁顺序：先 `Manager.Get`（读锁短），再调用 `Room` 方法（房间自锁）。**不要**在持有 `Room.mu` 时再去锁 `Manager`，当前 API 也未暴露此类需求。

## API 要点

- `StartBattle(ctx, calc.Calculator)`：`ctx` 为父上下文（关服、踢房、超时），内部 `WithCancel`；先 `InitBattle` 再启动 `Run`，避免空血进帧。
- `Settle()`：仅 `Fighting` → `Settled`；`cancel` + `Wait` 确保循环退出。
- `Manager.Destroy(id)`：`Room.Shutdown()` + `Remove`，用于撤房。

## 运行与测试

```bash
go test ./internal/battle/room/...
go run .
```

`main` 中 `day04Demo` 使用短超时 `context` 驱动循环结束，再 `Settle` 与 `Destroy` 演示完整闭环。

## 与后续天数衔接

| 天数 | 衔接 |
|------|------|
| 第 5 天 技能 | 在 `StartBattle` 后 `loop.Add(skillSystem)`；校验读 `entity` + `cooldown`。 |
| 第 7 天 Buff | 同 Subscriber；与房间同生命周期。 |
| 第 10 天 快照 | `Settle` 前导出 `[]*entity` 或组件快照再广播。 |
