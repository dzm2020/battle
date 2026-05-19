# 战斗内核优化建议（可执行清单）

> 与 [代码质量评估](../internal-code-quality-review.md)、[目录与结构约定](../STRUCTURE.md) 互补：本文侧重 **做什么、按什么顺序做**，便于排期与验收。  
> 范围：`internal/battle/`、`battle/ecs`、相关 `cmd/` / `test/`。

---

## 0. 优先级说明

- **P0**：不完成则无法稳定 CI / 无法交付联调。  
- **P1**：显著降本（维护成本、线上事故率）。  
- **P2**：体验与长期扩展（性能、可观测、工程舒适度）。

---

## 1. P0：工程基线（必须先做）

### 1.1 恢复全量可构建

| 建议 | 说明 |
|------|------|
| 对齐 **`component` 与子包类型假设** | 坐标见 `system/transform.XY`；阵营/属性见 `attrs.GetEntityCamp` 等，须与组件定义一致。 |
| **`go build ./...` 进 CI 门禁** | 合并前必过；可先拆阶段：`ecs` → `land` → `config` → `room` → `system` → 全量。 |
| **类型敏感处减少「隐式转换」** | 距离、排序若需 `float64`，在边界显式转换并单测边界值。 |

**验收**：本地与 CI 执行 `go build ./...` 无错误。

### 1.2 单一事实来源：阵营 / 格子 / 世界坐标

| 建议 | 说明 |
|------|------|
| 文档化 **「格子索引 vs 世界坐标」** | 在 `docs/guides/` 增加一页（如 `spatial-conventions.md`），并在 `Transform2D` 或 `land` 包注释引用。 |
| **`CampRelation` 与 `Team` 同源** | 避免 `uint8` 与 `SideType` 字符串混用；选一种并在 ADR 中记录（见 `docs/adr/`）。 |

**验收**：全仓库 `grep` 阵营/坐标相关无两套互斥约定；相关 utils 单测覆盖红蓝/无阵营。

---

## 2. P1：架构收敛与配置安全

### 2.1 网格与房间：`room.Grid` vs `*land.Grid`

| 建议 | 说明 |
|------|------|
| **二选一** | A）`Room` 只持有 `*room.Grid`，`SetGrid` 只接 `*land.Grid` 并在内部 `NewGrid`；B）删除 `room/grid.go` 包装，全部 `*land.Grid` + 独立 `placement` 小函数。 |
| **写 ADR** | 在 `docs/adr/` 记录最终选择与迁移步骤，避免再次分叉。 |

**验收**：`Room` / `room_bootstrap` / `SpawnSystem` 路径上仅一种网格访问风格。

### 2.2 配置加载与校验

| 建议 | 说明 |
|------|------|
| **`configvalidate` 增强** | 外键：`Dungeon.map_id` → `Map`、`Monster` → `Unit`、技能 → 效果 id 等；输出人类可读报告。 |
| **`Load` 错误策略** | 生产或 CI 使用 `Load() error` 或 `MustLoad` / `Load` 双 API；避免静默 `panic` 难以聚合。 |
| **表版本字段（可选）** | 顶层 `version` 或每表 `schema_version`，客户端与服务器校验不一致时快速失败。 |

**验收**：故意破坏 JSON 交叉引用时，`configvalidate` 非零退出且日志可定位行。

### 2.3 注释与实现一致

| 建议 | 说明 |
|------|------|
| 清理 **过时注释**（如已不存在的 `Room.mu`） | 若未来加锁，再写 ADR + 注释。 |
| **公开 API 一句话说明线程模型** | 例如：「`World` 修改仅在战斗 tick 协程」。 |

---

## 3. P1：测试策略

| 建议 | 说明 |
|------|------|
| **包内 `*_test.go` 优先** | `entity_factory`、`room_bootstrap`、`system` 中单系统表驱动测试。 |
| **减轻 `test/` 整包耦合** | 集成测试独立 package 或标签 `//go:build integration`，避免单文件损坏阻塞全部测试。 |
| **黄金配置最小集** | `test/battle_config` 保持最小可加载表；与 `docs/examples` 职责划分写进 `STRUCTURE.md`。 |

**验收**：`go test ./internal/battle/land/...` 等核心子包在无集成环境下可单独绿。

---

## 4. P2：性能与可扩展（有数据后再做）

| 建议 | 说明 |
|------|------|
| **Profiling 热点** | 单位量上来后：`GetNearbyUnits`、Buff 遍历、目标筛选；用 `pprof` 定数据再优化。 |
| **网格与 AOI** | 若同屏实体多，考虑稀疏结构或分桶；与策划地图尺寸联动。 |
| **ECS 分配** | `NewWorld` 初始容量、`RemoveEntity` 频率与碎片；按需调参并记录基准。 |

**验收**：有基准场景（N 单位、M 帧）与前后对比数据再合入优化 PR。

---

## 5. P2：并发与事件（若上多线程）

| 建议 | 说明 |
|------|------|
| **战斗写 `World` 单线程** | 网络线程只入队；tick 内 `flush` 生成/销毁实体，与现有 `Room` 生命周期一致。 |
| **房间关闭与泄漏** | `Shutdown` / `cancel` / `WaitGroup` 路径单测或竞态检测（`-race` 在 CI 子集）。 |

---

## 6. P2：可观测性

| 建议 | 说明 |
|------|------|
| **结构化日志接口** | 关键：房间创建、开战、结算、配置加载失败；带 `room_id`、`dungeon_id`。 |
| **指标（可选）** | 房间数、帧耗时分布、战斗异常退出次数。 |

---

## 7. 建议排期（示例）

| 周期 | 目标 |
|------|------|
| 第 1 周 | P0 构建修复 + 阵营/坐标约定文档 + 清理明显过时注释 |
| 第 2～3 周 | P1 网格方案 ADR + `configvalidate` 外键 + 包内关键单测 |
| 第 4 周起 | 按业务量启动 P2 profiling 与观测 |

---

## 8. 相关文档

- [目录与结构约定](../STRUCTURE.md)  
- [代码质量评估（背景与风险）](../internal-code-quality-review.md)  
- [ADR 索引](../adr/README.md)  

---

*本文可随迭代更新版本号或日期；重大方向变更请同步新增/修订 ADR。*
