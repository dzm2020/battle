# 仓库路径与文档结构（建议）

本文约定 **文档放哪、代码边界怎么划**，便于协作与 CI。若与现状不一致，以「建议目标」为准逐步迁移。

---

## 1. 顶层布局（建议）

```text
battleDemo/
├── cmd/                    # 可执行入口（configvalidate、server 等）
├── docs/                   # 所有人可读：架构、评审、ADR、接入说明
│   ├── STRUCTURE.md        # 本文件：路径与结构约定
│   ├── adr/                # 架构决策记录（一条决策一篇，长期不变更语义）
│   ├── reviews/            # 阶段性代码质量/安全评审（可按日期或版本命名）
│   ├── guides/             # 操作手册：如何配表、如何开房间、如何跑测
│   └── examples/           # 若有：示例配置片段（与 test 数据区分）
├── internal/               # 应用私有代码，不对外作为 library API 承诺
│   └── battle/             # 战斗内核（见 internal/battle/README.md）
├── ecs/                    # 通用 ECS（若希望更强边界，可未来改为 pkg/ecs）
├── test/                   # 黑盒/集成测试与 battle_config 样例（易整包编译失败时可拆包）
├── go.mod
└── todo.md                 # 个人任务可保留；团队路线图建议进 docs/guides/roadmap.md
```

---

## 2. 文档三类（不要混用）

| 类型 | 目录 | 命名示例 | 用途 |
|------|------|----------|------|
| **ADR** | `docs/adr/` | `001-room-grid-land-vs-wrapper.md` | 不可逆或高成本决策：为何用 `*land.Grid`、为何不用全局单例等 |
| **评审 / 质量报告** | `docs/reviews/` | `internal-battle-2026-05.md` | 某时间点的扫描结论、债清单、优先级 |
| **指南** | `docs/guides/` | `optimization-recommendations.md`, `config-tables.md` | 可执行优化清单、配表与流程说明 |

**建议迁移（可选）：**

- 将现有 `docs/internal-code-quality-review.md` **移动或复制**为：  
  `docs/reviews/internal-battle-code-quality.md`  
  原路径保留一个 **三行以内** 的 stub 指向新路径，避免外部死链。

---

## 3. ADR 约定（`docs/adr/`）

- **文件名**：`NNN-英文短横线主题.md`，`NNN` 三位序号，按合并时间递增。  
- **文内固定小节**：背景 → 决策 → 备选方案 → 后果（正面/负面）→ 状态（提议/已采纳/已废弃）。  
- **何时写**：涉及包边界、并发模型、持久化、网络协议、**破坏性改 public API** 时必写。  
- **索引**：维护 `docs/adr/README.md` 中的表格（序号、标题、状态、合并日期）。

---

## 4. `internal/battle/` 代码结构（建议目标）

```text
internal/battle/
├── README.md               # 子包地图 + 指向 docs/ 的链接（入口）
├── config/                 # 表结构与 Load；长期可加 loader 接口与校验子包
├── ecs/                    # 若仅 battle 使用可保持现状；与根 ecs/ 二选一见 ADR
├── component/              # 纯组件定义与 Register
├── entity_factory/         # 从配置/PB 造实体；不写格子策略
├── land/                   # 纯空间索引；不依赖 room
├── room/                   # 房间生命周期、阶段、与 World 同线程的 API
├── room_builder/           # 进房装配、副本类型分发；依赖 config + room + entity_factory
├── system/                 # ECS Systems；依赖顺序在 register 中显式化
├── skill/ / buff/ / target_selector/ …
├── utils/                  # 建议逐步瘦身：与 component 类型强耦合的进同包或 component 旁 helper
├── pb/                     # 与协议/存档对齐的 DTO
├── log/
└── event/
```

**边界原则（写入团队规范即可）：**

- **`land`**：不知道 `Room`、`pb.Player`。  
- **`room`**：不 `import entity_factory`（落点用 `PlaceOnGrid` 等薄 API；造怪在 `room_builder`）。  
- **`room_builder`**：唯一组装「表 + factory + grid」的场所之一（另一条可以是未来的 `spawn` 子包）。  
- **`utils`**：仅放与 ECS 无环或弱耦合的纯函数；否则易与 `component` 类型漂移。

---

## 5. 测试目录策略

| 位置 | 用途 |
|------|------|
| `internal/battle/<pkg>/*_test.go` | 包内单元测试、表驱动、快速反馈 |
| `test/` | 跨包集成、需 `battle_config` 全量的场景 |
| `test/battle_config/` | 最小可加载表集；与 `docs/examples/` 不重复时可只保留一处 |

**建议**：为 `test` 拆包（例如 `test/integration` 为单独 module 或单独 package）避免单个坏文件导致整目录无法 `go test`。

---

## 6. 执行顺序（落地 checklist）

1. 新建 `docs/adr/README.md`、`docs/guides/`（空或占位）、`docs/reviews/`。  
2. 将质量评审文移入 `docs/reviews/` 并在原路径留跳转（可选）。  
3. 新增 `internal/battle/README.md` 子包表。  
4. 第一个 ADR：记录 **`ecs` 根目录 vs `internal`** 或 **`room.Grid` vs `*land.Grid`** 的最终选择。  
5. CI：`go test ./internal/battle/...` 与 `go test ./ecs/...` 分步，修复后再加全仓库。

---

维护：结构变更时先改 **本文件** 与 **相关 ADR**，再动目录。
