# 第五阶段：战斗逻辑扩展（实现说明）

对应 `record.md` 第五阶段清单与验收口径。以下内容描述**当前仓库已实现行为**及用法入口。

---

## 1. 暴击 / 格挡 / 闪避 / 命中

### 流程顺序（[`DamageSystem`](/internal/battle/system/system_damage.go)）

仅当 **`PendingDamage.Source ≠ 0`** 且 **`DamageType` 非 True** 时启用「战斗判定链」：

1. **命中 vs 闪避**：千分比骰 `0–999`。有效命中几率约为 `命中 + 750 − 闪避`，夹在 `[50, 995]`。
2. **格挡**：按 `BlockChancePermille` 判定；成功后从当前伤害减去 `BlockAmount`。
3. **暴击**：按 `CritChancePermille` 判定；成功后伤害乘以 `(1000 + CritDamagePermille) / 1000`（未配置暴伤时使用内置默认加成）。
4. **护甲减免**：沿用 [`MitigatedDamage`]（物/魔抗与 Buff 增量）。

**`Source == 0`**（如部分 DoT）：**不参与**命中/格挡/暴击，仅做护甲减免，避免 DoT 误吃暴击链。

### 属性字段（[`component.Attributes`](/internal/battle/component/types.go)）

| 字段 | 含义 |
|------|------|
| `HitPermille` / `DodgePermille` / `CritChancePermille` / `CritDamagePermille` | 千分比；**0 表示使用系统默认**（见 `system_damage.go` 常量） |
| `BlockChancePermille` / `BlockAmount` | 格挡 |

[`StatModifiers`](/internal/battle/component/stat_modifiers.go) 增加 `HitDeltaPermille`、`DodgeDeltaPermille`、`CritChanceDeltaPermille`、`CritDamageDeltaPermille`，供 Buff 叠加。

### 事件

未命中时派发 **`EventDamageMissed`**（[`ecs.Event`](/ecs/events.go) 含 `Entity`=受击者、`Attacker`=来源）。

---

## 2. 治疗系统（[`HealSystem`](/internal/battle/system/system_heal.go)）

- **直接治疗**：技能 `EffectHeal` 改为写入 [`MergePendingHeal`](/internal/battle/component/pending_merge.go)，由 **HealSystem** 增加 `Health.Current` 并派发 **`EventHealApplied`**。
- **HoT**：[`BuffSystem`](/internal/battle/system/system_buff.go) 的 `EffectHoT` 同样汇总到 **`PendingHeal`**（不再直接改血量），与直接治疗走同一管线，便于日志与扩展。

注册顺序：**DamageSystem → HealSystem → HealthSystem**，保证先结算伤害链再应用治疗再扣血。

---

## 3. 仇恨 / 威胁值（[`ThreatSystem`](/internal/battle/system/system_threat.go)）

- 订阅 **`EventDamageApplied`**（由 HealthSystem 在扣血后发出，含 **`Attacker`** 字段）。
- 对 **受害者** 实体上的 [`ThreatBook`](/internal/battle/component/threat.go) 累加威胁（当前规则：`威胁 += 伤害量 × 1`，常量可调）。
- **AI 选敌**：使用 [`component.ThreatTopSource`](/internal/battle/component/threat.go) 取威胁最高的攻击源实体。

需在会拉怪的目标上 **`AddComponent(&ThreatBook{})`**。

---

## 4. 战斗日志（[`CombatLogSystem`](/internal/battle/system/system_combat_log.go)）

- 订阅 **`EventDamageMissed`**、**`EventDamageApplied`**、**`EventHealApplied`**、**`EventDeath`**。
- 环形缓冲（默认最多 512 条），条目见 `CombatLogEntry`（含 `Kind` / `Target` / `Source` / `Value` / `Message`）。

---

## 5. 多阵营（>2）

[`Team.Side`](/internal/battle/component/team.go) 为 **`uint8`**，技能 [`CampEnemy`](/internal/battle/skill/skill_target_spec.go) 语义为 **「与施法者 Side 不同」**，**不限制只能两个阵营**。三方、四方对战在同一 World 内只需分配不同 `Side`。

---

## 6. 战斗重置与清理（[`ClearCombatEntities`](/internal/battle/system/reset_combat.go)）

- 遍历所有带 **`Team`** 的实体并 **`RemoveEntity`**。

战斗结束后应调用该函数（或自行销毁 World）以满足「临时状态清理」验收。

---

## 7. 伤害来源与事件载荷

[`PendingDamage`](/internal/battle/component/types.go) / [`ResolvedDamage`](/internal/battle/component/types.go) 增加 **`Source`**；[`MergePendingDamage`](/internal/battle/component/pending_merge.go) 最后一参为攻击者。

[`HealthSystem`](/internal/battle/system/system_health.go) 派发 **`EventDamageApplied`** 时填写 **`Attacker: rd.Source`**。

---

## 8. 系统注册顺序（[`AddCombatSystems`](/internal/battle/system/register.go)）

```
Buff → Cooldown → SkillChannel → SkillIntent
→ Damage → Heal → CombatLog → Threat → Health → Death
```

---

## 9. 验收对照（record.md）

| 验收项 | 实现说明 |
|--------|-----------|
| 暴击链（判定→倍率→伤害） | `Source≠0`：命中/格挡/暴击后再 `MitigatedDamage` |
| 仇恨影响 AI 选敌 | `ThreatBook` + `ThreatTopSource` |
| 战斗结束清理 | `ClearCombatEntities` |

---

## 10. 相关源码索引

| 路径 | 说明 |
|------|------|
| `internal/battle/system/system_damage.go` | 命中/格挡/暴击/减免 |
| `internal/battle/system/system_heal.go` | 治疗 |
| `internal/battle/system/system_threat.go` | 仇恨 |
| `internal/battle/system/system_combat_log.go` | 日志 |
| `internal/battle/system/reset_combat.go` | 清理 |
| `internal/battle/component/types.go` | PendingDamage.Source、战斗属性 |
| `ecs/events.go` | `Attacker` 字段与新事件类型 |
