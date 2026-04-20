# 技能系统设计说明

本文描述第四阶段「可配置技能框架」的实现约定：静态模板 [skill.SkillConfig]、运行时组件、系统顺序、与 Buff/伤害链的衔接，以及 JSON/YAML 数据驱动方式。**实现以仓库源码为准**。目标选取的设计维度（范围 / 阵营 / 过滤规则等）可与仓库根目录 **skill_record.md** 对照阅读。

---

## 1. 目标与边界

- **数据驱动**：技能消耗、冷却、目标类型、效果列表完全由 [skill.CatalogConfig] 中的模板描述；可通过 [skill.LoadCatalogConfigFromJSON] / [skill.LoadCatalogConfigFromYAML] 加载，无需改代码即可增删技能（前提是逻辑已支持的 Effect 类型）。
- **资源**：支持 **无 / 法力 / 怒气 / 能量** 四种消耗语义，对应 [component.SkillUser] 上的三个资源槽（无消耗时 **Cost 必须为 0**）。
- **冷却**：按 **帧** 计数，存在 [component.SkillUser].CooldownRemaining；由 [CooldownSystem] 每帧递减。
- **目标**：按 [skill_record.md] 三维度在配置中写明 **`scope`**、**`camp`**、可选 **`pickRule`**，并配合 `maxTargets`、`aoeRadius`、Buff 条件。见 `internal/battle/skill/skill_target_spec.go` 与 [ResolveTargets]。
- **refinement**：`maxTargets` 截断；`aoeRadius` 对 Cone/Circle/Line、Multi、全屏 与 `Transform2D` 做距离球过滤（为 0 时不做距离裁剪）；`requireBuffDefId` / `forbidBuffDefId` 在排序截断前过滤。
- **瞬发 / 吟唱**：`CastFrames==0` 为瞬发；`CastFrames>0` 消耗资源当帧扣除，写入 [component.SkillCastState]，经过若干帧后结算效果并进入冷却。
- **与 Buff**：效果条目 `EffectApplyBuff` 调用 [buff.ApplyBuff]，共享同一份 [buff.DefinitionConfig]。
- **与伤害**：效果条目 `EffectDamage` 调用 [component.MergePendingDamage]，走既有 **Damage → Health → Death**。

不在本文范围：**打断吟唱**、技能队列、弹道飞行、精确扇形几何（当前无朝向；Cone 可与半径一起做近似）；可用玩法层不写 CastIntent 或扩展组件实现。

---

## 2. 静态模板：`skill.SkillConfig` / `skill.EffectConfig`

核心类型见 `internal/battle/skill/skill_config.go`。

| 字段 | 含义摘要 |
|------|-----------|
| `ID` | 全局技能 ID |
| `Resource` / `Cost` | 资源类型与数值 |
| `CooldownFrames` | 冷却帧长；从 **效果结算完毕** 的当帧写入冷却表 |
| **`scope`**（必填） | [TargetScope]，见下表；JSON 值为 **0 表示无效** |
| **`camp`**（必填） | [CampRelation]，见下表；敌方为 **0** |
| **`pickRule`** | [PickRule]，可选 |
| **`aoeRadius`** | 球半径；与 Multi/几何范围等组合 |
| **`campSide`** | 仅 `camp` = SpecificSide |
| `maxTargets` | 排序后截断；随机默认次数；链式总人数上限 |
| `chainJumps` | 链式额外跳跃 |
| `requireBuffDefId` / `forbidBuffDefId` | Buff 模板过滤 |
| `CastFrames` | 0 瞬发；>0 吟唱帧数 |
| `Effects` | 有序效果列表 |

**TargetScope（`scope`）**

| 值 | 常量 | 含义 |
|----|------|------|
| 0 | — | **无效**，施放解析失败 |
| 1 | `TargetScopeSelf` | 仅自身 |
| 2 | `TargetScopeSingle` | 单体：主目标 + `camp` |
| 3–5 | Cone / Circle / LineRect | 先按 `camp` 取候选，再按 `aoeRadius` + `Transform2D` 裁剪 |
| 6 | `TargetScopeMulti` | 全场符合 `camp` 的多目标 |
| 7 | `TargetScopeFullScreen` | 全场实体再按 `camp` 缩小 |
| 8 | `TargetScopeChain` | 链式 |
| 9 | `TargetScopeRandom` | 随机子集 |

**CampRelation（`camp`）**：`0` 敌方 · `1` 友方含自身 · `2` 友方不含自身 · `3` 不限阵营 · `4` 指定 Side（配合 `campSide`）。

**PickRule（`pickRule`）**：`0` 无 · `1` 最近 · `2` 最远 · `3` 当前血量升序 · `4` 血量百分比升序 · `5` 攻击力最高。

**组合示例**：「群体治疗、友方、血量百分比最低 5 个」→ `scope=6`、`camp=1`、`pickRule=4`、`maxTargets=5`。

**Effect 种类**：同前（伤害 / 治疗 / 挂 Buff）。

---

## 3. 目标解析与施法校验

- [skill.ResolveTargets]：**候选（scope×camp×空间）→ Buff 过滤 → PickRule 排序 → maxTargets**。
- [skill.ValidCastTargets]：意图阶段；`scope==0` 的配置恒失败。

---

## 4. 运行时组件（`internal/battle/component`）

| 组件 | 职责 |
|------|------|
| `Team` | `Side` 阵营 |
| `SkillUser` | 资源、授予列表、冷却 map |
| `CastIntent` | `SkillID` + 主目标；由 [SkillIntentSystem] 消费 |
| `SkillCastState` | 吟唱状态 |
| `Transform2D` | 可选坐标；用于 `aoeRadius` 与最近/最远 |

---

## 5. 系统顺序与单帧语义

```
BuffSystem → CooldownSystem → SkillChannelSystem → SkillIntentSystem → DamageSystem → HealthSystem → DeathSystem
```

（语义同前：先吟唱结算，再处理意图。）

---

## 6. API 使用流程（玩法层）

（与前一版相同：注册组件 → 加载配置 → AddCombatSystems → 添加单位 → 写 CastIntent → Update。）

---

## 7. JSON 示例

### 单体法术（敌方单体）

```json
[
  {
    "id": 10,
    "resource": 1,
    "cost": 30,
    "cooldownFrames": 2,
    "scope": 2,
    "camp": 0,
    "castFrames": 0,
    "effects": [{ "kind": 0, "amount": 40, "damageType": 1 }]
  }
]
```

### 友方群体治疗（百分比最低 2 人）

```json
{
  "id": 11,
  "resource": 0,
  "cost": 0,
  "cooldownFrames": 5,
  "scope": 6,
  "camp": 1,
  "pickRule": 4,
  "maxTargets": 2,
  "castFrames": 0,
  "effects": [{ "kind": 1, "amount": 20 }]
}
```

### 全场敌方 AOE（示例）

```yaml
- id: 12
  resource: 0
  cost: 0
  cooldownFrames: 0
  scope: 6
  camp: 0
  castFrames: 0
  effects:
    - kind: 0
      amount: 15
      damageType: 0
```

---

## 8. 验收对照

| 条目 | 实现要点 |
|------|-----------|
| 配置 | [SkillConfig]：`scope`/`camp`/`pickRule` + ResolveTargets |
| 管线 | SkillChannelSystem + SkillIntentSystem |

---

## 9. 相关源码索引

| 路径 | 说明 |
|------|------|
| `internal/battle/skill/skill_config.go` | 模板 |
| `internal/battle/skill/skill_target_spec.go` | TargetScope / CampRelation / PickRule |
| `internal/battle/skill/resolve.go` | 解析与效果执行 |
