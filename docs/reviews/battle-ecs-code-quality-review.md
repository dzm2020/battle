# battleDemo 代码质量与 ECS 规范评价

| 项 | 内容 |
|----|------|
| **文档类型** | 阶段性评审（Review） |
| **评价范围** | `internal/battle/`、`test/`、相关设计文档 |
| **排除范围** | `battle/ecs/` 底层实现（存储、Query 扫描、组件 ID 上限等） |
| **评价日期** | 2026-05-19 |
| **基线提交** | P0–P2 整改完成后的工作区快照 |
| **关联** | [STRUCTURE.md](../STRUCTURE.md) · [internal/battle/README.md](../../internal/battle/README.md) |

---

## 1. 执行摘要

### 1.1 一句话结论

**`internal/battle` 作为副本制帧同步战斗内核：在「广义 ECS」意义上规范、可演进；工程上已达「可协作的 Demo / 预生产骨架」，尚未达到「玩法完整 + 高覆盖率 + 生产运维」级别。**

### 1.2 综合评分（战斗层，满分 10）

| 维度 | 分值 | 说明 |
|------|------|------|
| 架构清晰度 | **9.0** | 分层、包边界、开房/帧管线可追踪 |
| ECS 用法（广义） | **8.0** | 数据/System/Resource/命令缓冲齐全；开战事件、规则 System 业务仍空 |
| 可维护性 | **8.5** | 目录语义化、`config.Load` 可失败返回、集成测与 API 对齐 |
| 可测试性 | **7.5** | 核心 System 已有单测；`system` 包覆盖率约 **44%**，随机伤害路径未测 |
| 可扩展性（玩法） | **8.0** | Installer/Spawner/表驱动/PVE·PVP 挂载点就绪 |
| 生产就绪 | **7.0** | 可预发；需补监控、玩法规则、更高测试覆盖 |
| **加权综合** | **8.0** | 适合继续投入；非「开箱即用上线」 |

### 1.3 是否符合 ECS？

| 口径 | 结论 |
|------|------|
| **广义 ECS**（组件 + System + Resource + 组合 + 帧更新） | **符合** |
| **狭义高性能 ECS**（SoA、Job、零分配 Query） | **不适用**；`battle/ecs` 为自研轻量 World，需另评 |
| 团队沟通建议 | 表述为 **「ECS 用法 + 自研 World」**，避免与 DOTS/EnTT 混称 |

---

## 2. 评价方法

1. **静态阅读**：包依赖、`component` 纯度、System 注册顺序、Resource 使用。  
2. **构建与测试**：`go build ./...`、`go test ./...`、子包 `-cover`。  
3. **对照设计文档**：`STRUCTURE.md`、skill/buff 设计说明与实现一致性。  
4. **不评价**：`ecs/` 包内 Query 性能、Entity API 细节。

---

## 3. 架构快照（当前）

```text
battleDemo/
├── ecs/                          # 通用 World（本评不评分）
├── internal/battle/
│   ├── component/                # comp_* + Register（含 Player）
│   ├── config/                   # Load / MustLoad
│   ├── resource/                 # TPS、SpawnQueue、RoomSpec、RoomPhase、SetPhase
│   ├── land/                     # 地图格子（≠ system/distance）
│   ├── room/                     # Create、Phase()、runLoop
│   ├── system/
│   │   ├── attrs/ transform/ action/ distance/ combatmath/
│   │   ├── room_bootstrap/       # Installer + Spawner（不 import 父包 system）
│   │   ├── entity_factory/ buff/ skill/ target_selector/target_filter/
│   │   ├── system_*.go           # 帧管线
│   │   ├── system_pve_rules.go   # PVE 扩展（占位）
│   │   ├── system_pvp_rules.go   # PVP 扩展（占位）
│   │   └── *_test.go
│   ├── pb/ event/ log/
│   └── README.md
├── test/                         # 集成测 + battle_config
└── cmd/configvalidate/
```

### 3.1 运行时序

```text
room.Create(RoomSpec)
  → 校验副本 / PVP 必须 Enemy
  → component.Register
  → BattleInitSystem + SpawnSystem
  → Resource：RoomPhase(Lobby→Fighting)、TPS、SpawnQueue、RoomSpec…
  → StartBattle → goroutine: world.Update(dt); Frame++

首帧 (TPS.Frame == 0)，BattleInitSystem 一次：
  land.CreateGridByID → AddResource(Grid)
  room_bootstrap.Bootstrap → Installer + Spawner
  FlushSpawnQueue（同帧刷怪）

后续帧：
  AddCoreCombatSystems 管线
  + PVERulesSystem 或 PVPRulesSystem（按副本类型）
```

### 3.2 依赖关系（无环）

```text
room → system
system.init → room_bootstrap.RegisterInstaller
room_bootstrap ↛ system
component ← attrs / transform / action / distance / combatmath
```

---

## 4. ECS 规范符合度（战斗层）

### 4.1 符合项（有代码依据）

