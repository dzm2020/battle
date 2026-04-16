# 第 2 天：战斗实体与属性模块（设计说明）

## 目标

- 建立 **基础属性 → 衍生战斗属性 → 局内临时状态** 三层数据模型。
- **计算与实体解耦**：实体只持有数据，`calc.Calculator` 负责公式，后续 Buff、配表、装备可替换实现。
- `attr` 包 **零业务依赖**（不引用 `entity` / `calc`），降低循环依赖风险。

## 目录与职责

| 路径 | 职责 |
|------|------|
| `internal/battle/attr/` | 纯数据结构：`Base`、`Derived`、`Runtime`。 |
| `internal/battle/calc/` | `Calculator` 接口与 `DefaultCalculator` 固定公式；可单测。 |
| `internal/battle/entity/` | `Entity` 聚合 ID、阵营、上述三类数据；`Recalculate` / `InitBattle`。 |

## 数据分层

1. **Base（基础）**  
   等级与四维：力量、敏捷、智力、体质。对应养成、配表入口，**不直接参与伤害公式读取**（避免客户端伪造攻防）。

2. **Derived（衍生 / 战斗属性）**  
   由 `Calculator.DerivedFromBase` 计算：生命/魔法上限、攻击、防御、暴击率/暴伤、物理减免等。  
   第 6 天伤害模块应只读 `Derived` + 运行时修正（Buff 可在未来作为 `Calculator` 的输入或单独 `Modifier` 链）。

3. **Runtime（临时）**  
   当前 HP/MP、护盾。可与 Redis 临时战斗数据对齐（学习计划第 1 阶段前置里提到的用途）。

## 默认公式（`DefaultCalculator`）

以下为教学用 **固定系数**，上线项目请改为策划表或脚本配置。

- `MaxHP = 100 + 20×Level + 15×VIT`
- `MaxMP = 50 + 5×Level + 8×INT`
- `ATK = 3×Level + 2×STR`
- `DEF = Level + VIT`
- `CritRate = clamp(0.05 + 0.001×AGI, 0, 0.5)`
- `CritDamage = clamp(1.5 + 0.01×STR, 1.25, 3.0)`
- `PhysMitigation = clamp(DEF / (DEF + 500), 0, 0.95)`（减伤系数，承伤比例约为 `1 - PhysMitigation`）

`ApplyMaxToRuntime`：当 `MaxHP/MaxMP` 变化时裁剪当前值，防止换装/降级后 **当前值超过上限**。

## 与后续天数的衔接

| 天数 | 衔接方式 |
|------|----------|
| 第 3 天 Tick | 可在 Tick 里对 `Runtime` 做持续效果；`Derived` 按帧重算或事件驱动重算。 |
| 第 4 天房间 | `BattleRoom` 持有 `[]*entity.Entity` 或实体 ID 映射，进房调用 `InitBattle`。 |
| 第 5–7 天技能/Buff | 校验读 `Derived`+`Runtime`；Buff 修改属性时可实现新 `Calculator` 或 **修饰器管道**（保持 `Entity` 结构稳定）。 |
| 第 8 天 ECS | 可将 `Base`/`Derived`/`Runtime` 迁为组件，`Calculator` 作为 System 依赖。 |

## 运行与测试

```bash
go test ./...
go run .
```

进程会先打印 `day02Demo` 再执行第 1 天 mock 网络流程。

## 扩展建议（不在本天代码中实现）

- 装备与 Buff：**不要**直接改 `Base` 公式路径，优先「修饰器累加后再进 `DerivedFromBase`」或独立 `CombatSnapshot`。
- 元素抗性：在 `Derived` 增加字段，公式集中在 `calc`。
- 服务器权威：任何客户端上报的「当前攻击/血量」仅作展示参考，**结算以服务端 `Entity` 为准**。
