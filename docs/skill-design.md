# 技能系统设计说明

本文描述第四阶段「可配置技能框架」的实现约定：静态模板 [skill.SkillConfig]、运行时组件、系统顺序、与 Buff/伤害链的衔接，以及 JSON/YAML 数据驱动方式。**实现以仓库源码为准**。

---

## 1. 目标与边界

- **数据驱动**：技能消耗、冷却、目标类型、效果列表完全由 [skill.CatalogConfig] 中的模板描述；可通过 [skill.LoadCatalogConfigFromJSON] / [skill.LoadCatalogConfigFromYAML] 加载，无需改代码即可增删技能（前提是逻辑已支持的 Effect 类型）。
- **资源**：支持 **无 / 法力 / 怒气 / 能量** 四种消耗语义，对应 [component.SkillUser] 上的三个资源槽（无消耗时 **Cost 必须为 0**）。
- **冷却**：按 **帧** 计数，存在 [component.SkillUser].CooldownRemaining；由 [CooldownSystem] 每帧递减。
- **目标**：自身、单体敌方、**全体异阵营**（需 [component.Team] + [component.Health]）。
- **瞬发 / 吟唱**：`CastFrames==0` 为瞬发；`CastFrames>0` 消耗资源当帧扣除，写入 [component.SkillCastState]，经过若干帧后结算效果并进入冷却。
- **与 Buff**：效果条目 `EffectApplyBuff` 调用 [buff.ApplyBuff]，共享同一份 [buff.DefinitionConfig]。
- **与伤害**：效果条目 `EffectDamage` 调用 [component.MergePendingDamage]，走既有 **Damage → Health → Death**。

不在本文范围：**打断吟唱**、技能队列、弹道飞行、精确范围几何；可用玩法层不写 CastIntent 或扩展组件实现。

---

## 2. 静态模板：`skill.SkillConfig` / `skill.EffectConfig`

核心类型见 `internal/battle/skill/skill_config.go`，字段均有中文注释。

| 字段 | 含义摘要 |
|------|-----------|
| `ID` | 全局技能 ID，与施放意图、授予列表一致 |
| `Resource` / `Cost` | 资源类型与数值 |
| `CooldownFrames` | 冷却帧长；从 **效果结算完毕** 的当帧写入冷却表 |
| `Target` | `Self` / `SingleEnemy` / `AllEnemySides` |
| `CastFrames` | 0 瞬发；>0 吟唱帧数 |
| `Effects` | 有序效果列表 |

**Effect 种类**：

| Kind | 行为 |
|------|------|
| `EffectDamage` | 对每个解析出的目标 `MergePendingDamage` |
| `EffectHeal` | 对每个目标增加 `Health.Current`（封顶 Max） |
| `EffectApplyBuff` | 对每个目标 `buff.ApplyBuff(..., BuffDefID)` |

---

## 3. 运行时组件（`internal/battle/component`）

| 组件 | 职责 |
|------|------|
| `Team` | `Side` 阵营；群体「敌方」筛选依赖施法者与目标的 Side 不等 |
| `SkillUser` | `Mana/Rage/Energy`、`GrantedSkillIDs`、`CooldownRemaining` |
| `CastIntent` | 外部写入：`SkillID` + `Target`（主目标）；由 [SkillSystem] 消费并移除 |
| `SkillCastState` | 吟唱中：`SkillID`、`PrimaryTarget`、`FramesLeft` |

施法者实体建议同时具备 `SkillUser`（与资源、冷却相关）及必要时 `Team`（阵营技能）。

---

## 4. 系统顺序与单帧语义

[AddCombatSystems] 注册顺序为：

```
BuffSystem → CooldownSystem → SkillSystem → DamageSystem → HealthSystem → DeathSystem
```

