# Buff（状态效果）系统设计说明

本文描述 `battle/internal/battle` 下 Buff/DeBuff 的运行时存储、静态配置与系统调度。**实现以仓库代码为准**；若与 `record.md` 措辞有出入，以本节 + 源码为准。

---

## 1. 设计目标

- **单实体多条 Buff**：一个 ECS 组件 `BuffList` 内持有一条切片，可挂载 **20+** 条互不影响的 Buff 实例（每条有独立 `DefID`、`Stacks`、`FramesLeft`、`TickCountdown`）。
- **静态与运行时分离**：效果数值、持续时间、叠层策略等来自 **`buff.Descriptor`**（可通过 JSON 加载）；运行时只保存实例状态 + 指向定义的 `DefID`。
- **与战斗管线衔接**：DoT 写入 `PendingDamage`，经 **Damage → Health → Death** 结算；Buff 系统在 Damage 之前运行。
- **属性与控制外置**：每帧根据 Buff 聚合出 **`StatModifiers`**（可被 `DamageSystem` 用于有效护甲/魔抗）与 **`ControlState`**（眩晕/沉默等，供 `action` 包校验）。

---

## 2. 运行时「缓冲区」：`BuffList` / `BuffInstance`

### 2.1 `component.BuffList`

| 字段 | 含义 |
|------|------|
| `Buffs []BuffInstance` | 当前实体上所有活跃的 Buff 实例列表（逻辑上的「缓冲表」） |

- 每个实体 **最多一个** `BuffList` 组件；多条 Buff 存在 **切片扩容** 中，而不是多个同名 ECS 组件类型。
- 优势：组件类型数量少（利于当前 ECS 的 `uint8` 组件 ID 空间），且便于一次遍历汇总。

### 2.2 `component.BuffInstance`

| 字段 | 含义 |
|------|------|
| `DefID` | 在 `buff.DefinitionRegistry` 中查找 `Descriptor` 的主键 |
| `Stacks` | 叠层数（参与数值类效果乘算） |
| `FramesLeft` | 剩余持续帧；`≥0` 时按帧递减；**`<0` 表示不限时**（不由时间轴移除） |
| `TickCountdown` | DoT/HoT 距离下一次结算的帧计数（见 §5） |

未知 `DefID`（注册表缺定义）的实例在本帧会被 **跳过**（不累计属性、不触发 DoT），并在存活列表中 **丢弃**（实例被过滤掉）。

---

## 3. 静态配置：`buff.Descriptor` / `EffectDef`

### 3.1 `Descriptor`（一条 Buff「模板」）

| 字段 | 含义 |
|------|------|
| `ID` | 全局唯一 Buff 定义 ID |
| `MaxStacks` | 叠层上限（至少按 1 处理） |
| `Policy` | 同名再次施加时的策略（见 §4） |
| `DurationFrames` | 持续帧数；`-1` 为不限时直到驱散/逻辑移除 |
| `Effects` | 效果列表，一条 Descriptor 可挂多种效果 |

### 3.2 `EffectKind` 与 `EffectDef`

| Kind | 作用 |
|------|------|
| `EffectStatMod` | 按层数叠加 **护甲/魔抗/物理强度** 增量到 `StatModifiers` |
| `EffectDoT` | 按 `TickIntervalFrames` 间隔向实体写入 **合并后的** `PendingDamage` |
| `EffectHoT` | 按间隔直接增加 `Health.Current`（不超过 `Max`） |
| `EffectControl` | 将 `control.Flags` **按位或** 并入 `ControlState`（不按层数倍增） |

数值类效果：`增量每栈 × Stacks`。DoT/HoT 单次结算量：`每跳数值 × Stacks`。

### 3.3 配置入口

- 代码：`DefinitionRegistry.Register(Descriptor)`。
- 数据：`buff.LoadDescriptorsFromJSON` 解析 JSON 数组并批量 `Register`（字段见 `Descriptor` / `EffectDef` 上的 `json` 标签）。

---

## 4. 叠层策略：`StackPolicy`

由 `buff.ApplyBuff` 实现：

| 策略 | 行为（当前实现） |
|------|-------------------|
| **Independent** | 每次施加追加 **一条新实例**，持续时间各自独立倒计时 |
| **Refresh** | 若已存在同名实例：**仅把 `FramesLeft` 重置为配置值**，**不改变 `Stacks`**；不存在则新建 1 层 |
| **Merge** | 若已存在同名实例：**`Stacks++`**（封顶 `MaxStacks`），并 **重置持续时间**；不存在则新建 |

