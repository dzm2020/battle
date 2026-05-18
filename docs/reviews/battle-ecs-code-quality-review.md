# battleDemo 代码质量与 ECS 规范评价（复审）

| 项 | 说明 |
|----|------|
| **范围** | `battle/ecs`、`internal/battle/` 及配置、房间、测试 |
| **评价日期** | 2026-05-18（复审，相对首轮评审已有多项架构收敛） |
| **相关文档** | [STRUCTURE.md](../STRUCTURE.md)、[skill-design.md](../skill-design.md)、[buff-and-skill-systems.md](../buff-and-skill-systems.md) |

---

## 1. 总评（一句话）

**战斗内核方向正确，ECS 边界收敛已可验证**：纯数据组件、`BattleContext`、刷怪队列、`test/` 与 `room_factory` 已对齐。当前 **`go build ./internal/battle/...` 通过**，**`go test ./ecs/... ./test/... ./internal/battle/...` 通过**。整体仍是 **Demo 级工程**：`abandoned/` 拖累全仓库 `go build ./...`，施法目标仍以 `TargetSelect` 表为准（`SkillCastState.TargetEntity` 未参与效果选取），Query 仍为全表扫描。

**ECS 定位**：符合 **广义 ECS**（E/C/S + 帧驱动 System + 命令缓冲）；不符合 **狭义高性能 ECS**（无 Archetype、全表扫描 Query）。

---

## 2. 自首轮评审以来的改进（已落地）

| 项 | 状态 |
|----|------|
| 组件带行为（`Attributes.Add`、`DamageQueue.Add`） | 已改为 `attr_ops` / `DamageQueueAppend` 包级函数，组件仅数据 |
| HP 双轨（`Attributes` + `Health`） | 已移除 `Health` 组件，`AttrHp` 为唯一生命源 |
| 施法双轨（`CastIntent` + `SkillCastRequest`） | 仅保留 `SkillCastRequest` + `RequestSkillCast` |
| 出生装配散落在 `unit` | 已迁至 `factory/entity_factory`，职责注释清晰 |
| Resource 分散注入 | 已统一为 `system/runtime.BattleContext`（`Grid` + `SpawnQueue`） |
| 进房刷单位 | `room_factory` 入队 → `SpawnSystem` 消费，符合 command-buffer 模式 |

---

## 3. 当前架构快照

```text
battle/ecs/                    # 通用 World、Query、Resource、事件
internal/battle/
├── component/                 # 纯组件 + attr_ops / skill_cast_ops / damage 访问函数
├── config/                    # JSON 表
├── factory/
│   ├── entity_factory/      # 出生装配（唯一允许初始 skill/buff 挂载处）
│   └── room_factory/          # 副本类型、入队 SpawnRequest（原 room_builder）
├── land/                      # 空间网格
├── room/                      # 房间生命周期、CreateRoom、tick
├── system/
│   ├── runtime/               # BattleContext 注入/访问
│   ├── skill/ buff/ target_selector/  # 战斗子域（部分已迁入 system 树）
│   └── system_*.go            # 帧管线
├── pb/ utils/ event/ tick/
```

**开房与战斗帧（推荐理解）**

```text
CreateRoom → component.Init → runtime.Install(BattleContext)
          → room_factory.Create（AddCombatSystems + EnqueueSpawn）
          → StartBattle（tick → World.Update）
每帧：SpawnSystem → Buff → 技能 CD/校验/阶段 → 伤/疗/血/死/结束
```

**房间创建不是 System**：开房属于局生命周期；单位生成在帧内由 `SpawnSystem` 完成——符合 ECS 常见分工。

---

## 4. 是否符合 ECS 规范？

### 4.1 符合或基本符合

| 要点 | 说明 |
|------|------|
| Entity 仅为 ID | `ecs.Entity`，无业务方法实体类 |
| System 驱动 | 伤害/治疗/Buff/技能/刷怪等均在 `Update` + Query 中处理 |
| 组合式能力 | `SkillSet`、`BuffList`、`Team`、`Attributes` 等拼装单位 |
| 命令缓冲 | `DamageQueue`→`ResolvedDamage`、`PendingHeal`、`SpawnRequestQueue` |
| 注册与顺序 | `component.InitCombatTypes` + `register.go` 显式系统顺序 |
| 单 World 隔离 | 每局 `Room` 独立 `ecs.World` |
| 玩法写、系统读 | 施法：`RequestSkillCast`；刷怪：`runtime.EnqueueSpawn` |

### 4.2 仍不符合或需注意