- **BuffSystem**：DoT、属性修正等先于技能，本帧已有的 `PendingDamage` / `StatModifiers` 仍有效。
- **CooldownSystem**：对所有 `SkillUser` 的冷却表 **先** 递减，便于「上一帧结束冷却」的技能在本帧 **[SkillSystem]** 中可用。
- **SkillSystem**：  
  1. **advanceChannels**：吟唱 `FramesLeft` 递减；归零当帧 **结算效果 + 写入冷却**，移除 `SkillCastState`。  
  2. **processIntents**：处理 `CastIntent`（眩晕由 [action.CanAct] 拦截）；校验授予、冷却、目标、资源；瞬发直接 `ExecuteEffects`；吟唱则写入 `SkillCastState`（**资源在发起吟唱当帧扣除**）。
- **DamageSystem**：合并技能与本帧其它来源的 `PendingDamage`。

---

## 5. API 使用流程（玩法层）

1. `component.RegisterCombatTypesWorld(w)`（含技能与 Team 等新组件）。
2. 构造 `buff.DefinitionConfig`、`skill.CatalogConfig`，加载 JSON/YAML。
3. `system.AddCombatSystems(w, buffConfig, skillConfig)`。
4. 给单位添加 `Team`、`SkillUser`（填充资源与 `GrantedSkillIDs`）、`Health` 等。
5. 施放：`AddComponent(caster, &CastIntent{SkillID, Target})`，再 `world.Update(dt)`。
6. 下一帧 Intent 应已被移除；吟唱则等待 `SkillCastState` 清空后再发起新 Intent（建议玩法层禁用重叠）。

---

## 6. JSON / YAML 示例

以下示例与测试用例一致，可作为配置起点（**枚举值为数值**，与 Go 中 `iota` 顺序一致）。

### 6.1 JSON 数组（单体法术）

```json
[
  {
    "id": 10,
    "resource": 1,
    "cost": 30,
    "cooldownFrames": 2,
    "target": 1,
    "castFrames": 0,
    "effects": [
      { "kind": 0, "amount": 40, "damageType": 1 }
    ]
  }
]
```

含义：`resource:1` 为法力；`target:1` 为单体敌方；`kind:0` 为伤害；`damageType:1` 为魔法（见 `component.DamageType`）。

### 6.2 YAML

结构与 JSON 相同，使用 [skill.LoadCatalogConfigFromYAML]，文件可为：

```yaml
- id: 11
  resource: 0
  cost: 0
  cooldownFrames: 0
  target: 2
  castFrames: 0
  effects:
    - kind: 0
      amount: 15
      damageType: 0
```

`target:2` 表示 `TargetAllEnemySides`（全场异阵营，需施法者有 `Team`）。

---

## 7. 验收对照（record 第四阶段）

| 条目 | 实现要点 |
|------|-----------|
| 配置结构 | [skill.SkillConfig] + [skill.EffectConfig]，JSON/YAML 加载 |
| SkillUser | 资源、授予列表、冷却 map |
| SkillSystem | CastIntent + SkillCastState + ExecuteEffects |
| 冷却 | [CooldownSystem] |
| 资源 | Mana/Rage/Energy + ResourceType |
| Buff 联动 | EffectApplyBuff |
| 瞬发 / 吟唱 | CastFrames 0 / >0 |
| 群体 | TargetAllEnemySides + Team |

---

## 8. 相关源码索引

| 路径 | 说明 |
|------|------|
| `internal/battle/skill/skill_config.go` | 模板类型与枚举 |
| `internal/battle/skill/catalog_config.go` | 模板配置表 |
| `internal/battle/skill/skill_config_json.go` / `skill_config_yaml.go` | 加载 |
| `internal/battle/skill/resolve.go` | 目标解析与效果执行 |
| `internal/battle/component/skill.go` / `team.go` | 运行时组件 |
| `internal/battle/system/system_skill.go` | 施法主逻辑 |
| `internal/battle/system/system_cooldown.go` | 冷却递减 |
| `internal/battle/system/register.go` | `AddCombatSystems` |
