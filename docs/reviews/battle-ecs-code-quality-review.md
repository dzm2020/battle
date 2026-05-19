# battleDemo 代码质量与 ECS 规范评价

| 项 | 说明 |
|----|------|
| **范围** | `internal/battle/`、`test/`、相关设计文档；**不评价** `battle/ecs/` 底层实现（存储、Query、位掩码等） |
| **评价日期** | 2026-05-19 |
| **相关文档** | [STRUCTURE.md](../STRUCTURE.md)、[skill-design.md](../skill-design.md)、[buff-and-skill-systems.md](../buff-and-skill-systems.md) |

---

## 1. 总评（一句话）

**作为战斗 Demo / 帧同步内核原型：架构清晰、广义 ECS 实践到位，近期重构（`resource`、`room_bootstrap`、同帧刷怪、属性/目标选取分层）显著改善可维护性。**  
**工程上仍是 Demo 级**：`test/` 集成测试未跟上 API、`internal/battle` 子包单测覆盖薄、PVE/PVP 战斗管线尚未分化。

**ECS 定位（仅就战斗层用法）**：**符合广义 ECS**（数据组件 + System 管线 + Resource + 命令缓冲 + 按副本扩展）；战斗层**不依赖**「高性能 ECS 引擎」能力，故不把 `ecs/` 实现细节纳入本评。

---

## 2. 评价范围说明

### 2.1 纳入

- 组件是否纯数据、System 职责与注册顺序  
- World Resource（`TPS`、`SpawnRequestQueue`、`RoomSpec`、`Grid` 等）  
- 开房生命周期（`room`）与帧内管线（`system`）分工  
- `room_bootstrap`（Installer / Spawner）、`entity_factory`、`attrs`、目标选取  
- 表驱动扩展、包边界与依赖方向  
- 构建/测试（`internal/battle` 子树）

### 2.2 不纳入（`battle/ecs/`）

以下仅作边界说明，**不评分、不建议改战斗层来迁就**：

- Archetype / SoA、Query 扫描复杂度、`compID` 上限  
- `AddComponent` 静默行为、Entity 生命周期 API 细节  

若未来要上大规模同屏，应单独评审或替换 `ecs/` 层，而非在 `internal/battle` 内打补丁。

---

## 3. 当前架构快照

```text
internal/battle/
├── component/              # 纯组件 + skill_cast_ops 等写入入口
├── config/                 # JSON 配表
├── resource/               # TPS、SpawnRequestQueue、RoomSpec、RoomPhase、Grid 访问
├── land/                   # 空间网格（含单测）
├── room/                   # Create、runLoop（动态 TPS）
├── system/
│   ├── attrs/              # Attributes 读写（System / entity_factory）
│   ├── room_bootstrap/     # Bootstrap = Installer + Spawner
│   ├── entity_factory/     # 出生装配
│   ├── buff/ skill/ target_selector/
│   ├── action/ distance/ utils/
│   └── system_*.go         # 帧管线
├── pb/ event/ log/
└── README.md
```

### 3.1 开房与战斗帧

```text
room.Create(spec)
  → component.Register
  → 仅挂 BattleInitSystem
  → Resource：RoomSpec、TPS(Frame=0)、SpawnRequestQueue…
  → runLoop：Update(dt) → Frame++

Frame == 0（首帧）：
  BattleInitSystem（仅执行一次）
    → land.CreateGridByID → AddResource(Grid)
    → room_bootstrap.Bootstrap
        → Installer：AddPVESystems / AddPVPSystems
        → Spawner：EnqueueSpawn
    → FlushSpawnQueue（同帧消费队列，不依赖新 System 进入当次遍历）

Frame >= 1：
  Spawn → Buff → 技能 CD/校验/阶段 → 伤/疗/血/死/结束
```

### 3.2 依赖关系（无包循环）

```text
room → system（BattleInitSystem）
system → room_bootstrap（Bootstrap、RegisterInstaller）
room_bootstrap ↛ system（Installer 由 system.init 注册）
```

---

## 4. 是否符合 ECS 规范？（战斗层视角）

### 4.1 符合或基本符合

| 要点 | 说明 |
|------|------|
| 组件以数据为主 | `Attributes`、`BuffList`、`SkillCastRequest` 等为 struct；运行时修改走 `system/attrs` |
| System 驱动 | 伤害/Buff/技能/刷怪/结束判定均在 `Update` + Query 中（Query 能力由 `ecs` 提供，此处不评实现） |
| 组合式单位 | `Team` + `Attributes` + `SkillSet` + `BuffList` 等拼装 |
| 命令缓冲 | `DamageQueue`→`ResolvedDamage`、`PendingHeal`、`SpawnRequestQueue` |
| World Resource | 网格、TPS、房间规格与实体解耦 |
| 单 World 隔离 | 每房间独立 World |
| 管线顺序显式 | `AddCoreCombatSystems` 注释与注册顺序一致 |
| 玩法写、系统读 | `RequestSkillCast` 入队；`EnqueueSpawn` 刷怪 |
| 按副本挂 System | `room_bootstrap.RegisterInstaller` + `AddPVESystems` / `AddPVPSystems` |
| 点选目标 | `SelectForCast` 在 `MaxCount==1` 时优先 `SkillCastState.TargetEntity` |
| 生命周期分工 | 开房在 `room` + 首帧 `BattleInitSystem`；帧规则在战斗 System |

### 4.2 仍需注意（战斗层）

