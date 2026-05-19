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
| `system/attrs` | `Attributes` 读写（System / 出生装配使用） |
| `system/skill` / `buff` / `target_selector` | 战斗规则子域 |
| `pb` | 与客户端/存档对齐的结构体 |
| `utils` | 通用小工具；**宜保持薄** |
| `log` / `event` | 日志与事件 |

## 文档入口

- 仓库路径与文档分类：[`docs/STRUCTURE.md`](../../docs/STRUCTURE.md)  
- **优化建议（清单）**：[`docs/guides/optimization-recommendations.md`](../../docs/guides/optimization-recommendations.md)  
- 架构决策：`docs/adr/`  
- 质量评审：[`docs/reviews/battle-ecs-code-quality-review.md`](../../docs/reviews/battle-ecs-code-quality-review.md)

## 依赖原则（简版）

1. `land` 不依赖 `room`、`entity_factory`。  
2. `room` 通过 `room_bootstrap` 触发刷怪入队，不直接 `import entity_factory`。  
3. `room_bootstrap` 不 import 父包 `system`（Installer 由 `system.init` 注册）。  
4. `config.Load` 行为与错误策略以 ADR 或 `docs/guides/config.md` 为准。