| 问题 | 严重度 | 说明 |
|------|--------|------|
| 存储模型 | 中 | 每实体 `map[uint8]Component`，Query 全表 O(N)，无 Archetype/SoA |
| `mask uint64` | 低～中 | `compID` 可达 255，**>64 种组件时位掩码失效** |
| `AddComponent` 静默 | 中 | 已有同类型组件时直接 return，难调试 |
| 属性写在包级函数 | 低 | 严格 ECS 希望 System 内改值；当前 `attr_ops` 是务实折中 |
| `runtime` 包路径 | 低 | 实现在 `system/runtime`，import 易写错（复审前曾导致无法编译） |
| 目录迁移中 | 低～中 | `test/` 已用 `room` + `room_factory`；`abandoned/room_builder` 仍无法编译；README 未完全同步 |
| 施法目标 | 中 | 校验写入 `SkillCastState.TargetEntity`，`ApplyEffects` 仅用 `target_selector.Select`，与玩家点选目标可能不一致 |

### 4.3 结论

| 标准 | 结论 |
|------|------|
| **广义 ECS** | **符合**，可作为帧同步战斗内核继续演进 |
| **狭义 ECS（DOTS/EnTT 类）** | **不符合** |

---

## 5. 代码质量（工程维度）

### 5.1 优点

- **管线清晰**：系统顺序有注释，伤害链、治疗、死亡、战斗结束分工明确。
- **表驱动扩展**：skill/buff 效果 registry、target_selector 过滤器可配置。
- **边界文档化**：`entity_factory`、`BattleContext`、`SkillCastRequest` 包注释说明职责。
- **`land` 有单测**；`component`、`config` 等子包可独立编译。

### 5.2 主要问题

| 优先级 | 问题 |
|--------|------|
| **P1** | `go build ./...` 因 `abandoned/room_builder` 失败；应删除、移出模块或加 `//go:build ignore` |
| **P1** | 技能效果未消费 `SkillCastState.TargetEntity`；表 `target_select_id` 与点选目标易漂移 |
| **P2** | `CreateRoom` → `room_factory.Create` → `StartBattle`；测试需 `World.Update` 才刷怪，文档需写清 |
| **P2** | `config.Load` panic、表间引用校验不足（见既有优化清单） |
| **P2** | `internal/battle` 内 System 级单测仍偏少（集成测在 `test/`） |

### 5.3 构建与测试验证（2026-05-18 本轮）

```text
go build ./internal/battle/...              → 通过
go test ./ecs/... ./test/... ./internal/battle/... → 通过
go build ./...                              → 失败（仅 abandoned/room_builder）
```

**本轮测试/战斗修复（摘要）**

| 项 | 处理 |
|----|------|
| `ecs` 单测 | 对齐 `NewWorld`、事件 Kind/Payload |
| `test/` | `room` + `room_factory`；`Skill.json` 消耗类型字符串化；`SkillEffect` 敌方选目标 `target_select_id: 10` |
| `DamageSystem` | 命中后基础伤害取 `RawDamage`；未配命中/闪避视为必中 |

---

## 6. 主观评分（复审）

| 项 | 首轮 | 复审 | 本轮 | 说明 |
|----|------|------|------|------|
| 架构清晰度 | 7 | 8 | **8** | factory + runtime + 管线稳定 |
| ECS 纯度 | 5 | 6.5 | **6.5** | 组件/施法/生命源已收敛；目标选取仍偏表驱动 |
| 可维护性 | 5 | 6 | **7** | ecs + test 可绿；`abandoned` 仍碍全仓构建 |
| 可扩展性（玩法） | 7 | 7.5 | **7.5** | 表驱动 + 队列模式 |
| 生产就绪 | 3 | 4 | **5** | battle 可编可测；CI 建议 `ecs` + `test` + 排除 abandoned |

---

## 7. 两个问题的直接回答

### 7.1 这份库代码怎么样？

**作为战斗 Demo / 内核原型：质量明显提升，值得继续投入。** 近期重构解决了评审中的多个「ECS 坏味道」。  
**作为上线战斗服：** 建议 CI 跑 `go test ./ecs/... ./test/... ./internal/battle/...` 与 `go build ./internal/battle/...`；清理 `abandoned`、让效果尊重点选目标或文档约定「仅表选目标」、并评估 Query 性能上限。

### 7.2 是否符合 ECS 规范？

- **符合广义 ECS**，且比首轮更接近规范实践（单一数据源、单一施法入口、统一 Context、刷怪走 System）。  
- **不符合狭义高性能 ECS**；若单位规模上千，需做 Query 索引或 Archetype，而非继续堆组件种类。

---

## 8. 建议的下一步（按优先级）

1. **P1**：删除/隔离 `abandoned/room_builder`，使 `go build ./...` 可过或 CI 显式排除。  
2. **P1**：`ApplyEffects` 优先 `SkillCastState.TargetEntity`（单目标技能），或与策划约定仅表选目标。  
3. **P2**：更新 `internal/battle/README.md`（`room_factory`、`system/runtime`、开房后需 `Update` 刷怪）。  
4. **P2**：为 `CastValidationSystem` / `DamageSystem` 增加包内单测（命中/必中/真伤）。  
5. **P3**：Query 按 `compID` 索引或文档明确组件数 ≤64。

---

## 9. 维护说明

架构变更请同步更新本文与 `docs/adr/`；可执行项见 [guides/optimization-recommendations.md](../guides/optimization-recommendations.md)。