| 问题 | 严重度 | 说明 |
|------|--------|------|
| `test/` 未同步 | 高 | 仍引用 `CreateRoom`、`room_bootstrap.Spec`、`component.Init` 等；`go test ./test/...` 无法编译 |
| PVE/PVP 管线相同 | 中 | `AddPVESystems` / `AddPVPSystems` 均仅 `AddCoreCombatSystems`，扩展点未用满 |
| System 单测少 | 中 | 除 `land`、`target_selector` 外多为 `?`（无测试文件） |
| `config.Load` | 中 | 失败路径偏硬（panic 等），表间引用校验有限 |
| 开战事件 | 低 | `KindBattleStart` 未在 Bootstrap 后统一派发 |
| 结束基线 | 低 | `BattleEndSystem` 在首次 `currSides>0` 时建基线，无独立开战 Resource |
| 命名笔误 | 低 | `skill_effect.Context.Word` 应为 `World` |

### 4.3 结论（战斗层 ECS 用法）

| 标准 | 结论 |
|------|------|
| **广义 ECS 用法** | **符合**，可作为副本制帧同步战斗内核继续演进 |
| **是否要求狭义高性能 ECS** | **不要求**；本评审已排除 `ecs/` 实现，不能据此声称「全栈高性能 ECS」 |

---

## 5. 自上一版评审以来的改进（战斗层）

| 项 | 状态 |
|----|------|
| `system/runtime` / `BattleContext` | 已移除；局内状态用 `resource` |
| `room_factory` 与 `system` 循环引用 | 已解决：`room_bootstrap` + Installer 注册表 |
| `component/attr_ops` | 已迁至 `system/attrs` |
| 施法目标双轨 | 已修复：`SelectForCast` |
| 首帧刷怪延迟 | 已修复：`FlushSpawnQueue` 同帧消费 |
| `BattleState` | 已删除；`BattleEnd` 用运行时存活阵营建基线 |
| 注释/文档陈旧 | 已批量对齐 `resource` / `room_bootstrap` |
| `AddCombatSystems` | 保留为单测别名 → `AddSystems` |

---

## 6. 代码质量（工程维度）

### 6.1 优点

- **管线清晰**：Spawn → Buff → 技能 → 伤害链 → 结束，职责可追踪。  
- **依赖倒置**：副本 System 由 `system.init` 注册到 `room_bootstrap`，子包不 import 父包。  
- **表驱动**：skill/buff effect、target_selector 可配置扩展。  
- **资源收敛**：`TPS` 同时驱动 `DeltaTime` 与 ticker。  
- **边界文档**：`entity_factory`、`room_bootstrap`、`attrs` 包注释说明职责。  
- **可编译子树**：`go build ./internal/battle/...` 通过。

### 6.2 主要问题

| 优先级 | 问题 |
|--------|------|
| **P0** | `test/` 编译失败，CI 若包含 `./test/...` 会红 |
| **P1** | 战斗 System（伤害/校验/Buff）缺包内表驱动单测 |
| **P2** | PVE/PVP 未挂差异化 System |
| **P2** | `room` 与 `RoomPhase` Resource 的同步关系未在代码中完全体现 |

### 6.3 构建与测试（2026-05-19）

```text
go build ./internal/battle/...              → 通过
go test  ./internal/battle/...              → 通过（land、target_selector 等）
go test  ./test/...                         → 编译失败（API 过期）
```

**建议 CI：**

```bash
go build ./internal/battle/...
go test ./internal/battle/...
# 修复 test/ 后：go test ./test/...
```

---

## 7. 主观评分（战斗层，不含 ecs/）

| 维度 | 分值 | 说明 |
|------|------|------|
| 架构清晰度 | **8.5** | resource + room_bootstrap + 首帧 Flush + 动态 TPS |
| ECS 用法（广义） | **7** | 数据/System/Resource/缓冲分工明确；缺开战/结束 Resource 闭环 |
| 可维护性 | **7** | 包边界好；test/ 与部分文档滞后 |
| 可扩展性（玩法） | **8** | Installer/Spawner/表驱动就绪 |
| 生产就绪 | **5.5** | 修 test/、补 System 单测、副本分化后可达 6+ |

---

## 8. 两个问题的直接回答

### 8.1 这份库（战斗层）代码怎么样？

**值得继续投入。** 开房流程、刷怪队列、战斗管线、副本扩展点已形成稳定骨架；近期重构解决了循环依赖、首帧刷怪、点选目标等实际问题。  
**上线前**需至少：修复 `test/`、为关键 System 补单测、按副本差异化 Installer。

### 8.2 是否符合 ECS 规范？

- **就战斗层如何用 ECS 而言：符合广义 ECS**，且比早期版本更规范（纯数据组件、`attrs`、命令缓冲、System 管线、Resource）。  
- **本评审不覆盖 `battle/ecs/`**，故不对「Query 性能、组件上限」等下结论。  
- 若团队口头说「ECS」，应明确是 **「ECS 用法 + 自研轻量 World」**，而非 DOTS/EnTT 一类引擎。

---

## 9. 建议的下一步（按优先级）

1. **P0**：修复 `test/`——`room.Create(&resource.RoomSpec{...})`、`component.Register`、`system.AddCombatSystems`。  
2. **P1**：为 `CastValidationSystem`、`DamageSystem`、`BuffSystem` 增加 `internal/battle/system/*_test.go`。  
3. **P2**：在 `AddPVESystems` / `AddPVPSystems` 中挂副本专用 System。  
4. **P2**（可选）：Bootstrap 成功后派发 `event.KindBattleStart`。  
5. **P3**：修正 `skill_effect.Context` 命名；评估 `RoomPhase` 与 `room` 状态机一致性。

---

## 10. 维护说明

- 架构变更请同步更新本文与 `docs/adr/`。  
- 若替换或重写 `battle/ecs/`，应**另起一篇**引擎层评审，勿与本战斗层评价混为一谈。  
- 下次评审触发条件建议：`test/` 全绿，或 PVE/PVP System 分化落地。