施加新实例时：`FramesLeft = Descriptor.DurationFrames`；`TickCountdown` 由第一个 DoT/HoT 的 `TickIntervalFrames` 推导（见 `apply.go` 中 `findFirstInterval`）。

---

## 5. 每帧时间轴：持续时间与 DoT/HoT

对 **每个** `BuffInstance`，`BuffSystem` 顺序为：

1. **聚合**本帧仍有效的属性/控制（`accumulateStatic`）。
2. **DoT/HoT 节拍**：若 Descriptor 中 **第一个** DoT/HoT 定义了间隔，则取该间隔；`TickCountdown` 每帧 `-1`，当 `<0` 时结算本跳所有 DoT/HoT 效果，并将 `TickCountdown` 重置为 `interval - 1`。
3. **持续时间**：若 `FramesLeft ≥ 0`，则本帧末 `FramesLeft--`；若 `≤0` 则本实例移除。`FramesLeft < 0` 的不做时间递减。

**要点**：多条 `EffectDoT` / `EffectHoT` 共用 **同一套** 实例级 `TickCountdown`；间隔以 **遍历 Effects 时首次遇到的** DoT/HoT 为准。

---

## 6. 与战斗系统、合并伤害

### 6.1 系统顺序（`system.AddCombatSystems`）

```
BuffSystem → DamageSystem → HealthSystem → DeathSystem
```

因此 **本帧 DoT** 写入的 `PendingDamage` 会在同一 `World.Update` 内进入减免与扣血。

### 6.2 `StatModifiers` 与 `DamageSystem`

- `BuffSystem` 每帧将 Buff 汇总写入 `StatModifiers`（全量覆盖式重算：先清零再累积）。
- `DamageSystem` 读取 `Attributes` + `StatModifiers` 得到有效物甲/魔抗后再走 `MitigatedDamage`。

### 6.3 `MergePendingDamage`

同一实体上 `PendingDamage` **仅有一份** ECS 组件时，多次写入会 **合并 `Amount`**；类型不一致时退化为 **`DamageTrue` 合并**（避免静默丢伤害）。DoT 与普攻同帧共存时依赖此合并。

### 6.4 `ControlState` 与行动

- `ControlState.Flags` 每帧由 Buff 控制类效果重算。
- `action.CanAct` 当前以 **眩晕** 为硬拦截；沉默/定身可在技能/位移分流中读取 `Flags` 扩展。

---

## 7. 初始化约定

1. `component.RegisterCombatTypesWorld(w)`（或等价地向 `Registry` 注册战斗相关组件类型）。
2. 构造 `buff.DefinitionRegistry`，注册或通过 JSON 加载全部 `Descriptor`。
3. `system.AddCombatSystems(w, registry)`。
4. 对实体 `ApplyBuff` / 添加 `Health`、`Attributes` 等；每帧调用 `world.Update(dt)`。

---

## 8. 已知限制与可选扩展

- **驱散 / 免疫 / 护盾**：未在 Buff 内核实现，可在施加前过滤或在独立系统中移除 `BuffList` 子集。
- **同类多条 DoT 不同间隔**：当前单实例仅一套 `TickCountdown`；若需要，可拆成多条 Independent 实例或扩展实例结构。
- **Refresh 策略**：当前实现为「只刷新时间、不涨层」；若策划需要「刷新且 +1 层」，需在 `ApplyBuff` 中单独调整。
- **组件 ID 上限**：自建 ECS 使用 `uint8` 组件类型计数；Buff 的大量种类应放在 **`BuffList` 切片 + DefID**，勿为每种 Buff 新建一种 ECS 组件类型。

---

## 9. 相关源码路径（索引）

| 模块 | 路径 |
|------|------|
| Buff 切片组件 | `internal/battle/component/buff_list.go` |
| Stat / 控制外挂组件 | `internal/battle/component/stat_modifiers.go`, `control_state.go` |
| 合并待结算伤害 | `internal/battle/component/pending_merge.go` |
| 定义 / 注册表 / Apply / JSON | `internal/battle/buff/` |
| BuffSystem | `internal/battle/system/system_buff.go` |
| 战斗系统注册顺序 | `internal/battle/system/register.go` |
| 控制位枚举 | `internal/battle/control/flags.go` |
| 能否行动（眩晕） | `internal/battle/action/act.go` |



DoT (Damage over Time):随时间持续伤害
HoT (Heal over Time):随时间持续治疗