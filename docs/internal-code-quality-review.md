# `internal/` 代码质量评估与优化方向

- 文档与目录约定见 **[STRUCTURE.md](./STRUCTURE.md)**；此类评审长文建议逐步迁入 **`docs/reviews/`** 并在本路径保留简短跳转（可选）。
- **可执行优化清单**见 **[guides/optimization-recommendations.md](./guides/optimization-recommendations.md)**。

> 范围：`internal/battle/` 及与之紧耦合的 ECS、配置与房间逻辑（抽样阅读 + 与当前仓库构建状态交叉验证）。  
> 日期：以仓库当前状态为准。

---

## 1. 总览

项目采用 **ECS（`battle/ecs`）+ 配置表驱动（`internal/battle/config`）+ 战斗系统管线（`internal/battle/system`）+ 房间隔离（`internal/battle/room`）** 的结构，方向清晰：数据与逻辑分层、单局 `World` 隔离、技能/Buff/选目标等子域拆分。整体属于 **可演进的战斗内核雏形**，但 **工程完整度与类型一致性** 尚未收敛，存在阻塞编译的债务与少量架构分叉。

---

## 2. 主要优点

| 维度 | 说明 |
|------|------|
| **域划分** | `component` / `entity_factory` / `skill` / `buff` / `target_selector` / `land` / `room` / `room_builder` 边界大体可读，符合继续长大的形态。 |
| **战斗管线** | `system/register.go` 对系统顺序有明确注释（Buff → 技能 → 伤害/治疗/生命/死亡/结束），利于排查帧内顺序问题。 |
| **房间模型** | `Room` 与 `World`、阶段、`tick` 绑定，注释中强调与网关解耦，职责意识较好。 |
| **空间网格** | `land` 包职责集中，具备单元测试（`grid_test.go`），并抽出 `PickFreeCell` 等可复用 API。 |
| **配置加载** | `config.Tables` + 分表 JSON 的模式简单直接，适合策划迭代。 |

---

## 3. 主要问题与风险（按优先级）

### 3.1 阻塞：无法通过完整构建（P0）

- `internal/battle/utils/attributes.go` 等与 **`Transform2D` / `Team.Side` 类型** 不一致的用法会导致 **`go build ./...` 失败**。  
- 根因典型为：**组件字段类型已演进（如 `Transform2D.X/Y` 为 `int`，`Team.Side` 为字符串枚举），工具函数仍按旧假设（`float64` / `uint8`）编写**。  
- **影响**：任何依赖 `utils` 的包（含 `system`、部分测试）无法稳定 CI，重构难以增量验证。

### 3.2 类型与语义一致性（P0～P1）

- **阵营/坐标**：`CampRelation`、`GetEntityCamp` 与 `component.Team`、`Transform2D` 必须 **单一事实来源**，避免「有的地方用格子索引当坐标、有的地方用 float」的隐式约定。  
- **注释与实现脱节**：例如 `Room` 注释提到 `Room.mu`，结构体中 **无互斥锁字段**，易误导后续并发改造。

### 3.3 架构分叉：`room.Grid` vs `*land.Grid`（P1）

- 仓库中同时存在 **`room/grid.go`（包装 `Grid`）** 与 **`Room.grid *land.Grid` 直连** 两种倾向时，易出现 **「一半走包装 API、一半走 land」** 的双轨维护。  
- **建议**：二选一并写进约定——要么 `Room` 始终持有 `*room.Grid` 并对外只暴露包装；要么删除包装层，统一在 `land` + 小函数上完成「与 ECS 同步」。

### 3.4 配置与运行时校验（P1）

- `config.Load` 对缺失 JSON 直接 **`panic`**：开发期干脆，生产/CI 需要 **可聚合错误**（或 `Load` 返回 `error` + 校验报告）。  
- **交叉引用**（副本 `map_id`、怪物 `unit_id`、技能引用效果 id）缺少 **启动时或 `cmd/configvalidate` 的强校验**，错误易延迟到运行时。

### 3.5 测试覆盖与包级耦合（P2）

- **`internal/battle` 内单测偏少**（除 `land`、`clock`、`tick` 等），大量逻辑依赖 `test/` 包；而 `test` 包易因个别文件损坏导致 **整包无法编译**。  
- **建议**：对 **`entity_factory` 关键路径、`room_builder` 创建流程、`system` 单系统** 增加可独立运行的表驱动测试。

### 3.6 可观测性与调试（P2）

- 战斗循环、技能施法、Buff 叠加等 **缺少结构化日志/指标钩子**（非必须过早做全链路 tracing，但关键帧事件建议预留接口）。

---

## 4. 后续优化方向（路线图）

### 4.1 短期（1～2 周）：恢复「可构建」与类型收敛

1. **修复 `utils` 与 `component` 的类型对齐**（`TransformXY`、`GetEntityCamp` 等），保证 `go test ./internal/...` 与主模块可编译。  
2. **统一阵营与坐标的表示**：文档化「格子索引 vs 世界坐标」；API 命名上区分 `Cell` / `World`。  
3. **`Room` 注释纠偏**：删除或实现 `mu`；若暂无双线程写 `phase`，注明「当前单线程 tick，不加锁」。

### 4.2 中期（2～6 周）：可验证性与装配一致性

1. **`configvalidate` 增强**：表间引用、枚举范围、`Map` 边界与 `cell_size` 合法性。  
2. **`config.Load` 策略**：开发 `panic`、生产 `error` 可配置，或分离 `MustLoad` / `Load`。  
3. **`room` / `room_builder` 单测**：`CreateRoom`、`PVP/PVE` 分支、无格/满格错误路径。  
4. **网格与实体生命周期**：明确 `Remove` 网格占用 vs `RemoveEntity` 的调用顺序（避免残留 `Transform2D`）。

### 4.3 中长期（按需）：性能、并发与扩展

1. **ECS 更新**：大规模单位时 profiling（`GetNearbyUnits`、查询器、Buff 系统）。  
2. **并发模型**：若未来网络线程与战斗线程分离，**刷怪/入房事件需入队**，在 tick 内 `flush`。  
3. **插件化**：`system` 注册由表驱动或构建标签注入，便于 MOD/分支玩法开关。  
4. **错误域**：战斗业务错误与 IO/配置错误分层（`errors.Is` 可判定），避免字符串散落。

---

## 5. 总体评价（一句话）

**架构意图与分包方向良好，具备继续演进的骨架；当前最大短板是类型/工具层与组件层未对齐导致的构建失败，以及少量双轨设计（网格包装、配置加载语义）需要收敛。** 优先恢复可编译与表校验后，再谈性能与并发，投入产出比最高。

---

## 6. 建议的「完成定义」（DoD）清单

- [ ] `go build ./...` 无错误  
- [ ] `go test ./internal/battle/...` 核心子包通过  
- [ ] `configvalidate` 覆盖主要外键引用  
- [ ] `Room` / `land` / `room_builder` 对网格与 ECS 的约定有 **一页架构说明**（可与本文合并）

---

*本文档随仓库演进可修订；若你希望把「阵营/坐标约定」或「房间状态机」单独拆成 ADR，可在 `docs/adr/` 下增量添加。*