| ECS 原则 | 实现 | 评价 |
|----------|------|------|
| **数据与行为分离** | 组件为 struct；修改走 `attrs`、`skill`、`buff`、各 System | 符合 |
| **System 驱动** | 伤害/Buff/技能/刷怪/结束均在 `Update` + Query | 符合 |
| **组合式实体** | `Team`+`Attributes`+`SkillSet`+`BuffList`+`Transform2D` 等拼装 | 符合 |
| **命令缓冲** | `DamageQueue`→`ResolvedDamage`、`PendingHeal`、`SpawnRequestQueue`、`SkillCastRequest` | 符合 |
| **World 单例资源** | `TPS`、`RoomSpec`、`land.Grid`、`RoomPhase` | 符合 |
| **单 World 隔离** | 每房间独立 `ecs.NewWorld` | 符合 |
| **显式系统顺序** | `AddCoreCombatSystems` 顺序与注释一致 | 符合 |
| **写请求、系统消费** | `RequestSkillCast`、`EnqueueSpawn` | 符合 |
| **按副本扩展** | `RegisterInstaller` / `RegisterSpawner` + PVE/PVP 规则 System | 符合 |
| **局内状态 Resource** | `RoomPhase` + `SetPhase`（Fighting / Settled / Closed） | 符合 |
| **玩家编队组件** | `Player` 已注册；`spawnPlayerUnits` 创建编队实体 | 符合（P0 已修） |

### 4.2 部分符合 / 折中

| 项 | 说明 |
|----|------|
| **组件「纯数据」** | `BuffControlState.Flags` 带 `HasStun()` 等查询方法 — 常见折中，可接受 |
| **`attrs` 职责** | 含阵营、最终属性（含 Buff 修正）；坐标已迁至 `transform` — 边界尚可 |
| **PVE/PVP 分化** | 挂载不同 System，但 `Update` 为空 — **结构分化完成，业务未分化** |

### 4.3 不符合或缺失（战斗层语义）

| 项 | 严重度 | 说明 |
|----|--------|------|
| **`KindBattleStart` 未派发** | 低 | 常量已定义，Bootstrap 后无统一开战事件 |
| **规则 System 无逻辑** | 低 | `PVERulesSystem` / `PVPRulesSystem` 仅占位 |
| **单测 `AddCombatSystems` 未含 PVE/PVP 规则 System** | 低 | 与正式开房路径略有差异，文档化即可 |

> 以上不构成「非 ECS」，属于 **玩法/事件闭环未完成**。

### 4.4 ECS 结论表

| 问题 | 答案 |
|------|------|
| 战斗层是否按 ECS 方式组织？ | **是** |
| 能否声称全栈高性能 ECS？ | **不能**（未评 `ecs/` 引擎） |
| 与早期版本相比？ | 循环依赖、utils 膨胀、test 编译、Player/Phase 等问题已解决 |

---

## 5. 代码质量（工程维度）

### 5.1 优势

1. **管线可读**：Spawn → Buff → 技能三阶段 → 伤/疗/血/死/结束，一目了然。  
2. **依赖倒置**：`room_bootstrap` 不 import `system`，Installer 由 `system.init` 注册。  
3. **属性单一数据源**：HP 仅 `Attributes`，经 `attrs` 读写。  
4. **首帧时序**：`FlushSpawnQueue` 与 Bootstrap 同帧，避免刷怪晚一帧。  
5. **子包语义**：`land`（格子）与 `distance`（实体距离）、`combatmath`、`target_filter` 命名清晰。  
6. **配表可失败**：`config.Load` 返回 `error`，`cmd/configvalidate` 可用于 CI。  
7. **测试门禁**：全仓库 `go test ./...` 通过。

### 5.2 不足与风险

| 类别 | 问题 | 风险 |
|------|------|------|
| **覆盖率** | `system` ~44%；buff_effect、entity_factory、room 无单测 | 回归依赖集成测 |
| **随机性** | `DamageSystem` 命中/暴击用 `rand`，无注入 | 单测难覆盖 miss/crit 分支 |
| **玩法** | PVE/PVP 规则 System 空实现 | 副本差异化仍靠配置表，非代码分支 |
| **事件** | 无 `KindBattleStart` | UI/录像订阅需自行约定帧号 |
| **粒度** | `action`、`distance` 各 1 个导出函数 | import 路径略多，非功能性问题 |

### 5.3 已关闭的历史问题（供对照）

| 历史问题 | 现状 |
|----------|------|
| `room_factory` ↔ `system` 循环依赖 | ✅ `room_bootstrap` |
| `utils` 职责混杂 | ✅ 拆为 attrs/transform/action/distance/combatmath |
| `test/` API 过期 | ✅ `room.Create`、`component.Register` |
| `Player` 未注册 | ✅ 已注册 + 单测 |
| `config.Load` panic | ✅ 返回 error |
| `target_fliter` 拼写 | ✅ `target_filter` |
| `RoomPhase` 不更新 | ✅ Settled / Closed |

---

## 6. 构建与测试（实测）

**执行时间**：2026-05-19  

```text
go build ./...     → 通过
go test ./...      → 通过
```

