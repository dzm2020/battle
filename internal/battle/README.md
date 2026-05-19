# `internal/battle` — 战斗内核

对外模块名：`battle/internal/battle/...`（应用私有，不作为稳定对外 SDK 承诺）。

## 子包一览

| 包路径 | 职责 |
|--------|------|
| `config` | JSON 表加载、`Tables`、按 id 查询 |
| `component` | ECS 组件定义与 `Register` |
| `resource` | 局内 World 资源：`TPS`、`SpawnRequestQueue`、`RoomSpec`、`RoomPhase` 等 |
| `land` | 空间网格、邻近查询；不依赖 `room` |
| `room` | 单局房间：`World`、`runLoop`、`room.Create` |
| `system` | 战斗帧内 Systems、`AddCoreCombatSystems`、注册顺序 |
| `system/room_bootstrap` | 按副本类型：`Installer` 挂 System + `Spawner` 入队刷怪 |
| `system/entity_factory` | 从单位表 / PB 创建实体与初始技能、Buff |
| `system/attrs` | 基础 `Attributes` 读写；`Final` / `Recompute`；World 级阵营查询 |
| `system`（`AttributeSystem`） | 每帧：最终属性 = 基础 + Buff 修正 → `FinalAttributes` |
| `system`（`ResourceSystem`） | 战斗资源：施法消耗队列 + 法力/怒气/能量自然恢复 |
| `system/skill` / `buff` / `target_selector` | 战斗规则子域 |
| `system`（`PVERulesSystem` / `PVPRulesSystem`） | 副本类型专用扩展（占位，按 PVE/PVP 挂载） |
| `system/transform` | `Transform2D` 坐标读取（`XY`） |
| `system/action` | 行动资格判定（`CanAct` 等） |
| `system/distance` | 实体间平面距离平方（`SquaredFromRef`，用于排序） |
| `system/combatmath` | 数值比较与换算常量（`CompareFloat64`、`Thousand` 等） |
| `pb` | 与客户端/存档对齐的结构体 |
| `log` / `event` | 日志与事件 |

## 文档入口

- 仓库路径与文档分类：[`docs/STRUCTURE.md`](../../docs/STRUCTURE.md)  
- **优化建议（清单）**：[`docs/guides/optimization-recommendations.md`](../../docs/guides/optimization-recommendations.md)  
- 架构决策：`docs/adr/`  
- **质量评审（ECS + 代码质量，最新）**：[`docs/reviews/battle-ecs-code-quality-review.md`](../../docs/reviews/battle-ecs-code-quality-review.md) — 综合 **8.0/10**，广义 ECS **符合**

## 依赖原则（简版）

1. `land` 不依赖 `room`、`entity_factory`。  
2. `room` 通过 `room_bootstrap` 触发刷怪入队，不直接 `import entity_factory`。  
3. `room_bootstrap` 不 import 父包 `system`（Installer 由 `system.init` 注册）。  
4. 配表使用 `config.Load`（返回 `error`）或启动路径 `config.MustLoad`；CI 可用 `go run ./cmd/configvalidate <dir>`。
