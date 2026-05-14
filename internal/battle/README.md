# `internal/battle` — 战斗内核

对外模块名：`battle/internal/battle/...`（应用私有，不作为稳定对外 SDK 承诺）。

## 子包一览

| 包路径 | 职责 |
|--------|------|
| `config` | JSON 表加载、`Tables`、按 id 查询 |
| `component` | ECS 组件定义与 `Register` |
| `entity_factory` | 从单位表 / PB 创建实体与初始技能、Buff |
| `land` | 空间网格、邻近查询；不依赖 `room` |
| `room` | 单局房间：`World`、阶段、`tick`、与网格相关的房间级 API |
| `room_builder` | 进房装配：副本类型、刷怪、玩家落位、`CreateRoom` 等 |
| `system` | 战斗帧内 Systems 及注册顺序 |
| `skill` / `buff` / `target_selector` | 战斗规则子域 |
| `pb` | 与客户端/存档对齐的结构体 |
| `utils` | 通用小工具；**宜保持薄**，避免与 `component` 类型不一致 |
| `log` / `event` | 日志与事件 |

## 文档入口

- 仓库路径与文档分类：[`docs/STRUCTURE.md`](../../docs/STRUCTURE.md)  
- **优化建议（清单）**：[`docs/guides/optimization-recommendations.md`](../../docs/guides/optimization-recommendations.md)  
- 架构决策：`docs/adr/`  
- 质量评审类长文：建议放在 `docs/reviews/`（见 STRUCTURE）

## 依赖原则（简版）

1. `land` 不依赖 `room`、`entity_factory`。  
2. `room` 尽量不依赖 `entity_factory`（造单位在 `room_builder`）。  
3. `config.Load` 行为与错误策略以 ADR 或 `docs/guides/config.md` 为准。