| 包 | 覆盖率（语句） | 测试文件 |
|----|----------------|----------|
| `internal/battle/config` | ~69% | `load_test.go` |
| `internal/battle/land` | ~71% | `grid_test.go` |
| `internal/battle/system` | ~44% | damage / cast_validation / buff / battle_end / spawn |
| `internal/battle/system/target_selector` | ~51% | `select_for_cast_test.go` |
| `test/` | 集成 | room / skill / target_selector 等 |

**建议 CI：**

```bash
go build ./...
go test ./...
go run ./cmd/configvalidate ./test/battle_config
```

---

## 7. 包级质量简评

| 包 | 质量 | 备注 |
|----|------|------|
| `component` | ★★★★★ | 纯数据 + 统一 `comp_` 前缀 |
| `resource` | ★★★★☆ | `SetPhase` 清晰；可补充 phase 转换单测 |
| `room` | ★★★★☆ | Create/Phase/Shutdown 完整；缺 room 包单测 |
| `land` | ★★★★☆ | 有单测；与 `distance` 文档需区分 |
| `attrs` | ★★★★☆ | 集中属性/阵营；`GetAttributeFinalValue` 略宽 |
| `room_bootstrap` | ★★★★★ | 解环 + Player 编队实体 |
| `entity_factory` | ★★★★☆ | 注释清楚；无单测 |
| `system`（帧管线） | ★★★★☆ | 核心路径有测；随机伤害未测 |
| `buff` / `skill` | ★★★★☆ | 表驱动；effect 子包无独立单测 |
| `target_selector` | ★★★★☆ | 有测；`SelectForCast` 覆盖点选 |
| `config` | ★★★★☆ | Load 可失败；表间引用校验仍弱 |

---

## 8. 两个核心问题的直接回答

### 8.1 代码质量怎么样？

**良好（8/10 档）。** 结构、依赖、ECS 用法、近期重构质量都达到「团队可在此基础上开发玩法」的水平。主要短板是 **测试覆盖率不均** 和 **PVE/PVP 规则 System 尚未承载业务**，而非架构性缺陷。

### 8.2 是否符合 ECS 规范？

- **就战斗层如何使用 ECS：符合广义 ECS 规范。**  
- **不满足**「组件绝对零方法」「全逻辑零分配」等极端口径。  
- **本评审不覆盖 `battle/ecs/`**，性能与扩展上限需单独评估。

---

## 9. 后续建议（按优先级）

| 优先级 | 建议 |
|--------|------|
| **P3** | `BattleInitSystem` 在 Bootstrap + `FlushSpawnQueue` 成功后派发 `event.KindBattleStart` |
| **P3** | `DamageSystem` 注入 `rand.Rand` 或固定种子，覆盖 miss/crit 单测 |
| **P2+** | 在 `PVERulesSystem` / `PVPRulesSystem` 实现波次、投降、同步等真实逻辑 |
| **P3** | 可选：合并 `action`+`distance` 为 `query`；或为 `room` 增加 phase 集成测 |
| **持续** | 将 `system` 包覆盖率提升至 **50%+**；`entity_factory` 补出生装配单测 |

---

## 10. 维护说明

- 架构或目录变更 → 同步更新本文、`docs/STRUCTURE.md`、`internal/battle/README.md`。  
- 替换 `battle/ecs/` → **另起**《ECS 引擎层评审》，勿与战斗层混谈。  
- **建议下次复审触发**：PVE/PVP 规则 System 有实质代码；或 `system` 覆盖率 ≥ 50%；或 `KindBattleStart` 落地。

---

## 附录 A：System 注册顺序

| # | System | 说明 |
|---|--------|------|
| — | `BattleInitSystem` | 仅 `room.Create` 挂载，首帧一次 |
| — | `SpawnSystem` | 常驻，消费刷怪队列 |
| 1 | `SpawnSystem` | （Installer 后同序重复注册于核心管线） |
| 2 | `BuffSystem` | |
| 3 | `CooldownSystem` | |
| 4 | `CastValidationSystem` | |
| 5 | `CastStateSystem` | |
| 6 | `DamageSystem` | |
| 7 | `HealSystem` | |
| 8 | `HealthSystem` | |
| 9 | `DeathSystem` | |
| 10 | `BattleEndSystem` | → `PhaseSettled` |
| 11 | `PVERulesSystem` 或 `PVPRulesSystem` | 按副本类型 |

---

## 附录 B：ECS 自检清单（评审用）

| # | 检查项 | 通过 |
|---|--------|------|
| 1 | 组件定义无业务帧循环 | ✅ |
| 2 | 帧逻辑在 System.Update | ✅ |
| 3 | 跨实体全局状态用 Resource | ✅ |
| 4 | 延迟写入用队列型组件 | ✅ |
| 5 | 开房与帧规则分包 | ✅ |
| 6 | 无包循环依赖 | ✅ |
| 7 | 副本可插 System | ✅ |
| 8 | 集成测可编译运行 | ✅ |
| 9 | 配表加载可失败返回 | ✅ |
| 10 | 生命周期事件基本闭环 | ⚠️ 缺 BattleStart |
